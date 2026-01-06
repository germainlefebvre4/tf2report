---
sidebar_position: 3
---

# Exit Codes

Reference for tf2report exit codes and error handling.

## Exit Code Summary

| Code | Meaning | Description |
|------|---------|-------------|
| `0` | Success | Command completed successfully |
| `1` | Error | An error occurred during execution |

## Exit Code 0: Success

Returned when tf2report completes successfully.

**Conditions:**
- Plan file parsed successfully
- Report generated successfully
- No errors encountered

**Example:**
```bash
tf2report --plan terraform.tfplan.json
echo $?
# Output: 0
```

## Exit Code 1: Error

Returned when an error occurs during execution.

### Common Error Scenarios

#### File Not Found

```bash
tf2report --plan nonexistent.json
# Exit code: 1
# Error: file not found
```

#### Invalid JSON

```bash
tf2report --plan invalid.json
# Exit code: 1
# Error: failed to parse JSON
```

#### Invalid Arguments

```bash
tf2report --format invalid_format
# Exit code: 1
# Error: invalid format
```

#### Missing Required Flag

```bash
tf2report
# Exit code: 1
# Error: required flag --plan not provided
```

#### Configuration Error

```bash
tf2report --config invalid-config.yaml
# Exit code: 1
# Error: failed to load configuration
```

## Error Messages

### File Errors

**File not found:**
```
Error: failed to open plan file: open terraform.tfplan.json: no such file or directory
```

**Permission denied:**
```
Error: failed to open plan file: open terraform.tfplan.json: permission denied
```

### Parsing Errors

**Invalid JSON:**
```
Error: failed to parse plan file: invalid character '}' after object key
```

**Unexpected format:**
```
Error: plan file does not match expected Terraform plan format
```

### Configuration Errors

**Invalid format:**
```
Error: invalid output format: must be markdown, text, or json
```

**Invalid action:**
```
Error: invalid action: must be create, update, delete, or replace
```

## Using Exit Codes

### Check for Success

```bash
if tf2report --plan plan.json > report.md; then
  echo "✅ Report generated successfully"
else
  echo "❌ Failed to generate report"
  exit 1
fi
```

### Script Integration

```bash
#!/bin/bash
set -e  # Exit on any error

tf2report --plan terraform.tfplan.json > report.md
echo "Report generated: report.md"
```

### Conditional Logic

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json --format json > report.json

if [ $? -eq 0 ]; then
  # Success - check for destructive changes
  DESTRUCTIVE=$(jq '.summary.to_destroy + .summary.to_replace' report.json)

  if [ "$DESTRUCTIVE" -gt 0 ]; then
    echo "⚠️  WARNING: Destructive changes detected"
    exit 1
  fi
else
  # tf2report failed
  echo "❌ Failed to generate report"
  exit 1
fi
```

### CI/CD Integration

```yaml
# GitHub Actions
- name: Generate Report
  run: tf2report --plan terraform.tfplan.json > report.md
  continue-on-error: false  # Fail workflow on error

- name: Check Exit Code
  run: |
    if ! tf2report --plan terraform.tfplan.json > report.md; then
      echo "Report generation failed"
      exit 1
    fi
```

## Error Handling Best Practices

### 1. Always Check Exit Codes

```bash
# Good
if tf2report --plan plan.json > report.md; then
  echo "Success"
fi

# Bad (ignores errors)
tf2report --plan plan.json > report.md
```

### 2. Provide Helpful Error Context

```bash
#!/bin/bash
if ! tf2report --plan terraform.tfplan.json > report.md; then
  echo "❌ Error: Failed to generate Terraform plan report"
  echo "Check that terraform.tfplan.json exists and is valid"
  exit 1
fi
```

### 3. Use set -e in Scripts

```bash
#!/bin/bash
set -e  # Exit immediately on error

terraform plan -out=tfplan
terraform show -json tfplan > terraform.tfplan.json
tf2report --plan terraform.tfplan.json > report.md
```

### 4. Clean Up on Failure

```bash
#!/bin/bash
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

tf2report --plan plan.json > "$TEMP_DIR/report.md" || {
  echo "Failed to generate report"
  exit 1
}

cp "$TEMP_DIR/report.md" ./final-report.md
```

## Debugging Errors

### Enable Verbose Mode

```bash
tf2report --plan plan.json --verbose
```

Shows detailed logging to help diagnose issues.

### Check File Existence

```bash
if [ ! -f "terraform.tfplan.json" ]; then
  echo "Error: Plan file not found"
  exit 1
fi

tf2report --plan terraform.tfplan.json
```

### Validate JSON

```bash
# Validate JSON before processing
if jq empty terraform.tfplan.json 2>/dev/null; then
  tf2report --plan terraform.tfplan.json
else
  echo "Error: Invalid JSON in plan file"
  exit 1
fi
```

## Exit Code in Different Shells

### Bash/Zsh

```bash
tf2report --plan plan.json
echo $?
```

### PowerShell

```powershell
tf2report --plan plan.json
echo $LASTEXITCODE
```

### Fish

```fish
tf2report --plan plan.json
echo $status
```

## Next Steps

- **[CLI Options](./cli-options.md)** - Command-line reference
- **[Configuration Schema](./configuration-schema.md)** - Configuration reference
- **[Usage Guide](../usage.md)** - Usage examples
