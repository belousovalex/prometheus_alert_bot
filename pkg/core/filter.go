package core

import (
	"time"
)

type SameSendFilter struct {
	sentAt time.Time
	dataHash string
	sendAfter time.Duration
}

func MakeSameSendFilter(seconds int64) *SameSendFilter {
	return &SameSendFilter{
		sentAt:time.Now().Add(-1 * time.Hour),
		dataHash:"",
		sendAfter: time.Duration(seconds) * time.Second,
	}
}

func (filter *SameSendFilter) IsNeedSend(dataHash string) bool {
	if filter.dataHash != dataHash {
		return true
	}
	if time.Now().After(filter.sentAt.Add(filter.sendAfter)) {
		return true
	}
	return false
}

func (filter *SameSendFilter) SetSentTime(dataHash string) {
	filter.sentAt = time.Now()
	filter.dataHash = dataHash
}

