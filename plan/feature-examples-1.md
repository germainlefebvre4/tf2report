---
goal: Example Files and Templates for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, examples]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan documents example Terraform plan files and configuration templates.

## 1. Requirements & Constraints

- **REQ-001**: Provide sample plan JSON files for testing.
- **REQ-002**: Provide example configuration template.
- **CON-001**: Store examples under `examples/`.
- **GUD-001**: Avoid embedding output examples in documentation.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Add sample plan files and configuration templates.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Add `examples/` directory with sample plan JSON files | Yes | 2026-01-13 |
| TASK-002 | Add root-level sample plan file for quick testing | Yes | 2026-01-13 |
| TASK-003 | Add `tf2report.yaml.example` with full configuration options | Yes | 2026-01-13 |

## 3. Alternatives

- **ALT-001**: Only provide minimal sample plan. Rejected to support varied scenarios.

## 4. Dependencies

- **DEP-001**: None beyond sample data files.

## 5. Files

- **FILE-001**: [examples/](examples/) - Sample plan JSON files.
- **FILE-002**: [terraform.tfplan.json](terraform.tfplan.json) - Root sample plan file.
- **FILE-003**: [tf2report.yaml.example](tf2report.yaml.example) - Example config.

## 6. Testing

- **TEST-001**: Run CLI against example plan files.

## 7. Risks & Assumptions

- **RISK-001**: Example plans may become outdated relative to Terraform versions.
- **ASSUMPTION-001**: Users can adapt examples to their environments.

## 8. Related Specifications / Further Reading

- [docs/EXAMPLES.md](../docs/EXAMPLES.md)