package kafka

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func (k *KafkaService) FetchLastMessages(
	topic string,
	partition int32,
	limit int,
	encoding Encoding,
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

	startOffset := newestOffset - int64(limit)
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
	messages := make([]*model.Message, limit)
	for message := range partitionConsumer.Messages() {
		timestamp := message.Timestamp
		decodedPayload, err := DecodeMessage(message.Value, encoding)
		if err != nil {
			k.logger.Error(
				"failed to decode message payload, skipping...",
				zap.String("encoding", string(encoding)),
				zap.Error(err),
			)
			continue
		}
		decodedKey, err := DecodeMessage(message.Key, encoding)
		if err != nil {
			k.logger.Error(
				"failed to decode message key, skipping...",
				zap.String("encoding", string(encoding)),
				zap.Error(err),
			)
			continue
		}
		var decodedKeyJSONPayload *model.JSONValue = nil
		if decodedKey.Type == model.JSONPayload {
			decodedKeyJSONPayload = &decodedKey.JSONPayload
		}
		var decodedValueJSONPayload *model.JSONValue = nil
		if decodedPayload.Type == model.JSONPayload {
			decodedValueJSONPayload = &decodedPayload.JSONPayload
		}
		messages[i] = &model.Message{
			Topic:            topic,
			Offset:           message.Offset,
			Partition:        message.Partition,
			Key:              decodedKey.Payload,
			KeyPayloadType:   decodedKey.Type,
			KeyJsonPayload:   decodedKeyJSONPayload,
			Value:            decodedPayload.Payload,
			ValuePayloadType: decodedPayload.Type,
			ValueJsonPayload: decodedValueJSONPayload,
			Timestamp:        timestamp,
		}
		i++
		if i >= limit {
			break
		}
	}

	return messages, nil
}
