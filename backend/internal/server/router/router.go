package router

import (
	"context"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/Avi18971911/kafka-window/backend/internal/server/handler"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func CreateRouter(
	ctx context.Context,
	kafkaService *kafka.KafkaService,
	logger *zap.Logger,
) http.Handler {
	r := mux.NewRouter()

	r.Handle(
		"/topics", handler.AllTopicsHandler(
			ctx,
			kafkaService,
			logger,
		),
	).Methods("GET")

	return r
}
