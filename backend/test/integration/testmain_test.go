package integration

import (
	"context"
	"github.com/Avi18971911/kafka-window/backend/test/containers"
	"github.com/IBM/sarama"
	"os"
	"testing"
	"time"
)

const timeoutMinutes = time.Minute * 5

var client sarama.Client
var admin sarama.ClusterAdmin

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes)
	defer cancel()

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	var cleanup func()
	client, admin, cleanup = containers.CreateKafkaRuntime(ctx, config)
	code := m.Run()
	cleanup()
	os.Exit(code)
}
