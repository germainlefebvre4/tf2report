package report

import (
	"fmt"
	"testing"
	"time"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

const scaleBudget = 2 * time.Second

// diagramScalePlan builds an in-memory plan of n resources over three providers and several
// nested modules, each module chaining its resources through configuration references.
func diagramScalePlan(n int) *terraform.Plan {
	types := []string{"aws_instance", "azurerm_subnet", "google_compute_instance", "aws_s3_bucket", "random_pet"}
	actions := [][]string{{"create"}, {"update"}, {"delete"}, {"delete", "create"}, {"no-op"}}
	const modules = 6

	plan := &terraform.Plan{FormatVersion: "1.2", Configuration: &terraform.Configuration{}}
	calls := map[string]terraform.ModuleCall{}
	for i := range n {
		typ := types[i%len(types)]
		mod := fmt.Sprintf("m%d", i%modules)
		name := fmt.Sprintf("r%d", i)
		plan.ResourceChanges = append(plan.ResourceChanges, terraform.ResourceChange{
			Address:      fmt.Sprintf("module.%s.%s.%s", mod, typ, name),
			Mode:         "managed",
			Type:         typ,
			Name:         name,
			ProviderName: "registry.terraform.io/hashicorp/" + providerOf(typ),
			Change:       terraform.Change{Actions: actions[(i/len(types))%len(actions)]},
		})
		call := calls[mod]
		res := terraform.ConfigResource{Mode: "managed", Type: typ, Name: name}
		if i >= modules {
			prev := i - modules
			res.Expressions = map[string]any{"dep": map[string]any{"references": []any{
				fmt.Sprintf("%s.r%d", types[prev%len(types)], prev),
			}}}
		}
		call.Module.Resources = append(call.Module.Resources, res)
		calls[mod] = call
	}
	plan.Configuration.RootModule.ModuleCalls = calls
	return plan
}

func providerOf(typ string) string {
	for i := range typ {
		if typ[i] == '_' {
			return typ[:i]
		}
	}
	return typ
}

// crossProductPlan builds groups of count-expanded resources where every instance of b depends on
// every instance of a, so the edges grow with the square of the count.
func crossProductPlan(groups, count int) *terraform.Plan {
	plan := &terraform.Plan{FormatVersion: "1.2", Configuration: &terraform.Configuration{}}
	for g := range groups {
		a := terraform.ConfigResource{Mode: "managed", Type: "aws_instance", Name: fmt.Sprintf("a%d", g)}
		b := terraform.ConfigResource{Mode: "managed", Type: "aws_subnet", Name: fmt.Sprintf("b%d", g),
			Expressions: map[string]any{"x": map[string]any{"references": []any{"aws_instance." + a.Name}}}}
		plan.Configuration.RootModule.Resources = append(plan.Configuration.RootModule.Resources, a, b)
		for i := range count {
			for _, c := range []terraform.ConfigResource{a, b} {
				plan.ResourceChanges = append(plan.ResourceChanges, terraform.ResourceChange{
					Address:      fmt.Sprintf("%s.%s[%d]", c.Type, c.Name, i),
					Mode:         "managed",
					Type:         c.Type,
					Name:         c.Name,
					ProviderName: "registry.terraform.io/hashicorp/aws",
					Change:       terraform.Change{Actions: []string{"create"}},
				})
			}
		}
	}
	return plan
}

func assertScale(t *testing.T, plan *terraform.Plan, opts SVGOptions) {
	t.Helper()

	start := time.Now()
	first := renderSVG(t, plan, opts)
	if elapsed := time.Since(start); elapsed > scaleBudget {
		t.Errorf("render took %v, budget is %v", elapsed, scaleBudget)
	}
	assertWellFormed(t, first)
	if second := renderSVG(t, plan, opts); first != second {
		t.Errorf("two renders of the same plan differ")
	}
}

func TestScale_300Resources(t *testing.T) {
	assertScale(t, diagramScalePlan(300), SVGOptions{})
}

func TestScale_300ResourcesWithEdges(t *testing.T) {
	tests := []struct {
		name string
		plan *terraform.Plan
		opts SVGOptions
	}{
		{"cross product", crossProductPlan(10, 15), SVGOptions{}},
		{"module chain", diagramScalePlan(300), SVGOptions{}},
		{"type filter", filterType(diagramScalePlan(300), "aws_instance"), SVGOptions{TypeFiltered: true}},
		{"action filter", filterAction(diagramScalePlan(300), terraform.ActionCreate, terraform.ActionNoOp), SVGOptions{ActionFiltered: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "cross product" && len(tt.plan.ResourceChanges) != 300 {
				t.Fatalf("cross product plan has %d resources, want 300", len(tt.plan.ResourceChanges))
			}
			assertScale(t, tt.plan, tt.opts)
		})
	}
}

func filterType(plan *terraform.Plan, types ...string) *terraform.Plan {
	plan.ResourceChanges = terraform.FilterByType(plan.ResourceChanges, types)
	return plan
}

func filterAction(plan *terraform.Plan, actions ...terraform.ActionType) *terraform.Plan {
	plan.ResourceChanges = terraform.FilterByAction(plan.ResourceChanges, actions)
	return plan
}
