package config

import (
	"reflect"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

func TestNewConfig_defaults(t *testing.T) {
	cfg := NewConfig()

	if cfg.TerraformPlanPath != "terraform.tfplan.json" {
		t.Errorf("TerraformPlanPath = %q", cfg.TerraformPlanPath)
	}
	if cfg.OutputFormat != "markdown" {
		t.Errorf("OutputFormat = %q", cfg.OutputFormat)
	}
	if cfg.Verbosity != "info" {
		t.Errorf("Verbosity = %q", cfg.Verbosity)
	}
	if cfg.Filters.ResourceTypes == nil || len(cfg.Filters.ResourceTypes) != 0 {
		t.Errorf("ResourceTypes = %#v, want a non-nil empty slice", cfg.Filters.ResourceTypes)
	}
	if cfg.Filters.Actions == nil || len(cfg.Filters.Actions) != 0 {
		t.Errorf("Actions = %#v, want a non-nil empty slice", cfg.Filters.Actions)
	}
}

func TestFilters_GetActionTypes(t *testing.T) {
	tests := []struct {
		name    string
		actions []string
		want    []terraform.ActionType
	}{
		{"nil", nil, []terraform.ActionType{}},
		{"empty", []string{}, []terraform.ActionType{}},
		{"known actions", []string{"create", "delete"}, []terraform.ActionType{terraform.ActionCreate, terraform.ActionDelete}},
		{"unknown value is passed through", []string{"bogus"}, []terraform.ActionType{terraform.ActionType("bogus")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filters{Actions: tt.actions}
			got := f.GetActionTypes()
			if got == nil {
				t.Fatalf("GetActionTypes returned nil")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActionTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
