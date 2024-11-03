package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type RegistryConfig struct {
    Profile string `yaml:"profile"`
}

type RegistryConfigs struct {
    RegistryConfigs map[string]RegistryConfig `yaml:"registryConfigs"`
}

var RegistryConfigPath = getRegistryConfigPath()
var RegistryConfigFilePath = filepath.Join(RegistryConfigPath, "registryConfig.yaml")

// Function to determine the RegistryConfigPath
func getRegistryConfigPath() string {
	// Get the path from the environment variable
	path := os.Getenv("AWS_ECR_REGISTRY_CONFIG_PATH")
	if path == "" {
		// Default to cache path; expand the tilde to the home directory
		expandedPath, err := os.UserHomeDir()
		if err != nil {
			expandedPath = "." // Fallback to the current directory if home directory cannot be resolved
		}
		return filepath.Join(expandedPath, ".ecr") // Combine with the .ecr folder
	}
	return path
}

// GetRegistryConfig attempts to retrieve a RegistryConfig for the specified registry.
// Returns a pointer to RegistryConfig if found, or nil if not found or if the file doesn't exist.
func getRegistryConfig(registry string) (*RegistryConfig, error) {
    registry = strings.TrimPrefix(registry, "https://")

    // Read the YAML configuration file
    fileData, err := os.ReadFile(RegistryConfigFilePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil // Return nil if file does not exist
        }
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    // Parse the YAML data
    var configs RegistryConfigs
    if err := yaml.Unmarshal(fileData, &configs); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }

    // Look for the registry configuration
    if config, exists := configs.RegistryConfigs[registry]; exists {
        return &config, nil // Return pointer to found config
    }

    return nil, nil // Return nil if registry is not found
}

// GetRegistryProfile attempts to retrieve a profile from the RegistryConfig for the specified registry.
// Returns a the profile string if found, or nil if not found.
func GetRegistryProfile(registry string) (string, error) {
    config, err := getRegistryConfig(registry)
    if err != nil {
        return "", err
    }

    if config == nil {
        return "", nil
    }

    return config.Profile, nil
}