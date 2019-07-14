package middleware

import (
	"go-proj/library/Logger"
	"go-proj/library/helper"
	"time"

	"github.com/daheige/thinkgo/common"

	"github.com/gin-gonic/gin"
)

type LogWare struct{}

func (ware *LogWare) Access() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		//uri := ctx.Request.RequestURI

		//性能分析后发现log.Println输出需要分配大量的内存空间,而且每次写入都需要枷锁处理
		//log.Println("request before")
		//log.Println("request uri: ", uri)

		//如果采用了nginx x-request-id功能，可以获得x-request-id
		logId := ctx.GetHeader("X-Request-Id")
		if logId == "" {
			logId = common.RndUuid() //日志id
		}

		//设置跟请求相关的ctx信息
		ctx.Request = helper.ContextSet(ctx.Request, "log_id", logId)
		ctx.Request = helper.ContextSet(ctx.Request, "client_ip", ctx.ClientIP())
		ctx.Request = helper.ContextSet(ctx.Request, "request_uri", ctx.Request.RequestURI)
		ctx.Request = helper.ContextSet(ctx.Request, "user_agent", ctx.GetHeader("User-Agent"))
		ctx.Request = helper.ContextSet(ctx.Request, "request_method", ctx.Request.Method)

		Logger.Info(ctx.Request.Context(), "exec start", nil)

		ctx.Next()

		//log.Println("request end")
		//请求结束记录日志
		c := map[string]interface{}{
			"exec_time": time.Now().Sub(t).Seconds(),
		}

		if code := ctx.Writer.Status(); code != 200 {
			c["response_code"] = code
		}

		Logger.Info(ctx.Request.Context(), "exec end", c)
	}
}

//请求处理中遇到异常或panic捕获
func (ware *LogWare) Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//log.Printf("error:%v", err)
				Logger.Emergency(ctx.Request.Context(), "exec panic", map[string]interface{}{
					"trace_error": err,
					"trace_info":  string(common.CatchStack()),
				})

				//响应状态
				ctx.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "server error",
				})

				return
			}
		}()

		ctx.Next()
	}
}
