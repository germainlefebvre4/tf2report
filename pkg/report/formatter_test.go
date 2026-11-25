package report

import (
	"reflect"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr string
	}{
		{input: "markdown", want: FormatMarkdown},
		{input: "md", want: FormatMarkdown},
		{input: "text", want: FormatText},
		{input: "txt", want: FormatText},
		{input: "json", want: FormatJSON},
		{input: "xml", wantErr: "invalid format: xml (supported: markdown, text, json)"},
		{input: "", wantErr: "invalid format:  (supported: markdown, text, json)"},
		{input: "MARKDOWN", wantErr: "invalid format: MARKDOWN (supported: markdown, text, json)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Errorf("ParseFormat(%q) error = %v, want %q", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseFormat(%q): %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		format Format
		want   Formatter
	}{
		{FormatMarkdown, &MarkdownFormatter{}},
		{FormatText, &TextFormatter{}},
		{FormatJSON, &JSONFormatter{}},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			got, err := NewFormatter(tt.format)
			if err != nil {
				t.Fatalf("NewFormatter: %v", err)
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewFormatter(%q) = %T, want %T", tt.format, got, tt.want)
			}
		})
	}

	t.Run("unsupported format", func(t *testing.T) {
		_, err := NewFormatter(Format("xml"))
		if err == nil || err.Error() != "unsupported format: xml" {
			t.Errorf("error = %v, want %q", err, "unsupported format: xml")
		}
	})
}
