package terraform

// Plan represents a Terraform plan file structure.
type Plan struct {
	FormatVersion    string            `json:"format_version"`
	TerraformVersion string            `json:"terraform_version"`
	ResourceChanges  []ResourceChange  `json:"resource_changes"`
	OutputChanges    map[string]Change `json:"output_changes,omitempty"`
}

// ResourceChange represents a single resource change in the plan.
type ResourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       Change `json:"change"`
	ActionReason string `json:"action_reason,omitempty"`
	PreviousAddr string `json:"previous_address,omitempty"`
}

// Change represents the before and after state of a resource.
type Change struct {
	Actions         []string               `json:"actions"`
	Before          map[string]interface{} `json:"before"`
	After           map[string]interface{} `json:"after"`
	AfterUnknown    map[string]interface{} `json:"after_unknown,omitempty"`
	BeforeSensitive map[string]interface{} `json:"before_sensitive,omitempty"`
	AfterSensitive  map[string]interface{} `json:"after_sensitive,omitempty"`
}

// ActionType represents possible Terraform actions.
type ActionType string

const (
	ActionCreate  ActionType = "create"
	ActionDelete  ActionType = "delete"
	ActionUpdate  ActionType = "update"
	ActionReplace ActionType = "replace"
	ActionRead    ActionType = "read"
	ActionNoOp    ActionType = "no-op"
)

// GetPrimaryAction returns the primary action for a change.
// Handles multi-action scenarios like ["delete", "create"] for replacements.
func (c *Change) GetPrimaryAction() ActionType {
	if len(c.Actions) == 0 {
		return ActionNoOp
	}

	// Handle replacement (delete + create)
	if len(c.Actions) == 2 {
		hasDelete := false
		hasCreate := false
		for _, action := range c.Actions {
			if action == "delete" {
				hasDelete = true
			}
			if action == "create" {
				hasCreate = true
			}
		}
		if hasDelete && hasCreate {
			return ActionReplace
		}
	}

	// Single action
	action := c.Actions[0]
	switch action {
	case "create":
		return ActionCreate
	case "delete":
		return ActionDelete
	case "update":
		return ActionUpdate
	case "read":
		return ActionRead
	case "no-op":
		return ActionNoOp
	default:
		return ActionNoOp
	}
}

// Summary contains aggregated statistics about plan changes.
type Summary struct {
	TotalChanges int
	ToAdd        int
	ToChange     int
	ToDestroy    int
	ToReplace    int
	ByType       map[string]TypeSummary
}

// TypeSummary contains change counts for a specific resource type.
type TypeSummary struct {
	Type      string
	ToAdd     int
	ToChange  int
	ToDestroy int
	ToReplace int
}

// NewSummary creates a Summary from a slice of ResourceChanges.
func NewSummary(changes []ResourceChange) *Summary {
	summary := &Summary{
		ByType: make(map[string]TypeSummary),
	}

	for _, rc := range changes {
		action := rc.Change.GetPrimaryAction()

		// Skip no-op and read actions
		if action == ActionNoOp || action == ActionRead {
			continue
		}

		summary.TotalChanges++

		switch action {
		case ActionCreate:
			summary.ToAdd++
		case ActionUpdate:
			summary.ToChange++
		case ActionDelete:
			summary.ToDestroy++
		case ActionReplace:
			summary.ToReplace++
		}

		// Update type summary
		ts := summary.ByType[rc.Type]
		ts.Type = rc.Type
		switch action {
		case ActionCreate:
			ts.ToAdd++
		case ActionUpdate:
			ts.ToChange++
		case ActionDelete:
			ts.ToDestroy++
		case ActionReplace:
			ts.ToReplace++
		}
		summary.ByType[rc.Type] = ts
	}

	return summary
}
