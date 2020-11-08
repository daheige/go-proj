package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daheige/go-proj/app/web/routes"
	config "github.com/daheige/go-proj/conf"
	"github.com/daheige/go-proj/conf/grpcconf"

	"github.com/daheige/thinkgo/gpprof"
	"github.com/daheige/thinkgo/logger"
	"github.com/daheige/thinkgo/monitor"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "go.uber.org/automaxprocs"
)

var (
	port      int
	logDir    string
	configDir string
	wait      time.Duration // 平滑重启的等待时间1s or 1m
)

func init() {
	flag.IntVar(&port, "port", 1338, "app listen port")
	flag.StringVar(&logDir, "log_dir", "./logs", "log dir")
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	// 日志文件设置
	logger.SetLogDir(logDir)
	logger.SetLogFile("go-web.log")
	logger.MaxSize(500)
	logger.TraceFileLine(true) // 开启文件名和行数追踪

	// 由于app/extensions/logger基于thinkgo/logger又包装了一层，所以这里是3
	logger.InitLogger(3)

	// 初始化配置文件
	config.InitConf(configDir)
	config.InitRedis()

	// web服务中是否要初始化hello gRpc client
	if config.WebHasGRPCService {
		grpcconf.InitGRPCClient()
	}

	// 添加prometheus性能监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)

	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	// 性能监控的端口port+1000,只能在内网访问
	httpMux := gpprof.New()

	// 添加prometheus metrics处理器
	httpMux.Handle("/metrics", promhttp.Handler())
	gpprof.Run(httpMux, port+1000)

	// gin mode设置
	switch config.AppEnv {
	case "local", "dev":
		gin.SetMode(gin.DebugMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// 这里推荐使用gin.New方法，默认的Default方法的logger,recovery中间件
	// 有些项目也许用不到，另一方面gin recovery 中间件
	// 对于broken pipe存在一些情形无法覆盖到
	// 具体请参考 go-proj/app/web/middleware/log.go#74
	router := gin.New()

	// 加载路由文件中的路由
	routes.WebRoute(router)

	// 服务server设置
	server := &http.Server{
		Handler:           router,
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		IdleTimeout:       20 * time.Second, // tcp idle time
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	// 在独立携程中运行
	log.Println("server run on: ", port)
	go func() {
		defer logger.Recover()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Info("server close error", map[string]interface{}{
					"trace_error": err.Error(),
				})

				log.Println("server close error: ", err.Error())
				return
			}

			log.Println("server will exit...")
		}
	}()

	// server平滑重启
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

	// 平滑退出之前，先停止接收请求，但老的请求，继续响应
	// SetKeepAlivesEnabled controls whether HTTP keep-alives are enabled.
	// By default, keep-alives are always enabled. Only very
	// resource-constrained environments or servers in the process of
	// shutting down should disable them.
	server.SetKeepAlivesEnabled(false)

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	done := make(chan error, 1)
	go func() {
		defer logger.Recover()
		defer close(done)

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			logger.Error("server shutdown error", map[string]interface{}{
				"trace_error": err.Error(),
			})

			log.Println("server shutdown error: ", err.Error())
		}
	}()

	<-done
	log.Println("shutdown success")
	logger.Info("server shutdown success", nil)
}
