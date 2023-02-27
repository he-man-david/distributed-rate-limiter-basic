package proxy

import (
	"log"

	"github.com/he-man-david/distributed-rate-limiter-basic/rate"
	"golang.org/x/net/context"
)

type Proxy struct {
	ratelimiter *rate.RateLimiter
	UnimplementedProxyServer
}

func NewProxy() *Proxy {
	rl := &rate.RateLimiter{}
	return &Proxy{ratelimiter: rl}
}

func (p *Proxy) AllowRequest(ctx context.Context, req *AllowRequestReq) (*AllowRequestResp, error) {
	log.Printf("[Proxy] rate limiting key: %d", req.ApiKey)
	res, err := p.ratelimiter.AllowRequest(ctx, req.ApiKey)
	if err != nil {
		return nil, err
	}
	return &AllowRequestResp{Res: res}, nil
}