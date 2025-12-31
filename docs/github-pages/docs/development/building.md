---
sidebar_position: 3
---

# Building tf2report

Comprehensive guide to building tf2report from source.

## Prerequisites

- Go 1.25 or later
- Make
- Git

## Basic Build

```bash
make build
```

Binary created at `bin/tf2report`.

## Build Commands

### Default Build

```bash
make
# or
make all
# or
make build
```

### Development Build

Build with race detection and debug symbols:

```bash
make dev
```

### Clean Build

Remove old artifacts and rebuild:

```bash
make clean
make build
```

## Manual Build

Without Make:

```bash
go build -o bin/tf2report ./cmd/tf2report
```

## Install System-Wide

```bash
sudo make install
```

Installs to `/usr/local/bin/tf2report`.

## Cross-Platform Builds

### Build for Linux

```bash
GOOS=linux GOARCH=amd64 go build -o bin/tf2report-linux-amd64 ./cmd/tf2report
```

### Build for macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o bin/tf2report-darwin-amd64 ./cmd/tf2report
GOOS=darwin GOARCH=arm64 go build -o bin/tf2report-darwin-arm64 ./cmd/tf2report
```

### Build for Windows

```bash
GOOS=windows GOARCH=amd64 go build -o bin/tf2report-windows-amd64.exe ./cmd/tf2report
```

### Build All Platforms

```bash
make release
```

## Build Options

### Optimize for Size

```bash
go build -ldflags="-s -w" -o bin/tf2report ./cmd/tf2report
```

### Add Version Information

```bash
VERSION=$(git describe --tags --always --dirty)
go build -ldflags="-X main.version=$VERSION" -o bin/tf2report ./cmd/tf2report
```

## Verify Build

```bash
./bin/tf2report --help
./bin/tf2report --plan examples/sample-plan.json
```

## Troubleshooting

### Missing Dependencies

```bash
make deps
# or
go mod download
```

### Out of Sync Modules

```bash
make tidy
# or
go mod tidy
```

### Build Errors

```bash
make clean
go mod verify
go mod tidy
make build
```

## Next Steps

- **[Testing](./testing.md)** - Test your build
- **[Getting Started](./getting-started.md)** - Development workflow
