package report

import (
	"strings"
	"testing"
)

// TestGolden_BroadFixtures pins the three report formats for fixtures that widen the dataset beyond
// a single provider and a single format version.
func TestGolden_BroadFixtures(t *testing.T) {
	fixtures := []struct {
		fixture string
		prefix  string
	}{
		{"multi-provider.json", "multi_provider"},
		{"legacy-format.json", "legacy_format"},
	}
	formatters := []struct {
		ext string
		f   Formatter
	}{
		{"markdown", &MarkdownFormatter{}},
		{"text", &TextFormatter{}},
		{"json", &JSONFormatter{}},
	}

	for _, fx := range fixtures {
		plan, summary := loadPlan(t, fx.fixture)
		for _, fm := range formatters {
			t.Run(fx.prefix+"/"+fm.ext, func(t *testing.T) {
				got := assertDeterministic(t, fm.f, plan, summary)
				assertGolden(t, fm.ext+"_"+fx.prefix+".golden", got)
			})
		}
	}
}

// TestMultiProvider_SameNameAcrossTypes checks that resources sharing the name "main" under
// different providers all show up, each under its own resource type.
func TestMultiProvider_SameNameAcrossTypes(t *testing.T) {
	plan, summary := loadPlan(t, "multi-provider.json")
	got := render(t, &MarkdownFormatter{}, plan, summary)

	for _, typ := range []string{
		"aws_instance", "azurerm_resource_group", "azurerm_linux_virtual_machine", "google_compute_instance",
		"google_storage_bucket", "kubernetes_deployment", "datadog_monitor", "cloudflare_record",
	} {
		if !strings.Contains(got, "`"+typ+"`") {
			t.Errorf("markdown report misses resource type %q", typ)
		}
	}
	if strings.Contains(got, "untouched") {
		t.Errorf("markdown report lists the no-op resource")
	}
}
