---
goal: Release Automation for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, release, automation]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan documents GoReleaser configuration and release workflow automation.

## 1. Requirements & Constraints

- **REQ-001**: Build release artifacts for Linux and Darwin.
- **REQ-002**: Generate checksums and changelog.
- **CON-001**: Use GoReleaser configuration in repository root.
- **GUD-001**: Inject version metadata via LDFLAGS.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement GoReleaser configuration and GitHub Actions release workflow.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Create `.goreleaser.yaml` with build, archive, and checksum configuration | Yes | 2026-01-08 |
| TASK-002 | Configure pre-release hooks and changelog filters | Yes | 2026-01-08 |
| TASK-003 | Add GitHub Actions release workflow to run GoReleaser | Yes | 2026-01-08 |

## 3. Alternatives

- **ALT-001**: Use custom shell scripts for releases. Rejected to leverage GoReleaser.

## 4. Dependencies

- **DEP-001**: GoReleaser.
- **DEP-002**: GitHub Actions.

## 5. Files

- **FILE-001**: [.goreleaser.yaml](.goreleaser.yaml) - GoReleaser config.
- **FILE-002**: [.github/workflows/release.yml](.github/workflows/release.yml) - Release workflow.

## 6. Testing

- **TEST-001**: Validate GoReleaser config via dry run.

## 7. Risks & Assumptions

- **RISK-001**: Missing release permissions could fail publishing.
- **ASSUMPTION-001**: Git tags follow `v*` convention.

## 8. Related Specifications / Further Reading

- [GoReleaser](https://goreleaser.com/)