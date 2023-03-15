package ratetracker

import (
	"container/heap"
	"container/list"
	"log"
	"sync"
)

type earliestTimeLLMinHeap []*list.List
type requestsByNode map[int64]*list.List

// Starting the RateTrackerInstance with one mutex for all reads and writes.
// The rate of operations are expected to be in the ns range and the number
// of operations are on sec/min scale. So its unlikely to be a bootleneck.
// This will revisited later if it presents a problem.
type RateTrackerInstance struct {
	requestsByNode        requestsByNode
	earliestTimeLLMinHeap earliestTimeLLMinHeap
	maxRequests           int64
	maxTimePeriodMs       int64
	mutex                 *sync.RWMutex
	apiKey                int64
	localNodeId           int64
}

func NewRateTrackerInstance(maxRequests int64, maxTimePeriodMs int64, apiKey int64, localNodeId int64) *RateTrackerInstance {
	requestsByNode := make(requestsByNode)
	earliestTimeLLMinHeap := make(earliestTimeLLMinHeap, 0)
	newRti := RateTrackerInstance{
		earliestTimeLLMinHeap: earliestTimeLLMinHeap,
		requestsByNode:        requestsByNode,
		maxRequests:           maxRequests,
		maxTimePeriodMs:       maxTimePeriodMs,
		mutex:                 &sync.RWMutex{},
		apiKey:                apiKey,
		localNodeId:           localNodeId,
	}
	return &newRti
}

// This function will check across all sliding windows to determine if a request should be allowed or not
// if the request is allowed, it is recorded for future look ups
// The return is a tuple of booleans indicated allowed and if the last T1 was removed
func (rti *RateTrackerInstance) AllowRequest(nodeId int64, timestamp int64) (bool, bool) {
	count := rti.getRecordedRequestsCount()
	if count < rti.maxRequests {
		log.Printf("getRecordedRequestsCount() is = %d", count)
		recorded := rti.RecordRequest(nodeId, timestamp, false)
		return recorded, false
	}

	timeDiff := timestamp - rti.getT1Lock()
	if timeDiff >= rti.maxTimePeriodMs {
		// log.Printf("timeDiff is = %d", timeDiff)
		recored := rti.RecordRequest(nodeId, timestamp, true)
		return recored, true
	}

	return false, false
}

// this function will just record a request in a sliding windows for a node
func (rti *RateTrackerInstance) RecordRequest(nodeId int64, timestamp int64, takeT1 bool) bool {
	rti.mutex.Lock()
	defer rti.mutex.Unlock()

	existingList := rti.requestsByNode[nodeId]
	if existingList == nil {
		newList := list.New()
		newList.PushBack(timestamp)
		heap.Push(&rti.earliestTimeLLMinHeap, newList)
		rti.requestsByNode[nodeId] = newList
		existingList = newList
	}

	timeDiff := timestamp - rti.getT1NoLock()
	if int64(existingList.Len()) >= rti.maxRequests && timeDiff < rti.maxTimePeriodMs {
		return false
	}

	if takeT1 {
		rti.takeT1NoLock()
	}

	existingList.PushBack(timestamp)
	for existingList.Len() > int(rti.maxRequests) {
		existingList.Remove(existingList.Front())
	}
	return true
}

// this function will just record a request with no checks in a sliding windows for a node
func (rti *RateTrackerInstance) UpdateRequests(nodeId int64, timestamp int64, takeT1 bool) {
	log.Printf("Updating for nodeId = %d", nodeId)
	rti.mutex.Lock()
	defer rti.mutex.Unlock()

	existingList := rti.requestsByNode[nodeId]
	if existingList == nil {
		newList := list.New()
		newList.PushBack(timestamp)
		heap.Push(&rti.earliestTimeLLMinHeap, newList)
		rti.requestsByNode[nodeId] = newList
		// existingList = newList
		return
	}

	if takeT1 {
		rti.takeT1NoLock()
	}

	rti.LogSelf()
	existingList.PushBack(timestamp)
	for existingList.Len() > int(rti.maxRequests) {
		existingList.Remove(existingList.Front())
	}
}

// Gets all the size of all requests accross all linked lists
func (rti *RateTrackerInstance) getRecordedRequestsCount() int64 {
	rti.mutex.RLock()
	defer rti.mutex.RUnlock()
	var total int64
	for _, list := range rti.requestsByNode {
		total += int64(list.Len())
	}
	return total
}

// Gets the earliest time across all linked lists
func (rti *RateTrackerInstance) getT1Lock() int64 {
	rti.mutex.RLock()
	defer rti.mutex.RUnlock()
	return rti.getT1NoLock()
}

// Gets the earliest time across all linked lists
func (rti *RateTrackerInstance) getT1NoLock() int64 {
	if len(rti.earliestTimeLLMinHeap) == 0 || rti.earliestTimeLLMinHeap[0].Len() == 0 {
		return 0
	}
	globalT1 := rti.earliestTimeLLMinHeap[0].Front().Value.(int64)

	var localT1 int64 = 0
	localReqRecord := rti.requestsByNode[rti.localNodeId]
	if localReqRecord != nil && localReqRecord.Len() > 0 {
		localT1 = localReqRecord.Front().Value.(int64)
	}

	if localT1 > globalT1 {
		return localT1
	}
	return globalT1
}

// Gets and removes the earliest time across all linked lists
func (rti *RateTrackerInstance) takeT1NoLock() int64 {
	// remove current t1, earliest time
	lowestList := heap.Pop(&rti.earliestTimeLLMinHeap).(*list.List)

	// this check is if the Heap is empty
	if lowestList == nil {
		return -1
	}

	// get the next t in line in current LL
	t1 := lowestList.Front()
	lowestList.Remove(lowestList.Front())

	// add back to our min heap if list is not empty to maintain min sliding windows
	heap.Push(&rti.earliestTimeLLMinHeap, lowestList)
	return t1.Value.(int64)
}

// Helper log
func (rti *RateTrackerInstance) LogSelf() {
	for nodeId, list := range rti.requestsByNode {
		log.Printf("Node ID = %d\n", nodeId)
		log.Printf("Node requests count = %d \n", list.Len())
		if list != nil && list.Len() > 0 {
			log.Printf("Node T1 = %d \n", list.Front().Value)
		}
	}
}

// Function implementing the heap Interface function Len
func (e earliestTimeLLMinHeap) Len() int {
	return len(e)
}

// Function implementing the heap Interface function Less
func (e earliestTimeLLMinHeap) Less(i, j int) bool {
	// special case since the lists inside the heap could be empty
	if e[i].Len() == 0 {
		return false
	} else if e[j].Len() == 0 {
		return true
	}
	return e[i].Front().Value.(int64) < e[j].Front().Value.(int64)
}

// Function implementing the heap Interface function Swap
func (e earliestTimeLLMinHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Function implementing the heap Interface function Push
func (e *earliestTimeLLMinHeap) Push(x any) {
	*e = append(*e, x.(*list.List))
}

// Function implementing the heap Interface function Pop
func (e *earliestTimeLLMinHeap) Pop() any {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}
