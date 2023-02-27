package main

import (
	"fmt"
	"log"
	"net"

	"github.com/he-man-david/distributed-rate-limiter-basic/api/gRPC/rate"
	"google.golang.org/grpc"
)

func main() {


	fmt.Println("Starting Rate Limiter Service....")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}

	s := rate.Server{}

	grpcServer := grpc.NewServer()

	rate.RegisterChatServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}