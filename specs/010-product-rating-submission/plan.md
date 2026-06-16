# Implementation Plan: Rate a product and see it on the product page

**Branch**: `010-product-rating-submission` | **Date**: 2026-06-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/010-product-rating-submission/spec.md`

## Summary

Let a shopper submit a 1–5 star rating on the product page and immediately see the product's average (rendered as stars) plus the number of ratings. Ratings have no current data source, so this feature both creates the data (from shopper submissions) and displays it. Per epic AIP-95, the work is confined to the two existing Go services — `productcatalogservice` (holds the rating aggregates in memory and exposes a new `RateProduct` RPC plus two new `Product` fields) and `frontend` (adds a submit control and renders the rating) — with no new service, no datastore, and no infrastructure/CI changes.

## Technical Context

**Language/Version**: Go (matches `frontend` and `productcatalogservice`; both already on the repo's Go toolchain)  
**Primary Dependencies**: existing only — gRPC + protobuf (`genproto`), `gorilla/mux` (frontend routing), Go `html/template`  
**Storage**: In-memory only. Rating aggregates held in process memory in `productcatalogservice` behind a mutex. No database, cache, or search engine (HC-002).  
**Testing**: `go test` — existing `product_catalog_test.go` pattern in `productcatalogservice`; table-driven unit tests for the aggregate logic.  
**Target Platform**: Linux containers, shipped through the existing pipeline unmodified (HC-004).  
**Project Type**: Web — multi-service (gRPC microservices + Go HTML frontend).  
**Performance Goals**: No regression to existing product-page / catalogue latency; rating submission is a single in-memory update.  
**Constraints**: HC-001..HC-005 from the spec (existing services only; no datastore; match Go patterns; no infra/CI changes; no review text or moderation).  
**Scale/Scope**: 9 products in the catalogue; single-replica demo; unbounded but low submission volume.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The project constitution (`.specify/memory/constitution.md`) is an unratified template — all principles are unfilled placeholders. There are therefore **no project-specific gates to enforce**. In their place, this plan is gated on the epic's hard constraints (HC-001..HC-005), which act as the binding non-negotiables:

| Gate | Status | Notes |
|------|--------|-------|
| HC-001 No new services | PASS | Changes only `frontend` + `productcatalogservice`. |
| HC-002 No new datastore | PASS | In-memory aggregates, mutex-guarded. |
| HC-003 Match service language/patterns | PASS | Go, following existing gRPC + template patterns. |
| HC-004 No infra/deploy/CI changes | PASS | Proto regen + code only; ships through existing pipeline. |
| HC-005 No review text / moderation | PASS | Star value only; no text, no dedup. |

**Initial gate: PASS.** Re-evaluated post-design below.

## Project Structure

### Documentation (this feature)

```text
specs/010-product-rating-submission/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── productcatalog-rating.md
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

```text
protos/
└── demo.proto                         # ADD: rating fields on Product; RateProduct RPC + messages

src/productcatalogservice/
├── genproto/                          # REGEN from demo.proto (demo.pb.go, demo_grpc.pb.go)
├── product_catalog.go                 # ADD: in-memory rating store + RateProduct + populate fields on reads
├── server.go                          # (unchanged wiring; mutex already present)
└── product_catalog_test.go            # ADD: unit tests for aggregate + validation

src/frontend/
├── genproto/                          # REGEN from demo.proto
├── rpc.go                             # ADD: rateProduct client helper
├── handlers.go                        # ADD: rateProductHandler (POST); pass rating data to product page
├── main.go                            # ADD: POST route /product/{id}/rate
└── templates/
    └── product.html                   # ADD: star display + submit form
```

**Structure Decision**: Web multi-service layout (already in place). The change touches exactly the two Go services named in the epic plus the shared `protos/demo.proto` (and the regenerated `genproto` in each service). No new directories or services are introduced.

## Complexity Tracking

No constitution violations. No entries required.
