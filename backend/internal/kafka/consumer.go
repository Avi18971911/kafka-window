package kafka

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func (k *KafkaService) GetLastMessagesForTopic(
	topic string,
	partitionData model.PartitionInput,
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
	for _, partition := range topicDetail.Partitions {
		partitionDetails, ok := partitionData.PartitionDetailsMap[partition.ID]
		if !ok {
			continue
		}

		partitionMessages, err := k.getMessagesForPartition(
			topic,
			partition.ID,
			partitionDetails.StartOffset,
			partitionDetails.EndOffset,
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

func (k *KafkaService) getMessagesForPartition(
	topic string,
	partition int32,
	startOffset int64,
	endOffset int64,
) ([]*model.Message, error) {
	if startOffset < 0 || endOffset < 0 {
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
		if startOffset < 0 {
			startOffset = max(1, newestOffset+startOffset)
		}
		if endOffset < 0 {
			endOffset = max(1, newestOffset+endOffset)
		}
		if startOffset > endOffset {
			k.logger.Error(
				"calculated start offset is greater than end offset",
				zap.String("topic", topic),
				zap.Int32("partition", partition),
				zap.Int64("startOffset", startOffset),
				zap.Int64("endOffset", endOffset),
			)
			return nil, fmt.Errorf(
				"calculated start offset %d is greater than end offset %d",
				startOffset,
				endOffset,
			)
		}
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
	numberMessages := int(endOffset - startOffset)
	messages := make([]*model.Message, 0)
	for message := range partitionConsumer.Messages() {
		decodedMessage, err := k.decodeKeyAndValue(message)
		if err != nil {
			k.logger.Error(
				"failed to decode message",
				zap.String("topic", topic),
				zap.Int32("partition", partition),
				zap.Error(err),
			)
		} else {
			messages = append(messages, decodedMessage)
		}
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
