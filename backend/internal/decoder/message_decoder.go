package decoder

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/avro"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/linkedin/goavro/v2"
)

type Encoding string

const (
	JSON           Encoding = "json"
	PlainText      Encoding = "plainText"
	Base64         Encoding = "base64"
	Avro           Encoding = "avro"
	ConsumerOffset Encoding = "consumerOffset"
)

type DecodedPayload struct {
	Payload     string
	JSONPayload model.JSONValue
	Type        model.PayloadType
}

type MessageDecoder struct {
	avroService *avro.AvroService
}

func NewMessageDecoder(avroService *avro.AvroService) *MessageDecoder {
	return &MessageDecoder{
		avroService: avroService,
	}
}

func (m *MessageDecoder) DecodeMessage(value []byte, encoding Encoding) (DecodedPayload, error) {
	switch encoding {
	case JSON:
		stringJson, decodedResult, err := m.decodeJSON(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode JSON: %w", err)
		}
		return DecodedPayload{
			Payload:     stringJson,
			JSONPayload: decodedResult,
			Type:        model.JSONPayload,
		}, nil
	case PlainText:
		decodedResult, err := m.decodePlainText(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode plain text: %w", err)
		}
		return DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	case Base64:
		decodedResult, err := m.decodeBase64(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode base64: %w", err)
		}
		return DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	case Avro:
		stringResult, decodedResult, err := m.decodeAvro(value)
		if err != nil {
			return DecodedPayload{}, fmt.Errorf("failed to decode avro: %w", err)
		}
		return DecodedPayload{
			Payload:     stringResult,
			JSONPayload: decodedResult,
			Type:        model.JSONPayload,
		}, nil
	default:
		return DecodedPayload{}, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func (m *MessageDecoder) decodeJSON(value []byte) (string, model.JSONValue, error) {
	jsonString := string(value)
	parsedJSON, err := ParseString(jsonString)
	return jsonString, parsedJSON, err
}

func (m *MessageDecoder) decodePlainText(value []byte) (string, error) {
	return string(value), nil
}

func (m *MessageDecoder) decodeBase64(value []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return string(decoded), nil
}

func (m *MessageDecoder) decodeAvro(value []byte) (string, model.JSONValue, error) {
	if len(value) < 5 || value[0] != 0x00 {
		return "", model.JSONValue{}, fmt.Errorf("invalid Avro message: missing magic byte or too short")
	}

	schemaID := int(binary.BigEndian.Uint32(value[1:5]))
	schema, err := m.avroService.GetSchema(schemaID)
	if err != nil {
		return "", model.JSONValue{}, fmt.Errorf("failed to fetch schema ID %d: %w", schemaID, err)
	}

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return "", model.JSONValue{}, fmt.Errorf("failed to create Avro codec: %w", err)
	}

	native, _, err := codec.NativeFromBinary(value[5:]) // skip magic byte + schema ID
	if err != nil {
		return "", model.JSONValue{}, fmt.Errorf("failed to decode Avro payload: %w", err)
	}

	jsonBytes, err := json.Marshal(native)
	if err != nil {
		return "", model.JSONValue{}, fmt.Errorf("failed to marshal decoded Avro to JSON: %w", err)
	}

	jsonString := string(jsonBytes)

	parsedJSON, err := ParseString(jsonString)
	if err != nil {
		return jsonString, model.JSONValue{}, fmt.Errorf("failed to parse decoded Avro JSON: %w", err)
	}

	return jsonString, parsedJSON, nil
}
