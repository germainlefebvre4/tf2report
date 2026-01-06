---
sidebar_position: 5
---

# API Reference

Complete API documentation for tf2report packages.

## Overview

tf2report consists of three main packages:

- `pkg/terraform` - Terraform plan parsing and analysis
- `pkg/report` - Report generation and formatting
- `pkg/config` - Configuration management

## pkg/terraform

Parse and analyze Terraform plans.

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

Represents the before/after state of a resource.

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

Terraform action types.

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

Aggregated change statistics.

#### TypeSummary

```go
type TypeSummary struct {
    Type      string
    ToAdd     int
    ToChange  int
    ToDestroy int
    ToReplace int
    Resources []string
}
```

Change counts for a specific resource type.

### Functions

#### NewParser

```go
func NewParser() *Parser
```

Creates a new Parser instance.

**Example:**
```go
parser := terraform.NewParser()
```

#### Parser.ParseFile

```go
func (p *Parser) ParseFile(path string) (*Plan, error)
```

Parses a Terraform plan JSON file.

**Example:**
```go
plan, err := parser.ParseFile("terraform.tfplan.json")
if err != nil {
    return err
}
```

#### Parser.Parse

```go
func (p *Parser) Parse(data []byte) (*Plan, error)
```

Parses Terraform plan JSON from bytes.

**Example:**
```go
data, _ := os.ReadFile("plan.json")
plan, err := parser.Parse(data)
```

#### NewSummary

```go
func NewSummary(changes []ResourceChange) *Summary
```

Creates a Summary from resource changes.

**Example:**
```go
summary := terraform.NewSummary(plan.ResourceChanges)
fmt.Printf("Total changes: %d\n", summary.TotalChanges)
```

#### FilterByType

```go
func FilterByType(changes []ResourceChange, types []string) []ResourceChange
```

Filters resource changes by type.

**Example:**
```go
filtered := terraform.FilterByType(
    plan.ResourceChanges,
    []string{"aws_instance", "aws_s3_bucket"},
)
```

#### FilterByAction

```go
func FilterByAction(changes []ResourceChange, actions []ActionType) []ResourceChange
```

Filters resource changes by action.

**Example:**
```go
filtered := terraform.FilterByAction(
    plan.ResourceChanges,
    []terraform.ActionType{terraform.ActionCreate, terraform.ActionDelete},
)
```

## pkg/report

Generate formatted reports from Terraform plans.

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

Output format types.

#### Formatter

```go
type Formatter interface {
    Format(plan *terraform.Plan, summary *terraform.Summary, w io.Writer) error
}
```

Interface for report formatters.

### Functions

#### NewFormatter

```go
func NewFormatter(format Format) (Formatter, error)
```

Creates a formatter for the specified format.

**Example:**
```go
formatter, err := report.NewFormatter(report.FormatMarkdown)
if err != nil {
    return err
}
```

#### ParseFormat

```go
func ParseFormat(s string) (Format, error)
```

Converts a string to a Format type.

**Example:**
```go
format, err := report.ParseFormat("markdown")
```

### Usage Example

```go
import (
    "os"
    "github.com/germainlefebvre4/tf2report/pkg/terraform"
    "github.com/germainlefebvre4/tf2report/pkg/report"
)

// Parse plan
parser := terraform.NewParser()
plan, _ := parser.ParseFile("terraform.tfplan.json")

// Generate summary
summary := terraform.NewSummary(plan.ResourceChanges)

// Create formatter
formatter, _ := report.NewFormatter(report.FormatMarkdown)

// Generate report
formatter.Format(plan, summary, os.Stdout)
```

## pkg/config

Configuration management.

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

Application configuration.

#### Filters

```go
type Filters struct {
    ResourceTypes []string
    Actions       []string
}
```

Resource filtering options.

### Functions

#### NewConfig

```go
func NewConfig() *Config
```

Creates a Config with defaults.

**Example:**
```go
cfg := config.NewConfig()
```

#### NewLoader

```go
func NewLoader() *Loader
```

Creates a configuration loader.

**Example:**
```go
loader := config.NewLoader()
```

#### Loader.Load

```go
func (l *Loader) Load(configPath string) (*Config, error)
```

Loads configuration from file and environment.

**Example:**
```go
cfg, err := loader.Load("tf2report.yaml")
```

#### Filters.GetActionTypes

```go
func (f *Filters) GetActionTypes() []terraform.ActionType
```

Converts string actions to ActionType.

**Example:**
```go
actionTypes := filters.GetActionTypes()
```

## Complete Example

```go
package main

import (
    "os"
    "github.com/germainlefebvre4/tf2report/pkg/config"
    "github.com/germainlefebvre4/tf2report/pkg/terraform"
    "github.com/germainlefebvre4/tf2report/pkg/report"
)

func main() {
    // Load configuration
    loader := config.NewLoader()
    cfg, err := loader.Load("tf2report.yaml")
    if err != nil {
        panic(err)
    }

    // Parse plan
    parser := terraform.NewParser()
    plan, err := parser.ParseFile(cfg.TerraformPlanPath)
    if err != nil {
        panic(err)
    }

    // Apply filters
    changes := plan.ResourceChanges
    if len(cfg.Filters.ResourceTypes) > 0 {
        changes = terraform.FilterByType(changes, cfg.Filters.ResourceTypes)
    }
    if len(cfg.Filters.Actions) > 0 {
        changes = terraform.FilterByAction(changes, cfg.Filters.GetActionTypes())
    }

    // Update plan with filtered changes
    plan.ResourceChanges = changes

    // Generate summary
    summary := terraform.NewSummary(changes)

    // Create formatter
    format, _ := report.ParseFormat(cfg.OutputFormat)
    formatter, err := report.NewFormatter(format)
    if err != nil {
        panic(err)
    }

    // Generate report
    formatter.Format(plan, summary, os.Stdout)
}
```

## Error Handling

All functions return errors that should be checked. Errors are wrapped using `fmt.Errorf` with `%w` for unwrapping support.

## Thread Safety

- **Parsers** - Stateless, safe for concurrent use
- **Formatters** - Stateless, safe for concurrent use
- **Loaders** - Not thread-safe, create new instance per goroutine

## Next Steps

- **[Contributing](./contributing.md)** - Contribution guidelines
- **[Getting Started](./getting-started.md)** - Development setup
