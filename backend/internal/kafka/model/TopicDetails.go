package model

type TopicDetails struct {
	Name              string            `json:"name" validate:"required"`
	NumPartitions     int32             `json:"numPartitions" validate:"required"`
	ReplicationFactor int16             `json:"replicationFactor" validate:"required"`
	IsInternal        bool              `json:"isInternal" validate:"required"`
	CleanupPolicy     CleanupPolicy     `json:"cleanupPolicy" validate:"required"`
	RetentionMs       RetentionMs       `json:"retentionMs" validate:"required"`
	RetentionBytes    int64             `json:"retentionBytes" validate:"required"`
	AdditionalConfigs map[string]string `json:"additionalConfigs" validate:"required"`
}

type CleanupPolicy string

const (
	CleanupPolicyDelete  CleanupPolicy = "delete"
	CleanupPolicyCompact CleanupPolicy = "compact"
	CleanupPolicyBoth    CleanupPolicy = "both"
	CleanupPolicyUnknown CleanupPolicy = "unknown"
)

type RetentionMs struct {
	Indefinite bool  `json:"indefinite" validate:"required"`
	Value      int64 `json:"value" validate:"required"`
}
