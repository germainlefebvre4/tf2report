---
sidebar_position: 7
---

# Filtering

Master tf2report's filtering capabilities to focus on relevant infrastructure changes.

## Overview

Filtering allows you to narrow down reports to specific resource types and actions, making it easier to review targeted changes.

## Filter Types

tf2report supports two types of filters:

1. **Resource Type Filters** - Filter by resource type (e.g., `aws_instance`)
2. **Action Filters** - Filter by change action (create, update, delete, replace)

Filters can be combined for precise targeting.

## Resource Type Filtering

### Single Resource Type

Show changes for one resource type:

```bash
tf2report --plan terraform.tfplan.json --type aws_s3_bucket
```

### Multiple Resource Types

Show changes for multiple resource types:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_s3_bucket \
  --type aws_security_group
```

### Common Resource Type Patterns

#### Compute Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_launch_template \
  --type aws_autoscaling_group \
  --type aws_lambda_function \
  --type aws_ecs_service
```

#### Storage Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --type aws_ebs_volume \
  --type aws_efs_file_system
```

#### Database Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_db_instance \
  --type aws_rds_cluster \
  --type aws_dynamodb_table \
  --type aws_elasticache_cluster
```

#### Network Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_vpc \
  --type aws_subnet \
  --type aws_route_table \
  --type aws_security_group \
  --type aws_network_acl
```

#### Security Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_security_group \
  --type aws_kms_key
```

## Action Filtering

### Single Action

Show only resources with a specific action:

```bash
# Only new resources
tf2report --plan terraform.tfplan.json --action create

# Only updated resources
tf2report --plan terraform.tfplan.json --action update

# Only deleted resources
tf2report --plan terraform.tfplan.json --action delete

# Only replaced resources
tf2report --plan terraform.tfplan.json --action replace
```

### Multiple Actions

Combine multiple actions:

```bash
# Destructive changes only
tf2report --plan terraform.tfplan.json --action delete --action replace

# Non-destructive changes
tf2report --plan terraform.tfplan.json --action create --action update
```

### Action Types

| Action | Description | Destructive |
|--------|-------------|-------------|
| `create` | Resources being added | No |
| `update` | Resources being modified in-place | No |
| `delete` | Resources being removed | Yes |
| `replace` | Resources being destroyed and recreated | Yes |

## Combined Filtering

Combine resource type and action filters for precise targeting.

### New S3 Buckets

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --action create
```

### Destructive Database Changes

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_db_instance \
  --type aws_rds_cluster \
  --action delete \
  --action replace
```

### Security Resource Updates

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_security_group \
  --action update
```

## Configuration File Filtering

Set default filters in configuration files:

```yaml title="tf2report.yaml"
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
    - aws_rds_cluster
  actions:
    - create
    - delete
```

Use the configuration:

```bash
tf2report --plan terraform.tfplan.json
```

Override configuration filters:

```bash
# Override to show all types
tf2report --plan terraform.tfplan.json --type aws_instance

# Override to show all actions
tf2report --plan terraform.tfplan.json --action create
```

## Use Case Examples

### Security Review

Focus on security-related resources:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_iam_user \
  --type aws_iam_group \
  --type aws_security_group \
  --type aws_security_group_rule \
  --type aws_kms_key \
  --type aws_kms_alias \
  --type aws_secretsmanager_secret \
  --format markdown > security-review.md
```

### Database Change Review

Focus on database changes:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_db_instance \
  --type aws_db_subnet_group \
  --type aws_db_parameter_group \
  --type aws_rds_cluster \
  --type aws_rds_cluster_instance \
  --type aws_dynamodb_table \
  --format markdown > database-review.md
```

### Destructive Changes Check

Identify all destructive changes:

```bash
tf2report --plan terraform.tfplan.json \
  --action delete \
  --action replace \
  --format json > destructive-changes.json

# Check if any destructive changes exist
if [ $(jq '.summary.to_destroy + .summary.to_replace' destructive-changes.json) -gt 0 ]; then
  echo "⚠️  WARNING: Destructive changes detected!"
  tf2report --plan terraform.tfplan.json --action delete --action replace
  exit 1
fi
```

### Network Changes Review

Focus on networking resources:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_vpc \
  --type aws_subnet \
  --type aws_route_table \
  --type aws_route \
  --type aws_internet_gateway \
  --type aws_nat_gateway \
  --type aws_network_acl \
  --format markdown > network-review.md
```

### Compute Changes Review

Focus on compute resources:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_launch_template \
  --type aws_autoscaling_group \
  --type aws_lambda_function \
  --type aws_ecs_cluster \
  --type aws_ecs_service \
  --type aws_ecs_task_definition \
  --format markdown > compute-review.md
```

### Compliance Review

Check changes to compliance-critical resources:

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --type aws_s3_bucket_public_access_block \
  --type aws_kms_key \
  --type aws_cloudtrail \
  --type aws_config_configuration_recorder \
  --type aws_guardduty_detector \
  --format markdown > compliance-review.md
```

## Advanced Filtering Techniques

### Using jq for Post-Filtering

Generate JSON and use jq for advanced filtering:

```bash
# Get all resources being deleted
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes[] | select(.action == "delete")'

# Get deletions of a specific type
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes[] | select(.action == "delete" and .type == "aws_s3_bucket")'

# Count changes by action
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes | group_by(.action) | map({action: .[0].action, count: length})'
```

### Multiple Filter Scenarios

Create different filtered reports:

```bash
#!/bin/bash
PLAN="terraform.tfplan.json"

# Security review
tf2report --plan $PLAN --type aws_iam_role --type aws_security_group \
  > security-review.md

# Database review
tf2report --plan $PLAN --type aws_db_instance --type aws_dynamodb_table \
  > database-review.md

# Destructive changes
tf2report --plan $PLAN --action delete --action replace \
  > destructive-changes.md

# All changes
tf2report --plan $PLAN > all-changes.md
```

### Environment-Specific Filters

Different filters for different environments:

```yaml title="prod-config.yaml"
filters:
  # Production: strict filtering
  resource_types:
    - aws_instance
    - aws_rds_cluster
    - aws_s3_bucket
  actions:
    - delete
    - replace
```

```yaml title="dev-config.yaml"
filters:
  # Development: minimal filtering
  actions: []
  resource_types: []
```

## Filter Behavior

### Empty Filters

When no filters are specified, all resources and actions are included:

```bash
# Shows all changes
tf2report --plan terraform.tfplan.json
```

### Multiple Values

Multiple filter values use OR logic:

```bash
# Shows aws_instance OR aws_s3_bucket
tf2report --plan terraform.tfplan.json --type aws_instance --type aws_s3_bucket

# Shows create OR delete
tf2report --plan terraform.tfplan.json --action create --action delete
```

### Combined Filters

Type and action filters use AND logic:

```bash
# Shows aws_instance AND create
# (new EC2 instances only)
tf2report --plan terraform.tfplan.json --type aws_instance --action create
```

## Tips and Best Practices

### 1. Start Broad, Narrow Down

Begin with no filters to see all changes, then narrow down:

```bash
# See everything first
tf2report --plan terraform.tfplan.json

# Then focus on specific areas
tf2report --plan terraform.tfplan.json --type aws_instance
```

### 2. Use Configuration for Common Filters

Save frequently used filters in configuration files:

```yaml title="tf2report.yaml"
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
```

### 3. Create Review Workflows

Establish filtered review workflows:

1. **Quick Review** - Destructive changes only
2. **Security Review** - Security resources only
3. **Database Review** - Database resources only
4. **Full Review** - All changes

### 4. Combine with Output Formats

Use appropriate formats for filtered output:

```bash
# Markdown for human review
tf2report --plan plan.json --type aws_s3_bucket --format markdown

# JSON for automation
tf2report --plan plan.json --action delete --format json
```

### 5. Document Filter Rationale

Document why specific filters are used:

```yaml title="security-review.yaml"
# Security team review configuration
# Focuses on IAM, security groups, and encryption
filters:
  resource_types:
    - aws_iam_role
    - aws_iam_policy
    - aws_security_group
    - aws_kms_key
```

## Next Steps

- **[Examples](./examples/basic-usage.md)** - See practical filtering examples
- **[CI/CD Integration](./examples/ci-cd-integration.md)** - Automate filtered reports
- **[Configuration](./configuration.md)** - Configure default filters
