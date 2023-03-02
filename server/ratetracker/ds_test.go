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

	if !(linkedList.TakeT1().Value.(int64) == 1) {
		t.Errorf("Expected to find 3")
	}
	if !(linkedList.TakeT1().Value.(int64) == 2) {
		t.Errorf("Expected to find 2")
	}
	if !(linkedList.TakeT1().Value.(int64) == 3) {
		t.Errorf("Expected to find 1")
	}
}
