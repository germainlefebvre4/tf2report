# Quickstart: validating the baseline

Run from the repository root. Prerequisites: Go 1.25 or later. No network, no Terraform, no cloud
account needed.

## 1. Run the whole suite

```bash
go test ./...
```

Expected: all packages pass. The CLI package builds the binary once at start-up, so the first run
takes a few seconds longer.

## 2. Check coverage of the baseline

```bash
make test-coverage        # writes coverage.out and coverage.html
go tool cover -func=coverage.out | tail -1
```

Expected: every function of `pkg/terraform`, `pkg/config` and `pkg/report` is exercised; `run()`
in `cmd/tf2report` is exercised through the binary (not counted by the in-process coverage tool).
No numeric threshold is set by the spec (SC-007 requires scenario coverage, not a percentage).

## 3. List the provisional tests (Known Gaps)

```bash
go test -run KnownGap -v ./...
```

Expected: one passing test per Known Gap, each named `TestKnownGap_<Slug>`. When a gap is fixed
in a later feature, the matching test is updated or removed in the same change.

## 4. Validate one story by hand

```bash
go build -o /tmp/tf2report-bin ./cmd/tf2report
mkdir -p /tmp/tf2report-work && cd /tmp/tf2report-work
/tmp/tf2report-bin --plan "$OLDPWD/testdata/plans/basic.json"
```

Expected (SC-001): Markdown report with 4 total changes: 1 add, 1 change, 1 destroy, 1 replace.
Repeat with `--format text` and `--format json`: same counts. Run the binary from a directory that
does not contain a file named `tf2report` (see the Known Gap on configuration lookup).

Failure checks (User Story 5), each ends with exit code 1 and no standard output:

```bash
/tmp/tf2report-bin -p missing.json; echo $?
/tmp/tf2report-bin -p "$OLDPWD/testdata/plans/invalid.json"; echo $?
/tmp/tf2report-bin -p "$OLDPWD/testdata/plans/basic.json" -f xml; echo $?
```

## 5. Refresh golden files (only on a deliberate output change)

```bash
go test ./pkg/report -update
git diff pkg/report/testdata      # review every changed line
```

## 6. Interim quality gate

```bash
gofmt -l .
go vet ./...
pre-commit run --all-files
```

Expected: no output from `gofmt -l`, `go vet` clean, hooks pass. `golangci-lint` is not part of the
gate until its configuration is migrated (see `plan.md`, Complexity Tracking).

## Success checklist

| Criterion | How to verify |
|-----------|---------------|
| SC-001 | Step 4, and the `basic` golden tests |
| SC-002 | `NewSummary` tests assert counts against fixtures with mixed actions |
| SC-003 | Golden tests render twice and compare bytes |
| SC-004 | CLI failure tests assert exit code 1 and a one-line cause on standard error |
| SC-005 | CLI test runs with only a plan in the working directory and no config |
| SC-006 | `sensitive.json` run through all formats; marker values must be absent |
| SC-007 | Traceability table in `plan.md`: each scenario maps to at least one passing test |
