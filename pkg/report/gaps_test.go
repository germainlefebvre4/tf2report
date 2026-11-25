package report

import (
	"strings"
	"testing"
)

// PROVISIONAL: pins Known Gap "Markdown and text reports show only a resource's type and name, not its
// module path or index" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when fixed.
func TestKnownGap_ModuleAndIndexNotShown(t *testing.T) {
	plan, summary := loadPlan(t, "modules-and-counts.json")

	formatters := map[string]struct {
		f      Formatter
		golden string
	}{
		"markdown": {&MarkdownFormatter{}, "markdown_modules_and_counts.golden"},
		"text":     {&TextFormatter{}, "text_modules_and_counts.golden"},
	}

	for name, tt := range formatters {
		t.Run(name, func(t *testing.T) {
			got := assertDeterministic(t, tt.f, plan, summary)
			assertGolden(t, tt.golden, got)

			for _, hidden := range []string{"module.net", "module.dr", "[0]", "[1]", `["eu"]`} {
				if strings.Contains(got, hidden) {
					t.Errorf("%s output shows %q; the gap seems fixed, update this test", name, hidden)
				}
			}
		})
	}

	t.Run("identical rows", func(t *testing.T) {
		got := render(t, &MarkdownFormatter{}, plan, summary)
		subnetRow := "| `aws_subnet` | `main` |"
		if n := strings.Count(got, subnetRow); n < 2 {
			t.Errorf("row %q appears %d times, want at least 2 identical rows", subnetRow, n)
		}
		webRow := "| `aws_instance` | `web` |"
		if n := strings.Count(got, webRow); n < 2 {
			t.Errorf("row %q appears %d times, want at least 2 identical rows", webRow, n)
		}
	})
}
