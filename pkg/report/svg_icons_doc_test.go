package report

import (
	"regexp"
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

func manyInstances(n int) *terraform.Plan {
	plan := &terraform.Plan{FormatVersion: "1.2"}
	for i := range n {
		name := "web" + itoa(i)
		plan.ResourceChanges = append(plan.ResourceChanges, terraform.ResourceChange{
			Address: "aws_instance." + name, Mode: "managed", Type: "aws_instance", Name: name,
			ProviderName: "registry.terraform.io/hashicorp/aws",
			Change:       terraform.Change{Actions: []string{"create"}},
		})
	}
	return plan
}

func TestSVG_IconDefinedOnce(t *testing.T) {
	doc := renderSVG(t, manyInstances(20), SVGOptions{})
	assertWellFormed(t, doc)
	if n := strings.Count(doc, `<symbol id="icon-aws-ec2"`); n != 1 {
		t.Errorf("symbol defined %d times, want 1", n)
	}
	if n := strings.Count(doc, `<use href="#icon-aws-ec2"`); n != 20 {
		t.Errorf("use count = %d, want 20", n)
	}
	small, large := len(renderSVG(t, manyInstances(10), SVGOptions{})), len(renderSVG(t, manyInstances(20), SVGOptions{}))
	symbol := len(`<symbol id="icon-aws-ec2" viewBox="0 0 24 24"></symbol>`) + 200
	if large-small < 10*100 || strings.Count(renderSVG(t, manyInstances(20), SVGOptions{}), "<symbol") != 1 {
		t.Errorf("document growth looks wrong: %d -> %d (symbol ~%d)", small, large, symbol)
	}
}

func TestSVG_UnknownTypeGetsGenericIcon(t *testing.T) {
	plan := &terraform.Plan{ResourceChanges: []terraform.ResourceChange{{
		Address: "weird_thing.x", Mode: "managed", Type: "weird_thing", Name: "x", ProviderName: "x/y/weird",
		Change: terraform.Change{Actions: []string{"create"}},
	}}}
	doc := renderSVG(t, plan, SVGOptions{})
	assertWellFormed(t, doc)
	if !strings.Contains(doc, `<use href="#icon-generic"`) {
		t.Errorf("unknown type does not use the generic icon")
	}
}

func TestSVG_HostilePackIsHarmless(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/aws-network.json")
	doc := renderSVG(t, plan, SVGOptions{IconsDir: testdataIcons("hostile")})
	assertWellFormed(t, doc)
	for _, re := range []string{`(?i)<script`, `(?i)\son[a-z]+\s*=`, `(?i)foreignObject`, `(?i)href="[^#]`} {
		if regexp.MustCompile(re).MatchString(doc) {
			t.Errorf("output matches forbidden pattern %s", re)
		}
	}
}

func TestSVG_PackGolden(t *testing.T) {
	plan := loadRealPlan(t, "diagram/aws-network.json")
	got := renderSVG(t, plan, SVGOptions{IconsDir: testdataIcons("valid")})
	assertWellFormed(t, got)
	if !strings.Contains(got, `<symbol id="icon-aws-ec2"`) {
		t.Errorf("pack icon aws-ec2 not used")
	}
	assertGolden(t, "svg_aws_network_pack.golden", got)
}

func TestSVG_InvalidIconsDirWritesNothing(t *testing.T) {
	plan, _ := loadPlan(t, "diagram/aws-network.json")
	var sb strings.Builder
	err := NewSVGFormatter(SVGOptions{IconsDir: "/nonexistent/icons"}).Format(plan, nil, &sb)
	if err == nil || !strings.Contains(err.Error(), "/nonexistent/icons") || sb.Len() != 0 {
		t.Errorf("err=%v written=%d", err, sb.Len())
	}
}
