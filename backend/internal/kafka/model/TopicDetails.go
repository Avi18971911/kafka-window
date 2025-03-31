package model

type TopicDetails struct {
	Name              string            `json:"name" validate:"required"`
	NumPartitions     int32             `json:"numPartitions" validate:"required"`
	ReplicationFactor int16             `json:"replicationFactor"`
	IsInternal        bool              `json:"isInternal"`
	CleanupPolicy     CleanupPolicy     `json:"cleanupPolicy"`
	RetentionMs       RetentionMs       `json:"retentionMs"`
	RetentionBytes    int64             `json:"retentionBytes"`
	AdditionalConfigs map[string]string `json:"additionalConfigs"`
}

type CleanupPolicy string

const (
	CleanupPolicyDelete  CleanupPolicy = "delete"
	CleanupPolicyCompact CleanupPolicy = "compact"
	CleanupPolicyBoth    CleanupPolicy = "both"
	CleanupPolicyUnknown CleanupPolicy = "unknown"
)

type RetentionMs struct {
	Indefinite bool  `json:"indefinite"`
	Value      int64 `json:"value"`
}
