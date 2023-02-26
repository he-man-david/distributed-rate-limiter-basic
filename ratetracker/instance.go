package ratetracker

type RateTrackerInstance struct {
	
}

// TODO: Jaser
// when instantiating this struct, we want to "find" all ratelimite_id,
// create N linked list sliding window (see ds.go)
// then create our min heap, and add N sliding window to min heap.
// this will all be done in init
// this struct should also have funcs and store data for below
// 1. get total requests in all RL
// 2. perform tn - t1 calculation
// 3. check logic whether or not tn can be added, return success/fail
// 4. ... whatever else u think is good