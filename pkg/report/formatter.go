package report

import (
	"fmt"
	"io"

	"github.com/germainlefebvre4/tf2report/pkg/terraform"
)

// Format represents the output format type.
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatText     Format = "text"
	FormatJSON     Format = "json"
)

// Formatter generates reports from Terraform plan data.
type Formatter interface {
	Format(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error
}

// NewFormatter creates a formatter based on the specified format.
func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatMarkdown:
		return &MarkdownFormatter{}, nil
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ParseFormat converts a string to a Format type.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "markdown", "md":
		return FormatMarkdown, nil
	case "text", "txt":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("invalid format: %s (supported: markdown, text, json)", s)
	}
}
