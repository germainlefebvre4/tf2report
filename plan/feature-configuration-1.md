---
goal: Configuration Management for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, configuration]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan documents configuration loading, defaults, and environment variable support.

## 1. Requirements & Constraints

- **REQ-001**: Support YAML configuration file loading.
- **REQ-002**: Provide defaults for plan path, output format, and verbosity.
- **REQ-003**: Support environment variable overrides with `TF2REPORT_` prefix.
- **CON-001**: Configuration file is optional.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement configuration types and loader.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Define `Config` and `Filters` structs | Yes | 2026-01-04 |
| TASK-002 | Implement `NewConfig()` with defaults | Yes | 2026-01-04 |
| TASK-003 | Implement `Filters.GetActionTypes()` conversion | Yes | 2026-01-04 |
| TASK-004 | Implement `Loader` with Viper integration | Yes | 2026-01-04 |
| TASK-005 | Load config from explicit path or default locations | Yes | 2026-01-04 |
| TASK-006 | Configure environment variable prefix and defaults | Yes | 2026-01-04 |
| TASK-007 | Provide example configuration file | Yes | 2026-01-04 |

## 3. Alternatives

- **ALT-001**: Use only CLI flags with no config file. Rejected to support repeatable usage.
- **ALT-002**: JSON configuration. Rejected in favor of YAML readability.

## 4. Dependencies

- **DEP-001**: `github.com/spf13/viper` for configuration loading.

## 5. Files

- **FILE-001**: [pkg/config/types.go](pkg/config/types.go) - Configuration types.
- **FILE-002**: [pkg/config/loader.go](pkg/config/loader.go) - Loader implementation.
- **FILE-003**: [tf2report.yaml.example](tf2report.yaml.example) - Example config.

## 6. Testing

- **TEST-001**: Load default config when no file is present.
- **TEST-002**: Load config from provided file path.
- **TEST-003**: Override config via environment variables.

## 7. Risks & Assumptions

- **RISK-001**: Misconfigured environment variables override intended defaults.
- **ASSUMPTION-001**: Config files are valid YAML.

## 8. Related Specifications / Further Reading

- [docs/CONFIGURATION.md](../docs/CONFIGURATION.md)