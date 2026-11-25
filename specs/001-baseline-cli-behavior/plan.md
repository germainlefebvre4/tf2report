# Implementation Plan: Baseline CLI Behavior

**Branch**: `chore/speckit-brownfield` (spec directory `001-baseline-cli-behavior`) | **Date**: 2026-11-25 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-baseline-cli-behavior/spec.md`

## Summary

Pin the current behavior of tf2report with automated **characterization tests**, without changing
application behavior (spec FR-016, SC-007). The approach is test-only:

- **Unit tests** next to each package (`pkg/terraform`, `pkg/config`, `pkg/report`), table-driven,
  with golden files for the three report formats.
- **Black-box CLI tests** in `cmd/tf2report` that build the real binary once and run it in an
  isolated working directory. This is the only way to cover flags, precedence, streams and exit
  codes, because `main.go` keeps its wiring in `init()` and package-level variables.
- **Synthetic fixtures** under `testdata/` (no real infrastructure data), taken from supported
  Terraform versions.
- **Known Gaps** (spec section) are pinned by tests named `TestKnownGap_*`, each marked
  provisional, so fixing a gap is a deliberate, visible test change.

Standard library only: no new dependency. Spec Kit artifacts (research, data model, contract,
quickstart) document what is being pinned.

## Technical Context

**Language/Version**: Go 1.25 (`go.mod`; CI uses 1.25; the local toolchain is newer, which is fine)

**Primary Dependencies**: none added. Existing: `cobra`, `viper`. Tests use `testing`, `os/exec`,
`encoding/json`, `os`, `path/filepath`, `bytes` only.

**Storage**: files only. Plan JSON in, report on standard output. Golden files are text.

**Testing**: `go test ./...` (`make test`, `make test-coverage`). Table-driven tests with `t.Run`,
`t.Helper()`, `t.Cleanup()`, `t.TempDir()`, `t.Setenv()`. Golden files refreshed with an
`-update` flag, never by default.

**Target Platform**: Linux and macOS (the shipped platforms); tests run on Linux in CI.

**Project Type**: single-binary CLI (`cmd/` + `pkg/`).

**Performance Goals**: none claimed (spec Assumptions). The full suite SHOULD finish in a few
seconds; the one-time binary build dominates.

**Constraints**: tests MUST NOT change application behavior or exported API; offline; no real
plan data; no network (Constitution, Technology constraints).

**Scale/Scope**: about 8 fixture plans, 3 formats, 5 user stories, 6 pinned Known Gaps (plus one
documentation-only gap). About 23 new `_test.go` files and a `testdata/` tree.

All Technical Context items are resolved; the one open point (which Terraform versions count as
"supported") is handled in `research.md` R-004 and does not block the design.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design (below).*

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. Idiomatic Go and Simplicity | Pass | Tests are plain stdlib Go, formatted with `gofmt`/`goimports`, no `init()` and no globals added (the `-update` flag lives in a test-only `flag.Bool`). No new dependency. |
| II. Layered Package Architecture | Pass | No production file changes. Tests respect the layering: `pkg/report` tests build `terraform.Plan` values or parse fixtures; no new cross-package dependency in production code. |
| III. Test-Backed Changes | Pass (this is the feature) | Table-driven, `t.Run`, success and error paths, malformed and empty plans, golden-pinned output per format, `t.Helper()`/`t.Cleanup()`. |
| IV. Automated Quality Gates | Deviation, tracked | `golangci-lint` 2.15 cannot load the repository's v1-format `.golangci.yml` ("unsupported version of the configuration"), so the gate cannot run. Interim gate: `gofmt -l`, `goimports`, `go vet ./...`, pre-commit hooks. See Complexity Tracking. |
| V. Conventional Commits and Automated Releases | Pass | Commits use `test:` / `docs:` / `chore:`; `CHANGELOG.md` untouched. `test:` commits do not trigger a release. |

Technology and Delivery Constraints: read-only and offline (tests run only local files); fixtures
are synthetic; no change to the Makefile contract beyond optional documentation.

**Gate result: PASS** (one tracked deviation, no unjustified violation).

## Project Structure

### Documentation (this feature)

```text
specs/001-baseline-cli-behavior/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-contract.md  # Phase 1 output: the pinned user-facing interface
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 output (/speckit-tasks command, not created here)
```

### Source Code (repository root)

```text
cmd/tf2report/
├── main.go                       # unchanged
├── main_test.go                  # NEW: TestMain builds the binary; runCLI helper
├── summary_test.go, format_test.go, filters_test.go,
│   config_test.go, errors_test.go,
│   diagnostics_test.go           # NEW: black-box CLI tests, one file per theme

pkg/terraform/
├── parser.go, types.go           # unchanged
├── helpers_test.go               # NEW: fixture loader
├── types_test.go                 # NEW: GetPrimaryAction table
├── summary_test.go               # NEW: NewSummary counts and invariants
├── parser_test.go                # NEW: Parse/ParseFile success cases
├── parser_errors_test.go         # NEW: malformed and empty input
└── filter_test.go                # NEW: FilterByType / FilterByAction

pkg/config/
├── loader.go, types.go           # unchanged
├── helpers_test.go               # NEW: isolate (cwd, HOME, env) and writeFile
├── loader_test.go                # NEW: defaults, explicit file, search paths, env, precedence, errors
└── types_test.go                 # NEW: NewConfig defaults, GetActionTypes

pkg/report/
├── formatter.go, json.go,
│   markdown.go, text.go          # unchanged
├── golden_test.go                # NEW: -update flag, loadPlan/render/assertGolden helpers
├── formatter_test.go             # NEW: NewFormatter, ParseFormat aliases and errors
├── markdown_test.go, text_test.go,
│   json_test.go                  # NEW: golden tests per format, determinism, no-changes
├── safety_test.go                # NEW: no attribute values in any format
├── gaps_test.go                  # NEW: provisional Known Gap tests for report output
└── testdata/                     # NEW: *.golden files (markdown, text, json per fixture)

testdata/
├── README.md                     # NEW: fixture provenance, Terraform versions, "synthetic only"
└── plans/                        # NEW: shared synthetic plans
    ├── basic.json                #   1 add, 1 change, 1 destroy, 1 replace
    ├── modules-and-counts.json   #   module paths, count index, for_each keys (Known Gap)
    ├── noop-and-read.json        #   no-op and data-source read mixed in
    ├── sensitive.json            #   sensitive attributes with synthetic marker values
    ├── no-changes.json           #   resource_changes: []
    ├── unknown-action.json       #   unrecognized action string
    ├── empty-object.json         #   {}
    └── invalid.json              #   not JSON
```

**Structure Decision**: keep tests beside the code they pin (Go convention, Constitution
Principle III, `go.instructions.md`). Plan fixtures are shared at the repository root `testdata/`
because the CLI, `pkg/terraform` and `pkg/report` tests need the same plans; `go` ignores
directories named `testdata` when building. Golden files live with the package that owns the
format (`pkg/report/testdata`). No production directory or package is added.

### Traceability: stories to tests

| Story | Primary tests |
|-------|---------------|
| 1 Summarize a plan | `pkg/terraform` (`GetPrimaryAction`, `NewSummary`), `pkg/report` goldens, CLI default run |
| 2 Choose the format | `pkg/report` `ParseFormat`/`NewFormatter`, three goldens, CLI `-f` incl. `md`/`txt`/`xml` |
| 3 Filters | `pkg/terraform` `FilterBy*`, CLI `-t`/`-a`, combined filters, zero-match |
| 4 Configuration | `pkg/config` loader tests, CLI flag-over-config and env cases |
| 5 Failures | `pkg/terraform` malformed input, CLI exit codes, stdout empty, stderr content, `-v`, `--version` |
| Known Gaps (FR-016) | `TestKnownGap_*` in CLI and `pkg/report` tests |

Run `go test -run KnownGap ./...` to list every provisional test.

## Phase 0 and Phase 1 outputs

- [research.md](research.md): decisions on test strategy, black-box CLI harness, golden files,
  fixtures, provisional-test convention, supported-versions policy, lint gate.
- [data-model.md](data-model.md): entities, action derivation rules, summary rules.
- [contracts/cli-contract.md](contracts/cli-contract.md): flags, configuration, streams, exit
  codes, and the three report layouts being pinned.
- [quickstart.md](quickstart.md): how to run, refresh goldens and validate SC-001 to SC-007.

## Constitution Check (post-design)

Re-evaluated after Phase 1: unchanged. The design adds only `_test.go` files, `testdata/` and
spec documents; no production code, no dependency, no new package. Principle IV remains a tracked
deviation caused by the pre-existing lint configuration, not by this feature.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Principle IV: `golangci-lint` gate not runnable (pre-existing config in v1 format, tool at 2.15) | The feature must still merge with a quality bar, and migrating the lint config is a separate change (Sync Impact Report TODO `LINT_CI`) | Migrating `.golangci.yml` here would mix a tooling change into a test-only feature and could force production edits to satisfy new linters, contradicting FR-016 "no application change". Interim gate (`gofmt`, `goimports`, `go vet`, pre-commit) is used instead. |
| Black-box CLI tests build and exec a binary | `cmd/tf2report/main.go` keeps flags in `init()` and package globals, so in-process tests would leak state between cases and could not observe exit codes | Refactoring `main.go` to be testable is a production change, out of scope here (tracked as TODO `INIT`/`CMD_LOGIC`). |
