# tf2report Configuration Guide

## Overview

`tf2report` supports configuration via YAML files. Configuration files allow you to set default values for frequently used options.

## Configuration File Locations

`tf2report` looks for configuration files in the following locations (in order):

1. Path specified by `--config` flag
2. `./tf2report.yaml` (current directory)
3. `$HOME/.config/tf2report/tf2report.yaml`

## Configuration Schema

```yaml
# Path to Terraform plan JSON file
terraform_plan_path: terraform.tfplan.json

# Output format: markdown, text, or json
output_format: markdown

# Filters for resource selection
filters:
  # Filter by resource types (empty array means no filtering)
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_security_group
  
  # Filter by actions (empty array means no filtering)
  # Valid values: create, update, delete, replace
  actions:
    - create
    - delete

# Logging verbosity: debug, info, warn, error
verbosity: info
```

## Configuration Options

### `terraform_plan_path`

- **Type:** String
- **Default:** `terraform.tfplan.json`
- **Description:** Path to the Terraform plan JSON file to analyze

### `output_format`

- **Type:** String
- **Default:** `markdown`
- **Valid Values:** `markdown`, `text`, `json`
- **Description:** Format for the generated report

### `filters.resource_types`

- **Type:** Array of strings
- **Default:** `[]` (no filtering)
- **Description:** List of resource types to include in the report. When empty, all resource types are included.

### `filters.actions`

- **Type:** Array of strings
- **Default:** `[]` (no filtering)
- **Valid Values:** `create`, `update`, `delete`, `replace`
- **Description:** List of actions to include in the report. When empty, all actions are included.

### `verbosity`

- **Type:** String
- **Default:** `info`
- **Valid Values:** `debug`, `info`, `warn`, `error`
- **Description:** Logging verbosity level for diagnostic output

## Example Configurations

### Minimal Configuration

```yaml
terraform_plan_path: plan.json
```

### Production Configuration

```yaml
terraform_plan_path: production.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_rds_cluster
    - aws_s3_bucket
  actions:
    - create
    - delete
    - replace
verbosity: warn
```

### Development Configuration

```yaml
terraform_plan_path: dev.tfplan.json
output_format: text
verbosity: debug
```

### CI/CD Configuration

```yaml
terraform_plan_path: /tmp/terraform.tfplan.json
output_format: json
filters:
  actions:
    - delete
    - replace
verbosity: info
```

## Precedence

Configuration values are resolved in the following order (highest to lowest precedence):

1. Command-line flags
2. Environment variables (prefixed with `TF2REPORT_`)
3. Configuration file
4. Default values

## Using Configuration Files

### Specify Configuration File

```bash
tf2report --config /path/to/config.yaml
```

### Use Default Configuration

```bash
# Uses ./tf2report.yaml if it exists
tf2report --plan myplan.json
```

### Override Configuration with Flags

```bash
# Uses config file but overrides output format
tf2report --config config.yaml --format json
```

## Environment Variables

All configuration options can be set via environment variables using the `TF2REPORT_` prefix:

```bash
export TF2REPORT_TERRAFORM_PLAN_PATH=plan.json
export TF2REPORT_OUTPUT_FORMAT=markdown
export TF2REPORT_VERBOSITY=debug

tf2report
```
