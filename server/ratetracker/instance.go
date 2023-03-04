package ratetracker

import (
	"container/heap"
	"container/list"
	"sync"
)

type earliestTimeLLMinHeap []*list.List
type requestsByNode map[int64]*list.List

type RateTrackerInstance struct {
	requestsByNode        requestsByNode
	earliestTimeLLMinHeap earliestTimeLLMinHeap
	maxRequests           int64
	maxTimePeriodMs       int64
	mutex                 *sync.RWMutex
	apiKey                int64
}

func NewRateTrackerInstance(maxRequests int64, maxTimePeriodMs int64, apiKey int64) *RateTrackerInstance {
	requestsByNode := make(requestsByNode)
	earliestTimeLLMinHeap := make(earliestTimeLLMinHeap, 0)
	newRti := RateTrackerInstance{
		earliestTimeLLMinHeap: earliestTimeLLMinHeap,
		requestsByNode:        requestsByNode,
		maxRequests:           maxRequests,
		maxTimePeriodMs:       maxTimePeriodMs,
		mutex:                 &sync.RWMutex{},
		apiKey:                apiKey,
	}
	return &newRti
}

func (rti *RateTrackerInstance) AllowRequest(nodeId int64, timestamp int64) bool {
	if rti.getRecordedRequestsCount() < rti.maxRequests {
		rti.RecordRequest(nodeId, timestamp)
		return true
	}
	timeDiff := timestamp - rti.getT1()
	if timeDiff < rti.maxTimePeriodMs {
		return false
	}

	rti.takeT1()
	rti.RecordRequest(nodeId, timestamp)
	return true
}

func (rti *RateTrackerInstance) RecordRequest(nodeId int64, timestamp int64) {
	rti.mutex.Lock()
	defer rti.mutex.Unlock()
	existingList := rti.requestsByNode[nodeId]
	if existingList == nil {
		newList := list.New()
		newList.PushBack(timestamp)
		heap.Push(&rti.earliestTimeLLMinHeap, newList)
		rti.requestsByNode[nodeId] = newList
	} else {
		existingList.PushBack(timestamp)
	}
}

func (rti *RateTrackerInstance) getRecordedRequestsCount() int64 {
	rti.mutex.RLock()
	defer rti.mutex.RUnlock()
	var total int64
	for _, list := range rti.requestsByNode {
		total += int64(list.Len())
	}
	return total
}

func (rti *RateTrackerInstance) getT1() int64 {
	rti.mutex.RLock()
	defer rti.mutex.RUnlock()
	if len(rti.earliestTimeLLMinHeap) == 0 || rti.earliestTimeLLMinHeap[0].Len() == 0 {
		return -1
	}
	return rti.earliestTimeLLMinHeap[0].Front().Value.(int64)
}

func (rti *RateTrackerInstance) takeT1() int64 {
	rti.mutex.Lock()
	defer rti.mutex.Unlock()

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

func (e earliestTimeLLMinHeap) Len() int {
	return len(e)
}

func (e earliestTimeLLMinHeap) Less(i, j int) bool {
	// special case since the lists inside the heap could be empty
	if e[i].Len() == 0 {
		return false
	} else if e[j].Len() == 0 {
		return true
	}
	return e[i].Front().Value.(int64) < e[j].Front().Value.(int64)
}

func (e earliestTimeLLMinHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *earliestTimeLLMinHeap) Push(x any) {
	*e = append(*e, x.(*list.List))
}

func (e *earliestTimeLLMinHeap) Pop() any {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}
