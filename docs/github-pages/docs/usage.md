---
sidebar_position: 4
---

# Usage Guide

Comprehensive guide to using tf2report for analyzing Terraform plans.

## Basic Usage

The simplest way to use tf2report:

```bash
tf2report --plan <path-to-plan.json>
```

## Command-Line Options

### Required Options

#### `--plan, -p <path>`

Path to the Terraform plan JSON file.

```bash
tf2report --plan terraform.tfplan.json
```

### Optional Options

#### `--format, -f <format>`

Output format: `markdown` (default), `text`, or `json`.

```bash
tf2report --plan terraform.tfplan.json --format markdown
tf2report --plan terraform.tfplan.json --format text
tf2report --plan terraform.tfplan.json --format json
```

#### `--type, -t <type>`

Filter by resource type. Can be specified multiple times.

```bash
tf2report --plan terraform.tfplan.json --type aws_instance
tf2report --plan terraform.tfplan.json --type aws_instance --type aws_s3_bucket
```

#### `--action, -a <action>`

Filter by action: `create`, `update`, `delete`, or `replace`. Can be specified multiple times.

```bash
tf2report --plan terraform.tfplan.json --action create
tf2report --plan terraform.tfplan.json --action delete --action replace
```

#### `--config <path>`

Path to configuration file (default: `./tf2report.yaml`).

```bash
tf2report --config /path/to/config.yaml
```

#### `--verbose, -v`

Enable verbose output for debugging.

```bash
tf2report --plan terraform.tfplan.json --verbose
```

#### `--help, -h`

Display help information.

```bash
tf2report --help
```

## Generating Terraform Plan JSON

Before using tf2report, generate a Terraform plan in JSON format:

### Step 1: Create Plan File

```bash
terraform plan -out=tfplan
```

This creates a binary plan file named `tfplan`.

### Step 2: Convert to JSON

```bash
terraform show -json tfplan > terraform.tfplan.json
```

This converts the binary plan to JSON format that tf2report can analyze.

### Combined Command

```bash
terraform plan -out=tfplan && terraform show -json tfplan > terraform.tfplan.json
```

## Output Formats

### Markdown (Default)

Human-readable format with tables and sections:

```bash
tf2report --plan terraform.tfplan.json
```

**Output:**
```markdown
# Terraform Plan Summary

**Terraform Version:** 1.10.5

## Summary

**Total Changes:** 3

| Action | Count |
|--------|-------|
| Add | 2 |
| Change | 1 |
```

Perfect for:
- Documentation
- Pull request comments
- Wiki pages
- Reports

### Plain Text

Simple text format for terminal display:

```bash
tf2report --plan terraform.tfplan.json --format text
```

**Output:**
```
Terraform Plan Summary
======================
Terraform Version: 1.10.5

Summary
-------
Total Changes: 3
  Add: 2
  Change: 1
```

Perfect for:
- Terminal output
- Log files
- Email notifications

### JSON

Structured data for programmatic consumption:

```bash
tf2report --plan terraform.tfplan.json --format json
```

**Output:**
```json
{
  "terraform_version": "1.10.5",
  "summary": {
    "total_changes": 3,
    "to_add": 2,
    "to_change": 1,
    "to_destroy": 0,
    "to_replace": 0
  },
  "changes": [...]
}
```

Perfect for:
- CI/CD pipelines
- Automation scripts
- Metrics collection
- Custom analysis

## Filtering

### By Resource Type

Filter to show only specific resource types:

```bash
# Single type
tf2report --plan terraform.tfplan.json --type aws_s3_bucket

# Multiple types
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_s3_bucket \
  --type aws_security_group
```

### By Action

Filter to show only specific actions:

```bash
# Show only new resources
tf2report --plan terraform.tfplan.json --action create

# Show only destructive changes
tf2report --plan terraform.tfplan.json --action delete --action replace

# Show only modifications
tf2report --plan terraform.tfplan.json --action update
```

### Combined Filters

Combine type and action filters:

```bash
# Show only new S3 buckets
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --action create

# Show only changes or deletions to instances and databases
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_db_instance \
  --action update \
  --action delete
```

## Configuration File

For repeated use, create a configuration file to avoid typing the same options:

### Create Configuration

```yaml title="tf2report.yaml"
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

### Use Configuration

```bash
# Uses ./tf2report.yaml by default
tf2report

# Or specify a custom config
tf2report --config /path/to/config.yaml
```

### Override Configuration

Command-line options override configuration file settings:

```bash
# Config says markdown, but output JSON
tf2report --format json
```

See [Configuration Guide](./configuration.md) for detailed configuration options.

## Output Redirection

Save output to files:

```bash
# Save Markdown
tf2report --plan terraform.tfplan.json > report.md

# Save text
tf2report --plan terraform.tfplan.json --format text > report.txt

# Save JSON
tf2report --plan terraform.tfplan.json --format json > report.json
```

## Common Workflows

### Code Review

Generate a report for pull request review:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > terraform.tfplan.json
tf2report --plan terraform.tfplan.json > pr-comment.md
```

Post `pr-comment.md` as a pull request comment.

### Security Review

Filter security-related resources:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_security_group \
  --type aws_kms_key \
  --format markdown > security-review.md
```

### Destructive Changes Check

Check for any destructive changes:

```bash
tf2report --plan terraform.tfplan.json \
  --action delete \
  --action replace \
  --format json > destructive-changes.json

# Check if any destructive changes exist
if [ $(jq '.summary.to_destroy + .summary.to_replace' destructive-changes.json) -gt 0 ]; then
  echo "Warning: Destructive changes detected!"
  exit 1
fi
```

### Multi-Environment

Analyze plans for different environments:

```bash
# Production
terraform plan -out=prod.tfplan
terraform show -json prod.tfplan > prod.tfplan.json
tf2report --plan prod.tfplan.json --config prod-config.yaml > prod-report.md

# Staging
terraform plan -out=staging.tfplan
terraform show -json staging.tfplan > staging.tfplan.json
tf2report --plan staging.tfplan.json --config staging-config.yaml > staging-report.md
```

## Exit Codes

- **0** - Success
- **1** - Error occurred (invalid arguments, file not found, parsing error, etc.)

### Example: Check for Errors

```bash
if tf2report --plan terraform.tfplan.json > report.md; then
  echo "Report generated successfully"
else
  echo "Error generating report"
  exit 1
fi
```

## Environment Variables

tf2report respects standard environment variables:

### `TF2REPORT_CONFIG`

Default configuration file path:

```bash
export TF2REPORT_CONFIG=/path/to/config.yaml
tf2report --plan terraform.tfplan.json
```

## Verbose Mode

Enable verbose logging for troubleshooting:

```bash
tf2report --plan terraform.tfplan.json --verbose
```

This shows:
- Configuration loading
- File parsing progress
- Filter application
- Report generation steps

## Tips and Best Practices

### 1. Use Configuration Files

For consistent results across your team, commit a `tf2report.yaml` configuration file to your repository.

### 2. Automate in CI/CD

Integrate tf2report into your CI/CD pipeline to automatically generate reports for every plan.

### 3. Filter Wisely

Use filters to focus on relevant changes:
- Security reviews: IAM, security groups, KMS
- Database reviews: RDS, DynamoDB
- Compute reviews: EC2, Lambda, ECS

### 4. Save JSON for Analysis

Use JSON output for programmatic analysis and metrics collection.

### 5. Combine with jq

Process JSON output with jq for advanced filtering:

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes[] | select(.action == "delete")'
```

## Next Steps

- **[Configuration](./configuration.md)** - Learn about configuration options
- **[Output Formats](./output-formats.md)** - Understand output formats in detail
- **[Filtering](./filtering.md)** - Advanced filtering techniques
- **[Examples](./examples/basic-usage.md)** - See practical examples
- **[CI/CD Integration](./examples/ci-cd-integration.md)** - Integrate with pipelines
