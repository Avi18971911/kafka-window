package avro

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AvroService struct {
	Config *Config
}

func NewAvroService(config *Config) *AvroService {
	return &AvroService{
		Config: config,
	}
}

func (s *AvroService) GetSchema(schemaId int) (string, error) {
	for _, url := range s.Config.SchemaRegistryURLs {
		fullURL := fmt.Sprintf("%s/schemas/ids/%d", url, schemaId)
		schema, err := fetchSchemaFromRegistry(fullURL)
		if err == nil {
			return schema, nil
		}
	}
	err := fmt.Errorf("failed to fetch schema with ID %d from all configured schema registries", schemaId)
	return "", err
}

func fetchSchemaFromRegistry(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema from registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch schema from registry: %s", resp.Status)
	}

	var payload struct {
		Schema string `json:"schema"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("failed to decode schema response: %w", err)
	}
	return payload.Schema, nil
}
