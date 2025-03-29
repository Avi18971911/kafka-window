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

var bootstrapAddress string

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMinutes)
	defer cancel()

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	var cleanup func()
	bootstrapAddress, cleanup = containers.CreateKafkaRuntime(ctx)
	code := m.Run()
	cleanup()
	os.Exit(code)
}
