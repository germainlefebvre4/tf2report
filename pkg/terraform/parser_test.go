package terraform

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()

	t.Run("basic plan", func(t *testing.T) {
		plan, err := p.Parse(fixture(t, "basic.json"))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if len(plan.ResourceChanges) != 4 {
			t.Errorf("ResourceChanges = %d, want 4", len(plan.ResourceChanges))
		}
		if plan.TerraformVersion != "1.17.2" {
			t.Errorf("TerraformVersion = %q, want 1.17.2", plan.TerraformVersion)
		}
	})

	t.Run("empty object is a valid plan", func(t *testing.T) {
		plan, err := p.Parse(fixture(t, "empty-object.json"))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if len(plan.ResourceChanges) != 0 {
			t.Errorf("ResourceChanges = %d, want 0", len(plan.ResourceChanges))
		}
	})

	t.Run("missing resource_changes", func(t *testing.T) {
		plan, err := p.Parse([]byte(`{"terraform_version": "1.17.2"}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if len(plan.ResourceChanges) != 0 {
			t.Errorf("ResourceChanges = %d, want 0", len(plan.ResourceChanges))
		}
	})

	t.Run("unknown top-level keys are ignored", func(t *testing.T) {
		plan, err := p.Parse([]byte(`{"terraform_version": "1.17.2", "prior_state": {"x": 1}, "whatever": [1,2]}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if plan.TerraformVersion != "1.17.2" {
			t.Errorf("TerraformVersion = %q", plan.TerraformVersion)
		}
	})

	t.Run("output_changes of the basic fixture is parsed", func(t *testing.T) {
		plan, err := p.Parse(fixture(t, "basic.json"))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if _, ok := plan.OutputChanges["endpoint"]; !ok {
			t.Errorf("OutputChanges = %+v, want an endpoint entry", plan.OutputChanges)
		}
	})

	t.Run("output_changes is parsed", func(t *testing.T) {
		plan, err := p.Parse([]byte(`{"output_changes": {"url": {"actions": ["create"]}}}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		oc, ok := plan.OutputChanges["url"]
		if !ok || !reflect.DeepEqual(oc.Actions, []string{"create"}) {
			t.Errorf("OutputChanges = %+v", plan.OutputChanges)
		}
	})
}

func TestParser_ParseFile(t *testing.T) {
	p := NewParser()
	data := fixture(t, "basic.json")

	fromBytes, err := p.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	fromFile, err := p.ParseFile(filepath.Join("..", "..", "testdata", "plans", "basic.json"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if !reflect.DeepEqual(fromBytes, fromFile) {
		t.Errorf("ParseFile result differs from Parse of the same bytes")
	}
}
