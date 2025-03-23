package integration

import (
	"github.com/Avi18971911/kafka-window/backend/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestKafkaService(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %s", err)
	}

	t.Run("Should be able to retrieve all topics", func(t *testing.T) {
		if admin == nil {
			t.Fatal("admin is nil")
		}
		if client == nil {
			t.Fatal("client is nil")
		}
		kafkaService := kafka.NewKafkaService(logger, client, admin)
		topicList := []string{"topic1", "topic2", "topic3"}
		for _, topic := range topicList {
			err := admin.CreateTopic(topic, &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 1,
			}, false)
			if err != nil {
				t.Fatalf("Failed to create topic: %s", err)
			}
		}
		timeout := 10 * time.Second
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		timeoutTimer := time.NewTimer(timeout)
		defer timeoutTimer.Stop()

		var topics []string

	WaitLoop:
		for {
			select {
			case <-ticker.C:
				topics, err = kafkaService.GetTopics()
				assert.NoError(t, err)

				if len(topics) == len(topicList) {
					break WaitLoop
				}
			case <-timeoutTimer.C:
				t.Fatalf("Timed out waiting for topics to appear")
			}
		}
		assert.ElementsMatch(t, topicList, topics)
	})

}
