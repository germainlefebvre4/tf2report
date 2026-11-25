package terraform

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestParser_Parse_errors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"not json", fixture(t, "invalid.json")},
		{"empty input", []byte{}},
		{"truncated json", []byte(`{"resource_changes": [`)},
		{"wrong type for resource_changes", []byte(`{"resource_changes": 3}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := NewParser().Parse(tt.data)
			if err == nil {
				t.Fatalf("Parse succeeded, want an error")
			}
			if plan != nil {
				t.Errorf("plan = %+v, want nil on error", plan)
			}
			if !strings.Contains(err.Error(), "failed to parse plan JSON") {
				t.Errorf("error = %q, want it to contain %q", err, "failed to parse plan JSON")
			}
		})
	}
}

func TestParser_ParseFile_missingFile(t *testing.T) {
	_, err := NewParser().ParseFile("does-not-exist.json")
	if err == nil {
		t.Fatalf("ParseFile succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "failed to read plan file") {
		t.Errorf("error = %q, want it to contain %q", err, "failed to read plan file")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("errors.Is(err, os.ErrNotExist) = false, the error is not wrapped with %%w")
	}
}

// PROVISIONAL: pins Known Gap "A plan in the exact shape Terraform exports is rejected"
// (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed.
func TestKnownGap_RealPlanShapeRejected(t *testing.T) {
	t.Run("boolean sensitivity flag", func(t *testing.T) {
		_, err := NewParser().Parse(fixture(t, "real-shape.json"))
		if err == nil {
			t.Fatalf("Parse succeeded; the gap seems fixed, update this test")
		}
		if !strings.Contains(err.Error(), "cannot unmarshal bool into Go struct field") {
			t.Errorf("error = %q, want a bool unmarshal error", err)
		}
	})

	t.Run("scalar output value", func(t *testing.T) {
		data := []byte(`{"output_changes": {"endpoint": {"actions": ["create"], "after": "x"}}}`)
		_, err := NewParser().Parse(data)
		if err == nil {
			t.Fatalf("Parse succeeded; the gap seems fixed, update this test")
		}
		if !strings.Contains(err.Error(), "cannot unmarshal string into Go struct field") {
			t.Errorf("error = %q, want a string unmarshal error", err)
		}
	})
}
