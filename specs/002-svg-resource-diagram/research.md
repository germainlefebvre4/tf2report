# Research: SVG Resource Diagram

Phase 0 decisions. Each entry: Decision, Rationale, Alternatives considered. The only open points
from the spec (icon pack naming, orientation, node labels) are decided here.

## R-001 Package layout

**Decision**: add `pkg/diagram` (selection, containers, edges, families, layout; imports
`pkg/terraform` only) and render in `pkg/report/svg*.go`. Amend the constitution (1.1.0).

**Rationale**: the spec requires a `Scene` independent of the drawing so an isometric renderer can
reuse it (Key Entities, Assumptions). Layout is the expensive part and must be testable without
parsing SVG. The `Formatter` stays in `pkg/report` (Constitution II).

**Alternatives**: everything in `pkg/report` (rejected: renderer reuse would import SVG code);
everything in `pkg/terraform` (rejected: the parser must not know about drawing).

## R-002 How the formatter learns the filters and the icon directory

**Decision**: keep `Formatter.Format(plan, summary, w)` unchanged. Add

```text
report.SVGOptions{ IconsDir string; TypeFiltered bool; ActionFiltered bool }
report.NewSVGFormatter(opts) *SVGFormatter
```

`NewFormatter(FormatSVG)` returns a formatter with zero options (no icon pack, no filters), so the
format is still reachable through the factory and `ParseFormat("svg")`. `cmd` builds the options
from the resolved configuration when the format is `svg` (helper `newFormatter` in
`cmd/tf2report/formatter.go`).

**Rationale**: the selection rules depend on which filters were given, not only on what is left in
the plan after filtering (spec FR-003, FR-019, FR-024):

| Type filter | Action filter | Drawn |
|-------------|---------------|-------|
| none | none | create, update, delete, replace, plus no-op dependency context (FR-024) |
| given | none | filtered resources of those four actions, no context, no no-op |
| any | given | everything left in the (already filtered) plan except data sources, no context |

Changing the interface would break every caller and test (`docs/API.md` documents it).

**Alternatives**: put the filter flags on `Plan` (rejected: pollutes the input model); change
`NewFormatter` to variadic options (rejected: touches the three existing call paths and tests);
re-run filtering inside the formatter (rejected: duplicates `cmd` logic).

## R-003 Reading relations from the plan

**Decision**: parse `configuration.root_module` into new types in `pkg/terraform`:

- `ConfigModule{ Resources []ConfigResource; ModuleCalls map[string]ModuleCall }`,
  `ModuleCall{ Module ConfigModule }` (child module is under `module_calls.<name>.module`).
- `ConfigResource{ Address, Mode, Type, Name string; DependsOn []string; Expressions map[string]any }`.
- `Plan.Configuration *Configuration` with `json:"configuration,omitempty"`: absent configuration
  is `nil`, no error (spec edge case "older format version").
- `References()` walks `Expressions` recursively (objects, arrays of blocks) and returns every
  string in any `references` array, plus `DependsOn`.

Verified on fixtures: `mixed-actions` has a deleted resource (`terraform_data.gone`) absent from
`configuration`; `destroy-all` has no references; `nested-modules` has module calls and only
`count.index` as a reference. Config addresses carry no index (`module.app.terraform_data.server`)
while `resource_changes` addresses do (`module.app[0].terraform_data.server`).

**Reference grammar** (a reference string, resolved within the module that contains the resource):

| Reference form | Meaning | Edge target |
|----------------|---------|-------------|
| `type.name`, `type.name.attr...`, `type.name[...]...` | managed resource in the same module | all instances of that resource |
| `module.call`, `module.call.output` | module call | the container(s) of that module (spec FR-008) |
| `data.type.name...` | data source | none (data sources are not drawn) |
| `var.`, `local.`, `count.`, `each.`, `path.`, `terraform.`, `self` | not a resource | none |

Unknown forms are ignored, never an error. `depends_on` entries use the same grammar. A reference
to a resource that is not in the plan changes is ignored, except that an unchanged resource may
become context (R-010).

**Alternatives**: decode `configuration` into `map[string]any` and walk it in the diagram package
(rejected: the plan format knowledge belongs to `pkg/terraform`); parse `prior_state` or
`planned_values` (rejected: no references there, only values, which the image must not use).

## R-004 Instances and containers

**Decision**:

- A resource instance in `resource_changes` is a node (address kept verbatim, so `a[0]` and `a[1]`
  are two nodes). Its **config address** is the address with every `[...]` removed.
- Module containers follow the **instance** path (`module.app[0]`, `module.app[1]` are two
  containers); edge resolution uses the index-free path.
- An edge between config resources A -> B expands to every instance of A towards every instance of
  B (spec FR-007). The cross product is accepted in v1: counts in real plans are small and the
  300-resource budget holds. A reference to a module call expands to every container instance of
  that module (one arrow per container; if the module holds several providers, one per provider
  container).
- Provider name: last path segment of `provider_name` (`registry.terraform.io/hashicorp/aws` ->
  `aws`, `registry.opentofu.org/hashicorp/random` -> `random`, `terraform.io/builtin/terraform` ->
  `terraform`). Empty provider name falls back to the resource type prefix before the first `_`.
- Container tree per provider: `provider > module instance > module instance ...`. A module that
  holds resources of two providers appears once under each provider.

**Rationale**: matches the spec (provider > module) and the plan data without inventing structure.

**Alternatives**: collapse instances into one node with a counter (rejected: spec says instances
are drawn individually); one container per module regardless of provider (rejected: spec fixes
provider > module).

## R-005 Layout

**Decision**: left-to-right, deterministic, integer coordinates, three steps per container:

1. **Layering** inside a container: only edges whose two ends are both direct nodes of that
   container count. Longest-path layering where a dependency sits in a lower column than its
   dependents (arrows point from dependent to dependency, towards the left). Cycles are broken by a
   DFS over nodes sorted by address; back edges are dropped from layering only (they are still
   drawn). Cost O(n + e).
2. **Ordering** inside a layer: alphabetical by address, then two barycenter sweeps (stable sort,
   ties by address) to reduce crossings.
3. **Fallback**: a container with no intra-container edge uses a regular grid of
   `ceil(sqrt(n))` columns, filled in address order.

Sizing: node boxes have a fixed size; a container's size is computed from its content (layer
columns, then child containers stacked beneath its nodes, then padding and a header). Provider
containers are stacked vertically in alphabetical order; the SVG `viewBox` is the bounding box plus
a margin; the legend sits above the diagram. Cross-container edges are drawn as cubic curves with
horizontal tangents from the left side of the dependent to the right side of the dependency, under
the nodes and above the containers.

Metrics (final values are tuned in implementation, but they are constants in one file): node
184 x 52, column gap 56, row gap 16, container padding 16, container header 24, monospace width
estimate 7 px per character at 11 px for truncation (no font metrics in the standard library).

**Rationale**: spec FR-011 asks for layers by dependency inside containers, a grid fallback,
alphabetical ties, content-sized image and cycle safety. Integer arithmetic keeps output
byte-stable across platforms (SC-003).

**Alternatives**: global layering across containers (rejected for v1: containers would span the
whole width and the nested boxes would no longer be rectangles); force-directed layout (rejected:
not deterministic and does not project to isometric); external Graphviz (rejected: a runtime
dependency, spec FR-021).

## R-006 Icons

**Decision**:

- **Family ids**: `aws-ec2`, `azure-vm`, `gcp-gce` style (`<cloud>-<service>`); category ids
  `compute`, `network`, `storage`, `database`, `security`; and `generic`. Each node resolves an
  ordered **chain**: service family, then provider family (`aws`, `azure` or `gcp`, selected by the
  type prefix `aws_`, `azurerm_` or `google_`), then category, then `generic`.
- **Service table**: a first table of about 25 services per cloud, kept in `pkg/diagram/family.go`
  (for example AWS: `aws_instance`, `aws_vpc`, `aws_subnet`, `aws_security_group`, `aws_s3_bucket`,
  `aws_db_instance`, `aws_lambda_function`, `aws_lb`, `aws_iam_role`, `aws_eks_cluster`; Azure:
  `azurerm_linux_virtual_machine`, `azurerm_virtual_network`, `azurerm_subnet`,
  `azurerm_storage_account`, `azurerm_key_vault`, `azurerm_kubernetes_cluster`; GCP:
  `google_compute_instance`, `google_compute_network`, `google_compute_subnetwork`,
  `google_storage_bucket`, `google_sql_database_instance`, `google_container_cluster`). Types not
  in the table fall to the category heuristic.
- **Category heuristic** on the `_`-separated tokens of the type name after the provider prefix
  (whole-token match, never substring: `aws_lb` matches `lb`, `aws_monkey_thing` does not match
  `key`), first match wins, in this order: security (`iam`, `role`, `policy`, `kms`, `key`, `secret`, `vault`, `security`, `firewall`,
  `certificate`), database (`db`, `database`, `sql`, `redis`, `dynamodb`, `cosmos`, `cache`),
  network (`vpc`, `subnet`, `network`, `route`, `gateway`, `lb`, `dns`, `nat`, `peering`),
  storage (`bucket`, `storage`, `disk`, `volume`, `efs`, `blob`), compute (`instance`,
  `virtual_machine`, `function`, `cluster`, `container`, `app_service`, `lambda`); otherwise
  `generic`. The provider prefix (`aws_`, `azurerm_`, `google_`) selects the cloud for the service
  table and the provider family level; the category heuristic is provider independent so unknown
  providers such as `random` still get an icon (their chain has no provider level). Because
  matching is by whole token, the service table lists `aws_lb`, `aws_alb` and `aws_elb`
  explicitly.
- **Built-in glyphs**: 6 original, hand-authored monochrome SVG glyphs (the five categories and
  `generic`), embedded with `go:embed`. No official artwork is committed (spec FR-013), so the
  provider family level (`aws`, `azure`, `gcp`) has no built-in glyph and applies only when the
  user pack provides it; without it the chain falls through to the category glyph.
- **Pack lookup**: `--icons-dir` holds one file per family id, named `<family-id>.svg`
  (`aws-ec2.svg`, `aws.svg`, `compute.svg`, `generic.svg`). For each level of the node chain the loader tries
  the pack first, then the built-in glyph; the first hit wins. A pack file that fails validation is
  skipped (the next candidate is used). A missing or unreadable directory is an error before any
  output is written (spec FR-016); a directory with no usable file is not an error.
- **Reuse in the document**: each used icon is emitted once as a `<symbol id="icon-<family>">` in
  `<defs>` and referenced with `<use href="#icon-...">` (spec FR-015).

**Alternatives**: one flat family per type (rejected: unmanageable table); fetch icons at run time
(rejected: offline, FR-021); embed the official packs (rejected: licenses not verified, spec
Assumptions).

## R-007 Visual style and labels

**Decision**: each status has a color, a symbol and a border style (spec FR-005):

| Status | Color (fill tint / stroke) | Symbol | Border |
|--------|----------------------------|--------|--------|
| create | green | `+` | solid, 2 px |
| update | orange | `~` | dashed `6 3` |
| delete | red | `-` | dotted `2 3` |
| replace | purple | `-/+` | dash-dot `8 3 2 3` |
| no-op (context or requested) | grey, 60 % opacity | none | solid, 1 px |

Colors are chosen so that a text label on the tint keeps a 4.5:1 contrast ratio; they are constants
in `svg_style.go` and the legend is generated from the same table. A node shows the symbol, the
resource name (with its index) on the first line and the resource type on the second line, both
truncated with an ellipsis to fit; the full address is the `<title>` tooltip. The icon sits at the
left of the node.

**Alternatives**: color only (rejected: spec FR-005); full address as the label (rejected: too
wide, tooltip keeps it available).

## R-008 Output safety

**Decision**:

- All text and attribute values go through one escaping helper (`encoding/xml.EscapeText` plus
  removal of characters that are illegal in XML 1.0), including addresses, names, container labels
  and the `<title>` tooltips (spec FR-018).
- The `Scene` has no field that can hold an attribute value; the renderer cannot leak what it does
  not receive (spec FR-017). A safety test searches the output for every value of the sensitive
  fixtures.
- **Icon pack sanitizer** (allowlist, spec FR-016): the file must parse as XML with root `svg`, be
  at most 256 KiB, and contain only these elements: `svg g path rect circle ellipse line polyline
  polygon defs linearGradient radialGradient stop clipPath title desc`. Attributes are checked
  against a presentation allowlist (geometry, fill, stroke, transform, opacity, `viewBox`, `id`,
  gradient attributes). Rejected: any `on*` attribute, `style` with `url(` or `@import`, `href`
  or `xlink:href` that is not a local `#fragment`, `<script>`, `<foreignObject>`, `<image>`,
  `<use>`, `<style>`, DOCTYPE and entity declarations. The accepted content is re-serialized (not
  copied byte for byte) with ids prefixed by the family id to avoid collisions.

**Alternatives**: denylist (rejected: easy to miss a vector); embedding the file as an `<image>`
data URI (rejected: loses `currentColor` styling and still needs validation).

## R-009 Determinism

**Decision**: no iteration over a map in any code path that produces output; maps are turned into
sorted slices first. Coordinates and sizes are `int`. Numbers are written with `strconv.Itoa`. No
timestamps, random ids or process-dependent data in the document. Ids are derived from sorted node
order (`n1`, `n2`, ...), not from addresses (addresses can contain characters invalid in an id).

**Rationale**: spec SC-003.

## R-010 Selection of resources, context and special actions

**Decision** (implements spec FR-003, FR-009, FR-019, FR-024):

- Data sources (`mode: data` or action `read`) are never drawn.
- `replace` (delete + create) is one node with status replace.
- Status comes from `GetPrimaryAction`; a resource with only `moved` or `import` metadata and a
  `no-op` action stays no-op (hidden by default). An import with an `update` action is an update.
- Edges whose source or target is a `delete` node are dropped (spec FR-009).
- Context (only when neither filter is given): after choosing the drawn changed nodes, every no-op
  node that is the direct target of an edge from a drawn node is added once, with `Context = true`;
  no edge starts from a context node and its own dependencies are not followed.
- With any filter, edges towards nodes that are not drawn are dropped.
- An empty selection produces a valid document with a centered message and the legend (spec FR-020).

## R-011 CLI behavior around `--action no-op`

**Decision**: `--action no-op` already works mechanically (the action filter compares strings and
`GetPrimaryAction` returns `no-op`). The help text and documentation of `--action` and the config
key list `no-op` as an accepted value for the `svg` format. For the three existing formats the
output is unchanged because they already exclude no-op changes. The baseline contract's
"Provisional" note for values outside the four listed gains a sentence for `no-op`.

`--icons-dir` is registered as a flag only (no config key, no environment variable in v1); with
any other format it is ignored, and a notice goes to standard error in verbose mode (spec FR-014).

## R-012 Tests and fixtures

**Decision**:

- Synthetic fixtures under `testdata/plans/diagram/` with a real `configuration` section:
  `aws-network` (vpc -> subnet -> instance, security group), `azure-network`, `gcp-network`,
  `modules-context` (changed resource referencing an unchanged one, module output reference,
  count instances), `cycle` (two resources referencing each other), `no-configuration`.
- Golden SVG files per fixture and option combination (no filter, type filter, action filter,
  action filter including `no-op`, icon pack).
- `testdata/icons/`: one valid pack, one pack with each hostile construct of R-008.
- Scale test: generated 300-resource plan, asserts time budget and determinism (two runs equal).
- Safety test: every attribute value of the sensitive fixtures is absent from the output.
- CLI tests: `--format svg`, `--icons-dir` with svg and with other formats (silent, verbose
  notice), invalid directory error and exit code, filters x statuses matrix.
- Existing goldens and examples tests are the regression guard for FR-022 (SC-008).

## R-013 Documentation and constitution

**Decision**: update `README.md`, `docs/USAGE.md`, `docs/OUTPUT_FORMATS.md`,
`docs/CONFIGURATION.md` (format value), `docs/API.md` (new formatter and options),
`docs/github-pages/docs/*` equivalents and regenerate `helper.out` (the `cli-helper` hook). Describe
the legend, the filter behavior table of R-002 and the limits (no isometric, no PNG, module outputs
point to containers, deleted resources have no arrows). Do not embed an example SVG
(Constitution, Technology constraints). Amend the constitution to 1.1.0: add `svg`, add
`pkg/diagram` to the package flow, allow the optional local icon directory input.

## R-014 Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Cross-product edges for large `count` on both ends | Many arrows, slow render | Accepted in v1; scale test at 300 resources; revisit with aggregation later |
| Cross-container arrows get long and cross | Readability | Curves drawn under nodes; stacked containers; revisit with global layering later |
| Text width is estimated (no font metrics) | Labels may clip in some fonts | Monospace estimate, ellipsis, full address in tooltip |
| Official icon licensing | Cannot ship official artwork | Only generic glyphs shipped; pack is user supplied; license check is a prerequisite for any later asset commit |
| Plan formats vary by provider (nested expressions) | Missed references | Recursive walk, unknown forms ignored; fixtures from several clouds |
