---
sidebar_position: 1
---

# CLI Options Reference

Complete reference for all tf2report command-line options.

## Synopsis

```bash
tf2report [flags]
```

## Flags

### Required Flags

#### `--plan, -p <path>`

Path to Terraform plan JSON file.

**Type:** String
**Required:** Yes (unless set in config)
**Example:**
```bash
tf2report --plan terraform.tfplan.json
tf2report -p /path/to/plan.json
```

### Optional Flags

#### `--format, -f <format>`

Output format.

**Type:** String
**Default:** `markdown`
**Valid Values:** `markdown`, `text`, `json`
**Example:**
```bash
tf2report --plan plan.json --format markdown
tf2report --plan plan.json -f text
```

#### `--type, -t <type>`

Filter by resource type (repeatable).

**Type:** String
**Default:** None (no filtering)
**Repeatable:** Yes
**Example:**
```bash
tf2report --plan plan.json --type aws_instance
tf2report --plan plan.json -t aws_instance -t aws_s3_bucket
```

#### `--action, -a <action>`

Filter by action (repeatable).

**Type:** String
**Default:** None (no filtering)
**Valid Values:** `create`, `update`, `delete`, `replace`
**Repeatable:** Yes
**Example:**
```bash
tf2report --plan plan.json --action create
tf2report --plan plan.json -a delete -a replace
```

#### `--config <path>`

Path to configuration file.

**Type:** String
**Default:** `./tf2report.yaml`
**Example:**
```bash
tf2report --config /path/to/config.yaml
tf2report --config prod-config.yaml
```

#### `--verbose, -v`

Enable verbose output.

**Type:** Boolean
**Default:** `false`
**Example:**
```bash
tf2report --plan plan.json --verbose
tf2report --plan plan.json -v
```

#### `--help, -h`

Show help message.

**Type:** Boolean
**Example:**
```bash
tf2report --help
tf2report -h
```

#### `--version`

Show version information.

**Type:** Boolean
**Example:**
```bash
tf2report --version
```

## Flag Combinations

### Basic Usage

```bash
# Minimal
tf2report --plan plan.json

# With format
tf2report --plan plan.json --format json

# With config
tf2report --config config.yaml
```

### Filtering

```bash
# Single type
tf2report --plan plan.json --type aws_instance

# Multiple types
tf2report --plan plan.json \
  --type aws_instance \
  --type aws_s3_bucket

# Single action
tf2report --plan plan.json --action create

# Multiple actions
tf2report --plan plan.json \
  --action delete \
  --action replace

# Combined filters
tf2report --plan plan.json \
  --type aws_instance \
  --action create
```

### Output Control

```bash
# Different formats
tf2report --plan plan.json --format markdown
tf2report --plan plan.json --format text
tf2report --plan plan.json --format json

# Verbose mode
tf2report --plan plan.json --verbose

# Redirect output
tf2report --plan plan.json > report.md
tf2report --plan plan.json --format json > report.json
```

## Flag Precedence

When the same option is specified multiple times, the following precedence applies:

1. **Command-line flags** (highest)
2. **Environment variables**
3. **Configuration file**
4. **Default values** (lowest)

**Example:**
```yaml
# config.yaml
output_format: markdown
```

```bash
# Uses JSON (command-line overrides config)
tf2report --config config.yaml --format json
```

## Environment Variables

All flags can be set via environment variables:

| Flag | Environment Variable |
|------|---------------------|
| `--plan` | `TF2REPORT_TERRAFORM_PLAN_PATH` |
| `--format` | `TF2REPORT_OUTPUT_FORMAT` |
| `--verbose` | `TF2REPORT_VERBOSITY` |

**Example:**
```bash
export TF2REPORT_TERRAFORM_PLAN_PATH=plan.json
export TF2REPORT_OUTPUT_FORMAT=json
tf2report
```

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (invalid arguments, file not found, parsing error, etc.) |

## Examples

### Simple Report

```bash
tf2report --plan terraform.tfplan.json
```

### JSON Output

```bash
tf2report --plan terraform.tfplan.json --format json > report.json
```

### Filter AWS Instances

```bash
tf2report --plan terraform.tfplan.json --type aws_instance
```

### Destructive Changes

```bash
tf2report --plan terraform.tfplan.json --action delete --action replace
```

### Security Review

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_security_group \
  --format markdown > security.md
```

### With Configuration

```bash
tf2report --config prod-config.yaml --plan prod-plan.json
```

## Next Steps

- **[Configuration Schema](./configuration-schema.md)** - Configuration file reference
- **[Exit Codes](./exit-codes.md)** - Detailed exit code information
- **[Usage Guide](../usage.md)** - Comprehensive usage guide
