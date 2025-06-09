package dto

// TopicMessagesInputDTO represents the input data structure for the Kafka events according to a particular topic
// @swagger:model TopicMessagesInputDTO
type TopicMessagesInputDTO struct {
	// The name of the topic to fetch messages from
	TopicName string `json:"topicName" validate:"required"`
	// The Partition request data of the topic to fetch messages from
	Partitions []TopicPartitionInputDTO `json:"partitions" validate:"required"`
}

// TopicPartitionInputDTO represents the partition request data of the topic to fetch messages from
// @swagger:model TopicPartitionInputDTO
type TopicPartitionInputDTO struct {
	// The ID of the partition to fetch messages from
	Partition int32 `json:"partition" validate:"required"`
	// The start offset of the partition to fetch messages from
	// Negative if the offset is to be from the latest offset, so for example -30 means 30 messages from the end
	// -1 means the latest offset
	// Inclusive
	StartOffset int64 `json:"startOffset" validate:"required"`
	// The end offset of the partition to fetch messages from
	// Negative if the offset is to be from the latest offset, so for example -30 means 30 messages from the end
	// -1 means the latest offset
	// Inclusive
	EndOffset int64 `json:"endOffset" validate:"required"`
}
