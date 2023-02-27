package rate

import (
	"log"

	"golang.org/x/net/context"
)

type RateLimiter struct {
	UnimplementedRateLimiterServer
}

func (r *RateLimiter) AllowRequest(ctx context.Context, in *AllowRequestParam) (*AllowRequestResponse, error) {
	log.Printf("[RateLimiter] rate limiting key: %d", in.ApiKey)
	//TODO: Jaser this is the Rate Module
	// We want to check if we can allow this request, return T/F below
	return &AllowRequestResponse{Res: true}, nil
}
