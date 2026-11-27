---

description: "Task list for the SVG resource diagram feature"
---

# Tasks: SVG Resource Diagram

**Input**: Design documents from `specs/002-svg-resource-diagram/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: INCLUDED. The constitution (Principle III) requires tests in the same change: table-driven,
`t.Run`, success and error paths, a golden-pinned output per representative plan, `t.Helper()` and
`t.Cleanup()`. Within each story, write the tests first and see them fail.

**Organization**: grouped by user story so each story is implemented, tested and demonstrated on
its own. Story order is the spec priority order (US1 P1 ... US4 P4).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: US1 to US4, only inside the story phases
- Every task names the exact file(s). Paths follow `plan.md` (Source Code section).

## Conventions that apply to every task

- Go style per `.github/instructions/go.instructions.md`; run `gofmt` and `goimports` (local prefix
  `github.com/germainlefebvre4/tf2report`); wrap errors with `fmt.Errorf("...: %w", err)`; document
  exported symbols starting with the symbol name; comments in English; no `init()`.
- Standard library only. No new entry in `go.mod`.
- Keep each function under the gocyclo limit of 15 (split layout, render and sanitizer into small
  functions).
- Fixtures are synthetic: invented ids, names and values only.
- Output must be deterministic: no map iteration on any output path (sort first), integer
  coordinates, no timestamps or random ids.
- No attribute value (before, after, unknown, sensitive) may reach the `Scene` or the document.
- Commit per logical group with Conventional Commits: `feat:` code, `test:` tests and fixtures,
  `docs:` documentation and the constitution amendment. Never edit `CHANGELOG.md`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: new package skeleton and the governance change the rest depends on.

- [X] T001 Create `pkg/diagram/doc.go` with the package comment: the package builds and lays out a render-independent `Scene` from a `terraform.Plan`; it imports `pkg/terraform` only and contains no SVG or XML knowledge
- [X] T002 [P] Amend `.specify/memory/constitution.md` to version 1.1.0 (MINOR, `Last Amended: 2026-11-26`): in Principle II add `pkg/diagram` (selection, containers, edges, icon families, layout; imports `pkg/terraform` only) and extend the dependency flow with `pkg/report` -> `pkg/diagram` -> `pkg/terraform`; in Technology and Delivery Constraints change the output formats line to `markdown` (default), `text`, `json`, `svg`, change the Principle II sentence listing the formats to include `svg`, and amend "reads a plan file and writes a report" to also allow one optional, read-only, local icon directory (still no network); in Principle II reword the rule "adding a format MUST NOT require changing the parser" to "Adding a format MUST NOT change how existing plans are parsed or the output of existing formats; extending the parsed model with new optional fields is allowed and MUST keep every previously valid plan working", and add to the amendment a migration note: "specs and plans in flight that relied on the parser never changing (spec 001) remain valid; the baseline characterization tests keep pinning the existing parse behavior". Commit as `docs:` with the rationale from plan.md Complexity Tracking (three tracked deviations)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: address parsing, the Scene types and the fixtures that every story uses.

**CRITICAL**: no user story work starts before this phase is complete.

- [X] T003 [P] Create `pkg/terraform/address.go`: `ParseAddress(s string) Address` and `Address.ConfigAddress() string` per data-model.md (`Address{ Modules []ModuleStep{Name, Key}, Mode, Type, Name, Key }`). Tokenize `module.<name>[<key>]` repeated, then an optional `data.`, then `<type>.<name>[<key>]`; keys may be integers or quoted strings that contain dots or brackets (`["a.b"]`), so scan with quote awareness instead of splitting on `.`. Never panic and never return an error: a malformed address yields a best-effort value. `ConfigAddress` returns the full address with every `[...]` removed (`module.app[0].aws_instance.web["a"]` -> `module.app.aws_instance.web`)
- [X] T004 [P] Create `pkg/terraform/address_test.go`: table-driven `TestParseAddress` and `TestConfigAddress` covering a root resource, an integer index, a quoted `for_each` key containing dots, nested modules with keys on both levels, a data source, and malformed input (empty string, unbalanced bracket, trailing dot) that must not panic
- [X] T005 [P] Create `pkg/diagram/scene.go` with `Scene`, `Container`, `Node`, `Edge`, `Status` (fixed order create, update, delete, replace, noop) and `Options{ TypeFiltered, ActionFiltered bool }` exactly as in data-model.md. No field may hold an attribute value. Constraints to keep verbatim in the doc comments: "`Address` unique across nodes; a node exists only for create, update, delete, replace or (requested or context) no-op; never for data sources." and "no edge starts or ends at a `delete` node; no edge starts at a context node; edges towards nodes not drawn are dropped; duplicates removed."
- [X] T006 [P] Create fixture `testdata/plans/diagram/aws-network.json`: synthetic plan in the real `terraform show -json` shape with a `configuration` section; `aws_vpc.main` (no-op), `aws_subnet.private` (create, references `aws_vpc.main.id`), `aws_instance.web` (update, references `aws_subnet.private.id` and `aws_security_group.web.id`), `aws_security_group.web` (replace, references `aws_vpc.main.id`), `aws_s3_bucket.old` (delete), `aws_instance.worker[0]` and `[1]` (create, `count = 2`, referencing the subnet); every status present at least once
- [X] T007 [P] Create fixture `testdata/plans/diagram/azure-network.json`: `azurerm_resource_group.main`, `azurerm_virtual_network.main`, `azurerm_subnet.app`, `azurerm_linux_virtual_machine.app`, `azurerm_storage_account.logs` with references between them and at least create, update and delete present
- [X] T008 [P] Create fixture `testdata/plans/diagram/gcp-network.json`: `google_compute_network.main`, `google_compute_subnetwork.app`, `google_compute_instance.app`, `google_storage_bucket.data`, `google_sql_database_instance.db` with references between them and at least create, update and replace present
- [X] T009 [P] Create fixture `testdata/plans/diagram/modules-context.json`: a root `aws_instance.app` (create) referencing `module.network.vpc_id` (a module output), a module `module.network` holding an unchanged `aws_vpc.main` (no-op) and a changed `aws_subnet.a` (create) referencing the unchanged vpc, an `aws_route_table.rt` (no-op) referenced by nothing, a `module.app[0]` and `module.app[1]` pair (`count`), and a `random_pet.name` (create, provider `random`)
- [X] T010 [P] Create fixture `testdata/plans/diagram/cycle.json`: two resources `aws_security_group.a` and `aws_security_group.b` (both create) whose configuration references each other, plus one independent create
- [X] T011 [P] Create fixture `testdata/plans/diagram/no-configuration.json`: a valid plan with `resource_changes` (create, update, delete) and no `configuration` key at all, mirroring an older format version
- [X] T012 [P] Update `testdata/README.md` to list the new fixtures under `testdata/plans/diagram/` with one line each (what shape each one exercises)

**Checkpoint**: address parsing, Scene types and fixtures exist; stories can start.

---

## Phase 3: User Story 1 - See what a plan changes at a glance (Priority: P1) MVP

**Goal**: `--format svg` produces a valid diagram with provider and module containers and one node
per changed resource, with the status shown by color, symbol and border, a legend, and an empty
message. No arrows and no icons yet (the node reserves space for the icon).

**Independent Test**: run `./bin/tf2report -p testdata/plans/diagram/aws-network.json -f svg`,
open the file, and check containers, one node per changed resource, statuses distinguishable in
greyscale, no unchanged resource without `--action no-op`, identical bytes on a second run
(quickstart.md sections 1, 4, 5).

### Tests for User Story 1 (write first, they must fail)

- [X] T013 [P] [US1] Create `pkg/diagram/select_test.go`: table-driven selection matrix from research.md R-002 and R-010: (no filter) create/update/delete/replace drawn and no-op not; (type filter only) no-op not drawn; (action filter given) everything left in the plan is drawn including no-op; data sources and `read` never drawn; `replace` (delete + create) yields one node; an unknown action is treated as no-op
- [X] T014 [P] [US1] Create `pkg/diagram/build_test.go`: provider short name (`registry.terraform.io/hashicorp/aws` -> `aws`, `registry.opentofu.org/hashicorp/random` -> `random`, `terraform.io/builtin/terraform` -> `terraform`, empty name falls back to the type prefix before the first `_`); containers nest provider > module instance (`module.app[0]` and `module.app[1]` are two containers); root-module resources sit directly in the provider container; two resources with the same name in different modules never collide; `Address` is unique and verbatim; node ids follow address order; a module holding two providers appears under each provider
- [X] T015 [P] [US1] Create `pkg/diagram/layout_test.go`: grid fallback uses `ceil(sqrt(n))` columns filled in address order; container size is computed from its children (nodes, then child containers stacked below, plus padding and header); provider containers stack in alphabetical order; `Scene.Width` and `Height` bound all boxes without overlap; an empty selection sets `Scene.Empty`; two runs on the same input give identical coordinates
- [X] T016 [P] [US1] Create `pkg/report/svg_test.go`: output parses with `encoding/xml` and has a single `svg` root with `xmlns`, `viewBox`, `role="img"`, `<title>` and `<desc>` (counts per status); every status has a distinct color, symbol and `stroke-dasharray` (assert the three cues are pairwise distinct, including in greyscale: no two statuses share symbol and dash pattern); the legend lists only the statuses present; each node has a `<title>` with the full address; an empty selection produces a valid document with an explicit message and the legend; long names are truncated with an ellipsis; `NewFormatter(FormatSVG)` returns `*SVGFormatter`
- [X] T017 [P] [US1] Create `pkg/report/svg_safety_test.go`: `TestSVG_NoAttributeValues` loads `testdata/plans/terraform/sensitive.json` and `testdata/plans/diagram/aws-network.json`, collects every string value under `before`, `after`, `before_sensitive`, `after_sensitive` and `after_unknown` and asserts none appears in the output; `TestSVG_Escaping` builds a plan whose resource name, module key and provider contain `<`, `>`, `&`, `"`, `'` and a control character and asserts the output still parses as XML and contains no unescaped markup from those values
- [X] T018 [P] [US1] Create `pkg/report/svg_golden_test.go` reusing the helpers of `pkg/report/golden_test.go` (`-update` flag, `loadPlan`, `render`, `assertGolden`): one golden per fixture `aws-network`, `azure-network`, `gcp-network`, `no-configuration`, `testdata/plans/terraform/empty-config.json`, `testdata/plans/terraform/provider-random.json`, `testdata/plans/terraform/moved.json` and `testdata/plans/terraform/import.json` (a moved-only resource is hidden by default; an imported resource with an update is drawn once as an update), plus option variants (type filter, action filter, action filter with no-op) on `aws-network`; golden files `pkg/report/testdata/svg_*.golden`
- [X] T019 [P] [US1] Create `pkg/report/svg_scale_test.go`: `TestScale_300Resources` generates a 300-resource plan in memory (several providers, nested modules, mixed statuses), asserts the render finishes in under 2 seconds (spec SC-006) and that two renders are byte-identical (spec SC-003)
- [X] T020 [P] [US1] Create `cmd/tf2report/svg_test.go` (black-box, same `runCLI` helper as `main_test.go`): `-f svg` writes a valid SVG on stdout with empty stderr and exit 0; `output_format: svg` in a config file and `TF2REPORT_OUTPUT_FORMAT=svg` both select it, with flag > environment > file precedence; an empty-selection plan exits 0 with a valid document; `-f xml` exits 1 with `invalid format: xml (supported: markdown, text, json, svg)` on stderr and empty stdout; the existing three formats still produce their usual output for `basic.json`

### Implementation for User Story 1

- [X] T021 [US1] Create `pkg/diagram/select.go`: `Select(plan *terraform.Plan, opts Options) []terraform.ResourceChange` implementing the table of research.md R-002 and R-010 (no filter: the four changing actions; type filter only: same, no no-op; action filter given: everything left in the plan; data sources and `read` always dropped). Context is added later in `build.go` (US2)
- [X] T022 [US1] Create `pkg/diagram/build.go`: `Build(plan *terraform.Plan, opts Options) (*Scene, error)` using `Select` and `terraform.ParseAddress`: provider short name per T014, container tree per provider > module instance path (`Kind`, `Label` such as `module.app[0]`, `Path`), nodes (`Address`, `ConfigAddress`, `Type`, `Name`, `Provider`, `Status` from `GetPrimaryAction`), all slices sorted, `ID`s assigned from sorted order (`n1`, `n2`, ... and container ids), `Scene.Statuses` in fixed order, `Scene.Empty` when no node. Edges and context are added by US2
- [X] T023 [US1] Create `pkg/diagram/layout.go`: `Layout(*Scene)` with the grid fallback (`ceil(sqrt(n))` columns, address order), fixed node box 184 x 52, column gap 56, row gap 16, container padding 16, header 24, content-sized containers (nodes first, then child containers stacked beneath), provider containers stacked alphabetically, `Scene.Width`/`Height`; all `int`, constants in one block
- [X] T024 [P] [US1] Create `pkg/report/svg_style.go`: one status table (color tint and stroke, symbol, `stroke-dasharray`, stroke width, opacity) following research.md R-007: create green `+` solid 2 px; update orange `~` dashed `6 3`; delete red `-` dotted `2 3`; replace purple `-/+` dash-dot `8 3 2 3`; no-op grey, 60 % opacity, solid 1 px, no symbol; colors chosen so the label keeps a 4.5:1 contrast on its tint; the legend reads the same table
- [X] T025 [US1] Create `pkg/report/svg_render.go`: render a laid-out `Scene` into a `strings.Builder` and write once: root attributes, `<title>`/`<desc>`, legend group, container groups (`<g class="container provider|module">` with rect and label), node groups (`<title>` full address, rect with the status style, symbol, name on line 1, type on line 2, both truncated with an ellipsis using the 7 px per character estimate, space reserved at the left for the icon), empty-selection message. One escaping helper (`encoding/xml.EscapeText` plus removal of characters illegal in XML 1.0) used for every text and attribute value; numbers with `strconv.Itoa`
- [X] T026 [US1] Create `pkg/report/svg.go`: `SVGOptions` (`IconsDir`, `TypeFiltered`, `ActionFiltered`, `OnSkip`) and `SVGFormatter` implementing the unchanged `Formatter` interface (`Format(plan, summary, w)`: build, layout, render, write once; `summary` unused), plus `NewSVGFormatter(opts SVGOptions) *SVGFormatter`
- [X] T027 [US1] Edit `pkg/report/formatter.go`: add `FormatSVG Format = "svg"`, `ParseFormat("svg")` (no alias), `NewFormatter(FormatSVG)` returning `&SVGFormatter{}` with zero options, and change the invalid-format message to `invalid format: %s (supported: markdown, text, json, svg)`. In the same commit update the pinned expectations: `pkg/report/formatter_test.go` (the three `wantErr` rows and a `NewFormatter` case for `*SVGFormatter`), `cmd/tf2report/format_test.go` (the stderr assertion at the invalid-format check) and the row at line 56 of `specs/001-baseline-cli-behavior/contracts/cli-contract.md`
- [X] T028 [US1] Create `cmd/tf2report/formatter.go` with `newFormatter(format report.Format, cfg *config.Config) (report.Formatter, error)`: returns `report.NewSVGFormatter(report.SVGOptions{TypeFiltered: len(cfg.Filters.ResourceTypes) > 0, ActionFiltered: len(cfg.Filters.Actions) > 0})` for `FormatSVG` and `report.NewFormatter(format)` otherwise. Edit `cmd/tf2report/main.go`: call `newFormatter` instead of `report.NewFormatter`, add `svg` to the `--format` help text and `no-op` to the `--action` help text (keeps `run()` below the gocyclo limit; no other logic in `main.go`)

**Checkpoint**: `./bin/tf2report -p ... -f svg` works end to end without arrows or icons; T013 to T020 pass.

---

## Phase 4: User Story 2 - Understand how resources depend on each other (Priority: P2)

**Goal**: arrows between dependent resources, the layered layout, module-output arrows towards
containers, one-hop dependency context for unchanged resources, safe handling of cycles.

**Independent Test**: `./bin/tf2report -p testdata/plans/diagram/modules-context.json -f svg`
shows an arrow from each dependent resource to its dependency, the muted unchanged resource a
changed resource depends on, and no other unchanged resource; `cycle.json` and
`no-configuration.json` render without failure (quickstart.md sections 2 and 3).

### Tests for User Story 2 (write first, they must fail)

- [X] T029 [P] [US2] Create `pkg/terraform/configuration_test.go`: the plan parses `configuration` on `testdata/plans/terraform/nested-modules.json` and `mixed-actions.json` (module calls nested under `module_calls.<name>.module`), `Plan.Configuration` is `nil` for `testdata/plans/legacy-format.json` and `no-configuration.json` without error, `References()` finds references in nested expression blocks and arrays and in `depends_on`, returns them without duplicates in a stable order, and tolerates unexpected expression shapes (string, number, null) without panic
- [X] T030 [P] [US2] Create `pkg/terraform/reference_test.go`: table-driven `ParseReference` per research.md R-003: `aws_subnet.a.id` and `aws_subnet.a[0].id` -> resource; `module.net.vpc_id` and `module.net` -> module; `data.aws_ami.x.id` -> data; `var.x`, `local.x`, `count.index`, `each.key`, `path.module`, `terraform.workspace`, `self` and garbage -> other; never an error
- [X] T031 [P] [US2] Create `pkg/diagram/edges_test.go`: an edge links a dependent to its dependency from `references` and from `depends_on`; a repeated resource expands to every instance on each side (spec FR-007); a module reference yields an edge to each container of that module (spec FR-008); no edge starts or ends at a delete node (spec FR-009); duplicates are removed; edges towards nodes not drawn are dropped; with no filter a no-op resource referenced by a drawn changed node becomes a context node (`Context = true`) exactly once, with one edge per dependent, no edge starting at it, and its own dependencies not followed (spec FR-024, scenario 7); with a type filter or an action filter no context is added; no unchanged resource other than a referenced one is drawn
- [X] T032 [P] [US2] Create `pkg/diagram/layering_test.go`: a dependency sits in a lower column than its dependents; layers are ordered alphabetically then by two barycenter sweeps with address tie-breaks; a cycle terminates, still draws both edges and layers the nodes deterministically; a container with no intra-container edge keeps the grid; layering only uses edges whose two ends are direct nodes of the same container; two runs are identical; `routing` produces curves from the left side of the dependent to the right side of the dependency (or the container edge for a module target), with integer points
- [X] T033 [P] [US2] Create `pkg/report/svg_edges_test.go`: arrow marker defined once in `<defs>`; one `<path>` per edge inside `<g class="edges">` placed after containers and before nodes; context nodes carry the `context` class and the muted style and show no status symbol; no arrow in `no-configuration.json`; the cycle fixture renders without failure

### Implementation for User Story 2

- [X] T034 [US2] Edit `pkg/terraform/types.go` to add `Configuration *Configuration` with tag `json:"configuration,omitempty"` to `Plan`, and create `pkg/terraform/configuration.go` with `Configuration{RootModule ConfigModule}`, `ConfigModule{Resources []ConfigResource; ModuleCalls map[string]ModuleCall}`, `ModuleCall{Module ConfigModule}`, `ConfigResource{Address, Mode, Type, Name string; DependsOn []string; Expressions map[string]any}` and `(ConfigResource).References() []string` (recursive walk of `Expressions` collecting every string of any `references` array, plus `DependsOn`, de-duplicated, stable order). Confirm the existing JSON, text and markdown goldens are unchanged
- [X] T035 [P] [US2] Create `pkg/terraform/reference.go` with `ParseReference(s string) Reference{Kind, ModuleCall, Type, Name}` and `Kind` values `resource`, `module`, `data`, `other` per research.md R-003; strip `[...]` and attribute suffixes; never fail
- [X] T036 [US2] Create `pkg/diagram/edges.go`: walk `Plan.Configuration` module by module (relative references resolve inside the module that holds the resource, `module.X` resolves to the child module call), map each config resource to its nodes through `ConfigAddress`, and expand edges to every instance on both sides; a module reference targets every container whose index-free path matches; skip `data`/`other` references; drop edges that start or end at a delete node, that start at a context node, that point to nodes not drawn, and duplicates; sort edges by (from, to). Plan without `Configuration` yields no edges
- [X] T037 [US2] Edit `pkg/diagram/build.go`: after nodes and containers, call the edge resolution over all candidate resources of the plan (including unchanged ones); when `!opts.TypeFiltered && !opts.ActionFiltered`, add each unchanged node that is the direct target of an edge from a drawn changed node as a context node (`Context = true`, in its own provider and module containers), once; then drop edges whose ends are not drawn; keep `Scene.Edges` sorted
- [X] T038 [US2] Create `pkg/diagram/layering.go`: per container, longest-path layering over intra-container edges (a dependency in a lower column than its dependents, arrows point left); cycle breaking by a DFS over nodes sorted by address, dropping back edges from layering only; ordering inside a layer alphabetical then two barycenter sweeps with stable sort and address tie-break; O(nodes + edges). Edit `pkg/diagram/layout.go` to use layering when a container has at least one intra-container edge and the grid otherwise
- [X] T039 [US2] Create `pkg/diagram/routing.go`: for each edge compute integer cubic Bezier control points with horizontal tangents from the left side of the dependent node to the right side of the dependency node (or to the right edge of the target container for a module target); deterministic; called at the end of `Layout`
- [X] T040 [US2] Edit `pkg/report/svg_render.go`: define the arrow `<marker>` once in `<defs>`, emit `<g class="edges">` with one `<path>` per edge between the containers and the nodes, and render context nodes with the `context` class, the muted style from `svg_style.go` and no status symbol
- [X] T041 [US2] Regenerate the goldens with the existing `-update` flag now that arrows and the layered layout exist (`go test ./pkg/report -run Golden -update`), review every golden diff by hand (US1 goldens change on purpose), add goldens for `modules-context`, `cycle` and `no-configuration`, and re-run without `-update`

**Checkpoint**: arrows, context and cycles work; T029 to T033 pass; US1 still passes.

---

## Phase 5: User Story 3 - Recognize resources by their cloud service icon (Priority: P3)

**Goal**: an icon per resource (service family, then provider family, then category, then generic), embedded generic
glyphs, an optional user icon pack through `--icons-dir`, each icon defined once in the document,
a strict sanitizer for pack files.

**Independent Test**: `--format svg` on `aws-network`, `azure-network`, `gcp-network` and
`provider-random` shows service, category or generic icons; with `--icons-dir testdata/icons/valid`
matching families use the pack; `testdata/icons/hostile` produces a clean document; a missing
directory fails before any output (quickstart.md section 6).

### Pack fixtures (needed by the tests below)

- [X] T042 [P] [US3] Create `testdata/icons/valid/` with small synthetic icon files that pass the sanitizer: `aws-ec2.svg`, `aws.svg`, `azure-vm.svg`, `compute.svg` and `generic.svg` (original shapes, `viewBox`, only allowed elements), plus `notes.txt` and a subdirectory `nested/aws-vpc.svg` that must be ignored (non-recursive, `*.svg` only)
- [X] T043 [P] [US3] Create `testdata/icons/hostile/` with one file per rejected construct, named after it: `script.svg` (`<script>`), `onload.svg` (`onload` attribute), `foreign.svg` (`<foreignObject>`), `image.svg` (`<image>`), `style.svg` (`<style>` and a `style` attribute with `url(` and `@import`), `use-external.svg` (`<use href="http://...">`), `href-external.svg` (`xlink:href` not starting with `#`), `doctype.svg` (DOCTYPE with an entity), `not-svg.svg` (valid XML with another root), `malformed.svg` (unclosed tag); the oversize case (more than 256 KiB) is generated inside the test with `t.TempDir()`, not committed

### Tests for User Story 3 (write first, they must fail)

- [X] T044 [P] [US3] Create `pkg/diagram/family_test.go`: service table hits (`aws_instance` -> chain `aws-ec2`, `aws`, `compute`, `generic`; `azurerm_linux_virtual_machine`; `google_compute_instance`) for the first 25 services of each cloud; the category heuristic of research.md R-006 with its priority order (security, database, network, storage, compute) on unknown types such as `aws_foo_bucket` (chain `aws`, `storage`, `generic`) and `google_foo_firewall` (chain `gcp`, `security`, `generic`); category keywords match whole `_`-separated tokens, never substrings (`aws_monkey_thing` and `azurerm_albumin` match no keyword and end on `aws`/`azure` then `generic`; `aws_lb`, `aws_alb` and `aws_elb` are resolved by the service table); a provider without a table (`random_pet`, `vault_foo`) has no provider level and ends on `generic`; the chain is never empty and always ends with `generic`; a type that matches several categories takes the first in priority order
- [X] T045 [P] [US3] Create `pkg/report/svg_icons_test.go`: the six built-in glyphs (`compute`, `network`, `storage`, `database`, `security`, `generic`) exist, pass the sanitizer and parse as XML; pack loading is non-recursive and only reads `*.svg`; for each level of a chain the pack wins over the built-in glyph and the first hit wins; an unknown `aws_foo` type uses the pack's `aws.svg` before the category glyph and falls back to the built-in category glyph when the pack has no `aws.svg` (the provider level has no built-in glyph); every file of `testdata/icons/hostile` is rejected with a reason and reported through `OnSkip`; a file over 256 KiB (generated in `t.TempDir()`) is rejected; accepted content is re-serialized (not copied byte for byte) with ids prefixed by the family id; a missing directory, a path that is a file and an unreadable directory each return an error that names the directory; a directory with no usable file is not an error
- [X] T046 [P] [US3] Create `pkg/report/svg_icons_doc_test.go`: for a plan with many resources of the same family the icon appears once as `<symbol id="icon-...">` in `<defs>` and every node uses `<use href="#icon-...">` (spec FR-015, scenario 4); the document grows by less than one symbol when the node count doubles; with `testdata/icons/hostile` the output contains no `<script`, `on...=`, `foreignObject` or non-fragment `href` (case-insensitive regexp); an unknown type still gets a valid node with the generic icon; golden `pkg/report/testdata/svg_aws_network_pack.golden` for `aws-network` with `testdata/icons/valid`
- [X] T047 [P] [US3] Create `cmd/tf2report/svg_icons_test.go`: `-f svg --icons-dir testdata/icons/valid` exits 0 and uses pack icons; `--icons-dir /nonexistent` exits non-zero with an error naming the directory and empty stdout (no partial document); `-f json --icons-dir /nonexistent` exits 0 with the normal JSON, empty stderr; the same with `--verbose` prints exactly one notice line saying `--icons-dir` is ignored for this format; `-f svg -v --icons-dir testdata/icons/hostile` lists each skipped file on stderr and still exits 0

### Implementation for User Story 3

- [X] T048 [P] [US3] Create `pkg/diagram/family.go`: `Family(resourceType, provider string) []string` returning the ordered chain (service family, provider family `aws`/`azure`/`gcp` when the type prefix is `aws_`/`azurerm_`/`google_`, category, `generic`). Service table of at least 25 entries per cloud in `<cloud>-<service>` ids (AWS e.g. `aws_instance` -> `aws-ec2`, `aws_vpc` -> `aws-vpc`, `aws_subnet`, `aws_security_group`, `aws_s3_bucket`, `aws_db_instance`, `aws_lambda_function`, `aws_lb`, `aws_iam_role`, `aws_eks_cluster`; Azure e.g. `azurerm_linux_virtual_machine` -> `azure-vm`, `azurerm_virtual_network`, `azurerm_subnet`, `azurerm_storage_account`, `azurerm_key_vault`, `azurerm_kubernetes_cluster`; GCP e.g. `google_compute_instance` -> `gcp-gce`, `google_compute_network`, `google_compute_subnetwork`, `google_storage_bucket`, `google_sql_database_instance`, `google_container_cluster`); the category heuristic of research.md R-006 in its fixed priority order on the `_`-separated tokens of the type name after the provider prefix, matching whole tokens and never substrings (`aws_lb` matches `lb`, `aws_monkey_thing` does not match `key`), with `aws_lb`, `aws_alb` and `aws_elb` listed explicitly in the service table; tables are read-only package-level values, no mutable globals
- [X] T049 [US3] Edit `pkg/diagram/build.go` to set `Node.Family` from `Family(type, provider)` for every node, including context nodes
- [X] T050 [P] [US3] Create the six glyphs `pkg/report/icons/compute.svg`, `network.svg`, `storage.svg`, `database.svg`, `security.svg`, `generic.svg`: original monochrome drawings authored from scratch (no third-party or provider artwork), `viewBox="0 0 24 24"`, `currentColor`, only elements allowed by the sanitizer
- [X] T051 [US3] Create `pkg/report/svg_icons.go`: `//go:embed icons/*.svg`; `loadPack(dir string, onSkip func(file, reason string)) (*iconPack, error)` reading `*.svg` non-recursively; the sanitizer per research.md R-008 and the contract: at most 256 KiB per file, root `svg`, allowed elements only `svg g path rect circle ellipse line polyline polygon defs linearGradient radialGradient stop clipPath title desc`, attributes limited to a presentation allowlist, reject any `on*` attribute, `<style>`, `<image>`, `<use>`, `<script>`, `<foreignObject>`, DOCTYPE and entity declarations, `style` containing `url(` or `@import`, and any `href`/`xlink:href` that is not a local `#fragment`; re-serialize the accepted content into a `<symbol>` with ids prefixed by the family id; `resolve(chain []string)` tries the pack then the built-in glyph at each level (built-in glyphs exist only for the five categories and `generic`, so the service and provider levels, including `aws`, `azure` and `gcp`, are pack-only); an unreadable or missing directory returns an error naming it
- [X] T052 [US3] Edit `pkg/report/svg_render.go` and `pkg/report/svg.go`: `Format` loads the pack (when `IconsDir` is set) before producing any output so that an invalid directory leaves the writer untouched; collect the families used, emit each icon once as a `<symbol id="icon-...">` in `<defs>`, and draw `<use href="#icon-...">` at the left of each node in the reserved space; pass `OnSkip` through
- [X] T053 [US3] Edit `cmd/tf2report/main.go` and `cmd/tf2report/formatter.go`: add the `--icons-dir` flag (string, no short form, no config key, no environment variable), pass `IconsDir` and an `OnSkip` that writes `skipped icon <file>: <reason>` to stderr only when `--verbose` is set; when the format is not `svg` and `--icons-dir` is given, ignore it and, with `--verbose`, print one notice line to stderr; keep the new logic in `formatter.go` so `run()` stays under the gocyclo limit

**Checkpoint**: icons and packs work; T044 to T047 pass; US1 and US2 still pass.

---

## Phase 6: User Story 4 - Focus and share the diagram safely (Priority: P4)

**Goal**: the existing filters behave exactly as the contract table says on diagrams with arrows
and context, and no attribute value ever reaches an image, whatever the options.

**Independent Test**: quickstart.md sections 2 and 5: type and action filters with the `modules-context`
and `sensitive` fixtures, and a search of the output for attribute values.

### Tests for User Story 4 (write first)

- [X] T054 [P] [US4] Create `cmd/tf2report/svg_filters_test.go`: the filters x statuses table of contracts/cli-contract.md on `modules-context.json` and `aws-network.json`: no filter (changed plus one-hop muted context), `--type aws_subnet` (only subnets, no muted node, no dangling arrow), `--action create` (only creates, no context, no arrow to a node not drawn), `--action no-op` and `--action create,no-op` (unchanged resources drawn muted), `--type` together with `--action`; assert node counts, absence of arrows to undrawn nodes and exit code 0
- [X] T055 [P] [US4] Create `pkg/report/svg_safety_matrix_test.go`: for every fixture under `testdata/plans/diagram/` and the sensitive fixtures `testdata/plans/terraform/sensitive.json`, `provider-sensitive.json`, `testdata/plans/opentofu/sensitive.json`, `testdata/plans/terraform/moved.json` and `testdata/plans/terraform/import.json`, and for the option sets (no filter, type filter, action filter, action with no-op, valid pack, hostile pack), assert the output is well-formed XML, contains no attribute value (spec FR-017, SC-005) and is identical across two runs (SC-003)
- [X] T056 [P] [US4] Edit `pkg/report/svg_scale_test.go`: add `TestScale_300ResourcesWithEdges` with `count`-expanded resources referencing each other (the cross-product case of research.md R-014) staying under 2 seconds and deterministic; add a table case with a type filter and with an action filter over the same generated plan

### Implementation for User Story 4

- [X] T057 [US4] Run T054 to T056 and fix every gap they reveal in `pkg/diagram/select.go`, `pkg/diagram/build.go`, `pkg/diagram/edges.go` or `pkg/report/svg_render.go` (for example an arrow left towards a resource removed by a filter, or context leaking under a filter). If no gap exists, record "no gap found" as the task note and change nothing

**Checkpoint**: all four stories pass independently and together.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: documentation, generated help, examples and the final gates (Constitution IV).

- [X] T058 [P] Update `README.md`: add `svg` to the features and to `--format` (line listing markdown, text, json), document `--icons-dir`, the status legend (color + symbol + border), the filters behavior table, the limits (no isometric, no PNG, module outputs point to containers, deleted resources have no arrows) and a `--format svg` example command; no example output embedded
- [X] T059 [P] Update `docs/USAGE.md`, `docs/OUTPUT_FORMATS.md` and `docs/CONFIGURATION.md` (the `output_format` values list and `TF2REPORT_OUTPUT_FORMAT=svg`, no new config key) for the `svg` format, `--icons-dir` and `--action no-op`, with the icon pack contract (file names, allowed content, fallback)
- [X] T060 [P] Update `docs/API.md` and `docs/DEVELOPMENT.md`: `report.FormatSVG`, `report.NewSVGFormatter`, `report.SVGOptions`, the new `pkg/diagram` package and its role, the `Plan.Configuration` field, and the architecture note that renderers consume the `Scene`
- [X] T061 [P] Update the Docusaurus pages `docs/github-pages/docs/output-formats.md`, `docs/github-pages/docs/usage.md`, `docs/github-pages/docs/configuration.md`, `docs/github-pages/docs/filtering.md`, `docs/github-pages/docs/reference/cli-options.md`, `docs/github-pages/docs/reference/configuration-schema.md` and `docs/github-pages/docs/development/api-reference.md` with the same content as T058 to T060
- [X] T062 Regenerate the help output with `go run ./cmd/tf2report/main.go -h > docs/github-pages/static/helper.out` (the `cli-helper` pre-commit hook does the same) and copy the new flags block into `docs/dockerhub/README.md` (it holds a copy at the `-f, --format` line); commit the result
- [X] T063 [P] Copy `testdata/plans/diagram/aws-network.json` byte for byte to `examples/svg-network.json`, add `"svg-network.json": "diagram/aws-network.json"` to `exampleSources` in `cmd/tf2report/examples_test.go` (mandatory: `TestExamples` compares the number of files in `examples/` with `exampleSources` and each file with its source fixture, and runs the three existing formats on it), and add a line for the file in `examples/README.md`. Do not embed any example output anywhere
- [X] T064 [P] Edit `tf2report.yaml.example` and `specs/001-baseline-cli-behavior/contracts/cli-contract.md`: list `svg` among the `output_format` values in the sample (comment line only), and add to the "Provisional" note on `--action` that `no-op` is accepted and matches unchanged resources (the existing formats exclude them, so their output is unchanged)
- [X] T065 Run the gates and fix findings: `make fmt && go vet ./...`, `gofmt -l .` is empty, `goimports -l -local github.com/germainlefebvre4/tf2report .` is empty, `make test`, `make test-coverage` (every new package has tests), `pre-commit run --all-files`; confirm `go.mod` and `go.sum` are unchanged (no new dependency) and that all existing goldens and example tests passed without edits other than the pinned format-list message of T027 (spec SC-008)
- [ ] T066 Run every command of `quickstart.md` (sections 1 to 9) against the built binary and compare with the stated expectations; open the SVGs of `aws-network`, `modules-context` and `cycle` in a browser and check legibility (containers, arrows, muted context, legend); record deviations as follow-up tasks rather than editing the quickstart to match (NOTE: CLI commands run and pass; browser legibility check not done, Chrome tool cannot open file:// URLs)
- [ ] T067 Check the planned commit split and titles: `docs: amend constitution to 1.1.0 (svg format, pkg/diagram)`, `test: add diagram fixtures and icon pack samples`, `feat: add svg diagram format`, `feat: draw dependencies and context in svg diagram`, `feat: add icons and --icons-dir`, `docs: document the svg format`; none is breaking, so none carries a breaking marker

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependency. T002 (constitution) can run at any time before merge but should land first.
- **Foundational (Phase 2)**: depends on T001 for T005; T003 to T012 are otherwise independent. Blocks all stories.
- **US1 (Phase 3)**: depends on Phase 2. MVP.
- **US2 (Phase 4)**: depends on Phase 2 and on the US1 code (it edits `build.go`, `layout.go`, `svg_render.go`); its tests (T029 to T033) can be written in parallel with the end of US1.
- **US3 (Phase 5)**: depends on Phase 2 and on US1; it edits `build.go` and `svg_render.go` after US2, so run it after US2 or merge carefully. Its pure parts (`family.go`, glyphs, `svg_icons.go`, fixtures, their tests) are independent and can run in parallel with US2.
- **US4 (Phase 6)**: depends on US1 and US2 (arrows and context) and on US3 for the pack cases of T055.
- **Polish (Phase 7)**: after the stories you ship. T062 needs the final flag set (after T053).

### User Story Dependencies

- **US1**: independent after Foundational.
- **US2**: builds on the US1 Scene and renderer; testable alone with `modules-context.json`.
- **US3**: builds on the US1 renderer; testable alone with the pack fixtures.
- **US4**: verification story over US1 to US3.

### Within Each User Story

- Tests first, and they must fail; then data types and selection, then build, layout, render, then CLI wiring.
- Same-file edits are sequential: `pkg/diagram/build.go` (T022, T037, T049), `pkg/report/svg_render.go` (T025, T040, T052), `pkg/diagram/layout.go` (T023, T038), `pkg/report/svg.go` (T026, T052), `cmd/tf2report/main.go` and `formatter.go` (T028, T053).

### Parallel Opportunities

- Phase 2: T003 to T012 are all [P] (T004 after T003).
- US1 tests: T013 to T020 are all [P]; T024 (`svg_style.go`) is [P] with T021 to T023.
- US2 tests: T029 to T033 are all [P]; T035 is [P] with T034.
- US3: T042, T043, T044, T045, T046, T047, T048, T050 can run in parallel; T051 then T052 and T053 are sequential.
- US4 tests: T054 to T056 are [P].
- Polish: T058 to T061, T063 and T064 are [P].

### Parallel Example: User Story 1

```text
# Launch all tests for User Story 1 together:
Task: "pkg/diagram/select_test.go"        (T013)
Task: "pkg/diagram/build_test.go"         (T014)
Task: "pkg/diagram/layout_test.go"        (T015)
Task: "pkg/report/svg_test.go"            (T016)
Task: "pkg/report/svg_safety_test.go"     (T017)
Task: "pkg/report/svg_golden_test.go"     (T018)

# Then in parallel with the first implementation tasks:
Task: "pkg/report/svg_style.go"           (T024) alongside select.go / build.go / layout.go
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Phase 1 and Phase 2 (constitution, address parsing, Scene types, fixtures).
2. Phase 3 (US1): `--format svg` draws containers and statused nodes.
3. **Stop and validate**: quickstart.md sections 1, 4 and 5; run `go test ./...`.
4. Demo: the diagram is already useful as a grouped picture of the changes.

### Incremental Delivery

1. Foundation -> US1 (MVP, grouped picture, no arrows or icons).
2. US2: arrows, layered layout and dependency context (the architecture view).
3. US3: icons and the optional icon pack.
4. US4: filter and safety verification across everything.
5. Polish: documentation, help output, examples, gates.

Each story adds value without breaking the previous ones; the existing three formats stay
byte-identical throughout (spec FR-022, SC-008).

### Out of scope (spec Assumptions)

Isometric rendering, PNG export, resolving module outputs to the producing resource, recovering the
relations of resources to delete, committing official AWS, Azure or GCP artwork (license check
required first), and a configuration key or environment variable for `--icons-dir`.

---

## Notes

- Spec traceability: US1 -> FR-001 to FR-005, FR-011 (grid part), FR-017, FR-018, FR-020, FR-022;
  US2 -> FR-006 to FR-011, FR-024; US3 -> FR-012 to FR-016; US4 -> FR-017, FR-019; Polish -> FR-021 to
  FR-023.
- Success criteria: SC-001, SC-002, SC-007 are checked by tests T016, T054, T066; SC-003 by
  T016/T019/T055; SC-004 by T055 and the fixture set; SC-005 by T017 and T055; SC-006 by T019 and
  T056; SC-008 by T065.
- `-update` refreshes goldens only on purpose and every golden diff is reviewed by hand.
- Stop at any checkpoint to validate a story independently.

---

## Phase 8: Convergence

**Purpose**: remaining work found by `/speckit-converge` on 2026-11-26. The functional requirements are
satisfied (build, vet, gofmt and the 622 tests pass, quickstart sections 1 to 6 behave as specified);
what remains is code quality against Constitution IV on the new code, and the quickstart commands.

- [ ] T068 CRITICAL: split `sanitizeIcon` in `pkg/report/svg_icons.go` (cyclomatic complexity 21, limit 15) into small helpers (element allowlist check, attribute allowlist check, `style`/`href` check, re-serialization) with unchanged behavior, and reduce `TestBuild_ContainersNestProviderThenModule` in `pkg/diagram/build_test.go` (complexity 16) and `TestEdges_ContextNodes` in `pkg/diagram/edges_test.go` (complexity 17) to 15 or less with subtests or helpers; `golangci-lint` gocyclo must report nothing on these files per Constitution IV (contradicts)
- [ ] T069 CRITICAL: wrap the two lines longer than 140 characters in `pkg/report/svg_render.go` (line 116 with 142 characters, line 144 with 150 characters) without changing the generated SVG bytes (the goldens in `pkg/report/testdata/svg_*.golden` must stay identical) per Constitution IV (contradicts)
- [ ] T070 CRITICAL: fix the remaining lint findings on new code per Constitution IV (contradicts): handle the unchecked errors of `fmt.Fprintln` and `fmt.Fprintf` at `cmd/tf2report/formatter.go` lines 18 and 29 (assign to `_, _` with a one-line reason, since a failed write to standard error has no recovery), handle `defer f.Close()` at `pkg/report/svg_icons.go:151` (check or ignore the error with an explained `//nolint:errcheck`, read-only file), rename or drop the local `max` that redefines the builtin at `pkg/report/svg_render.go` lines 217 and 219 (Go 1.25 has `max`), and apply the De Morgan simplification flagged at `pkg/report/svg_icons.go:283`
- [ ] T071 Fix the two quickstart commands that hit the baseline Known Gap of spec 001 (a plan with `before_sensitive: false` fails to parse) in `specs/002-svg-resource-diagram/quickstart.md`: in section 3 replace `testdata/plans/terraform/provider-random.json` with `testdata/plans/diagram/modules-context.json` (it holds a `random_pet` resource, which must get the generic icon), and in section 5 replace `testdata/plans/terraform/sensitive.json` with `testdata/plans/sensitive.json`; re-run both commands and confirm exit 0 and, for section 5, that no `SYNTHETIC-SECRET` value appears in the output per plan: quickstart.md (partial)
- [ ] T072 Resolve the `goconst` findings of the repository configuration (min length 2, min occurrences 2) on new code: the repeated icon family ids in `pkg/diagram/family.go`, and the repeated `"data"` mode string in `pkg/diagram/edges.go`, `pkg/diagram/select.go` and `pkg/terraform/address.go` (an exported constant `RefData` already exists in `pkg/terraform`); use named constants where it reads better, and where a data table is clearer as literals, add a `//nolint:goconst` with an explanation instead (the same finding already exists on unchanged code such as `pkg/terraform/types.go`) per Constitution IV (partial)
- [ ] T073 Record in the pull request description that `Makefile` (`build` target), `.goreleaser.yaml` (`main:`) and `.pre-commit-config.yaml` (`cli-helper` entry) now build and run `./cmd/tf2report` instead of `cmd/tf2report/main.go`, because `cmd/tf2report` gained `formatter.go` and a single-file build no longer compiles; keep the three edits and check that `make build`, `go run ./cmd/tf2report -h` and `goreleaser check` (`make goreleaser-check`) still work per plan: Source Code structure (unrequested)

**Checkpoint**: after T068 to T073, run `go test ./...`, `gofmt -l .`, the lint command used by this phase and `quickstart.md` sections 3 and 5; then T066 and T067 above can be checked.
