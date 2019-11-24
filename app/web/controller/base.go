package controller

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	HTTP_SUCCESS_CODE = 200
	HTTP_ERROR_CODE   = 500
	API_SUCCESS_CODE  = 0
)

//用作空[]返回
type EmptyArray []struct{}

type BaseController struct{}

func (ctrl *BaseController) ajaxReturn(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(HTTP_SUCCESS_CODE, gin.H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

func (ctrl *BaseController) Success(ctx *gin.Context, message string, data interface{}) {
	if len([]rune(message)) == 0 {
		message = "ok"
	}

	ctrl.ajaxReturn(ctx, API_SUCCESS_CODE, message, data)
}

//错误处理code,message
func (ctrl *BaseController) Error(ctx *gin.Context, code int, message string) {
	if code <= 0 {
		code = HTTP_ERROR_CODE
	}

	ctrl.ajaxReturn(ctx, code, message, nil)
}
