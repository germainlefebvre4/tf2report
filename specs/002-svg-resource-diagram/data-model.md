# Data Model: SVG Resource Diagram

Types introduced or extended by the feature. Field names are indicative; the contract is the
behavior described in `spec.md` and `research.md`. None of these types holds a resource attribute
value (spec FR-017).

## pkg/terraform (extended)

### Plan (edited)

| Field | Type | Notes |
|-------|------|-------|
| `Configuration` | `*Configuration` | JSON `configuration`, `omitempty`. `nil` when the plan has no configuration section. |

The existing fields and the behavior of `Parse`, `ParseFile`, `FilterByType`, `FilterByAction` and
`NewSummary` are unchanged. The JSON, text and markdown formatters do not read `Configuration`.

### Configuration, ConfigModule, ModuleCall, ConfigResource (new)

| Type | Fields | Notes |
|------|--------|-------|
| `Configuration` | `RootModule ConfigModule` | JSON `root_module`. |
| `ConfigModule` | `Resources []ConfigResource`, `ModuleCalls map[string]ModuleCall` | A module body. |
| `ModuleCall` | `Module ConfigModule` | JSON `module`. Only the nested body is read. |
| `ConfigResource` | `Address, Mode, Type, Name string`; `DependsOn []string`; `Expressions map[string]any` | `Address` is index-free, relative to its module. `Expressions` is walked, never exposed to the diagram. |

`ConfigResource.References() []string` returns the `references` strings found anywhere in
`Expressions` plus `DependsOn`, in a stable order, without duplicates.

### Address helpers (new)

| Function | Input | Output |
|----------|-------|--------|
| `ParseAddress(s)` | `module.app[0].aws_instance.web["a"]` | `Address{ Modules []ModuleStep{Name, Key}, Mode, Type, Name, Key }` |
| `Address.ConfigAddress()` | parsed address | index-free full address (`module.app.aws_instance.web`) |
| `ParseReference(s)` | `aws_subnet.a.id`, `module.net.vpc_id`, `data.x.y`, `var.z` | `Reference{ Kind, ModuleCall, Type, Name }` with `Kind` in `resource`, `module`, `data`, `other` |

Parsing never fails the run: an address or reference that does not match the grammar yields
`Kind: other` (references) or a best-effort address (resources), and is ignored by the diagram.

## pkg/diagram (new)

### Scene

The picture of one plan: everything a renderer needs, in integer coordinates.

| Field | Type | Notes |
|-------|------|-------|
| `Containers` | `[]*Container` | Top level: one per provider, sorted by provider name. |
| `Nodes` | `[]*Node` | All nodes, sorted by address. |
| `Edges` | `[]Edge` | Sorted by (from, to). |
| `Statuses` | `[]Status` | Statuses present, in fixed order, for the legend. |
| `Width`, `Height` | `int` | Content bounding box without margin. |
| `Empty` | `bool` | True when nothing is drawn (spec FR-020). |

### Container

| Field | Type | Notes |
|-------|------|-------|
| `ID` | `string` | Stable, derived from sorted order. |
| `Kind` | `provider` or `module` | |
| `Label` | `string` | Provider name, or module step with key (`module.app[0]`). |
| `Path` | `string` | Full instance path, used for edge targets. |
| `Children` | `[]*Container` | Sorted by label. |
| `Nodes` | `[]*Node` | Direct nodes, sorted by address. |
| `X, Y, W, H` | `int` | Filled by layout. |

Relations: a container belongs to at most one parent. Identity: (provider, instance path).

### Node

| Field | Type | Notes |
|-------|------|-------|
| `ID` | `string` | `n1`, `n2`, ... in address order. |
| `Address` | `string` | Verbatim, with index. Tooltip text. |
| `ConfigAddress` | `string` | Index-free, used to expand edges. |
| `Type`, `Name` | `string` | Resource type and name (with index). |
| `Provider` | `string` | Short provider name. |
| `Status` | `Status` | create, update, delete, replace, no-op. |
| `Context` | `bool` | Unchanged and drawn only as dependency context (FR-024). |
| `Family` | `[]string` | Ordered icon chain: service, provider (`aws`, `azure` or `gcp`, only for the `aws_`, `azurerm_` and `google_` prefixes), category, `generic`. |
| `Layer`, `X, Y, W, H` | `int` | Filled by layout. |

Validation: `Address` unique across nodes; a node exists only for create, update, delete, replace
or (requested or context) no-op; never for data sources.

### Edge

| Field | Type | Notes |
|-------|------|-------|
| `From` | `*Node` | The dependent resource. |
| `To` | `*Node` or `*Container` | The dependency: a node, or a module container (FR-008). |
| `Points` | `[]Point` | Bezier control points, filled by layout. |

Rules: no edge starts or ends at a `delete` node; no edge starts at a context node; edges towards
nodes not drawn are dropped; duplicates removed.

### Status

`create`, `update`, `delete`, `replace`, `noop`, with fixed order for sorting and the legend.
Color, symbol and border are renderer concerns (`pkg/report`), not Scene fields, so another
renderer can restyle them.

### Options for building

`diagram.Options{ TypeFiltered bool; ActionFiltered bool }`, passed by the formatter. It decides
no-op visibility and context as tabulated in `research.md` R-002.

## pkg/report (extended)

### SVGOptions

| Field | Type | Notes |
|-------|------|-------|
| `IconsDir` | `string` | Empty means built-in glyphs only. |
| `TypeFiltered` | `bool` | A type filter was given. |
| `ActionFiltered` | `bool` | An action filter was given. |
| `OnSkip` | `func(file, reason string)` | Optional. Called for each icon pack file that is skipped; `cmd` prints it to standard error in verbose mode. The formatter itself never writes to standard error. |

### Icon resolution

`family chain` (from the node) x (`pack`, `built-in`) -> first valid SVG found. Built-in glyphs
exist for the five categories and `generic` only; service and provider families exist only in a user pack.

## State and lifecycle

No persisted state. A run is: parse plan -> filter (existing) -> build Scene -> layout -> render ->
write. Errors before the first byte of output (invalid icon directory) leave standard output empty
and exit non-zero; once rendering starts the document is built in memory and written once, so a
failure never leaves a truncated SVG.
