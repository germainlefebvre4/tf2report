package terraform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// realPlan describes an authentic `terraform show -json` plan in testdata/plans/terraform/, generated
// by scripts/gen-plans.sh from testdata/scenarios/<name>. The expected counts were derived from the
// JSON itself (jq over resource_changes[].change.actions), not from this tool.
type realPlan struct {
	name string
	// parses is false when the plan carries a boolean sensitivity flag (every plan that creates or
	// deletes a resource does), which the parser rejects today: see TestKnownGap_RealPlansRejected.
	parses                        bool
	add, change, destroy, replace int
	ignored                       int // no-op and read entries, excluded from every count
}

var realPlans = []realPlan{
	{name: "count-for-each", add: 8},
	{name: "create-before-destroy", parses: true, replace: 1},
	{name: "data-source", add: 1, ignored: 1},
	{name: "destroy-all", destroy: 3},
	{name: "empty-config", parses: true},
	{name: "forced-replace", parses: true, replace: 2, ignored: 1},
	{name: "import", parses: true, change: 1},
	{name: "mixed-actions", add: 1, change: 1, destroy: 1, replace: 1, ignored: 1},
	{name: "moved", parses: true, change: 1, ignored: 1},
	{name: "nested-modules", add: 7},
	{name: "no-changes", parses: true, ignored: 2},
	{name: "sensitive", add: 1},
}

func realPlanPath(name string) string {
	return filepath.Join("..", "..", "testdata", "plans", "terraform", name+".json")
}

func realPlanData(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(realPlanPath(name))
	if err != nil {
		t.Fatalf("read real plan %s: %v", name, err)
	}
	return data
}

// withoutSensitivityFlags removes what the parser cannot read today (boolean sensitivity and unknown
// flags, and the output changes), so the rest of the real plan can still exercise the summary logic.
func withoutSensitivityFlags(t *testing.T, data []byte) []byte {
	t.Helper()

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("decode plan: %v", err)
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
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("encode plan: %v", err)
	}
	return out
}

// Every file in testdata/plans/terraform must be classified in realPlans, so a regenerated or newly
// added scenario cannot silently escape the tests.
func TestRealPlans_AllClassified(t *testing.T) {
	files, err := filepath.Glob(realPlanPath("*"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}

	known := map[string]bool{}
	for _, rp := range realPlans {
		known[rp.name] = true
	}
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".json")
		if !known[name] {
			t.Errorf("%s is not listed in realPlans", name)
		}
		delete(known, name)
	}
	for name := range known {
		t.Errorf("realPlans lists %s but the fixture file is missing", name)
	}
}

// Regular test: summary counts of the plans the parser reads as they are.
func TestRealPlans_SummaryAsExported(t *testing.T) {
	for _, rp := range realPlans {
		if !rp.parses {
			continue
		}
		t.Run(rp.name, func(t *testing.T) {
			plan, err := NewParser().Parse(realPlanData(t, rp.name))
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			assertRealSummary(t, rp, plan)
		})
	}
}

// Regular test: the counting logic is right on every real plan, once the unreadable flags are removed.
// This covers create, delete and mixed plans that the parser rejects as exported.
func TestRealPlans_SummaryWithoutSensitivityFlags(t *testing.T) {
	for _, rp := range realPlans {
		t.Run(rp.name, func(t *testing.T) {
			plan, err := NewParser().Parse(withoutSensitivityFlags(t, realPlanData(t, rp.name)))
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			assertRealSummary(t, rp, plan)
		})
	}
}

func assertRealSummary(t *testing.T, rp realPlan, plan *Plan) {
	t.Helper()

	s := NewSummary(plan.ResourceChanges)
	wantTotal := rp.add + rp.change + rp.destroy + rp.replace
	if s.TotalChanges != wantTotal {
		t.Errorf("TotalChanges = %d, want %d", s.TotalChanges, wantTotal)
	}
	if s.ToAdd != rp.add || s.ToChange != rp.change || s.ToDestroy != rp.destroy || s.ToReplace != rp.replace {
		t.Errorf("counts add/change/destroy/replace = %d/%d/%d/%d, want %d/%d/%d/%d",
			s.ToAdd, s.ToChange, s.ToDestroy, s.ToReplace, rp.add, rp.change, rp.destroy, rp.replace)
	}
	if got := len(plan.ResourceChanges) - wantTotal; got != rp.ignored {
		t.Errorf("ignored (no-op/read) entries = %d, want %d", got, rp.ignored)
	}
}

// PROVISIONAL: pins Known Gap "A plan in the exact shape Terraform exports is rejected"
// (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed; then every
// realPlans entry becomes parses: true.
func TestKnownGap_RealPlansRejected(t *testing.T) {
	for _, rp := range realPlans {
		if rp.parses {
			continue
		}
		t.Run(rp.name, func(t *testing.T) {
			_, err := NewParser().Parse(realPlanData(t, rp.name))
			if err == nil {
				t.Fatalf("Parse succeeded; the gap seems fixed, update this test")
			}
			if !strings.Contains(err.Error(), "cannot unmarshal bool into Go struct field") {
				t.Errorf("error = %q, want a bool unmarshal error", err)
			}
		})
	}
}

// Regular test: attributes of real plans that the report logic relies on or must tolerate.
func TestRealPlans_Attributes(t *testing.T) {
	parse := func(t *testing.T, name string) *Plan {
		t.Helper()
		plan, err := NewParser().Parse(withoutSensitivityFlags(t, realPlanData(t, name)))
		if err != nil {
			t.Fatalf("Parse %s: %v", name, err)
		}
		return plan
	}
	find := func(t *testing.T, plan *Plan, address string) ResourceChange {
		t.Helper()
		for _, rc := range plan.ResourceChanges {
			if rc.Address == address {
				return rc
			}
		}
		t.Fatalf("address %q not found in plan", address)
		return ResourceChange{}
	}

	t.Run("create-before-destroy order is a replace", func(t *testing.T) {
		rc := find(t, parse(t, "create-before-destroy"), "terraform_data.cbd")
		if got := rc.Change.Actions; len(got) != 2 || got[0] != "create" || got[1] != "delete" {
			t.Fatalf("actions = %v, want [create delete]", got)
		}
		if rc.Change.GetPrimaryAction() != ActionReplace {
			t.Errorf("primary action = %q, want %q", rc.Change.GetPrimaryAction(), ActionReplace)
		}
	})

	t.Run("destroy-first order is a replace", func(t *testing.T) {
		rc := find(t, parse(t, "mixed-actions"), "terraform_data.replace")
		if got := rc.Change.Actions; len(got) != 2 || got[0] != "delete" || got[1] != "create" {
			t.Fatalf("actions = %v, want [delete create]", got)
		}
		if rc.Change.GetPrimaryAction() != ActionReplace {
			t.Errorf("primary action = %q, want %q", rc.Change.GetPrimaryAction(), ActionReplace)
		}
	})

	t.Run("action reasons", func(t *testing.T) {
		forced := parse(t, "forced-replace")
		for address, want := range map[string]string{
			"terraform_data.requested": "replace_by_request",
			"terraform_data.tainted":   "replace_because_tainted",
		} {
			if got := find(t, forced, address).ActionReason; got != want {
				t.Errorf("%s action_reason = %q, want %q", address, got, want)
			}
		}
		if got := find(t, parse(t, "mixed-actions"), "terraform_data.gone").ActionReason; got != "delete_because_no_resource_config" {
			t.Errorf("gone action_reason = %q, want delete_because_no_resource_config", got)
		}
	})

	t.Run("moved resource keeps its previous address", func(t *testing.T) {
		plan := parse(t, "moved")
		rc := find(t, plan, "terraform_data.renamed_and_changed")
		if rc.PreviousAddr != "terraform_data.moved_and_changed" {
			t.Errorf("previous_address = %q, want terraform_data.moved_and_changed", rc.PreviousAddr)
		}
		if rc.Change.GetPrimaryAction() != ActionUpdate {
			t.Errorf("primary action = %q, want %q", rc.Change.GetPrimaryAction(), ActionUpdate)
		}
		// A pure move is a no-op and is not counted.
		pure := find(t, plan, "terraform_data.new_name")
		if got := pure.Change.GetPrimaryAction(); got != ActionNoOp {
			t.Errorf("pure move primary action = %q, want %q", got, ActionNoOp)
		}
	})

	t.Run("deferred data source is a read and is not counted", func(t *testing.T) {
		plan := parse(t, "data-source")
		rc := find(t, plan, "data.terraform_remote_state.other")
		if rc.Mode != "data" {
			t.Errorf("mode = %q, want data", rc.Mode)
		}
		if rc.Change.GetPrimaryAction() != ActionRead {
			t.Errorf("primary action = %q, want %q", rc.Change.GetPrimaryAction(), ActionRead)
		}
		if _, counted := NewSummary(plan.ResourceChanges).ByType["terraform_remote_state"]; counted {
			t.Errorf("data source read appears in the per-type summary")
		}
	})

	t.Run("imported resource is an update", func(t *testing.T) {
		rc := find(t, parse(t, "import"), "terraform_data.adopted")
		if rc.Change.GetPrimaryAction() != ActionUpdate {
			t.Errorf("primary action = %q, want %q", rc.Change.GetPrimaryAction(), ActionUpdate)
		}
	})

	t.Run("module and index addresses are kept verbatim", func(t *testing.T) {
		plan := parse(t, "nested-modules")
		for _, address := range []string{
			"module.network.terraform_data.vpc",
			"module.network.module.subnets.terraform_data.subnet[1]",
			"module.app[0].terraform_data.server",
			`module.region["eu"].terraform_data.server`,
		} {
			find(t, plan, address)
		}
	})

	t.Run("special keys in addresses", func(t *testing.T) {
		plan := parse(t, "count-for-each")
		for _, address := range []string{
			"terraform_data.counted[2]",
			`terraform_data.keyed["with space"]`,
			`terraform_data.keyed["quote\"d"]`,
			`terraform_data.keyed["ünïcode-é"]`,
		} {
			find(t, plan, address)
		}
	})

	t.Run("empty configuration has no resource_changes", func(t *testing.T) {
		plan, err := NewParser().Parse(realPlanData(t, "empty-config"))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if len(plan.ResourceChanges) != 0 {
			t.Errorf("resource_changes = %d entries, want 0", len(plan.ResourceChanges))
		}
	})
}
