package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type RegistryConfig struct {
    Profile string `yaml:"profile"`
}

type RegistryConfigs struct {
    RegistryConfigs map[string]RegistryConfig `yaml:"registryConfigs"`
}

// TODO: Change this into an env var and default to cache location
var RegistryConfigPath = "..\\..\\..\\registryConfigs.yml"

// FindRegistryConfig attempts to retrieve a RegistryConfig for the specified registry.
// Returns a pointer to RegistryConfig if found, or nil if not found or if the file doesn't exist.
func FindRegistryConfig(registry string) (*RegistryConfig, error) {
    registry = strings.TrimPrefix(registry, "https://")

    // Read the YAML configuration file
    fileData, err := os.ReadFile(RegistryConfigPath)
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
    config, err := FindRegistryConfig(registry)
    if err != nil {
        return "", err
    }

    if config == nil {
        return "", nil
    }

    return config.Profile, nil
}