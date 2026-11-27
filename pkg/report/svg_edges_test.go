package report

import (
	"strings"
	"testing"
)

func TestSVG_Edges(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/modules-context.json")
	doc := renderSVG(t, plan, SVGOptions{})
	assertWellFormed(t, doc)

	if n := strings.Count(doc, `<marker id="arrow"`); n != 1 {
		t.Errorf("arrow marker defined %d times, want 1", n)
	}
	edges := strings.Index(doc, `<g class="edges">`)
	lastContainer := strings.LastIndex(doc, `<g class="container`)
	nodes := strings.Index(doc, `<g class="nodes">`)
	if edges < 0 || nodes < 0 || lastContainer >= edges || edges >= nodes {
		t.Errorf("edges must come after containers and before nodes (%d, %d, %d)", lastContainer, edges, nodes)
	}
	group := doc[edges:nodes]
	if got := strings.Count(group, "<path "); got != 2 {
		t.Errorf("edge paths = %d, want 2 (one to the context vpc, one to the module container)", got)
	}
	if !strings.Contains(doc, `class="node status-noop context"`) {
		t.Errorf("context node is missing the context class")
	}
	_, ctx, _ := strings.Cut(doc, `class="node status-noop context"`)
	ctx, _, _ = strings.Cut(ctx, "</g>")
	if strings.Contains(ctx, "font-weight=\"bold\" fill=\"#6e7781\"") {
		t.Errorf("context node shows a status symbol")
	}
}

func TestSVG_NoArrowsWithoutConfiguration(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/no-configuration.json")
	doc := renderSVG(t, plan, SVGOptions{})
	if strings.Contains(doc, `<g class="edges">`) && strings.Contains(section2(doc, `<g class="edges">`), "<path ") {
		t.Errorf("document without configuration has arrows")
	}
}

func TestSVG_CycleRenders(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/cycle.json")
	doc := renderSVG(t, plan, SVGOptions{})
	assertWellFormed(t, doc)
	if got := strings.Count(section2(doc, `<g class="edges">`), "<path "); got != 2 {
		t.Errorf("cycle edges = %d, want 2", got)
	}
}

func section2(doc, marker string) string {
	i := strings.Index(doc, marker)
	if i < 0 {
		return ""
	}
	rest := doc[i:]
	if j := strings.Index(rest, `<g class="nodes">`); j >= 0 {
		return rest[:j]
	}
	return rest
}
