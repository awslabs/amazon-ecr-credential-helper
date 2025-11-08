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
func TestMatchesPattern(t *testing.T) {
    testCases := []struct {
        name        string
        registry    string
        pattern     string
        shouldMatch bool
    }{
        {
            name: "Exact match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            shouldMatch: true,
        },
        {
            name: "No match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "987654321000.dkr.ecr.ap-southeast-2.amazonaws.com",
            shouldMatch: false,
        },
        {
            name: "Suffix wildcard match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "123456789000.*",
            shouldMatch: true,
        },
        {
            name: "Prefix wildcard match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "*.dkr.ecr.ap-southeast-2.amazonaws.com",
            shouldMatch: true,
        },
        {
            name: "Wildcard",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "*",
            shouldMatch: true,
        },
        {
            name: "Wildcard no match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "*.dkr.ecr.us-east-1.amazonaws.com",
            shouldMatch: false,
        },
        {
            name: "Not supported double wildcard match",
            registry: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            pattern: "*.dkr.ecr.*.amazonaws.com",
            shouldMatch: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            match := matchesPattern(tc.pattern, tc.registry)
            assert.Equal(t, tc.shouldMatch, match)
        })
    }
}

func TestGetRegistryConfigPath_NoEnvVar(t *testing.T) {
	expectedPath := "~/.ecr"

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

func TestGetRegistryConfig_ExactRegistryAndConfig(t *testing.T) {
    registry := "some-registry"
    profile := "another-profile"

    testRegistryConfigs := NewRegistryConfigBuilder().
        AddRegistryConfigWithProfile(registry, profile).
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

func TestGetRegistryProfileWildcards(t *testing.T) {
    testCases := []struct {
        name            string
        registryPattern string
        registryConfigs *RegistryConfigs
        expectedProfile string
    }{
        {
            name: "No config",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: nil,
            expectedProfile: "",
        },
        {
            name: "Exact match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("123456789000.dkr.ecr.ap-southeast-2.amazonaws.com", "production").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "No match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("987654321000.dkr.ecr.ap-southeast-2.amazonaws.com", "production").
                Build(),
            expectedProfile: "",
        },
        {
            name: "Duplicate match use first match",
            registryPattern: "987654321000.us-east-1.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("987654321000.us-east-1.amazonaws.com", "production").
                AddRegistryConfigWithProfile("987654321000.us-east-1.amazonaws.com", "other_profile").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "Unsupported complex wildcard, no match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("*.dkr.ecr.*.amazonaws.com", "production").
                Build(),
            expectedProfile: "",
        },
        {
            name: "Wildcard prefix single match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("*.dkr.ecr.ap-southeast-2.amazonaws.com", "production").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "Wildcard prefix first match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("*.dkr.ecr.ap-southeast-2.amazonaws.com", "production").
                AddRegistryConfigWithProfile("123456789000.dkr.ecr.ap-southeast-2.amazonaws.com", "some_other_profile").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "Wildcard suffix single match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("123456789000.*", "production").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "Wildcard suffix first match",
            registryPattern: "123456789000.dkr.ecr.ap-southeast-2.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("123456789000.*", "production").
                AddRegistryConfigWithProfile("123456789000.dkr.ecr.ap-southeast-2.amazonaws.com", "some_other_profile").
                AddRegistryConfigWithProfile("*.dkr.ecr.ap-southeast-2.amazonaws.com", "yet_another_profile").
                Build(),
            expectedProfile: "production",
        },
        {
            name: "Wildcard fallback match",
            registryPattern: "123456789000.dkr.ecr.us-east-1.amazonaws.com",
            registryConfigs: NewRegistryConfigBuilder().
                AddRegistryConfigWithProfile("987654321000.us-east-1.amazonaws.com", "production").
                AddRegistryConfigWithProfile("*.us-east-1.amazonaws.com", "fallback_profile").
                Build(),
            expectedProfile: "fallback_profile",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            if (tc.registryConfigs != nil && len(tc.registryConfigs.RegistryConfigs) > 0) {
                setupTestConfig(t, tc.registryConfigs)
            }

            resultProfile, err := GetRegistryProfile(tc.registryPattern)
            
            assert.Equal(t, tc.expectedProfile, resultProfile)
            assert.NoError(t, err, "expected no error")
        })
    }
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