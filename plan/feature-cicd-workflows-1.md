---
goal: CI/CD Workflows for tf2report
version: 1.0
date_created: 2026-02-01
last_updated: 2026-02-01
owner: Germain Lefebvre
status: 'Completed'
tags: [feature, cicd, workflows]
---

# Introduction

![Status: Completed](https://img.shields.io/badge/status-Completed-brightgreen)

This plan covers CI/CD workflows for pull requests, Docker publishing, and release automation.

## 1. Requirements & Constraints

- **REQ-001**: Validate builds and tests on pull requests.
- **REQ-002**: Publish Docker images on release or tag events.
- **REQ-003**: Automate versioning with release-please.
- **CON-001**: Use GitHub Actions workflows.

## 2. Implementation Steps

### Implementation Phase 1

- GOAL-001: Implement CI/CD workflows for validation and publishing.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Create pull request workflow to run build, test, and lint | Yes | 2026-01-08 |
| TASK-002 | Create Docker publish workflow with multi-platform builds | Yes | 2026-01-08 |
| TASK-003 | Create release-please workflow and config | Yes | 2026-01-08 |

## 3. Alternatives

- **ALT-001**: Use external CI provider. Rejected to keep GitHub-native automation.

## 4. Dependencies

- **DEP-001**: GitHub Actions.
- **DEP-002**: Docker Buildx.
- **DEP-003**: release-please.

## 5. Files

- **FILE-001**: [.github/workflows/pull-requests.yml](.github/workflows/pull-requests.yml)
- **FILE-002**: [.github/workflows/docker-publish.yml](.github/workflows/docker-publish.yml)
- **FILE-003**: [.github/workflows/release-please.yml](.github/workflows/release-please.yml)
- **FILE-004**: [release-please-config.json](release-please-config.json)

## 6. Testing

- **TEST-001**: Validate workflow syntax and triggers.

## 7. Risks & Assumptions

- **RISK-001**: Secrets or permissions misconfiguration could block publishing.
- **ASSUMPTION-001**: Repository has required GitHub Actions permissions enabled.

## 8. Related Specifications / Further Reading

- [GitHub Actions](https://docs.github.com/en/actions)