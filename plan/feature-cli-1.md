---
goal: CLI Implementation for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, cli]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan covers the Cobra-based CLI implementation, flag handling, and execution flow.

## 1. Requirements & Constraints

- **REQ-001**: Provide CLI flags for plan path, output format, filters, and verbosity.
- **REQ-002**: Support config file override and environment variable defaults.
- **CON-001**: Print errors to stderr and return non-zero exit codes.
- **GUD-001**: Provide concise help and version information.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement CLI entry point and command execution.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Define root command with descriptions and version template | Yes | 2026-01-06 |
| TASK-002 | Add flags: `--config`, `--plan`, `--format`, `--type`, `--action`, `--verbose` | Yes | 2026-01-06 |
| TASK-003 | Load configuration and override with CLI flags | Yes | 2026-01-06 |
| TASK-004 | Validate plan path and file existence | Yes | 2026-01-06 |
| TASK-005 | Parse plan, filter changes, compute summary, and format output | Yes | 2026-01-06 |

## 3. Alternatives

- **ALT-001**: Use a custom flag parser. Rejected to leverage Cobra ecosystem.

## 4. Dependencies

- **DEP-001**: `github.com/spf13/cobra` for CLI framework.
- **DEP-002**: `pkg/config`, `pkg/terraform`, `pkg/report` packages.

## 5. Files

- **FILE-001**: [cmd/tf2report/main.go](cmd/tf2report/main.go) - CLI implementation.

## 6. Testing

- **TEST-001**: Run CLI with sample plan and default format.
- **TEST-002**: Run CLI with `--format json` and filters.

## 7. Risks & Assumptions

- **RISK-001**: Invalid flag combinations could lead to confusing errors.
- **ASSUMPTION-001**: Users provide valid plan paths.

## 8. Related Specifications / Further Reading

- [README.md](../README.md)