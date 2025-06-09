package model

type PartitionInput struct {
	PartitionDetailsMap map[int32]PartitionDetails
}

type PartitionDetails struct {
	// The start offset of the partition to fetch messages from
	// Negative if the offset is to be from the latest offset, so for example -30 means 30 messages from the end
	// -1 means the latest offset
	// Inclusive
	StartOffset int64
	// The end offset of the partition to fetch messages from
	// Negative if the offset is to be from the latest offset, so for example -30 means 30 messages from the end
	// -1 means the latest offset
	// Inclusive
	EndOffset int64
}
