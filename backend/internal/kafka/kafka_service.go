package kafka

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaService struct {
	client sarama.Client
	admin  sarama.ClusterAdmin
	logger *zap.Logger
}

func NewKafkaService(
	logger *zap.Logger,
) *KafkaService {
	return &KafkaService{
		logger: logger,
	}
}

func (k *KafkaService) ConnectToCluster(brokers []string, config *sarama.Config) error {
	if config == nil {
		config = sarama.NewConfig()
		config.ClientID = "kafka-ui"
		config.Version = sarama.V3_6_0_0
	}
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin: %w", err)
	}
	k.client = client
	k.admin = admin
	return nil
}

func (k *KafkaService) GetTopics() ([]string, error) {
	topicMap, err := k.admin.ListTopics()
	if err != nil {
		k.logger.Error("KafkaService can't get topics: failed to get topics", zap.Error(err))
		return nil, err
	}
	topics := make([]string, 0, len(topicMap))
	for topic, _ := range topicMap {
		topics = append(topics, topic)
	}
	return topics, nil
}

func (k *KafkaService) GetConsumerGroups() ([]string, error) {
	consumerToTypeMap, err := k.admin.ListConsumerGroups()
	if err != nil {
		k.logger.Error("KafkaService can't get consumers: failed to get consumers", zap.Error(err))
		return nil, fmt.Errorf("failed to get consumer groups: %w", err)
	}
	consumers := make([]string, len(consumerToTypeMap))
	i := 0
	for consumer, _ := range consumerToTypeMap {
		consumers[i] = consumer
		i++
	}
	return consumers, nil
}

func (k *KafkaService) GetConsumerGroupsDetailsListeningToTopic(topic string) ([]model.ConsumerGroupDetails, error) {
	consumerGroups, err := k.GetConsumerGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get consumer groups listening to topic %s: %w", topic, err)
	}

	consumerGroupDetailsMap := make(map[string]*model.ConsumerGroupDetails)

	consumerGroupToConsumersToTopicToPartitionsMap, err :=
		k.getConsumerGroupToConsumersToTopicsToPartitionsMap(consumerGroups, map[string]bool{topic: true})
	if err != nil {
		k.logger.Error(
			"%s failed to get consumer group details for topic",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get consumer group details for topic %s: %w", topic, err)
	}

	for consumerGroup, consumersToTopicToPartitionsMap := range consumerGroupToConsumersToTopicToPartitionsMap {
		topicPartitionsMap := k.getTopicPartitionsMap(consumersToTopicToPartitionsMap)
		consumerToTopicToPartitionOffsetsMap, err := k.getConsumerToTopicToPartitionOffsetsMap(
			consumerGroup,
			consumersToTopicToPartitionsMap,
			topicPartitionsMap,
		)
		if err != nil {
			k.logger.Error(
				"failed to get consumer group details",
				zap.String("consumerGroup", consumerGroup),
				zap.Error(err),
			)
			continue
		}

		consumerDetails := make([]model.ConsumerDetails, 0, len(consumerToTopicToPartitionOffsetsMap))
		for consumer, topicToPartitionOffsetsMap := range consumerToTopicToPartitionOffsetsMap {
			for _, partitionOffsetsMap := range topicToPartitionOffsetsMap {
				for _, offsets := range partitionOffsetsMap {
					consumerDetails = append(consumerDetails, model.ConsumerDetails{
						MemberId:            consumer,
						LastCommittedOffset: offsets.LastCommittedOffset,
						HighWaterMark:       offsets.HighWaterMark,
					})
				}
			}
		}
		consumerGroupDetailsMap[consumerGroup] = &model.ConsumerGroupDetails{
			GroupId:         consumerGroup,
			ConsumerDetails: consumerDetails,
		}
	}

	consumerGroupDetails := make([]model.ConsumerGroupDetails, 0, len(consumerGroupDetailsMap))
	for _, details := range consumerGroupDetailsMap {
		consumerGroupDetails = append(consumerGroupDetails, *details)
	}

	return consumerGroupDetails, nil
}

func (k *KafkaService) Close() error {
	err := k.admin.Close()
	if err != nil {
		k.logger.Error("KafkaService can't close: failed to close admin and client", zap.Error(err))
		return err
	}
	return nil
}

func (k *KafkaService) getConsumerGroupToConsumersToTopicsToPartitionsMap(
	consumerGroups []string,
	topicsToInclude map[string]bool,
) (map[string]map[string]map[string]map[int32]bool, error) {
	returnMap := make(map[string]map[string]map[string]map[int32]bool)

	consumerGroupDescription, err := k.admin.DescribeConsumerGroups(consumerGroups)
	if err != nil {
		k.logger.Error("failed to get consumer group descriptions", zap.Error(err))
		return nil, fmt.Errorf("failed to get consumer group descriptions: %w", err)
	}
	for _, description := range consumerGroupDescription {
		if _, exists := returnMap[description.GroupId]; !exists {
			returnMap[description.GroupId] = make(map[string]map[string]map[int32]bool)
		}
		for _, member := range description.Members {
			if _, exists := returnMap[description.GroupId][member.MemberId]; !exists {
				returnMap[description.GroupId][member.MemberId] = make(map[string]map[int32]bool)
			}
			assignment, err := member.GetMemberAssignment()
			if err != nil {
				k.logger.Error("failed to get member metadata", zap.Error(err))
				continue
			}
			if assignment == nil {
				k.logger.Warn(
					"consumer has no member assignment",
					zap.String("consumerGroup", description.GroupId),
					zap.String("consumer", member.MemberId),
				)
				continue
			}
			topicToPartitionMap := assignment.Topics
			if len(topicToPartitionMap) == 0 {
				k.logger.Warn(
					"consumer has no assigned partitions",
					zap.String("consumerGroup", description.GroupId),
					zap.String("consumer", member.MemberId),
				)
				continue
			}
			for topic, partitions := range topicToPartitionMap {
				if topicsToInclude != nil && !topicsToInclude[topic] {
					continue
				}
				for _, partition := range partitions {
					if _, exists := returnMap[description.GroupId][member.MemberId][topic]; !exists {
						returnMap[description.GroupId][member.MemberId][topic] = make(map[int32]bool)
					}
					returnMap[description.GroupId][member.MemberId][topic][partition] = true
				}
			}
		}
	}
	return returnMap, nil
}

func (k *KafkaService) getTopicPartitionsMap(
	consumerToTopicToPartitionMap map[string]map[string]map[int32]bool,
) map[string]map[int32]bool {
	topicPartitionsMap := make(map[string]map[int32]bool)
	for _, topicToPartitionMap := range consumerToTopicToPartitionMap {
		for topic, partitions := range topicToPartitionMap {
			for partition, _ := range partitions {
				if _, exists := topicPartitionsMap[topic]; !exists {
					topicPartitionsMap[topic] = make(map[int32]bool)
				}
				topicPartitionsMap[topic][partition] = true
			}
		}
	}
	return topicPartitionsMap
}

func (k *KafkaService) getConsumerToTopicToPartitionOffsetsMap(
	consumerGroupId string,
	consumerToTopicToPartitionMap map[string]map[string]map[int32]bool,
	topicPartitionMap map[string]map[int32]bool,
) (map[string]map[string]map[int32]model.ConsumerDetails, error) {
	topicToPartitionList := getTopicToPartitionListFromMap(topicPartitionMap)
	committedOffsets, err := k.admin.ListConsumerGroupOffsets(consumerGroupId, topicToPartitionList)
	if err != nil {
		k.logger.Error(
			"failed to fetch committed offsets",
			zap.String("consumerGroup", consumerGroupId),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch committed offsets: %w", err)
	}

	consumerToTopicToPartitionOffsetMap := make(map[string]map[string]map[int32]model.ConsumerDetails)

	for topic, partitions := range topicPartitionMap {
		for partition, _ := range partitions {
			highWatermark, err := k.client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				k.logger.Error(
					"failed to fetch high watermark",
					zap.String("topic", topic),
					zap.Int32("partition", partition),
					zap.Error(err),
				)
				continue
			}

			committedOffset := committedOffsets.Blocks[topic][partition].Offset

			for consumer, topicToPartitionMap := range consumerToTopicToPartitionMap {
				if partitionsMap, exists := topicToPartitionMap[topic]; exists {
					if _, exists := partitionsMap[partition]; exists {
						if _, exists := consumerToTopicToPartitionOffsetMap[topic]; !exists {
							consumerToTopicToPartitionOffsetMap[topic] = make(map[string]map[int32]model.ConsumerDetails)
						}
						if _, exists := consumerToTopicToPartitionOffsetMap[topic][consumer]; !exists {
							consumerToTopicToPartitionOffsetMap[topic][consumer] = make(map[int32]model.ConsumerDetails)
						}
						if _, exists := consumerToTopicToPartitionOffsetMap[topic][consumer][partition]; !exists {
							consumerToTopicToPartitionOffsetMap[topic][consumer][partition] = model.ConsumerDetails{}
						}
						consumerToTopicToPartitionOffsetMap[topic][consumer][partition] = model.ConsumerDetails{
							LastCommittedOffset: committedOffset,
							HighWaterMark:       highWatermark,
						}
					}
				}
			}
		}
	}
	return consumerToTopicToPartitionOffsetMap, nil
}

func getTopicToPartitionListFromMap(topicPartitionMap map[string]map[int32]bool) map[string][]int32 {
	topicToPartitionList := make(map[string][]int32, len(topicPartitionMap))
	for topic, partitions := range topicPartitionMap {
		if _, exists := topicToPartitionList[topic]; !exists {
			topicToPartitionList[topic] = make([]int32, 0, len(partitions))
		}
		for partition, _ := range partitions {
			topicToPartitionList[topic] = append(topicToPartitionList[topic], partition)
		}
	}
	return topicToPartitionList
}
