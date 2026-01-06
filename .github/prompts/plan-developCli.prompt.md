# TF2Report CLI Development Plan

## Project Overview

TF2Report is a Terraform plan analysis CLI tool that parses JSON plan files and generates resource change reports. The project is starting from scratch and needs the complete Go project structure, CLI framework (Cobra), configuration support (Viper), build automation (Makefile), and documentation. Development will follow Go 1.25 standards and idiomatic practices.

## Development Steps

### 1. Initialize Go Module and Project Structure

Create foundational Go project with:
- `go.mod` - Go module file with Go 1.25 dependency
- `go.sum` - Dependency lock file
- `cmd/tf2report/main.go` - Main CLI entry point
- `pkg/` - Core parsing and report generation packages
- `docs/` - Documentation directory

Follow Go conventions:
- Single-word package names in lowercase
- `cmd/` for executables
- `pkg/` or `internal/` for library code
- Standard directory structure for CLI tools

### 2. Implement Core Packages

Build foundational packages in `pkg/`:

**`pkg/terraform/`**
- Parse Terraform JSON plan files
- Extract resource changes (additions, deletions, modifications)
- Data structures representing plan resources and changes

**`pkg/report/`**
- Generate reports in multiple formats (Markdown, plain text, JSON)
- Support resource filtering by type and change action
- Format output according to specified format

**`pkg/config/`**
- YAML configuration file parsing using Viper
- Configuration schema definition
- Default values and validation

### 3. Build CLI Command Structure with Cobra

Implement `cmd/tf2report/main.go` with:
- Root command with help text
- Subcommands for different output formats if needed
- Flags for:
  - Input terraform plan file path
  - Output format (markdown, text, json)
  - Resource type filtering
  - Change action filtering (add, modify, delete)
  - Verbosity level
  - Config file path
- Proper error handling and user feedback

### 4. Integrate Viper for Configuration

Support YAML configuration with:
- `tf2report.yaml` file in current directory or specified path
- Configuration schema:
  - `terraform_plan_path` - Path to terraform plan file
  - `output_format` - Default output format (markdown, text, json)
  - `filters` - Resource type and action filters
  - `verbosity` - Logging level (debug, info, warn, error)
- CLI flags override configuration file values
- Sensible defaults for all options

### 5. Create Build Automation with Makefile

Implement `Makefile` with standard targets:
- `make build` - Compile the CLI binary
- `make test` - Run all tests
- `make clean` - Remove build artifacts
- `make install` - Install binary to system
- Optional: `make release` - Build for multiple platforms

Follow GNU Make best practices:
- Use `.PHONY` for non-file targets
- Tab-indented recipes
- Clear target descriptions
- Proper dependency management

### 6. Write Comprehensive Documentation

Create documentation in `docs/` directory:
- `USAGE.md` - Command-line usage guide
- `CONFIGURATION.md` - Configuration file reference
- `OUTPUT_FORMATS.md` - Details on each output format
- `EXAMPLES.md` - Usage examples with different scenarios
- `API.md` - Package API reference for developers
- `DEVELOPMENT.md` - Contributing and development setup

Keep documentation synchronized with code changes using the update-docs-on-code-change instructions.

## Key Features to Implement

### Output Formats

**Markdown (default)**
- Human-readable summary of resource changes
- Tables showing resources by action
- Section organization

**Plain Text**
- Simple text-based summary
- Formatted for easy terminal reading

**JSON**
- Structured format for tool integration
- Complete change information

### Filtering Capabilities

- Filter by resource type (aws_instance, aws_s3_bucket, etc.)
- Filter by change action (add, modify, delete)
- Combine multiple filters with AND/OR logic

### Configuration File Support

Create YAML configuration schema:
```yaml
terraform_plan_path: ./terraform.tfplan
output_format: markdown
filters:
  resource_types:
    - aws_instance
    - aws_s3_bucket
  actions:
    - add
    - modify
verbosity: info
```

## Constraints and Standards

### Code Quality
- Follow idiomatic Go practices
- Single-word package names
- MixedCaps naming conventions
- Proper error handling with `%w` format
- No emojis in code, comments, or documentation

### CLI Standards
- Use Cobra for command structure
- Comprehensive help text for all commands
- Clear error messages for user feedback
- Support for bash/zsh completion (optional)

### Documentation Standards
- Documentation only in `docs/` and README.md
- No example outputs in code comments or help text
- Keep documentation synchronized with code changes

## Testing Strategy

- Unit tests for parsing, formatting, and filtering logic
- Use Go's `testing` package
- Example Terraform files from `examples/` for integration tests
- Test coverage for all public functions

## File Structure After Completion

```
tf2report_02/
├── README.md
├── Makefile
├── go.mod
├── go.sum
├── cmd/
│   └── tf2report/
│       └── main.go
├── pkg/
│   ├── terraform/
│   │   ├── parser.go
│   │   └── types.go
│   ├── report/
│   │   ├── formatter.go
│   │   ├── markdown.go
│   │   ├── text.go
│   │   └── json.go
│   └── config/
│       ├── loader.go
│       └── types.go
├── docs/
│   ├── USAGE.md
│   ├── CONFIGURATION.md
│   ├── OUTPUT_FORMATS.md
│   ├── EXAMPLES.md
│   ├── API.md
│   └── DEVELOPMENT.md
├── examples/
│   ├── iac-Domains_baff-d.json
│   ├── iac-Domains_ledgers-d.json
│   └── iac-Stream_kafka-d.json
└── .github/
    └── instructions/
        ├── cli.instructions.md
        ├── main.instructions.md
        ├── go.instructions.md
        ├── makefile.instructions.md
        └── update-docs-on-code-change.instructions.md
```

## Next Steps

1. Create `go.mod` with Go 1.25 and add Cobra and Viper dependencies
2. Initialize project structure directories
3. Implement terraform plan parser in `pkg/terraform/`
4. Build report formatters in `pkg/report/`
5. Create configuration loader in `pkg/config/`
6. Build CLI command structure with Cobra in `cmd/tf2report/main.go`
7. Create Makefile with build targets
8. Write comprehensive documentation in `docs/`
9. Add tests for all packages
10. Test with example Terraform plans in `examples/`
