package diagram

import (
	"sort"
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// cfgPlan builds a plan whose root module holds resources with the given references.
func cfgPlan(changes []terraform.ResourceChange, root terraform.ConfigModule) *terraform.Plan {
	p := planOf(changes...)
	p.Configuration = &terraform.Configuration{RootModule: root}
	return p
}

func cres(typ, name string, refs ...string) terraform.ConfigResource {
	list := make([]any, len(refs))
	for i, r := range refs {
		list[i] = r
	}
	return terraform.ConfigResource{
		Address: typ + "." + name, Mode: "managed", Type: typ, Name: name,
		Expressions: map[string]any{"x": map[string]any{"references": list}},
	}
}

// edgeSet lists edges as "from->to" (container targets as "from->[container path]").
func edgeSet(s *Scene) []string {
	out := []string{}
	for _, e := range s.Edges {
		if e.To != nil {
			out = append(out, e.From.Address+"->"+e.To.Address)
		} else {
			out = append(out, e.From.Address+"->["+e.ToContainer.Path+"]")
		}
	}
	sort.Strings(out)
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func nodeByAddress(s *Scene, addr string) *Node {
	for _, n := range s.Nodes {
		if n.Address == addr {
			return n
		}
	}
	return nil
}

func TestEdges_References(t *testing.T) {
	root := terraform.ConfigModule{Resources: []terraform.ConfigResource{
		cres("aws_subnet", "s", "aws_vpc.v.id"),
		cres("aws_instance", "i", "aws_subnet.s.id"),
		{Address: "aws_s3_bucket.b", Mode: "managed", Type: "aws_s3_bucket", Name: "b", DependsOn: []string{"aws_instance.i"}},
		cres("aws_vpc", "v"),
	}}
	changes := []terraform.ResourceChange{
		rc("aws_vpc.v", awsProv, "create"), rc("aws_subnet.s", awsProv, "create"),
		rc("aws_instance.i", awsProv, "update"), rc("aws_s3_bucket.b", awsProv, "create"),
	}
	s, _ := Build(cfgPlan(changes, root), Options{})
	want := []string{"aws_instance.i->aws_subnet.s", "aws_s3_bucket.b->aws_instance.i", "aws_subnet.s->aws_vpc.v"}
	if got := edgeSet(s); !equalStrings(got, want) {
		t.Errorf("edges = %v, want %v", got, want)
	}
}

func TestEdges_InstanceExpansionAndDuplicates(t *testing.T) {
	root := terraform.ConfigModule{Resources: []terraform.ConfigResource{
		cres("aws_instance", "w", "aws_subnet.s.id", "aws_subnet.s", "aws_subnet.s.arn"),
		cres("aws_subnet", "s"),
	}}
	changes := []terraform.ResourceChange{
		rc("aws_instance.w[0]", awsProv, "create"), rc("aws_instance.w[1]", awsProv, "create"),
		rc("aws_subnet.s[0]", awsProv, "create"), rc("aws_subnet.s[1]", awsProv, "create"),
	}
	s, _ := Build(cfgPlan(changes, root), Options{})
	if len(s.Edges) != 4 {
		t.Errorf("want 4 edges (2x2, no duplicates), got %v", edgeSet(s))
	}
}

func TestEdges_ModuleReferenceTargetsEveryContainer(t *testing.T) {
	root := terraform.ConfigModule{
		Resources: []terraform.ConfigResource{cres("aws_instance", "app", "module.net.vpc_id")},
		ModuleCalls: map[string]terraform.ModuleCall{
			"net": {Module: terraform.ConfigModule{Resources: []terraform.ConfigResource{cres("aws_vpc", "v")}}},
		},
	}
	changes := []terraform.ResourceChange{
		rc("aws_instance.app", awsProv, "create"),
		rc("module.net[0].aws_vpc.v", awsProv, "create"), rc("module.net[1].aws_vpc.v", awsProv, "create"),
	}
	s, _ := Build(cfgPlan(changes, root), Options{})
	want := []string{"aws_instance.app->[module.net[0]]", "aws_instance.app->[module.net[1]]"}
	if got := edgeSet(s); !equalStrings(got, want) {
		t.Errorf("edges = %v, want %v", got, want)
	}
}

func TestEdges_DeleteNodesHaveNoEdges(t *testing.T) {
	root := terraform.ConfigModule{Resources: []terraform.ConfigResource{
		cres("aws_instance", "i", "aws_subnet.s"), cres("aws_subnet", "s", "aws_vpc.v"), cres("aws_vpc", "v"),
	}}
	changes := []terraform.ResourceChange{
		rc("aws_instance.i", awsProv, "delete"), rc("aws_subnet.s", awsProv, "create"), rc("aws_vpc.v", awsProv, "delete"),
	}
	s, _ := Build(cfgPlan(changes, root), Options{})
	if len(s.Edges) != 0 {
		t.Errorf("edges = %v, want none", edgeSet(s))
	}
}

func TestEdges_DataAndOtherReferencesIgnored(t *testing.T) {
	root := terraform.ConfigModule{Resources: []terraform.ConfigResource{
		cres("aws_instance", "i", "data.aws_ami.x.id", "var.a", "local.b", "count.index", "aws_ghost.z"),
	}}
	s, _ := Build(cfgPlan([]terraform.ResourceChange{rc("aws_instance.i", awsProv, "create")}, root), Options{})
	if len(s.Edges) != 0 {
		t.Errorf("edges = %v", edgeSet(s))
	}
}

func TestEdges_NoConfiguration(t *testing.T) {
	s, _ := Build(planOf(rc("aws_instance.i", awsProv, "create")), Options{})
	if len(s.Edges) != 0 {
		t.Error("a plan without configuration must have no edges")
	}
}

func TestEdges_ContextNodes(t *testing.T) {
	root := terraform.ConfigModule{Resources: []terraform.ConfigResource{
		cres("aws_subnet", "a", "aws_vpc.v"), cres("aws_subnet", "b", "aws_vpc.v"),
		cres("aws_vpc", "v", "aws_route_table.rt"), cres("aws_route_table", "rt"), cres("aws_route_table", "lonely"),
	}}
	changes := []terraform.ResourceChange{
		rc("aws_subnet.a", awsProv, "create"), rc("aws_subnet.b", awsProv, "update"),
		rc("aws_vpc.v", awsProv, "no-op"), rc("aws_route_table.rt", awsProv, "no-op"),
		rc("aws_route_table.lonely", awsProv, "no-op"),
	}
	plan := cfgPlan(changes, root)

	t.Run("no filter adds one-hop context once", func(t *testing.T) { assertOneHopContext(t, plan) })
	t.Run("type filter adds no context", func(t *testing.T) {
		s, _ := Build(plan, Options{TypeFiltered: true})
		if nodeByAddress(s, "aws_vpc.v") != nil || len(s.Edges) != 0 {
			t.Errorf("nodes=%d edges=%v", len(s.Nodes), edgeSet(s))
		}
	})
	t.Run("action filter adds no context and drops dangling edges", func(t *testing.T) {
		filtered := cfgPlan(changes[:2], root)
		s, _ := Build(filtered, Options{ActionFiltered: true})
		if len(s.Nodes) != 2 || len(s.Edges) != 0 {
			t.Errorf("nodes=%d edges=%v", len(s.Nodes), edgeSet(s))
		}
	})
	t.Run("requested no-op is a normal node", func(t *testing.T) { assertRequestedNoop(t, plan) })
}

// assertOneHopContext checks that an unchanged referenced resource is drawn once as context, with its
// dependents' edges, and that its own dependencies are not followed.
func assertOneHopContext(t *testing.T, plan *terraform.Plan) {
	t.Helper()

	s, _ := Build(plan, Options{})
	vpc := nodeByAddress(s, "aws_vpc.v")
	if vpc == nil || !vpc.Context {
		t.Fatalf("vpc must be a context node: %+v", vpc)
	}
	if nodeByAddress(s, "aws_route_table.rt") != nil || nodeByAddress(s, "aws_route_table.lonely") != nil {
		t.Error("context must not be followed, and unreferenced no-op must not be drawn")
	}
	want := []string{"aws_subnet.a->aws_vpc.v", "aws_subnet.b->aws_vpc.v"}
	if got := edgeSet(s); !equalStrings(got, want) {
		t.Errorf("edges = %v, want %v", got, want)
	}
	count := 0
	for _, n := range s.Nodes {
		if n.Address == "aws_vpc.v" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("vpc drawn %d times", count)
	}
}

func assertRequestedNoop(t *testing.T, plan *terraform.Plan) {
	t.Helper()

	s, _ := Build(plan, Options{ActionFiltered: true})
	vpc := nodeByAddress(s, "aws_vpc.v")
	if vpc == nil || vpc.Context || vpc.Status != StatusNoop {
		t.Fatalf("vpc = %+v", vpc)
	}
	if nodeByAddress(s, "aws_route_table.rt") == nil {
		t.Error("all no-op resources must be drawn with an action filter")
	}
}

func TestEdges_FixtureModulesContext(t *testing.T) {
	s, err := Build(loadDiagramFixture(t, "diagram/modules-context.json"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	vpc := nodeByAddress(s, "module.network.aws_vpc.main")
	if vpc == nil || !vpc.Context {
		t.Fatalf("module vpc must be context: %+v", vpc)
	}
	if nodeByAddress(s, "aws_route_table.rt") != nil {
		t.Error("unreferenced no-op must not be drawn")
	}
	foundModuleEdge := false
	for _, e := range s.Edges {
		if e.From.Address == "aws_instance.app" && e.ToContainer != nil && e.ToContainer.Path == "module.network" {
			foundModuleEdge = true
		}
		if e.From.Context {
			t.Errorf("edge starts at context node %s", e.From.Address)
		}
	}
	if !foundModuleEdge {
		t.Errorf("missing module output edge: %v", edgeSet(s))
	}
}
