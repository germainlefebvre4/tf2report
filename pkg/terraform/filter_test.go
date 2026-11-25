package terraform

import (
	"reflect"
	"testing"
)

func basicChanges(t *testing.T) []ResourceChange {
	t.Helper()

	plan, err := NewParser().Parse(fixture(t, "basic.json"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return plan.ResourceChanges
}

func addresses(changes []ResourceChange) []string {
	out := make([]string, 0, len(changes))
	for _, rc := range changes {
		out = append(out, rc.Address)
	}
	return out
}

func TestFilterByType(t *testing.T) {
	tests := []struct {
		name  string
		types []string
		want  []string
	}{
		{"nil types returns everything", nil, []string{"aws_s3_bucket.data", "aws_instance.web", "aws_instance.app", "aws_security_group.old"}},
		{"empty types returns everything", []string{}, []string{"aws_s3_bucket.data", "aws_instance.web", "aws_instance.app", "aws_security_group.old"}},
		{"one type", []string{"aws_instance"}, []string{"aws_instance.web", "aws_instance.app"}},
		{"several types", []string{"aws_s3_bucket", "aws_security_group"}, []string{"aws_s3_bucket.data", "aws_security_group.old"}},
		{"absent type matches nothing", []string{"aws_lambda_function"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByType(basicChanges(t), tt.types)
			if got == nil {
				t.Fatalf("FilterByType returned nil, want a non-nil slice")
			}
			if !reflect.DeepEqual(addresses(got), tt.want) {
				t.Errorf("addresses = %v, want %v", addresses(got), tt.want)
			}
		})
	}
}

func TestFilterByAction(t *testing.T) {
	tests := []struct {
		name    string
		actions []ActionType
		want    []string
	}{
		{"nil actions returns everything", nil, []string{"aws_s3_bucket.data", "aws_instance.web", "aws_instance.app", "aws_security_group.old"}},
		{"create", []ActionType{ActionCreate}, []string{"aws_s3_bucket.data"}},
		{"update", []ActionType{ActionUpdate}, []string{"aws_instance.web"}},
		{"delete", []ActionType{ActionDelete}, []string{"aws_security_group.old"}},
		{"replace is judged on the primary action", []ActionType{ActionReplace}, []string{"aws_instance.app"}},
		{"delete and replace", []ActionType{ActionDelete, ActionReplace}, []string{"aws_instance.app", "aws_security_group.old"}},
		{"unknown action matches nothing", []ActionType{ActionType("bogus")}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByAction(basicChanges(t), tt.actions)
			if got == nil {
				t.Fatalf("FilterByAction returned nil, want a non-nil slice")
			}
			if !reflect.DeepEqual(addresses(got), tt.want) {
				t.Errorf("addresses = %v, want %v", addresses(got), tt.want)
			}
		})
	}
}

func TestFilters_combinedAreAnd(t *testing.T) {
	changes := FilterByType(basicChanges(t), []string{"aws_instance"})
	changes = FilterByAction(changes, []ActionType{ActionUpdate})

	want := []string{"aws_instance.web"}
	if !reflect.DeepEqual(addresses(changes), want) {
		t.Errorf("addresses = %v, want %v", addresses(changes), want)
	}
}
