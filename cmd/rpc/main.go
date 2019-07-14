package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"go-proj/pb"

	"go-proj/app/rpc/service"

	"google.golang.org/grpc"
)

var port int

func init() {
	flag.IntVar(&port, "port", 50051, "grpc port")
	flag.Parse()
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &service.GreeterService{})

	//其他grpc拦截器用法，看go grpc源代码，里面都有对应的方法
	// Go-gRPC 实践指南 https://www.bookstack.cn/read/go-grpc/chapter2-interceptor.md
	log.Println("go-proj grpc run on:", port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
