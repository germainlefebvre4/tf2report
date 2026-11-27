package main

import (
	"strings"
	"testing"
)

// The filters x statuses table of the CLI contract.
func TestCLI_SVG_Filters(t *testing.T) {
	tests := []struct {
		name      string
		fixture   string
		args      []string
		nodes     int
		context   int
		edges     int
		wantTitle []string
		absent    []string
	}{
		{"no filter", "diagram/modules-context.json", nil, 6, 1, 2,
			[]string{"module.network.aws_vpc.main"}, []string{"aws_route_table.rt"}},
		{"type filter", "diagram/modules-context.json", []string{"--type", "aws_subnet"}, 1, 0, 0,
			[]string{"module.network.aws_subnet.a"}, []string{"aws_vpc.main", "aws_instance"}},
		{"action create", "diagram/aws-network.json", []string{"--action", "create"}, 3, 0, 2,
			[]string{"aws_subnet.private"}, []string{"aws_vpc.main", "aws_instance.web", "aws_s3_bucket.old"}},
		{"action no-op", "diagram/aws-network.json", []string{"--action", "no-op"}, 1, 0, 0,
			[]string{"aws_vpc.main"}, []string{"aws_subnet.private"}},
		{"create and no-op", "diagram/aws-network.json", []string{"--action", "create,no-op"}, 4, 0, 3,
			[]string{"aws_vpc.main", "aws_subnet.private"}, []string{"aws_instance.web"}},
		{"type and action", "diagram/aws-network.json", []string{"--type", "aws_instance", "--action", "create"}, 2, 0, 0,
			[]string{"aws_instance.worker[0]"}, []string{"aws_subnet.private"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"-p", fixturePath(t, tt.fixture), "-f", "svg"}, tt.args...)
			res := runCLI(t, cliOpts{Args: args})
			if res.Code != 0 {
				t.Fatalf("exit %d: %s", res.Code, res.Stderr)
			}
			assertSVG(t, res.Stdout)
			if got := strings.Count(res.Stdout, `<g class="node `); got != tt.nodes {
				t.Errorf("nodes = %d, want %d", got, tt.nodes)
			}
			if got := strings.Count(res.Stdout, ` context"`); got != tt.context {
				t.Errorf("context nodes = %d, want %d", got, tt.context)
			}
			if got := strings.Count(res.Stdout, `marker-end="url(#arrow)"`); got != tt.edges {
				t.Errorf("arrows = %d, want %d", got, tt.edges)
			}
			for _, a := range tt.wantTitle {
				if !strings.Contains(res.Stdout, "<title>"+a+"</title>") {
					t.Errorf("missing node %s", a)
				}
			}
			for _, a := range tt.absent {
				if strings.Contains(res.Stdout, "<title>"+a) {
					t.Errorf("node %s must not be drawn", a)
				}
			}
		})
	}
}
