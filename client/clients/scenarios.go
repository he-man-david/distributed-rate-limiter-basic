package clients

import (
	"log"
	"sync"
	"sync/atomic"
)

func TestExceedLimit(c *Clients) {
	apiKey := 1111
	iterations := 100
	iterationsPerGoRoutine := 100

	var count int32 = 0
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go exceedLimit(&wg, c, apiKey, iterationsPerGoRoutine, &count)
	}
	wg.Wait()
	// time.Sleep(time.Second * 5)
	log.Printf("Total count = %v\n", count)
}

func exceedLimit(wg *sync.WaitGroup, c *Clients, apiKey int, iterrations int, count *int32) error {
	defer wg.Done()
	var totalAllowedCount int32
	for i := 0; i < iterrations; i++ {
		allowedCount, err := c.AllowRequests(apiKey)
		if err != nil {
			return err
		}
		totalAllowedCount += allowedCount
	}
	atomic.AddInt32(count, totalAllowedCount)
	return nil
}
