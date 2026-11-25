# Data Model: Baseline CLI Behavior

Entities as they exist today (source: `pkg/terraform/types.go`, `pkg/config/types.go`,
`pkg/report/*.go`). This feature documents and pins them; it changes none of them.

## Plan

| Field | Source key | Notes |
|-------|------------|-------|
| Format version | `format_version` | Read, never validated. |
| Terraform version | `terraform_version` | Optional; omitted from reports when empty. |
| Resource changes | `resource_changes` | List; absent or empty means "no changes". |
| Output changes | `output_changes` | Parsed, never reported. |

Unknown keys are ignored. An empty object `{}` is a valid plan with no changes. Anything that is
not valid JSON is an error.

## Resource Change

| Field | Source key | Used in reports |
|-------|------------|-----------------|
| Address | `address` | JSON output only (`address`), and for sorting. |
| Mode | `mode` | No. |
| Type | `type` | Yes (all formats), filter key. |
| Name | `name` | Yes (all formats). |
| Provider | `provider_name` | No. |
| Change | `change` | Action only. |
| Action reason, previous address | `action_reason`, `previous_address` | No. |

Identity: the address is unique in a plan, but Markdown and text reports show only type and name,
so two resources can look identical (Known Gap "module path and index not shown").

## Change

| Field | Notes |
|-------|-------|
| Actions | Ordered list of strings from the plan. |
| Before / After | Attribute maps. **Never written to any report** (FR-013, SC-006). |
| After unknown, before/after sensitive | Parsed, never reported. |

### Primary action derivation (pinned by `GetPrimaryAction` tests)

| `actions` | Primary action |
|-----------|----------------|
| `[]` (empty) | no-op |
| `["create"]` | create |
| `["update"]` | update |
| `["delete"]` | delete |
| `["read"]` | read |
| `["no-op"]` | no-op |
| `["delete","create"]` or `["create","delete"]` | replace |
| any other single or unknown value | no-op |
| two actions that are not delete and create | first action decides |

Read and no-op are excluded from counts and listings (FR-003).

## Summary

| Field | Rule |
|-------|------|
| Total changes | Number of changes whose primary action is not read or no-op. |
| To add / change / destroy / replace | Count of create / update / delete / replace. |
| By type | Per resource type: the same four counters; only types with at least one counted change appear. |

Invariant: `total = add + change + destroy + replace`, and the sum over `by type` equals the same
totals (SC-002).

## Configuration

| Key (file) | Env variable | Flag | Default |
|------------|--------------|------|---------|
| `terraform_plan_path` | `TF2REPORT_TERRAFORM_PLAN_PATH` | `--plan`, `-p` | `terraform.tfplan.json` |
| `output_format` | `TF2REPORT_OUTPUT_FORMAT` | `--format`, `-f` | `markdown` |
| `verbosity` | `TF2REPORT_VERBOSITY` | `--verbose`, `-v` (sets debug) | `info` |
| `filters.resource_types` | `TF2REPORT_FILTERS.RESOURCE_TYPES` (dotted name) | `--type`, `-t` | empty |
| `filters.actions` | `TF2REPORT_FILTERS.ACTIONS` (dotted name) | `--action`, `-a` | empty |

Precedence, highest first: flag, environment variable, configuration file, built-in default.
File search when no `--config`: `./tf2report.yaml` (any extension or none, see Known Gap), then
`$HOME/.config/tf2report`. A missing default file is not an error; an explicit missing file is.

Nested filter keys are settable from the environment only through a variable name containing a
dot (`TF2REPORT_FILTERS.ACTIONS`), which most shells cannot export; `TF2REPORT_FILTERS_ACTIONS`
is ignored (verified by running the binary). Tests pin this and mark it as provisional.

## Report

| Format | Names accepted | Sections |
|--------|----------------|----------|
| Markdown | `markdown`, `md` | title, Terraform version, summary table, resources by change type, changes by resource type |
| Text | `text`, `txt` | same sections, aligned columns, 50-character rules |
| JSON | `json` | `terraform_version` (omitted when empty), `summary`, `changes`, `by_type` |

Ordering: Markdown and text sort types alphabetically and resources by name or address inside
each section; JSON lists `changes` in plan order and `by_type` keys in sorted order (encoder
behavior). Identical input and options give identical bytes (FR-014).

## State transitions

None. The tool is stateless: one run reads inputs and writes one report.
