package diagram

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// rc builds a resource change from an address and its actions.
func rc(address, provider string, actions ...string) terraform.ResourceChange {
	a := terraform.ParseAddress(address)
	return terraform.ResourceChange{
		Address: address, Mode: "managed", Type: a.Type, Name: a.Name, ProviderName: provider,
		Change: terraform.Change{Actions: actions},
	}
}

func planOf(changes ...terraform.ResourceChange) *terraform.Plan {
	return &terraform.Plan{ResourceChanges: changes}
}

// loadDiagramFixture parses testdata/plans/<name>.
func loadDiagramFixture(t *testing.T, name string) *terraform.Plan {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "plans", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	plan, err := terraform.NewParser().Parse(data)
	if err != nil {
		t.Fatalf("parse fixture %s: %v", name, err)
	}
	return plan
}

func addresses(changes []terraform.ResourceChange) []string {
	out := []string{}
	for _, c := range changes {
		out = append(out, c.Address)
	}
	return out
}

func TestSelect(t *testing.T) {
	const aws = "registry.terraform.io/hashicorp/aws"
	data := rc("data.aws_ami.x", aws, "read")
	data.Mode = "data"
	plan := planOf(
		rc("aws_a.c", aws, "create"),
		rc("aws_a.u", aws, "update"),
		rc("aws_a.d", aws, "delete"),
		rc("aws_a.r", aws, "delete", "create"),
		rc("aws_a.n", aws, "no-op"),
		rc("aws_a.unknown", aws, "weird"),
		rc("aws_a.read", aws, "read"),
		data,
	)
	tests := []struct {
		name string
		opts Options
		want []string
	}{
		{"no filter", Options{}, []string{"aws_a.c", "aws_a.u", "aws_a.d", "aws_a.r"}},
		{"type filter only", Options{TypeFiltered: true}, []string{"aws_a.c", "aws_a.u", "aws_a.d", "aws_a.r"}},
		{"action filter keeps no-op and unknown", Options{ActionFiltered: true},
			[]string{"aws_a.c", "aws_a.u", "aws_a.d", "aws_a.r", "aws_a.n", "aws_a.unknown"}},
		{"both filters", Options{TypeFiltered: true, ActionFiltered: true},
			[]string{"aws_a.c", "aws_a.u", "aws_a.d", "aws_a.r", "aws_a.n", "aws_a.unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addresses(Select(plan, tt.opts)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect_ReplaceIsOneNode(t *testing.T) {
	scene, err := Build(planOf(rc("aws_a.r", "aws", "delete", "create")), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Nodes) != 1 || scene.Nodes[0].Status != StatusReplace {
		t.Errorf("want one replace node, got %+v", scene.Nodes)
	}
}

type planType = terraform.Plan

func itoa(i int) string { return strconv.Itoa(i) }
