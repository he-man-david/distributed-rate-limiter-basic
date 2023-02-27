package main

import (
	"fmt"
	"log"
	"net"

	"github.com/he-man-david/distributed-rate-limiter-basic/proxy"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Rate Limiter Service....")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	proxyImpl := proxy.NewProxy();
	proxy.RegisterProxyServer(grpcServer, proxyImpl)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}