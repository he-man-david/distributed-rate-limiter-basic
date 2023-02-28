package ratetracker

// TODO: Jaser
// here we want a struct with methods and store value/metadata
// maybe we want to know total API_KEYS we are tracking...etc
// This is where our mapping of KEY : ratetrackerInstance will be
// Our Rate module will be calling with this module

type SyncMsg struct {
	RatelimiterId int64 
	ApiKey		  int64
	Timestamp     int64
}

type RateTracker struct {
	InSyncChan chan *SyncMsg
	OutSyncChan chan *SyncMsg
}