# tf2report Output Formats

## Overview

`tf2report` supports three output formats for generating Terraform plan change reports:

- **Markdown** - Human-readable format with tables and sections
- **Text** - Plain text format for terminal display
- **JSON** - Structured format for programmatic consumption

## Markdown Format

The default output format designed for human readability and documentation.

### Characteristics

- Headers and sections for clear organization
- Markdown tables for summary statistics
- Resource lists with code formatting
- Suitable for documentation, pull requests, and reports

### Structure

1. **Header** - Title and Terraform version
2. **Summary** - Total changes and counts by action type
3. **Changes by Resource Type** - Detailed breakdown by resource type
   - Resources grouped by action (Add, Change, Replace, Destroy)
   - Alphabetically sorted resource addresses

### Selection

```bash
tf2report --plan plan.json --format markdown
```

## Text Format

Plain text format optimized for terminal display and simple text processing.

### Characteristics

- Simple text formatting with separators
- Indented hierarchical structure
- No special characters or markup
- Suitable for terminal output and log files

### Structure

1. **Header** - Title with separator and Terraform version
2. **Summary** - Change counts with aligned formatting
3. **Changes by Resource Type** - Grouped and indented resource lists

### Selection

```bash
tf2report --plan plan.json --format text
```

## JSON Format

Structured JSON output for programmatic analysis and integration with other tools.

### Characteristics

- Machine-readable structured data
- Complete change information
- Suitable for CI/CD pipelines and automation
- Can be parsed by other tools

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
    }
  ],
  "by_type": {
    "aws_s3_bucket": {
      "type": "aws_s3_bucket",
      "to_add": 1,
      "to_change": 0,
      "to_destroy": 0,
      "to_replace": 0,
      "resources": ["aws_s3_bucket.data"]
    }
  }
}
```

### Fields

- `terraform_version` - Version of Terraform that generated the plan
- `summary` - Aggregate statistics
  - `total_changes` - Total number of resource changes
  - `to_add` - Resources to be created
  - `to_change` - Resources to be updated
  - `to_destroy` - Resources to be deleted
  - `to_replace` - Resources to be replaced
- `changes` - Array of individual resource changes
  - `address` - Terraform resource address
  - `type` - Resource type
  - `name` - Resource name
  - `action` - Action type (create, update, delete, replace)
- `by_type` - Statistics grouped by resource type
  - Per-type action counts
  - List of affected resources

### Selection

```bash
tf2report --plan plan.json --format json
```

## Format Selection

### Via Command Line

```bash
# Markdown (default)
tf2report --plan plan.json

# Text
tf2report --plan plan.json --format text

# JSON
tf2report --plan plan.json --format json
```

### Via Configuration File

```yaml
output_format: json
```

### Via Environment Variable

```bash
export TF2REPORT_OUTPUT_FORMAT=text
tf2report --plan plan.json
```

## Use Cases by Format

### Markdown

- Documentation generation
- Pull request comments
- Wiki pages
- Human review

### Text

- Terminal display
- Simple log files
- Email notifications
- Quick reviews

### JSON

- CI/CD pipeline integration
- Automated analysis
- Metrics collection
- Integration with other tools
- Custom reporting scripts

## Output Redirection

All formats can be redirected to files:

```bash
# Save Markdown to file
tf2report --plan plan.json > report.md

# Save Text to file
tf2report --plan plan.json --format text > report.txt

# Save JSON to file
tf2report --plan plan.json --format json > report.json
```
