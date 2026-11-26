package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// scalePlan builds an in-memory plan of n changes spread over several providers, resource types and
// actions, with names that stress ordering and encoding. Nothing is written to disk, so the dataset
// stays small while the volume stays configurable.
func scalePlan(n int) (*terraform.Plan, map[terraform.ActionType]int) {
	types := []string{
		"aws_instance", "aws_s3_bucket", "azurerm_resource_group", "google_compute_instance",
		"kubernetes_deployment", "zz_last_type", "aa_first_type",
	}
	actions := [][]string{
		{"create"}, {"update"}, {"delete"}, {"delete", "create"}, {"create", "delete"}, {"no-op"}, {"read"}, {},
	}
	names := []string{"main", "web", "ünïcode-é", "with space", `quote"d`, "日本語"}

	want := map[terraform.ActionType]int{}
	changes := make([]terraform.ResourceChange, 0, n)
	for i := range n {
		typ := types[i%len(types)]
		name := fmt.Sprintf("%s-%04d", names[i%len(names)], (n-i)%97)
		rc := terraform.ResourceChange{
			Address: fmt.Sprintf("module.m%d.%s.%s[%d]", i%5, typ, name, i),
			Mode:    "managed",
			Type:    typ,
			Name:    name,
			Change:  terraform.Change{Actions: actions[i%len(actions)]},
		}
		want[rc.Change.GetPrimaryAction()]++
		changes = append(changes, rc)
	}
	return &terraform.Plan{FormatVersion: "1.2", TerraformVersion: "1.17.2", ResourceChanges: changes}, want
}

func TestScale_LargePlan(t *testing.T) {
	const n = 5000
	plan, want := scalePlan(n)
	summary := terraform.NewSummary(plan.ResourceChanges)

	counted := want[terraform.ActionCreate] + want[terraform.ActionUpdate] + want[terraform.ActionDelete] + want[terraform.ActionReplace]
	if summary.TotalChanges != counted {
		t.Fatalf("TotalChanges = %d, want %d", summary.TotalChanges, counted)
	}
	if summary.ToAdd != want[terraform.ActionCreate] || summary.ToChange != want[terraform.ActionUpdate] ||
		summary.ToDestroy != want[terraform.ActionDelete] || summary.ToReplace != want[terraform.ActionReplace] {
		t.Errorf("summary = %+v, want per-action counts %v", summary, want)
	}

	t.Run("deterministic formats", func(t *testing.T) {
		for name, f := range map[string]Formatter{"markdown": &MarkdownFormatter{}, "text": &TextFormatter{}, "json": &JSONFormatter{}} {
			if got := assertDeterministic(t, f, plan, summary); got == "" {
				t.Errorf("%s output is empty", name)
			}
		}
	})

	t.Run("json lists every counted change in plan order", func(t *testing.T) {
		var rep JSONReport
		if err := json.Unmarshal([]byte(render(t, &JSONFormatter{}, plan, summary)), &rep); err != nil {
			t.Fatalf("json output is not valid JSON: %v", err)
		}
		if rep.Summary.TotalChanges != counted || len(rep.Changes) != counted {
			t.Errorf("json total/changes = %d/%d, want %d", rep.Summary.TotalChanges, len(rep.Changes), counted)
		}
		if len(rep.ByType) != 7 {
			t.Errorf("by_type has %d entries, want 7", len(rep.ByType))
		}
	})

	t.Run("markdown types sorted alphabetically", func(t *testing.T) {
		md := render(t, &MarkdownFormatter{}, plan, summary)
		section := md[strings.Index(md, "## Changes by Resource Type"):]

		var headings []string
		for _, line := range strings.Split(section, "\n") {
			if strings.HasPrefix(line, "### ") {
				headings = append(headings, strings.Fields(strings.TrimPrefix(line, "### "))[0])
			}
		}
		if len(headings) != 7 {
			t.Fatalf("found %d type sections, want 7: %v", len(headings), headings)
		}
		if !sort.StringsAreSorted(headings) {
			t.Errorf("type sections are not sorted: %v", headings)
		}
	})
}

func BenchmarkFormatLargePlan(b *testing.B) {
	plan, _ := scalePlan(10000)
	summary := terraform.NewSummary(plan.ResourceChanges)

	for name, f := range map[string]Formatter{"markdown": &MarkdownFormatter{}, "text": &TextFormatter{}, "json": &JSONFormatter{}} {
		b.Run(name, func(b *testing.B) {
			for range b.N {
				if err := f.Format(plan, summary, discard{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }
