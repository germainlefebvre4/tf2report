package terraform

import "testing"

func TestNewSummary(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		total   int
		add     int
		change  int
		destroy int
		replace int
		byType  map[string]TypeSummary
	}{
		{
			name: "basic", fixture: "basic.json",
			total: 4, add: 1, change: 1, destroy: 1, replace: 1,
			byType: map[string]TypeSummary{
				"aws_s3_bucket":      {Type: "aws_s3_bucket", ToAdd: 1},
				"aws_instance":       {Type: "aws_instance", ToChange: 1, ToReplace: 1},
				"aws_security_group": {Type: "aws_security_group", ToDestroy: 1},
			},
		},
		{
			name: "no-op and read are excluded", fixture: "noop-and-read.json",
			total: 1, add: 1,
			byType: map[string]TypeSummary{
				"aws_s3_bucket": {Type: "aws_s3_bucket", ToAdd: 1},
			},
		},
		{
			name: "unknown action is left out", fixture: "unknown-action.json",
			total: 1, add: 1,
			byType: map[string]TypeSummary{
				"aws_s3_bucket": {Type: "aws_s3_bucket", ToAdd: 1},
			},
		},
		{name: "no changes", fixture: "no-changes.json", byType: map[string]TypeSummary{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := NewParser().Parse(fixture(t, tt.fixture))
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			s := NewSummary(plan.ResourceChanges)

			if s.TotalChanges != tt.total || s.ToAdd != tt.add || s.ToChange != tt.change ||
				s.ToDestroy != tt.destroy || s.ToReplace != tt.replace {
				t.Errorf("counts = total %d add %d change %d destroy %d replace %d, want %d %d %d %d %d",
					s.TotalChanges, s.ToAdd, s.ToChange, s.ToDestroy, s.ToReplace,
					tt.total, tt.add, tt.change, tt.destroy, tt.replace)
			}

			if sum := s.ToAdd + s.ToChange + s.ToDestroy + s.ToReplace; sum != s.TotalChanges {
				t.Errorf("total %d != add+change+destroy+replace %d", s.TotalChanges, sum)
			}

			if len(s.ByType) != len(tt.byType) {
				t.Errorf("ByType has %d entries, want %d: %v", len(s.ByType), len(tt.byType), s.ByType)
			}
			var byTypeSum int
			for name, want := range tt.byType {
				got, ok := s.ByType[name]
				if !ok || got != want {
					t.Errorf("ByType[%q] = %+v (present=%v), want %+v", name, got, ok, want)
				}
			}
			for _, ts := range s.ByType {
				byTypeSum += ts.ToAdd + ts.ToChange + ts.ToDestroy + ts.ToReplace
			}
			if byTypeSum != s.TotalChanges {
				t.Errorf("sum over ByType %d != total %d", byTypeSum, s.TotalChanges)
			}
		})
	}
}
