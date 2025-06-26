package fileconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	testEnvConfigFileName = "test_env.cfg.json"
)

type ConfigForAcceptance struct {
	NewDomainId    string `json:"NEW_DOMAIN_ID"`
	UpdateDomainId string `json:"UPDATE_DOMAIN_ID"`
}

// GetTestConfig Get the test configuration
func GetTestConfig() (*ConfigForAcceptance, error) {
	// 1. Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get current working directory: %w", err)
	}

	// 2. the configuration file path
	configFilePath := filepath.Join(workingDir, testEnvConfigFileName)

	// 3. Read the configuration file content (including auto-close file)
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %s: %w", configFilePath, err)
	}

	// 4. Parse JSON configuration
	var config ConfigForAcceptance
	if err := json.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON configuration: %w", err)
	}

	return &config, nil
}
