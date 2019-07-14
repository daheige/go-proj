package Logger

import (
	"context"
	"go-proj/library/helper"
	"runtime"
	"strings"

	"github.com/daheige/thinkgo/common"
	"github.com/daheige/thinkgo/logger"
)

/**
{
    "level":"info",
    "time_local":"2019-07-13T22:13:00.934+0800",
    "line":"/web/go/go-proj/app/extensions/Logger/log.go:61",
    "msg":"exec end",
    "tag":"index",
    "request_uri":"/index",
    "log_id":"6a03c58d-c324-bb24-304a-bc319d652e49",
    "options":{
        "exec_time":0.000566554
    },
    "plat":"web",
    "trace_line":51,
    "trace_file":"/web/go/go-proj/app/web/middleware/LogWare.go",
    "ip":"::1",
    "ua":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36",
    "request_method":"GET"
}
*/

func writeLog(ctx context.Context, levelName string, message string, options map[string]interface{}) {
	reqUri := getStringByCtx(ctx, "request_uri")
	tag := strings.Replace(reqUri, "/", "_", -1)
	tag = strings.Replace(tag, ".", "_", -1)
	tag = strings.TrimLeft(tag, "_")

	logId := getStringByCtx(ctx, "log_id")
	if logId == "" {
		logId = common.RndUuid()
	}

	ua := getStringByCtx(ctx, "user_agent")

	//函数调用
	_, file, line, _ := runtime.Caller(2)
	logInfo := map[string]interface{}{
		"tag":            tag,
		"request_uri":    reqUri,
		"log_id":         logId,
		"options":        options,
		"ip":             getStringByCtx(ctx, "client_ip"),
		"ua":             ua,
		"plat":           helper.GetDeviceByUa(ua), //当前设备匹配
		"request_method": getStringByCtx(ctx, "request_method"),
		"trace_line":     line,
		"trace_file":     file,
	}

	switch levelName {
	case "info":
		logger.Info(message, logInfo)
	case "debug":
		logger.Debug(message, logInfo)
	case "warn":
		logger.Warn(message, logInfo)
	case "error":
		logger.Error(message, logInfo)
	case "emergency":
		logger.DPanic(message, logInfo)
	default:
		logger.Info(message, logInfo)
	}
}

func getStringByCtx(ctx context.Context, key string) string {
	return helper.GetStringByCtx(ctx, key)
}

func Info(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "info", message, context)
}

func Debug(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "debug", message, context)
}

func Warn(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "warn", message, context)
}

func Error(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "error", message, context)
}

//致命错误或panic捕获
func Emergency(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "emergency", message, context)
}

//异常捕获处理
func Recover(c interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if ctx, ok := c.(context.Context); ok {
				Emergency(ctx, "exec panic", map[string]interface{}{
					"error":       err,
					"error_trace": string(logger.Stack()),
				})

				return
			}

			logger.DPanic("exec panic", map[string]interface{}{
				"error":       err,
				"error_trace": string(logger.Stack()),
			})
		}
	}()

}
