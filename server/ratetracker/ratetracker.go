package ratetracker

import (
	"context"
	"log"
	"sync"
	"time"
)

type SyncMsg struct {
	RatelimiterId int64
	ApiKey        int64
	Timestamp     int64
	PopT1         bool
}

type RateTracker struct {
	id                   int64
	rtiByApiKey          map[int64]*RateTrackerInstance
	InSyncChan           chan *SyncMsg
	OutSyncChan          chan *SyncMsg
	defaultMaxRequests   int64
	defaultMaxTimePeriod int64
	mutex                 *sync.Mutex
}

var (
	ctx = context.Background()
)

func NewRateTracker(id int64, defaultMaxRequests int64, defaultMaxTimePeriod int64) *RateTracker {
	rt := RateTracker{
		id:                   id,
		rtiByApiKey:          make(map[int64]*RateTrackerInstance),
		InSyncChan:           make(chan *SyncMsg),
		OutSyncChan:          make(chan *SyncMsg),
		defaultMaxRequests:   defaultMaxRequests,
		defaultMaxTimePeriod: defaultMaxTimePeriod,
		mutex:				  &sync.Mutex{},
	}
	go rt.openInChannel()
	return &rt
}

// This is the main function used to determine if a request for an apiKey is allwoed or not
// If the request is allowed a message is sent to the OutSyncChannel
func (rt *RateTracker) AllowRequest(ctx context.Context, apiKey int64, timestamp int64) bool {
	rti := rt.getOrCreateRti(ctx, apiKey)
	allow, t1Popped := rti.AllowRequest(rt.id, timestamp)
	if !allow {
		return false
	}
	msg := &SyncMsg{ApiKey: apiKey, RatelimiterId: rt.id, Timestamp: timestamp, PopT1: t1Popped}
	rt.OutSyncChan <- msg
	return true
}

// Function to open InSync channel and continue to listen to it for sync messages
func (rt *RateTracker) openInChannel() {
	for {
		select {
		case msg, ok := <-rt.InSyncChan:
			if !ok {
				log.Println("In channel is closed!")
				return
			}
			rt.syncRequest(ctx, msg)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// used to sync message by recording request in a RateTrackerInstance
func (rt *RateTracker) syncRequest(ctx context.Context, syncMsg *SyncMsg) {
	rti := rt.getOrCreateRti(ctx, syncMsg.ApiKey)
	rti.RecordRequest(syncMsg.RatelimiterId, syncMsg.Timestamp, syncMsg.PopT1)
}

// gets (or creates if not existing) a RateTrackerInstance for an apiKey
func (rt *RateTracker) getOrCreateRti(ctx context.Context, apiKey int64) *RateTrackerInstance {
	rti := rt.rtiByApiKey[apiKey]

	if rti == nil {
		return rt.lockAndMakeNewInstance(apiKey)
	}
	return rti
}

func (rt *RateTracker) lockAndMakeNewInstance(apiKey int64) *RateTrackerInstance {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	// When rate instance does not exist yet, and we are setting the key for the first time
	// problem can occur where many concurrent req can come in to create Rti
	// We need to have optimistic locking here, where we lock and when commiting a new Rti
	// we check that it has not already been created.
	if _rti := rt.rtiByApiKey[apiKey]; _rti != nil {
		return _rti
	}
	newRti := NewRateTrackerInstance(rt.defaultMaxRequests, rt.defaultMaxTimePeriod, apiKey)
	rt.rtiByApiKey[apiKey] = newRti
	return newRti
}
