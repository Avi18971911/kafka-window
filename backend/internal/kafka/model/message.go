package model

import (
	"time"
)

type Message struct {
	Offset           int64
	Partition        int32
	Topic            string
	Timestamp        time.Time
	Key              string
	KeyPayloadType   PayloadType
	Value            string
	ValuePayloadType PayloadType
}

type PayloadType string

const (
	JSONPayload   PayloadType = "json"
	StringPayload PayloadType = "string"
)
