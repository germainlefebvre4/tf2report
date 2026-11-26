# Quickstart: validating the SVG resource diagram

Runnable checks that prove the feature end to end once implemented. They use only local files and
the synthetic fixtures; nothing needs a cloud account. Contracts: see
[contracts/cli-contract.md](contracts/cli-contract.md); entities: [data-model.md](data-model.md).

## Prerequisites

- Go as declared in `go.mod`, repository checked out on this feature's branch.
- A browser or any SVG viewer to look at the result.
- No Terraform CLI, network or credentials.

```bash
make build            # builds ./bin/tf2report (see the Makefile)
make test             # go test ./...
```

## 1. Draw a plan with dependencies (user stories 1 and 2)

```bash
./bin/tf2report --plan testdata/plans/diagram/aws-network.json --format svg > /tmp/aws-network.svg
```

Expected: exit 0, standard error empty, a valid SVG. Open it and check:

- one `aws` container, resources inside it, one node per changed resource;
- create nodes green with `+`, update orange with `~`, delete red with `-`, replace purple with `-/+`;
- an arrow from each dependent resource to its dependency (instance to subnet to vpc);
- the legend lists only the statuses present;
- no attribute value anywhere in the file (see check 5).

## 2. Unchanged resources and dependency context (FR-003, FR-024)

```bash
./bin/tf2report -p testdata/plans/diagram/modules-context.json -f svg > /tmp/ctx.svg
./bin/tf2report -p testdata/plans/diagram/modules-context.json -f svg --action create,no-op > /tmp/noop.svg
./bin/tf2report -p testdata/plans/diagram/modules-context.json -f svg --type aws_subnet > /tmp/type.svg
```

Expected: `ctx.svg` shows the changed resources plus the unchanged resource they reference, muted, with
its arrow; `noop.svg` shows every unchanged resource muted; `type.svg` shows only subnets, no muted
nodes and no dangling arrow.

## 3. Cycles, unknown providers, no configuration (spec edge cases)

```bash
for p in cycle no-configuration; do
  ./bin/tf2report -p testdata/plans/diagram/$p.json -f svg > /tmp/$p.svg || echo "FAIL $p"
done
./bin/tf2report -p testdata/plans/terraform/provider-random.json -f svg > /tmp/random.svg
```

Expected: all exit 0; `cycle.svg` is readable (no endless processing); `no-configuration.svg` has
no arrows; `random.svg` uses the generic icon.

## 4. Empty plan and determinism (FR-020, SC-003)

```bash
./bin/tf2report -p testdata/plans/terraform/empty-config.json -f svg > /tmp/empty.svg
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg | sha256sum
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg | sha256sum
```

Expected: `empty.svg` is valid and says there is nothing to show; the two checksums are identical.

## 5. No attribute value in the image (FR-017, SC-005)

```bash
./bin/tf2report -p testdata/plans/terraform/sensitive.json -f svg > /tmp/sensitive.svg
go test ./pkg/report -run 'TestSVG_NoAttributeValues'
```

Expected: the test passes; it searches the output for every value of the sensitive fixtures.

## 6. Icons and icon packs (user story 3)

```bash
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg --icons-dir testdata/icons/valid > /tmp/pack.svg
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg --icons-dir testdata/icons/hostile > /tmp/hostile.svg
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg --icons-dir /nonexistent; echo "exit=$?"
./bin/tf2report -p testdata/plans/diagram/aws-network.json -f json --icons-dir /nonexistent -v 2>&1 >/dev/null | tail -1
```

Expected: `pack.svg` uses the pack icons where families match; `hostile.svg` is produced and contains
no script, event attribute or external reference; the third command prints an error naming the
directory and exits non-zero with empty standard output; the fourth prints one verbose notice that
`--icons-dir` is ignored and still exits 0.

```bash
grep -ciE '<script|onload=|foreignObject|href="http' /tmp/hostile.svg   # expected: 0
```

## 7. Regression guard for the other formats (FR-022, SC-008)

```bash
go test ./...        # existing goldens and examples tests must stay green, unchanged
```

## 8. Scale (SC-006)

```bash
go test ./pkg/report ./pkg/diagram -run 'TestScale' -v
```

Expected: the generated 300-resource plan renders in under 2 seconds and two runs are identical.

## 9. Documentation and gates (Constitution IV)

```bash
make fmt && go vet ./...             # go fmt, then go vet (no vet target in the Makefile)
pre-commit run --all-files           # regenerates docs/github-pages/static/helper.out
git diff --stat -- docs README.md .specify/memory/constitution.md
```

Expected: `helper.out` lists `--icons-dir` and the `svg` format; README and the pages under `docs/`
describe the format, the option, the legend and the limits; the constitution is at 1.1.0.
