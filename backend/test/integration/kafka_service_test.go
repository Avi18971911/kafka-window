package integration

import (
	"context"
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
	kafkaService := kafka.NewKafkaService(logger, client, admin)

	t.Run("Should be able to retrieve all topics", func(t *testing.T) {
		if admin == nil {
			t.Fatal("admin is nil")
		}
		if client == nil {
			t.Fatal("client is nil")
		}
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

	t.Run("Should be able to retrieve all consumer groups", func(t *testing.T) {
		if admin == nil {
			t.Fatal("admin is nil")
		}
		if client == nil {
			t.Fatal("client is nil")
		}
		consumerList := []string{"consumer1", "consumer2", "consumer3"}
		topic := "test-topic"
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		assert.NoError(t, err)
		consumerCtx, cancel := context.WithCancel(context.Background())
		cgList := []sarama.ConsumerGroup{}

		for _, consumer := range consumerList {
			cg, err := sarama.NewConsumerGroupFromClient(consumer, client)
			if err != nil {
				t.Fatalf("Failed to create consumer group: %s", err)
			}
			cgList = append(cgList, cg)
			go func(cg sarama.ConsumerGroup) {
				defer cg.Close()
				for {
					err := cg.Consume(consumerCtx, []string{topic}, &noopConsumer{})
					if err != nil {
						return
					}
					if consumerCtx.Err() != nil {
						return
					}
					if consumerCtx.Done() != nil {
						return
					}
				}
			}(cg)
		}
		timeout := 10 * time.Second
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		timeoutTimer := time.NewTimer(timeout)
		defer timeoutTimer.Stop()

		var consumers []string

	WaitLoop:
		for {
			select {
			case <-ticker.C:
				consumers, err = kafkaService.GetConsumerGroups()
				assert.NoError(t, err)

				if len(consumers) == len(consumerList) {
					break WaitLoop
				}
			case <-timeoutTimer.C:
				t.Fatalf("Timed out waiting for consumers to appear")
			}
		}
		assert.ElementsMatch(t, consumerList, consumers)
		cancel()
		for _, cg := range cgList {
			cg.Close()
		}
	})
}

type noopConsumer struct{}

func (n *noopConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (n *noopConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (n *noopConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		sess.MarkMessage(msg, "")
		return nil
	}
	return nil
}
