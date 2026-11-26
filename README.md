# tf2report

Summarize Terraform plan changes into human-readable reports.

## Overview

`tf2report` is a command-line tool that analyzes Terraform plan files and generates resource change reports in multiple formats. It helps you understand and communicate infrastructure changes before applying them.

## Features

- Parse Terraform plan files in JSON format
- Extract and summarize resource changes
- Multiple output formats: Markdown, plain text, JSON, SVG diagram
- SVG architecture diagram with provider and module containers, dependency arrows, status legend and icons
- Filter by resource type and change action
- Configuration file support
- CI/CD integration friendly

## Installation

### From Source

```bash
git clone https://github.com/germainlefebvre4/tf2report.git
cd tf2report
make build
sudo make install
```

### Using Go

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

## Quick Start

### Generate a Terraform Plan

```bash
terraform plan -out=tfplan
terraform show -json tfplan > terraform.tfplan.json
```

### Generate a Report

```bash
tf2report --plan terraform.tfplan.json
```

### Example Output

```markdown
# Terraform Plan Summary

**Terraform Version:** 1.10.5

## Summary

**Total Changes:** 4

| Action | Count |
|--------|-------|
| Add | 1 |
| Change | 1 |
| Destroy | 1 |
| Replace | 1 |

## Changes by Resource Type

### aws_instance (2)

**To Change (1):**

- `aws_instance.web[0]`

**To Replace (1):**

- `aws_instance.web[1]`
```

## Usage

### Basic Usage

```bash
tf2report --plan <path-to-plan.json>
```

### Output Formats

```bash
# Markdown (default)
tf2report --plan terraform.tfplan.json

# Plain text
tf2report --plan terraform.tfplan.json --format text

# JSON
tf2report --plan terraform.tfplan.json --format json

# SVG diagram
tf2report --plan terraform.tfplan.json --format svg > diagram.svg
```

### SVG diagram

Each status is shown by color, symbol and border style, never by color alone:

| Status | Color | Symbol | Border |
|--------|-------|--------|--------|
| create | green | `+` | solid |
| update | orange | `~` | dashed |
| delete | red | `-` | dotted |
| replace | purple | `-/+` | dash-dot |
| no-op | grey, muted | none | thin solid |

What is drawn depends on the filters:

| `--type` | `--action` | create, update, delete, replace | no-op | Dependency context |
|----------|------------|---------------------------------|-------|--------------------|
| not given | not given | all | no | yes, one hop, muted |
| given | not given | matching types | no | no |
| any | given | matching, among those left after filtering | only if `no-op` is in the list | no |

Use `--icons-dir <dir>` to supply your own icons (see `docs/USAGE.md`).

Limits:

- No isometric rendering and no PNG export: the output is one 2D SVG document.
- A reference to a module output points to the module container, not to the producing resource.
- Deleted resources have no arrows, because their relations are no longer in the configuration.
- Data sources are never drawn. Arrows towards resources that are not drawn are omitted.
- The diagram carries addresses, types, names, providers and statuses only, never attribute values.

### Filtering

```bash
# Filter by resource type
tf2report --plan terraform.tfplan.json --type aws_instance --type aws_s3_bucket

# Filter by action
tf2report --plan terraform.tfplan.json --action create --action delete

# Combine filters
tf2report --plan terraform.tfplan.json --type aws_instance --action create
```

### Configuration File

Create `tf2report.yaml`:

```yaml
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
  actions:
    - create
    - delete
verbosity: info
```

Run with configuration:

```bash
tf2report
```

## Command-Line Options

- `--plan, -p <path>` - Path to Terraform plan JSON file (required)
- `--format, -f <format>` - Output format: markdown, text, json, svg (default: markdown)
- `--type, -t <type>` - Filter by resource type (repeatable)
- `--action, -a <action>` - Filter by action: create, update, delete, replace, no-op (repeatable; `no-op` is meaningful for svg)
- `--icons-dir <dir>` - Directory of SVG icons for the svg format (ignored for other formats)
- `--config <path>` - Configuration file path (default: ./tf2report.yaml)
- `--verbose, -v` - Enable verbose output
- `--help, -h` - Display help

## Documentation

- [Usage Guide](docs/USAGE.md) - Detailed usage instructions
- [Configuration](docs/CONFIGURATION.md) - Configuration file reference
- [Output Formats](docs/OUTPUT_FORMATS.md) - Details on each output format
- [Examples](docs/EXAMPLES.md) - Usage examples and integrations
- [API Reference](docs/API.md) - Package API documentation
- [Development](docs/DEVELOPMENT.md) - Contributing and development setup

## Use Cases

### Code Review

Generate reports for pull request comments to help reviewers understand infrastructure changes.

### CI/CD Integration

Integrate into your CI/CD pipeline to automatically generate and analyze Terraform plan reports.

### Change Documentation

Create documentation of planned infrastructure changes for compliance and audit purposes.

### Security Review

Filter changes to security-related resources for security team review.

## Requirements

- Go 1.25 or later (for building from source)
- Terraform plan files in JSON format

## Technology

- Go 1.25
- Cobra (CLI framework)
- Viper (configuration management)

## Contributing

Contributions are welcome. Please read the [Development Guide](docs/DEVELOPMENT.md) for details on the development process and coding standards.

## License

Copyright Germain LEFEBVRE. All rights reserved.
