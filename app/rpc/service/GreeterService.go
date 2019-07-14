package service

import (
	"context"
	"go-proj/pb"
)

type GreeterService struct {
	BaseService
}

func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}
