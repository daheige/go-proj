package utils

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/daheige/thinkgo/common"

	"go-proj/library/Logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// request interceptor 请求拦截器，记录请求的基本信息
func RequestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	clientIp, _ := GetClientIp(ctx)
	//log.Println("client_ip: ", clientIp)

	//b, _ := json.Marshal(info)
	//log.Println("info: ")
	//log.Println(string(b))

	//log.Println("req: ", req)

	logId := common.RndUuid()
	ctx = context.WithValue(ctx, "log_id", logId)
	ctx = context.WithValue(ctx, "client_ip", clientIp)
	ctx = context.WithValue(ctx, "request_method", info.FullMethod)
	ctx = context.WithValue(ctx, "request_uri", info.FullMethod)
	ctx = context.WithValue(ctx, "user_agent", "grpc")

	Logger.Info(ctx, "exec begin", nil)

	// 继续处理请求
	res, err := handler(ctx, req)
	ttd := time.Now().Sub(t).Seconds()
	if err != nil {
		Logger.Error(ctx, "exec error", map[string]interface{}{
			"reply":       res,
			"trace_error": err,
			"exec_time":   ttd,
		})

		return nil, err
	}

	Logger.Info(ctx, "exec end", map[string]interface{}{
		"exec_time": ttd,
	})

	return res, err
}

//从上下文中获取客户端ip地址
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
