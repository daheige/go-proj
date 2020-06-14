package main

import (
	"context"
	"flag"
	config "go-proj/conf"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/gpprof"
	"github.com/daheige/thinkgo/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"

	"github.com/daheige/thinkgo/logger"

	"go-proj/app/worker/job"
	"go-proj/app/worker/task"

	_ "go.uber.org/automaxprocs"
)

var port int
var logDir string
var configDir string
var wait time.Duration // 平滑重启的等待时间1s or 1m

func init() {
	flag.IntVar(&port, "port", 30031, "app listen port")
	flag.StringVar(&logDir, "log_dir", "./logs", "log dir")
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful_timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	// 日志文件设置
	logger.SetLogDir(logDir)
	logger.SetLogFile("go-job.log")
	logger.MaxSize(500)
	logger.TraceFileLine(true) // 开启文件名和行数追踪

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
	log.Println("===worker service start===")
	j := &job.TestJob{}

	// 设置id到上下文上
	jCtx := context.WithValue(context.Background(), "id", "heige")
	j.SetCtx(jCtx)

	testTask := &task.TestTask{}
	c := cron.New(cron.WithSeconds())          // 具体用法可以看github.com/robfig/cron
	c.AddFunc("*/3 * * * * *", j.Info)         // 每隔3s执行
	c.AddFunc("*/2 * * * * *", testTask.Hello) // 每隔2s执行
	c.Start()

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

	c.Stop()
	<-ctx.Done()

	log.Println("shutting down")
}
