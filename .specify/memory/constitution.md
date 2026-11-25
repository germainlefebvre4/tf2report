# tf2report Constitution

## Core Principles

### I. Idiomatic Go and Simplicity

Code MUST be simple, clear and idiomatic Go, following Effective Go, Go Code Review Comments
and `.github/instructions/go.instructions.md`.

- Code MUST be formatted with `gofmt` and `goimports` (local prefix
  `github.com/germainlefebvre4/tf2report`).
- Errors MUST be checked immediately and wrapped with context using `fmt.Errorf` and `%w`.
  An error MUST be either logged or returned, never both.
- The standard library MUST be preferred over custom code or new dependencies when it
  covers the need. Direct dependencies stay minimal (today: `cobra`, `viper`).
- `init()` functions and unnecessary globals MUST NOT be introduced (`gochecknoinits`).
- Exported symbols MUST be documented, starting with the symbol name. Comments are written
  in English and explain why, not what.

Rationale: a small CLI maintained by few people stays reliable only if it stays boring.

### II. Layered Package Architecture

The current separation of concerns MUST be preserved:

- `cmd/tf2report` holds the CLI entry point only (cobra wiring, flags, version injection).
  It MUST NOT contain parsing or formatting logic.
- `pkg/terraform` parses Terraform plan JSON and computes the change summary.
- `pkg/report` renders output. Every output format (markdown, text, json) MUST be a
  `Formatter` implementation selected through `NewFormatter`; adding a format MUST NOT
  require changing the parser.
- `pkg/config` loads configuration (defaults, YAML file, `TF2REPORT_*` environment
  variables, via viper); command-line flags override it in `cmd`.

Dependencies flow `cmd` -> `pkg/*`, `pkg/report` -> `pkg/terraform` and `pkg/config` ->
`pkg/terraform`; circular dependencies are forbidden. Interfaces are small and defined where they are consumed.
The input contract is the output of `terraform show -json`; changes to how it is read MUST
keep existing valid plans working.

Rationale: a stable parse -> summarize -> format pipeline lets formats and filters evolve
independently.

### III. Test-Backed Changes

Every new or modified behavior MUST ship with tests in the same change.

- Tests are table-driven, use `t.Run` subtests and live next to the code (`_test.go`).
- Tests cover success and error paths, including malformed or empty plan files.
- Each output format MUST have tests that pin its output for representative plans.
- Test helpers call `t.Helper()` and clean up with `t.Cleanup()`.
- `go test ./...` MUST pass before merge (`make test`; `make test-coverage` for coverage).

Current behavior is pinned by the baseline characterization tests. A behavior listed as a Known
Gap in a spec is pinned by a `TestKnownGap_<Slug>` test; fixing it MUST change that test in the
same commit. Code touched by a change MUST gain tests.

Rationale: report output is a contract consumed by CI pipelines and humans; silent drift
is the main risk.

### IV. Automated Quality Gates

Quality is enforced by tools, not by memory.

- `golangci-lint` with the repository `.golangci.yml` MUST report no issue on changed
  code (gocyclo limit 15, line length 140).
- Pre-commit hooks (`.pre-commit-config.yaml`) MUST pass.
- The `cli-helper` hook regenerates `docs/github-pages/static/helper.out` from `-h`. A
  change to CLI flags MUST commit that regenerated file and update the README and the
  affected pages under `docs/`.
- A `//nolint` directive MUST carry an explanation (`nolintlint`).

Rationale: the same checks must apply to every contributor and to CI.

### V. Conventional Commits and Automated Releases

Releases are produced by tooling and MUST NOT be hand-edited.

- Commit messages MUST follow Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`,
  `ci:`, `test:`). `release-please` derives versions and `CHANGELOG.md` from them.
- Versions follow Semantic Versioning. A breaking change to CLI flags, config keys or
  output formats MUST be flagged as breaking in the commit.
- `CHANGELOG.md` MUST NOT be edited by hand.
- GoReleaser (`.goreleaser.yaml`) builds static binaries (`CGO_ENABLED=0`) with version,
  commit and build time injected through `-ldflags`.
- Commits are focused and atomic.

Rationale: predictable versioning for users who pin tf2report in CI pipelines.

## Technology and Delivery Constraints

- Language and toolchain: Go as declared in `go.mod` (currently 1.25). A toolchain bump
  MUST update `go.mod`, CI workflows and this section together.
- Platforms: Linux and macOS binaries via GoReleaser. Windows is currently disabled and
  MUST be re-enabled explicitly, not implicitly.
- The tool is read-only and offline: it reads a plan file and writes a report to stdout or
  a file. It MUST NOT call the Terraform CLI, cloud APIs or any network endpoint.
- Plan files and reports may contain sensitive infrastructure values. The tool MUST NOT
  log or echo raw attribute values beyond what the selected report format specifies, and
  fixtures committed to the repository MUST contain synthetic data only.
- Output formats: `markdown` (default), `text`, `json`. Filters: resource type and action.
- Build and local tasks go through the `Makefile`, written per
  `.github/instructions/makefile.instructions.md`.
- Documentation lives in `README.md` and `docs/`, without emojis. Example plans live in
  `examples/`, and example outputs MUST NOT be embedded anywhere.

## Development Workflow

- Work happens on a branch and lands through a pull request. The `Pull Requests` workflow
  MUST be green before merge.
- Spec Kit flow: `/speckit-specify` -> `/speckit-clarify` -> `/speckit-plan` ->
  `/speckit-tasks` -> `/speckit-implement`. Plans MUST include a Constitution Check.
- Reviewers MUST verify: tests present (III), lint and hooks clean (IV), docs and `helper.out`
  updated when the CLI changes (IV), Conventional Commit title (V).
- Documentation is updated in the same change as the behavior it describes.

## Governance

This constitution supersedes other project practices. Where `.github/instructions/` files
or other docs conflict with it, the constitution wins and the conflicting file MUST be
corrected.

- **Amendments**: proposed by pull request modifying this file, with rationale and version
  bump, and approved by a maintainer. Removing or redefining a principle MUST include a
  migration note for in-flight specs and plans.
- **Versioning**: semantic. MAJOR for removal or incompatible redefinition of a principle;
  MINOR for a new principle or section, or materially expanded guidance; PATCH for
  clarifications and wording.
- **Compliance review**: every pull request and every `/speckit-plan` Constitution Check
  verifies compliance. A deviation MUST be justified in the plan or PR description and
  recorded as accepted debt.
- **Runtime guidance**: `.github/instructions/go.instructions.md` and
  `.github/instructions/makefile.instructions.md` hold day-to-day conventions.

**Version**: 1.0.0 | **Ratified**: 2026-11-22 | **Last Amended**: 2026-11-22
