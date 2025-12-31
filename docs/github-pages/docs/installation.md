---
sidebar_position: 3
---

# Installation

Complete installation guide for tf2report across different platforms and methods.

## System Requirements

- **Operating System:** Linux, macOS, or Windows
- **Go:** 1.25 or later (for building from source)
- **Terraform:** Any recent version (for generating plans)

## Installation Methods

### Using Go Install

The simplest and recommended method:

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

This installs the latest version to `$GOPATH/bin/tf2report`.

#### Verify Installation

```bash
tf2report --help
```

#### Update to Latest Version

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

### From Source

Build from source for the latest development version or to contribute:

#### Clone Repository

```bash
git clone https://github.com/germainlefebvre4/tf2report.git
cd tf2report
```

#### Build

```bash
make build
```

This creates the binary at `bin/tf2report`.

#### Install System-Wide (Optional)

```bash
sudo make install
```

This copies the binary to `/usr/local/bin/tf2report`.

#### Development Build

For development with additional checks:

```bash
make dev
```

### Using Docker (Coming Soon)

Docker images will be available soon:

```bash
# Pull the image
docker pull germainlefebvre4/tf2report:latest

# Run tf2report
docker run --rm -v $(pwd):/workspace germainlefebvre4/tf2report:latest \
  --plan /workspace/terraform.tfplan.json
```

### Download Pre-built Binaries (Coming Soon)

Pre-built binaries for each release will be available on the [GitHub Releases](https://github.com/germainlefebvre4/tf2report/releases) page.

## Platform-Specific Instructions

### Linux

#### Using Go Install

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

#### Add to PATH

If the command isn't found, add Go's bin directory to your PATH:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

For other shells (zsh, fish), modify the appropriate rc file.

### macOS

#### Using Go Install

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

#### Add to PATH

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

#### Using Homebrew (Coming Soon)

```bash
brew install tf2report
```

### Windows

#### Using Go Install

```powershell
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

The binary will be installed to `%GOPATH%\bin\tf2report.exe`.

#### Add to PATH

Add `%GOPATH%\bin` to your system PATH:

1. Open System Properties → Advanced → Environment Variables
2. Edit the `Path` variable
3. Add `C:\Users\<YourUsername>\go\bin` (or your GOPATH)

## Verify Installation

Check the version and that tf2report is working:

```bash
tf2report --help
```

You should see the help output with available commands and options.

## Updating tf2report

### Go Install

Simply reinstall with the latest tag:

```bash
go install github.com/germainlefebvre4/tf2report/cmd/tf2report@latest
```

### From Source

Pull the latest changes and rebuild:

```bash
cd tf2report
git pull origin main
make build
sudo make install
```

## Uninstalling

### Go Install

Remove the binary from your Go bin directory:

```bash
rm $(go env GOPATH)/bin/tf2report
```

### System Installation

If installed system-wide:

```bash
sudo rm /usr/local/bin/tf2report
```

## Build Options

When building from source, several Make targets are available:

### Standard Build

```bash
make build
```

Creates optimized binary in `bin/tf2report`.

### Development Build

```bash
make dev
```

Builds with race detection and additional checks.

### Clean Build

```bash
make clean
make build
```

Removes old binaries and builds fresh.

### Install

```bash
make install
```

Installs to `/usr/local/bin` (requires sudo).

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

## Troubleshooting

### Command Not Found

**Problem:** `tf2report: command not found`

**Solution:** Ensure `$GOPATH/bin` is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Add this to your shell profile (`.bashrc`, `.zshrc`, etc.) to make it permanent.

### Permission Denied

**Problem:** Permission denied when running `make install`

**Solution:** Use `sudo`:

```bash
sudo make install
```

### Go Version Error

**Problem:** `go: golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64 requires go >= 1.25`

**Solution:** Update Go to version 1.25 or later:

```bash
# Using gvm (Go Version Manager)
gvm install go1.25
gvm use go1.25

# Or download from https://golang.org/dl/
```

### Build Fails

**Problem:** Build errors or missing dependencies

**Solution:** Clean and rebuild:

```bash
make clean
make deps
make build
```

## Next Steps

- **[Getting Started](./getting-started.md)** - Generate your first report
- **[Usage Guide](./usage.md)** - Learn how to use tf2report
- **[Configuration](./configuration.md)** - Set up configuration files
