package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/Avi18971911/kafka-window/backend/internal/server/dto"
	"go.uber.org/zap"
	"io"
	"net/http"
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
// @Param topicMessagesInput body dto.TopicMessagesInputDTO true "Topic messages input"
// @Success 200 {array} model.Message "List of messages"
// @Failure 400 {object} ErrorMessage "Bad request"
// @Failure 500 {object} ErrorMessage "Internal server error"
// @Router /topics/messages [post]
func TopicMessagesHandler(
	ctx context.Context,
	kafkaService *kafka.KafkaService,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.TopicMessagesInputDTO
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			HttpError(w, "Invalid request payload", http.StatusBadRequest, logger)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Error("Failed to close request body", zap.Error(err))
			}
		}(r.Body)

		if err := validateRequest(&req); err != nil {
			logger.Error("Validation failed for request", zap.Error(err))
			HttpError(w, err.Error(), http.StatusBadRequest, logger)
			return
		}

		partitionModel := mapTopicPartitionInputDtoToModel(req.Partitions)

		messages, err := kafkaService.GetLastMessagesForTopic(ctx, req.TopicName, partitionModel)
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

func validateRequest(req *dto.TopicMessagesInputDTO) error {
	if req.TopicName == "" {
		return errors.New("topic name is required, but was not provided")
	}
	if len(req.Partitions) == 0 {
		return errors.New("at least one partition is required, but none were provided")
	}
	for _, partition := range req.Partitions {
		if partition.Partition < 0 {
			return errors.New("partition ID must be a non-negative integer")
		}
		if partition.StartOffset > partition.EndOffset {
			return errors.New("start offset must be less than or equal to end offset")
		}
	}
	return nil
}

func mapTopicPartitionInputDtoToModel(
	input []dto.TopicPartitionInputDTO,
) model.PartitionInput {
	partitionInput := model.PartitionInput{
		PartitionDetailsMap: make(map[int32]model.PartitionDetails, len(input)),
	}
	for _, partitionInputDto := range input {
		partitionInput.PartitionDetailsMap[partitionInputDto.Partition] = model.PartitionDetails{
			StartOffset: partitionInputDto.StartOffset,
			EndOffset:   partitionInputDto.EndOffset,
		}
	}
	return partitionInput
}
