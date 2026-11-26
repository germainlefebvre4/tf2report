# Feature Specification: SVG Resource Diagram

**Feature Branch**: `feat/show-visual-diagram`

**Created**: 2026-11-26

**Status**: Draft

**Input**: User description: "Ajouter un format de sortie SVG (--format svg) qui dessine le diagramme des ressources d'un plan Terraform/OpenTofu, avec leurs changements. Objectif : produire une image qui permet de visualiser les ressources de l'infrastructure et leur statut (create / update / delete / replace / no-op), sous forme de diagramme d'architecture 2D."

## Clarifications

### Session 2026-11-26

- Q: By default, with no filter, are unchanged (no-op) resources drawn? → A: No. Only changed resources are drawn by default; no-op resources appear only when explicitly requested through the action filter, apart from the dependency context described in the next two answers.
- Q: When a changed resource depends on an unchanged resource that is not drawn, what does the diagram show? → A: The unchanged resource is drawn as muted context, with its arrow, only because a changed resource depends on it.
- Q: Does the dependency context still appear when a type or action filter is given? → A: No. Context appears only when no filter is given; with a filter, only the selected resources are drawn and arrows to resources not drawn are omitted.
- Q: What happens when `--icons-dir` is given with a format other than `svg`? → A: It is ignored without failing; a notice is written to the error output only in verbose mode.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See what a plan changes at a glance (Priority: P1)

An infrastructure engineer or reviewer has a plan file and wants a picture of the
infrastructure it touches, instead of reading a list. They request the `svg` output format
and get a diagram in which every changed resource appears once, grouped by provider and by
module, and visibly marked as created, updated, deleted or replaced. Unchanged resources are
left out unless the user asks for them or a changed resource depends on them.

**Why this priority**: This is the whole point of the feature. Without a readable, correctly
classified picture of the resources and their change status, nothing else has value. It is
also the smallest slice that can be demonstrated on its own.

**Independent Test**: Run the tool with the `svg` format on a plan that mixes all five
statuses, once without filter and once asking for unchanged resources too, open the results in
a browser or an image viewer, and check that every expected resource is present, grouped
under its provider and module, and distinguishable by status.

**Acceptance Scenarios**:

1. **Given** a plan with resources to create, update, delete, replace and leave unchanged,
   **When** the user runs the tool with the `svg` format and no filter, **Then** a single valid
   SVG image is produced in which each resource to create, update, delete or replace appears
   exactly once with a visual marker for its status, and no unchanged resource appears except
   as dependency context (see User Story 2).
2. **Given** the same plan, **When** the user runs the tool with the `svg` format and asks for
   unchanged resources through the action filter, **Then** those resources also appear, drawn
   muted, so the changed resources stay the visual focus.
3. **Given** a plan whose resources belong to several providers and to nested modules,
   **When** the diagram is produced, **Then** resources are drawn inside a container per
   provider, and inside nested containers for their module path, so the grouping matches the
   resource addresses.
4. **Given** the same plan and the same options, **When** the tool is run twice, **Then** both
   outputs are byte-for-byte identical.
5. **Given** a diagram printed in greyscale or viewed by a user who cannot distinguish the
   status colors, **When** the user looks at it, **Then** each status is still identifiable
   from its symbol and its border style alone.
6. **Given** a plan with no resource changes, or only unchanged resources and no filter asking
   for them, **When** the diagram is produced, **Then** a valid SVG image is still produced
   and it tells the viewer that the plan contains nothing to show.

---

### User Story 2 - Understand how resources depend on each other (Priority: P2)

The reviewer wants to see which resources reference which others (for example a subnet that
belongs to a network, an instance that sits in a subnet), so they can judge the impact of a
change, not only its list. The diagram draws arrows between dependent resources.

**Why this priority**: Dependencies turn a grouped list into an architecture diagram. They
need the grouping from story 1 but are not required for it to be useful.

**Independent Test**: Run the tool with the `svg` format on a plan where resources reference
each other, and check that an arrow links each dependent resource to the resource it depends
on, and that resources without relations are laid out without arrows.

**Acceptance Scenarios**:

1. **Given** a plan where a resource references another through its configuration or through
   an explicit dependency, **When** the diagram is produced, **Then** an arrow links the two
   resources.
2. **Given** a resource created several times (repeated by count or by key), **When** it is
   referenced by another resource, **Then** the relation is drawn for every instance.
3. **Given** a resource that references the output of a module, **When** the diagram is
   produced, **Then** the arrow points to the container of that module.
4. **Given** a resource to delete, **When** the diagram is produced, **Then** it is shown
   without arrows, even if it was related to other resources before the change.
5. **Given** resources whose references form a cycle, **When** the diagram is produced, **Then**
   it is still produced, without failure and without endless processing, and remains readable.
6. **Given** a plan with no information about relations between resources, **When** the
   diagram is produced, **Then** resources are shown in a regular arrangement without arrows.
7. **Given** a new resource that references an unchanged resource (for example a new subnet in an
   existing network), **When** the diagram is produced with no type or action
   filter, **Then** the unchanged resource is drawn muted, in its own provider and module
   containers, with an arrow from the new resource, and no other unchanged resource is drawn.

---

### User Story 3 - Recognize resources by their cloud service icon (Priority: P3)

The reviewer recognizes services faster from an icon than from a resource type name. Each
resource shows an icon matching its kind (compute, network, storage, database, security, or
a specific service for the common AWS, Azure and GCP services). A user who has the official
icon packs can supply them to replace the built-in icons.

**Why this priority**: Icons improve recognition but the diagram is already useful with plain
labeled boxes. They can be delivered after stories 1 and 2.

**Independent Test**: Produce a diagram from a plan containing common AWS, Azure and GCP
resources and from one containing an unknown resource type, with and without an icon pack,
and check which icon each resource receives.

**Acceptance Scenarios**:

1. **Given** a resource of a common AWS, Azure or GCP service, **When** the diagram is
   produced, **Then** the resource shows the icon of that service family.
2. **Given** a resource of a type the tool does not know, **When** the diagram is produced,
   **Then** the resource is still drawn, with an icon chosen by provider prefix, then by
   general category, then a generic one.
3. **Given** the user points to a directory of icons, **When** the diagram is produced,
   **Then** resources whose family has a matching icon in that directory use it, and the
   others keep the built-in icon.
4. **Given** an icon appears on many resources, **When** the diagram is produced, **Then** the
   icon content is included once in the file and reused, so the file does not grow with the
   number of resources.
5. **Given** the icon directory does not exist or contains unusable files, **When** the diagram
   is produced, **Then** the user gets a clear message for an invalid directory, and a
   single unusable file does not prevent the diagram: the built-in icon is used instead.

---

### User Story 4 - Focus and share the diagram safely (Priority: P4)

The user restricts the diagram to some resource types or some actions with the existing
filters, and shares the resulting image in a ticket or a pull request without leaking
attribute values.

**Why this priority**: Filters already exist in the tool and must simply be honored, and
the no-leak rule is a constraint on stories 1 to 3 rather than a separate capability, but it
deserves its own verification.

**Independent Test**: Run the tool with the `svg` format together with a type filter and an
action filter on a plan with sensitive attributes, then check the resources shown and search
the output for attribute values.

**Acceptance Scenarios**:

1. **Given** a type filter or an action filter, **When** the diagram is produced, **Then**
   only the matching resources are drawn, no dependency context is added, and the groupings and
   arrows between the remaining resources stay correct.
2. **Given** a plan whose resources carry attribute values, including values marked as
   sensitive, **When** the diagram is produced, **Then** no attribute value of the plan appears
   anywhere in the image; only addresses, types, names, providers, modules and statuses do.
3. **Given** a plan with several hundred resources, **When** the diagram is produced, **Then**
   it is produced in a reasonable time and the image size adapts to its content.

---

### Edge Cases

- A plan with no resource changes, or with only unchanged resources, produces a valid image
  with an explicit message, and never an error.
- Resource names, module keys and icon files can contain characters that are special in SVG
  (such as `<`, `>`, `&`, quotes): the output remains a well-formed image and nothing in it
  can be interpreted as active content.
- An icon pack file contains scripts or external references: the produced image must not
  execute or load anything from it. Such a file is rejected and the built-in icon is used.
- A resource of an unknown provider (for example `random`) or of an unknown type is still
  drawn, with the generic icon.
- A resource that is only moved or imported, without an update, has the no-op status: it is
  hidden unless the action filter asks for no-op or a changed resource depends on it. A moved or
  imported resource that also changes is drawn once, with the status of that change.
- A plan produced by an older format version, with no information about relations: the
  diagram is produced without arrows.
- Resources to replace are drawn as one resource with the replace status, not as two.
- Data sources (read actions) are not drawn, consistently with the other formats.
- Several changed resources depend on the same unchanged resource: it is drawn once as
  context, with one arrow per dependent resource.
- Very long addresses are shortened or wrapped in the image and remain readable; the full
  address is still available to a viewer that supports it (for example as a tooltip).
- Two resources with the same display name in different modules never collide.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tool MUST accept `svg` as an output format value, from the command line,
  from the configuration file and from the environment, like the existing formats.
- **FR-002**: With the `svg` format the tool MUST write a single well-formed SVG image to its
  normal output, and MUST NOT mix other text into it.
- **FR-003**: Without an action filter, every resource whose action is create, update, delete
  or replace MUST appear in the image exactly once, and no-op resources MUST NOT appear,
  except as dependency context under FR-024. A no-op resource MUST also appear, exactly once,
  when the action filter explicitly asks for the no-op action. Replaced resources MUST appear
  as one resource.
- **FR-004**: Resources MUST be grouped in nested containers following their provider and
  then their module path; resources in the root module appear directly in their provider
  container.
- **FR-005**: Each status MUST be conveyed by three redundant cues: a color, a symbol and a
  border style. Create is green with `+`, update is orange with `~`, delete is red with `-`,
  replace is purple with `-/+`, and no-op, when it is requested, is grey and visually muted.
- **FR-006**: Relations between resources MUST be drawn as arrows, deduced from the references
  between resources and the explicit dependencies declared in the plan.
- **FR-007**: A relation declared on a resource that is repeated (count or key) MUST apply to
  every instance of that resource.
- **FR-008**: A reference to a module output MUST be drawn towards the container of that
  module; resolving it to the resource that produces the output is not required.
- **FR-009**: Resources to delete MUST be drawn without any arrow.
- **FR-010**: When the plan carries no relation information, or a resource has no relation, the
  layout MUST fall back to a regular arrangement without arrows.
- **FR-011**: The arrangement MUST be automatic, readable and deterministic: related
  resources are placed in successive columns following their dependencies, ties are
  ordered alphabetically by address, the image size is computed from its content, and
  cycles of references are handled without failure or endless processing.
- **FR-012**: Every resource MUST show an icon chosen by the resource type: first from the
  list of common AWS, Azure and GCP services, then by provider prefix and general category
  (compute, network, storage, database, security), then a generic icon. A type the tool does
  not know MUST NOT cause an error.
- **FR-013**: The tool MUST include built-in generic icons so that a diagram needs no external
  file. Official cloud provider icons MUST NOT be shipped with the tool in this version.
- **FR-014**: The user MUST be able to provide a directory of icons with an `--icons-dir`
  option; icons found there take precedence over the built-in ones for the families they
  match. With any format other than `svg`, `--icons-dir` MUST be ignored without failing, and
  a notice MUST be written to the error output only in verbose mode.
- **FR-015**: An icon used by several resources MUST be included once in the image and
  referenced by the others.
- **FR-016**: Icon files supplied by the user MUST be validated; a file that contains active
  content or external references MUST be ignored for the diagram, and the built-in icon
  used instead. A missing or unreadable icon directory MUST give an actionable error.
- **FR-017**: The image MUST NOT contain any attribute value (before, after, unknown or
  sensitive) of the plan; it MAY contain resource addresses, types, names, providers,
  modules and statuses only.
- **FR-018**: All text drawn in the image MUST be escaped so that arbitrary resource names
  cannot break the image or inject content.
- **FR-019**: The existing type and action filters MUST apply to the diagram; the containers
  and arrows of the remaining resources MUST stay consistent after filtering. When a type or
  action filter is given, only the selected resources are drawn, no dependency context is
  added, and arrows towards resources that are not drawn are omitted.
- **FR-020**: A plan with nothing to show MUST produce a valid image with an explicit message
  and a success exit status.
- **FR-021**: The `svg` format MUST read everything it needs from the plan file, MUST work
  offline and MUST NOT require the user to install anything beyond the tool itself.
- **FR-022**: Existing formats (markdown, text, json) and their outputs MUST NOT change.
- **FR-023**: The documentation MUST describe the `svg` format, the `--icons-dir` option, the
  status legend and the limits, and the `-h` helper output MUST be regenerated.
- **FR-024**: When neither a type filter nor an action filter is given, an unchanged (no-op)
  resource that a drawn changed resource depends on MUST be drawn as muted context, exactly
  once, in its own provider and module containers, with the arrow from the resource that
  depends on it. Only direct dependencies are added: the dependencies of a context resource
  are not followed, and no arrow starts from it.

### Key Entities

- **Diagram**: The picture of one plan: the containers, the resource nodes, the arrows and
  a legend, independent of how it is drawn so that other renderings can reuse it.
- **Container**: A provider or a module that groups resources and other containers; it has a
  label and a nesting level.
- **Resource node**: One resource instance with its address, type, name, status, icon
  family and its container. A node is a context node when it is unchanged and drawn only
  because a changed resource depends on it.
- **Status**: One of create, update, delete, replace, no-op, with its color, symbol and border
  style.
- **Relation**: A directed link from a resource to a resource or a module container it
  depends on.
- **Icon family**: A kind of service (specific service, provider, category or generic) that
  determines which icon a resource displays.
- **Icon pack**: An optional user-provided directory of icons overriding the built-in ones.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For a plan mixing the five statuses, with unchanged resources requested, a
  reviewer can state the status of every drawn resource from the image alone, including in a
  greyscale rendering, with no mistake.
- **SC-002**: 100% of the resources selected by the filters appear exactly once in the image,
  and no resource that the filters exclude appears.
- **SC-003**: Two runs on the same plan and the same options produce identical output in
  100% of the cases.
- **SC-004**: Every plan that the tool accepts today for the other formats produces a valid image
  with the `svg` format, including plans with unknown providers, unknown types, no
  relations, cycles and no changes.
- **SC-005**: Searching the image for any attribute value of the test plans, including values
  marked as sensitive, finds none.
- **SC-006**: A plan of 300 resources is turned into an image in under 2 seconds on a
  standard developer machine.
- **SC-007**: A reviewer finds the container (provider and module) of a given resource in the
  image in under 10 seconds for a plan of up to 50 resources.
- **SC-008**: The output of the existing formats is unchanged for the whole existing test and
  example corpus.

## Assumptions

- "OpenSpec" in the original request is read as OpenTofu: both Terraform and OpenTofu plans in
  the JSON form already accepted by the tool are in scope.
- The image is written to the normal output like the other formats; writing to a file
  relies on the existing shell redirection, and no new output-file option is added in this
  version.
- Data sources (read actions) are not drawn, consistently with how the summary treats them.
- Unchanged (no-op) resources are drawn only on explicit request through the existing action
  filter, or, when no filter is given, as one-hop dependency context of a drawn changed
  resource, so the changes stay the focus of the diagram, as in the plan summary.
- Instances of a repeated resource are drawn individually, each with its own status.
- An icon family is identified by a short name and an icon pack provides one SVG file per
  family, named after the family; the exact naming scheme is fixed during planning.
- The status legend is part of the image.
- Out of scope for this version: isometric rendering, PNG export, resolving module outputs to
  the resource that produces them, and recovering the relations of resources to delete.
  The diagram model is kept separate from the drawing so that an isometric rendering can
  reuse it later.
- Official AWS, Azure and GCP icons are not committed or shipped until their license terms
  are checked; this version ships generic icons only.
- Existing fixtures contain almost no relations between resources, so new synthetic fixtures
  for AWS, Azure and GCP (network, subnet, instance) are needed for the tests, in line with
  the rule that committed fixtures contain synthetic data only.
- The project constitution lists the supported output formats; adding `svg` requires updating
  that list in the same change.
