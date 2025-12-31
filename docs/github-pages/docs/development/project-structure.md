---
sidebar_position: 2
---

# Project Structure

Detailed overview of tf2report's codebase organization.

## Directory Layout

```
tf2report/
├── cmd/
│   └── tf2report/          # CLI application
├── pkg/
│   ├── config/             # Configuration management
│   ├── report/             # Report formatters
│   └── terraform/          # Terraform plan parsing
├── docs/                   # Documentation
├── examples/               # Sample Terraform plans
├── .github/
│   ├── workflows/          # CI/CD pipelines
│   └── instructions/       # Coding guidelines
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Build automation
└── README.md               # Main documentation
```

## Package Overview

### `cmd/tf2report`

CLI application entry point.

**Key files:**
- `main.go` - Application entry, CLI setup, command handling

**Responsibilities:**
- Parse command-line arguments
- Load configuration
- Initialize components
- Execute report generation
- Handle errors and exit codes

### `pkg/config`

Configuration loading and management.

**Key files:**
- `types.go` - Configuration struct definitions
- `loader.go` - Configuration file loading logic

**Responsibilities:**
- Load YAML configuration files
- Merge CLI flags with config
- Validate configuration
- Provide defaults

### `pkg/terraform`

Terraform plan parsing and analysis.

**Key files:**
- `types.go` - Terraform plan data structures
- `parser.go` - Plan parsing and filtering logic

**Responsibilities:**
- Parse Terraform plan JSON
- Extract resource changes
- Filter by type and action
- Generate summary statistics

### `pkg/report`

Report generation and formatting.

**Key files:**
- `formatter.go` - Formatter interface and factory
- `markdown.go` - Markdown formatter
- `text.go` - Plain text formatter
- `json.go` - JSON formatter

**Responsibilities:**
- Format reports in multiple formats
- Generate human-readable output
- Structure data for machines

## Data Flow

```
terraform.tfplan.json
        ↓
[pkg/terraform/Parser]
        ↓
Plan struct
        ↓
[pkg/terraform/NewSummary]
        ↓
Summary struct
        ↓
[pkg/report/Formatter]
        ↓
Formatted output
```

## Key Concepts

### Plan Parsing

Terraform plan JSON → Go structs representing resources and changes.

### Filtering

Apply resource type and action filters to narrow results.

### Summary Generation

Aggregate change statistics by type and action.

### Formatting

Convert structured data to user-friendly formats.

## Configuration Flow

```
1. Load defaults
2. Load config file (YAML)
3. Apply environment variables
4. Apply CLI flags
5. Validate final config
```

## Next Steps

- **[Building](./building.md)** - Build instructions
- **[Testing](./testing.md)** - Testing guide
- **[API Reference](./api-reference.md)** - Detailed API docs
