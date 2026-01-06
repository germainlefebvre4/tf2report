---
sidebar_position: 6
---

# Output Formats

Detailed guide to tf2report's output formats and when to use each one.

## Overview

tf2report supports three output formats, each optimized for different use cases:

- **Markdown** - Human-readable format with tables and sections
- **Text** - Plain text format for terminal display
- **JSON** - Structured format for programmatic consumption

## Markdown Format

The default output format designed for maximum readability and documentation.

### Characteristics

- Headers and sections for clear organization
- Markdown tables for summary statistics
- Resource lists with code formatting (backticks)
- Suitable for documentation, pull requests, and reports
- Compatible with GitHub, GitLab, and other Markdown renderers

### Structure

1. **Header** - Title and Terraform version
2. **Summary Section**
   - Total changes count
   - Action breakdown table
3. **Changes by Resource Type**
   - Resources grouped by type
   - Sub-grouped by action (Add, Change, Replace, Destroy)
   - Alphabetically sorted resource addresses

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

### aws_s3_bucket (2)

**To Add (1):**

- `aws_s3_bucket.data`

**To Destroy (1):**

- `aws_s3_bucket.old_data`
```

### Use Cases

- **Pull Request Comments** - Clear, formatted change summaries
- **Documentation** - Infrastructure change logs
- **Wiki Pages** - Team documentation
- **Emails** - Stakeholder updates
- **Reports** - Compliance and audit trails

### Selection

```bash
# Explicit
tf2report --plan plan.json --format markdown

# Default (no format flag needed)
tf2report --plan plan.json
```

## Text Format

Plain text format optimized for terminal display and simple text processing.

### Characteristics

- Simple text formatting with separators
- Indented hierarchical structure
- No special characters or markup
- ANSI-safe (works in all terminals)
- Easy to parse with basic text tools

### Structure

1. **Header** - Title with separator line and Terraform version
2. **Summary** - Change counts with aligned formatting
3. **Changes by Resource Type** - Grouped and indented resource lists

### Example Output

```
Terraform Plan Summary
======================
Terraform Version: 1.10.5

Summary
-------
Total Changes: 4
  Add: 1
  Change: 1
  Destroy: 1
  Replace: 1

Changes by Resource Type
------------------------

aws_instance (2)
  To Change (1):
    - aws_instance.web[0]
  To Replace (1):
    - aws_instance.web[1]

aws_s3_bucket (2)
  To Add (1):
    - aws_s3_bucket.data
  To Destroy (1):
    - aws_s3_bucket.old_data
```

### Use Cases

- **Terminal Output** - Direct viewing in console
- **Log Files** - Simple text logs
- **Email Notifications** - Plain text emails
- **Quick Reviews** - Fast command-line analysis
- **CI/CD Logs** - Pipeline output logs

### Selection

```bash
tf2report --plan plan.json --format text
```

## JSON Format

Structured JSON output for programmatic analysis and tool integration.

### Characteristics

- Machine-readable structured data
- Complete change information
- Easy to parse with standard tools (jq, Python, etc.)
- Suitable for automation and integration
- Includes all available metadata

### Structure

```json
{
  "terraform_version": "1.10.5",
  "summary": {
    "total_changes": 4,
    "to_add": 1,
    "to_change": 1,
    "to_destroy": 1,
    "to_replace": 1
  },
  "changes": [
    {
      "address": "aws_s3_bucket.data",
      "type": "aws_s3_bucket",
      "name": "data",
      "action": "create"
    },
    {
      "address": "aws_instance.web[0]",
      "type": "aws_instance",
      "name": "web",
      "action": "update"
    },
    {
      "address": "aws_instance.web[1]",
      "type": "aws_instance",
      "name": "web",
      "action": "replace"
    },
    {
      "address": "aws_s3_bucket.old_data",
      "type": "aws_s3_bucket",
      "name": "old_data",
      "action": "delete"
    }
  ],
  "by_type": {
    "aws_instance": {
      "type": "aws_instance",
      "to_add": 0,
      "to_change": 1,
      "to_destroy": 0,
      "to_replace": 1,
      "resources": [
        "aws_instance.web[0]",
        "aws_instance.web[1]"
      ]
    },
    "aws_s3_bucket": {
      "type": "aws_s3_bucket",
      "to_add": 1,
      "to_change": 0,
      "to_destroy": 1,
      "to_replace": 0,
      "resources": [
        "aws_s3_bucket.data",
        "aws_s3_bucket.old_data"
      ]
    }
  }
}
```

### Fields

#### Root Level

- `terraform_version` (string) - Version of Terraform that generated the plan
- `summary` (object) - Aggregate statistics
- `changes` (array) - List of individual resource changes
- `by_type` (object) - Statistics grouped by resource type

#### Summary Object

- `total_changes` (number) - Total number of resource changes
- `to_add` (number) - Resources to be created
- `to_change` (number) - Resources to be updated
- `to_destroy` (number) - Resources to be deleted
- `to_replace` (number) - Resources to be replaced

#### Change Object

- `address` (string) - Terraform resource address
- `type` (string) - Resource type (e.g., `aws_instance`)
- `name` (string) - Resource name
- `action` (string) - Action type: `create`, `update`, `delete`, or `replace`

#### By Type Object

For each resource type:
- `type` (string) - Resource type
- `to_add` (number) - Count of resources to add
- `to_change` (number) - Count of resources to change
- `to_destroy` (number) - Count of resources to destroy
- `to_replace` (number) - Count of resources to replace
- `resources` (array) - List of resource addresses

### Use Cases

- **CI/CD Pipelines** - Automated analysis and decisions
- **Metrics Collection** - Track infrastructure changes over time
- **Custom Reports** - Build custom reporting tools
- **Integration** - Feed into other tools and systems
- **Analysis Scripts** - Programmatic change analysis
- **Dashboards** - Real-time infrastructure change monitoring

### Selection

```bash
tf2report --plan plan.json --format json
```

## Format Selection

### Via Command Line

```bash
# Markdown (default)
tf2report --plan plan.json

# Markdown (explicit)
tf2report --plan plan.json --format markdown

# Text
tf2report --plan plan.json --format text

# JSON
tf2report --plan plan.json --format json
```

### Via Configuration File

```yaml title="tf2report.yaml"
output_format: json
```

```bash
tf2report --plan plan.json
```

### Via Environment Variable

```bash
export TF2REPORT_OUTPUT_FORMAT=text
tf2report --plan plan.json
```

## Output Redirection

All formats support output redirection to files:

### Save to File

```bash
# Markdown
tf2report --plan plan.json > report.md

# Text
tf2report --plan plan.json --format text > report.txt

# JSON
tf2report --plan plan.json --format json > report.json
```

### Pipe to Other Commands

```bash
# View with pager
tf2report --plan plan.json | less

# Filter with grep
tf2report --plan plan.json --format text | grep aws_s3

# Process JSON with jq
tf2report --plan plan.json --format json | jq '.summary'

# Count lines
tf2report --plan plan.json | wc -l
```

## Format Comparison

| Feature | Markdown | Text | JSON |
|---------|----------|------|------|
| Human Readable | ✅ Excellent | ✅ Good | ❌ No |
| Terminal Display | ⚠️ OK | ✅ Excellent | ❌ No |
| Documentation | ✅ Excellent | ⚠️ OK | ❌ No |
| Programmatic | ❌ No | ⚠️ Limited | ✅ Excellent |
| CI/CD Integration | ⚠️ OK | ⚠️ OK | ✅ Excellent |
| File Size | Medium | Small | Large |
| Parsing Complexity | High | Medium | Low |

## Advanced Usage

### Combine Formats

Generate multiple formats:

```bash
# Generate both Markdown and JSON
tf2report --plan plan.json > report.md
tf2report --plan plan.json --format json > report.json
```

### Conditional Formatting

Choose format based on environment:

```bash
#!/bin/bash
if [ "$CI" = "true" ]; then
  FORMAT="json"
else
  FORMAT="markdown"
fi

tf2report --plan plan.json --format $FORMAT
```

### Process JSON Output

```bash
# Get summary statistics
tf2report --plan plan.json --format json | jq '.summary'

# List all deletions
tf2report --plan plan.json --format json | \
  jq '.changes[] | select(.action == "delete") | .address'

# Count changes by type
tf2report --plan plan.json --format json | \
  jq '.by_type | to_entries[] | "\(.key): \(.value.to_add + .value.to_change + .value.to_destroy + .value.to_replace)"'

# Check for destructive changes
tf2report --plan plan.json --format json | \
  jq 'if (.summary.to_destroy + .summary.to_replace) > 0 then "WARNING: Destructive changes!" else "OK" end'
```

## Best Practices

### 1. Choose the Right Format

- **Markdown** for humans (reviews, documentation)
- **Text** for terminals (quick views, logs)
- **JSON** for machines (automation, integration)

### 2. Use JSON for Automation

Always use JSON in CI/CD pipelines and automation scripts for reliable parsing.

### 3. Markdown for Communication

Use Markdown when communicating changes to stakeholders, teams, or in pull requests.

### 4. Text for Quick Checks

Use text format for quick terminal reviews and debugging.

### 5. Save JSON for Analysis

Save JSON output for historical analysis and metrics:

```bash
tf2report --plan plan.json --format json > reports/$(date +%Y%m%d-%H%M%S).json
```

## Next Steps

- **[Usage Guide](./usage.md)** - Learn how to use tf2report
- **[Filtering](./filtering.md)** - Filter resources and actions
- **[Examples](./examples/basic-usage.md)** - See practical examples
- **[CI/CD Integration](./examples/ci-cd-integration.md)** - Automate reports
