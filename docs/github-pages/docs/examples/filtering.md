---
sidebar_position: 2
---

# Filtering Examples

Advanced filtering techniques and real-world scenarios.

## Security Reviews

### Filter Security Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --type aws_iam_user \
  --type aws_security_group \
  --type aws_security_group_rule \
  --type aws_kms_key \
  --format markdown > security-review.md
```

### Configuration-Based Security Review

```yaml title="security-config.yaml"
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_iam_role
    - aws_iam_policy
    - aws_security_group
    - aws_kms_key
    - aws_secretsmanager_secret
```

```bash
tf2report --config security-config.yaml > security-report.md
```

## Database Reviews

### All Database Changes

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_db_instance \
  --type aws_rds_cluster \
  --type aws_dynamodb_table \
  --type aws_elasticache_cluster \
  --format markdown
```

### Destructive Database Changes Only

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_db_instance \
  --type aws_rds_cluster \
  --action delete \
  --action replace \
  --format markdown > database-destructive.md
```

## Network Reviews

### VPC and Subnet Changes

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_vpc \
  --type aws_subnet \
  --type aws_route_table \
  --type aws_internet_gateway \
  --type aws_nat_gateway \
  --format markdown
```

### Security Group Changes Only

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_security_group \
  --type aws_security_group_rule \
  --format markdown > security-groups.md
```

## Compute Reviews

### EC2 and Auto Scaling

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_launch_template \
  --type aws_autoscaling_group \
  --format markdown
```

### Lambda Functions

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_lambda_function \
  --type aws_lambda_layer_version \
  --type aws_lambda_permission \
  --format markdown
```

## Action-Based Filtering

### New Resources Only

```bash
tf2report --plan terraform.tfplan.json --action create --format markdown
```

### Destructive Changes Alert

```bash
#!/bin/bash
tf2report --plan terraform.tfplan.json \
  --action delete \
  --action replace \
  --format json > destructive.json

COUNT=$(jq '.summary.to_destroy + .summary.to_replace' destructive.json)

if [ "$COUNT" -gt 0 ]; then
  echo "⚠️  WARNING: $COUNT destructive changes detected!"
  tf2report --plan terraform.tfplan.json --action delete --action replace
  exit 1
else
  echo "✅ No destructive changes"
fi
```

### Updates Only

```bash
tf2report --plan terraform.tfplan.json --action update --format markdown
```

## Combined Filtering

### New S3 Buckets

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --action create \
  --format markdown
```

### IAM Changes (Non-Destructive)

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_iam_role \
  --type aws_iam_policy \
  --action create \
  --action update \
  --format markdown
```

### Deleted Compute Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_instance \
  --type aws_lambda_function \
  --action delete \
  --format json
```

## Advanced jq Filtering

### List All Deleted Resources

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq -r '.changes[] | select(.action == "delete") | .address'
```

### Group Changes by Type and Action

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes | group_by(.type) | map({type: .[0].type, count: length, resources: map(.address)})'
```

### Find Specific Resource

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes[] | select(.address | contains("web"))'
```

### Count by Action

```bash
tf2report --plan terraform.tfplan.json --format json | \
  jq '.changes | group_by(.action) | map({action: .[0].action, count: length})'
```

## Multi-Report Workflow

```bash
#!/bin/bash
PLAN="terraform.tfplan.json"

# Generate multiple focused reports
echo "Generating filtered reports..."

# Security
tf2report --plan $PLAN \
  --type aws_iam_role --type aws_security_group --type aws_kms_key \
  > reports/security-review.md

# Database
tf2report --plan $PLAN \
  --type aws_db_instance --type aws_dynamodb_table \
  > reports/database-review.md

# Destructive changes
tf2report --plan $PLAN \
  --action delete --action replace \
  > reports/destructive-changes.md

# All changes
tf2report --plan $PLAN > reports/full-report.md

echo "✅ Reports generated in reports/"
ls -lh reports/
```

## Environment-Specific Filters

### Production (Strict)

```yaml title="prod-filter.yaml"
filters:
  resource_types:
    - aws_instance
    - aws_rds_cluster
    - aws_s3_bucket
  actions:
    - delete
    - replace
```

### Staging (Moderate)

```yaml title="staging-filter.yaml"
filters:
  actions:
    - delete
    - replace
```

### Development (Permissive)

```yaml title="dev-filter.yaml"
filters:
  resource_types: []
  actions: []
```

## Compliance Workflows

### GDPR-Relevant Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_s3_bucket \
  --type aws_kms_key \
  --type aws_rds_cluster \
  --type aws_dynamodb_table \
  --format markdown > gdpr-review.md
```

### PCI-DSS Resources

```bash
tf2report --plan terraform.tfplan.json \
  --type aws_kms_key \
  --type aws_cloudtrail \
  --type aws_config_configuration_recorder \
  --type aws_guardduty_detector \
  --format markdown > pci-review.md
```

## Next Steps

- **[CI/CD Integration](./ci-cd-integration.md)** - Automate filtered reports
- **[Filtering Guide](../filtering.md)** - Complete filtering reference
