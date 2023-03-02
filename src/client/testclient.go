package main

import (
	"log"

	"github.com/he-man-david/distributed-rate-limiter-basic/src/server/proxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[testclient] grpc did not connect: %s", err)
	}
	defer conn.Close()

	c := proxy.NewProxyClient(conn)

	response, err := c.AllowRequest(context.Background(), &proxy.AllowRequestReq{ApiKey: 12345})
	if err != nil {
		log.Fatalf("[testclient] Error when calling AllowRequest: %s", err)
	}
	log.Printf("[testclient] Response from server: %t", response.Res)

}
