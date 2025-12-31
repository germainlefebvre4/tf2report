# tf2report Examples

This document provides practical examples of using `tf2report` in various scenarios.

## Prerequisites

All examples assume you have generated a Terraform plan in JSON format:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > terraform.tfplan.json
```

## Basic Examples

### Generate Default Markdown Report

```bash
tf2report --plan terraform.tfplan.json
```

### Generate Text Report

```bash
tf2report --plan terraform.tfplan.json --format text
```

### Generate JSON Report

```bash
tf2report --plan terraform.tfplan.json --format json
```

## Filtering Examples

### Filter by Single Resource Type

Show only AWS S3 bucket changes:

```bash
tf2report --plan terraform.tfplan.json --type aws_s3_bucket
```

### Filter by Multiple Resource Types

Show only EC2 instances and security groups:

```bash
tf2report --plan terraform.tfplan.json --type aws_instance --type aws_security_group
```

### Filter by Action

Show only resources being created:

```bash
tf2report --plan terraform.tfplan.json --action create
```

Show only destructive changes:

```bash
tf2report --plan terraform.tfplan.json --action delete --action replace
```

### Combine Type and Action Filters

Show only new S3 buckets being created:

```bash
tf2report --plan terraform.tfplan.json --type aws_s3_bucket --action create
```

## Output Redirection Examples

### Save to File

```bash
tf2report --plan terraform.tfplan.json > changes.md
```

### Save JSON for Further Processing

```bash
tf2report --plan terraform.tfplan.json --format json > changes.json
jq '.summary.total_changes' changes.json
```

## Configuration File Examples

### Using a Configuration File

Create `tf2report.yaml`:

```yaml
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
```

Run with configuration:

```bash
tf2report
```

### Override Configuration with Flags

```bash
tf2report --format json
```

## CI/CD Integration Examples

### GitHub Actions

```yaml
- name: Generate Terraform Plan
  run: |
    terraform plan -out=tfplan
    terraform show -json tfplan > plan.json

- name: Generate Report
  run: |
    tf2report --plan plan.json --format markdown > report.md

- name: Comment on PR
  uses: actions/github-script@v6
  with:
    script: |
      const fs = require('fs');
      const report = fs.readFileSync('report.md', 'utf8');
      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body: report
      });
```

### GitLab CI

```yaml
terraform:plan:
  script:
    - terraform plan -out=tfplan
    - terraform show -json tfplan > plan.json
    - tf2report --plan plan.json --format markdown > report.md
  artifacts:
    paths:
      - report.md
```

### Check for Destructive Changes

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json --format json > report.json

deletions=$(jq '.summary.to_destroy' report.json)
replacements=$(jq '.summary.to_replace' report.json)

if [ "$deletions" -gt 0 ] || [ "$replacements" -gt 0 ]; then
  echo "Warning: Destructive changes detected!"
  echo "Deletions: $deletions"
  echo "Replacements: $replacements"
  exit 1
fi
```

## Advanced Examples

### Multiple Environments

```bash
# Production
tf2report --plan prod.tfplan.json --config prod-config.yaml > prod-report.md

# Staging
tf2report --plan staging.tfplan.json --config staging-config.yaml > staging-report.md
```

### Filtering for Security Review

Show only IAM and security group changes:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_security_group \
  --type aws_security_group_rule \
  --format markdown > security-review.md
```

### Verbose Debugging

```bash
tf2report --plan terraform.tfplan.json --verbose 2>&1 | tee debug.log
```

## Shell Script Integration

### Automated Report Generation

```bash
#!/bin/bash
set -e

PLAN_FILE="terraform.tfplan.json"
REPORT_FILE="terraform-changes-$(date +%Y%m%d-%H%M%S).md"

# Generate report
tf2report --plan "$PLAN_FILE" > "$REPORT_FILE"

# Display summary
echo "Report generated: $REPORT_FILE"
echo "---"
head -20 "$REPORT_FILE"
```

### Conditional Approval

```bash
#!/bin/bash

# Generate JSON report
tf2report --plan terraform.tfplan.json --format json > report.json

# Extract counts
total=$(jq '.summary.total_changes' report.json)
deletions=$(jq '.summary.to_destroy' report.json)

if [ "$total" -eq 0 ]; then
  echo "No changes detected."
  exit 0
fi

echo "Total changes: $total"
echo "Deletions: $deletions"

if [ "$deletions" -gt 5 ]; then
  echo "Warning: More than 5 resources will be deleted."
  echo "Manual approval required."
  exit 1
fi

echo "Changes approved."
```

## Makefile Integration

```makefile
.PHONY: plan-report
plan-report:
	terraform plan -out=tfplan
	terraform show -json tfplan > plan.json
	tf2report --plan plan.json > PLAN_REPORT.md
	cat PLAN_REPORT.md
```

## Docker Integration

```dockerfile
FROM hashicorp/terraform:latest

COPY --from=builder /app/tf2report /usr/local/bin/tf2report

WORKDIR /workspace

CMD terraform plan -out=tfplan && \
    terraform show -json tfplan > plan.json && \
    tf2report --plan plan.json
```
