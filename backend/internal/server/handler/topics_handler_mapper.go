package handler

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka"
)

func mapEncodingStringToEnum(encoding string) (kafka.Encoding, error) {
	switch encoding {
	case "json":
		return kafka.JSON, nil
	case "plaintext":
		return kafka.PlainText, nil
	case "base64":
		return kafka.Base64, nil
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
