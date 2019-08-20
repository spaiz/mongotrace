package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Hosts []*HostConfig `json:"hosts"`
}

type HostConfig struct {
	ConnectionString string `json:"connectionString"`
	DBName string `json:"dbName"`
}

func LoadConfig(filePath string) (*Configuration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
