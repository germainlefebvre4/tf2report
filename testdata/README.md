# Test fixtures

Fixtures are synthetic only. They contain no real infrastructure data. Values use obvious
markers such as `SYNTHETIC-...`.

Golden files (in `pkg/report/testdata/`) are refreshed only with `go test ./pkg/report -update`,
and every diff must be reviewed against `specs/001-baseline-cli-behavior/contracts/cli-contract.md`.

## Terraform version

All fixtures use `terraform_version` `1.17.2`, the newest stable release listed on
https://releases.hashicorp.com/terraform/ (checked 2026-11-25).

Supported series could not be confirmed from HashiCorp's official source: the support-policy
page (`developer.hashicorp.com/terraform/language/docs/support-period-and-end-of-life-policy`)
returned HTTP 404 on 2026-11-25. The releases index lists the series 1.17, 1.16, 1.15, 1.14,
1.13, 1.12, 1.11, 1.10 and 1.9 (latest patches 1.17.2, 1.16.7, 1.15.9, 1.14.9, 1.13.5, 1.12.2,
1.11.4, 1.10.5, 1.9.8) but does not state which are supported. The choice of 1.17.2 was
confirmed by the maintainer.

## Plans

| Fixture | terraform_version | format_version | purpose |
|---------|-------------------|----------------|---------|
| `plans/basic.json` | 1.17.2 | 1.2 | 4 changes: 1 add, 1 change, 1 destroy, 1 replace |
| `plans/modules-and-counts.json` | 1.17.2 | 1.2 | module paths and indexes (Known Gap: not shown in reports) |
| `plans/noop-and-read.json` | 1.17.2 | 1.2 | 1 counted change plus no-op, read and empty-actions entries |
| `plans/sensitive.json` | 1.17.2 | 1.2 | marker strings that must never appear in any report |
| `plans/no-changes.json` | 1.17.2 | 1.2 | valid plan with empty `resource_changes` |
| `plans/unknown-action.json` | 1.17.2 | 1.2 | one valid change and one unknown action |
| `plans/real-shape.json` | 1.17.2 | 1.2 | exact Terraform export shape (boolean sensitivity flags, scalar output); Known Gap: rejected by the parser |
| `plans/empty-object.json` | none | none | `{}`: valid plan without changes or version |
| `plans/invalid.json` | none | none | not JSON, must fail to parse |
