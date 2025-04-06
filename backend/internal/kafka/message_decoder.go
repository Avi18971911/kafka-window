package kafka

import (
	"encoding/base64"
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
)

type Encoding string

const (
	JSON      Encoding = "json"
	PlainText Encoding = "plainText"
	Base64    Encoding = "base64"
)

type DecodedPayload struct {
	Payload string
	Type    model.PayloadType
}

func DecodeMessage(value []byte, encoding Encoding) (DecodedPayload, error) {
	switch encoding {
	case JSON:
		decodedResult, err := decodeJSON(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode JSON: %w", err)
		}
		return DecodedPayload{
			Payload: decodedResult,
			Type:    model.JSONPayload,
		}, nil
	case PlainText:
		decodedResult, err := decodePlainText(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode plain text: %w", err)
		}
		return DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	case Base64:
		decodedResult, err := decodeBase64(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode base64: %w", err)
		}
		return DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	default:
		return DecodedPayload{}, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func decodeJSON(value []byte) (string, error) {
	return string(value), nil
}

func decodePlainText(value []byte) (string, error) {
	return string(value), nil
}

func decodeBase64(value []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return string(decoded), nil
}
