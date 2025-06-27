package fileconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	testEnvConfigFileName = "test_env.cfg.json"
)

type ConfigForAcceptance struct {
	NewDomainId    string `json:"new_domain_id"`
	UpdateDomainId string `json:"update_domain_id"`
}

var config ConfigForAcceptance
var once sync.Once

// GetTestConfig Get the test configuration, this method panics on error
func GetTestConfig() *ConfigForAcceptance {
	once.Do(func() {
		// 1. Get the current working directory
		workingDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("unable to get current working directory"))
		}

		// 2. the configuration file path
		configFilePath := filepath.Join(workingDir, testEnvConfigFileName)

		// 3. Read the configuration file content (including auto-close file)
		configContent, err := os.ReadFile(configFilePath)
		if err != nil {
			panic(fmt.Errorf("failed to read configuration file: %s", testEnvConfigFileName))
		}

		// 4. Parse JSON configuration
		var c ConfigForAcceptance
		if err := json.Unmarshal(configContent, &c); err != nil {
			panic(fmt.Errorf("failed to parse JSON configuration: %s", testEnvConfigFileName))
		}
		config = c
	})
	return &config
}
