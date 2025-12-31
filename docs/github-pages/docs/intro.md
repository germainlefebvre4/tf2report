---
sidebar_position: 1
---

# Introduction to tf2report

Welcome to **tf2report** - a powerful command-line tool that analyzes Terraform plan files and generates human-readable reports.

## What is tf2report?

tf2report transforms Terraform plan JSON output into comprehensive reports that help you understand and communicate infrastructure changes before applying them. Whether you're reviewing changes in a pull request, documenting infrastructure updates, or analyzing security implications, tf2report makes it easy.

## Key Features

### 🎯 Multiple Output Formats
Generate reports in:
- **Markdown** - Perfect for documentation and pull requests
- **Plain Text** - Ideal for terminal output and logs
- **JSON** - Machine-readable for programmatic analysis

### 🔍 Powerful Filtering
Focus on what matters:
- Filter by resource type (e.g., only show S3 buckets)
- Filter by action (create, update, delete, replace)
- Combine filters for targeted analysis

### ⚙️ Configuration Driven
- YAML configuration files for repeatable setups
- Command-line overrides for flexibility
- Environment variable support

### 🚀 CI/CD Integration
- Seamless integration with GitHub Actions, GitLab CI, and other platforms
- Automated report generation on every plan
- Perfect for compliance and security reviews

## Why tf2report?

### Code Reviews
Generate reports for pull request comments to help reviewers understand infrastructure changes without reading raw Terraform output.

### Compliance & Audit
Create documentation of planned infrastructure changes for compliance requirements and audit trails.

### Security Reviews
Filter changes to security-related resources (IAM roles, security groups, etc.) for dedicated security team review.

### Communication
Transform technical Terraform plans into reports that stakeholders can understand.

## Quick Example

```bash
# Generate a Terraform plan
terraform plan -out=tfplan
terraform show -json tfplan > terraform.tfplan.json

# Create a report
tf2report --plan terraform.tfplan.json
```

**Output:**
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
```

## Next Steps

Ready to get started? Follow our [Getting Started Guide](./getting-started.md) to install tf2report and generate your first report!

### User Guides
- [Installation](./installation.md) - Install tf2report
- [Usage Guide](./usage.md) - Learn how to use tf2report
- [Configuration](./configuration.md) - Configure tf2report for your needs
- [Output Formats](./output-formats.md) - Understand different output formats

### Examples
- [Basic Usage](./examples/basic-usage.md) - Simple examples
- [Filtering](./examples/filtering.md) - Filter resources and actions
- [CI/CD Integration](./examples/ci-cd-integration.md) - Integrate with pipelines

### For Developers
- [Development Guide](./development/getting-started.md) - Contribute to tf2report
- [API Reference](./development/api-reference.md) - Package documentation
