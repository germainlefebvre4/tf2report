# Tasks: Baseline CLI Behavior

**Input**: Design documents from `specs/001-baseline-cli-behavior/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: This feature **is** the tests (spec FR-016, SC-007). Every user-story task writes
characterization tests that pin the current behavior, so they pass on first run; there is no
"write failing test first" step. **No production file may be modified** (research R-007): all
paths below are `_test.go` files, `testdata/` files, or documentation.

**Organization**: Tasks are grouped by user story. Each story's tests are independent of the
others once Phase 2 is done.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an incomplete task)
- **[Story]**: US1 to US5, mapping to the user stories of `spec.md`
- Module: `github.com/germainlefebvre4/tf2report`; goimports local prefix is that module path
- All Go code: stdlib only, `gofmt`/`goimports`, `t.Run` subtests, `t.Helper()`, `t.Cleanup()`,
  no `init()` and no new package-level mutable state except the test-only `-update` flag
- Provisional tests are named `TestKnownGap_<Slug>` and start with the comment
  `// PROVISIONAL: pins Known Gap "<title>" (specs/001-baseline-cli-behavior/spec.md). Update deliberately when the gap is fixed.`

## Path Conventions

Single Go module. Tests live next to the code (`pkg/*/`, `cmd/tf2report/`). Shared synthetic plans
live in `testdata/plans/` at the repository root; golden files in `pkg/report/testdata/`.

---

## Phase 1: Setup

**Purpose**: fixture area and the supported-version decision.

- [x] T001 Create `testdata/README.md` stating: fixtures are synthetic only (no real infrastructure data, Constitution Technology constraints), values use obvious markers such as `SYNTHETIC-...`, golden files are refreshed only with `-update`, and the complete fixture table `Fixture | terraform_version | format_version | purpose` with its 8 rows already listed (`basic.json`, `modules-and-counts.json`, `noop-and-read.json`, `sensitive.json`, `no-changes.json`, `unknown-action.json`, `empty-object.json`, `invalid.json`, with the purpose given in T003 to T009 and the version column left as `(set in T002)`); this is the only task that creates the table, so the fixture tasks below never edit this file; create the empty directory `testdata/plans/`
- [x] T002 Confirm from HashiCorp's official source (WebFetch on developer.hashicorp.com / releases.hashicorp.com, not third-party pages) which Terraform versions are currently supported, then record the list and the check date in `testdata/README.md`; pick the `terraform_version` used by every fixture from that list (research R-004) and fill the `terraform_version` and `format_version` columns of the fixture table created in T001. If the official page cannot be reached, stop and ask the user instead of guessing

---

## Phase 2: Foundational (blocking prerequisites)

**Purpose**: fixtures and shared test helpers that every story uses.

**CRITICAL**: no story phase starts before this phase is complete.

- [x] T003 [P] Create `testdata/plans/basic.json`: `format_version` "1.2", `terraform_version` from T002, exactly 4 `resource_changes` derived from the repository sample `terraform.tfplan.json` but with synthetic values: `aws_s3_bucket.data` actions `["create"]`; `aws_instance.web` (distinct `name` "web", address `aws_instance.web`) `["update"]`; `aws_instance.app` `["delete","create"]`; `aws_security_group.old` `["delete"]`; plus an `output_changes` entry shaped as an object, `{"endpoint": {"actions": ["create"], "before": null, "after": {"value": "SYNTHETIC-OUTPUT-VALUE"}}}` (it is parsed but never reported; an object is used because the parser rejects the real scalar shape, see T009b). Result: 1 add, 1 change, 1 destroy, 1 replace (SC-001)
- [x] T004 [P] Create `testdata/plans/modules-and-counts.json`: 5 changes all `["create"]` or `["update"]` with the addresses `module.net.aws_subnet.main`, `module.dr.aws_subnet.main`, `aws_instance.web[0]`, `aws_instance.web[1]`, `aws_s3_bucket.logs["eu"]` (matching `type`/`name` fields: same `type` and `name` for the two subnets and for the two web instances). Purpose: Known Gap "module path and index not shown"
- [x] T005 [P] Create `testdata/plans/noop-and-read.json`: 1 real `["create"]` change plus one `["no-op"]` change, one `["read"]` data-source change (`"mode": "data"`) and one change with `"actions": []`. Expected counted total: 1
- [x] T006 [P] Create `testdata/plans/sensitive.json`: 2 changes (`create`, `update`) whose `before`/`after` hold the marker strings `SYNTHETIC-SECRET-PASSWORD` and `SYNTHETIC-SECRET-TOKEN`, with `after_sensitive` flagging them. Purpose: SC-006 / FR-013
- [x] T007 [P] Create `testdata/plans/no-changes.json`: valid plan with `terraform_version` and `"resource_changes": []`
- [x] T008 [P] Create `testdata/plans/unknown-action.json`: 1 valid `["create"]` change plus one change with `"actions": ["frobnicate"]`. Used by T015 and T019 (an unknown action is left out of counts and reports, spec Edge Cases)
- [x] T009b [P] Create `testdata/plans/real-shape.json`: one `["create"]` resource whose change carries `"before_sensitive": false` and `"after_sensitive": {}`, plus `"output_changes": {"endpoint": {"actions": ["create"], "before": null, "after": "SYNTHETIC-OUTPUT-VALUE", "after_unknown": false, "before_sensitive": false, "after_sensitive": false}}`, i.e. the shape `terraform show -json` really produces; purpose: Known Gap `RealPlanShapeRejected`; add its row to the fixture table (the README already lists the others)
- [x] T009 [P] Create `testdata/plans/empty-object.json` containing exactly `{}` and `testdata/plans/invalid.json` containing the text `this is not json`
- [x] T010 [P] Create `pkg/terraform/helpers_test.go` (package `terraform`) with `func fixture(t *testing.T, name string) []byte` that reads `../../testdata/plans/<name>` and fails the test with `t.Fatalf` on error
- [x] T011 [P] Create `pkg/report/golden_test.go` (package `report`) with: a test-only `var update = flag.Bool("update", false, "rewrite golden files")`; `func loadPlan(t *testing.T, name string) (*terraform.Plan, *terraform.Summary)` that reads `../../testdata/plans/<name>`, parses with `terraform.NewParser().Parse` and builds `terraform.NewSummary(plan.ResourceChanges)`; `func render(t *testing.T, f Formatter, plan *terraform.Plan, summary *terraform.Summary) string` that formats into a `bytes.Buffer` **twice**, fails the test if the two outputs differ (determinism, FR-014, SC-003) and returns the output; `func assertGolden(t *testing.T, name, got string)` that compares `got` with `testdata/<name>` byte for byte and rewrites the file only when `*update` is set
- [x] T012 [P] Create `pkg/config/helpers_test.go` (package `config`) with `func isolate(t *testing.T) (workDir, homeDir string)` that creates two `t.TempDir()`s, calls `t.Chdir(workDir)` and `t.Setenv("HOME", homeDir)`, and unsets every existing `TF2REPORT_*` variable via `os.Environ()` + `t.Setenv(name, "")` + `os.Unsetenv`; plus `func writeFile(t *testing.T, path, content string)` creating parent directories
- [x] T013 Create `cmd/tf2report/main_test.go` (package `main`): `TestMain` builds the binary once with `go build -ldflags "-X main.appVersion=test-1.2.3 -X main.buildCommit=abc1234 -X main.buildTime=2026-01-01T00:00:00Z" -o <tmp>/tf2report .` into an `os.MkdirTemp` directory that is **not** a working directory of any test (research R-002 pitfall), runs `m.Run()`, then removes the temp directory; plus helper `runCLI(t *testing.T, o cliOpts) cliResult` where `cliOpts{Args []string; Files map[string]string; HomeFiles map[string]string; Env []string}` and `cliResult{Stdout, Stderr string; Code int}`: it creates a `t.TempDir()` cwd and a separate `HOME` temp dir, writes `Files` into the cwd and `HomeFiles` into the temporary `HOME` (paths relative to `HOME`, for example `.config/tf2report/tf2report.yaml`), starts from a copy of `os.Environ()` with every `TF2REPORT_*` removed and `HOME` replaced, appends `Env`, runs the binary with `os/exec` capturing stdout and stderr separately, and maps `*exec.ExitError` to `Code`; plus `func fixturePath(t *testing.T, name string) string` returning the absolute path of `../../testdata/plans/<name>`

**Checkpoint**: fixtures and helpers ready; user-story test tasks can start in parallel.

---

## Phase 3: User Story 1 - Summarize a Terraform plan (Priority: P1) MVP

**Goal**: pin plan parsing, action classification, counting and the default (Markdown) report.

**Independent Test**: `go test ./pkg/terraform ./pkg/report ./cmd/tf2report -run 'Summary|PrimaryAction|Parse|Markdown|KnownGap_ModuleAndIndex|DefaultPlanPath'` passes, and the CLI test on `basic.json` shows 4 changes (1 add, 1 change, 1 destroy, 1 replace).

- [x] T014 [P] [US1] Create `pkg/terraform/types_test.go`: table-driven `TestChange_GetPrimaryAction` with one row per line of the "Primary action derivation" table in `data-model.md` (`[]`, `create`, `update`, `delete`, `read`, `no-op`, `["delete","create"]`, `["create","delete"]` both replace, two actions that are not delete+create such as `["update","delete"]` where the first action decides, unknown value `["frobnicate"]` no-op)
- [x] T015 [P] [US1] Create `pkg/terraform/summary_test.go`: table-driven `TestNewSummary` on plans built from `basic.json`, `noop-and-read.json`, `unknown-action.json` (total 1: the unknown action is left out) and `no-changes.json` (parse with `NewParser().Parse(fixture(...))`): assert `TotalChanges`, `ToAdd`, `ToChange`, `ToDestroy`, `ToReplace`, `ByType` entries, that read/no-op are excluded, and the invariants "total = add+change+destroy+replace" and "sum of ByType equals the totals" (SC-002)
- [x] T016 [P] [US1] Create `pkg/terraform/parser_test.go`: `TestParser_Parse` and `TestParser_ParseFile` success cases: `basic.json` gives 4 `ResourceChanges` and the `TerraformVersion` recorded in the fixture; `empty-object.json` and a JSON without `resource_changes` give 0 changes and no error; unknown top-level keys are ignored; `output_changes` of `basic.json` is parsed into `OutputChanges`; `ParseFile` on `../../testdata/plans/basic.json` equals `Parse` of the same bytes (use `reflect.DeepEqual`)
- [x] T017 [P] [US1] Create `pkg/report/markdown_test.go`: table-driven `TestMarkdownFormatter` over `basic.json`, `no-changes.json`, `noop-and-read.json` using `assertGolden` with golden files `markdown_basic.golden`, `markdown_no_changes.golden`, `markdown_noop_and_read.golden`; add a subtest rendering a plan with an empty `TerraformVersion` and asserting there is no `**Terraform Version:**` line. Generate the goldens with `go test ./pkg/report -update`, then review each against `contracts/cli-contract.md` (Markdown layout) before committing them
- [x] T018 [P] [US1] Create `pkg/report/gaps_test.go` with `TestKnownGap_ModuleAndIndexNotShown`: render `modules-and-counts.json` as Markdown and as text via `assertGolden` (`markdown_modules_and_counts.golden`, `text_modules_and_counts.golden`) and additionally assert that the two `main` subnet rows and the two `web` rows are textually identical (the gap). Mark PROVISIONAL per the convention above
- [x] T019 [P] [US1] Create `cmd/tf2report/summary_test.go`: using `runCLI` and `fixturePath`: (a) `-p basic.json` exits 0, stdout contains `**Total Changes:** 4` and rows `| Add | 1 |`, `| Change | 1 |`, `| Destroy | 1 |`, `| Replace | 1 |`, stderr empty (SC-001, SC-005); (b) the replace resource is counted once and not as add+destroy; (c) `noop-and-read.json` reports 1 change only; (d) `no-changes.json` and `empty-object.json` exit 0 with `No changes. Infrastructure is up-to-date.`; (e) `unknown-action.json` reports `**Total Changes:** 1` and exits 0; (f) `TestSummary_DefaultPlanPath` (regular test, **not** provisional, FR-001 and US1 scenario 5): no `--plan`, with `Files{"terraform.tfplan.json": <content of basic.json read with os.ReadFile>}` in the working directory, exits 0 and reports 4 changes. The README's claim that `--plan` is required is a documentation gap and is not pinned by a test

**Checkpoint**: US1 is fully testable on its own: the default report is pinned.

---

## Phase 4: User Story 2 - Choose the report format (Priority: P2)

**Goal**: pin format selection, the text and JSON layouts, and the "no attribute value" guarantee.

**Independent Test**: `go test ./pkg/report ./cmd/tf2report -run 'Format|Text|JSON|Safety'` passes; the three formats show the same counts for `basic.json`.

- [x] T020 [P] [US2] Create `pkg/report/formatter_test.go`: table-driven `TestParseFormat` (`markdown`, `md`, `text`, `txt`, `json` succeed and return `FormatMarkdown`/`FormatText`/`FormatJSON`; `xml`, `""`, `MARKDOWN` fail with the exact message `invalid format: <x> (supported: markdown, text, json)`); table-driven `TestNewFormatter` asserting the concrete types `*MarkdownFormatter`, `*TextFormatter`, `*JSONFormatter` and the error `unsupported format: <x>` for an unknown `Format` value
- [x] T021 [P] [US2] Create `pkg/report/text_test.go`: table-driven `TestTextFormatter` over `basic.json`, `no-changes.json`, `noop-and-read.json` with goldens `text_basic.golden`, `text_no_changes.golden`, `text_noop_and_read.golden`; assert no Markdown syntax (`|`, `#`, backticks) in the basic output. Generate with `-update`, review against `contracts/cli-contract.md` (Text layout)
- [x] T022 [P] [US2] Create `pkg/report/json_test.go`: `TestJSONFormatter` with goldens `json_basic.golden`, `json_no_changes.golden`, `json_noop_and_read.golden`; additionally `json.Unmarshal` the basic output into `JSONReport` and assert the `summary` counters, `changes` in plan order with the fields `address`, `type`, `name`, `action`, `by_type` with `resources` addresses, `terraform_version` present; for `empty-object.json` assert `terraform_version` is omitted and `changes` is `[]` and `by_type` is `{}` (decode into `map[string]json.RawMessage` to check key absence); assert two-space indentation and a trailing newline
- [x] T023 [P] [US2] Create `pkg/report/safety_test.go`: `TestFormatters_NeverOutputAttributeValues` rendering `sensitive.json` with all three formatters and asserting neither `SYNTHETIC-SECRET-PASSWORD` nor `SYNTHETIC-SECRET-TOKEN` appears (SC-006, FR-013); also render `basic.json` and assert the marker `SYNTHETIC-OUTPUT-VALUE` of `output_changes` appears in none of the formats
- [x] T024 [P] [US2] Create `cmd/tf2report/format_test.go`: default format is Markdown; `-f text`, `-f json` (stdout parses as JSON); `-f md` output equals `-f markdown` and `-f txt` equals `-f text` byte for byte; the counts (4 total, 1/1/1/1) are identical in the three formats for `basic.json` (FR-006, SC-001); `-f xml` exits 1 with `invalid format: xml (supported: markdown, text, json)` on stderr and empty stdout; `sensitive.json` through the CLI in the three formats never prints the two marker strings

**Checkpoint**: US1 and US2 both pass independently.

---

## Phase 5: User Story 3 - Focus with filters (Priority: P3)

**Goal**: pin type and action filtering, AND semantics, zero-match behavior and the unknown-action gap.

**Independent Test**: `go test ./pkg/terraform ./cmd/tf2report -run 'Filter|KnownGap_UnknownActionFilter'` passes.

- [x] T025 [P] [US3] Create `pkg/terraform/filter_test.go`: table-driven `TestFilterByType` (nil/empty types returns the input unchanged; one type; several types; type absent from the plan returns an empty non-nil slice) and `TestFilterByAction` (empty actions returns input; each of `ActionCreate`, `ActionUpdate`, `ActionDelete`, `ActionReplace` selects the matching resources of `basic.json` with replace judged on the primary action; `ActionType("bogus")` matches nothing); a subtest chaining `FilterByType` then `FilterByAction` to pin AND semantics
- [x] T026 [P] [US3] Create `cmd/tf2report/filters_test.go`: on `basic.json`: `-t aws_instance` yields only that type and totals reflect it; repeated `-t a -t b` and comma form `-t a,b`; `-a create`; `-a delete -a replace`; `-t aws_instance -a create` (both must match); a filter matching nothing exits 0 with `No changes. Infrastructure is up-to-date.`; plus `TestKnownGap_UnknownActionFilter`: `-a bogus` exits 0 and reports no changes instead of failing (PROVISIONAL)

**Checkpoint**: US1 to US3 pass independently.

---

## Phase 6: User Story 4 - Configure once, run anywhere (Priority: P3)

**Goal**: pin defaults, file discovery, environment variables, precedence and the three configuration gaps.

**Independent Test**: `go test ./pkg/config ./cmd/tf2report -run 'Config|Loader|KnownGap_(Extensionless|Nested|Verbosity)'` passes.

- [x] T027 [P] [US4] Create `pkg/config/types_test.go`: `TestNewConfig_Defaults` (`terraform.tfplan.json`, `markdown`, `info`, non-nil empty filter slices) and table-driven `TestFilters_GetActionTypes` (nil/empty gives an empty non-nil slice; `["create","delete"]` maps to the matching `terraform.ActionType`; an unknown string is passed through unchanged)
- [x] T028 [P] [US4] Create `pkg/config/loader_test.go` using `isolate`/`writeFile`: defaults when no file exists; explicit file (`Load(path)`) values applied for all five keys; explicit missing file returns an error containing `failed to open config file`; invalid YAML returns an error containing `failed to read config`; `tf2report.yaml` in the working directory is found; `$HOME/.config/tf2report/tf2report.yaml` is found; `TF2REPORT_OUTPUT_FORMAT` applies and wins over the file value; `TF2REPORT_VERBOSITY` applies. Then three provisional tests: `TestKnownGap_NestedFilterEnvVariable` (`t.Setenv("TF2REPORT_FILTERS.ACTIONS","delete")` is honored while `TF2REPORT_FILTERS_ACTIONS` is ignored), `TestKnownGap_ExtensionlessConfigFile` (a file named `tf2report`, without extension, holding YAML in the working directory is loaded as configuration)
- [x] T029 [P] [US4] Create `cmd/tf2report/config_test.go`: `Files{"tf2report.yaml": ...}` applied with no flags (format, plan path, filters); `--config <file>` explicit; a flag overrides the file (`-f json` over `output_format: text`; `-p` over `terraform_plan_path`); `TF2REPORT_OUTPUT_FORMAT=json` via `Env`; a config under `$HOME/.config/tf2report/` (use `HomeFiles` from T013; this task does not edit `main_test.go`); defaults with no configuration produce the Markdown report for `terraform.tfplan.json`; plus `TestKnownGap_ExtensionlessConfigFile` (a text file named `tf2report` in the cwd, containing the single line `hello world`, makes the run exit 1 with `failed to load config`) and `TestKnownGap_VerbosityOnlyViaVerboseFlag` (`verbosity: debug` in the file, no `-v`: stderr stays empty; with `-v`: stderr shows `Parsed plan with 4 resource changes`). Mark both PROVISIONAL

**Checkpoint**: US1 to US4 pass independently.

---

## Phase 7: User Story 5 - Understand and trust failures (Priority: P3)

**Goal**: pin error handling, streams, exit codes, verbose output, help and version.

**Independent Test**: `go test ./pkg/terraform ./cmd/tf2report -run 'Error|Diagnostic|KnownGap_ErrorPrinted'` passes.

- [x] T030 [P] [US5] Create `pkg/terraform/parser_errors_test.go`: table-driven failure cases for `Parser.Parse` (`invalid.json` content, empty bytes, truncated JSON, `{"resource_changes": 3}`) asserting an error containing `failed to parse plan JSON` and no non-nil plan; `ParseFile` on a missing path returns an error containing `failed to read plan file` and `errors.Is(err, os.ErrNotExist)` is true (proves `%w` wrapping, Principle I); plus `TestKnownGap_RealPlanShapeRejected` in the same file: `real-shape.json` fails to parse with `cannot unmarshal bool into Go struct field`, and a plan whose only issue is a scalar `output_changes[...].after` fails with `cannot unmarshal string` (PROVISIONAL)
- [x] T031 [P] [US5] Create `cmd/tf2report/errors_test.go`: missing plan file (`plan file not found: <path>`), `invalid.json` (`failed to parse plan:`), unsupported format, unreadable explicit config (`failed to load config`); each asserts exit code 1, empty stdout and the expected cause on stderr (SC-004); plus `TestKnownGap_ErrorPrintedTwiceWithUsage`: for the missing-file case, stderr contains the cause line twice and the `Usage:` text (PROVISIONAL)
- [x] T032 [P] [US5] Create `cmd/tf2report/diagnostics_test.go`: `-v` writes `Parsed plan with 4 resource changes` on stderr while stdout equals the run without `-v`; `-v -t aws_instance` adds `Filtered by type: N resources remain`; `-v -a create` adds `Filtered by action: N resources remain`; without `-v` stderr is empty; `--help` exits 0 with `Usage:` on stdout; `--version` exits 0 and prints the four lines `tf2report`, `Version: test-1.2.3`, `Commit: abc1234`, `Build time: 2026-01-01T00:00:00Z` (values injected by T013)

**Checkpoint**: all five stories are pinned.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: quality gate, documentation and proof that the baseline is test-only.

- [x] T033 [P] Update `docs/DEVELOPMENT.md` (check its current Testing section first; no emojis; no example outputs, per `.github/instructions/main.instructions.md`): document the three test layers, `go test ./pkg/report -update` for goldens, `testdata/`, the `TestKnownGap_*` convention and `go test -run KnownGap ./...`
- [x] T034 Run the interim quality gate (research R-006) and fix findings in the new files only: `gofmt -l .` (no output), `goimports -local github.com/germainlefebvre4/tf2report -l .` (no output), `go vet ./...`
- [x] T035 Run `go test ./...` and `make test-coverage`; all green (`coverage.out` and `coverage.html` are git-ignored, see `.gitignore`; delete them afterwards). Then run `go test -run KnownGap -v ./...` and check that exactly 9 provisional tests (7 distinct gaps) exist and pass: `UnknownActionFilter`, `VerbosityOnlyViaVerboseFlag`, `ErrorPrintedTwiceWithUsage`, `ExtensionlessConfigFile` (two tests, one in `pkg/config` and one in `cmd/tf2report`, counted as one gap), `ModuleAndIndexNotShown`, `NestedFilterEnvVariable`, `RealPlanShapeRejected` (also two tests, `pkg/terraform` and `cmd/tf2report`, counted as one gap); the README gap is documentation-only and has no provisional test
- [x] T036 Prove the change is test-only (R-007): `git diff --name-only main -- cmd pkg | grep -v '_test.go$'` must print nothing, and so must `git status --porcelain -- cmd pkg | grep -v -e '_test.go$' -e 'testdata/$'` (catches new untracked production files); also run `grep -rEn '"net/http"|"os/exec"' cmd pkg --include='*.go' --exclude='*_test.go'`, which must print nothing (FR-013: no network call, no Terraform execution, verified by inspection). If a production file shows up, revert it and record a Known Gap instead
- [x] T037 Execute every step of `specs/001-baseline-cli-behavior/quickstart.md` (steps 1 to 6) and tick the success checklist (SC-001 to SC-007); if a test name or path differs from the plan, fix `plan.md` (project structure and traceability table) so documentation matches the code
- [x] T038 Run `SKIP=hadolint-docker pre-commit run --all-files` (Constitution Principle IV: hooks MUST pass; only the `hadolint-docker` hook, which needs Docker and has no Dockerfile change to check, is skipped, and that skip is stated in the pull request description). If Docker is available, run `pre-commit run --all-files` without `SKIP`. Done 2026-11-25: all hooks pass except `update-dockerhub-readme`, which rewrites `docs/dockerhub/README.md` as a side effect of `--all-files` (unrelated to this feature; reverted, do not commit that change)
- [x] T039 Prepare the governance artifacts before committing: in `.specify/memory/constitution.md` delete the temporary `SYNC IMPACT REPORT` comment (it must not be committed); `TODO(TESTS)` is closed by this feature, so move the still-open items (`LINT_CI`, `INIT`, `CMD_LOGIC`) to the pull request description or to GitHub issues so they are not lost with the report; confirm `.gitignore` ignores `/coverage.out` and `/coverage.html`; list the files to commit with `git status --short` and check that none exceeds 500 KB (pre-commit `check-added-large-files`)
- [x] T040 Commit in focused Conventional Commit units (Principle V), for example `chore: add spec kit scaffolding` (`.specify/`, `.claude/skills/`, `.gitignore`), `docs: ratify constitution v1.0.0 and add baseline spec, plan and tasks` (constitution without its Sync Impact Report, `specs/`), `test: add shared fixtures and test helpers`, `test(terraform): pin parser, actions, summary and filters`, `test(report): pin markdown, text and json output`, `test(config): pin loader and precedence`, `test(cli): pin end-to-end behavior and known gaps`, `docs: document characterization tests`; never edit `CHANGELOG.md`

---

## Dependencies & Execution Order

### Phase dependencies

- **Phase 1 (Setup)**: none. T002 gates the `terraform_version` value of the fixtures.
- **Phase 2 (Foundational)**: after T001 and T002. T003 to T012 are mutually parallel; T013 is independent of them except for the fixture files it reads at test time.
- **Phases 3 to 7 (stories)**: all depend on Phase 2 only; they can run in any order or in parallel. Recommended order: P1, P2, then P3 stories.
- **Phase 8 (Polish)**: after the stories you intend to ship. T034 to T038 run strictly in order; T033 can run any time after Phase 2.

### Cross-story file notes

- `cmd/tf2report/main_test.go` (T013, including `HomeFiles`) is the only shared CLI file and no story task edits it; the story files `summary_test.go`, `format_test.go`, `filters_test.go`, `config_test.go`, `errors_test.go`, `diagnostics_test.go` never overlap.
- `pkg/report/` helpers (T011) are shared by T017, T018, T021, T022, T023.
- `TestKnownGap_ExtensionlessConfigFile` exists at both levels on purpose (T028 loader, T029 CLI).

### Within each story

Write the test file, run it, and when it uses goldens generate them with `-update` and **review the
diff against `contracts/cli-contract.md`** before accepting. A characterization test that fails on
the unmodified application means the test is wrong (or exposes an undocumented behavior: record it
in `spec.md` Known Gaps), never that the application must change.

### Parallel opportunities

- Phase 2: T003 to T012 all [P].
- Phase 3: T014 to T019 all [P] (six different files).
- Phase 4: T020 to T024 all [P]. Phase 5: T025 and T026. Phase 6: T027 to T029. Phase 7: T030 to T032.
- After Phase 2, the five story phases can proceed in parallel across people or agents.

### Parallel example: User Story 1

```text
Task: "TestChange_GetPrimaryAction in pkg/terraform/types_test.go"
Task: "TestNewSummary in pkg/terraform/summary_test.go"
Task: "TestParser_Parse in pkg/terraform/parser_test.go"
Task: "TestMarkdownFormatter goldens in pkg/report/markdown_test.go"
Task: "TestKnownGap_ModuleAndIndexNotShown in pkg/report/gaps_test.go"
Task: "CLI summary tests in cmd/tf2report/summary_test.go"
```

---

## Implementation Strategy

### MVP first (User Story 1 only)

1. Phase 1 and Phase 2.
2. Phase 3 (US1): the default report is pinned, which is the product's core value.
3. **Stop and validate**: `go test ./...` green, quickstart step 4 (SC-001) by hand.
4. Commit and open a pull request if a smaller first increment is preferred.

### Incremental delivery

1. Setup + Foundational, then US1 (MVP).
2. US2 formats and safety, then US3 filters, then US4 configuration, then US5 failures.
3. Each story is mergeable on its own because it only adds test files and fixtures.
4. Finish with Phase 8 (documentation, gate, test-only proof).

## Notes

- 40 tasks in total; Known Gap coverage matches `spec.md` (FR-016) and `research.md` R-005.
- No task edits `cmd/tf2report/main.go` or any `pkg/**/*.go` production file.
- Verify the "no application change" rule after every story, not only at the end (T036).
