package avro

import "fmt"

type Config struct {
	Enabled            bool     `yaml:"enabled"`
	SchemaRegistryURLs []string `yaml:"schemaRegistryUrls"`
	Username           string   `yaml:"username,omitempty"`
	Password           string   `yaml:"password,omitempty"`
}

func (c *Config) Validate() error {
	if c.Enabled == false {
		return nil
	}

	if len(c.SchemaRegistryURLs) == 0 {
		return fmt.Errorf("schema registry is enabled but no URL is configured")
	}

	return nil
}
