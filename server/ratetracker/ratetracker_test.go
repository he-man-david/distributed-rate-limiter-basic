package ratetracker_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
)

type rtTestData struct {
	nodeId        int64
	timestamp     int64
	expectedValue bool
	t1Popped      bool
	isSyncMsg     bool
}

func TestRateTracker(t *testing.T) {
	var wg sync.WaitGroup

	apiKey := int64(1)
	rt := ratetracker.NewRateTracker(1, 2, 1000)

	// sequnce of events
	testSequece := []rtTestData{
		{isSyncMsg: true, nodeId: 2, timestamp: 100, t1Popped: false},
		{isSyncMsg: false, nodeId: 1, timestamp: 100, t1Popped: false, expectedValue: true},
		{isSyncMsg: true, nodeId: 2, timestamp: 100, t1Popped: true},
		{isSyncMsg: false, nodeId: 1, timestamp: 100, t1Popped: false, expectedValue: false},
		{isSyncMsg: false, nodeId: 1, timestamp: 2000, t1Popped: true, expectedValue: true},
	}

	for _, sequenceData := range testSequece {
		if sequenceData.isSyncMsg {
			msg := createSynMsg(&sequenceData, apiKey)
			rt.InSyncChan <- msg
		} else {
			// assert outchannel because Allow Request will block until some one recieves
			// only call the go routine if we are expecting Out Channel to have a sync message
			// other wise it will keep listening to channels for ever
			if sequenceData.expectedValue {
				expectedSyncMsg := createSynMsg(&sequenceData, apiKey)
				wg.Add(1)
				go assertOutChannelValue(&wg, t, rt.OutSyncChan, expectedSyncMsg)
			}

			// assert AllowRequest call
			actualValue := rt.AllowRequest(apiKey, sequenceData.timestamp)
			expectedValue := sequenceData.expectedValue
			if !(actualValue == expectedValue) {
				t.Errorf("Expected to find %v but found %v", expectedValue, actualValue)
			}
		}
	}
	fmt.Printf("Done with loop")
	wg.Wait()
	close(rt.InSyncChan)
	close(rt.OutSyncChan)
}

func createSynMsg(sequenceData *rtTestData, apiKey int64) *ratetracker.SyncMsg {
	return &ratetracker.SyncMsg{
		Timestamp:     sequenceData.timestamp,
		RatelimiterId: sequenceData.nodeId,
		ApiKey:        apiKey,
		PopT1:         sequenceData.t1Popped,
	}
}

func assertOutChannelValue(wg *sync.WaitGroup, t *testing.T, outSynChan chan *ratetracker.SyncMsg, expectedValue *ratetracker.SyncMsg) {
	defer wg.Done()
	for {
		fmt.Printf("Loop is running")
		select {
		case msg, ok := <-outSynChan:
			fmt.Printf("Msg from out chan = %v\n", msg)
			if !ok || msg.PopT1 != expectedValue.PopT1 || msg.Timestamp != expectedValue.Timestamp {
				t.Errorf("Expected to find %v but found %v", expectedValue, msg)
			}
			// we are only asserting one value
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
