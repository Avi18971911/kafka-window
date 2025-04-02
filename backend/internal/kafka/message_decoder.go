package kafka

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Encoding string

const (
	JSON      Encoding = "json"
	PlainText Encoding = "plainText"
	Base64    Encoding = "base64"
)

func DecodeMessage(value []byte, encoding Encoding) (string, error) {
	switch encoding {
	case JSON:
		return decodeJSON(value)
	case PlainText:
		return decodePlainText(value)
	case Base64:
		return decodeBase64(value)
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func decodeJSON(value []byte) (string, error) {
	var message string
	err := json.Unmarshal(value, &message)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}
	return message, nil
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
