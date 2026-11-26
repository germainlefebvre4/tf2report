# Contract: SVG output format, `--icons-dir` and icon pack

User-facing interfaces added by this feature. It extends the baseline contract in
`specs/001-baseline-cli-behavior/contracts/cli-contract.md`; everything not listed here is
unchanged.

## CLI

| Flag | Short | Value | Effect |
|------|-------|-------|--------|
| `--format` | `-f` | `svg` (new value) | Writes one SVG document to standard output. |
| `--icons-dir` | | directory path | Optional icon pack read when the format is `svg`. |
| `--action` | `-a` | `create`, `update`, `delete`, `replace`, and `no-op` | `no-op` asks for unchanged resources in the diagram. |

`svg` is also accepted from the configuration file (`output_format: svg`) and from
`TF2REPORT_OUTPUT_FORMAT=svg`, with the usual precedence (flag > environment > file > default).
There is no alias for `svg`. `ParseFormat` error text lists `markdown, text, json, svg`.

`--icons-dir` has no configuration key and no environment variable in this version.

### Streams and exit codes

| Situation | Standard output | Standard error | Exit |
|-----------|-----------------|----------------|------|
| Plan with changes, `--format svg` | the SVG document only | empty (verbose adds progress) | 0 |
| Nothing to draw (empty selection) | a valid SVG with an explicit message and the legend | empty | 0 |
| `--icons-dir` missing or unreadable, format svg | empty | error naming the directory | non-zero |
| `--icons-dir` with a usable or partly unusable pack | the SVG; unusable files fall back to built-in icons | empty (verbose may list skipped files) | 0 |
| `--icons-dir` with a format other than `svg` | the requested format, unchanged | empty; with `--verbose` one notice line that the option is ignored | 0 |
| Plan file missing, invalid JSON, invalid format | unchanged baseline behavior | unchanged | non-zero |

Existing formats: output byte-for-byte unchanged for every plan (spec FR-022, SC-008).

### What is drawn (filters x statuses)

| `--type` | `--action` | create / update / delete / replace | no-op | Dependency context (no-op) |
|----------|------------|------------------------------------|-------|-----------------------------|
| not given | not given | all | no | yes, one hop, muted |
| given | not given | matching types | no | no |
| any | given | matching, among those left after filtering | only if `no-op` is in the action list | no |

Data sources (`read`) are never drawn. A replaced resource is one node. Arrows towards nodes that
are not drawn are omitted. Deleted resources never have arrows.

## SVG document contract

The document is UTF-8, a single `<svg>` root, no XML declaration dependency on external resources
(no external href, no scripts, no fonts to load, no `<foreignObject>`).

| Part | Contract |
|------|----------|
| Root | `xmlns="http://www.w3.org/2000/svg"`, `viewBox` computed from the content, `role="img"`, `<title>` and `<desc>` summarizing the plan (counts per status). |
| Legend | A group `class="legend"` with one entry per status present: color, symbol, border style and label. |
| Containers | `<g class="container provider|module">` with a `<rect>`, a label, nested in provider > module order. |
| Edges | `<g class="edges">` containing one `<path>` per edge with an arrow marker; drawn above containers and below nodes. |
| Nodes | `<g class="node status-create|update|delete|replace|noop" [context]>` with a `<title>` holding the full address, a `<rect>` (border style by status), the status symbol, the name and type text, and a `<use>` of the icon. |
| Icons | Each used icon appears once as `<symbol id="icon-...">` in `<defs>`; nodes reference it with `<use href="#icon-...">`. |
| Text | All text and attribute values are XML-escaped; no characters invalid in XML 1.0. |
| Content | Only addresses, types, names, providers, modules and statuses. Never an attribute value (before, after, unknown, sensitive). |
| Determinism | Same plan and options give identical bytes. Coordinates are integers. |
| Status cues | create green `+` solid; update orange `~` dashed; delete red `-` dotted; replace purple `-/+` dash-dot; no-op grey, muted, no symbol. Color is never the only cue. |

The exact class names and layout metrics are an internal detail pinned by golden files, not a
compatibility promise; the table above is the stable part.

## Icon pack contract (`--icons-dir`)

| Item | Rule |
|------|------|
| Location | One flat directory, read non-recursively. Only `*.svg` files are considered. |
| File name | `<family-id>.svg`, for example `aws-ec2.svg`, `azure-vm.svg`, `gcp-gce.svg`, a provider name (`aws.svg`, `azure.svg`, `gcp.svg`), or a category/generic name: `compute.svg`, `network.svg`, `storage.svg`, `database.svg`, `security.svg`, `generic.svg`. |
| Resolution | For a resource, the first existing file along its chain (service family, then provider family, then category, then `generic`) is used; at each level the pack wins over the built-in glyph. |
| Size | At most 256 KiB per file. |
| Allowed elements | `svg g path rect circle ellipse line polyline polygon defs linearGradient radialGradient stop clipPath title desc`. |
| Rejected content | Scripts, event attributes (`on*`), `<style>`, `<image>`, `<use>`, `<foreignObject>`, external or non-fragment `href`/`xlink:href`, `url(` or `@import` in `style`, DOCTYPE and entities. |
| On rejection | The file is skipped silently (listed in verbose mode) and the next candidate is used; the diagram is still produced. |
| Invalid directory | A path that does not exist, is not a directory, or cannot be read is a fatal error before any output. |
| Output | Accepted content is re-serialized into a `<symbol>`; ids are prefixed to avoid collisions. |

Built-in glyphs (shipped in the binary): `compute`, `network`, `storage`, `database`, `security`,
`generic`. No official AWS, Azure or GCP artwork is shipped, so the service and provider levels of
the chain apply only when a pack provides them.

## Library contract (`pkg/report`, `pkg/diagram`, `pkg/terraform`)

- `report.FormatSVG` and `report.ParseFormat("svg")` exist; `report.NewFormatter(report.FormatSVG)`
  returns a formatter with default options (no pack, no filters).
- `report.NewSVGFormatter(report.SVGOptions{...})` returns a formatter with options. It implements
  the unchanged `Formatter` interface.
- `diagram.Build(plan, diagram.Options) (*Scene, error)` and `diagram.Layout(*Scene)` are
  render-independent; the Scene carries no attribute value.
- `terraform.Plan` gains an optional `Configuration`; all existing exported symbols keep their
  signatures.
