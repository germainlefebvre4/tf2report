---
goal: Plain Text Report Formatting for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, report, text]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan describes the plain text formatter for terminal-friendly output.

## 1. Requirements & Constraints

- **REQ-001**: Produce plain text output with clear sections and separators.
- **REQ-002**: Include summary counts and grouped resource lists.
- **CON-001**: Use consistent indentation for readability.
- **GUD-001**: Keep output deterministic and readable in terminals.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement plain text formatter.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Add Text formatter type and registration in factory | Yes | 2026-01-05 |
| TASK-002 | Implement `Format()` flow matching Markdown formatter | Yes | 2026-01-05 |
| TASK-003 | Implement `writeHeader()` and `writeSummary()` with separators | Yes | 2026-01-05 |
| TASK-004 | Implement `writeResourcesByChangeType()` and `writeDetailsByType()` | Yes | 2026-01-05 |
| TASK-005 | Sort resource types alphabetically for consistent output | Yes | 2026-01-05 |

## 3. Alternatives

- **ALT-001**: Use a table formatting library. Rejected to keep dependencies minimal.

## 4. Dependencies

- **DEP-001**: `pkg/report/formatter.go` for formatter interface.
- **DEP-002**: `pkg/terraform/types.go` for summary structures.

## 5. Files

- **FILE-001**: [pkg/report/text.go](pkg/report/text.go) - Text formatter.

## 6. Testing

- **TEST-001**: Verify output for empty change set.
- **TEST-002**: Verify output for mixed change types.

## 7. Risks & Assumptions

- **RISK-001**: Wide output lines may wrap on narrow terminals.
- **ASSUMPTION-001**: Consumers prefer concise text output.

## 8. Related Specifications / Further Reading

- [docs/OUTPUT_FORMATS.md](../docs/OUTPUT_FORMATS.md)