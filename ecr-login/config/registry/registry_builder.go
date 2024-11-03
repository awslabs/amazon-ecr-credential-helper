package config

// Build a RegistryConfigs instance with a fluent interface.
type RegistryConfigBuilder struct {
    config *RegistryConfigs
}

// Initializes a new builder instance for RegistryConfigs.
func NewRegistryConfigBuilder() *RegistryConfigBuilder {
    return &RegistryConfigBuilder{
        config: &RegistryConfigs{
            RegistryConfigs: make(map[string]RegistryConfig),
        },
    }
}

// Adds a new registry configuration with the given name and credential.
func (b *RegistryConfigBuilder) AddRegistryWithProfile(registry string, profile string) *RegistryConfigBuilder {
    b.config.RegistryConfigs[registry] = RegistryConfig{Profile: profile}
    return b
}

// Build finalizes and returns the configured RegistryConfigs instance.
func (b *RegistryConfigBuilder) Build() *RegistryConfigs {
    return b.config
}