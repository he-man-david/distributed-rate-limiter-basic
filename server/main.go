package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/proxy"
	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
	sync "github.com/he-man-david/distributed-rate-limiter-basic/server/synchronization"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9000, "the port number for the TCP server")
)

func main() {
	flag.Parse()

	fmt.Printf("Starting Rate Limiter Service on port :: %d \n", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	// creating RateTrackerImpl
	rt := ratetracker.RateTracker{}

	// registering Syncrhonization implementation
	s := sync.NewSyncImpl(&rt)
	defer s.Shutdown()
	sync.RegisterSyncServer(grpcServer, s)

	// registering Proxy implementation
	p := proxy.NewProxy(&rt, s)
	proxy.RegisterProxyServer(grpcServer, p)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}
