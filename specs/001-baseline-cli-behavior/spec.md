# Feature Specification: Baseline CLI Behavior

**Feature Branch**: `chore/speckit-brownfield` (spec directory `001-baseline-cli-behavior`; no feature branch created)

**Created**: 2026-11-22

**Status**: Draft

**Input**: User description: "Baseline de l'existant: specify the current behavior of the tf2report CLI (plan parsing, filters, markdown/text/json formats, configuration) from the Repomix snapshot, as the brownfield reference for future changes."

> This is a **baseline specification**: it documents what tf2report does today, verified by reading
> the code and running the tool against the repository's sample plan. It is the reference that
> future feature specs and plans are compared against (Constitution, Principle II). Where current
> behavior is questionable, it is recorded in *Known Gaps* and not silently "fixed" here.

## Clarifications

### Session 2026-11-25

- Q: Should the automated tests derived from this spec also pin the questionable behaviors listed
  in Known Gaps, or only the intended behavior? → A: Pin the current behavior, including the gaps,
  each marked as provisional so that fixing one is a deliberate, visible test change.
- Q: For which Terraform plan versions must the baseline be guaranteed? → A: The Terraform
  versions currently supported by HashiCorp; plans from other versions are still accepted, with
  no guarantee.
- Q: Does this feature include writing the characterization tests, or only the spec? → A: It
  delivers the spec and the tests (parser, filters, configuration, the three formats, errors),
  without modifying application code.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Summarize a Terraform plan (Priority: P1)

A DevOps engineer exports a Terraform plan as JSON and runs the tool on it. They get a
human-readable summary of what will change: total number of changes, counts per action
(add, change, destroy, replace), the resources grouped by action, and the changes grouped by
resource type. The Terraform version recorded in the plan is shown when present.

**Why this priority**: this is the product's reason to exist; every other story refines it.

**Independent Test**: run the tool on the sample plan with only the plan path given. The report
lists 4 total changes (1 add, 1 change, 1 destroy, 1 replace) and exits successfully.

**Acceptance Scenarios**:

1. **Given** a valid plan with changes, **When** the tool runs with the plan path, **Then** a
   report is written to standard output with the Terraform version, a summary table of counts,
   resources grouped by action, and changes grouped by resource type, and the tool exits with
   code 0.
2. **Given** a plan with no resource changes (or an empty JSON object), **When** the tool runs,
   **Then** the report states that there are no changes and that infrastructure is up-to-date,
   and exits with code 0.
3. **Given** a plan containing resources that are only read or unchanged (no-op), **When** the
   tool runs, **Then** those resources are not counted nor listed.
4. **Given** a resource that is deleted then created in the same plan, **When** the tool runs,
   **Then** it is counted once, as a replacement, not as a delete plus an add.
5. **Given** no plan path is provided anywhere, **When** the tool runs, **Then** it looks for
   `terraform.tfplan.json` in the current directory.

---

### User Story 2 - Choose the report format (Priority: P2)

A user wants the same summary as Markdown (to paste in a pull request), plain text (for a
terminal or a log) or JSON (for another program to consume).

**Why this priority**: it makes the report usable in the three main contexts (review, console,
automation). The default format keeps the P1 story working with no extra input.

**Independent Test**: run the tool on the sample plan three times, once per format, and check
that each output carries the same counts and resources in its own style.

**Acceptance Scenarios**:

1. **Given** no format is requested, **When** the tool runs, **Then** the report is Markdown.
2. **Given** the format `text`, **When** the tool runs, **Then** the report is plain text with
   aligned columns and no Markdown syntax.
3. **Given** the format `json`, **When** the tool runs, **Then** the output is a single valid JSON
   document with the Terraform version (when known), a summary of counts, the list of changes
   (address, type, name, action) and per-type counts with the resource addresses.
4. **Given** the aliases `md` or `txt`, **When** the tool runs, **Then** they behave like
   `markdown` and `text`.
5. **Given** an unsupported format such as `xml`, **When** the tool runs, **Then** it fails
   with a message listing the supported formats (markdown, text, json) and exit code 1.

---

### User Story 3 - Focus on part of the plan with filters (Priority: P3)

A reviewer only cares about some resource types (for example security-related ones) or some
actions (for example only deletions). They restrict the report to those.

**Why this priority**: valuable for large plans, but the tool is already useful without it.

**Independent Test**: run the tool with one type and one action on the sample plan and verify
that only matching resources are counted, listed, and summed in the totals.

**Acceptance Scenarios**:

1. **Given** one or more resource types, **When** the tool runs, **Then** only changes of those
   types are reported, and totals reflect only those changes.
2. **Given** one or more actions among create, update, delete, replace, **When** the tool runs,
   **Then** only changes whose action matches are reported.
3. **Given** both a type filter and an action filter, **When** the tool runs, **Then** a change
   must satisfy both to be reported.
4. **Given** filters that match nothing, **When** the tool runs, **Then** the report states that
   there are no changes and the tool exits with code 0 (this is not an error).

---

### User Story 4 - Configure once, run anywhere (Priority: P3)

A team commits a configuration file (or sets environment variables in CI) holding the plan path,
output format, filters and verbosity, so that the command line stays short and consistent.

**Why this priority**: convenience and consistency; every setting is also available as a flag.

**Independent Test**: put a configuration file in the working directory, run the tool with no
flags and check it behaves as configured; then add a flag and check the flag wins.

**Acceptance Scenarios**:

1. **Given** a configuration file `tf2report.yaml` in the working directory (or in
   `~/.config/tf2report`), **When** the tool runs without flags, **Then** its plan path, format,
   filters and verbosity are applied.
2. **Given** a configuration file passed explicitly, **When** the tool runs, **Then** that file is
   used, and an unreadable or missing file is an error.
3. **Given** environment variables prefixed `TF2REPORT_`, **When** the tool runs, **Then** they
   can provide the plan path, the format and the verbosity (filters: see Known Gaps).
4. **Given** a setting present both in configuration and as a flag, **When** the tool runs,
   **Then** the flag takes precedence.
5. **Given** no configuration at all, **When** the tool runs, **Then** defaults apply: plan
   `terraform.tfplan.json`, format Markdown, verbosity info, no filters.

---

### User Story 5 - Understand and trust failures (Priority: P3)

A user (or a CI job) needs to know when the tool could not produce a report, and why.

**Why this priority**: automation depends on a reliable exit code and an actionable message.

**Independent Test**: run the tool on a missing file, an invalid JSON file and an invalid format,
and verify that each exits with code 1 and a distinct, readable message on standard error.

**Acceptance Scenarios**:

1. **Given** a plan path that does not exist, **When** the tool runs, **Then** it fails with a
   "plan file not found" message and exit code 1, and writes nothing to standard output.
2. **Given** a file that is not valid JSON, **When** the tool runs, **Then** it fails with a
   message that the plan could not be parsed, and exit code 1.
3. **Given** verbose mode, **When** the tool runs, **Then** progress is written to standard error
   (number of resource changes parsed, resources remaining after each filter) and the report on
   standard output is unchanged.
4. **Given** the version flag, **When** the tool runs, **Then** it prints version, commit and build
   time.

---

### Edge Cases

- A plan whose `resource_changes` is absent or empty is treated as "no changes", not as an error.
- A plan with no Terraform version produces a report without the version line (and JSON without
  the version field).
- A change with no action, or an unknown action, is treated as a no-op and left out of the report.
- Unknown JSON fields in the plan are ignored; plan output changes (`output_changes`) are parsed
  but not reported.
- An action filter containing an unknown value (for example `bogus`) is accepted and matches
  nothing, so the report shows no changes (see Known Gaps).
- A resource type filter containing a type absent from the plan matches nothing and is not an
  error.
- Report ordering is stable: in Markdown and text, resource types are sorted alphabetically and
  resources are sorted by name or address within each section; JSON lists changes in plan order.
- Attribute values of resources (including sensitive ones) are never written to the report; only
  address, type, name and action are.
- The tool never contacts Terraform, a cloud provider or the network; it only reads files.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST read a Terraform plan in the JSON form produced by the Terraform plan
  export, from a path given by flag, configuration or default (`terraform.tfplan.json`).
- **FR-002**: System MUST classify each resource change as create, update, delete, replace, read
  or no-op, treating a delete combined with a create as a replace.
- **FR-003**: System MUST exclude read and no-op changes from counts and listings.
- **FR-004**: System MUST compute a total, per-action counts, and per-resource-type counts.
- **FR-005**: System MUST produce reports in Markdown (default), plain text and JSON, accepting
  `markdown`/`md`, `text`/`txt` and `json` as format names.
- **FR-006**: Every format MUST convey the same information: Terraform version when known,
  summary counts, resources by action, and changes by resource type.
- **FR-007**: System MUST state "no changes" in text and Markdown reports (and emit zero counts
  with empty lists in JSON) when nothing is left to report.
- **FR-008**: System MUST support filtering by one or more resource types and by one or more
  actions (create, update, delete, replace), combined with AND semantics, applied before the
  summary is computed.
- **FR-009**: System MUST accept settings from a YAML configuration file (explicit path, working
  directory, or `~/.config/tf2report`), from `TF2REPORT_`-prefixed environment variables and from
  flags, with flags overriding the rest and built-in defaults applying last.
- **FR-010**: System MUST write the report to standard output and diagnostics (errors, verbose
  progress) to standard error.
- **FR-011**: System MUST exit with code 0 on success, including when no changes remain, and
  with code 1 on any failure (configuration, missing plan file, invalid JSON, unsupported format).
- **FR-012**: System MUST fail with an explicit message when the plan file is missing, is not
  valid JSON, or when the requested format is unsupported, without writing a partial report.
- **FR-013**: System MUST NOT output resource attribute values and MUST NOT perform any network
  call or run Terraform.
- **FR-014**: System MUST produce identical output for identical input and options (stable
  ordering).
- **FR-015**: System MUST provide help and a version output including version, commit and build
  time.
- **FR-016**: Characterization tests for this baseline MUST pin the current behavior, including
  each item of Known Gaps that is observable behavior (the documentation-only gap is excluded); every such test MUST be marked as provisional and reference its gap, so
  that correcting a gap is a deliberate, visible change to that test.

### Key Entities

- **Plan**: the exported Terraform plan; carries a format version, the Terraform version, the
  list of resource changes and output changes.
- **Resource Change**: one planned change; has an address, type, name, provider, and a change.
- **Change**: the list of actions planned (for example create, or delete then create) with the
  before/after state, from which the primary action is derived.
- **Action**: create, update, delete, replace, read or no-op.
- **Summary**: total and per-action counts, plus counts per resource type.
- **Configuration**: plan path, output format, filters (resource types, actions) and verbosity.
- **Report**: the rendered result, in one of three formats.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For the baseline sample plan (a synthetic plan with 4 changes), the report shows exactly 1 add,
  1 change, 1 destroy and 1 replace in all three formats.
- **SC-002**: For any valid plan, the counts shown equal the number of resources per action in
  the plan, excluding read and no-op (100% agreement).
- **SC-003**: Running the tool twice on the same plan with the same options produces
  byte-identical output.
- **SC-004**: Every failure scenario in User Story 5 ends with exit code 1 and a one-line cause
  a user can act on without reading the source code.
- **SC-005**: A new user obtains a first report from an existing plan with one command and no
  configuration.
- **SC-006**: No attribute value from any sample plan appears in any report format.
- **SC-007**: Every acceptance scenario of User Stories 1 to 5 is covered by at least one
  automated test that passes on the unmodified application, and the test suite runs green in the
  pull request check.

## Assumptions

- The baseline describes behavior at the current commit of the `main` line; documentation
  (README, `docs/`) is treated as secondary to observed behavior when they disagree.
- Behavior is guaranteed for plans produced by the Terraform versions currently supported by
  HashiCorp, in JSON form. Plans from other versions are accepted as before (the format version is
  not validated) but carry no guarantee, and the sample fixtures used by tests MUST come from
  supported versions.
- The tool is run locally or in CI on a trusted machine; plan files may hold sensitive values and
  are handled by the caller.
- Windows is out of scope (binaries are built for Linux and macOS only).
- Plan size is within what a CI runner handles in memory; no performance target is claimed.
- This feature delivers this spec and the characterization tests that pin it (FR-016), covering
  the parser, filters, configuration, the three report formats and the error paths. It MUST NOT
  modify application behavior: any code change is limited to what tests strictly need to run.

## Known Gaps

Observed while building the baseline; each deserves its own spec rather than a silent change.
Until corrected, the behavior below is the baseline and is pinned by provisional tests (FR-016).

- An unknown value for the action filter is accepted silently, and yields an empty report instead
  of an error.
- The `verbosity` setting only has an effect through verbose mode; levels debug/info/warn/error
  are otherwise not differentiated.
- When the tool fails, the error message is printed twice on standard error, followed by the
  usage text.
- A file named `tf2report` (no extension) in the working directory, such as the tool's own
  binary, is read as the configuration file and breaks the run.
- The filter settings (`filters.resource_types`, `filters.actions`) can be given through the
  environment only with a dotted variable name (for example `TF2REPORT_FILTERS.ACTIONS`), which
  most shells cannot export; the underscore form `TF2REPORT_FILTERS_ACTIONS` is silently ignored.
- The README states the plan path is required and shows an example report whose layout no longer
  matches the actual output; the project instructions forbid example outputs in documentation.
  This is a documentation gap and is not pinned by a test; the default plan path itself is covered
  by a regular test (User Story 1, scenario 5).
- A plan in the exact shape Terraform exports is rejected with a parse error when a change carries a
  boolean sensitivity flag (`"before_sensitive": false`, which Terraform emits for created
  resources) or when an output change holds a plain value instead of an object. Only plans whose
  sensitivity flags are objects and whose outputs are objects are read.
- Markdown and text reports show only a resource's type and name, not its module path or index, so
  `web[0]` and `web[1]`, or the same name in two modules, appear as identical rows.
- No automated tests cover any scenario above.
