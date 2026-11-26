package report

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// TestSVG_SafetyMatrix draws every diagram and sensitive fixture with every option set and checks
// well-formedness, absence of attribute values and determinism.
func TestSVG_SafetyMatrix(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("..", "..", "testdata", "plans", "diagram", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range fixtures {
		names = append(names, "diagram/"+filepath.Base(f))
	}
	names = append(names, "terraform/sensitive.json", "terraform/provider-sensitive.json", "opentofu/sensitive.json",
		"terraform/moved.json", "terraform/import.json")

	optionSets := map[string]SVGOptions{
		"no filter":     {},
		"type filter":   {TypeFiltered: true},
		"action filter": {ActionFiltered: true},
		"valid pack":    {IconsDir: testdataIcons("valid")},
		"hostile pack":  {IconsDir: testdataIcons("hostile")},
	}
	for _, name := range names {
		values := attributeValues(t, name)
		for optName, opts := range optionSets {
			t.Run(name+"/"+optName, func(t *testing.T) {
				plan := loadRealPlan(t, name)
				first := renderSVG(t, plan, opts)
				assertWellFormed(t, first)
				ids := identifiers(plan)
				for _, v := range values {
					if len(v) >= 6 && !strings.Contains(ids, v) && strings.Contains(first, v) {
						t.Errorf("output contains attribute value %q", v)
					}
				}
				if second := renderSVG(t, plan, opts); first != second {
					t.Errorf("two renders differ")
				}
			})
		}
	}
}

// identifiers joins the addresses, types and names of a plan: a value that equals one of them is
// legitimately drawn and is not a leak.
func identifiers(plan *terraform.Plan) string {
	var parts []string
	for _, rc := range plan.ResourceChanges {
		parts = append(parts, rc.Address, rc.Type, rc.Name, rc.ProviderName)
	}
	return strings.Join(parts, "\n")
}
