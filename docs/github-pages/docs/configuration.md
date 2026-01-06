---
sidebar_position: 5
---

# Configuration

Complete guide to configuring tf2report using YAML configuration files.

## Overview

tf2report supports configuration via YAML files to set default values and avoid repetitive command-line options.

## Configuration File Locations

tf2report looks for configuration files in the following order:

1. Path specified by `--config` flag
2. `./tf2report.yaml` (current directory)
3. `$HOME/.config/tf2report/tf2report.yaml`

The first file found is used.

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

**Example:**
```yaml
terraform_plan_path: /path/to/plan.json
```

**Command-line override:**
```bash
tf2report --plan /path/to/other-plan.json
```

### `output_format`

- **Type:** String
- **Default:** `markdown`
- **Valid Values:** `markdown`, `text`, `json`
- **Description:** Format for the generated report

**Example:**
```yaml
output_format: json
```

**Command-line override:**
```bash
tf2report --format text
```

### `filters.resource_types`

- **Type:** Array of strings
- **Default:** `[]` (no filtering)
- **Description:** List of resource types to include in the report. When empty, all resource types are included.

**Example:**
```yaml
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_rds_cluster
    - aws_security_group
```

**Command-line override:**
```bash
tf2report --type aws_instance --type aws_s3_bucket
```

### `filters.actions`

- **Type:** Array of strings
- **Default:** `[]` (no filtering)
- **Valid Values:** `create`, `update`, `delete`, `replace`
- **Description:** List of actions to include in the report. When empty, all actions are included.

**Example:**
```yaml
filters:
  actions:
    - create
    - delete
    - replace
```

**Command-line override:**
```bash
tf2report --action create --action delete
```

### `verbosity`

- **Type:** String
- **Default:** `info`
- **Valid Values:** `debug`, `info`, `warn`, `error`
- **Description:** Logging verbosity level for diagnostic output

**Example:**
```yaml
verbosity: debug
```

**Command-line override:**
```bash
tf2report --verbose
```

## Example Configurations

### Minimal Configuration

Simplest possible configuration:

```yaml title="tf2report.yaml"
terraform_plan_path: plan.json
```

### Production Configuration

Configuration for production environment reviews:

```yaml title="tf2report.yaml"
terraform_plan_path: production.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_rds_cluster
    - aws_s3_bucket
    - aws_elasticache_cluster
  actions:
    - create
    - delete
    - replace
verbosity: warn
```

### Development Configuration

Configuration for development with verbose output:

```yaml title="tf2report.yaml"
terraform_plan_path: dev.tfplan.json
output_format: text
verbosity: debug
```

### Security Review Configuration

Focus on security-related resources:

```yaml title="security-review.yaml"
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_iam_role
    - aws_iam_policy
    - aws_iam_user
    - aws_security_group
    - aws_security_group_rule
    - aws_kms_key
    - aws_kms_alias
verbosity: info
```

### CI/CD Configuration

Configuration for CI/CD pipelines:

```yaml title="ci-config.yaml"
terraform_plan_path: /tmp/terraform.tfplan.json
output_format: json
filters:
  actions:
    - delete
    - replace
verbosity: info
```

### Database Changes Configuration

Focus on database-related changes:

```yaml title="database-review.yaml"
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_db_instance
    - aws_rds_cluster
    - aws_rds_cluster_instance
    - aws_dynamodb_table
    - aws_elasticache_cluster
    - aws_elasticache_replication_group
verbosity: info
```

## Using Configuration Files

### Default Configuration

Place a `tf2report.yaml` in your project directory:

```bash
# tf2report automatically finds ./tf2report.yaml
tf2report
```

### Custom Configuration

Specify a custom configuration file:

```bash
tf2report --config /path/to/custom-config.yaml
```

### Multiple Environments

Use different configurations for different environments:

```bash
# Production
tf2report --config prod-config.yaml

# Staging
tf2report --config staging-config.yaml

# Development
tf2report --config dev-config.yaml
```

### Per-Repository Configuration

Commit a `tf2report.yaml` to your repository for consistent team configuration:

```yaml title=".github/tf2report.yaml"
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_rds_cluster
verbosity: info
```

```bash
tf2report --config .github/tf2report.yaml --plan terraform.tfplan.json
```

## Configuration Precedence

Configuration values are resolved in the following order (highest to lowest):

1. **Command-line flags** (highest priority)
2. **Environment variables** (prefixed with `TF2REPORT_`)
3. **Configuration file**
4. **Default values** (lowest priority)

### Example

```yaml title="tf2report.yaml"
output_format: markdown
verbosity: info
```

```bash
# Uses JSON format (command-line override)
tf2report --format json

# Uses debug verbosity (environment variable)
export TF2REPORT_VERBOSITY=debug
tf2report

# Uses configuration file values
tf2report
```

## Environment Variables

All configuration options can be set via environment variables using the `TF2REPORT_` prefix:

```bash
# Set plan path
export TF2REPORT_TERRAFORM_PLAN_PATH=plan.json

# Set output format
export TF2REPORT_OUTPUT_FORMAT=json

# Set verbosity
export TF2REPORT_VERBOSITY=debug

# Run with environment variables
tf2report
```

### Environment Variable Mapping

| Configuration Option | Environment Variable |
|---------------------|---------------------|
| `terraform_plan_path` | `TF2REPORT_TERRAFORM_PLAN_PATH` |
| `output_format` | `TF2REPORT_OUTPUT_FORMAT` |
| `verbosity` | `TF2REPORT_VERBOSITY` |

## Validation

tf2report validates configuration files on load:

### Valid Configuration

```yaml
output_format: markdown
verbosity: info
```

✅ Loads successfully

### Invalid Configuration

```yaml
output_format: invalid_format
verbosity: info
```

❌ Error: Invalid output format. Must be markdown, text, or json.

## Tips and Best Practices

### 1. Commit Configuration Files

Commit configuration files to your repository for consistent team usage:

```bash
git add tf2report.yaml
git commit -m "Add tf2report configuration"
```

### 2. Environment-Specific Configurations

Create separate configurations for each environment:

```
configs/
├── prod.yaml
├── staging.yaml
└── dev.yaml
```

### 3. Use Filters Strategically

Configure filters based on your review process:

- **Security reviews** - IAM, security groups, KMS
- **Database reviews** - RDS, DynamoDB, ElastiCache
- **Compute reviews** - EC2, Lambda, ECS

### 4. CI/CD Integration

Use JSON format for CI/CD pipelines:

```yaml title=".github/tf2report.yaml"
output_format: json
filters:
  actions:
    - delete
    - replace
```

### 5. Override When Needed

Keep configuration general, override specific cases:

```bash
# Config uses markdown, override for JSON
tf2report --format json

# Config filters types, override to see all
tf2report --type ""
```

## Global Configuration

Create a global configuration in your home directory:

```bash
mkdir -p ~/.config/tf2report
```

```yaml title="~/.config/tf2report/tf2report.yaml"
output_format: markdown
verbosity: info
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
```

This applies to all tf2report invocations unless overridden.

## Next Steps

- **[Usage Guide](./usage.md)** - Learn how to use tf2report
- **[Output Formats](./output-formats.md)** - Understand different formats
- **[Filtering](./filtering.md)** - Advanced filtering techniques
- **[Examples](./examples/basic-usage.md)** - See practical examples
