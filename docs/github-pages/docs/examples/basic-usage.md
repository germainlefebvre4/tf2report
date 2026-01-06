---
sidebar_position: 1
---

# Basic Usage Examples

Practical examples to get you started with tf2report.

## Prerequisites

All examples assume you have:
1. Terraform installed
2. tf2report installed
3. A Terraform configuration ready

## Example 1: First Report

Generate your first report:

```bash
# Create a plan
terraform plan -out=tfplan

# Convert to JSON
terraform show -json tfplan > terraform.tfplan.json

# Generate report
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

## Example 2: Different Output Formats

### Markdown (Default)

```bash
tf2report --plan terraform.tfplan.json
```

### Plain Text

```bash
tf2report --plan terraform.tfplan.json --format text
```

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

### JSON

```bash
tf2report --plan terraform.tfplan.json --format json
```

```json
{
  "terraform_version": "1.10.5",
  "summary": {
    "total_changes": 3,
    "to_add": 2,
    "to_change": 1,
    "to_destroy": 0,
    "to_replace": 0
  }
}
```

## Example 3: Save to File

Save report output to files:

```bash
# Markdown
tf2report --plan terraform.tfplan.json > report.md

# Text
tf2report --plan terraform.tfplan.json --format text > report.txt

# JSON
tf2report --plan terraform.tfplan.json --format json > report.json
```

## Example 4: View with Pager

For large reports, use a pager:

```bash
tf2report --plan terraform.tfplan.json | less
```

## Example 5: Quick Summary

Get just the summary using jq:

```bash
tf2report --plan terraform.tfplan.json --format json | jq '.summary'
```

**Output:**
```json
{
  "total_changes": 3,
  "to_add": 2,
  "to_change": 1,
  "to_destroy": 0,
  "to_replace": 0
}
```

## Example 6: Check for Changes

Check if there are any changes:

```bash
#!/bin/bash
CHANGES=$(tf2report --plan terraform.tfplan.json --format json | jq '.summary.total_changes')

if [ "$CHANGES" -eq 0 ]; then
  echo "✅ No changes detected"
else
  echo "📝 Found $CHANGES changes"
  tf2report --plan terraform.tfplan.json
fi
```

## Example 7: Using Configuration File

Create a configuration file:

```yaml title="tf2report.yaml"
terraform_plan_path: terraform.tfplan.json
output_format: markdown
verbosity: info
```

Use it:

```bash
# Uses ./tf2report.yaml automatically
tf2report
```

## Example 8: Simple Filtering

### Show Only S3 Buckets

```bash
tf2report --plan terraform.tfplan.json --type aws_s3_bucket
```

### Show Only New Resources

```bash
tf2report --plan terraform.tfplan.json --action create
```

## Example 9: Count Resources by Type

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq '.by_type | to_entries[] | "\(.key): \(.value.to_add + .value.to_change + .value.to_destroy + .value.to_replace)"'
```

**Output:**
```
aws_instance: 2
aws_s3_bucket: 1
```

## Example 10: Verbose Output

See detailed logging:

```bash
tf2report --plan terraform.tfplan.json --verbose
```

## Complete Workflow Example

A complete workflow from Terraform plan to report:

```bash
#!/bin/bash
set -e

echo "📋 Creating Terraform plan..."
terraform plan -out=tfplan

echo "📄 Converting plan to JSON..."
terraform show -json tfplan > terraform.tfplan.json

echo "📊 Generating report..."
tf2report --plan terraform.tfplan.json > report.md

echo "✅ Report generated: report.md"

# Display summary
echo ""
echo "Summary:"
tf2report --plan terraform.tfplan.json --format json | jq '.summary'

# Open report
if command -v bat &> /dev/null; then
  bat report.md
else
  cat report.md
fi
```

## Using with Make

Add tf2report to your Makefile:

```makefile title="Makefile"
.PHONY: plan report

plan:
	terraform plan -out=tfplan
	terraform show -json tfplan > terraform.tfplan.json

report: plan
	tf2report --plan terraform.tfplan.json > report.md
	@echo "Report generated: report.md"

clean:
	rm -f tfplan terraform.tfplan.json report.md
```

Usage:

```bash
make report
```

## Common Patterns

### Pattern 1: Daily Report

```bash
#!/bin/bash
DATE=$(date +%Y-%m-%d)
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
tf2report --plan plan.json > reports/daily-${DATE}.md
```

### Pattern 2: Email Report

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json --format text | \
  mail -s "Terraform Plan Report" team@example.com
```

### Pattern 3: Slack Notification

```bash
#!/bin/bash
REPORT=$(tf2report --plan terraform.tfplan.json --format text)
curl -X POST -H 'Content-type: application/json' \
  --data "{\"text\":\"$REPORT\"}" \
  YOUR_SLACK_WEBHOOK_URL
```

## Next Steps

- **[Filtering Examples](./filtering.md)** - Learn advanced filtering
- **[CI/CD Integration](./ci-cd-integration.md)** - Automate reports
