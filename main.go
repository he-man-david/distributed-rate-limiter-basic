package main

import (
	"fmt"
	"log"
	"net"

	"github.com/he-man-david/distributed-rate-limiter-basic/proxy"
	"github.com/he-man-david/distributed-rate-limiter-basic/rate"
	sync "github.com/he-man-david/distributed-rate-limiter-basic/synchronization"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Rate Limiter Service....")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// registering Rate Limiter implementation
	rl := rate.NewRatelimiterImpl()
	rate.RegisterRateServer(grpcServer, rl)

	// registering Syncrhonization implementation
	s := sync.NewSyncImpl()
	sync.RegisterSyncServer(grpcServer, s)

	// registering Proxy implementation
	p := proxy.NewProxy(rl, s)
	proxy.RegisterProxyServer(grpcServer, p)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}