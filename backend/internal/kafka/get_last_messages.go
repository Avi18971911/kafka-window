package kafka

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"go.uber.org/zap"
	"unicode"
	"unicode/utf8"
)

func (k *KafkaService) GetLastMessages(
	topic string,
	keyEncoding Encoding,
	messageEncoding Encoding,
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
		partitionMessages, err := k.FetchLastMessages(
			topic,
			partition.ID,
			pageNumber*pageSize,
			pageSize,
			keyEncoding,
			messageEncoding,
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
		lastMessages = append(lastMessages, partitionMessages...)
	}
	return lastMessages, nil
}

func (k *KafkaService) getEncodingType(rawMessage []byte) (Encoding, error) {
	if len(rawMessage) == 0 {
		return "", fmt.Errorf("empty message")
	}

	/*
		if rawMessage[0] == 0 && len(rawMessage) >= 5 {
			return Avro, nil
		}
	*/

	trimmed := bytes.TrimSpace(rawMessage)

	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return JSON, nil
	}

	if utf8.Valid(trimmed) && isMostlyPrintable(trimmed) {
		return PlainText, nil
	}

	if isValidBase64(trimmed) {
		return Base64, nil
	}

	k.logger.Warn("unable to determine encoding type, defaulting to plain text")
	return PlainText, nil
}

func isMostlyPrintable(b []byte) bool {
	printableCount := 0
	for _, r := range string(b) {
		if unicode.IsPrint(r) {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(b)) > 0.9
}

func isValidBase64(data []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return false
	}
	// Heuristic: if decoding worked but result is garbage, it's not really base64
	return len(decoded) > 0 && utf8.Valid(decoded)
}
