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
	KeyJsonPayload   *JSONValue
	KeyPayloadType   PayloadType
	Value            string
	ValueJsonPayload *JSONValue
	ValuePayloadType PayloadType
}

type PayloadType string

const (
	JSONPayload   PayloadType = "json"
	StringPayload PayloadType = "string"
)
