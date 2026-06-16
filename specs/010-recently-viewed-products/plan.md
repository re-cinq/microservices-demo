# Implementation Plan: Recently Viewed Products (Product Page)

**Branch**: `010-recently-viewed-products` | **Date**: 2026-06-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/010-recently-viewed-products/spec.md` (derived from Jira [AIP-199](https://odevo.atlassian.net/browse/AIP-199) under epic [AIP-97](https://odevo.atlassian.net/browse/AIP-97))

## Summary

Add a "Recently viewed" strip to the product page. As a shopper views products, the frontend records each viewed product's ID in an in-memory, per-session list. On a product page, the strip renders the most-recently-viewed products (most-recent-first, current product excluded, capped at 4) as half-size thumbnails at the bottom of the page; selecting one navigates to that product. When nothing has been viewed, the page is unchanged.

The implementation lives entirely in the existing **frontend** Go service, reusing the established "recommendations strip" pattern: an ordered list of product IDs hydrated to full products via the existing `productCatalogService` `GetProduct` RPC, then rendered through an html/template partial. No new service, datastore, dependency, or infrastructure change is introduced — satisfying the epic's hard constraints.

## Technical Context

**Language/Version**: Go 1.25.0 (existing `src/frontend` module — no toolchain change)

**Primary Dependencies**: `gorilla/mux` (routing), `html/template` (server-side rendering), existing `ProductCatalogServiceClient` gRPC stub — all already in use. **No new dependencies.**

**Storage**: In-process, per-session map held on the `frontendServer` struct (`map[sessionID][]productID`, mutex-guarded). No database, cache, or search engine — per constraint C-002. Catalogue data continues to come from the existing `productcatalogservice` (backed by `products.json`).

**Testing**: `go test` unit tests for the in-memory store (ordering, cap, exclusion, dedup-absence); manual UI verification via quickstart.

**Target Platform**: Linux container, existing frontend deployment (unchanged manifests/pipeline — per C-004).

**Project Type**: Web service — server-side-rendered Go frontend within the Online Boutique microservices demo.

**Performance Goals**: Negligible. The strip adds at most 4 `GetProduct` RPCs per product-page render — identical in shape to the existing recommendations strip. Record-on-view is an O(1) in-memory operation.

**Constraints** (hard, from epic AIP-97):
- C-001: Use only existing services; no new services.
- C-002: No new datastore (no DB/cache/search); work in memory over the existing catalogue (`productcatalogservice/products.json`).
- C-003: Match the language/patterns of the service changed; frontend and product catalogue are Go.
- C-004: No changes to infrastructure, deployment manifests, or CI; ships through the existing pipeline unmodified.

**Scale/Scope**: Demo scale. One service touched (frontend). One new Go source file + one new template partial + one CSS rule + edits to `productHandler` and the `frontendServer` struct.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` contains only unpopulated template placeholders — there are **no ratified project-specific principles or gates** to enforce. Applying the default spirit (simplicity / YAGNI / no unjustified complexity):

| Default principle | Status | Notes |
|---|---|---|
| Simplicity / YAGNI | ✅ Pass | Single in-memory map + one template partial; reuses existing strip pattern. No new abstractions. |
| No new infrastructure | ✅ Pass | Frontend-only; C-004 honoured. |
| Match existing patterns | ✅ Pass | Mirrors `recommendations` strip and `sessionID` usage. |

**Result**: PASS (no violations; Complexity Tracking not required).

## Project Structure

### Documentation (this feature)

```text
specs/010-recently-viewed-products/
├── plan.md              # This file
├── spec.md              # Feature spec (from AIP-199)
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── recently-viewed.md
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

All changes are confined to the frontend service:

```text
src/frontend/
├── main.go                         # [EDIT] add recentlyViewed store field to frontendServer; init in main()
├── handlers.go                     # [EDIT] productHandler: record view + build strip view-model
├── recentlyviewed.go               # [NEW]  in-memory per-session store (record, list)
├── recentlyviewed_test.go          # [NEW]  unit tests for the store
├── templates/
│   ├── recentlyviewed.html         # [NEW]  strip partial (mirrors recommendations.html)
│   └── product.html                # [EDIT] include the partial at the bottom of the page
└── static/styles/
    └── styles.css                  # [EDIT] one rule: half-size thumbnails for the strip
```

**Structure Decision**: Single-service change inside `src/frontend`. No new packages or services. The store is a small struct on `frontendServer` (same lifecycle as the existing gRPC client fields). Rendering reuses the existing template-partial mechanism (`{{ template ... }}`) already used for `recommendations` at `product.html:77`.

> **Scope note**: This story is the product-page strip only. The home-page strip is [AIP-200](https://odevo.atlassian.net/browse/AIP-200) and de-duplication/move-to-front is [AIP-201](https://odevo.atlassian.net/browse/AIP-201). To keep AIP-200 cheap, the template partial and the store list are built to be reusable (the home handler will reuse both unchanged).

## Complexity Tracking

No constitution violations; no entries required.
