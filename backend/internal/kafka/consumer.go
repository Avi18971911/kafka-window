package kafka

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/decoder"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"github.com/twmb/franz-go/pkg/kmsg"
	"go.uber.org/zap"
)

func (k *KafkaService) GetLastMessage(
	topic string,
	partition int32,
) (*sarama.ConsumerMessage, error) {
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
	startOffset := newestOffset - 1
	if startOffset < 0 {
		return nil, nil
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
	message := <-partitionConsumer.Messages()
	if message == nil {
		k.logger.Error(
			"failed to get message",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
		)
		return nil, fmt.Errorf("failed to get message")
	}
	return message, nil
}

func (k *KafkaService) FetchLastMessages(
	topic string,
	partition int32,
	distanceFromLatestOffset int,
	numberMessages int,
	keyEncoding decoder.Encoding,
	messageEncoding decoder.Encoding,
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
		decodedMessage, err := k.decodeKeyAndValue(message, messageEncoding, keyEncoding)
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
	messageEncoding decoder.Encoding,
	keyEncoding decoder.Encoding,
) (decodedMessage *model.Message, err error) {
	if keyEncoding == decoder.ConsumerOffset || messageEncoding == decoder.ConsumerOffset {
		if messageEncoding != decoder.ConsumerOffset || keyEncoding != decoder.ConsumerOffset {
			return nil, fmt.Errorf("both key and value encoding must be ConsumerOffset, not one or the other")
		}
		return k.DecodeConsumerOffset(message.Key, message.Value)
	} else {
		decodedPayload, err := k.decoder.DecodeMessage(message.Value, messageEncoding)
		if err != nil {
			k.logger.Error(
				"failed to decode message payload, skipping...",
				zap.String("encoding", string(messageEncoding)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to decode message payload: %w", err)
		}
		decodedKey, err := k.decoder.DecodeMessage(message.Key, keyEncoding)
		if err != nil {
			k.logger.Error(
				"failed to decode message key, skipping...",
				zap.String("encoding", string(keyEncoding)),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to decode message key: %w", err)
		}
		var decodedKeyJSONPayload *model.JSONValue = nil
		if decodedKey.Type == model.JSONPayload {
			decodedKeyJSONPayload = &decodedKey.JSONPayload
		}
		var decodedValueJSONPayload *model.JSONValue = nil
		if decodedPayload.Type == model.JSONPayload {
			decodedValueJSONPayload = &decodedPayload.JSONPayload
		}
		return &model.Message{
			Topic:            message.Topic,
			Offset:           message.Offset,
			Partition:        message.Partition,
			Key:              decodedKey.Payload,
			KeyPayloadType:   decodedKey.Type,
			KeyJsonPayload:   decodedKeyJSONPayload,
			Value:            decodedPayload.Payload,
			ValuePayloadType: decodedPayload.Type,
			ValueJsonPayload: decodedValueJSONPayload,
			Timestamp:        message.Timestamp,
		}, nil
	}
}

func (k *KafkaService) DecodeConsumerOffset(keyBytes, valueBytes []byte) (*model.Message, error) {
	if keyBytes == nil {
		k.logger.Warn("key is nil, skipping record")
		return nil, nil
	}

	keyReader := bytes.NewReader(keyBytes)

	var version int16
	if err := binary.Read(keyReader, binary.BigEndian, &version); err != nil {
		k.logger.Error("failed to read key version: %w", zap.Error(err))
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	msg := &model.Message{}

	switch version {
	case 0, 1: // Offset Commit Key
		offsetCommitKey := kmsg.NewOffsetCommitKey()
		err := offsetCommitKey.ReadFrom(keyBytes)
		if err == nil {
			key, _ := json.Marshal(offsetCommitKey)
			msg.Key = string(key)
			jsonVal, err := decoder.ParseString(msg.Key)
			if err != nil {
				k.logger.Error("failed to parse consumer offsets key as JSON: %w", zap.Error(err))
				return nil, fmt.Errorf("failed to parse consumer offsets key as JSON: %w", err)
			}
			msg.KeyJsonPayload = &jsonVal
			msg.KeyPayloadType = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read key:", err)
			return nil, fmt.Errorf("failed to read key: %w", err)
		}

		if valueBytes == nil {
			break
		}
		offsetCommitValue := kmsg.NewOffsetCommitValue()
		err = offsetCommitValue.ReadFrom(valueBytes)
		if err == nil {
			val, _ := json.Marshal(offsetCommitValue)
			msg.Value = string(val)
			jsonVal, err := decoder.ParseString(msg.Value)
			if err != nil {
				k.logger.Error("failed to parse consumer offsets value as JSON: %w", zap.Error(err))
				return nil, fmt.Errorf("failed to parse consumer offsets value as JSON: %w", err)
			}
			msg.ValueJsonPayload = &jsonVal
			msg.ValuePayloadType = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read value:", err)
			return nil, fmt.Errorf("failed to read value: %w", err)
		}

	case 2: // Group Metadata Key
		metadataKey := kmsg.NewGroupMetadataKey()
		err := metadataKey.ReadFrom(keyBytes)
		if err == nil {
			key, _ := json.Marshal(metadataKey)
			msg.Key = string(key)
			jsonVal, err := decoder.ParseString(msg.Key)
			if err != nil {
				k.logger.Error("failed to parse consumer offsets key as JSON: %w", zap.Error(err))
				return nil, fmt.Errorf("failed to parse consumer offsets key as JSON: %w", err)
			}
			msg.KeyJsonPayload = &jsonVal
			msg.KeyPayloadType = model.ConsumerOffsetPayload
		}

		if valueBytes == nil {
			break
		}
		metadataValue := kmsg.NewGroupMetadataValue()
		err = metadataValue.ReadFrom(valueBytes)
		if err == nil {
			val, _ := json.Marshal(metadataValue)
			msg.Value = string(val)
			jsonVal, err := decoder.ParseString(msg.Value)
			if err != nil {
				k.logger.Error("failed to parse consumer offsets value as JSON: %w", zap.Error(err))
				return nil, fmt.Errorf("failed to parse consumer offsets value as JSON: %w", err)
			}
			msg.ValueJsonPayload = &jsonVal
			msg.ValuePayloadType = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read value:", err)
			return nil, fmt.Errorf("failed to read value: %w", err)
		}
	default:
		fmt.Println("Unknown key version:", version)
	}

	return msg, nil
}
