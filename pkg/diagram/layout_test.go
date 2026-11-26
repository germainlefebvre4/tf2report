package diagram

import (
	"testing"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

const awsProv = "registry.terraform.io/hashicorp/aws"

func laidOut(t *testing.T, plan *planType, opts Options) *Scene {
	t.Helper()

	s, err := Build(plan, opts)
	if err != nil {
		t.Fatal(err)
	}
	Layout(s)
	return s
}

func manyCreates(n int) *planType {
	var changes []terraform.ResourceChange
	for i := 0; i < n; i++ {
		changes = append(changes, rc("aws_s3_bucket.b"+pad(i), awsProv, "create"))
	}
	return planOf(changes...)
}

func pad(i int) string { return string(rune('a'+i/26)) + string(rune('a'+i%26)) }

func TestLayout_GridFallback(t *testing.T) {
	tests := []struct{ n, cols int }{{1, 1}, {2, 2}, {4, 2}, {5, 3}, {9, 3}, {10, 4}}
	for _, tt := range tests {
		t.Run(itoa(tt.n), func(t *testing.T) {
			s := laidOut(t, manyCreates(tt.n), Options{})
			c := s.Containers[0]
			if got := ceilSqrt(tt.n); got != tt.cols {
				t.Fatalf("ceilSqrt(%d) = %d, want %d", tt.n, got, tt.cols)
			}
			for i, n := range c.Nodes {
				wantX := c.X + ContainerPad + (i%tt.cols)*(NodeW+ColGap)
				wantY := c.Y + HeaderH + ContainerPad + (i/tt.cols)*(NodeH+RowGap)
				if n.X != wantX || n.Y != wantY {
					t.Errorf("node %d at (%d,%d), want (%d,%d)", i, n.X, n.Y, wantX, wantY)
				}
			}
		})
	}
}

func TestLayout_ContainerSizeFromChildren(t *testing.T) {
	plan := planOf(
		rc("aws_s3_bucket.root", awsProv, "create"),
		rc("module.m.aws_s3_bucket.in", awsProv, "create"),
	)
	s := laidOut(t, plan, Options{})
	p := s.Containers[0]
	m := p.Children[0]
	wantMH := HeaderH + NodeH + 2*ContainerPad
	if m.H != wantMH || m.W != NodeW+2*ContainerPad {
		t.Errorf("module size = %dx%d, want %dx%d", m.W, m.H, NodeW+2*ContainerPad, wantMH)
	}
	wantPH := HeaderH + NodeH + RowGap + m.H + 2*ContainerPad
	if p.H != wantPH {
		t.Errorf("provider H = %d, want %d", p.H, wantPH)
	}
	if m.Y < p.Nodes[0].Y+NodeH {
		t.Errorf("module container (y=%d) must sit below the provider nodes (bottom %d)", m.Y, p.Nodes[0].Y+NodeH)
	}
	if m.X != p.X+ContainerPad {
		t.Errorf("module x = %d", m.X)
	}
}

func TestLayout_ProvidersStackAlphabeticallyWithoutOverlap(t *testing.T) {
	plan := planOf(
		rc("random_pet.x", "registry.terraform.io/hashicorp/random", "create"),
		rc("aws_s3_bucket.b", awsProv, "create"),
		rc("google_storage_bucket.g", "registry.terraform.io/hashicorp/google", "create"),
	)
	s := laidOut(t, plan, Options{})
	labels := []string{}
	for i, c := range s.Containers {
		labels = append(labels, c.Label)
		if i > 0 && c.Y < s.Containers[i-1].Y+s.Containers[i-1].H {
			t.Errorf("container %s overlaps %s", c.Label, s.Containers[i-1].Label)
		}
		if c.X+c.W > s.Width || c.Y+c.H > s.Height {
			t.Errorf("container %s outside the scene box %dx%d", c.Label, s.Width, s.Height)
		}
	}
	if labels[0] != "aws" || labels[1] != "google" || labels[2] != "random" {
		t.Errorf("order = %v", labels)
	}
	for _, n := range s.Nodes {
		if n.X+n.W > s.Width || n.Y+n.H > s.Height {
			t.Errorf("node %s outside the scene box", n.Address)
		}
	}
}

func TestLayout_EmptyAndDeterministic(t *testing.T) {
	s := laidOut(t, planOf(), Options{})
	if !s.Empty || s.Width != 0 || s.Height != 0 {
		t.Errorf("empty scene = %+v", s)
	}
	plan := loadDiagramFixture(t, "diagram/modules-context.json")
	a, b := laidOut(t, plan, Options{}), laidOut(t, plan, Options{})
	if len(a.Nodes) != len(b.Nodes) {
		t.Fatal("different node counts")
	}
	for i := range a.Nodes {
		if a.Nodes[i].X != b.Nodes[i].X || a.Nodes[i].Y != b.Nodes[i].Y {
			t.Errorf("node %s differs between runs", a.Nodes[i].Address)
		}
	}
	if a.Width != b.Width || a.Height != b.Height {
		t.Error("scene size differs between runs")
	}
}
