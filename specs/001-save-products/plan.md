# Implementation Plan: Save products and view them later

**Branch**: `001-save-products` | **Date**: 2026-06-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-save-products/spec.md` (Jira AIP-192, epic AIP-94)

## Summary

Let a shopper save a product from its page without signing in, and return to a
header-linked saved-products view that lists what they've saved (with click-through
to each product) or shows a clear empty state. Saved state is scoped to the shopper's
existing browsing session and held **in memory inside the frontend service** — no new
service and no new datastore — with product details resolved through the existing
product catalogue. The implementation mirrors the store's existing add-to-cart flow
(POST form → handler keyed by `sessionID(r)` → redirect) and the existing
cart-size/header patterns, in idiomatic Go.

## Technical Context

**Language/Version**: Go 1.25.0 — changes confined to the `src/frontend` service.

**Primary Dependencies**: All existing — `gorilla/mux` (routing), `html/template`
(server-rendered pages), `logrus` (logging), and the existing product catalogue
service accessed via the frontend's current helper (`fe.getProduct`). No new
dependencies.

**Storage**: In-memory, process-local map keyed by session ID, living in the frontend
process. No database, cache, or search engine (Principle II). The product catalogue
(`src/productcatalogservice/products.json`) remains the single source of product data,
read through the existing catalogue service; it is not modified.

**Testing**: Go `testing` with `net/http/httptest`, matching the repo's existing Go
test style (e.g. `productcatalogservice/product_catalog_test.go`). Unit tests for the
in-memory store; handler tests for save and view.

**Target Platform**: Linux container, existing frontend deployment; server-rendered
HTML pages.

**Project Type**: Web service (single server-rendered Go frontend within the
microservices demo). Single-project structure — no new project.

**Performance Goals**: Parity with existing pages. Save/view are O(1)/O(n-saved)
in-memory operations; no measurable added latency.

**Constraints**: Bound by Constitution Principles I–V (below). **Known limitation:**
because state is in memory per Principle II, saved products are *process-local* — not
shared across frontend replicas. See research.md; acceptable for the session-scoped,
demo-scale target, but recorded explicitly rather than hidden.

**Scale/Scope**: Demo scale; small per-session collections.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle | Status | How this plan complies |
|-----------|--------|------------------------|
| **I. Existing Services Only** | ✅ PASS | All changes live in the existing `frontend` service; product details come from the existing product catalogue service. No new service. |
| **II. No New Datastore** | ✅ PASS (with noted limitation) | Saved state is an in-memory map in the frontend; no DB/cache/search added. Product data still sourced from `products.json` via the catalogue. Process-local nature documented in research.md. |
| **III. Match the Service You Change** | ✅ PASS | Idiomatic Go mirroring existing patterns: a POST form like add-to-cart, a `mux` route in `main.go`, an `html/template` view, `sessionID(r)` for identity. |
| **IV. Infra/Deploy/CI Frozen** | ✅ PASS | Only Go source and templates change. New routes are registered in application code (`main.go`), not infrastructure. No Kubernetes/Helm or CI edits. |
| **V. Ship a Vertical Slice** | ✅ PASS | Delivers save + saved view + header link + empty state end to end — demonstrable to a shopper on its own. |

**Result**: PASS — no violations. Complexity Tracking left empty.

## Project Structure

### Documentation (this feature)

```text
specs/001-save-products/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (HTTP/UI contracts)
│   └── wishlist-routes.md
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

All changes are within the existing frontend service; no new top-level structure.

```text
src/frontend/
├── main.go              # CHANGE: register routes  GET/POST /wishlist, POST /wishlist/remove (remove is a later story; route added only when needed)
├── handlers.go          # CHANGE: add saveProductHandler + viewWishlistHandler; set saved-state flag on product page
├── wishlist_store.go     # NEW: in-memory, session-keyed saved-products store (add / list / contains)
├── wishlist_store_test.go# NEW: unit tests for the store
└── templates/
    ├── product.html     # CHANGE: add "Save" form + saved/unsaved indicator (mirrors the add-to-cart form)
    ├── header.html      # CHANGE: add a saved-products link (no count in this story)
    └── wishlist.html    # NEW: saved-products view + empty state
```

**Structure Decision**: Single existing service (`src/frontend`). This story adds one
new file (the in-memory store) plus one new template, and edits three existing files,
keeping strictly within the frontend per Principles I, III, and IV. Saving from the
list, removal, and the count badge belong to AIP-193/194/195 and are out of scope here
(the `/wishlist/remove` route is noted but not added by this story).

## Complexity Tracking

> No constitution violations — nothing to justify.
