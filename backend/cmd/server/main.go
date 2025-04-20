package main

import (
	"context"
	"github.com/Avi18971911/kafka-window/backend/internal/avro"
	messageDecoder "github.com/Avi18971911/kafka-window/backend/internal/decoder"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/Avi18971911/kafka-window/backend/internal/server/router"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// @title Kafka Window API
// @version 1.0
// @description This is a monitoring and analytics tool for Kafka.
// termsOfService: http://swagger.io/terms/
// contact:
//   name: API Support
//   url: http://www.swagger.io/support
//   email: support@swagger.io

// license:
//   name: Apache 2.0
//   url: http://www.apache.org/licenses/LICENSE-2.0.html

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
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	avroConfig := avro.NewConfig(true, []string{"http://schema-registry:8081"})
	avroService := avro.NewAvroService(avroConfig)
	decoder := messageDecoder.NewMessageDecoder(avroService)

	kafkaService := kafka.NewKafkaService(decoder, logger)
	err = kafkaService.ConnectToCluster(brokers, config)
	if err != nil {
		logger.Fatal("could not connect to broker", zap.Error(err))
	}
	defer kafkaService.Close()
	r := router.CreateRouter(context.Background(), kafkaService, logger)
	logger.Info("Starting query server at :8085")
	if err := http.ListenAndServe(":8085", r); err != nil {
		logger.Fatal("Failed to serve: %v", zap.Error(err))
	}
}
