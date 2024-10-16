package service

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaService struct {
	client sarama.Client
	logger *zap.Logger
}

func NewKafkaService(logger *zap.Logger, client sarama.Client) *KafkaService {
	return &KafkaService{
		logger: logger,
		client: client,
	}
}

func (k *KafkaService) Start() error {
	return nil
}

func (k *KafkaService) GetTopics() ([]string, error) {
	topics, err := k.client.Topics()
	if err != nil {
		k.logger.Error("KafkaService can't get topics: failed to get topics", zap.Error(err))
		return nil, err
	}
	return topics, nil
}

func (k *KafkaService) Close() error {
	err := k.client.Close()
	if err != nil {
		k.logger.Error("KafkaService can't close: failed to close client", zap.Error(err))
		return err
	}
	return nil
}
