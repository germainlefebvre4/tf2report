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

- `--format, -f <format>` - Output format: `markdown` (default), `text`, or `json`
- `--type, -t <type>` - Filter by resource type (can be specified multiple times)
- `--action, -a <action>` - Filter by action: `create`, `update`, `delete`, or `replace`
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
