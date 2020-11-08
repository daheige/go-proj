package service

import (
	"context"
	"log"

	"github.com/daheige/go-proj/pb"
)

type GreeterService struct {
	BaseService
}

func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	//panic(111) 这里模拟的panic可以在请求拦截器中，自动捕获，记录操作日志
	log.Println("req data: ", in)
	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}
