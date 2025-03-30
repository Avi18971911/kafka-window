package handler

import (
	"context"
	"encoding/json"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"go.uber.org/zap"
	"net/http"
)

// AllTopicsHandler creates a handler for getting all topics from a list of brokers.
// @Summary Get a list of all topics.
// @Tags analytics
// @Accept json
// @Produce json
// @Param - body string false "No parameters required"
// @Success 200 {array} string "List of topic names"
// @Failure 500 {object} ErrorMessage "Internal server error"
// @Router /topics [get]
func AllTopicsHandler(
	ctx context.Context,
	kafkaService *kafka.KafkaService,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topics, err := kafkaService.GetTopics()
		if err != nil {
			logger.Error("Error encountered when getting all topics", zap.Error(err))
			HttpError(w, "Couldn't query for topics.", http.StatusInternalServerError, logger)
			return
		}
		err = json.NewEncoder(w).Encode(topics)
		if err != nil {
			logger.Error("Error encountered when encoding response", zap.Error(err))
			HttpError(w, "Couldn't encode response.", http.StatusInternalServerError, logger)
		}
	}
}
