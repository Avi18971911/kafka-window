package handler

import (
	"context"
	"encoding/json"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// AllTopicsHandler creates a handler for getting all topics from a list of brokers.
// @Summary Get a list of all topics.
// @Tags topics
// @Accept json
// @Produce json
// @Param - body string false "No parameters required"
// @Success 200 {array} model.TopicDetails "List of topic names"
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

// TopicMessagesHandler creates a handler for getting messages from a topic.
// @Summary Get messages from a topic.
// @Tags topics
// @Accept json
// @Produce json
// @Param topic path string true "Topic name"
// @Param keyEncoding path string true "Key encoding (json, plaintext, base64)"
// @Param messageEncoding path string true "Message encoding (json, plaintext, base64)"
// @Param pageSize path string true "Page size"
// @Param pageNumber path string true "Page number"
// @Success 200 {array} model.Message "List of messages"
// @Failure 400 {object} ErrorMessage "Bad request"
// @Failure 500 {object} ErrorMessage "Internal server error"
// @Router /topics/{topic}/messages [get]
func TopicMessagesHandler(
	ctx context.Context,
	kafkaService *kafka.KafkaService,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := mux.Vars(r)["topic"]
		pageSize := r.URL.Query().Get("pageSize")
		pageNumber := r.URL.Query().Get("pageNumber")
		keyEncoding := r.URL.Query().Get("keyEncoding")
		messageEncoding := r.URL.Query().Get("messageEncoding")

		mappedKeyEncoding, err := mapEncodingStringToEnum(keyEncoding)
		if err != nil {
			logger.Error("Error encountered when getting encoding", zap.Error(err))
			HttpError(w, "Couldn't get encoding.", http.StatusBadRequest, logger)
			return
		}
		mappedMessageEncoding, err := mapEncodingStringToEnum(messageEncoding)
		if err != nil {
			logger.Error("Error encountered when getting encoding", zap.Error(err))
			HttpError(w, "Couldn't get encoding.", http.StatusBadRequest, logger)
			return
		}
		intPageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			logger.Error("Error encountered when converting page size to int", zap.Error(err))
			HttpError(w, "Couldn't convert page size to int.", http.StatusBadRequest, logger)
			return
		}
		intPageNumber, err := strconv.Atoi(pageNumber)
		if err != nil {
			logger.Error("Error encountered when converting page number to int", zap.Error(err))
			HttpError(w, "Couldn't convert page number to int.", http.StatusBadRequest, logger)
			return
		}
		messages, err := kafkaService.GetLastMessages(
			topic,
			mappedKeyEncoding,
			mappedMessageEncoding,
			intPageSize,
			intPageNumber,
		)
		if err != nil {
			logger.Error("Error encountered when getting messages", zap.Error(err))
			HttpError(w, "Couldn't get messages.", http.StatusInternalServerError, logger)
			return
		}
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			logger.Error("Error encountered when encoding response", zap.Error(err))
			HttpError(w, "Couldn't encode response.", http.StatusInternalServerError, logger)
			return
		}
	}
}
