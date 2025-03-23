package containers

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestKafkaContainer(t *testing.T) {
	t.Run("Should be able to create a topic with Kafka container", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes)
		defer cancel()

		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		_, admin, cleanup := CreateKafkaRuntime(ctx, config)
		defer cleanup()

		topic := "test-topic"
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		assert.NoError(t, err)
	})

	t.Run("Should be able to use producers and consumers with Kafka container", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes)
		defer cancel()

		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
		client, admin, cleanup := CreateKafkaRuntime(ctx, config)
		defer cleanup()

		topic := "test-topic"
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		assert.NoError(t, err)

		producer, err := sarama.NewSyncProducerFromClient(client)
		assert.NoError(t, err)
		defer producer.Close()

		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("test-message"),
		})
		assert.NoError(t, err)

		consumer, err := sarama.NewConsumerGroupFromClient("test-group", client)
		assert.NoError(t, err)
		defer consumer.Close()

		outChannel := make(chan string)
		testConsumer := &TestConsumer{out: outChannel}
		consumerCtx, consumerCancel := context.WithTimeout(context.Background(), 10*time.Second)

		go func() {
			for {
				err = consumer.Consume(consumerCtx, []string{topic}, testConsumer)
				if err != nil {
					t.Error("Error consuming message", err)
					return
				}
				if ctx.Err() != nil {
					return
				}
				if ctx.Done() != nil {
					return
				}
			}
		}()

		select {
		case msg := <-outChannel:
			assert.Equal(t, "test-message", msg)
			consumerCancel()
		case <-time.After(5 * time.Second):
			t.Error("Timed out waiting for message")
			consumerCancel()
		}
	})
}

type TestConsumer struct {
	out chan string
}

func (c *TestConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *TestConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (c *TestConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		c.out <- string(msg.Value)
		sess.MarkMessage(msg, "")
		return nil
	}
	return nil
}
