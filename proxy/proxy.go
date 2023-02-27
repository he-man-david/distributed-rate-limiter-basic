package proxy

import (
	"log"

	"github.com/he-man-david/distributed-rate-limiter-basic/rate"
	"golang.org/x/net/context"
)

type Proxy struct {
	rate.RateLimiter
	UnimplementedProxyServer
}

func NewProxy() *Proxy {
	rl := rate.RateLimiter{}
	unimplProxyS := UnimplementedProxyServer{}
	return &Proxy{rl, unimplProxyS}
}


func (p *Proxy) AllowRequest(ctx context.Context, req *AllowRequestReq) (*AllowRequestResp, error) {
	log.Printf("[Proxy] rate limiting key: %d", req.ApiKey)
	//TODO: Jaser this is the Rate Module
	// We want to check if we can allow this request, return T/F below
	return &AllowRequestResp{Res: true}, nil
}