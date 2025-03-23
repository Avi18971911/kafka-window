package main

import (
	"github.com/Avi18971911/kafka-window/backend/pkg/kafka"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	defer logger.Sync()
	// TODO: Get broker URL from config
	brokers := []string{"localhost:9092"}
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	config := sarama.NewConfig()
	config.ClientID = "kafka-ui"
	config.Version = sarama.V3_6_0_0
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		logger.Fatal("failed to create client", zap.Error(err))
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	kafkaService := kafka.NewKafkaService(logger, client, admin)
	defer kafkaService.Close()
}
