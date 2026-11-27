package report

import (
	"reflect"
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/diagram"
	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

func TestSVG_Document(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/aws-network.json")
	doc := renderSVG(t, plan, SVGOptions{})

	attrs := assertWellFormed(t, doc)
	for _, name := range []string{"xmlns", "viewBox", "role"} {
		if attrs[name] == "" {
			t.Errorf("root element has no %s attribute", name)
		}
	}
	if attrs["xmlns"] != "http://www.w3.org/2000/svg" || attrs["role"] != "img" {
		t.Errorf("root attributes = %v", attrs)
	}
	if !strings.Contains(doc, "<title>Terraform plan diagram</title>") {
		t.Errorf("document has no title")
	}
	if !strings.Contains(doc, "<desc>7 resources: create 3, update 1, delete 1, replace 1, no-op 1</desc>") {
		t.Errorf("desc does not summarize the statuses:\n%s", doc)
	}
	for _, addr := range []string{"aws_instance.web", "aws_instance.worker[0]", "aws_s3_bucket.old"} {
		if !strings.Contains(doc, "<title>"+addr+"</title>") {
			t.Errorf("node tooltip for %s is missing", addr)
		}
	}
	if !strings.Contains(doc, `class="node status-noop context"`) {
		t.Errorf("the unchanged vpc referenced by changed resources must be drawn as muted context")
	}
	if strings.Contains(doc, "aws_route_table") {
		t.Errorf("unreferenced unchanged resources must not be drawn")
	}
}

func TestStatusStylesAreDistinct(t *testing.T) {
	seenColor := map[string]string{}
	seenPair := map[[2]string]string{}
	for st := diagram.StatusCreate; st <= diagram.StatusNoop; st++ {
		s := statusStyles[st]
		if other, ok := seenColor[s.stroke]; ok {
			t.Errorf("%s and %s share color %s", s.label, other, s.stroke)
		}
		seenColor[s.stroke] = s.label

		pair := [2]string{s.symbol, s.dash}
		if other, ok := seenPair[pair]; ok {
			t.Errorf("%s and %s share symbol %q and dash %q: not distinguishable in greyscale", s.label, other, s.symbol, s.dash)
		}
		seenPair[pair] = s.label
	}
	if len(seenPair) != 5 {
		t.Errorf("styles = %d, want 5", len(seenPair))
	}
}

func TestSVG_Legend(t *testing.T) {
	tests := []struct {
		fixture string
		present []string
		absent  []string
	}{
		{"diagram/no-configuration.json", []string{"create", "update", "delete"}, []string{"replace", "no-op"}},
		{"diagram/aws-network.json", []string{"create", "update", "delete", "replace", "no-op"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			plan, _ := loadPlan(t, tt.fixture)
			legend := legendOf(t, renderSVG(t, plan, SVGOptions{}))
			for _, label := range tt.present {
				if !strings.Contains(legend, ">"+label+"</text>") {
					t.Errorf("legend misses %q", label)
				}
			}
			for _, label := range tt.absent {
				if strings.Contains(legend, ">"+label+"</text>") {
					t.Errorf("legend lists %q although it is not present", label)
				}
			}
		})
	}
}

func TestSVG_EmptySelection(t *testing.T) {
	plan := loadRealPlan(t, "terraform/empty-config.json")
	doc := renderSVG(t, plan, SVGOptions{})

	assertWellFormed(t, doc)
	if !strings.Contains(doc, "No resources to display") {
		t.Errorf("empty document has no explicit message")
	}
	if !strings.Contains(doc, `<g class="legend">`) {
		t.Errorf("empty document has no legend group")
	}
}

func TestSVG_Truncation(t *testing.T) {
	long := strings.Repeat("very_long_resource_name_", 6)
	plan := &terraform.Plan{ResourceChanges: []terraform.ResourceChange{{
		Address: "aws_instance." + long, Mode: "managed", Type: "aws_instance", Name: long,
		ProviderName: "registry.terraform.io/hashicorp/aws",
		Change:       terraform.Change{Actions: []string{"create"}},
	}}}
	doc := renderSVG(t, plan, SVGOptions{})

	assertWellFormed(t, doc)
	if !strings.Contains(doc, "…</text>") {
		t.Errorf("long name is not truncated with an ellipsis")
	}
	if !strings.Contains(doc, "<title>aws_instance."+long+"</title>") {
		t.Errorf("tooltip must keep the full address")
	}
}

func TestSVG_NodeStyles(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/aws-network.json")
	doc := renderSVG(t, plan, SVGOptions{})
	for _, class := range []string{"status-create", "status-update", "status-delete", "status-replace"} {
		if !strings.Contains(doc, `class="node `+class+`"`) {
			t.Errorf("no node with class %s", class)
		}
	}
	for _, dash := range []string{`stroke-dasharray="6 3"`, `stroke-dasharray="2 3"`, `stroke-dasharray="8 3 2 3"`} {
		if !strings.Contains(doc, dash) {
			t.Errorf("document misses border style %s", dash)
		}
	}
}

func TestNewFormatter_SVG(t *testing.T) {
	f, err := NewFormatter(FormatSVG)
	if err != nil {
		t.Fatalf("NewFormatter: %v", err)
	}
	if reflect.TypeOf(f) != reflect.TypeOf(&SVGFormatter{}) {
		t.Errorf("NewFormatter(svg) = %T, want *SVGFormatter", f)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		in   string
		max  int
		want string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"abcdefghijk", 10, "abcdefghi…"},
		{"日本語日本語", 4, "日本語…"},
		{"abc", 0, "…"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := truncate(tt.in, tt.max); got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
			}
		})
	}
}
