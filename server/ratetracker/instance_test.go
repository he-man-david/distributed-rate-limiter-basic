package ratetracker_test

import (
	"testing"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
)

type rtiTestData struct {
	nodeId                   int64
	timestamp                int64
	expectedValueForAllowed  bool
	expectedValueForPoppedT1 bool
}

func TestRti(t *testing.T) {
	rti := ratetracker.NewRateTrackerInstance(3, 1000, 123456, 0)

	// sequnce of events
	testSequece := []rtiTestData{
		{nodeId: 0, timestamp: 1, expectedValueForAllowed: true, expectedValueForPoppedT1: false},
		{nodeId: 1, timestamp: 100, expectedValueForAllowed: true, expectedValueForPoppedT1: false},
		{nodeId: 0, timestamp: 100, expectedValueForAllowed: true, expectedValueForPoppedT1: false},
		{nodeId: 1, timestamp: 200, expectedValueForAllowed: false, expectedValueForPoppedT1: false},
		{nodeId: 2, timestamp: 2000, expectedValueForAllowed: true, expectedValueForPoppedT1: true},
		{nodeId: 1, timestamp: 2001, expectedValueForAllowed: true, expectedValueForPoppedT1: false},
	}

	for _, testData := range testSequece {
		actualValue, _ := rti.AllowRequest(testData.nodeId, testData.timestamp)
		expectedValue := testData.expectedValueForAllowed
		if !(actualValue == expectedValue) {
			t.Errorf("Expected to find %v but found %v", expectedValue, actualValue)
		}
	}
}
