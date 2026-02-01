---
goal: Project Foundation and Core Architecture for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, foundation, architecture, build]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan defines the foundational project structure, core architecture layout, and baseline build configuration for the tf2report CLI.

## 1. Requirements & Constraints

- **REQ-001**: Initialize Go module with correct module path and dependency management.
- **REQ-002**: Establish standard project layout with `cmd/` and `pkg/` directories.
- **CON-001**: Use Go 1.21+ with module support.
- **GUD-001**: Follow idiomatic Go project layout conventions.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Establish project structure, dependencies, and build system.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Initialize Go module with module path `github.com/germainlefebvre4/tf2report` | Yes | 2025-12-31 |
| TASK-002 | Create package structure: `cmd/tf2report`, `pkg/config`, `pkg/terraform`, `pkg/report` | Yes | 2025-12-31 |
| TASK-003 | Set up Makefile with build, test, install, clean, and run targets | Yes | 2025-12-31 |
| TASK-004 | Configure build with LDFLAGS for version, commit, and build time injection | Yes | 2025-12-31 |
| TASK-005 | Add dependencies: `github.com/spf13/cobra` and `github.com/spf13/viper` | Yes | 2025-12-31 |
| TASK-006 | Create LICENSE file (appropriate open-source license) | Yes | 2025-12-31 |
| TASK-007 | Initialize Git repository with `.gitignore` for Go projects | Yes | 2025-12-31 |

## 3. Alternatives

- **ALT-001**: Use a monolithic `main.go` without packages. Rejected to preserve modularity.
- **ALT-002**: Use a non-standard directory layout. Rejected to align with Go conventions.

## 4. Dependencies

- **DEP-001**: Go 1.21+ toolchain for module management.
- **DEP-002**: Cobra and Viper for CLI and configuration scaffolding.

## 5. Files

- **FILE-001**: [go.mod](go.mod) - Module definition.
- **FILE-002**: [Makefile](Makefile) - Build automation.
- **FILE-003**: [cmd/tf2report/main.go](cmd/tf2report/main.go) - CLI entry point.
- **FILE-004**: [pkg/](pkg/) - Core packages for config, report, terraform.

## 6. Testing

- **TEST-001**: Validate `go mod tidy` succeeds.
- **TEST-002**: Validate `make build` produces executable binary.

## 7. Risks & Assumptions

- **RISK-001**: Incorrect module path could break imports.
- **ASSUMPTION-001**: Git is available for version metadata during builds.

## 8. Related Specifications / Further Reading

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Makefile](../Makefile)