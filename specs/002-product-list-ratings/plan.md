# Implementation Plan: Star ratings on the product list

**Branch**: `attendee/daniel-tufvander` (work continues on this branch; spec dir `002-product-list-ratings` is independent of branch name) | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/002-product-list-ratings/spec.md`
**Source story**: [AIP-161](https://odevo.atlassian.net/browse/AIP-161) under epic [AIP-95](https://odevo.atlassian.net/browse/AIP-95)

## Summary

Add a `rating` field to every product in the Online Boutique catalogue and surface it as stars on every product card on the home page (the product list). The change is a vertical slice through three layers: the catalogue data (`products.json`), the gRPC contract (`protos/demo.proto`'s `Product` message), and the frontend rendering (`src/frontend/handlers.go` + `templates/home.html`). No new services, no new datastores, no infrastructure changes — strictly within constraints C-001…C-004 from the spec.

The sibling story [AIP-162](https://odevo.atlassian.net/browse/AIP-162) (rating on the product page) reuses the proto/data work from this plan and only adds rendering in `templates/product.html`. That story is explicitly out of scope here.

## Technical Context

**Language/Version**: Go (matching the versions already pinned in `src/productcatalogservice/go.mod` and `src/frontend/go.mod`; no upgrade needed).
**Primary Dependencies**: existing only — `google.golang.org/grpc`, the in-repo `genproto` package, the existing `html/template` rendering in the frontend. No new module added.
**Storage**: none added. Rating values are seeded into the existing `src/productcatalogservice/products.json` file and held in memory via the existing `productCatalog.catalog` field. Strictly satisfies C-002.
**Testing**: extend the existing `src/productcatalogservice/product_catalog_test.go` with a unit test that asserts the rating is loaded and returned via `ListProducts`. No new test framework. Manual frontend verification covered by `quickstart.md`.
**Target Platform**: existing cohort deployment — `https://daniel-tufvander.training.gcp.re-cinq.com`. No platform change.
**Project Type**: web-service (frontend Go service + gRPC backend services). Established layout under `src/`.
**Performance Goals**: no measurable regression versus today's product list render. The change adds one float field per product and a static stars render per card — negligible.
**Constraints**: hard constraints C-001…C-004 from the spec — no new services, no new datastore, Go only, no infra/CI changes. The proto regeneration is an in-repo code change (runs `genproto.sh` locally and commits the result); it is **not** infrastructure or CI.
**Scale/Scope**: catalogue size is the existing ~10 products; no fan-out or batch concerns.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`/Users/danieltufvander/Lovable app Wheel/microservices-demo/.specify/memory/constitution.md` is the unfilled Spec Kit template — every principle is still a `[PRINCIPLE_*_NAME]` / `[PRINCIPLE_*_DESCRIPTION]` placeholder, and the version is `[CONSTITUTION_VERSION]`. There are no ratified principles to gate against.

**Initial gate**: vacuously passes (no principles defined).
**Post-design gate**: vacuously passes (re-evaluated below after Phase 1).

This is recorded explicitly so a reader doesn't mistake the empty constitution for "the plan ignored governance."

## Project Structure

### Documentation (this feature)

```text
specs/002-product-list-ratings/
├── plan.md              # This file (/speckit-plan output)
├── spec.md              # Feature spec (already created by /speckit-specify)
├── research.md          # Phase 0 output — research & decision log
├── data-model.md        # Phase 1 output — rating field on Product entity
├── quickstart.md        # Phase 1 output — how to run / verify locally
├── contracts/
│   ├── service-contract.md   # Phase 1 — gRPC Product message change
│   └── ui-contract.md        # Phase 1 — product-list rating widget contract
├── checklists/
│   └── requirements.md  # Already created by /speckit-specify validation step
└── tasks.md             # Phase 2 output — created by /speckit-tasks (NOT this command)
```

### Source Code (repository root)

This feature touches three existing locations only. No new directories.

```text
protos/
└── demo.proto                                 # add `float rating = 7;` to message Product

src/productcatalogservice/
├── products.json                              # seed a `rating` field on every product
├── genproto/                                  # regenerated from demo.proto via genproto.sh
└── product_catalog_test.go                    # extend with rating coverage

src/frontend/
├── genproto/                                  # regenerated from demo.proto via genproto.sh
├── handlers.go                                # extend productView with Rating; pass it through homeHandler
├── templates/home.html                        # render stars in the existing hot-product-card block
└── static/styles/                             # add styles for the rating widget (existing dir; no new asset pipeline)
```

**Structure Decision**: extend the existing web-service layout in place. Two services are touched (`productcatalogservice` for the data + contract; `frontend` for the rendering). The shared `protos/demo.proto` is the source of truth for the contract; regenerated `genproto/*.pb.go` files in each service are committed alongside the proto change. No new service, no new directory, no new module. This is the smallest layout that satisfies the spec's hard constraints.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations to justify — constitution is unfilled, gates are vacuous, and the design takes the simplest path through the existing layers.

## Phase 0 — Outline & Research

See [research.md](research.md). The PRD/spec left six items open (A1–A6 in the spec's Assumptions). research.md works through each as a *Decision / Rationale / Alternatives* triple and resolves any technical unknowns about how to extend the proto contract, the data file, and the template without violating C-001…C-004.

## Phase 1 — Design & Contracts

Outputs:

- [data-model.md](data-model.md) — `Product` entity gains a `rating` field (float, 0.0–5.0, half-star granularity). Validation rules and seeding strategy captured.
- [contracts/service-contract.md](contracts/service-contract.md) — the proto change as a versioned contract: new field number, default-value behaviour, backward-compatibility note for any unmigrated consumer.
- [contracts/ui-contract.md](contracts/ui-contract.md) — the product-list rating widget: position on card, render rules for whole/half/empty stars, no-rating fallback, what does *not* change.
- [quickstart.md](quickstart.md) — concrete local steps: edit proto → regen → seed products.json → run skaffold/locally → verify on the cohort URL.

## Phase 1 — Re-evaluate Constitution Check (post-design)

Constitution is still unfilled; gates are still vacuous. The post-design re-check does not change the initial verdict. If a real constitution is ratified later, this section must be revisited.
