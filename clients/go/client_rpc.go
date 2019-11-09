package main

import (
	"log"
	"os"

	pb "go-proj/pb"

	"context"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50050" //连接nginx grpc端口
	defaultName = "golang grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	res, err := c.SayHello(context.Background(), &pb.HelloReq{
		Name: name,
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("name:%s,message:%s", res.Name, res.Message)
}
