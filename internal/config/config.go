package config

import (
	"fmt"
	"os"

	"github.com/w-h-a/interrogo/api/config/v1alpha1"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*v1alpha1.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg v1alpha1.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config yaml: %w", err)
	}

	// TODO: validate config

	return &cfg, nil
}
