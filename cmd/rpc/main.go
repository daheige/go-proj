package main

import (
	"context"
	"flag"
	"fmt"
	"go-proj/app/rpc/middleware"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/gpprof"
	"github.com/daheige/thinkgo/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	flag.IntVar(&port, "port", 50051, "grpc port")
	flag.StringVar(&logDir, "log_dir", "./logs", "log dir")
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	// 日志文件设置
	logger.SetLogDir(logDir)
	logger.SetLogFile("go-grpc.log")
	logger.MaxSize(500)
	logger.TraceFileLine(true) //开启文件名和行数追踪

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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	// 设置超时10s
	opts = append(opts, grpc.ConnectionTimeout(10*time.Second))

	// 注册interceptor和中间件
	opts = append(opts, grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			middleware.RequestInterceptor,
			middleware.Limit(&middleware.MockPassLimiter{}),
		)))

	server := grpc.NewServer(opts...)
	pb.RegisterGreeterServiceServer(server, &service.GreeterService{})

	// 其他grpc拦截器用法，看go grpc源代码，里面都有对应的方法
	// Go-gRPC 实践指南 https://www.bookstack.cn/read/go-grpc/chapter2-interceptor.md
	log.Println("go-proj grpc run on:", port)

	go func() {
		defer logger.Recover()

		if err = server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
