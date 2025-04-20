package kafka

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/decoder"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaService struct {
	client  sarama.Client
	admin   sarama.ClusterAdmin
	decoder *decoder.MessageDecoder
	logger  *zap.Logger
}

func NewKafkaService(
	decoder *decoder.MessageDecoder,
	logger *zap.Logger,
) *KafkaService {
	return &KafkaService{
		decoder: decoder,
		logger:  logger,
	}
}

func (k *KafkaService) ConnectToCluster(brokers []string, config *sarama.Config) error {
	if config == nil {
		config = sarama.NewConfig()
		config.ClientID = "kafka-ui"
		config.Version = sarama.V3_6_0_0
	}
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin: %w", err)
	}
	k.client = client
	k.admin = admin
	return nil
}

func (k *KafkaService) Close() error {
	err := k.admin.Close()
	if err != nil {
		k.logger.Error("KafkaService can't close: failed to close admin and client", zap.Error(err))
		return err
	}
	return nil
}
