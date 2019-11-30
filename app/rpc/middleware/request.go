package middleware

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/daheige/thinkgo/gutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-proj/library/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// request interceptor 请求拦截器，记录请求的基本信息
func RequestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	defer func() {
		// grpc recover异常捕获
		if r := recover(); r != nil {
			//这里必须采用grpc定义的错误格式返回code,desc
			err = status.Errorf(codes.Internal, "%s", "server inner error")

			//log.Println("exec error:", err)
			logger.Emergency(ctx, "exec error", map[string]interface{}{
				"reply":       res,
				"panic_error": r,
				"grpc_error":  err,
				"trace_error": string(debug.Stack()),
			})
		}
	}()

	t := time.Now()
	clientIp, _ := GetClientIp(ctx)
	//log.Println("client_ip: ", clientIp)

	//b, _ := json.Marshal(info)
	//log.Println("info: ")
	//log.Println(string(b))

	//log.Println("req: ", req)

	logId := gutils.RndUuid()
	ctx = context.WithValue(ctx, "log_id", logId)
	ctx = context.WithValue(ctx, "client_ip", clientIp)
	ctx = context.WithValue(ctx, "request_method", info.FullMethod)
	ctx = context.WithValue(ctx, "request_uri", info.FullMethod)
	ctx = context.WithValue(ctx, "user_agent", "grpc")

	logger.Info(ctx, "exec begin", nil)

	// 继续处理请求
	res, err = handler(ctx, req)
	ttd := time.Now().Sub(t).Seconds()
	if err != nil {
		logger.Error(ctx, "exec error", map[string]interface{}{
			"reply":       res,
			"trace_error": err,
			"exec_time":   ttd,
		})

		return nil, err
	}

	logger.Info(ctx, "exec end", map[string]interface{}{
		"exec_time": ttd,
	})

	return res, err
}

// GetClientIp 从上下文中获取客户端ip地址
func GetClientIp(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}

	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}
