---
sidebar_position: 1
---

# Development Getting Started

Guide for developers who want to contribute to tf2report.

## Prerequisites

- **Go** 1.25 or later
- **Git**
- **Make**
- **Terraform** (for testing)

Optional:
- **GVM** (Go Version Manager)
- **Docker** (for containerized builds)

## Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/germainlefebvre4/tf2report.git
cd tf2report
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Build

```bash
make build
```

Binary created at `bin/tf2report`.

### 4. Run Tests

```bash
make test
```

### 5. Try It Out

```bash
./bin/tf2report --plan examples/sample-plan.json
```

## Development Workflow

### Create a Feature Branch

```bash
git checkout -b feat/your-feature-name
```

### Make Changes

Edit files following the coding guidelines in `.github/instructions/`.

### Test Changes

```bash
# Build
make build

# Run tests
make test

# Test manually
./bin/tf2report --plan examples/sample-plan.json
```

### Format Code

```bash
make fmt
make tidy
```

### Commit

```bash
git add .
git commit -m "feat: add your feature description"
git push origin feat/your-feature-name
```

### Create Pull Request

Open a pull request on GitHub with:
- Clear description of changes
- Reference to any related issues
- Test results

## Available Make Targets

```bash
make          # Build (default)
make build    # Build binary
make test     # Run tests
make fmt      # Format code
make tidy     # Tidy dependencies
make clean    # Clean build artifacts
make install  # Install system-wide
make deps     # Download dependencies
make run      # Build and run
make help     # Show help
```

## Project Structure

```
tf2report/
├── cmd/
│   └── tf2report/          # CLI entry point
│       └── main.go
├── pkg/
│   ├── config/             # Configuration
│   │   ├── loader.go
│   │   └── types.go
│   ├── report/             # Report generation
│   │   ├── formatter.go
│   │   ├── markdown.go
│   │   ├── text.go
│   │   └── json.go
│   └── terraform/          # Terraform parsing
│       ├── parser.go
│       └── types.go
├── docs/                   # Documentation
├── examples/               # Example plans
├── .github/instructions/   # Coding guidelines
├── go.mod
├── Makefile
└── README.md
```

## Running Tests

### All Tests

```bash
make test
```

### Specific Package

```bash
go test ./pkg/terraform
go test ./pkg/report
go test ./pkg/config
```

### With Coverage

```bash
make test-coverage
```

Opens `coverage.html` in browser.

### Verbose

```bash
go test -v ./...
```

## Coding Guidelines

### Follow Instructions

Read and follow:
- `.github/instructions/go.instructions.md` - Go coding standards
- `.github/instructions/cli.instructions.md` - CLI practices
- `.github/instructions/makefile.instructions.md` - Makefile standards

### Key Principles

1. **Idiomatic Go** - Write clear, idiomatic Go code
2. **Error Handling** - Always handle errors properly
3. **Documentation** - Document all exported items
4. **Testing** - Write tests for new features
5. **Single Package Declaration** - One `package` statement per file

## Adding Features

### Add Output Format

1. Create `pkg/report/yourformat.go`
2. Implement `Formatter` interface
3. Add to `NewFormatter()` in `formatter.go`
4. Add tests
5. Update documentation

### Add Configuration Option

1. Add field to `Config` in `pkg/config/types.go`
2. Set default in `loader.go`
3. Add CLI flag in `cmd/tf2report/main.go`
4. Update docs

### Add Filter

1. Add logic to `pkg/terraform/parser.go`
2. Add config option if needed
3. Wire up CLI
4. Add tests
5. Update docs

## Testing with Examples

```bash
# Test all formats
./bin/tf2report --plan examples/sample-plan.json
./bin/tf2report --plan examples/sample-plan.json --format text
./bin/tf2report --plan examples/sample-plan.json --format json
```

## Debugging

### Verbose Mode

```bash
./bin/tf2report --plan examples/sample-plan.json --verbose
```

### Use Debugger

```bash
# Delve debugger
dlv debug ./cmd/tf2report -- --plan examples/sample-plan.json
```

## Documentation Updates

When changing code:
1. Update relevant docs in `docs/`
2. Update `README.md` if needed
3. Follow `.github/instructions/update-docs-on-code-change.instructions.md`

## Continuous Integration

Tests run automatically on:
- Pull requests
- Pushes to main
- Release tags

Check `.github/workflows/` for CI configuration.

## Getting Help

- Read existing code
- Check `docs/` directory
- Review `.github/instructions/`
- Open an issue on GitHub

## Next Steps

- **[Project Structure](./project-structure.md)** - Detailed structure
- **[Building](./building.md)** - Build instructions
- **[Testing](./testing.md)** - Testing guide
- **[API Reference](./api-reference.md)** - API documentation
- **[Contributing](./contributing.md)** - Contribution guidelines
