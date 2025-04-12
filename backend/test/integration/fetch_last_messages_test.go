package integration

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

func TestFetchLastMessages(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %s", err)
	}
	kafkaService := kafka.NewKafkaService(logger)

	t.Run("Should be able to fetch last 100 plaintext encoded messages", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true

		client, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)

		topic := "test-topic-fetch-plaintext"
		numPartitions := int32(1)
		replicationFactor := int16(1)

		err := createTopic(admin, topic, numPartitions, replicationFactor)
		assert.NoError(t, err)
		plaintextMessages, err := createInitialMessages(topic, 0, kafka.PlainText, 100)
		assert.NoError(t, err)
		err = produceMessages(client, plaintextMessages)
		assert.NoError(t, err)

		messages, err := kafkaService.FetchLastMessages(
			topic,
			0,
			100,
			100,
			kafka.PlainText,
			kafka.PlainText,
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 100)
		for i, message := range messages {
			plainTextValue, err := plaintextMessages[i].Value.Encode()
			assert.NoError(t, err)
			plainTextKey, err := plaintextMessages[i].Key.Encode()
			assert.NoError(t, err)

			assert.Equal(t, string(plainTextValue), message.Value)
			assert.Equal(t, string(plainTextKey), message.Key)
			assert.Equal(t, plaintextMessages[i].Partition, message.Partition)
			assert.Equal(t, plaintextMessages[i].Topic, message.Topic)
			assert.Equal(t, int64(i), message.Offset)
			assert.Equal(t, model.StringPayload, message.KeyPayloadType)
			assert.Equal(t, model.StringPayload, message.ValuePayloadType)
			assert.Nil(t, message.KeyJsonPayload)
			assert.Nil(t, message.ValueJsonPayload)
		}
		teardown(t, kafkaService, admin, []string{topic})
	})

	t.Run("Should be able to fetch last 100 base64 encoded messages", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true

		client, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)

		topic := "test-topic-fetch-base64"
		numPartitions := int32(1)
		replicationFactor := int16(1)

		err := createTopic(admin, topic, numPartitions, replicationFactor)
		assert.NoError(t, err)
		base64Messages, err := createInitialMessages(topic, 0, kafka.Base64, 100)
		assert.NoError(t, err)
		err = produceMessages(client, base64Messages)
		assert.NoError(t, err)

		messages, err := kafkaService.FetchLastMessages(
			topic,
			0,
			100,
			100,
			kafka.Base64,
			kafka.Base64,
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 100)
		for i, message := range messages {
			encodedVal, err := base64Messages[i].Value.Encode()
			assert.NoError(t, err)
			base64Value, err := base64.StdEncoding.DecodeString(string(encodedVal))
			assert.NoError(t, err)
			encodedKey, err := base64Messages[i].Key.Encode()
			assert.NoError(t, err)
			base64Key, err := base64.StdEncoding.DecodeString(string(encodedKey))
			assert.NoError(t, err)

			assert.Equal(t, string(base64Value), message.Value)
			assert.Equal(t, string(base64Key), message.Key)
			assert.Equal(t, base64Messages[i].Partition, message.Partition)
			assert.Equal(t, base64Messages[i].Topic, message.Topic)
			assert.Equal(t, int64(i), message.Offset)
			assert.Equal(t, model.StringPayload, message.KeyPayloadType)
			assert.Equal(t, model.StringPayload, message.ValuePayloadType)
			assert.Nil(t, message.KeyJsonPayload)
			assert.Nil(t, message.ValueJsonPayload)
		}
		teardown(t, kafkaService, admin, []string{topic})
	})

	t.Run("Should be able to fetch last 100 json encoded messages", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true

		client, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)

		topic := "test-topic-fetch-json"
		numPartitions := int32(1)
		replicationFactor := int16(1)

		err := createTopic(admin, topic, numPartitions, replicationFactor)
		assert.NoError(t, err)
		jsonMessages, err := createInitialMessages(topic, 0, kafka.JSON, 100)
		assert.NoError(t, err)
		err = produceMessages(client, jsonMessages)
		assert.NoError(t, err)

		messages, err := kafkaService.FetchLastMessages(
			topic,
			0,
			100,
			100,
			kafka.JSON,
			kafka.JSON,
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 100)
		for i, message := range messages {
			encodedVal, err := jsonMessages[i].Value.Encode()
			assert.NoError(t, err)
			encodedKey, err := jsonMessages[i].Key.Encode()
			assert.NoError(t, err)

			assert.Equal(t, string(encodedVal), message.Value)
			assert.Equal(t, string(encodedKey), message.Key)
			assert.Equal(t, jsonMessages[i].Partition, message.Partition)
			assert.Equal(t, jsonMessages[i].Topic, message.Topic)
			assert.Equal(t, int64(i), message.Offset)

			assert.Equal(t, model.JSONPayload, message.KeyPayloadType)
			assert.Equal(t, model.JSONPayload, message.ValuePayloadType)

			assert.NotNil(t, message.KeyJsonPayload)
			assert.NotNil(t, message.KeyJsonPayload.ObjectVal["key"])
			assert.NotNil(t, message.ValueJsonPayload)
			assert.NotNil(t, message.ValueJsonPayload.ObjectVal["value"])

			assert.Contains(t, string(encodedKey), *message.KeyJsonPayload.ObjectVal["key"].StringVal)
			assert.Contains(t, string(encodedVal), *message.ValueJsonPayload.ObjectVal["value"].StringVal)
		}
		teardown(t, kafkaService, admin, []string{topic})
	})
}

func createTopic(
	admin sarama.ClusterAdmin,
	topic string,
	numPartition int32,
	replicationFactor int16,
) error {
	err := admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     numPartition,
		ReplicationFactor: replicationFactor,
	}, false)
	if err != nil {
		return err
	}
	return nil
}

func createInitialMessages(
	topic string,
	partition int32,
	encoding kafka.Encoding,
	numMessages int,
) ([]*sarama.ProducerMessage, error) {
	messages := make([]*sarama.ProducerMessage, numMessages)
	for i := 0; i < numMessages; i++ {
		timestamp := time.Now()
		messageString := "Test message " + strconv.FormatInt(int64(i), 10)
		var encodedKey, encodedValue []byte
		var err error
		switch encoding {
		case kafka.JSON:
			jsonMessage := ValueJSON{
				Value: messageString,
			}
			jsonKey := KeyJSON{
				Key: messageString,
			}
			encodedKey, err = json.Marshal(jsonKey)
			if err != nil {
				return nil, err
			}
			encodedValue, err = json.Marshal(jsonMessage)
			if err != nil {
				return nil, err
			}
		case kafka.PlainText:
			encodedKey = []byte(messageString)
			encodedValue = []byte(messageString)
		case kafka.Base64:
			encodedKey = []byte(base64.StdEncoding.EncodeToString([]byte(messageString)))
			encodedValue = []byte(base64.StdEncoding.EncodeToString([]byte(messageString)))
		}

		message := &sarama.ProducerMessage{
			Topic:     topic,
			Partition: partition,
			Key:       sarama.ByteEncoder(encodedKey),
			Value:     sarama.ByteEncoder(encodedValue),
			Timestamp: timestamp,
		}
		messages[i] = message
	}
	return messages, nil
}

func produceMessages(
	client sarama.Client,
	messages []*sarama.ProducerMessage,
) error {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return err
	}
	defer producer.Close()

	for _, message := range messages {
		_, _, err := producer.SendMessage(message)
		if err != nil {
			return err
		}
	}
	return nil
}

type KeyJSON struct {
	Key string `json:"key"`
}

type ValueJSON struct {
	Value string `json:"value"`
}
