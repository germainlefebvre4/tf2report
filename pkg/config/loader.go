package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Loader handles loading configuration from files and environment.
type Loader struct {
	v *viper.Viper
}

// NewLoader creates a new configuration loader.
func NewLoader() *Loader {
	return &Loader{
		v: viper.New(),
	}
}

// Load reads configuration from a file and environment variables.
func (l *Loader) Load(configPath string) (*Config, error) {
	cfg := NewConfig()

	// Set defaults
	l.setDefaults()

	// Configure viper
	l.v.SetConfigType("yaml")
	l.v.SetEnvPrefix("TF2REPORT")
	l.v.AutomaticEnv()

	// Load config file if specified
	if configPath != "" {
		if err := l.loadConfigFile(configPath); err != nil {
			return nil, err
		}
	} else {
		// Try to find config in common locations
		l.v.SetConfigName("tf2report")
		l.v.AddConfigPath(".")
		l.v.AddConfigPath("$HOME/.config/tf2report")

		if err := l.v.ReadInConfig(); err != nil {
			// Config file not required, continue with defaults
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// Unmarshal configuration
	if err := l.v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func (l *Loader) loadConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if err := l.v.ReadConfig(file); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	return nil
}

func (l *Loader) setDefaults() {
	l.v.SetDefault("terraform_plan_path", "terraform.tfplan.json")
	l.v.SetDefault("output_format", "markdown")
	l.v.SetDefault("verbosity", "info")
	l.v.SetDefault("filters.resource_types", []string{})
	l.v.SetDefault("filters.actions", []string{})
}
