# Implementation Plan: SVG Resource Diagram

**Branch**: `feat/show-visual-diagram` (spec directory `002-svg-resource-diagram`) | **Date**: 2026-11-26 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/002-svg-resource-diagram/spec.md`

## Summary

Add a fourth output format, `svg`, that draws the resources of a plan as a 2D architecture
diagram: provider and module containers, one node per resource with a status (color, symbol,
border), arrows for dependencies, and an icon per resource type.

The approach separates three concerns so that a later isometric renderer reuses the first two:

1. **Read the relations** (`pkg/terraform`): parse the plan `configuration` section (resource
   expressions' references, `depends_on`, module calls) and expose address helpers. The parser
   stays tolerant: a plan without `configuration` still parses.
2. **Build and lay out a render-independent `Scene`** (`pkg/diagram`, new): select the resources
   to draw, group them in containers, resolve edges, assign an icon family, and compute integer
   coordinates with a deterministic layered layout. No XML or SVG knowledge.
3. **Render the Scene as SVG** (`pkg/report`): a new `SVGFormatter` writes the document with the
   standard library only, embeds generic glyphs with `go:embed`, and loads an optional
   user icon pack through a strict allowlist sanitizer.

`cmd/tf2report` gains the `--icons-dir` flag and passes filter knowledge to the formatter, because
the dependency context and no-op rules (spec FR-003, FR-019, FR-024) depend on which filters the
user gave. The `Formatter` interface and the three existing formats do not change (FR-022).

No new third-party dependency: `encoding/xml`, `embed`, `sort`, `strings`, `strconv`,
`path/filepath` and `io/fs` only.

## Technical Context

**Language/Version**: Go 1.25 (`go.mod`).

**Primary Dependencies**: none added. Existing: `cobra`, `viper`.

**Storage**: files only. Plan JSON in (plus an optional icon directory), SVG on standard output.
Generic glyphs are embedded in the binary.

**Testing**: `go test ./...`. Table-driven tests with `t.Run`; golden files for the SVG output
(refreshed only with the existing `-update` flag convention); well-formedness checked by parsing
the output with `encoding/xml`; a safety test searching the output for every attribute value of
the sensitive fixtures; layout unit tests on the `Scene` (no SVG involved).

**Target Platform**: Linux and macOS (shipped platforms).

**Project Type**: single-binary CLI (`cmd/` + `pkg/`).

**Performance Goals**: 300 resources rendered in under 2 s (spec SC-006). Layering and edge
resolution are O(nodes + edges); ordering sweeps are bounded (two passes). A benchmark-style test
with a generated 300-resource plan guards it.

**Constraints**: offline and read-only; deterministic output (integers only, sorted iteration, no
map ranging in output paths); no attribute value in the image; text escaped; user SVG icons
sanitized by allowlist; no official AWS/Azure/GCP artwork committed (spec FR-013); gocyclo limit
15 (layout and renderer are split into small functions).

**Scale/Scope**: 4 user stories, 24 functional requirements. One new package (`pkg/diagram`)
and three extended ones (`pkg/terraform`, `pkg/report`, `cmd/tf2report`), 6 synthetic fixture
plans, one embedded glyph set of 6 icons, about 75 icon-family mappings (about 25 per cloud) as a
first table.

All Technical Context items are resolved. Open points from the spec (icon pack naming, orientation,
node labels) are decided in `research.md` (R-006, R-005, R-007).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design (below).*

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. Idiomatic Go and Simplicity | Pass | Stdlib only, no `init()`, no mutable globals (the embedded glyph `embed.FS` and the family tables are read-only). Small interfaces defined where consumed. Exported symbols documented. |
| II. Layered Package Architecture | Deviation, tracked | Two points. (1) A new package `pkg/diagram` is added (imports `pkg/terraform` only; `pkg/report` imports it): the constitution lists the allowed package flow. (2) `pkg/terraform` is extended (optional `Plan.Configuration`, address and reference helpers) because the `svg` format needs the relations of `configuration`, which the parser ignored: the constitution says adding a format MUST NOT require changing the parser. Both are resolved by amending the constitution in the same change (see Complexity Tracking); the extension is additive and every previously valid plan keeps parsing and formatting as before. `cmd` still holds wiring only: the format selection moves to a small helper file, parsing and drawing stay in `pkg/*`. The `Formatter` interface is unchanged and the format is still selected through `NewFormatter`. |
| III. Test-Backed Changes | Pass | Every new package ships table-driven tests; golden-pinned SVG per representative plan; success and error paths (empty plan, no configuration, cycles, unknown types, hostile icon files, invalid icon directory); `t.Helper()`/`t.Cleanup()`. |
| IV. Automated Quality Gates | Pass with inherited deviation | `--icons-dir` is a CLI change: `docs/github-pages/static/helper.out`, the README and `docs/` are regenerated and updated. The `golangci-lint` v1/v2 configuration issue recorded in spec 001 is unchanged; the interim gate (`gofmt`, `goimports`, `go vet`, pre-commit) applies. |
| V. Conventional Commits and Automated Releases | Pass | `feat:` for the format and the flag, `docs:` for the constitution amendment and documentation, `test:` for fixtures. The change is additive (new format value, new flag), not breaking. `CHANGELOG.md` untouched. |

Technology and Delivery Constraints:

- Output formats: the list `markdown`, `text`, `json` MUST gain `svg` (amendment).
- "Reads a plan file and writes a report": the optional icon directory is a second, read-only,
  local input; the sentence is amended to allow it. No network access is added.
- Sensitive values: the renderer receives only addresses, types, names, providers and statuses
  (the `Scene` carries no attribute value by construction), which is what FR-017 requires.
- Fixtures: all new fixtures are synthetic and hand-written. Examples: one input plan added under
  `examples/`; example outputs are not embedded anywhere (golden files under `testdata/` are test
  data, not documentation).

**Gate result: PASS** (one tracked deviation to II, with its amendment part of the change).

## Project Structure

### Documentation (this feature)

```text
specs/002-svg-resource-diagram/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-contract.md  # Phase 1 output: CLI, SVG document and icon pack contracts
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 output (/speckit-tasks command, not created here)
```

### Source Code (repository root)

```text
cmd/tf2report/
├── main.go                  # EDIT: --icons-dir flag, call newFormatter, notice in verbose mode
├── formatter.go             # NEW: newFormatter(format, options) wiring helper
├── svg_test.go              # NEW: black-box tests (flags, filters x statuses, streams, exit codes)

pkg/terraform/
├── types.go                 # EDIT: Plan.Configuration (optional, omitted when absent)
├── configuration.go         # NEW: Configuration, ConfigModule, ConfigResource, ModuleCall,
│                            #      reference extraction (walks expressions, depends_on)
├── address.go               # NEW: ParseAddress (module path with keys, type, name, key),
│                            #      Address.ConfigAddress (index-free)
├── reference.go             # NEW: ParseReference (resource, module, data, other)
├── configuration_test.go, address_test.go, reference_test.go   # NEW

pkg/diagram/                 # NEW package, imports pkg/terraform only
├── doc.go                   # package documentation (no rendering knowledge)
├── scene.go                 # Scene, Container, Node, Edge, Status, Family
├── select.go                # which resources are drawn (filters, no-op)
├── build.go                 # Build: containers and nodes from Plan, then edges and context
├── edges.go                 # edge resolution from configuration, instance expansion, context
├── family.go                # type -> icon family (table, prefix, category, generic)
├── layout.go                # Layout: grid fallback, container sizing, stacking
├── layering.go              # longest-path layering, cycle breaking, barycenter ordering
├── routing.go               # edge curves between nodes and containers
├── select_test.go, build_test.go, edges_test.go, family_test.go, layout_test.go,
│   layering_test.go

pkg/report/
├── formatter.go             # EDIT: FormatSVG, ParseFormat ("svg"), NewFormatter(FormatSVG)
├── svg.go                   # NEW: SVGFormatter, SVGOptions, NewSVGFormatter
├── svg_render.go            # NEW: Scene -> SVG document (escape, styles, edges, legend)
├── svg_style.go             # NEW: status palette, symbols, border styles, metrics
├── svg_icons.go             # NEW: embedded glyphs, icon pack loader, allowlist sanitizer
├── icons/*.svg              # NEW: 6 original generic glyphs (compute, network, storage,
│                            #      database, security, generic), embedded with go:embed
├── svg_test.go, svg_icons_test.go, svg_safety_test.go, svg_golden_test.go
└── testdata/svg_*.golden    # NEW: golden files

testdata/plans/diagram/      # NEW synthetic fixtures with a configuration section
├── aws-network.json, azure-network.json, gcp-network.json,
├── modules-context.json, cycle.json, no-configuration.json
testdata/icons/              # NEW: valid and hostile sample icon packs for tests
examples/svg-network.json    # NEW: example input plan (no example output embedded)

docs/, README.md, docs/github-pages/docs/*   # EDIT: svg format, --icons-dir, legend, limits
docs/github-pages/static/helper.out          # REGENERATED by the cli-helper hook
.specify/memory/constitution.md              # EDIT: amendment 1.1.0 (see Complexity Tracking)
```

**Structure Decision**: one new package `pkg/diagram` holds everything that is independent of the
drawing technology (selection, containers, edges, families, layout), so the future isometric
renderer consumes the same `Scene` and only the projection and drawing differ. The SVG renderer
stays in `pkg/report` because it is a `Formatter` implementation (Constitution II). Icon assets
live next to the renderer that embeds them.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| New package `pkg/diagram` (Constitution II lists the package flow) | The Scene and layout must be reusable by a future isometric renderer and testable without SVG. A render-agnostic model keeps `pkg/report` a pure renderer and lets several renderers share one layout, which Constitution II's separation of concerns asks for. | Putting selection, layout and drawing inside `pkg/report` would make `pkg/report` hold a graph-layout engine and would force the isometric renderer to import SVG code or duplicate the layout. |
| Constitution amendment 1.1.0 (MINOR) in the same change | Output formats list, package flow, and the "reads a plan file" sentence must describe the new format, the new package and the optional icon directory, and the parser rule must be reworded (see next row) with a migration note for in-flight specs, because redefining a principle requires one. Governance: amendments go through a pull request with rationale and version bump. | Leaving the constitution as is would contradict it on the first merge. |
| Additive extension of `pkg/terraform` (`Plan.Configuration`, `ParseAddress`, `ParseReference`) against the rule "adding a format MUST NOT require changing the parser" | The `svg` format needs the dependency relations that live in the plan `configuration` section, which the parser ignored. The proposed reworded rule: "Adding a format MUST NOT change how existing plans are parsed or the output of existing formats; extending the parsed model with new optional fields is allowed and MUST keep every previously valid plan working." | Re-reading the raw plan JSON inside `pkg/diagram` would duplicate the plan-format knowledge outside the package that owns it and would make `pkg/diagram` depend on the input encoding instead of the parsed model. |

## Post-Design Constitution Re-check

After Phase 1 (research, data model, contracts, quickstart) the table above still holds:
the `Scene` carries no attribute values (data-model.md), `Formatter` is unchanged (contract), all
rules are testable with the planned golden and safety tests (quickstart.md), and the only
deviations remain the three tracked in Complexity Tracking. **Gate result: PASS.**
