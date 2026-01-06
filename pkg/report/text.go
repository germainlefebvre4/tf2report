package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// TextFormatter generates plain text reports.
type TextFormatter struct{}

// Format generates a plain text report.
func (f *TextFormatter) Format(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error {
	if err := f.writeHeader(plan, summary, w); err != nil {
		return err
	}

	if err := f.writeSummary(summary, w); err != nil {
		return err
	}

	if err := f.writeResourcesByChangeType(plan, summary, w); err != nil {
		return err
	}

	if err := f.writeDetailsByType(plan, summary, w); err != nil {
		return err
	}

	return nil
}

func (f *TextFormatter) writeHeader(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error {
	_, err := fmt.Fprintf(w, "TERRAFORM PLAN SUMMARY\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n\n", strings.Repeat("=", 50))
	if err != nil {
		return err
	}

	if plan.TerraformVersion != "" {
		_, err = fmt.Fprintf(w, "Terraform Version: %s\n\n", plan.TerraformVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *TextFormatter) writeSummary(summary *terraform.Summary, w io.Writer) error {
	_, err := fmt.Fprintf(w, "SUMMARY\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n\n", strings.Repeat("-", 50))
	if err != nil {
		return err
	}

	if summary.TotalChanges == 0 {
		_, err = fmt.Fprintf(w, "No changes. Infrastructure is up-to-date.\n\n")
		return err
	}

	_, err = fmt.Fprintf(w, "Total Changes: %d\n\n", summary.TotalChanges)
	if err != nil {
		return err
	}

	if summary.ToAdd > 0 {
		_, err = fmt.Fprintf(w, "  Add:     %d\n", summary.ToAdd)
		if err != nil {
			return err
		}
	}
	if summary.ToChange > 0 {
		_, err = fmt.Fprintf(w, "  Change:  %d\n", summary.ToChange)
		if err != nil {
			return err
		}
	}
	if summary.ToDestroy > 0 {
		_, err = fmt.Fprintf(w, "  Destroy: %d\n", summary.ToDestroy)
		if err != nil {
			return err
		}
	}
	if summary.ToReplace > 0 {
		_, err = fmt.Fprintf(w, "  Replace: %d\n", summary.ToReplace)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "\n")
	return err
}

func (f *TextFormatter) writeResourcesByChangeType(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error {
	if summary.TotalChanges == 0 {
		return nil
	}

	_, err := fmt.Fprintf(w, "RESOURCES BY CHANGE TYPE\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n\n", strings.Repeat("-", 50))
	if err != nil {
		return err
	}

	// Group resources by action
	byAction := make(map[terraform.ActionType][]struct {
		resourceType string
		resourceName string
	})

	for _, rc := range plan.ResourceChanges {
		action := rc.Change.GetPrimaryAction()
		if action == terraform.ActionNoOp || action == terraform.ActionRead {
			continue
		}
		byAction[action] = append(byAction[action], struct {
			resourceType string
			resourceName string
		}{
			resourceType: rc.Type,
			resourceName: rc.Name,
		})
	}

	// Write each action section
	actions := []terraform.ActionType{
		terraform.ActionCreate,
		terraform.ActionUpdate,
		terraform.ActionReplace,
		terraform.ActionDelete,
	}

	firstSection := true
	for _, action := range actions {
		resources := byAction[action]
		if len(resources) == 0 {
			continue
		}

		if !firstSection {
			_, err = fmt.Fprintf(w, "\n")
			if err != nil {
				return err
			}
		}
		firstSection = false

		actionLabel := f.actionLabel(action)
		_, err = fmt.Fprintf(w, "%s (%d)\n", actionLabel, len(resources))
		if err != nil {
			return err
		}

		// Sort resources
		sort.Slice(resources, func(i, j int) bool {
			if resources[i].resourceType != resources[j].resourceType {
				return resources[i].resourceType < resources[j].resourceType
			}
			return resources[i].resourceName < resources[j].resourceName
		})

		// Calculate column widths
		maxTypeLen := len("Resource Type")
		maxNameLen := len("Resource Name")
		for _, r := range resources {
			if len(r.resourceType) > maxTypeLen {
				maxTypeLen = len(r.resourceType)
			}
			if len(r.resourceName) > maxNameLen {
				maxNameLen = len(r.resourceName)
			}
		}
		maxTypeLen += 2
		maxNameLen += 2

		// Write table
		_, err = fmt.Fprintf(w, "\n  %-*s %s\n", maxTypeLen, "Resource Type", "Resource Name")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "  %s %s\n", strings.Repeat("-", maxTypeLen), strings.Repeat("-", maxNameLen))
		if err != nil {
			return err
		}

		for _, r := range resources {
			_, err = fmt.Fprintf(w, "  %-*s %s\n", maxTypeLen, r.resourceType, r.resourceName)
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintf(w, "\n")
	return err
}

func (f *TextFormatter) writeDetailsByType(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error {
	if len(summary.ByType) == 0 {
		return nil
	}

	_, err := fmt.Fprintf(w, "CHANGES BY RESOURCE TYPE\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n\n", strings.Repeat("-", 50))
	if err != nil {
		return err
	}

	// Sort types alphabetically
	types := make([]string, 0, len(summary.ByType))
	for t := range summary.ByType {
		types = append(types, t)
	}
	sort.Strings(types)

	for i, typeName := range types {
		if i > 0 {
			_, err = fmt.Fprintf(w, "\n")
			if err != nil {
				return err
			}
		}

		ts := summary.ByType[typeName]
		if err := f.writeTypeSection(typeName, ts, plan, w); err != nil {
			return err
		}
	}

	return nil
}

func (f *TextFormatter) writeTypeSection(typeName string, ts terraform.TypeSummary, plan *terraform.Plan, w io.Writer) error {
	total := ts.ToAdd + ts.ToChange + ts.ToDestroy + ts.ToReplace

	_, err := fmt.Fprintf(w, "%s (%d)\n", typeName, total)
	if err != nil {
		return err
	}

	// Collect resources for this type
	type resourceInfo struct {
		address string
		name    string
		action  terraform.ActionType
	}
	var resources []resourceInfo

	for _, rc := range plan.ResourceChanges {
		if rc.Type != typeName {
			continue
		}
		action := rc.Change.GetPrimaryAction()
		if action == terraform.ActionNoOp || action == terraform.ActionRead {
			continue
		}
		resources = append(resources, resourceInfo{
			address: rc.Address,
			name:    rc.Name,
			action:  action,
		})
	}

	if len(resources) == 0 {
		return nil
	}

	// Sort by address
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].address < resources[j].address
	})

	// Calculate column widths
	maxNameLen := len("Resource Name")
	maxTypeLen := len("Change Type")
	for _, r := range resources {
		if len(r.name) > maxNameLen {
			maxNameLen = len(r.name)
		}
		label := f.actionLabel(r.action)
		if len(label) > maxTypeLen {
			maxTypeLen = len(label)
		}
	}

	// Add padding
	maxNameLen += 2
	maxTypeLen += 2

	// Write table
	_, err = fmt.Fprintf(w, "\n  %-*s %-*s %s\n", len(typeName)+2, "Resource Type", maxNameLen, "Resource Name", "Change Type")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "  %s %s %s\n", strings.Repeat("-", len(typeName)+2), strings.Repeat("-", maxNameLen), strings.Repeat("-", maxTypeLen))
	if err != nil {
		return err
	}

	for _, r := range resources {
		changeType := f.actionLabel(r.action)
		_, err = fmt.Fprintf(w, "  %-*s %-*s %s\n", len(typeName)+2, typeName, maxNameLen, r.name, changeType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *TextFormatter) actionLabel(action terraform.ActionType) string {
	switch action {
	case terraform.ActionCreate:
		return "To Add"
	case terraform.ActionUpdate:
		return "To Change"
	case terraform.ActionDelete:
		return "To Destroy"
	case terraform.ActionReplace:
		return "To Replace"
	default:
		return strings.Title(string(action))
	}
}
