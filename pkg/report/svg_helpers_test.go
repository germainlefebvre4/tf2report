package report

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// loadRealPlan parses a fixture, dropping the boolean sensitivity flags and the output changes that
// the parser cannot read today (Known Gap of spec 001), so real plans can still be drawn.
func loadRealPlan(t *testing.T, name string) *terraform.Plan {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("decode fixture %s: %v", name, err)
	}
	delete(doc, "output_changes")
	changes, _ := doc["resource_changes"].([]any)
	for _, c := range changes {
		change := c.(map[string]any)["change"].(map[string]any)
		for _, key := range []string{"before_sensitive", "after_sensitive", "after_unknown"} {
			if _, isBool := change[key].(bool); isBool {
				delete(change, key)
			}
		}
	}
	data, err = json.Marshal(doc)
	if err != nil {
		t.Fatalf("encode fixture %s: %v", name, err)
	}
	plan, err := terraform.NewParser().Parse(data)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", name, err)
	}
	return plan
}

// renderSVG renders a plan with the given options.
func renderSVG(t *testing.T, plan *terraform.Plan, opts SVGOptions) string {
	t.Helper()

	return render(t, NewSVGFormatter(opts), plan, terraform.NewSummary(plan.ResourceChanges))
}

// assertWellFormed parses the whole document and returns the attributes of the root element.
func assertWellFormed(t *testing.T, doc string) map[string]string {
	t.Helper()

	dec := xml.NewDecoder(strings.NewReader(doc))
	roots := 0
	depth := 0
	attrs := map[string]string{}
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("document is not well-formed XML: %v", err)
		}
		switch e := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				roots++
				if e.Name.Local != "svg" {
					t.Fatalf("root element is %q, want svg", e.Name.Local)
				}
				for _, a := range e.Attr {
					attrs[a.Name.Local] = a.Value
				}
			}
			depth++
		case xml.EndElement:
			depth--
		}
	}
	if roots != 1 {
		t.Fatalf("document has %d root elements, want 1", roots)
	}
	return attrs
}

// legendOf returns the legend group of the document.
func legendOf(t *testing.T, doc string) string {
	t.Helper()

	i := strings.Index(doc, `<g class="legend">`)
	if i < 0 {
		t.Fatalf("document has no legend")
	}
	rest := doc[i:]
	end := len(rest)
	markers := []string{`<g class="container`, `<g class="edges">`, `<g class="nodes">`, "</svg>", "No resources to display</text>"}
	for _, marker := range markers {
		if j := strings.Index(rest, marker); j >= 0 && j < end {
			end = j
		}
	}
	return rest[:end]
}

// hostilePlan builds a plan whose names carry markup, quotes and a control character.
func hostilePlan() *terraform.Plan {
	name := "a<b>&\"c'\x07d"
	return &terraform.Plan{
		FormatVersion: "1.2",
		ResourceChanges: []terraform.ResourceChange{{
			Address:      `module.m<x>["k&\"1"].aws_instance.` + name,
			Mode:         "managed",
			Type:         "aws_in<stance>",
			Name:         name,
			ProviderName: "registry.terraform.io/hash<i>corp/a&ws",
			Change:       terraform.Change{Actions: []string{"create"}},
		}},
	}
}

// countNodes counts the node groups of a document.
func countNodes(doc string) int {
	return strings.Count(doc, `<g class="node `)
}
