package integration

import (
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
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
		plaintextMessages := createInitialMessages(topic, 0, kafka.PlainText, 100)
		err = produceMessages(client, plaintextMessages)
		assert.NoError(t, err)

		messages, err := kafkaService.FetchLastMessages(topic, 0, 100, kafka.PlainText)
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
		}
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
) []*sarama.ProducerMessage {
	messages := make([]*sarama.ProducerMessage, numMessages)
	for i := 0; i < numMessages; i++ {
		timestamp := time.Now()
		message := &sarama.ProducerMessage{
			Topic:     topic,
			Partition: partition,
			Key:       sarama.StringEncoder("key"),
			Value:     sarama.StringEncoder("Test message " + strconv.FormatInt(int64(i), 10)),
			Timestamp: timestamp,
		}
		messages[i] = message
	}
	return messages
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
