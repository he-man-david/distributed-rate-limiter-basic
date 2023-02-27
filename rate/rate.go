package rate

import (
	"context"
	"log"
)

type RateLimiter struct {
}

func (r *RateLimiter) allowRequest(ctx context.Context, apiKey int32) (bool, error) {
	log.Printf("[rateLimiter] : %d", apiKey)
	//TODO: Jaser this is the Rate Module
	// We want to check if we can allow this request, return T/F below
	return true, nil
}