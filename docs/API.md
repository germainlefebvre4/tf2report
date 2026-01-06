# tf2report API Reference

This document describes the packages and public APIs in the `tf2report` project for developers who want to use the libraries programmatically or contribute to the project.

## Package Structure

```
github.com/germainlefebvre4/tf2report
├── cmd/tf2report          # CLI application entry point
├── pkg/config             # Configuration loading
├── pkg/report             # Report generation and formatting
└── pkg/terraform          # Terraform plan parsing
```

## pkg/terraform

Package `terraform` provides functionality for parsing and analyzing Terraform plan files.

### Types

#### Plan

```go
type Plan struct {
    FormatVersion    string
    TerraformVersion string
    ResourceChanges  []ResourceChange
    OutputChanges    map[string]Change
}
```

Represents a Terraform plan file structure.

#### ResourceChange

```go
type ResourceChange struct {
    Address       string
    Mode          string
    Type          string
    Name          string
    ProviderName  string
    Change        Change
    ActionReason  string
    PreviousAddr  string
}
```

Represents a single resource change in the plan.

#### Change

```go
type Change struct {
    Actions         []string
    Before          map[string]interface{}
    After           map[string]interface{}
    AfterUnknown    map[string]interface{}
    BeforeSensitive map[string]interface{}
    AfterSensitive  map[string]interface{}
}
```

Represents the before and after state of a resource.

#### ActionType

```go
type ActionType string

const (
    ActionCreate  ActionType = "create"
    ActionDelete  ActionType = "delete"
    ActionUpdate  ActionType = "update"
    ActionReplace ActionType = "replace"
    ActionRead    ActionType = "read"
    ActionNoOp    ActionType = "no-op"
)
```

Represents possible Terraform actions.

#### Summary

```go
type Summary struct {
    TotalChanges int
    ToAdd        int
    ToChange     int
    ToDestroy    int
    ToReplace    int
    ByType       map[string]TypeSummary
}
```

Contains aggregated statistics about plan changes.

#### TypeSummary

```go
type TypeSummary struct {
    Type      string
    ToAdd     int
    ToChange  int
    ToDestroy int
    ToReplace int
}
```

Contains change counts for a specific resource type.

### Functions

#### NewParser

```go
func NewParser() *Parser
```

Creates a new Parser instance for parsing Terraform plans.

#### Parser.ParseFile

```go
func (p *Parser) ParseFile(path string) (*Plan, error)
```

Reads and parses a Terraform plan JSON file from the given path.

#### Parser.Parse

```go
func (p *Parser) Parse(data []byte) (*Plan, error)
```

Parses Terraform plan JSON data from bytes.

#### Change.GetPrimaryAction

```go
func (c *Change) GetPrimaryAction() ActionType
```

Returns the primary action for a change. Handles multi-action scenarios like delete+create for replacements.

#### NewSummary

```go
func NewSummary(changes []ResourceChange) *Summary
```

Creates a Summary from a slice of ResourceChanges.

#### FilterByType

```go
func FilterByType(changes []ResourceChange, types []string) []ResourceChange
```

Filters resource changes by resource type.

#### FilterByAction

```go
func FilterByAction(changes []ResourceChange, actions []ActionType) []ResourceChange
```

Filters resource changes by action type.

## pkg/report

Package `report` provides report formatting functionality.

### Types

#### Format

```go
type Format string

const (
    FormatMarkdown Format = "markdown"
    FormatText     Format = "text"
    FormatJSON     Format = "json"
)
```

Represents the output format type.

#### Formatter

```go
type Formatter interface {
    Format(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error
}
```

Interface for generating reports from Terraform plan data.

### Functions

#### NewFormatter

```go
func NewFormatter(format Format) (Formatter, error)
```

Creates a formatter based on the specified format.

#### ParseFormat

```go
func ParseFormat(s string) (Format, error)
```

Converts a string to a Format type. Supports "markdown", "md", "text", "txt", and "json".

## pkg/config

Package `config` handles configuration loading and management.

### Types

#### Config

```go
type Config struct {
    TerraformPlanPath string
    OutputFormat      string
    Filters           Filters
    Verbosity         string
}
```

Represents the application configuration.

#### Filters

```go
type Filters struct {
    ResourceTypes []string
    Actions       []string
}
```

Contains filtering options for resources.

### Functions

#### NewConfig

```go
func NewConfig() *Config
```

Creates a Config with default values.

#### NewLoader

```go
func NewLoader() *Loader
```

Creates a new configuration loader.

#### Loader.Load

```go
func (l *Loader) Load(configPath string) (*Config, error)
```

Reads configuration from a file and environment variables.

#### Filters.GetActionTypes

```go
func (f *Filters) GetActionTypes() []terraform.ActionType
```

Converts string actions to ActionType slice.

## Usage Examples

### Parse a Plan File

```go
import "github.com/germainlefebvre4/tf2report/pkg/terraform"

parser := terraform.NewParser()
plan, err := parser.ParseFile("terraform.tfplan.json")
if err != nil {
    return err
}
```

### Generate a Summary

```go
summary := terraform.NewSummary(plan.ResourceChanges)
fmt.Printf("Total changes: %d\n", summary.TotalChanges)
```

### Filter Resources

```go
filtered := terraform.FilterByType(
    plan.ResourceChanges,
    []string{"aws_instance", "aws_s3_bucket"},
)

filtered = terraform.FilterByAction(
    filtered,
    []terraform.ActionType{terraform.ActionCreate},
)
```

### Generate a Report

```go
import (
    "os"
    "github.com/germainlefebvre4/tf2report/pkg/report"
)

formatter, err := report.NewFormatter(report.FormatMarkdown)
if err != nil {
    return err
}

err = formatter.Format(plan, summary, os.Stdout)
```

### Load Configuration

```go
import "github.com/germainlefebvre4/tf2report/pkg/config"

loader := config.NewLoader()
cfg, err := loader.Load("config.yaml")
if err != nil {
    return err
}
```

## Error Handling

All parsing and formatting functions return errors that should be checked. Errors are wrapped with context using `fmt.Errorf` with the `%w` verb, allowing error unwrapping with `errors.Unwrap`.

## Thread Safety

The parser and formatters are stateless and safe for concurrent use. The configuration loader creates a new Viper instance per loader, so each Loader instance should not be shared across goroutines.
