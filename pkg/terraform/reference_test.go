package terraform

import "testing"

func TestParseReference(t *testing.T) {
	tests := []struct {
		in   string
		want Reference
	}{
		{"aws_subnet.a.id", Reference{Kind: RefResource, Type: "aws_subnet", Name: "a"}},
		{"aws_subnet.a[0].id", Reference{Kind: RefResource, Type: "aws_subnet", Name: "a"}},
		{"module.net.vpc_id", Reference{Kind: RefModule, ModuleCall: "net"}},
		{"module.net", Reference{Kind: RefModule, ModuleCall: "net"}},
		{"data.aws_ami.x.id", Reference{Kind: RefData}},
		{"var.x", Reference{Kind: RefOther}},
		{"local.x", Reference{Kind: RefOther}},
		{"count.index", Reference{Kind: RefOther}},
		{"each.key", Reference{Kind: RefOther}},
		{"path.module", Reference{Kind: RefOther}},
		{"terraform.workspace", Reference{Kind: RefOther}},
		{"self", Reference{Kind: RefOther}},
		{"module", Reference{Kind: RefOther}},
		{"garbage", Reference{Kind: RefOther}},
		{"", Reference{Kind: RefOther}},
		{"...", Reference{Kind: RefOther}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := ParseReference(tt.in); got != tt.want {
				t.Errorf("ParseReference(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}
