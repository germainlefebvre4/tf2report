package terraform

import "testing"

func TestChange_GetPrimaryAction(t *testing.T) {
	tests := []struct {
		name    string
		actions []string
		want    ActionType
	}{
		{"nil actions", nil, ActionNoOp},
		{"empty actions", []string{}, ActionNoOp},
		{"create", []string{"create"}, ActionCreate},
		{"update", []string{"update"}, ActionUpdate},
		{"delete", []string{"delete"}, ActionDelete},
		{"read", []string{"read"}, ActionRead},
		{"no-op", []string{"no-op"}, ActionNoOp},
		{"delete then create is replace", []string{"delete", "create"}, ActionReplace},
		{"create then delete is replace", []string{"create", "delete"}, ActionReplace},
		{"two other actions: first decides", []string{"update", "delete"}, ActionUpdate},
		{"unknown action", []string{"frobnicate"}, ActionNoOp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Change{Actions: tt.actions}
			if got := c.GetPrimaryAction(); got != tt.want {
				t.Errorf("GetPrimaryAction(%v) = %q, want %q", tt.actions, got, tt.want)
			}
		})
	}
}
