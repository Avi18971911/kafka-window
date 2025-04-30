package kafka

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func (k *KafkaService) GetLastMessagesForTopic(
	topic string,
	pageSize int,
	pageNumber int,
) ([]*model.Message, error) {
	topicMetaData, err := k.admin.DescribeTopics([]string{topic})
	if err != nil {
		k.logger.Error(
			"failed to get topic metadata",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return nil, err
	}
	if len(topicMetaData) == 0 {
		k.logger.Warn(
			"topic not found",
			zap.String("topic", topic),
		)
		return nil, nil
	}
	topicDetail := topicMetaData[0]
	if len(topicDetail.Partitions) == 0 {
		k.logger.Warn(
			"topic has no partitions",
			zap.String("topic", topic),
		)
		return nil, nil
	}
	lastMessages := make([]*model.Message, 0)
	// TODO: figure out a proper pagination strategy
	for _, partition := range topicDetail.Partitions {
		partitionMessages, err := k.getLastMessagesForPartition(
			topic,
			partition.ID,
			pageNumber*pageSize,
			pageSize,
		)
		if err != nil {
			k.logger.Error(
				"failed to fetch last messages",
				zap.String("topic", topic),
				zap.Int32("partition", partition.ID),
				zap.Error(err),
			)
			return nil, err
		}
		if partitionMessages != nil && len(partitionMessages) > 0 {
			lastMessages = append(lastMessages, partitionMessages...)
		}
	}
	return lastMessages, nil
}

func (k *KafkaService) getLastMessagesForPartition(
	topic string,
	partition int32,
	distanceFromLatestOffset int,
	numberMessages int,
) ([]*model.Message, error) {
	newestOffset, err := k.client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		k.logger.Error(
			"failed to get newest offset",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get newest offset: %w", err)
	}

	startOffset := newestOffset - int64(distanceFromLatestOffset)
	if newestOffset == 0 {
		return nil, nil
	}
	if startOffset < 0 {
		startOffset = 0
	}

	consumer, err := sarama.NewConsumerFromClient(k.client)
	if err != nil {
		k.logger.Error(
			"failed to create consumer",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, startOffset)
	if err != nil {
		k.logger.Error(
			"failed to create partition consumer",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	i := 0
	messages := make([]*model.Message, numberMessages)
	for message := range partitionConsumer.Messages() {
		decodedMessage, err := k.decodeKeyAndValue(message)
		if err != nil {
			k.logger.Error(
				"failed to decode message",
				zap.String("topic", topic),
				zap.Int32("partition", partition),
				zap.Error(err),
			)
			continue
		}
		messages[i] = decodedMessage
		i++
		if i >= numberMessages {
			break
		}
	}

	return messages, nil
}

func (k *KafkaService) decodeKeyAndValue(
	message *sarama.ConsumerMessage,
) (decodedMessage *model.Message, err error) {
	decodedKeyAndValue, err := k.decoder.DecodeKeyAndValue(message.Topic, message.Key, message.Value)
	if err != nil {
		k.logger.Error(
			"failed to decode key and value",
			zap.String("topic", message.Topic),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode key and value: %w", err)
	}
	var decodedKeyJSONPayload *model.JSONValue = nil
	if decodedKeyAndValue.Key.Type == model.JSONPayload {
		decodedKeyJSONPayload = &decodedKeyAndValue.Key.JSONPayload
	}
	var decodedValueJSONPayload *model.JSONValue = nil
	if decodedKeyAndValue.Value.Type == model.JSONPayload {
		decodedValueJSONPayload = &decodedKeyAndValue.Value.JSONPayload
	}
	return &model.Message{
		Topic:            message.Topic,
		Offset:           message.Offset,
		Partition:        message.Partition,
		Key:              decodedKeyAndValue.Key.Payload,
		KeyPayloadType:   decodedKeyAndValue.Key.Type,
		KeyJsonPayload:   decodedKeyJSONPayload,
		Value:            decodedKeyAndValue.Value.Payload,
		ValuePayloadType: decodedKeyAndValue.Value.Type,
		ValueJsonPayload: decodedValueJSONPayload,
		Timestamp:        message.Timestamp,
	}, nil
}
