package ratetracker_test

import (
	"testing"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
)

func TestThreadSafeLL(t *testing.T) {
	linkedList := ratetracker.NewThreadSafeLL()
	linkedList.AddReq(1)
	linkedList.AddReq(2)
	linkedList.AddReq(3)

	for _, expectedValue := range []int64{1, 2, 3} {
		actualValue := linkedList.TakeT1().Value.(int64)
		if !(actualValue == expectedValue) {
			t.Errorf("Expected to find %d but found %d", expectedValue, actualValue)
		}
	}
	emptyListTakeT1 := linkedList.TakeT1()
	if emptyListTakeT1 != nil {
		t.Errorf("Expected to find %v but found %v", nil, emptyListTakeT1)
	}
}

func TestThreadSafeHeap(t *testing.T) {
	linkedListOne := ratetracker.NewThreadSafeLL()
	linkedListOne.AddReq(1)
	linkedListOne.AddReq(2)
	linkedListOne.AddReq(6)
	linkedListOne.AddReq(10)

	linkedListTwo := ratetracker.NewThreadSafeLL()
	linkedListTwo.AddReq(3)
	linkedListTwo.AddReq(8)
	linkedListTwo.AddReq(9)

	minHeap := ratetracker.NewThreadSafeMinHeap()
	minHeap.AddList(linkedListOne.GetList())
	minHeap.AddList(linkedListTwo.GetList())

	for _, expectedValue := range []int64{1, 2, 3, 6, 8, 9, 10} {
		actualValue := minHeap.TakeT1()
		if !(actualValue == expectedValue) {
			t.Errorf("Expected to find %d but found %d", expectedValue, actualValue)
		}
	}

	emptyHeapTakeT1 := minHeap.TakeT1()
	if emptyHeapTakeT1 != -1 {
		t.Errorf("Expected to find %d but found %d", -1, emptyHeapTakeT1)
	}
}
