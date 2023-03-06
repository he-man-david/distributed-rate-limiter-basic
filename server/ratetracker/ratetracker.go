package ratetracker

import (
	"log"
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
}

func NewRateTracker(id int64, defaultMaxRequests int64, defaultMaxTimePeriod int64) *RateTracker {
	rt := RateTracker{
		id:                   id,
		rtiByApiKey:          make(map[int64]*RateTrackerInstance),
		InSyncChan:           make(chan *SyncMsg),
		OutSyncChan:          make(chan *SyncMsg),
		defaultMaxRequests:   defaultMaxRequests,
		defaultMaxTimePeriod: defaultMaxTimePeriod,
	}
	go rt.openInChannel()
	return &rt
}

// This is the main function used to determine if a request for an apiKey is allwoed or not
// If the request is allowed a message is sent to the OutSyncChannel
func (rt *RateTracker) AllowRequest(apiKey int64, timestamp int64) bool {
	rti := rt.getOrCreateRti(apiKey)
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
			rt.syncRequest(msg)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// used to sync message by recording request in a RateTrackerInstance
func (rt *RateTracker) syncRequest(syncMsg *SyncMsg) {
	rti := rt.getOrCreateRti(syncMsg.ApiKey)
	rti.RecordRequest(syncMsg.RatelimiterId, syncMsg.Timestamp, syncMsg.PopT1)
}

// gets (or creates if not existing) a RateTrackerInstance for an apiKey
func (rt *RateTracker) getOrCreateRti(apiKey int64) *RateTrackerInstance {
	rti := rt.rtiByApiKey[apiKey]
	if rti == nil {
		newRti := NewRateTrackerInstance(rt.defaultMaxRequests, rt.defaultMaxTimePeriod, apiKey)
		rt.rtiByApiKey[apiKey] = newRti
		return newRti
	}
	return rti
}
