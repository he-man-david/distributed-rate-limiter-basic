package clients

import (
	"fmt"
	"log"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/proxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ctx = context.Background()
)

type Clients struct {
	Conns   map[int]*grpc.ClientConn
	Proxies map[int]proxy.ProxyClient
}

// range of port nums, 8001 --- 8003
func NewClients(start, end int) (*Clients, error) {
	c := &Clients{
		Conns:   make(map[int]*grpc.ClientConn),
		Proxies: make(map[int]proxy.ProxyClient),
	}

	// create connection to our 3 test server endpoints 8001,8002,8003
	for port := start; port <= end; port++ {
		p, err := c.connect(port)
		if err != nil {
			log.Printf("[Testclient] grpc did not connect: %v", err)
			return nil, err
		}
		c.Proxies[port] = p
	}
	// register all the nodes
	if start != end {
		err := c.registerNodes(ctx, start, end)
		if err != nil {
			log.Printf("[Testclient] failed to register node: %v", err)
			return nil, err
		}
	}

	return c, nil
}

func (c *Clients) connect(port int) (proxy.ProxyClient, error) {
	// Connect to grpc port
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.Conns[port] = conn
	// create proxy client to Rate Limiter
	p := proxy.NewProxyClient(conn)
	return p, nil
}

func (c *Clients) registerNodes(ctx context.Context, start, end int) error {
	for port := start; port <= end; port++ {
		p := c.Proxies[port]
		for peer := start; peer <= end; peer++ {
			if peer != port {
				log.Printf("[Testclient] registering node at port: %d with: %d", peer, port)
				res, err := p.RegisterNode(ctx, &proxy.RegisterNodeReq{RateLimiterId: int64(peer), Port: int64(peer)})
				if err != nil || !res.Res {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Clients) AllowRequest(port int, apiKey int) (bool, error) {
	p, ok := c.Proxies[port]
	if !ok {
		log.Printf("[testclient] invalid port number for proxy: %d", port)
		return false, fmt.Errorf("[testclient] invalid port number for proxy: %d", port)
	}
	response, err := p.AllowRequest(context.Background(), &proxy.AllowRequestReq{ApiKey: int64(apiKey)})
	if err != nil {
		log.Printf("[testclient] Error when calling AllowRequest for proxy: %d, ERR: %s", port, err)
		return false, err
	}
	log.Printf("[testclient] Response from server for proxy: %d, RES: %t", port, response.Res)
	return response.Res, nil
}

func (c *Clients) AllowRequests(apiKey int) (int32, error) {
	var allowedCount int32
	for id := range c.Proxies {
		allowed, err := c.AllowRequest(id, apiKey)
		if err != nil {
			return allowedCount, err
		}
		if allowed {
			allowedCount++
		}
	}
	return allowedCount, nil
}
