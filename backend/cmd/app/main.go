package app

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	// TODO: Get broker URL from config
	brokerUrl := "localhost:9092"
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	broker := sarama.NewBroker(brokerUrl)

	config := sarama.NewConfig()
	config.ClientID = "kafka-ui"
	config.Version = sarama.V2_5_0_0

	err = broker.Open(config)
	if err != nil {
		logger.Error("failed to open broker", zap.Error(err))
	}

	defer logger.Sync()

}
