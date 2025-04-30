package handler

import (
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/decoder"
)

func mapEncodingStringToEnum(encoding string) (decoder.Encoding, error) {
	switch encoding {
	case "json":
		return decoder.JSON, nil
	case "plaintext":
		return decoder.PlainText, nil
	case "base64":
		return decoder.Base64, nil
	case "avro":
		return decoder.Avro, nil
	case "consumerOffset":
		return decoder.ConsumerOffset, nil
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
