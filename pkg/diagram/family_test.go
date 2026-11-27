package diagram

import (
	"reflect"
	"strings"
	"testing"
)

func TestFamily(t *testing.T) {
	tests := []struct {
		typ  string
		want []string
	}{
		{"aws_instance", []string{"aws-ec2", "aws", "compute", "generic"}},
		{"azurerm_linux_virtual_machine", []string{"azure-vm", "azure", "compute", "generic"}},
		{"google_compute_instance", []string{"gcp-gce", "gcp", "compute", "generic"}},
		{"aws_lb", []string{"aws-elb", "aws", "network", "generic"}},
		{"aws_alb", []string{"aws-elb", "aws", "generic"}},
		{"aws_elb", []string{"aws-elb", "aws", "generic"}},
		{"aws_foo_bucket", []string{"aws", "storage", "generic"}},
		{"google_foo_firewall", []string{"gcp", "security", "generic"}},
		{"aws_monkey_thing", []string{"aws", "generic"}},
		{"azurerm_albumin", []string{"azure", "generic"}},
		{"random_pet", []string{"generic"}},
		{"vault_foo", []string{"generic"}},
		{"aws_foo_role_bucket", []string{"aws", "security", "generic"}}, // security beats storage
		{"aws_foo_db_instance", []string{"aws", "database", "generic"}}, // database beats compute
		{"", []string{"generic"}},
	}
	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			got := Family(tt.typ, "")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Family(%q) = %v, want %v", tt.typ, got, tt.want)
			}
			if len(got) == 0 || got[len(got)-1] != "generic" {
				t.Errorf("chain must end with generic: %v", got)
			}
		})
	}
}

func TestServiceTableSize(t *testing.T) {
	counts := map[string]int{}
	for typ := range serviceFamilies {
		switch {
		case strings.HasPrefix(typ, "aws_"):
			counts["aws"]++
		case strings.HasPrefix(typ, "azurerm_"):
			counts["azure"]++
		case strings.HasPrefix(typ, "google_"):
			counts["gcp"]++
		}
	}
	for cloud, n := range counts {
		if n < 25 {
			t.Errorf("%s service table has %d entries, want at least 25", cloud, n)
		}
	}
	if len(counts) != 3 {
		t.Errorf("clouds in table: %v", counts)
	}
}

func TestServiceTableHits(t *testing.T) {
	for typ, svc := range serviceFamilies {
		if got := Family(typ, ""); got[0] != svc {
			t.Errorf("Family(%q)[0] = %q, want %q", typ, got[0], svc)
		}
	}
}
