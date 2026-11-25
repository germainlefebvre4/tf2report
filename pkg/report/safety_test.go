package report

import (
	"strings"
	"testing"
)

func allFormatters(t *testing.T) map[string]Formatter {
	t.Helper()

	formatters := make(map[string]Formatter)
	for _, name := range []string{"markdown", "text", "json"} {
		format, err := ParseFormat(name)
		if err != nil {
			t.Fatalf("ParseFormat(%q): %v", name, err)
		}
		f, err := NewFormatter(format)
		if err != nil {
			t.Fatalf("NewFormatter(%q): %v", name, err)
		}
		formatters[name] = f
	}
	return formatters
}

func TestFormatters_NeverOutputAttributeValues(t *testing.T) {
	tests := []struct {
		fixture string
		markers []string
	}{
		{"sensitive.json", []string{"SYNTHETIC-SECRET-PASSWORD", "SYNTHETIC-SECRET-TOKEN", "SYNTHETIC-SECRET-OLD"}},
		{"basic.json", []string{"SYNTHETIC-OUTPUT-VALUE", "synthetic-data-bucket"}},
	}

	for _, tt := range tests {
		plan, summary := loadPlan(t, tt.fixture)
		for name, f := range allFormatters(t) {
			t.Run(tt.fixture+"/"+name, func(t *testing.T) {
				got := render(t, f, plan, summary)
				for _, marker := range tt.markers {
					if strings.Contains(got, marker) {
						t.Errorf("%s output contains attribute value %q", name, marker)
					}
				}
			})
		}
	}
}
