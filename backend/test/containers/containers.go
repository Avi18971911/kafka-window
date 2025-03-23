package containers

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"log"
	"time"
)

const port = "9093"
const ContainerClusterId = "test-cluster"
const timeoutMinutes = time.Minute * 2

func startKafkaContainer(ctx context.Context) (bootstrapAddress string, stopContainer func(), err error) {
	containerCtx, cancel := context.WithTimeout(ctx, timeoutMinutes)
	defer cancel()
	kafkaContainer, err := kafka.Run(
		containerCtx,
		"confluentinc/cp-kafka:7.6.1",
		kafka.WithClusterID(ContainerClusterId),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to start Kafka container: %v", err)
	}

	stopContainer = func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to stop Kafka container: %v", err)
		}
	}

	host, err := kafkaContainer.Host(containerCtx)
	if err != nil {
		return "", stopContainer, fmt.Errorf("failed to get host: %v", err)
	}

	mappedPort, err := kafkaContainer.MappedPort(containerCtx, port)
	if err != nil {
		return "", stopContainer, fmt.Errorf("failed to get mapped port: %v", err)
	}

	bootstrap := fmt.Sprintf("%s:%s", host, mappedPort.Port())
	return bootstrap, stopContainer, nil
}

func CreateKafkaRuntime(
	ctx context.Context,
	config *sarama.Config,
) (client sarama.Client, admin sarama.ClusterAdmin, cleanup func()) {
	bootstrapAddress, stopContainer, err := startKafkaContainer(ctx)
	if err != nil {
		if stopContainer != nil {
			stopContainer()
		}
		log.Fatalf("Failed to start Kafka container: %v", err)
	}

	client, err = sarama.NewClient([]string{bootstrapAddress}, config)
	if err != nil {
		stopContainer()
		log.Fatalf("Failed to create client: %v", err)
	}

	admin, err = sarama.NewClusterAdminFromClient(client)
	if err != nil {
		stopContainer()
		log.Fatalf("Failed to create admin: %v", err)
	}

	cleanup = func() {
		if err := admin.Close(); err != nil {
			log.Fatalf("Failed to close admin: %v", err)
		}
		stopContainer()
	}

	return client, admin, cleanup
}
