# tf2report Output Formats

## Overview

`tf2report` supports four output formats for generating Terraform plan change reports:

- **Markdown** - Human-readable format with tables and sections
- **Text** - Plain text format for terminal display
- **JSON** - Structured format for programmatic consumption
- **SVG** - 2D architecture diagram of the changed resources

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

## SVG Format

### Characteristics

- One SVG document on standard output, no external resource, script or font
- Provider and module containers, one node per resource, dependency arrows, a legend and icons
- Deterministic: the same plan and options give identical bytes
- Only addresses, types, names, providers and statuses; never attribute values

### Status cues

Each status is shown by color, symbol and border style, never by color alone:

| Status | Color | Symbol | Border |
|--------|-------|--------|--------|
| create | green | `+` | solid |
| update | orange | `~` | dashed |
| delete | red | `-` | dotted |
| replace | purple | `-/+` | dash-dot |
| no-op | grey, muted | none | thin solid |

### Selection

```bash
tf2report --plan plan.json --format svg > diagram.svg
```

### What is drawn

| `--type` | `--action` | create, update, delete, replace | no-op | Dependency context |
|----------|------------|---------------------------------|-------|--------------------|
| not given | not given | all | no | yes, one hop, muted |
| given | not given | matching types | no | no |
| any | given | matching, among those left after filtering | only if `no-op` is in the list | no |

### Icons

Each resource gets an icon chosen along a chain: service family, provider family, category
(`compute`, `network`, `storage`, `database`, `security`), then `generic`. The binary embeds six
generic glyphs (the five categories and `generic`). Service and provider icons come only from a
user pack given with `--icons-dir <dir>`:

- One flat directory, read non-recursively; only `*.svg` files are used.
- File names are family ids: `aws-ec2.svg`, `azure-vm.svg`, `gcp-gce.svg`, `aws.svg`, `azure.svg`,
  `gcp.svg`, `compute.svg`, `network.svg`, `storage.svg`, `database.svg`, `security.svg`, `generic.svg`.
- For each resource the first existing file along its chain wins; at each level the pack wins over the built-in glyph.
- At most 256 KiB per file. Allowed elements: `svg g path rect circle ellipse line polyline polygon defs
  linearGradient radialGradient stop clipPath title desc`.
- Rejected content: scripts, `on*` attributes, `<style>`, `<image>`, `<use>`, `<foreignObject>`, external
  or non-fragment `href`, `url(` or `@import` in `style`, DOCTYPE and entities.
- A rejected file is skipped (listed with `--verbose`) and the next candidate is used.
- A missing or unreadable directory is an error before any output. With a format other than `svg`, `--icons-dir` is ignored.

### Limits

- No isometric rendering and no PNG export: the output is one 2D SVG document.
- A reference to a module output points to the module container, not to the producing resource.
- Deleted resources have no arrows, because their relations are no longer in the configuration.
- Data sources are never drawn. Arrows towards resources that are not drawn are omitted.
- The diagram carries addresses, types, names, providers and statuses only, never attribute values.

## Format Selection

### Via Command Line

```bash
# Markdown (default)
tf2report --plan plan.json

# Text
tf2report --plan plan.json --format text

# JSON
tf2report --plan plan.json --format json

# SVG
tf2report --plan plan.json --format svg > diagram.svg
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
