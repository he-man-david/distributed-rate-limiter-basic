package main

import (
	"fmt"
	"log"
	"net"

	"github.com/he-man-david/distributed-rate-limiter-basic/rate"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Starting Rate Limiter Service....")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}

	s := rate.RateLimiter{}

	grpcServer := grpc.NewServer()

	rate.RegisterRateLimiterServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}