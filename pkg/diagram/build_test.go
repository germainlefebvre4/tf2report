package diagram

import (
	"reflect"
	"sort"
	"testing"
)

func TestProviderName(t *testing.T) {
	tests := []struct{ full, typ, want string }{
		{"registry.terraform.io/hashicorp/aws", "aws_instance", "aws"},
		{"registry.opentofu.org/hashicorp/random", "random_pet", "random"},
		{"terraform.io/builtin/terraform", "terraform_data", "terraform"},
		{"", "azurerm_subnet", "azurerm"},
		{`provider["registry.terraform.io/hashicorp/aws"]`, "aws_vpc", "aws"},
	}
	for _, tt := range tests {
		t.Run(tt.want+"/"+tt.full, func(t *testing.T) {
			if got := providerName(tt.full, tt.typ); got != tt.want {
				t.Errorf("providerName(%q, %q) = %q, want %q", tt.full, tt.typ, got, tt.want)
			}
		})
	}
}

func TestBuild_ContainersNestProviderThenModule(t *testing.T) {
	plan := planOf(
		rc("aws_instance.root", "registry.terraform.io/hashicorp/aws", "create"),
		rc("module.app[0].aws_instance.web", "registry.terraform.io/hashicorp/aws", "create"),
		rc("module.app[1].aws_instance.web", "registry.terraform.io/hashicorp/aws", "create"),
		rc("module.app[1].module.db.aws_db_instance.main", "registry.terraform.io/hashicorp/aws", "create"),
		rc("module.app[0].random_pet.x", "registry.terraform.io/hashicorp/random", "create"),
	)
	scene, err := Build(plan, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(scene.Containers) != 2 || scene.Containers[0].Label != "aws" || scene.Containers[1].Label != "random" {
		t.Fatalf("providers = %+v", scene.Containers)
	}
	aws := scene.Containers[0]
	if len(aws.Nodes) != 1 || aws.Nodes[0].Address != "aws_instance.root" {
		t.Errorf("root-module resource must sit in the provider container: %+v", aws.Nodes)
	}
	var labels []string
	for _, c := range aws.Children {
		labels = append(labels, c.Label)
	}
	if want := []string{"module.app[0]", "module.app[1]"}; !reflect.DeepEqual(labels, want) {
		t.Errorf("module containers = %v, want %v", labels, want)
	}
	if aws.Children[0].Kind != KindModule || aws.Children[0].Path != "module.app[0]" {
		t.Errorf("module container = %+v", aws.Children[0])
	}
	nested := aws.Children[1].Children
	if len(nested) != 1 || nested[0].Path != "module.app[1].module.db" || nested[0].configPath != "module.app.module.db" {
		t.Errorf("nested container = %+v", nested)
	}
	// a module holding two providers appears under each provider
	rnd := scene.Containers[1]
	if len(rnd.Children) != 1 || rnd.Children[0].Label != "module.app[0]" {
		t.Errorf("random provider modules = %+v", rnd.Children)
	}
}

func TestBuild_SameNameDifferentModulesNeverCollide(t *testing.T) {
	const aws = "registry.terraform.io/hashicorp/aws"
	plan := planOf(
		rc("aws_instance.web", aws, "create"),
		rc("module.a.aws_instance.web", aws, "create"),
		rc("module.b.aws_instance.web", aws, "create"),
	)
	scene, _ := Build(plan, Options{})
	seen := map[string]bool{}
	for _, n := range scene.Nodes {
		if seen[n.Address] {
			t.Errorf("duplicate address %s", n.Address)
		}
		seen[n.Address] = true
	}
	if len(scene.Nodes) != 3 {
		t.Errorf("nodes = %d, want 3", len(scene.Nodes))
	}
}

func TestBuild_IDsAndOrder(t *testing.T) {
	scene, err := Build(loadDiagramFixture(t, "diagram/aws-network.json"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !sort.SliceIsSorted(scene.Nodes, func(i, j int) bool { return scene.Nodes[i].Address < scene.Nodes[j].Address }) {
		t.Error("nodes not sorted by address")
	}
	for i, n := range scene.Nodes {
		if want := "n" + itoa(i+1); n.ID != want {
			t.Errorf("node %s id = %s, want %s", n.Address, n.ID, want)
		}
	}
	if scene.Containers[0].ID != "c1" {
		t.Errorf("first container id = %s", scene.Containers[0].ID)
	}
	if scene.Empty {
		t.Error("scene must not be empty")
	}
	want := []Status{StatusCreate, StatusUpdate, StatusDelete, StatusReplace, StatusNoop}
	if !reflect.DeepEqual(scene.Statuses, want) {
		t.Errorf("statuses = %v, want %v (vpc is drawn as no-op context)", scene.Statuses, want)
	}
}

func TestBuild_NodeFields(t *testing.T) {
	scene, _ := Build(planOf(rc(`module.app[0].aws_instance.web["a"]`, "registry.terraform.io/hashicorp/aws", "update")), Options{})
	n := scene.Nodes[0]
	if n.Address != `module.app[0].aws_instance.web["a"]` || n.ConfigAddress != "module.app.aws_instance.web" ||
		n.Name != `web["a"]` || n.Type != "aws_instance" || n.Provider != "aws" || n.Status != StatusUpdate {
		t.Errorf("node = %+v", n)
	}
	if len(n.Family) == 0 || n.Family[len(n.Family)-1] != "generic" {
		t.Errorf("family = %v", n.Family)
	}
}

func TestBuild_Empty(t *testing.T) {
	for name, plan := range map[string]*planType{"no changes": planOf(), "only no-op": planOf(rc("aws_a.n", "aws", "no-op"))} {
		t.Run(name, func(t *testing.T) {
			scene, err := Build(plan, Options{})
			if err != nil {
				t.Fatal(err)
			}
			if !scene.Empty || len(scene.Nodes) != 0 || len(scene.Statuses) != 0 {
				t.Errorf("scene = %+v", scene)
			}
		})
	}
}
