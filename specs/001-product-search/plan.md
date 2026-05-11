# Implementation Plan: Product Search

**Branch**: `attendee/kate-payne-fix` | **Date**: 2026-05-11 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `specs/001-product-search/spec.md`

---

## Summary

Shoppers cannot find products without scrolling the full catalogue. This plan
adds a search input to the home/catalogue page (`/`) that filters the product
list by name and description, highlights matching substrings, shows a clear
empty state, and blocks queries shorter than 2 non-whitespace characters.

The `SearchProducts` gRPC RPC and its proto messages already exist in the repo.
The only productcatalogservice change is extending the filter to include
`Description` alongside `Name`. All remaining work is in the frontend service.

---

## Technical Context

**Language/Version**: Go 1.25 (frontend), Go 1.25 (productcatalogservice)  
**Primary Dependencies**: gorilla/mux (routing), html/template (templating), gRPC/protobuf (service calls), Bootstrap 4 (CSS — already in use)  
**Storage**: In-memory product catalogue loaded from `products.json` at startup  
**Testing**: Go standard `testing` package; existing `product_catalog_test.go` pattern  
**Target Platform**: Kubernetes (existing deployment, no new infra)  
**Project Type**: Web service (HTTP frontend) + gRPC backend service  
**Performance Goals**: Results visible within 5 seconds on standard connection (SC-001); in-memory filter is sub-millisecond  
**Constraints**: No new services, no external search engines, no new env vars, no proto changes, no build pipeline changes  
**Scale/Scope**: 10 products in catalogue; filter runs in-process on every search request

---

## Constitution Check

The project constitution (`CONSTITUTION.md`) is a blank template — no project-
specific principles have been ratified. No gates to evaluate.

*Re-check after Phase 1 design*: no new abstractions, no new services, no new
dependencies introduced. Design is consistent with existing patterns throughout
the codebase.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-product-search/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output
├── contracts/
│   └── search-rpc.md   ← Phase 1 output
├── quickstart.md        ← Phase 1 output
└── tasks.md             ← Phase 2 output (/speckit-tasks, NOT created here)
```

### Source Code (repository root)

```text
src/productcatalogservice/
└── product_catalog.go   ← fix SearchProducts to also filter on Description

src/frontend/
├── main.go              ← no change needed (homeHandler already registered)
├── handlers.go          ← extend homeHandler to handle ?q= param; add
│                           highlightTerm template function to FuncMap
├── rpc.go               ← add searchProducts() gRPC wrapper
└── templates/
    └── home.html        ← add search form + conditional results/empty state
```

No new files needed in either service. No proto changes. No new routes.

**Structure Decision**: Unified single-page approach — search form and results
both live on `/` via query param `?q=`. This avoids a new route, keeps the
existing `homeHandler` as the single entry point for the catalogue, and
naturally preserves the search term in the browser URL bar.

---

## Complexity Tracking

No constitution violations. No complexity justification required.
