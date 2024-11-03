package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// Helper function to create a temporary YAML file with given content
func createTempYAMLFile(t *testing.T, content interface{}) string {
    t.Helper() // Marks this as a helper function for clearer error reporting

    // Marshal the content to YAML format
    yamlData, err := yaml.Marshal(content)
    if err != nil {
        t.Fatalf("Failed to marshal content to YAML: %v", err)
    }

    // Create a temporary file
    tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }

    // Write YAML data to the temporary file
    if _, err := tmpFile.Write(yamlData); err != nil {
        t.Fatalf("Failed to write to temp file: %v", err)
    }

    // Close the file to ensure data is flushed
    if err := tmpFile.Close(); err != nil {
        t.Fatalf("Failed to close temp file: %v", err)
    }

    // Schedule the file for deletion after the test completes
    t.Cleanup(func() { os.Remove(tmpFile.Name()) })

    return tmpFile.Name()
}

func setupTestConfig(t *testing.T, content interface{}) string {
    t.Helper()
    tempFilePath := createTempYAMLFile(t, content)
    RegistryConfigFilePath = tempFilePath
    return tempFilePath
}

func setupEnvVar(t *testing.T, key, value string) {
    t.Helper()
    original := os.Getenv(key)
    os.Setenv(key, value)
    
    // Restore value after test finishes
    t.Cleanup(func() { os.Setenv(key, original) })
}

func TestGetRegistryConfigPath_NoEnvVar(t *testing.T) {
    homedir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homedir, ".ecr")

    setupEnvVar(t, "AWS_ECR_REGISTRY_CONFIG_PATH", "")

	path := getRegistryConfigPath()

	assert.Equal(t, expectedPath, path, "Expected path to default to the ~/.ecr directory")
}

func TestGetRegistryConfigPath_WithEnvVar(t *testing.T) {
    expectedPath := "/custom/path"
    setupEnvVar(t, "AWS_ECR_REGISTRY_CONFIG_PATH", expectedPath)

	path := getRegistryConfigPath()

	assert.Equal(t, expectedPath, path, "Expected path to match environment variable")
}

func TestGetRegistryConfig_ValidRegistryAndConfig(t *testing.T) {
    registry := "some-registry"
    profile := "another-profile"

    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryWithProfile(registry, profile).
        Build()

    setupTestConfig(t, testRegistryConfigs)

	config, err := getRegistryConfig(registry)

	assert.Equal(t, config.Profile, profile, "Expect returned config to match YAML configuration, got mismatch instead")
    assert.NoError(t, err, "Expected no errors for registry with valid configuration, got error")
}

func TestGetRegistryConfig_InvalidYamlConfig(t *testing.T) {
    invalidYAML := `
		# This is an invalid YAML structure
		registryConfigs:
		  validRegistry:
		    profile: "exampleProfile"
		  invalidRegistry: "missingProfile
	`
    setupTestConfig(t, invalidYAML)

	config, err := getRegistryConfig("validRegistry")

    assert.Nil(t, config, "Expected nil config for invalid YAML format, got config instead")
    assert.Error(t, err, "Expected an error when YAML format is invalid, got no error instead")
}

func TestGetRegistryConfig_ConfigFileNotFound(t *testing.T) {
    RegistryConfigFilePath = "some_invalid_path"

	config, err := getRegistryConfig("someRegistry")

    assert.Nil(t, config, "Expected nil config for missing config file, found config instead")
    assert.NoError(t, err, "Expected no error when the config file cannot be found, got error instead")
}

func TestGetRegistryProfile_NoRegistryConfiguration(t *testing.T) {
    profile, err := GetRegistryProfile("some_registry")

	assert.Equal(t, profile, "", "Expected no profile to be returned when no configuration is set, found profile instead")
    assert.NoError(t, err, "Expected no error when no configuration is set, got error instead")
}

func TestGetRegistryProfile_ValidRegistry(t *testing.T) {
    registry := "some-registry"
    profile := "another-profile"

    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryWithProfile(registry, profile).
        Build()

    setupTestConfig(t, testRegistryConfigs)

    resultProfile, err := GetRegistryProfile(registry)

    assert.Equal(t, profile, resultProfile, "Expect returned profile to match YAML configuration, got mismatch instead")
    assert.NoError(t, err, "Expected not error for a valid registry profile configuration, got error")
}

func TestGetRegistryProfile_NoProfileForRepo(t *testing.T) {
    registry := "some-registry"
    profile := "another-profile"
    
    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryWithProfile(registry, profile).
        Build()

    setupTestConfig(t, testRegistryConfigs)

    profile, err := GetRegistryProfile("some-other-registry")

    assert.Equal(t, profile, "", "Expected no explicit profile (i.e. \"\") for registry, got a profile")
    assert.NoError(t, err, "Expected no error for registry without configuration, got error")
}

func TestGetRegistryProfile_InvalidYaml(t *testing.T) {
    invalidYAML := `
		# This is an invalid YAML structure
		registryConfigs:
		  validRegistry:
		    profile: "exampleProfile"
		  invalidRegistry: "missingProfile
	`

    setupTestConfig(t, invalidYAML)

    resultProfile, err := GetRegistryProfile("validRegistry")

    assert.Equal(t, "", resultProfile, "Expected no profile for invalid YAML, got a profile")
    assert.Error(t, err, "Expected error for invalid YAML format, got no error")
}