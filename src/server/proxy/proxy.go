package proxy

import (
	"context"

	"github.com/he-man-david/distributed-rate-limiter-basic/src/server/rate"
	sync "github.com/he-man-david/distributed-rate-limiter-basic/src/server/synchronization"
)

type Proxy struct {
	ratelimiter *rate.RateLimiter
	sync        *sync.Sync

	UnimplementedProxyServer
}

func NewProxy(ratelimiter *rate.RateLimiter, sync *sync.Sync) *Proxy {
	return &Proxy{ratelimiter: ratelimiter, sync: sync}
}

func (p *Proxy) RegisterNode(ctx context.Context, node *RegisterNodeReq) (*RegisterNodeResp, error) {
	err := p.sync.RegisterNode(ctx, node.RateLimiterId, node.Port)
	if err != nil {
		return &RegisterNodeResp{Res: false}, err
	}
	return &RegisterNodeResp{Res: true}, nil
}

func (p *Proxy) AllowRequest(ctx context.Context, req *AllowRequestReq) (*AllowRequestResp, error) {
	return &AllowRequestResp{Res: true}, nil
}
