package terraform

import (
	"encoding/json"
	"reflect"
	"testing"
)

// parseWithoutSensitiveFlags parses a real-shape fixture. Real plans carry boolean before/after_sensitive flags,
// which the parser rejects today (known gap pinned elsewhere), so the flags are removed first.
func parseWithoutSensitiveFlags(t *testing.T, name string) *Plan {
	t.Helper()

	var raw map[string]any
	if err := json.Unmarshal(fixture(t, name), &raw); err != nil {
		t.Fatalf("decode %s: %v", name, err)
	}
	changes, _ := raw["resource_changes"].([]any)
	for _, c := range changes {
		if ch, ok := c.(map[string]any)["change"].(map[string]any); ok {
			delete(ch, "before_sensitive")
			delete(ch, "after_sensitive")
		}
	}
	outputs, _ := raw["output_changes"].(map[string]any)
	for _, o := range outputs {
		if ch, ok := o.(map[string]any); ok {
			delete(ch, "before_sensitive")
			delete(ch, "after_sensitive")
			delete(ch, "after_unknown")
			delete(ch, "before_unknown")
		}
	}
	data, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := NewParser().Parse(data)
	if err != nil {
		t.Fatalf("Parse %s: %v", name, err)
	}
	return plan
}

func TestPlanConfiguration(t *testing.T) {
	tests := []struct {
		fixture string
		wantNil bool
	}{
		{"terraform/nested-modules.json", false},
		{"terraform/mixed-actions.json", false},
		{"legacy-format.json", true},
		{"diagram/no-configuration.json", true},
	}
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			plan := parseWithoutSensitiveFlags(t, tt.fixture)
			if (plan.Configuration == nil) != tt.wantNil {
				t.Errorf("Configuration nil = %v, want %v", plan.Configuration == nil, tt.wantNil)
			}
		})
	}
}

func TestConfigurationNestedModuleCalls(t *testing.T) {
	plan := parseWithoutSensitiveFlags(t, "terraform/nested-modules.json")
	sub, ok := plan.Configuration.RootModule.ModuleCalls["network"].Module.ModuleCalls["subnets"]
	if !ok {
		t.Fatal("module.network.module.subnets not parsed")
	}
	if len(sub.Module.Resources) != 1 || sub.Module.Resources[0].Name != "subnet" {
		t.Errorf("unexpected nested resources: %+v", sub.Module.Resources)
	}
	if got := sub.Module.Resources[0].References(); !reflect.DeepEqual(got, []string{"count.index"}) {
		t.Errorf("References = %v", got)
	}
}

func TestConfigResourceReferences(t *testing.T) {
	tests := []struct {
		name string
		res  ConfigResource
		want []string
	}{
		{"nested blocks and arrays", ConfigResource{Expressions: map[string]any{
			"a": map[string]any{"references": []any{"x.a", "x.a.id"}},
			"b": []any{map[string]any{"inner": map[string]any{"references": []any{"y.b"}}}},
		}}, []string{"x.a", "x.a.id", "y.b"}},
		{"depends_on and duplicates", ConfigResource{
			DependsOn:   []string{"z.c", "x.a"},
			Expressions: map[string]any{"a": map[string]any{"references": []any{"x.a"}}},
		}, []string{"x.a", "z.c"}},
		{"odd shapes", ConfigResource{Expressions: map[string]any{
			"s": "text", "n": 3.0, "nil": nil, "r": map[string]any{"references": "notalist"},
			"l": map[string]any{"references": []any{1.0, nil, "ok.v"}},
		}}, []string{"ok.v"}},
		{"none", ConfigResource{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.References(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("References() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurationLenient(t *testing.T) {
	plan, err := NewParser().Parse([]byte(`{"resource_changes":[],"configuration":{"root_module":"bad"}}`))
	if err != nil {
		t.Fatalf("a malformed configuration must not fail the parse: %v", err)
	}
	if plan.Configuration == nil || len(plan.Configuration.RootModule.Resources) != 0 {
		t.Errorf("unexpected configuration: %+v", plan.Configuration)
	}
}
