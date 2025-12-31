# tf2report Development Guide

This guide is for developers who want to contribute to or modify the `tf2report` project.

## Prerequisites

- Go 1.25 or later
- Git
- Make

### Optional

- GVM (Go Version Manager) for managing Go versions

## Project Structure

```
tf2report/
├── cmd/
│   └── tf2report/          # CLI entry point
│       └── main.go
├── pkg/
│   ├── config/             # Configuration loading
│   │   ├── loader.go
│   │   └── types.go
│   ├── report/             # Report generation
│   │   ├── formatter.go
│   │   ├── markdown.go
│   │   ├── text.go
│   │   └── json.go
│   └── terraform/          # Terraform plan parsing
│       ├── parser.go
│       └── types.go
├── docs/                   # Documentation
├── examples/               # Example Terraform plans
├── .github/
│   └── instructions/       # Coding guidelines
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Build automation
└── README.md
```

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/germainlefebvre4/tf2report.git
cd tf2report
```

### Install Dependencies

```bash
make deps
```

### Build the Project

```bash
make build
```

The binary will be created in `bin/tf2report`.

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Tidy Dependencies

```bash
make tidy
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feat/your-feature-name
```

### 2. Make Changes

Edit the relevant files following the coding guidelines in `.github/instructions/`.

### 3. Test Your Changes

```bash
make build
./bin/tf2report --plan examples/sample-plan.json
```

### 4. Format and Verify

```bash
make fmt
make tidy
make verify
```

### 5. Commit Changes

```bash
git add .
git commit -m "Add your feature description"
```

## Coding Standards

### Go Code

Follow the guidelines in `.github/instructions/go.instructions.md`:

- Use idiomatic Go practices
- Single-word package names in lowercase
- MixedCaps for exported names, mixedCaps for unexported
- Document all exported functions, types, and packages
- Handle errors properly with `fmt.Errorf` and `%w`
- Keep the happy path left-aligned

### Package Declarations

**CRITICAL:** Each Go file must have exactly ONE `package` declaration. When editing existing files, preserve the existing package declaration. When creating new files in an existing directory, use the same package name as other files in that directory.

### CLI Code

Follow the guidelines in `.github/instructions/cli.instructions.md`:

- Use Cobra for command structure
- Use Viper for configuration
- Provide clear help text
- Support both flags and configuration files

### Makefile

Follow the guidelines in `.github/instructions/makefile.instructions.md`:

- Use tabs for indentation in recipes
- Declare phony targets with `.PHONY`
- Use descriptive target names
- Include help text for targets

## Adding New Features

### Adding a New Output Format

1. Create a new formatter in `pkg/report/`
2. Implement the `Formatter` interface
3. Add the format to `NewFormatter()` in `formatter.go`
4. Update documentation

### Adding New Filters

1. Add filter logic in `pkg/terraform/parser.go`
2. Add configuration options in `pkg/config/types.go`
3. Wire up CLI flags in `cmd/tf2report/main.go`
4. Update documentation

### Adding New Configuration Options

1. Add fields to `Config` in `pkg/config/types.go`
2. Set defaults in `loader.go`
3. Add CLI flags in `main.go`
4. Update configuration documentation

## Testing

### Running Tests

```bash
make test
```

### Running Tests with Coverage

```bash
make test-coverage
```

This generates `coverage.html` which can be opened in a browser.

### Testing with Example Files

```bash
./bin/tf2report --plan examples/sample-plan.json
./bin/tf2report --plan examples/sample-plan.json --format text
./bin/tf2report --plan examples/sample-plan.json --format json
```

## Documentation

### Updating Documentation

When you change code that affects user-facing behavior or APIs:

1. Update relevant documentation in `docs/`
2. Update `README.md` if necessary
3. Follow `.github/instructions/update-docs-on-code-change.instructions.md`

### Documentation Files

- `USAGE.md` - Command-line usage
- `CONFIGURATION.md` - Configuration file reference
- `OUTPUT_FORMATS.md` - Output format details
- `EXAMPLES.md` - Usage examples
- `API.md` - Package API reference
- `DEVELOPMENT.md` - This file

## Building for Release

### Build for Multiple Platforms

```bash
make release
```

This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Binaries are placed in `bin/` with platform-specific names.

## Makefile Targets

- `make all` - Build the binary (default)
- `make build` - Build the binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make clean` - Remove build artifacts
- `make install` - Install binary to system
- `make uninstall` - Remove binary from system
- `make fmt` - Format Go code
- `make tidy` - Tidy Go modules
- `make verify` - Verify dependencies
- `make deps` - Download dependencies
- `make run` - Build and run
- `make release` - Build for multiple platforms
- `make help` - Show help

## Troubleshooting

### Build Errors

If you encounter build errors:

```bash
make clean
make deps
make tidy
make build
```

### Import Errors

```bash
go mod tidy
go mod verify
```

### Test Failures

Run tests with verbose output:

```bash
go test -v ./...
```

## Contributing Guidelines

1. Follow the coding standards in `.github/instructions/`
2. Write tests for new functionality
3. Update documentation
4. Format code with `make fmt`
5. Ensure tests pass with `make test`
6. Create clear, descriptive commit messages
7. Keep changes focused and atomic

## Getting Help

- Check existing documentation in `docs/`
- Review code examples in the repository
- Examine existing tests for patterns
- Read the Go coding instructions in `.github/instructions/go.instructions.md`
