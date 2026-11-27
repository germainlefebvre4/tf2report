package report

import (
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// TestSVG_Golden pins the svg output for representative plans and option sets. Refresh the golden
// files only on purpose, with -update, and review every diff.
func TestSVG_Golden(t *testing.T) {
	fixtures := []struct {
		name    string
		fixture string
	}{
		{"aws_network", "diagram/aws-network.json"},
		{"azure_network", "diagram/azure-network.json"},
		{"gcp_network", "diagram/gcp-network.json"},
		{"modules_context", "diagram/modules-context.json"},
		{"cycle", "diagram/cycle.json"},
		{"no_configuration", "diagram/no-configuration.json"},
		{"empty_config", "terraform/empty-config.json"},
		{"provider_random", "terraform/provider-random.json"},
		{"moved", "terraform/moved.json"},
		{"import", "terraform/import.json"},
	}
	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			plan := loadRealPlan(t, fx.fixture)
			got := render(t, NewSVGFormatter(SVGOptions{}), plan, terraform.NewSummary(plan.ResourceChanges))
			assertWellFormed(t, got)
			assertDeterministic(t, NewSVGFormatter(SVGOptions{}), plan, terraform.NewSummary(plan.ResourceChanges))
			assertGolden(t, "svg_"+fx.name+".golden", got)
		})
	}
}

func TestSVG_GoldenOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   SVGOptions
		filter func(*terraform.Plan)
	}{
		{"aws_network_type_filter", SVGOptions{TypeFiltered: true}, func(p *terraform.Plan) {
			p.ResourceChanges = terraform.FilterByType(p.ResourceChanges, []string{"aws_instance"})
		}},
		{"aws_network_action_filter", SVGOptions{ActionFiltered: true}, func(p *terraform.Plan) {
			p.ResourceChanges = terraform.FilterByAction(p.ResourceChanges, []terraform.ActionType{terraform.ActionCreate})
		}},
		{"aws_network_action_noop", SVGOptions{ActionFiltered: true}, func(p *terraform.Plan) {
			p.ResourceChanges = terraform.FilterByAction(p.ResourceChanges,
				[]terraform.ActionType{terraform.ActionCreate, terraform.ActionNoOp})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := loadRealPlan(t, "diagram/aws-network.json")
			tt.filter(plan)
			got := render(t, NewSVGFormatter(tt.opts), plan, terraform.NewSummary(plan.ResourceChanges))
			assertWellFormed(t, got)
			assertGolden(t, "svg_"+tt.name+".golden", got)
		})
	}
}

// A moved-only resource is hidden by default, and an imported resource with an update is drawn
// once, as an update.
func TestSVG_MovedAndImported(t *testing.T) {
	for _, fixture := range []string{"terraform/moved.json", "terraform/import.json"} {
		plan := loadRealPlan(t, fixture)
		doc := renderSVG(t, plan, SVGOptions{})
		drawn := 0
		for _, rc := range plan.ResourceChanges {
			if rc.Change.GetPrimaryAction() != terraform.ActionNoOp && rc.Mode != "data" {
				drawn++
			}
		}
		if got := countNodes(doc); got != drawn {
			t.Errorf("%s: %d nodes drawn, want %d", fixture, got, drawn)
		}
	}
}
