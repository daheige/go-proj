package grpcconf

import (
	"context"
	"log"
	"time"

	"github.com/daheige/thinkgo/logger"
	"google.golang.org/grpc"
)

var (
	// HellServiceClient hell grpc service client
	HellServiceClient *grpc.ClientConn
)

// InitGRPCClient 初始化gRpc客户端
func InitGRPCClient() {
	helloGRPCAddress := "localhost:50051"

	// 客户端3s连接超时
	/**
	建立连接默认只是返回一个ClientConn的指针,相当于new了一个ClientConn 把指针返回给你。
	并不是一定要建立真实的h2连接.至于真实的连接建立实际上是一个异步的过程。
	当然如果你想等真实的链接完全建立再返回ClientConn可以通过WithBlock传入Options来实现,
	这样的话链接如果建立不成功就会一直阻塞直到Context超时.
	如果你使用withBlock 但是不使用超时的话会不断的重试下去。中途断掉也会不断重联。
	当然了重连的过程中是使用了backoff算法来重连。
	而且默认会在grpc的配置中有个默认最大重试间隔时间。默认是120.
	grpc源码
	var DefaultBackoffConfig = BackoffConfig{
	    MaxDelay:  120 * time.Second,
	    baseDelay: 1.0 * time.Second,
	    factor:    1.6,
	    jitter:    0.2,
	}

	go grpc 底层重连机制，使用resetTransport 主要内容就是一个for 循环,可以看到在这个for循环中会尝试建立链接。
	如果建立成功就返回一个nil。如果不成功会不断重试下去。
	实际上不管是开头的Dial或者Dial完了关闭服务器后都是由这段代码来建立真实的链接。
	这也就是如果你使用withBlock 但是不使用超时的话会不断的重试下去。中途断掉也会不断重联。
	当然了重连的过程中是使用了backoff算法来重连。
	而且默认会在grpc的配置中有个默认最大重试间隔时间。默认是120

	参考link: https://www.jianshu.com/p/a5dec04d042b
	*/

	// 这里的超时时间，仅当存在WithBlock（）时，此方法才有效
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	/*
		WithBlock返回一个DialOption，它使Dial块的调用者直到
		基础连接已建立。 没有此功能，Dial将立即返回并
		连接服务器在后台进行。
	*/
	hellClient, err := SetGRPCClient(ctx, helloGRPCAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		//log.Println("hell grpc service connection error: ", err)
		return
	}

	HellServiceClient = hellClient
}

// SetGRPCClient 初始化gRpc client connection.
func SetGRPCClient(ctx context.Context, gRpcAddress string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if len(opts) == 0 {
		opts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
		}
	}

	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, gRpcAddress, opts...)
	if err != nil {
		log.Println("current grpc service connection error: ", err)

		// 记录日志到文件中
		logger.Fatal("grpc client connection error", map[string]interface{}{
			"trace_error":  err.Error(),
			"grpc_address": gRpcAddress,
		})

		return nil, err
	}

	// 如果是使用client复用的话，这里不需要关闭
	// 一般放在main函数中关闭 defer HellServiceClient.Close()
	// defer conn.Close()

	return conn, nil
}
