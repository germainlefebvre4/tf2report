---
goal: Markdown Report Formatting for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, report, markdown]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan describes the Markdown formatter implementation for human-readable reports.

## 1. Requirements & Constraints

- **REQ-001**: Generate Markdown output with headers, tables, and lists.
- **REQ-002**: Include summary counts and details grouped by change type and resource type.
- **CON-001**: Sort resource types for stable output.
- **GUD-001**: Keep output deterministic and readable.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement Markdown formatter.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Add Markdown formatter type and registration in formatter factory | Yes | 2026-01-05 |
| TASK-002 | Implement `Format()` flow: header, summary, change groups, type details | Yes | 2026-01-05 |
| TASK-003 | Implement `writeHeader()` and `writeSummary()` methods | Yes | 2026-01-05 |
| TASK-004 | Implement `writeResourcesByChangeType()` and `writeDetailsByType()` | Yes | 2026-01-05 |
| TASK-005 | Sort resource types alphabetically for consistent output | Yes | 2026-01-05 |

## 3. Alternatives

- **ALT-001**: Use a template engine. Rejected to keep dependencies minimal.

## 4. Dependencies

- **DEP-001**: `pkg/report/formatter.go` for formatter interface.
- **DEP-002**: `pkg/terraform/types.go` for summary structures.

## 5. Files

- **FILE-001**: [pkg/report/markdown.go](pkg/report/markdown.go) - Markdown formatter.

## 6. Testing

- **TEST-001**: Verify Markdown output for empty changes.
- **TEST-002**: Verify Markdown output for mixed actions and types.

## 7. Risks & Assumptions

- **RISK-001**: Large change sets could produce verbose output.
- **ASSUMPTION-001**: Consumers can parse Markdown tables.

## 8. Related Specifications / Further Reading

- [docs/OUTPUT_FORMATS.md](../docs/OUTPUT_FORMATS.md)