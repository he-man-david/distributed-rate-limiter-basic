package rate

import (
	"context"
	"log"
)

type RateLimiter struct {
	UnimplementedRateServer
}

func NewRatelimiterImpl() *RateLimiter {
	return &RateLimiter{}
}

func (r *RateLimiter) AllowRequest(ctx context.Context, req *AllowRequestReq) (*AllowRequestResp, error) {
	log.Printf("[RateLimiter] rate limiting key: %d", req.ApiKey)
	// TODO: Jaser add rate limiting here
	return &AllowRequestResp{Res: true}, nil
}