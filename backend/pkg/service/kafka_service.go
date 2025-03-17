package service

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaService struct {
	client sarama.Client
	admin  sarama.ClusterAdmin
	logger *zap.Logger
}

func NewKafkaService(
	logger *zap.Logger,
	client sarama.Client,
	admin sarama.ClusterAdmin,
) *KafkaService {
	return &KafkaService{
		logger: logger,
		client: client,
		admin:  admin,
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

func (k *KafkaService) GetConsumerGroups() ([]string, error) {
	consumerToTypeMap, err := k.admin.ListConsumerGroups()
	if err != nil {
		k.logger.Error("KafkaService can't get consumers: failed to get consumers", zap.Error(err))
		return nil, err
	}
	consumers := make([]string, 0, len(consumerToTypeMap))
	i := 0
	for consumer, _ := range consumerToTypeMap {
		// consumerToType, the _ variable, is not needed
		consumers[i] = consumer
		i++
	}
	return consumers, nil
}

func (k *KafkaService) Close() error {
	err := k.admin.Close()
	if err != nil {
		k.logger.Error("KafkaService can't close: failed to close admin and client", zap.Error(err))
		return err
	}
	return nil
}
