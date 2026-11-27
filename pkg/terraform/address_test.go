package terraform

import (
	"reflect"
	"testing"
)

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Address
	}{
		{"root resource", "aws_instance.web", Address{Mode: "managed", Type: "aws_instance", Name: "web"}},
		{"integer index", "aws_instance.web[0]", Address{Mode: "managed", Type: "aws_instance", Name: "web", Key: "[0]"}},
		{"quoted key with dots", `aws_instance.web["a.b"]`, Address{Mode: "managed", Type: "aws_instance", Name: "web", Key: `["a.b"]`}},
		{"quoted key with brackets", `aws_instance.web["a[1]"]`, Address{Mode: "managed", Type: "aws_instance", Name: "web", Key: `["a[1]"]`}},
		{
			"nested modules with keys",
			`module.app[0].module.db["eu.west"].aws_instance.web["a"]`,
			Address{
				Modules: []ModuleStep{{"app", "[0]"}, {"db", `["eu.west"]`}},
				Mode:    "managed", Type: "aws_instance", Name: "web", Key: `["a"]`,
			},
		},
		{"data source", "data.aws_ami.x", Address{Mode: "data", Type: "aws_ami", Name: "x"}},
		{"empty", "", Address{Mode: "managed"}},
		{"unbalanced bracket", "aws_instance.web[0", Address{Mode: "managed", Type: "aws_instance", Name: "web", Key: "[0"}},
		{"trailing dot", "aws_instance.", Address{Mode: "managed", Type: "aws_instance"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAddress(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAddress(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestConfigAddress(t *testing.T) {
	tests := []struct{ in, want string }{
		{`module.app[0].aws_instance.web["a"]`, "module.app.aws_instance.web"},
		{"aws_instance.web", "aws_instance.web"},
		{`module.a["x.y"].module.b[1].data.aws_ami.c`, "module.a.module.b.data.aws_ami.c"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := ParseAddress(tt.in).ConfigAddress(); got != tt.want {
				t.Errorf("ConfigAddress(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
