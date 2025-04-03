package integration

import (
	"context"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestKafkaService(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %s", err)
	}
	kafkaService := kafka.NewKafkaService(logger)

	t.Run("Should be able to retrieve all topics", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		_, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)

		topicList := []string{"topic1", "topic2", "topic3"}
		err := createTopics(admin, topicList)
		assert.NoError(t, err)
		timeout := 10 * time.Second
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		timeoutTimer := time.NewTimer(timeout)
		defer timeoutTimer.Stop()

		var topics []model.TopicDetails

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
		topicNames := make([]string, len(topics))
		for i, topic := range topics {
			topicNames[i] = topic.Name
		}
		assert.ElementsMatch(t, topicList, topicNames)
		teardown(t, kafkaService, admin, topicList)
	})

	t.Run("Should be able to retrieve all consumer groups", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true
		client, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)
		consumerList := []string{"consumer1", "consumer2", "consumer3"}
		topic := "test-topic"
		err := createTopics(admin, []string{topic})
		assert.NoError(t, err)
		consumerCtx, cancel := context.WithCancel(context.Background())
		cgList := createAndListenConsumerGroups(t, client, consumerCtx, consumerList, topic)
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
		teardown(t, kafkaService, admin, []string{topic})
	})

	t.Run("Should be able to retrieve all consumer groups listening to a topic", func(t *testing.T) {
		assertPrerequisites(t)
		config := sarama.NewConfig()
		config.Version = sarama.V3_6_0_0
		config.Producer.Return.Successes = true
		client, admin := getClientAndAdmin(t, bootstrapAddress, config)
		initializeKafkaService(t, kafkaService, bootstrapAddress, config)
		consumerList := []string{"consumerGroup1", "consumerGroup2", "consumerGroup3"}
		topic := "test-topic2"
		err := createTopics(admin, []string{topic})
		assert.NoError(t, err)
		consumerCtx, cancel := context.WithCancel(context.Background())
		_ = createAndListenConsumerGroups(t, client, consumerCtx, consumerList, topic)

		timeout := 20 * time.Second
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		timeoutTimer := time.NewTimer(timeout)
		defer timeoutTimer.Stop()

		details := []model.ConsumerGroupDetails{}

	WaitLoop:
		for {
			select {
			case <-ticker.C:
				details, err = kafkaService.GetConsumerGroupsDetailsListeningToTopic(topic)

				if len(details) == len(consumerList) &&
					len(details[0].ConsumerDetails) > 0 &&
					len(details[1].ConsumerDetails) > 0 &&
					len(details[2].ConsumerDetails) > 0 {
					break WaitLoop
				}
			case <-timeoutTimer.C:
				t.Fatalf("Timed out waiting for topics to appear")
			}
		}
		cancel()

		actualConsumerGroups := make([]string, len(details))
		offsets := make([][]int64, len(details))
		highWaterMarks := make([][]int64, len(details))

		for i, detail := range details {
			actualConsumerGroups[i] = detail.GroupId
			offsets[i] = make([]int64, len(detail.ConsumerDetails))
			highWaterMarks[i] = make([]int64, len(detail.ConsumerDetails))
			for j, consumer := range detail.ConsumerDetails {
				offsets[i][j] = consumer.LastCommittedOffset
				highWaterMarks[i][j] = consumer.HighWaterMark
			}
		}

		assert.ElementsMatch(t, consumerList, actualConsumerGroups)
		assert.ElementsMatch(t, offsets, [][]int64{{-1}, {-1}, {-1}})
		assert.ElementsMatch(t, highWaterMarks, [][]int64{{0}, {0}, {0}})
		teardown(t, kafkaService, admin, []string{topic})
	})
}

func assertPrerequisites(t *testing.T) {
	if bootstrapAddress == "" {
		t.Fatal("bootstrapAddress is nil")
	}
}

func getClientAndAdmin(
	t *testing.T,
	bootstrapAddress string,
	config *sarama.Config,
) (sarama.Client, sarama.ClusterAdmin) {
	client, err := sarama.NewClient([]string{bootstrapAddress}, config)
	if err != nil {
		t.Fatalf("Failed to create client: %s", err)
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		t.Fatalf("Failed to create cluster admin: %s", err)
	}
	return client, admin
}

func initializeKafkaService(
	t *testing.T,
	kafkaService *kafka.KafkaService,
	bootstrapAddress string,
	config *sarama.Config,
) {
	err := kafkaService.ConnectToCluster([]string{bootstrapAddress}, config)
	if err != nil {
		t.Fatalf("Failed to connect to cluster: %s", err)
	}
}

func createAndListenConsumerGroups(
	t *testing.T,
	client sarama.Client,
	consumerCtx context.Context,
	consumerList []string,
	topic string,
) []sarama.ConsumerGroup {
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
	return cgList
}

func createTopics(admin sarama.ClusterAdmin, topics []string) error {
	for _, topic := range topics {
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func teardown(t *testing.T, kafkaService *kafka.KafkaService, admin sarama.ClusterAdmin, topics []string) {
	for _, topic := range topics {
		err := admin.DeleteTopic(topic)
		if err != nil {
			t.Fatalf("Failed to delete topic: %s", err)
		}
	}
	err := admin.Close()
	if err != nil {
		t.Fatalf("Failed to close admin: %s", err)
	}
	err = kafkaService.Close()
	if err != nil {
		t.Fatalf("Failed to close kafka service: %s", err)
	}
}

type noopConsumer struct{}

func (n *noopConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (n *noopConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (n *noopConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received message: %s", msg.Value)
		sess.MarkMessage(msg, "")
		return nil
	}
	return nil
}
