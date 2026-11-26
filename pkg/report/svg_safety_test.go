package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// attributeValues collects every string found under before, after, before_sensitive,
// after_sensitive and after_unknown of every resource change of a fixture.
func attributeValues(t *testing.T, name string) []string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	var doc struct {
		ResourceChanges []struct {
			Change map[string]any `json:"change"`
		} `json:"resource_changes"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("decode fixture %s: %v", name, err)
	}
	var out []string
	for _, rc := range doc.ResourceChanges {
		for _, key := range []string{"before", "after", "before_sensitive", "after_sensitive", "after_unknown"} {
			collectStrings(rc.Change[key], &out)
		}
	}
	return out
}

func collectStrings(v any, out *[]string) {
	switch t := v.(type) {
	case string:
		*out = append(*out, t)
	case map[string]any:
		for _, child := range t {
			collectStrings(child, out)
		}
	case []any:
		for _, child := range t {
			collectStrings(child, out)
		}
	}
}

func TestSVG_NoAttributeValues(t *testing.T) {
	for _, fixture := range []string{"terraform/sensitive.json", "diagram/aws-network.json"} {
		t.Run(fixture, func(t *testing.T) {
			values := attributeValues(t, fixture)
			if len(values) == 0 {
				t.Fatalf("fixture %s has no attribute value to look for", fixture)
			}
			doc := renderSVG(t, loadRealPlan(t, fixture), SVGOptions{})
			for _, v := range values {
				if len(v) < 6 {
					continue // too short to be told apart from markup and labels
				}
				if strings.Contains(doc, v) {
					t.Errorf("output contains attribute value %q", v)
				}
			}
		})
	}
}

func TestSVG_Escaping(t *testing.T) {
	doc := renderSVG(t, hostilePlan(), SVGOptions{})

	assertWellFormed(t, doc)
	for _, raw := range []string{"a<b>", "<x>", "<i>", "aws_in<stance>", "\x07"} {
		if strings.Contains(doc, raw) {
			t.Errorf("output contains unescaped %q", raw)
		}
	}
	if !strings.Contains(doc, "&lt;") || !strings.Contains(doc, "&amp;") {
		t.Errorf("output does not escape markup characters")
	}
}
