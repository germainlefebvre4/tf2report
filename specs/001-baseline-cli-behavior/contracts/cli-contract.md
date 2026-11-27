# CLI Contract: tf2report (baseline)

The user-facing interface pinned by the characterization tests. It describes the behavior at the
time of the baseline, including behavior listed as a Known Gap in `spec.md` (those are marked
**provisional**).

## Invocation

```text
tf2report [flags]
```

No positional arguments are used. The tool reads one plan file and writes one report.

## Flags

| Flag | Short | Value | Effect |
|------|-------|-------|--------|
| `--plan` | `-p` | path | Plan JSON file. |
| `--format` | `-f` | `markdown`, `md`, `text`, `txt`, `json` | Report format. |
| `--type` | `-t` | resource type, repeatable or comma separated | Keep only these types. |
| `--action` | `-a` | `create`, `update`, `delete`, `replace`, repeatable (`no-op` also accepted, see below) | Keep only these actions. |
| `--config` | | path | Explicit configuration file. |
| `--verbose` | `-v` | none | Sets verbosity to debug and prints progress on standard error. |
| `--help` | `-h` | none | Usage on standard output, exit 0. |
| `--version` | | none | Prints `tf2report`, `Version:`, `Commit:`, `Build time:` lines. |

Provisional: an `--action` value outside the four listed is accepted and matches nothing, except `no-op`, which is accepted and matches unchanged resources (the markdown, text and json formats already exclude them, so their output is unchanged; the svg format draws them muted).

## Configuration file and environment

Keys: `terraform_plan_path`, `output_format`, `verbosity`, `filters.resource_types`,
`filters.actions`. Sample: `tf2report.yaml.example`.

Environment variables are prefixed `TF2REPORT_` (`TF2REPORT_OUTPUT_FORMAT`, ...). Nested keys work
only with a dotted name (provisional, see data model).

Precedence: flag > environment > file > default. Defaults: plan `terraform.tfplan.json`, format
`markdown`, verbosity `info`, no filters.

File lookup without `--config`: `tf2report` in the working directory, then
`$HOME/.config/tf2report`. Provisional: a file named exactly `tf2report` (no extension) in the
working directory is read as the configuration.

## Streams and exit codes

| Situation | Standard output | Standard error | Exit code |
|-----------|-----------------|----------------|-----------|
| Report produced (including "no changes") | the report | empty (progress lines with `-v`) | 0 |
| `--help` | usage | empty | 0 |
| `--version` | version lines | empty | 0 |
| Config file unreadable or invalid | empty | message | 1 |
| Plan path empty after resolution | empty | message | 1 |
| Plan file missing | empty | `plan file not found: <path>` | 1 |
| Plan not valid JSON | empty | `failed to parse plan: ...` | 1 |
| Unsupported format | empty | `invalid format: <x> (supported: markdown, text, json, svg)` | 1 |

Provisional: on error, the message is printed twice (once prefixed `Error:` followed by the usage
text, once alone) on standard error.

Verbose progress lines on standard error: `Parsed plan with N resource changes`, then
`Filtered by type: N resources remain` and `Filtered by action: N resources remain` when the
matching filter is set.

## Report layouts

### Markdown

```text
# Terraform Plan Summary
**Terraform Version:** <v>            (only when known)
## Summary
**Total Changes:** <n>               (or "No changes. Infrastructure is up-to-date.")
| Action | Count |                   (rows only for non-zero: Add, Change, Destroy, Replace)
## Resources by Change Type
### To Add (n) | To Change (n) | To Replace (n) | To Destroy (n)    (non-empty sections, this order)
| Resource Type | Resource Name |
## Changes by Resource Type
### <type> (n)
| Resource Type | Resource Name | Change Type |
```

### Text

Same sections with upper-case titles, 50-character `=` and `-` rules, aligned columns and no
Markdown syntax. "No changes" wording identical to Markdown.

### JSON

```text
{
  "terraform_version": string,     // omitted when unknown
  "summary": { total_changes, to_add, to_change, to_destroy, to_replace },
  "changes": [ { address, type, name, action } ],   // plan order, no-op/read excluded
  "by_type": { "<type>": { type, to_add, to_change, to_destroy, to_replace, resources: [address] } }
}
```

With no changes: counters at 0, `changes` is `[]` and `by_type` is `{}`. Indented with two spaces,
terminated by a newline.

## Guarantees

- Attribute values (including sensitive ones) never appear in any format.
- No network access and no Terraform execution.
- Same input and options give byte-identical output.
