package config

import (
	"os"
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

func TestGetRegistryProfile_NoRegistryConfiguration(t *testing.T) {
	// Act
    profile, err := GetRegistryProfile("some_registry")

    // Assert
	assert.Equal(t, profile, "")
    assert.Nil(t, err)
}

func TestGetRegistryProfile_NoProfileForRepo(t *testing.T) {
    // Setup
    repository := "some-repository"
    profile := "another-profile"
    
    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryWithProfile(repository, profile).
        Build()

    tempFilePath := createTempYAMLFile(t, testRegistryConfigs)
    RegistryConfigFilePath = tempFilePath

	// Act
    profile, err := GetRegistryProfile("some-other-repository")

    // Assert
    assert.Equal(t, profile, "", "Expected no explicit profile (i.e. \"\") for repository, got profile")
    assert.Nil(t, err, "Expected no error for repository without configuration, got error")
}

func TestGetRegistryProfile_ValidRepo(t *testing.T) {
    // Setup
    repository := "some-repository"
    profile := "another-profile"

    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryWithProfile(repository, profile).
        Build()

    tempFilePath := createTempYAMLFile(t, testRegistryConfigs)
    RegistryConfigFilePath = tempFilePath

	// Act
    resultProfile, err := GetRegistryProfile(repository)

    // Assert
    assert.Equal(t, profile, resultProfile)
    assert.Nil(t, err)
}