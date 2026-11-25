package report

import (
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

func TestMarkdownFormatter(t *testing.T) {
	tests := []struct {
		fixture string
		golden  string
	}{
		{"basic.json", "markdown_basic.golden"},
		{"no-changes.json", "markdown_no_changes.golden"},
		{"noop-and-read.json", "markdown_noop_and_read.golden"},
		{"unknown-action.json", "markdown_unknown_action.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			plan, summary := loadPlan(t, tt.fixture)
			got := assertDeterministic(t, &MarkdownFormatter{}, plan, summary)
			assertGolden(t, tt.golden, got)
		})
	}
}

func TestMarkdownFormatter_withoutTerraformVersion(t *testing.T) {
	plan, summary := loadPlan(t, "empty-object.json")

	got := render(t, &MarkdownFormatter{}, plan, summary)

	if strings.Contains(got, "**Terraform Version:**") {
		t.Errorf("output contains a Terraform Version line for a plan without version:\n%s", got)
	}
	if !strings.Contains(got, "No changes. Infrastructure is up-to-date.") {
		t.Errorf("output does not state that there are no changes:\n%s", got)
	}
}

func TestMarkdownFormatter_emptyPlanValue(t *testing.T) {
	plan := &terraform.Plan{}

	got := render(t, &MarkdownFormatter{}, plan, terraform.NewSummary(plan.ResourceChanges))

	if !strings.HasPrefix(got, "# Terraform Plan Summary\n") {
		t.Errorf("output does not start with the title:\n%s", got)
	}
}
