package report

import (
	"strings"
	"testing"
)

func TestTextFormatter(t *testing.T) {
	tests := []struct {
		fixture string
		golden  string
	}{
		{"basic.json", "text_basic.golden"},
		{"no-changes.json", "text_no_changes.golden"},
		{"noop-and-read.json", "text_noop_and_read.golden"},
		{"unknown-action.json", "text_unknown_action.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			plan, summary := loadPlan(t, tt.fixture)
			got := assertDeterministic(t, &TextFormatter{}, plan, summary)
			assertGolden(t, tt.golden, got)
		})
	}
}

func TestTextFormatter_hasNoMarkdownSyntax(t *testing.T) {
	plan, summary := loadPlan(t, "basic.json")

	got := render(t, &TextFormatter{}, plan, summary)

	for _, syntax := range []string{"|", "#", "`", "**"} {
		if strings.Contains(got, syntax) {
			t.Errorf("text output contains Markdown syntax %q:\n%s", syntax, got)
		}
	}
}
