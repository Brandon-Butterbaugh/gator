package config

import (
	"encoding/json"
	"log"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	// Make a Config
	var cfg Config

	// Get file path
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting file path: %v", err)
		return cfg
	}

	// Read file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file '%s': %v\n", filePath, err)
		return cfg
	}

	// Unmarshal to cfg
	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v\n", err)
		return cfg
	}

	return cfg
}

func (cfg *Config) SetUser(name string) error {
	// Set name in Config
	cfg.CurrentUserName = name

	// Write Config to json
	err := write(*cfg)
	if err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
		return err
	}

	return err
}

func getConfigFilePath() (string, error) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
		return "", err
	}

	// Add working directory and file name to string
	path := home + configFileName

	return path, nil
}

func write(cfg Config) error {
	// Get file path
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Error getting file path: %v", err)
		return err
	}

	// Marshal cfg
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
		return err
	}

	// Write to file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
		return err
	}

	return nil
}
