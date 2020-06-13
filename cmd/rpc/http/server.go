package main

import (
	"context"
	"flag"
	"fmt"
	"go-proj/app/rpc/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/gpprof"
	"github.com/daheige/thinkgo/monitor"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	config "go-proj/conf"
	"go-proj/pb"

	"go-proj/app/rpc/service"

	"github.com/daheige/thinkgo/logger"

	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
)

var port int
var logDir string
var configDir string
var wait time.Duration // 平滑重启的等待时间1s or 1m

func init() {
	flag.IntVar(&port, "port", 1339, "grpc http gw port")
	flag.StringVar(&logDir, "log_dir", "./logs", "log dir")
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful_timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	// 日志文件设置
	logger.SetLogDir(logDir)
	logger.SetLogFile("go-grpc.log")
	logger.MaxSize(500)

	// 由于app/extensions/logger基于thinkgo/logger又包装了一层，所以这里是3
	logger.InitLogger(3)

	// 初始化配置文件
	config.InitConf(configDir)
	config.InitRedis()

	// 添加prometheus性能监控指标
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	// 性能监控的端口port+1000,只能在内网访问
	httpMux := gpprof.New()

	// 添加prometheus metrics处理器
	httpMux.Handle("/metrics", promhttp.Handler())
	gpprof.Run(httpMux, port+1000)

}

// grpcHandlerFunc 使用标准库 h2c 将http请求进行转化为http2
// 通过这种方式就可以把go grpc和http 服务共存，用一个端口就可以做到同时提供grpc服务，又可以提供http 服务
// 在 2018 年 6 月，代表 “h2c” 标志的 golang.org/x/net/http2/h2c 标准库正式合并进来，自此我们就可以使用官方标准库（h2c）
// 这个标准库实现了 HTTP/2 的未加密模式，因此我们就可以利用该标准库在同个端口上既提供 HTTP/1.1 又提供 HTTP/2 的功能了
// h2c.NewHandler 方法进行了特殊处理，h2c.NewHandler 会返回一个 http.handler
// 主要的内部逻辑是拦截了所有 h2c 流量，然后根据不同的请求流量类型将其劫持并重定向到相应的 Hander 中去处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func main() {
	var grpcAddress = fmt.Sprintf("0.0.0.0:%d", port)
	var opts []grpc.ServerOption

	// 设置grpc服务参数
	// 设置超时10s
	opts = append(opts, grpc.ConnectionTimeout(10*time.Second))

	// 注册interceptor和中间件
	opts = append(opts, grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			middleware.RequestInterceptor,
			middleware.Limit(&middleware.MockPassLimiter{}),
		)))

	// 注册grpc服务
	server := grpc.NewServer(opts...)
	pb.RegisterGreeterServiceServer(server, &service.GreeterService{})

	dopts := []grpc.DialOption{grpc.WithInsecure()}
	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux() // grpc-gateway/runtime mux
	err := pb.RegisterGreeterServiceHandlerFromEndpoint(context.Background(), gwmux, grpcAddress, dopts)
	if err != nil {
		log.Fatalln("grpc register http gw err: ", err)
	}

	mux.Handle("/", gwmux)

	log.Println("go-proj grpc run on:", port)

	go func() {
		defer logger.Recover()

		if err := http.ListenAndServe(grpcAddress, grpcHandlerFunc(server, mux)); err != nil {
			log.Println("server run error: ", err)
		}
	}()

	// 平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal,please use this in production.
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	server.GracefulStop()

	<-ctx.Done()

	log.Println("shutting down")
}
