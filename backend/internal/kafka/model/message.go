package model

import (
	"time"
)

type Message struct {
	Offset           int64       `json:"offset" validate:"required"`
	Partition        int32       `json:"partition" validate:"required"`
	Topic            string      `json:"topic" validate:"required"`
	Timestamp        time.Time   `json:"timestamp" validate:"required"`
	Key              string      `json:"key" validate:"required"`
	KeyJsonPayload   *JSONValue  `json:"keyJsonPayload"`
	KeyPayloadType   PayloadType `json:"keyPayloadType" validate:"required"`
	Value            string      `json:"value" validate:"required"`
	ValueJsonPayload *JSONValue  `json:"valueJsonPayload"`
	ValuePayloadType PayloadType `json:"valuePayloadType" validate:"required"`
}

type PayloadType string

const (
	JSONPayload           PayloadType = "json"
	StringPayload         PayloadType = "string"
	ConsumerOffsetPayload PayloadType = "consumerOffset"
)
