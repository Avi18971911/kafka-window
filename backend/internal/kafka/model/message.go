package model

import "time"

type Message struct {
	Offset    int64
	Partition int32
	Topic     string
	Timestamp time.Time
	Key       string
	Value     string
}
