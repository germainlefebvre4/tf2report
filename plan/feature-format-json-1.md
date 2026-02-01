---
goal: JSON Report Formatting for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, report, json]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan documents the JSON formatter for machine-readable output and CI/CD integration.

## 1. Requirements & Constraints

- **REQ-001**: Output structured JSON with summary, changes, and per-type breakdown.
- **REQ-002**: Skip `no-op` and `read` actions in output.
- **CON-001**: Use standard library JSON encoder with indentation.
- **GUD-001**: Keep JSON schema stable for automation.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement JSON formatter and schema.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Define JSON report structures for summary and type breakdowns | Yes | 2026-01-05 |
| TASK-002 | Implement `Format()` using JSON encoder with indentation | Yes | 2026-01-05 |
| TASK-003 | Build change list and by-type resources | Yes | 2026-01-05 |
| TASK-004 | Register JSON formatter in factory and parse format strings | Yes | 2026-01-05 |

## 3. Alternatives

- **ALT-001**: Use YAML output instead. Rejected to keep machine parsing simple.

## 4. Dependencies

- **DEP-001**: Go standard library `encoding/json`.
- **DEP-002**: `pkg/terraform/types.go` for summaries and action types.

## 5. Files

- **FILE-001**: [pkg/report/json.go](pkg/report/json.go) - JSON formatter.

## 6. Testing

- **TEST-001**: Verify JSON output schema and indentation.
- **TEST-002**: Ensure `no-op` and `read` actions are excluded.

## 7. Risks & Assumptions

- **RISK-001**: Schema changes could break automation consumers.
- **ASSUMPTION-001**: Consumers can parse JSON reliably.

## 8. Related Specifications / Further Reading

- [docs/OUTPUT_FORMATS.md](../docs/OUTPUT_FORMATS.md)