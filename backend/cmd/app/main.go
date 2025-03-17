package main

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"kafka-window.com/pkg/service"
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
	config.Version = sarama.V2_5_0_0
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		logger.Fatal("failed to create client", zap.Error(err))
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	kafkaService := service.NewKafkaService(logger, client, admin)
	defer kafkaService.Close()
}
