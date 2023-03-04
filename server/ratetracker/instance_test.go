package ratetracker_test

import (
	"testing"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
)

type testData struct {
	nodeId        int64
	timestamp     int64
	expectedValue bool
}

func TestRti(t *testing.T) {
	rti := ratetracker.NewRateTrackerInstance(3, 1000, 123456)

	// sequnce of events
	testSequece := []testData{
		{nodeId: 0, timestamp: 1, expectedValue: true},
		{nodeId: 1, timestamp: 100, expectedValue: true},
		{nodeId: 0, timestamp: 100, expectedValue: true},
		{nodeId: 1, timestamp: 200, expectedValue: false},
		{nodeId: 2, timestamp: 2000, expectedValue: true},
		{nodeId: 1, timestamp: 2001, expectedValue: true},
	}

	for _, testData := range testSequece {
		actualValue := rti.AllowRequest(testData.nodeId, testData.timestamp)
		expectedValue := testData.expectedValue
		if !(actualValue == expectedValue) {
			t.Errorf("Expected to find %v but found %v", expectedValue, actualValue)
		}
	}
}
