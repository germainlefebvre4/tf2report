# Main Instructions

This project is a CLI tool to analyze Terraform plan files and generate resource change reports.

## Branding

- Name: tf2report
- Description: Summarize Terraform plan changes into human-readable reports

## Features

- Parse Terraform plan files in JSON format
- Extract resource changes from the plan
- Summarize resource additions, deletions, and modifications
- Output summary in multiple formats: Markdown (default), Plain text, JSON
- Support filtering by resource type and change action

## Configuration

- Provide (optional) configuration file in YAML format
- Provide the file of the terraform plan as cli argument or in config file. Set a default path if not provided.
- Configure output format
- Configure resource type and action filters
- Configure verbosity level for logging

## Technology

### CLI

- Golang

## Documentation

- README.md file in the root of the repository
- Documentation directory: docs/
- Do not write documentation in root directory except README.md
- Do not use emojis in the documentation

## Examples

- Example Terraform plan files in examples/ directory
- Do not include example outputs anywhere (documentation, code comments, tests, configuration files, CLI help messages, code)
