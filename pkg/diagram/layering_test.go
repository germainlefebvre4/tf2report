package diagram

import (
	"reflect"
	"testing"
)

func mkNodes(addrs ...string) []*Node {
	out := make([]*Node, len(addrs))
	for i, a := range addrs {
		out[i] = &Node{Address: a}
	}
	return out
}

func dep(from, to *Node) Edge { return Edge{From: from, To: to} }

func layerAddrs(cols [][]*Node) [][]string {
	out := [][]string{}
	for _, c := range cols {
		var l []string
		for _, n := range c {
			l = append(l, n.Address)
		}
		out = append(out, l)
	}
	return out
}

func TestLayerNodes(t *testing.T) {
	n := mkNodes("a", "b", "c", "d", "e")
	a, b, c, d, e := n[0], n[1], n[2], n[3], n[4]
	tests := []struct {
		name  string
		nodes []*Node
		edges []Edge
		want  [][]string
	}{
		{"dependency in a lower column", n[:2], []Edge{dep(a, b)}, [][]string{{"b"}, {"a"}}},
		{"chain", n[:3], []Edge{dep(a, b), dep(b, c)}, [][]string{{"c"}, {"b"}, {"a"}}},
		{"alphabetical within a layer", n[:4], []Edge{dep(a, d), dep(b, d), dep(c, d)}, [][]string{{"d"}, {"a", "b", "c"}}},
		{"longest path", n[:3], []Edge{dep(a, b), dep(a, c), dep(b, c)}, [][]string{{"c"}, {"b"}, {"a"}}},
		{"barycenter reorders", []*Node{a, b, c, d}, []Edge{dep(a, d), dep(b, c)}, [][]string{{"c", "d"}, {"b", "a"}}},
		{"isolated node in first column", []*Node{a, b, e}, []Edge{dep(a, b)}, [][]string{{"b", "e"}, {"a"}}},
		{"edge to container ignored", n[:2], []Edge{{From: a, ToContainer: &Container{}}}, [][]string{{"a", "b"}}},
		{"edge leaving the container ignored", n[:1], []Edge{dep(a, b)}, [][]string{{"a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := layerAddrs(layerNodes(tt.nodes, tt.edges))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("layers = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayerNodes_CycleTerminatesDeterministically(t *testing.T) {
	first := func() [][]string {
		n := mkNodes("a", "b", "c")
		return layerAddrs(layerNodes(n, []Edge{dep(n[0], n[1]), dep(n[1], n[0]), dep(n[2], n[0])}))
	}
	a, b := first(), first()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("runs differ: %v vs %v", a, b)
	}
	total := 0
	for _, l := range a {
		total += len(l)
	}
	if total != 3 {
		t.Errorf("all nodes must be layered: %v", a)
	}
}

func TestLayout_CycleFixtureKeepsBothEdges(t *testing.T) {
	s := laidOut(t, loadDiagramFixture(t, "diagram/cycle.json"), Options{})
	if len(s.Edges) != 2 {
		t.Fatalf("edges = %v", edgeSet(s))
	}
	a, b := nodeByAddress(s, "aws_security_group.a"), nodeByAddress(s, "aws_security_group.b")
	if a.X == b.X {
		t.Error("cycle members must sit in different columns")
	}
	for _, e := range s.Edges {
		if len(e.Points) != 4 {
			t.Errorf("edge %s has %d points", e.From.Address, len(e.Points))
		}
	}
}

func TestLayout_LayeredColumnsAndGridFallback(t *testing.T) {
	s := laidOut(t, loadDiagramFixture(t, "diagram/aws-network.json"), Options{})
	sub, vpc := nodeByAddress(s, "aws_subnet.private"), nodeByAddress(s, "aws_vpc.main")
	if vpc.X >= sub.X {
		t.Errorf("dependency vpc (x=%d) must be left of subnet (x=%d)", vpc.X, sub.X)
	}
	// a container without intra edge keeps the grid
	g := laidOut(t, manyCreates(4), Options{})
	if g.Nodes[1].X != g.Containers[0].X+ContainerPad+NodeW+ColGap {
		t.Error("grid expected without intra-container edges")
	}
}

func TestRouting(t *testing.T) {
	s := laidOut(t, loadDiagramFixture(t, "diagram/aws-network.json"), Options{})
	for _, e := range s.Edges {
		p := e.Points
		if len(p) != 4 {
			t.Fatalf("points = %v", p)
		}
		if p[0].X != e.From.X || p[0].Y != e.From.Y+e.From.H/2 {
			t.Errorf("edge must start at the left side of %s: %v", e.From.Address, p[0])
		}
		if e.To != nil && (p[3].X != e.To.X+e.To.W || p[3].Y != e.To.Y+e.To.H/2) {
			t.Errorf("edge must end at the right side of %s: %v", e.To.Address, p[3])
		}
		if p[1].Y != p[0].Y || p[2].Y != p[3].Y || p[1].X >= p[0].X || p[2].X <= p[3].X {
			t.Errorf("horizontal tangents expected: %v", p)
		}
	}
}

func TestRouting_ContainerTarget(t *testing.T) {
	s := laidOut(t, loadDiagramFixture(t, "diagram/modules-context.json"), Options{})
	found := false
	for _, e := range s.Edges {
		if e.ToContainer != nil {
			found = true
			c := e.ToContainer
			if end := e.Points[3]; end.X != c.X+c.W || end.Y != c.Y+c.H/2 {
				t.Errorf("container edge end = %v", end)
			}
		}
	}
	if !found {
		t.Error("expected a container edge in modules-context")
	}
}
