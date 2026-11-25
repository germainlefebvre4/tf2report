# Research: Baseline CLI Behavior

All items from the plan's Technical Context are resolved here. Findings come from the Repomix
snapshot (2026-11-22), the source, and runs of the built binary.

## R-001 Test layers

- **Decision**: three layers: unit tests per package, golden-file tests for formatters, and
  black-box CLI tests.
- **Rationale**: the spec's acceptance scenarios mix pure logic (action classification, counts),
  rendered text (formats) and process behavior (flags, streams, exit codes). Each layer is the
  cheapest one that can observe its behavior.
- **Alternatives considered**: only CLI-level tests (slow to diagnose, cannot cover
  `GetPrimaryAction` edge cases cheaply); only unit tests (cannot observe exit codes, flag
  precedence or the double error output).

## R-002 CLI harness without touching `main.go`

- **Decision**: `TestMain` in `cmd/tf2report` runs `go build` into a temporary directory (with
  `-ldflags "-X main.appVersion=... -X main.buildCommit=... -X main.buildTime=..."`), and each
  test runs the binary with `os/exec`, with its own `t.TempDir()` as working directory, `HOME`
  set to another temporary directory, and every `TF2REPORT_*` variable removed from the
  environment. Stdout, stderr and exit code are captured separately.
- **Rationale**: `main.go` declares `init()` and package-level flag variables; calling
  `rootCmd.Execute()` repeatedly in-process would share state between tests. A real process also
  gives true exit codes (FR-011).
- **Alternatives considered**: in-process `rootCmd.SetArgs(...)` (state leaks, `os.Exit` in
  `main` unreachable); refactoring `main.go` (production change, out of scope, tracked as
  `INIT`/`CMD_LOGIC`).
- **Pitfall found**: viper looks for `tf2report` with any or no extension in `.`. The binary must
  therefore live outside the test working directory, otherwise it is parsed as the config file
  (a Known Gap, pinned by a dedicated test that creates such a file on purpose).
- **Pitfall found**: the default config search also reads `$HOME/.config/tf2report`; isolating
  `HOME` keeps the developer's real configuration out of the tests.

## R-003 Golden files

- **Decision**: one golden file per fixture and format under `pkg/report/testdata/`, compared
  byte for byte; refreshed only with `go test ./pkg/report -update`. CLI tests compare a few
  full outputs against the same goldens and assert structure elsewhere.
- **Rationale**: report output is the contract (Constitution IV: "pin its output"). Byte-equality
  also proves determinism (SC-003) because each golden test renders twice and compares.
- **Alternatives considered**: substring assertions only (miss ordering and spacing regressions);
  a golden library (new dependency, rejected by Principle I).

## R-004 Fixtures and "supported Terraform versions"

- **Decision**: fixtures are synthetic, hand-written in the `terraform show -json` shape
  (`format_version` 1.x, resource addresses, `change.actions`, `before`/`after`, sensitive
  markers), with a `terraform_version` taken from a list kept in `testdata/README.md`. Values are
  obviously fake (`SYNTHETIC-...`). The list starts with the version already used by the
  repository's sample plan (1.10.5) and MUST be checked against HashiCorp's published support
  policy when the fixtures are written.
- **Rationale**: spec clarification 2 guarantees "versions currently supported by HashiCorp".
  A web search during planning did not return HashiCorp's own policy page (only third-party
  pages), so the exact list is **not asserted here**; it is a task to confirm it from the source
  at implementation time and record the date in `testdata/README.md`. The JSON plan format
  exposes its own `format_version`, which the tool does not validate; tests do not assume it does.
- **Alternatives considered**: capturing real plans (risk of real infrastructure data, forbidden
  by the Constitution); one fixture per Terraform minor version (no behavioral difference is
  expected since the tool reads only a stable subset of fields).

## R-005 Pinning Known Gaps

- **Decision**: every Known Gap gets one test named `TestKnownGap_<Slug>` with a comment
  `PROVISIONAL: pins Known Gap "<title>" (spec.md). Update deliberately when the gap is fixed.`
  Slugs: `UnknownActionFilter`, `VerbosityOnlyViaVerboseFlag`, `ErrorPrintedTwiceWithUsage`,
  `ExtensionlessConfigFile`, `ModuleAndIndexNotShown`, `NestedFilterEnvVariable` (only the dotted
  name works), `RealPlanShapeRejected` (boolean sensitivity flags and scalar output values fail to
  parse; verified by running the binary on a real-shaped plan). The README gap (plan path described as required, outdated example output) is a
  documentation gap: it has no provisional test, and the default plan path is covered by a regular
  test. The missing-tests gap is closed by this feature itself.
- **Rationale**: spec FR-016. A common prefix lets `go test -run KnownGap ./...` enumerate them.
- **Alternatives considered**: build tags or skipping (hides the pinned behavior); marking
  expected failures (spec clarification 1 chose pinning current behavior, not intended).

## R-006 Quality gate for test code

- **Decision**: interim gate is `gofmt -l .`, `goimports -local github.com/germainlefebvre4/tf2report`,
  `go vet ./...`, `go test ./...` and the pre-commit hooks. `golangci-lint` is not used because
  v2.15 rejects the repository's v1 `.golangci.yml`.
- **Rationale**: observed: `golangci-lint run ./...` fails with "unsupported version of the
  configuration". Fixing the config is tracked in the constitution report (`LINT_CI`).
- **Alternatives considered**: migrating the lint config in this feature (scope creep, may force
  production changes); running golangci-lint with `--no-config` (would not apply project rules).

## R-007 Scope of "no application change"

- **Decision**: no production file is edited. If a test cannot be written without a production
  change, it is recorded as a Known Gap and skipped from the baseline, not worked around by a
  code edit (spec Assumptions: "any code change is limited to what tests strictly need").
- **Rationale**: keeps the baseline honest and the diff reviewable as test-only.
- **Alternatives considered**: exposing internals for testing (changes the public surface).
