package config

import "github.com/germainlefebvre4/tf2report/pkg/terraform"

// Config represents the application configuration.
type Config struct {
	TerraformPlanPath string  `mapstructure:"terraform_plan_path"`
	OutputFormat      string  `mapstructure:"output_format"`
	Filters           Filters `mapstructure:"filters"`
	Verbosity         string  `mapstructure:"verbosity"`
}

// Filters contains filtering options for resources.
type Filters struct {
	ResourceTypes []string `mapstructure:"resource_types"`
	Actions       []string `mapstructure:"actions"`
}

// NewConfig creates a Config with default values.
func NewConfig() *Config {
	return &Config{
		TerraformPlanPath: "terraform.tfplan.json",
		OutputFormat:      "markdown",
		Verbosity:         "info",
		Filters: Filters{
			ResourceTypes: []string{},
			Actions:       []string{},
		},
	}
}

// GetActionTypes converts string actions to ActionType slice.
func (f *Filters) GetActionTypes() []terraform.ActionType {
	if len(f.Actions) == 0 {
		return []terraform.ActionType{}
	}

	actions := make([]terraform.ActionType, 0, len(f.Actions))
	for _, a := range f.Actions {
		actions = append(actions, terraform.ActionType(a))
	}
	return actions
}
