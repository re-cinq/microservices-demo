# Online Boutique Feature Constitution

This constitution governs features built into the Online Boutique demo store for the
AI product training (epics under the Jira `AIP` project, label `m6-epic`). Its
principles are **hard constraints**: every `/speckit-plan` and `/speckit-implement`
MUST be checked against them, and any violation blocks the work until resolved.

## Core Principles

### I. Existing Services Only (NON-NEGOTIABLE)

Features MUST be built using only the services that already exist in this repository.
Adding a new service is prohibited. Work that appears to need a new service MUST be
re-scoped into an existing service (typically the frontend or the product catalogue
service) or escalated as out of scope — never solved by introducing a new deployable.

### II. No New Datastore (NON-NEGOTIABLE)

No new datastore of any kind may be introduced — no database, no cache, no search
engine, no external persistence. Feature state MUST be held in memory within an
existing service, working over the existing product catalogue in
`src/productcatalogservice/products.json`. Any requirement for durable or
cross-instance persistence is out of scope unless the constitution is amended.

### III. Match the Service You Change

Changes MUST follow the language, structure, and patterns of the service being
modified — no new languages, frameworks, or architectural styles. The frontend and
the product catalogue service are written in Go; changes there MUST be idiomatic Go
consistent with the surrounding code (routing, handlers, templates, naming).

### IV. Infrastructure, Deployment, and CI Are Frozen

Features MUST NOT change infrastructure, deployment manifests (Kubernetes/Helm), or
CI configuration. The change MUST ship through the existing pipeline unmodified. If a
feature appears to require an infra, manifest, or CI change, that is a signal the
feature is out of scope as specified — stop and surface it rather than editing those
files.

### V. Ship a Vertical Slice

Each story MUST deliver an independently shippable, customer-visible slice (UI +
behaviour + in-memory state end to end), not a horizontal or purely technical layer.
A story that cannot be demonstrated to a shopper on its own does not satisfy this
principle and MUST be re-cut.

## Constraint Provenance

These principles are transcribed from the **Technical constraints** block carried in
every Module 6 epic description (Jira `AIP`, label `m6-epic`), e.g. epic AIP-94
*Wishlist / Save for later*. They apply uniformly to every story under those epics.
Specs remain technology-agnostic; these constraints bind the **plan** and
**implementation** phases, where the "how" is decided.

## Governance

- This constitution supersedes convenience and individual preference. Where a plan or
  implementation conflicts with a principle, the principle wins.
- `/speckit-plan` MUST include a Constitution Check confirming the plan honours
  Principles I–V, and MUST record any tension explicitly rather than proceeding
  silently.
- Amendments require an explicit decision, a version bump per the policy below, and an
  update to the "Last Amended" date. Amendments that relax a NON-NEGOTIABLE principle
  must state the justification.
- Versioning policy: **MAJOR** for removing or redefining a principle in a
  backward-incompatible way; **MINOR** for adding a new principle or section; **PATCH**
  for clarifications and wording that do not change meaning.

**Version**: 1.0.0 | **Ratified**: 2026-06-16 | **Last Amended**: 2026-06-16
