package controller

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// HTTPSuccess 请求成功的时候 http code = 200
	// 当然也可以直接用http 包的200状态码
	HTTPSuccess = 200

	// APIError 业务code !=0的时候，默认API error code
	APIError = 500

	// APISuccess 业务成功code = 0,非0表示错误或异常
	APISuccess = 0
)

// EmptyArray 用作空[]返回
type EmptyArray []struct{}

// EmptyObject 空对象{}格式返回
type EmptyObject struct{}

// BaseController 基础控制器
type BaseController struct{}

func (ctrl *BaseController) ajaxReturn(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(HTTPSuccess, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

// Success returns code,data,message if crtl response success.
func (ctrl *BaseController) Success(ctx *gin.Context, message string, data interface{}) {
	if len([]rune(message)) == 0 {
		message = "ok"
	}

	ctrl.ajaxReturn(ctx, APISuccess, message, data)
}

// Error returns code,data,message if crtl response error.
func (ctrl *BaseController) Error(ctx *gin.Context, code int, message string, data interface{}) {
	if code <= 0 {
		code = APIError
	}

	ctrl.ajaxReturn(ctx, code, message, data)
}
