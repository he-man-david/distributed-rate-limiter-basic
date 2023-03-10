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
	maxrate = flag.Int("maxrate", 1000, "the max request per time period")
	timeout = flag.Int("timeout", 60000, "the rate limite time period")
)

func main() {
	flag.Parse()

	fmt.Println("Starting Rate Limiter Service.....")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("[main] TCP failed to listen: %v", err)
	}
	defer lis.Close()

	addr := lis.Addr().String()
	fmt.Printf("\n Listening address %s \n", addr)

	grpcServer := grpc.NewServer()

	// creating RateTrackerImpl
	rt := ratetracker.NewRateTracker(int64(*port), int64(*maxrate), int64(*timeout))

	// registering Syncrhonization implementation
	s := sync.NewSyncImpl(rt, port)
	defer s.Shutdown()
	sync.RegisterSyncServer(grpcServer, s)

	// registering Proxy implementation
	p := proxy.NewProxy(rt, s)
	proxy.RegisterProxyServer(grpcServer, p)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[main] gRPC failed to serve: %s", err)
	}
}
