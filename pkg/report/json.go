package report

import (
	"encoding/json"
	"io"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// JSONFormatter generates JSON-formatted reports.
type JSONFormatter struct{}

// JSONReport represents the JSON output structure.
type JSONReport struct {
	TerraformVersion string                     `json:"terraform_version,omitempty"`
	Summary          JSONSummary                `json:"summary"`
	Changes          []JSONResourceChange       `json:"changes"`
	ByType           map[string]JSONTypeSummary `json:"by_type"`
}

// JSONSummary represents the summary section in JSON output.
type JSONSummary struct {
	TotalChanges int `json:"total_changes"`
	ToAdd        int `json:"to_add"`
	ToChange     int `json:"to_change"`
	ToDestroy    int `json:"to_destroy"`
	ToReplace    int `json:"to_replace"`
}

// JSONTypeSummary represents per-type statistics.
type JSONTypeSummary struct {
	Type      string   `json:"type"`
	ToAdd     int      `json:"to_add"`
	ToChange  int      `json:"to_change"`
	ToDestroy int      `json:"to_destroy"`
	ToReplace int      `json:"to_replace"`
	Resources []string `json:"resources"`
}

// JSONResourceChange represents a single resource change.
type JSONResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Action  string `json:"action"`
}

// Format generates a JSON report.
func (f *JSONFormatter) Format(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error {
	report := f.buildReport(plan, summary)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func (f *JSONFormatter) buildReport(plan *terraform.Plan, summary *terraform.Summary) *JSONReport {
	report := &JSONReport{
		TerraformVersion: plan.TerraformVersion,
		Summary: JSONSummary{
			TotalChanges: summary.TotalChanges,
			ToAdd:        summary.ToAdd,
			ToChange:     summary.ToChange,
			ToDestroy:    summary.ToDestroy,
			ToReplace:    summary.ToReplace,
		},
		Changes: make([]JSONResourceChange, 0),
		ByType:  make(map[string]JSONTypeSummary),
	}

	// Build changes list and by-type map
	typeResources := make(map[string][]string)

	for _, rc := range plan.ResourceChanges {
		action := rc.Change.GetPrimaryAction()
		if action == terraform.ActionNoOp || action == terraform.ActionRead {
			continue
		}

		report.Changes = append(report.Changes, JSONResourceChange{
			Address: rc.Address,
			Type:    rc.Type,
			Name:    rc.Name,
			Action:  string(action),
		})

		typeResources[rc.Type] = append(typeResources[rc.Type], rc.Address)
	}

	// Build by-type summaries
	for typeName, ts := range summary.ByType {
		report.ByType[typeName] = JSONTypeSummary{
			Type:      ts.Type,
			ToAdd:     ts.ToAdd,
			ToChange:  ts.ToChange,
			ToDestroy: ts.ToDestroy,
			ToReplace: ts.ToReplace,
			Resources: typeResources[typeName],
		}
	}

	return report
}
