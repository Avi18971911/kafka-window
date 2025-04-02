package kafka

import (
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"strconv"
)

func (k *KafkaService) GetTopics() ([]model.TopicDetails, error) {
	topicMap, err := k.admin.ListTopics()
	if err != nil {
		k.logger.Error("KafkaService can't get topics: failed to get topics", zap.Error(err))
		return nil, err
	}

	topicDetails := k.getTopicDetailsFromTopicMap(topicMap)
	return topicDetails, nil
}

func (k *KafkaService) getTopicDetailsFromTopicMap(
	topicMap map[string]sarama.TopicDetail,
) []model.TopicDetails {
	topicDetails := make([]model.TopicDetails, len(topicMap))
	i := 0
	for topic, topicDetail := range topicMap {
		additionalConfigs := make(map[string]string)
		var cleanupPolicy = model.CleanupPolicyUnknown
		var retentionMs *int64 = nil
		var retentionBytes *int64 = nil
		for configKey, config := range topicDetail.ConfigEntries {
			if config == nil {
				k.logger.Warn(
					"nil config encountering during getTopicDetailsFromTopicMap. Skipping...",
					zap.String("topic", topic),
				)
				continue
			}
			switch configKey {
			case "cleanup.policy":
				cleanupPolicy = getCleanupPolicy(config)
				if cleanupPolicy == model.CleanupPolicyUnknown {
					k.logger.Warn(
						"encountered unknown cleanup policy",
						zap.String("topic", topic),
						zap.String("policy", *config),
					)
				}
			case "retention.ms":
				rMs, err := strconv.ParseInt(*config, 10, 64)
				if err != nil {
					k.logger.Error(
						"error converting retention.ms to number",
						zap.String("topic", topic),
						zap.Error(err),
					)
					continue
				}
				retentionMs = &rMs
			case "retention.bytes":
				rBs, err := strconv.ParseInt(*config, 10, 64)
				if err != nil {
					k.logger.Error(
						"error converting retention.bytes to number",
						zap.String("topic", topic),
						zap.Error(err),
					)
					continue
				}
				retentionBytes = &rBs
			default:
				additionalConfigs[configKey] = *config
			}
		}
		isInternal := isInternalTopic(topic)
		var retentionMsModel *model.RetentionMs = nil
		if retentionMs != nil {
			retentionMsModel = &model.RetentionMs{
				Indefinite: *retentionMs == -1,
				Value:      *retentionMs,
			}
		}
		topicDetails[i] = model.TopicDetails{
			Name:              topic,
			NumPartitions:     topicDetail.NumPartitions,
			ReplicationFactor: topicDetail.ReplicationFactor,
			IsInternal:        isInternal,
			CleanupPolicy:     cleanupPolicy,
			RetentionMs:       retentionMsModel,
			RetentionBytes:    retentionBytes,
			AdditionalConfigs: additionalConfigs,
		}
		i += 1
	}
	return topicDetails
}

func getCleanupPolicy(config *string) model.CleanupPolicy {
	var cleanupPolicy = model.CleanupPolicyUnknown
	if config == nil {
		cleanupPolicy = model.CleanupPolicyUnknown
	} else {
		policy := *config
		switch policy {
		case "delete":
			cleanupPolicy = model.CleanupPolicyDelete
		case "compact":
			cleanupPolicy = model.CleanupPolicyCompact
		case "delete,compact", "compact,delete":
			cleanupPolicy = model.CleanupPolicyBoth
		default:
			cleanupPolicy = model.CleanupPolicyUnknown
		}
	}
	return cleanupPolicy
}

func isInternalTopic(topic string) bool {
	return len(topic) > 2 && topic[:2] == "__"
}
