package report

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestJSONFormatter(t *testing.T) {
	tests := []struct {
		fixture string
		golden  string
	}{
		{"basic.json", "json_basic.golden"},
		{"no-changes.json", "json_no_changes.golden"},
		{"noop-and-read.json", "json_noop_and_read.golden"},
		{"unknown-action.json", "json_unknown_action.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			plan, summary := loadPlan(t, tt.fixture)
			got := assertDeterministic(t, &JSONFormatter{}, plan, summary)
			assertGolden(t, tt.golden, got)
		})
	}
}

func TestJSONFormatter_structure(t *testing.T) {
	plan, summary := loadPlan(t, "basic.json")
	got := render(t, &JSONFormatter{}, plan, summary)

	var report JSONReport
	if err := json.Unmarshal([]byte(got), &report); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, got)
	}

	if report.TerraformVersion != "1.17.2" {
		t.Errorf("TerraformVersion = %q", report.TerraformVersion)
	}
	wantSummary := JSONSummary{TotalChanges: 4, ToAdd: 1, ToChange: 1, ToDestroy: 1, ToReplace: 1}
	if report.Summary != wantSummary {
		t.Errorf("Summary = %+v, want %+v", report.Summary, wantSummary)
	}

	wantChanges := []JSONResourceChange{
		{Address: "aws_s3_bucket.data", Type: "aws_s3_bucket", Name: "data", Action: "create"},
		{Address: "aws_instance.web", Type: "aws_instance", Name: "web", Action: "update"},
		{Address: "aws_instance.app", Type: "aws_instance", Name: "app", Action: "replace"},
		{Address: "aws_security_group.old", Type: "aws_security_group", Name: "old", Action: "delete"},
	}
	if !reflect.DeepEqual(report.Changes, wantChanges) {
		t.Errorf("Changes = %+v, want plan order %+v", report.Changes, wantChanges)
	}

	instances, ok := report.ByType["aws_instance"]
	if !ok {
		t.Fatalf("ByType has no aws_instance entry: %+v", report.ByType)
	}
	if instances.ToChange != 1 || instances.ToReplace != 1 ||
		!reflect.DeepEqual(instances.Resources, []string{"aws_instance.web", "aws_instance.app"}) {
		t.Errorf("ByType[aws_instance] = %+v", instances)
	}

	if !strings.HasSuffix(got, "}\n") {
		t.Errorf("output does not end with a closing brace and a newline")
	}
	if !strings.Contains(got, "\n  \"summary\": {") {
		t.Errorf("output is not indented with two spaces:\n%s", got)
	}
}

func TestJSONFormatter_withoutChanges(t *testing.T) {
	plan, summary := loadPlan(t, "empty-object.json")
	got := render(t, &JSONFormatter{}, plan, summary)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, present := raw["terraform_version"]; present {
		t.Errorf("terraform_version is present for a plan without version")
	}
	if !bytes.Equal(bytes.TrimSpace(raw["changes"]), []byte("[]")) {
		t.Errorf("changes = %s, want []", raw["changes"])
	}
	if !bytes.Equal(bytes.TrimSpace(raw["by_type"]), []byte("{}")) {
		t.Errorf("by_type = %s, want {}", raw["by_type"])
	}
}
