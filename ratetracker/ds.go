package ratetracker

import (
	"container/heap"
	"container/list"
	"sync"
)

// Thread safe linked list for sliding window, storing timestamps for each API_KEY
// Since our rate limiter will be distributed, we need to maintain N sliding windows,
// for N replications. Each API_KEY will have N linked list sliding window, which
// the combine will constitute the current state/status of rate limiting for given KEY.
type ThreadSafeLL struct {
    list  *list.List
    mutex sync.RWMutex
}

func NewThreadSafeLL() *ThreadSafeLL {
    return &ThreadSafeLL{list: list.New()}
}

// get total # of request in current window
func (t *ThreadSafeLL) TotalReqs() int {
    t.mutex.RLock()
    defer t.mutex.RUnlock()
    return t.list.Len()
}

// add request to our current window, Tn
 func (t *ThreadSafeLL) AddReq(value int32) {
    t.mutex.Lock()
    defer t.mutex.Unlock()
    t.list.PushBack(value)
}

// get T1 of current window (earliest time)
func (t *ThreadSafeLL) GetT1() *list.Element {
    t.mutex.Lock()
    defer t.mutex.Unlock()
    return t.list.Front()
}

// remove T1 (earliest time)
func (t *ThreadSafeLL) TakeT1(value int32) *list.Element {
    t.mutex.Lock()
    defer t.mutex.Unlock()
	t1 := t.GetT1()
    t.list.Remove(t1)
	return t1
}


// Thread safe min heap, used to maintain our N linked list sliding windows. We need
// to perform Tn - T1 operation, to determine the current window size in ms. To do this,
// we need to know "what is T1" across all our replications, which will be synchronized.
// Since N replications are distributed, they will be eventually consistent, with some
// latency in synchronization. A small rate of error and inconsistency is tolerable.
type ThreadSafeMinHeap struct {
    heap  heapElements
    mutex sync.RWMutex
}

func NewThreadSafeMinHeap(t1s heapElements) *ThreadSafeMinHeap {
    return &ThreadSafeMinHeap{heap: t1s}
}

// return # of node in our min heap
func (h *ThreadSafeMinHeap) Len() int {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    return h.heap.Len()
}

// Add new LL head to our min heap, LL head will be from our N sliding windows
func (h *ThreadSafeMinHeap) AddT1(t1 *list.Element) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    heap.Push(h.heap, t1)
}

// removes t1 from our min heap
func (t *ThreadSafeMinHeap) TakeT1() *list.Element {
    t.mutex.Lock()
    defer t.mutex.Unlock()
	// remove current t1, earliest time
    t1 := heap.Pop(&t.heap).(*list.Element)
	// get the next t in line in current LL
	nextT := t1.Next()
	// add back to our min heap to maintain N sliding windows
	t.AddT1(nextT)
	return t1
}

func (t *ThreadSafeMinHeap) GetT1() *list.Element {
    t.mutex.RLock()
    defer t.mutex.RUnlock()
    return t.heap[0]
}

// wrapper element for linked list elements
// below are interface func required by container/heap implementation
// Our heap is just an array of elements, and we using container/heap
// heap.Push(h.heap, t1) to perform heapify. Similar to Java's heap
type heapElements []*list.Element

func (e heapElements) Len() int {
    return len(e)
}

func (e heapElements) Less(i, j int) bool {
    return e[i].Value.(int32) < e[j].Value.(int32)
}

func (e heapElements) Swap(i, j int) {
    e[i], e[j] = e[j], e[i]
}

func (e heapElements) Push(x interface{}) {
    e = append(e, x.(*list.Element))
}

func (e heapElements) Pop() interface{} {
    old := e
    n := len(old)
    x := old[n-1]
    e = old[0 : n-1]
    return x
}