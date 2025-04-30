package decoder

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Avi18971911/kafka-window/backend/internal/avro"
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/linkedin/goavro/v2"
	"github.com/twmb/franz-go/pkg/kmsg"
	"unicode"
	"unicode/utf8"
)

type Encoding string

const consumerOffsetEncoding = "__consumer_offsets"

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

type DecodedKeyAndValue struct {
	Key   *DecodedPayload
	Value *DecodedPayload
}

type MessageDecoder struct {
	avroService *avro.AvroService
}

func NewMessageDecoder(avroService *avro.AvroService) *MessageDecoder {
	return &MessageDecoder{
		avroService: avroService,
	}
}

func (m *MessageDecoder) DecodeKeyAndValue(topic string, keyBytes, valueBytes []byte) (*DecodedKeyAndValue, error) {
	keyEncoding, err := m.getEncodingType(topic, keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get key encoding: %w", err)
	}

	valueEncoding, err := m.getEncodingType(topic, valueBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get value encoding: %w", err)
	}

	if keyEncoding == ConsumerOffset || valueEncoding == ConsumerOffset {
		if keyEncoding != ConsumerOffset || valueEncoding != ConsumerOffset {
			return nil, fmt.Errorf(
				"key and value must both be consumer offset, not one or the other",
			)
		}
		keyAndValue, err := m.decodeConsumerOffset(keyBytes, valueBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decode consumer offset: %w", err)
		}
		return keyAndValue, nil
	} else {
		keyDecoded, err := m.decodeMessage(keyBytes, keyEncoding)
		if err != nil {
			return nil, fmt.Errorf("failed to decode key: %w", err)
		}

		valueDecoded, err := m.decodeMessage(valueBytes, valueEncoding)
		if err != nil {
			return nil, fmt.Errorf("failed to decode value: %w", err)
		}
		return &DecodedKeyAndValue{
			Key:   keyDecoded,
			Value: valueDecoded,
		}, nil
	}
}

func (m *MessageDecoder) decodeMessage(value []byte, encoding Encoding) (*DecodedPayload, error) {
	switch encoding {
	case JSON:
		stringJson, decodedResult, err := m.decodeJSON(value)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}
		return &DecodedPayload{
			Payload:     stringJson,
			JSONPayload: decodedResult,
			Type:        model.JSONPayload,
		}, nil
	case PlainText:
		decodedResult, err := m.decodePlainText(value)
		if err != nil {
			return nil, fmt.Errorf("failed to decode plain text: %w", err)
		}
		return &DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	case Base64:
		decodedResult, err := m.decodeBase64(value)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		return &DecodedPayload{
			Payload: decodedResult,
			Type:    model.StringPayload,
		}, nil
	case Avro:
		stringResult, decodedResult, err := m.decodeAvro(value)
		if err != nil {
			return nil, fmt.Errorf("failed to decode avro: %w", err)
		}
		return &DecodedPayload{
			Payload:     stringResult,
			JSONPayload: decodedResult,
			Type:        model.JSONPayload,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func (m *MessageDecoder) decodeJSON(value []byte) (string, model.JSONValue, error) {
	jsonString := string(value)
	parsedJSON, err := parseString(jsonString)
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

	parsedJSON, err := parseString(jsonString)
	if err != nil {
		return jsonString, model.JSONValue{}, fmt.Errorf("failed to parse decoded Avro JSON: %w", err)
	}

	return jsonString, parsedJSON, nil
}

func (m *MessageDecoder) getEncodingType(topic string, rawMessage []byte) (Encoding, error) {
	if len(rawMessage) == 0 {
		return "", fmt.Errorf("empty message")
	}

	if topic == consumerOffsetEncoding {
		return ConsumerOffset, nil
	}

	if rawMessage[0] == 0 && len(rawMessage) >= 5 {
		return Avro, nil
	}

	trimmed := bytes.TrimSpace(rawMessage)

	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return JSON, nil
	}

	if utf8.Valid(trimmed) && isMostlyPrintable(trimmed) {
		return PlainText, nil
	}

	if isValidBase64(trimmed) {
		return Base64, nil
	}

	return "", fmt.Errorf("unknown encoding type")
}

func isMostlyPrintable(b []byte) bool {
	printableCount := 0
	for _, r := range string(b) {
		if unicode.IsPrint(r) {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(b)) > 0.9
}

func isValidBase64(data []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return false
	}
	// Heuristic: if decoding worked but result is garbage, it's not really base64
	return len(decoded) > 0 && utf8.Valid(decoded)
}

func (m *MessageDecoder) decodeConsumerOffset(keyBytes, valueBytes []byte) (*DecodedKeyAndValue, error) {
	if keyBytes == nil {
		return nil, nil
	}

	keyReader := bytes.NewReader(keyBytes)

	var version int16
	if err := binary.Read(keyReader, binary.BigEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	msg := &DecodedKeyAndValue{}
	msg.Key = &DecodedPayload{}
	msg.Value = &DecodedPayload{}

	switch version {
	case 0, 1: // Offset Commit Key
		offsetCommitKey := kmsg.NewOffsetCommitKey()
		err := offsetCommitKey.ReadFrom(keyBytes)
		if err == nil {
			key, _ := json.Marshal(offsetCommitKey)
			msg.Key.Payload = string(key)
			jsonVal, err := parseString(msg.Key.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to parse consumer offsets key as JSON: %w", err)
			}
			msg.Key.JSONPayload = jsonVal
			msg.Key.Type = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read key:", err)
			return nil, fmt.Errorf("failed to read key: %w", err)
		}

		if valueBytes == nil {
			break
		}
		offsetCommitValue := kmsg.NewOffsetCommitValue()
		err = offsetCommitValue.ReadFrom(valueBytes)
		if err == nil {
			val, _ := json.Marshal(offsetCommitValue)
			msg.Value.Payload = string(val)
			jsonVal, err := parseString(msg.Value.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to parse consumer offsets value as JSON: %w", err)
			}
			msg.Value.JSONPayload = jsonVal
			msg.Value.Type = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read value:", err)
			return nil, fmt.Errorf("failed to read value: %w", err)
		}

	case 2: // Group Metadata Key
		metadataKey := kmsg.NewGroupMetadataKey()
		err := metadataKey.ReadFrom(keyBytes)
		if err == nil {
			key, _ := json.Marshal(metadataKey)
			msg.Key.Payload = string(key)
			jsonVal, err := parseString(msg.Key.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to parse consumer offsets key as JSON: %w", err)
			}
			msg.Key.JSONPayload = jsonVal
			msg.Key.Type = model.ConsumerOffsetPayload
		}

		if valueBytes == nil {
			break
		}
		metadataValue := kmsg.NewGroupMetadataValue()
		err = metadataValue.ReadFrom(valueBytes)
		if err == nil {
			val, _ := json.Marshal(metadataValue)
			msg.Value.Payload = string(val)
			jsonVal, err := parseString(msg.Value.Payload)
			if err != nil {
				return nil, fmt.Errorf("failed to parse consumer offsets value as JSON: %w", err)
			}
			msg.Value.JSONPayload = jsonVal
			msg.Value.Type = model.ConsumerOffsetPayload
		} else {
			fmt.Println("Failed to read value:", err)
			return nil, fmt.Errorf("failed to read value: %w", err)
		}
	default:
		fmt.Println("Unknown key version:", version)
	}

	return msg, nil
}
