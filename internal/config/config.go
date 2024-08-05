package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogLevel    string `json:"log_level"`
	ServerPort  string `json:"server_port"`
	RegistryURL string `json:"registry_url"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
