package diagram

import "github.com/germainlefebvre4/tf2report/pkg/terraform"

// Select returns the resource changes that are drawn as nodes.
//
// Data sources and `read` actions are never drawn. Without an action filter only
// create, update, delete and replace are drawn; with an action filter everything
// left in the (already filtered) plan is drawn, including no-op. Dependency
// context is added later by Build.
func Select(plan *terraform.Plan, opts Options) []terraform.ResourceChange {
	var out []terraform.ResourceChange
	for _, rc := range plan.ResourceChanges {
		if !drawable(rc) {
			continue
		}
		if !opts.ActionFiltered && rc.Change.GetPrimaryAction() == terraform.ActionNoOp {
			continue
		}
		out = append(out, rc)
	}
	return out
}

// drawable reports whether a change can ever be a node (managed, not read).
func drawable(rc terraform.ResourceChange) bool {
	return rc.Mode != "data" && rc.Change.GetPrimaryAction() != terraform.ActionRead
}

// statusOf maps the primary action of a change to a Status.
func statusOf(rc terraform.ResourceChange) Status {
	switch rc.Change.GetPrimaryAction() {
	case terraform.ActionCreate:
		return StatusCreate
	case terraform.ActionUpdate:
		return StatusUpdate
	case terraform.ActionDelete:
		return StatusDelete
	case terraform.ActionReplace:
		return StatusReplace
	default:
		return StatusNoop
	}
}
