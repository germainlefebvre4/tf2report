---
sidebar_position: 2
---

# Configuration Schema

Complete reference for tf2report configuration files.

## File Format

Configuration files use YAML format.

## File Locations

tf2report searches for configuration in order:

1. Path specified by `--config` flag
2. `./tf2report.yaml` (current directory)
3. `$HOME/.config/tf2report/tf2report.yaml`

## Schema

```yaml
# Path to Terraform plan JSON file
terraform_plan_path: string

# Output format
output_format: string

# Resource filters
filters:
  resource_types: []string
  actions: []string

# Logging verbosity
verbosity: string
```

## Fields

### `terraform_plan_path`

Path to Terraform plan JSON file.

**Type:** String  
**Required:** No  
**Default:** `terraform.tfplan.json`  
**Example:**
```yaml
terraform_plan_path: /path/to/plan.json
```

### `output_format`

Output format for reports.

**Type:** String  
**Required:** No  
**Default:** `markdown`  
**Valid Values:** `markdown`, `text`, `json`  
**Example:**
```yaml
output_format: json
```

### `filters`

Resource filtering options.

**Type:** Object  
**Required:** No  
**Default:** Empty (no filtering)

#### `filters.resource_types`

List of resource types to include.

**Type:** Array of strings  
**Required:** No  
**Default:** `[]` (all types included)  
**Example:**
```yaml
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_rds_cluster
```

#### `filters.actions`

List of actions to include.

**Type:** Array of strings  
**Required:** No  
**Default:** `[]` (all actions included)  
**Valid Values:** `create`, `update`, `delete`, `replace`  
**Example:**
```yaml
filters:
  actions:
    - create
    - delete
```

### `verbosity`

Logging verbosity level.

**Type:** String  
**Required:** No  
**Default:** `info`  
**Valid Values:** `debug`, `info`, `warn`, `error`  
**Example:**
```yaml
verbosity: debug
```

## Complete Example

```yaml title="tf2report.yaml"
# Terraform plan file location
terraform_plan_path: terraform.tfplan.json

# Output format
output_format: markdown

# Filters
filters:
  # Resource types to include
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_rds_cluster
    - aws_security_group
  
  # Actions to include
  actions:
    - create
    - delete
    - replace

# Logging verbosity
verbosity: info
```

## Environment-Specific Examples

### Production

```yaml title="prod-config.yaml"
terraform_plan_path: production.tfplan.json
output_format: markdown
filters:
  actions:
    - delete
    - replace
verbosity: warn
```

### Development

```yaml title="dev-config.yaml"
terraform_plan_path: dev.tfplan.json
output_format: text
verbosity: debug
```

### CI/CD

```yaml title="ci-config.yaml"
terraform_plan_path: /tmp/terraform.tfplan.json
output_format: json
filters:
  actions:
    - delete
    - replace
verbosity: info
```

## Validation

tf2report validates configuration on load:

### Valid Configuration

```yaml
output_format: markdown
verbosity: info
```
✅ Loads successfully

### Invalid Format

```yaml
output_format: invalid_format
```
❌ Error: Invalid output format

### Invalid Verbosity

```yaml
verbosity: invalid_level
```
❌ Error: Invalid verbosity level

### Invalid Action

```yaml
filters:
  actions:
    - invalid_action
```
❌ Error: Invalid action type

## Merging Configuration

Configuration is merged from multiple sources:

1. Load defaults
2. Apply configuration file
3. Apply environment variables
4. Apply command-line flags

**Example:**

```yaml title="config.yaml"
output_format: markdown
verbosity: info
```

```bash
# Environment variable
export TF2REPORT_OUTPUT_FORMAT=text

# Command-line flag (highest priority)
tf2report --config config.yaml --format json

# Result: Uses JSON format
```

## Best Practices

### 1. Commit Configuration

Commit configuration files to version control:

```yaml title=".github/tf2report.yaml"
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
```

### 2. Environment-Specific Configs

Use separate configs for each environment:

```
configs/
├── prod.yaml
├── staging.yaml
└── dev.yaml
```

### 3. Document Custom Configs

Add comments to explain filters:

```yaml
filters:
  # Security review: IAM and network resources only
  resource_types:
    - aws_iam_role
    - aws_iam_policy
    - aws_security_group
    - aws_kms_key
```

### 4. Use Minimal Config

Only specify non-default values:

```yaml
# Good: Only overrides
output_format: json

# Unnecessary: Same as defaults
terraform_plan_path: terraform.tfplan.json
verbosity: info
```

## Next Steps

- **[CLI Options](./cli-options.md)** - Command-line reference
- **[Configuration Guide](../configuration.md)** - Usage guide
- **[Exit Codes](./exit-codes.md)** - Exit code reference
