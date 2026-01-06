---
sidebar_position: 2
---

# Getting Started

Get up and running with tf2report in 5 minutes!

## Prerequisites

Before you start, ensure you have:
- **Terraform** installed (any recent version)
- **Go 1.25+** (if building from source)
- A Terraform configuration ready to plan

## Installation

Choose your preferred installation method:

### Using Go Install (Recommended)

The easiest way to install tf2report:

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

Verify the installation:

```bash
tf2report --help
```

### From Source

Build from the latest source code:

```bash
# Clone the repository
git clone https://github.com/germainlefebvre4/tf2report.git
cd tf2report

# Build the binary
make build

# Install to /usr/local/bin (requires sudo)
sudo make install
```

### Docker (Coming Soon)

```bash
docker pull germainlefebvre4/tf2report:latest
```

## Your First Report

Let's generate your first Terraform plan report!

### Step 1: Create a Terraform Plan

Navigate to your Terraform project and create a plan:

```bash
cd /path/to/your/terraform/project
terraform plan -out=tfplan
```

### Step 2: Convert Plan to JSON

Terraform's JSON output is what tf2report analyzes:

```bash
terraform show -json tfplan > terraform.tfplan.json
```

### Step 3: Generate a Report

Now run tf2report:

```bash
tf2report --plan terraform.tfplan.json
```

🎉 **Congratulations!** You've generated your first report.

## Understanding the Output

The default Markdown output includes:

1. **Header** - Terraform version information
2. **Summary** - Total changes and counts by action type
3. **Changes by Resource Type** - Detailed breakdown of each resource

Example output:

```markdown
# Terraform Plan Summary

**Terraform Version:** 1.10.5

## Summary

**Total Changes:** 3

| Action | Count |
|--------|-------|
| Add | 2 |
| Change | 1 |

## Changes by Resource Type

### aws_s3_bucket (2)

**To Add (2):**
- `aws_s3_bucket.data`
- `aws_s3_bucket.logs`

### aws_instance (1)

**To Change (1):**
- `aws_instance.web[0]`
```

## Try Different Formats

### Plain Text

Great for terminal viewing:

```bash
tf2report --plan terraform.tfplan.json --format text
```

### JSON

Perfect for automation and scripts:

```bash
tf2report --plan terraform.tfplan.json --format json
```

## Save Output to a File

Redirect output to save your report:

```bash
# Markdown
tf2report --plan terraform.tfplan.json > report.md

# JSON
tf2report --plan terraform.tfplan.json --format json > report.json
```

## Filter Changes

Focus on specific resources or actions:

```bash
# Only show S3 bucket changes
tf2report --plan terraform.tfplan.json --type aws_s3_bucket

# Only show resources being created
tf2report --plan terraform.tfplan.json --action create

# Combine filters
tf2report --plan terraform.tfplan.json --type aws_s3_bucket --action create
```

## Use a Configuration File

For repeated use, create a configuration file:

```yaml title="tf2report.yaml"
terraform_plan_path: terraform.tfplan.json
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
  actions:
    - create
    - delete
verbosity: info
```

Then simply run:

```bash
tf2report
```

## Next Steps

Now that you've generated your first report, explore more features:

- **[Usage Guide](./usage.md)** - Comprehensive usage instructions
- **[Configuration](./configuration.md)** - Configure tf2report for your needs
- **[Output Formats](./output-formats.md)** - Learn about different output formats
- **[Examples](./examples/basic-usage.md)** - See more practical examples
- **[CI/CD Integration](./examples/ci-cd-integration.md)** - Integrate with your pipeline

## Troubleshooting

### Command not found

If `tf2report` is not found after installation:

1. Check your `$GOPATH/bin` is in your `$PATH`:
   ```bash
   echo $PATH | grep -q "$GOPATH/bin" || echo "Add $GOPATH/bin to PATH"
   ```

2. Add to your shell profile (`.bashrc`, `.zshrc`, etc.):
   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

### Invalid JSON format

Make sure you're using the JSON output from `terraform show -json`, not the plan file directly.

### Need Help?

- Check the [Usage Guide](./usage.md) for detailed documentation
- Report issues on [GitHub](https://github.com/germainlefebvre4/tf2report/issues)
