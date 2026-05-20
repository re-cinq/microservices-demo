# Implementation Plan: Save for Later / Wishlist

**Branch**: `005-save-wishlist` | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-save-wishlist/spec.md`

## Summary

Add a session-scoped wishlist to the Online Boutique frontend. Shoppers can save products from the product detail page and view them on a dedicated `/wishlist` page. No new services or datastores are introduced — wishlist state is held in an in-memory `sync.Map` on the existing `frontendServer` struct, keyed by the session ID already maintained by the frontend's cookie middleware. Product details are resolved at display time from the existing product catalogue gRPC connection.

## Technical Context

**Language/Version**: Go 1.25

**Primary Dependencies**: `github.com/gorilla/mux` (routing), `html/template` (server-side rendering), `github.com/sirupsen/logrus` (logging), gRPC to `productcatalogservice` and `currencyservice` (existing connections on `frontendServer`)

**Storage**: In-memory `sync.Map` field added to `frontendServer` — no new datastore

**Testing**: `go test ./...` (standard Go toolchain); new `wishlist_test.go` for handler unit tests

**Target Platform**: Linux container (existing Docker/Skaffold pipeline, unchanged)

**Project Type**: Web service (HTTP frontend, server-side rendered HTML)

**Performance Goals**: Save confirmation visible in under 1 second (SC-001); wishlist page renders up to 20 items without pagination (SC-002)

**Constraints**: No new services; no new datastores; no infra/CI changes; match Go patterns of `src/frontend/`

**Scale/Scope**: Single service change (`src/frontend/`); 5 files modified or created

## Constitution Check

*Constitution not yet ratified for this project — template in place, no gates defined. No violations to evaluate.*

## Project Structure

### Documentation (this feature)

```text
specs/001-save-wishlist/
├── plan.md                        # This file
├── research.md                    # Phase 0: storage, confirmation, nav decisions
├── data-model.md                  # Phase 1: WishlistStore + WishlistItem entities
├── quickstart.md                  # Phase 1: how to run and test
├── contracts/
│   └── http-endpoints.md          # Phase 1: POST /wishlist/save, GET /wishlist contracts
└── tasks.md                       # Phase 2 output (/speckit-tasks — not yet created)
```

### Source Code

```text
src/frontend/
├── main.go                  # +wishlists sync.Map field on frontendServer; +2 routes
├── handlers.go              # +saveWishlistHandler; +viewWishlistHandler; productHandler reads ?saved=1
└── templates/
    ├── product.html         # +Save for later form button; +conditional confirmation banner
    ├── wishlist.html        # NEW — wishlist view (product cards + empty state)
    └── header.html          # +wishlist icon nav link (mirrors cart link pattern)
```

## Complexity Tracking

*No constitution violations — table not required.*
