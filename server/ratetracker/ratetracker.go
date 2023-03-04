package ratetracker

type SyncMsg struct {
	RatelimiterId int64
	ApiKey        int64
	Timestamp     int64
}

type RateTracker struct {
	id          int64
	rtiByApiKey map[int64]*RateTrackerInstance
	InSyncChan  chan *SyncMsg
	OutSyncChan chan *SyncMsg
}

func NewRateTracker(id int64) *RateTracker {
	return &RateTracker{id: id, rtiByApiKey: make(map[int64]*RateTrackerInstance)}
}

func (rt *RateTracker) AllowRequest(apiKey int64, timestamp int64) bool {
	rti := rt.rtiByApiKey[apiKey]
	if rti == nil {
		newRti := NewRateTrackerInstance(10, 1000, apiKey)
		rti = newRti
	}
	return rti.AllowRequest(rt.id, timestamp)
}

func (rt *RateTracker) SyncRequest(syncMsg SyncMsg) {
	rti := rt.rtiByApiKey[syncMsg.ApiKey]
	if rti == nil {
		newRti := NewRateTrackerInstance(10, 1000, syncMsg.ApiKey)
		rti = newRti
	}
	rti.RecordRequest(syncMsg.RatelimiterId, syncMsg.Timestamp)
}
