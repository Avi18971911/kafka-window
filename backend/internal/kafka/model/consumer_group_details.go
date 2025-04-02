package model

type ConsumerGroupDetails struct {
	GroupId         string
	ConsumerDetails []ConsumerDetails
}

type ConsumerDetails struct {
	MemberId            string
	LastCommittedOffset int64
	HighWaterMark       int64
}
