package config

// Build a RegistryConfigs instance with a fluent interface.
type RegistryConfigBuilder struct {
    config *RegistryConfigs
}

// Initializes a new builder instance for RegistryConfigs.
func NewRegistryConfigBuilder() *RegistryConfigBuilder {
    return &RegistryConfigBuilder{
        config: &RegistryConfigs{
            RegistryConfigs: []RegistryConfigEntry{}, // Initialize the slice
        },
    }
}

// Adds a new registry configuration with the given name and credential.
func (b *RegistryConfigBuilder) AddRegistryConfigWithProfile(pattern string, profile string) *RegistryConfigBuilder {
    entry := RegistryConfigEntry{
        Pattern: pattern,
        Config:  RegistryConfig{Profile: profile},
    }
    b.config.RegistryConfigs = append(b.config.RegistryConfigs, entry) // Append to the slice
    return b
}

// Build finalizes and returns the configured RegistryConfigs instance.
func (b *RegistryConfigBuilder) Build() *RegistryConfigs {
    return b.config
}