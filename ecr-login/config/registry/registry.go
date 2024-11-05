package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type RegistryConfig struct {
    Profile string `yaml:"profile"`
}

type RegistryConfigEntry struct {
    Pattern string         `yaml:"pattern"`
    Config  RegistryConfig `yaml:"config"`
}

type RegistryConfigs struct {
    RegistryConfigs []RegistryConfigEntry `yaml:"registryConfigs"`
}

const ENV_AWS_ECR_REGISTRY_CONFIG_PATH = "AWS_ECR_REGISTRY_CONFIG_PATH"

var (
    RegistryConfigPath = getRegistryConfigPath()
    RegistryConfigFilePath = filepath.Join(RegistryConfigPath, "registryConfig.yaml")
    GetRegistryProfile = getRegistryProfile // Provide override for mocking
)

// Helper to match registry with wildard patterns
func matchesPattern(pattern, registry string) bool {
    if pattern == "*" {
        return true
    }
    if strings.HasPrefix(pattern, "*") && strings.HasSuffix(registry, pattern[1:]) {
        return true
    }
    if strings.HasSuffix(pattern, "*") && strings.HasPrefix(registry, pattern[:len(pattern)-1]) {
        return true
    }
    return pattern == registry // Exact match
}

// Function to determine the RegistryConfigPath
func getRegistryConfigPath() string {
	// Get the path from the environment variable
    if configPath := os.Getenv(ENV_AWS_ECR_REGISTRY_CONFIG_PATH); configPath != "" {
        logrus.WithField(ENV_AWS_ECR_REGISTRY_CONFIG_PATH, configPath).Debug("Using custom registry config path from environment variables.")
		return configPath
	}
	return "~/.ecr"
}

// GetRegistryConfig attempts to retrieve a RegistryConfig for the specified registry.
// Returns a pointer to RegistryConfig if found, or nil if not found or if the file doesn't exist.
func getRegistryConfig(registry string) (*RegistryConfig, error) {
    registry = strings.TrimPrefix(registry, "https://")

    // Read the YAML configuration file
    fileData, err := os.ReadFile(RegistryConfigFilePath)
    if err != nil {
        if os.IsNotExist(err) {
            logrus.WithField(ENV_AWS_ECR_REGISTRY_CONFIG_PATH, RegistryConfigFilePath).
                Debug("No custom registry config file found. Using default credentials.")
            return nil, nil
        }
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    // Parse the YAML data
    var configs RegistryConfigs
    if err := yaml.Unmarshal(fileData, &configs); err != nil {
        logrus.Error("Failed to parse YAML for custom registry config file.")
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }

    // Look for the registry configuration with wildcards support in file order
    for _, entry := range configs.RegistryConfigs {
        if matchesPattern(entry.Pattern, registry) {
            return &entry.Config, nil
        }
    }

    return nil, nil // Return nil if registry is not found
}

// GetRegistryProfile attempts to retrieve a profile from the RegistryConfig for the specified registry.
// Returns a the profile string if found, or nil if not found.
func getRegistryProfile(registry string) (string, error) {
    config, err := getRegistryConfig(registry)
    if err != nil {
        return "", err
    }

    if config == nil {
        return "", nil
    }

    logrus.WithFields(logrus.Fields{
        "Registry": registry,
        "AWS Profile": config.Profile,
      }).Debug("Using explicit AWS Profile for registry.")
    return config.Profile, nil
}