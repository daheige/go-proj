package main

import (
	"context"
	"flag"
	"fmt"
	config "go-proj/conf"
	"go-proj/healthCheck/pprofCheck"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron"

	"github.com/daheige/thinkgo/monitor"

	"github.com/daheige/thinkgo/logger"

	"go-proj/app/worker/job"
	"go-proj/app/worker/task"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var port int
var log_dir string
var config_dir string
var wait time.Duration //平滑重启的等待时间1s or 1m

//go:generate sh ../../bin/web-check-version.sh
func init() {
	flag.IntVar(&port, "port", 30031, "app listen port")
	flag.StringVar(&log_dir, "log_dir", "./logs", "log dir")
	flag.StringVar(&config_dir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	//日志文件设置
	logger.SetLogDir(log_dir)
	logger.SetLogFile("go-job.log")
	logger.MaxSize(500)
	logger.InitLogger()

	//初始化配置文件
	config.InitConf(config_dir)
	config.InitRedis()

	//注册监控指标
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	//性能监控的端口port+1000,只能在内网访问
	go func() {
		defer logger.Recover()

		log.Println("server pprof run on: ", port)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", pprofCheck.HealthHandler)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), httpMux); err != nil {
			log.Println(err)
		}
	}()
}

func main() {
	log.Println("===worker service start===")
	j := &job.TestJob{}

	//设置id到上下文上
	jCtx := context.WithValue(context.Background(), "id", "heige")
	j.SetCtx(jCtx)

	testTask := &task.TestTask{}
	c := cron.New()
	c.AddFunc("*/3 * * * * *", j.Info)         //每隔3s执行
	c.AddFunc("*/2 * * * * *", testTask.Hello) //每隔2s执行
	c.Start()

	//平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

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
