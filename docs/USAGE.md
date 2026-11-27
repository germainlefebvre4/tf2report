# tf2report Usage Guide

## Overview

`tf2report` is a command-line tool that analyzes Terraform plan files and generates resource change reports in multiple formats.

## Basic Usage

```bash
tf2report --plan <path-to-plan.json>
```

## Command-Line Options

### Required Options

- `--plan, -p <path>` - Path to the Terraform plan JSON file

### Optional Options

- `--format, -f <format>` - Output format: `markdown` (default), `text`, `json`, or `svg`
- `--type, -t <type>` - Filter by resource type (can be specified multiple times)
- `--action, -a <action>` - Filter by action: `create`, `update`, `delete`, `replace`, or `no-op` (unchanged resources, drawn muted by the svg format)
- `--icons-dir <dir>` - Directory of SVG icons used by the `svg` format only
- `--config <path>` - Path to configuration file (default: `./tf2report.yaml`)
- `--verbose, -v` - Enable verbose output for debugging

### Help

- `--help, -h` - Display help information

## Common Use Cases

### Generate Markdown Report

```bash
tf2report --plan terraform.tfplan.json
```

### Generate Plain Text Report

```bash
tf2report --plan terraform.tfplan.json --format text
```

### Generate JSON Report

```bash
tf2report --plan terraform.tfplan.json --format json
```

### Generate an SVG Diagram

```bash
tf2report --plan terraform.tfplan.json --format svg > diagram.svg
tf2report --plan terraform.tfplan.json --format svg --icons-dir ./icons > diagram.svg
```

Status legend:

Each status is shown by color, symbol and border style, never by color alone:

| Status | Color | Symbol | Border |
|--------|-------|--------|--------|
| create | green | `+` | solid |
| update | orange | `~` | dashed |
| delete | red | `-` | dotted |
| replace | purple | `-/+` | dash-dot |
| no-op | grey, muted | none | thin solid |

What is drawn:

| `--type` | `--action` | create, update, delete, replace | no-op | Dependency context |
|----------|------------|---------------------------------|-------|--------------------|
| not given | not given | all | no | yes, one hop, muted |
| given | not given | matching types | no | no |
| any | given | matching, among those left after filtering | only if `no-op` is in the list | no |

#### Icon packs

Each resource gets an icon chosen along a chain: service family, provider family, category
(`compute`, `network`, `storage`, `database`, `security`), then `generic`. The binary embeds six
generic glyphs (the five categories and `generic`) and a selection of about 60 official AWS, Azure and
Google Cloud service icons (see `pkg/report/icons/vendor/NOTICE.md`). Other service and provider icons come
from a user pack given with `--icons-dir <dir>`, which wins over the embedded ones:

- One flat directory, read non-recursively; only `*.svg` files are used.
- File names are family ids: `aws-ec2.svg`, `azure-vm.svg`, `gcp-gce.svg`, `aws.svg`, `azure.svg`,
  `gcp.svg`, `compute.svg`, `network.svg`, `storage.svg`, `database.svg`, `security.svg`, `generic.svg`.
- For each resource the first existing file along its chain wins; at each level the pack wins over the embedded service icon, which wins over the built-in glyph.
- At most 256 KiB per file. Allowed elements: `svg g path rect circle ellipse line polyline polygon defs
  linearGradient radialGradient stop clipPath title desc`.
- Rejected content: scripts, `on*` attributes, `<style>`, `<image>`, `<use>`, `<foreignObject>`, external
  or non-fragment `href`, `url(` or `@import` in `style`, DOCTYPE and entities.
- A rejected file is skipped (listed with `--verbose`) and the next candidate is used.
- A missing or unreadable directory is an error before any output. With a format other than `svg`, `--icons-dir` is ignored.

Limits:

- No isometric rendering and no PNG export: the output is one 2D SVG document.
- A reference to a module output points to the module container, not to the producing resource.
- Deleted resources have no arrows, because their relations are no longer in the configuration.
- Data sources are never drawn. Arrows towards resources that are not drawn are omitted.
- The diagram carries addresses, types, names, providers and statuses only, never attribute values.

### Filter by Resource Type

```bash
tf2report --plan terraform.tfplan.json --type aws_instance --type aws_s3_bucket
```

### Filter by Action

```bash
tf2report --plan terraform.tfplan.json --action create --action delete
```

### Combine Filters

```bash
tf2report --plan terraform.tfplan.json --type aws_instance --action create
```

### Use Custom Configuration File

```bash
tf2report --config my-config.yaml
```

### Enable Verbose Output

```bash
tf2report --plan terraform.tfplan.json --verbose
```

## Generating Terraform Plan JSON

Before using `tf2report`, you need to generate a Terraform plan in JSON format:

```bash
# Create a plan
terraform plan -out=tfplan

# Convert plan to JSON
terraform show -json tfplan > terraform.tfplan.json

# Analyze with tf2report
tf2report --plan terraform.tfplan.json
```

## Exit Codes

- `0` - Success
- `1` - Error occurred

## Environment Variables

Environment variables can be used to configure `tf2report`. All variables use the `TF2REPORT_` prefix:

- `TF2REPORT_TERRAFORM_PLAN_PATH` - Default plan file path
- `TF2REPORT_OUTPUT_FORMAT` - Default output format
- `TF2REPORT_VERBOSITY` - Logging verbosity level

Command-line flags override environment variables and configuration file settings.
