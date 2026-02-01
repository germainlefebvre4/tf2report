---
goal: Resource Filtering and Summarization for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, filtering, summary]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan documents filtering by resource type/action and the summarization of Terraform plan changes.

## 1. Requirements & Constraints

- **REQ-001**: Filter resource changes by type and action.
- **REQ-002**: Provide summary counts for create, update, delete, replace.
- **CON-001**: Skip `no-op` and `read` actions in summaries.
- **GUD-001**: Keep filtering logic deterministic and order-preserving.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement filters and summary aggregation.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Implement `FilterByType()` for resource type filtering | Yes | 2026-01-03 |
| TASK-002 | Implement `FilterByAction()` for action filtering | Yes | 2026-01-03 |
| TASK-003 | Define `Summary` and `TypeSummary` structs | Yes | 2026-01-03 |
| TASK-004 | Implement `NewSummary()` aggregation with per-type counts | Yes | 2026-01-03 |
| TASK-005 | Skip `no-op` and `read` actions in counts | Yes | 2026-01-03 |

## 3. Alternatives

- **ALT-001**: Pre-filter during parsing. Rejected to keep parser responsibilities limited.
- **ALT-002**: Use map-based aggregation only. Rejected to preserve stable list output ordering.

## 4. Dependencies

- **DEP-001**: `pkg/terraform/types.go` for action types.

## 5. Files

- **FILE-001**: [pkg/terraform/parser.go](pkg/terraform/parser.go) - Filter functions.
- **FILE-002**: [pkg/terraform/types.go](pkg/terraform/types.go) - Summary types and aggregation.

## 6. Testing

- **TEST-001**: Filter by type with single and multiple values.
- **TEST-002**: Filter by action with create/update/delete/replace.
- **TEST-003**: Summaries exclude `no-op` and `read` actions.

## 7. Risks & Assumptions

- **RISK-001**: Mixed action arrays could be misclassified.
- **ASSUMPTION-001**: Action values match Terraform JSON spec.

## 8. Related Specifications / Further Reading

- [pkg/terraform/types.go](../pkg/terraform/types.go)