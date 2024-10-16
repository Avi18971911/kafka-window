package service

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaService struct {
	broker *sarama.Broker
	logger *zap.Logger
}

func NewKafkaService(broker *sarama.Broker, logger *zap.Logger) *KafkaService {
	return &KafkaService{
		broker: broker,
		logger: logger,
	}
}

func (k *KafkaService) Close() error {
	err := k.broker.Close()
	if err != nil {
		k.logger.Error("failed to close broker", zap.Error(err))
	}
	return err
}
