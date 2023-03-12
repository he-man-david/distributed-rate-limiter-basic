package proxy

import (
	"context"
	"time"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
	sync "github.com/he-man-david/distributed-rate-limiter-basic/server/synchronization"
)

type Proxy struct {
	rt   *ratetracker.RateTracker
	sync *sync.Sync

	UnimplementedProxyServer
}

func NewProxy(rt *ratetracker.RateTracker, sync *sync.Sync) *Proxy {
	return &Proxy{rt: rt, sync: sync}
}

func (p *Proxy) RegisterNode(ctx context.Context, node *RegisterNodeReq) (*RegisterNodeResp, error) {
	err := p.sync.RegisterNode(ctx, node.RateLimiterId, node.Port)
	if err != nil {
		return &RegisterNodeResp{Res: false}, err
	}
	return &RegisterNodeResp{Res: true}, nil
}

func (p *Proxy) AllowRequest(ctx context.Context, req *AllowRequestReq) (*AllowRequestResp, error) {
	// log.Printf("[Proxy] rate limiting key: %d", req.ApiKey)
	// epoch time in ms
	t := time.Now().UnixNano() / int64(time.Millisecond)
	// calling rate tracker allow request
	res := p.rt.AllowRequest(ctx, req.ApiKey, t)
	return &AllowRequestResp{Res: res}, nil
}

// TODO: Add resetRateLimiter -- for testing purpose, so after each test we can clear state
