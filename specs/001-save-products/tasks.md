---
description: "Task list for: Save products and view them later"
---

# Tasks: Save products and view them later

**Input**: Design documents from `specs/001-save-products/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: INCLUDED — the plan (Testing) and research.md (Decision 7) explicitly call for
Go unit tests for the store and `net/http/httptest` handler tests. Written test-first.

**Organization**: One user story (US1, P1). All work lives under the existing
`src/frontend` service per the constitution (no new service/datastore/infra).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[US1]**: Belongs to User Story 1

## Path Conventions

- All paths are under the existing frontend service: `src/frontend/`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish a green baseline before changes.

- [ ] T001 Establish baseline: from `src/frontend`, confirm `go build ./...` and `go test ./...` succeed before any edits.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Cross-story prerequisites.

**Note**: This feature has a single user story, so there are **no foundational-only
tasks** — the in-memory saved-products store is created within US1 (it serves only this
story). Proceed directly to Phase 3.

---

## Phase 3: User Story 1 - Save a product and come back to it (Priority: P1) 🎯 MVP

**Goal**: A shopper can save a product from its page (no login), reach a header-linked
saved-products view listing what they've saved with click-through, see saved/unsaved
state on the product page, and get a clear empty state — all session-scoped.

**Independent Test**: Save a product from its page, browse elsewhere, open the saved
view from the header, confirm the product is listed and links back — no sign-in, within
one session. Empty state shows for a fresh session.

### Tests for User Story 1 (write first; ensure they FAIL before implementation) ⚠️

- [ ] T002 [P] [US1] Write failing unit tests for the in-memory store in `src/frontend/wishlist_store_test.go`: Add is idempotent (no duplicate, FR-007), List returns saved IDs in save order, Contains reflects saved state (FR-002), and collections are isolated per session ID.
- [ ] T003 [P] [US1] Write failing handler tests in `src/frontend/wishlist_handlers_test.go` using `net/http/httptest`: `POST /wishlist` with `product_id` saves and 302-redirects (FR-001); `GET /wishlist` renders saved products with links when present (FR-003/FR-005) and the empty-state message when none (FR-008).

### Implementation for User Story 1

- [ ] T004 [US1] Implement the in-memory store in `src/frontend/wishlist_store.go`: a mutex-guarded `map[sessionID]orderedSet` with `Add`, `List`, `Contains` per data-model.md; makes T002 pass. (depends on T002)
- [ ] T005 [US1] Add a `*wishlistStore` field to `frontendServer` and initialize it in `src/frontend/main.go`. (depends on T004)
- [ ] T006 [US1] Implement `saveProductHandler` in `src/frontend/handlers.go`: read `product_id` via `r.FormValue`, validate it resolves via the existing catalogue helper (`fe.getProduct`), call `store.Add(sessionID(r), id)`, and `http.Redirect` back to the product page — mirroring `addToCartHandler`. (depends on T004, T005)
- [ ] T007 [US1] Implement `viewWishlistHandler` in `src/frontend/handlers.go`: `store.List(sessionID(r))`, resolve each ID via `fe.getProduct` (skip IDs no longer in the catalogue), and render `wishlist.html` with the products and an empty-state flag. (depends on T004, T005)
- [ ] T008 [US1] In `productHandler` (`src/frontend/handlers.go`), pass a `product_saved` boolean from `store.Contains(sessionID(r), id)` into the product template data (FR-002). (depends on T004, T005)
- [ ] T009 [US1] Register routes in `src/frontend/main.go`: `POST {baseUrl}/wishlist` → `saveProductHandler`, `GET {baseUrl}/wishlist` → `viewWishlistHandler`. (depends on T006, T007)
- [ ] T010 [P] [US1] Create `src/frontend/templates/wishlist.html`: list saved products each linking to `{baseUrl}/product/{id}`, plus an empty-state message with a prompt to browse (FR-003/FR-005/FR-008).
- [ ] T011 [P] [US1] In `src/frontend/templates/product.html`, add a "Save" `POST` form (hidden `product_id`) to `{baseUrl}/wishlist` styled like Add To Cart, and show saved vs. unsaved using `product_saved` (FR-001/FR-002).
- [ ] T012 [P] [US1] In `src/frontend/templates/header.html`, add a saved-products link to `{baseUrl}/wishlist` present on every page — **no count** (count is AIP-195) (FR-004).
- [ ] T013 [US1] Run `go test ./...` in `src/frontend` and confirm the store and handler tests (T002, T003) now pass. (depends on T004–T012)

**Checkpoint**: User Story 1 fully functional and independently testable = the MVP.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [ ] T014 [P] Run `gofmt -w` and `go vet ./...` on the changed files in `src/frontend`.
- [ ] T015 Execute `specs/001-save-products/quickstart.md` manual validation (scenarios 1–6, incl. idempotency and the unavailable-product edge case).
- [ ] T016 Constitution spot-check: confirm `git status` shows changes only under `src/frontend/` (plus the spec docs) — no new service, no new datastore, no infra/manifest/CI edits (Principles I, II, IV).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies.
- **Foundational (Phase 2)**: none for this feature.
- **User Story 1 (Phase 3)**: starts after Setup.
- **Polish (Phase 4)**: after US1 complete.

### Within User Story 1

- Tests (T002, T003) written first and failing → store (T004) → server wiring (T005) →
  handlers (T006–T008) → routes (T009) → templates (T010–T012) → green test run (T013).
- Templates (T010–T012) are independent files and can be done in parallel with each
  other and alongside handler work; routes (T009) depend on the handlers existing.

### Parallel Opportunities

- T002 and T003 (separate test files) can be written in parallel.
- T010, T011, T012 (separate template files) can run in parallel.

---

## Parallel Example: User Story 1

```bash
# Write both failing test files together:
Task: "Unit tests for the store in src/frontend/wishlist_store_test.go"
Task: "Handler tests in src/frontend/wishlist_handlers_test.go"

# Build the three templates together:
Task: "Create src/frontend/templates/wishlist.html"
Task: "Edit src/frontend/templates/product.html (Save form + indicator)"
Task: "Edit src/frontend/templates/header.html (saved-products link)"
```

---

## Implementation Strategy

### MVP = User Story 1 (the whole feature)

1. Phase 1: Setup (green baseline).
2. Phase 3: User Story 1 — tests first, then store → wiring → handlers → routes →
   templates → green tests.
3. **STOP and VALIDATE** with quickstart.md.
4. Demo / hand off.

Saving from the product list (AIP-193), removal (AIP-194), and the saved-count badge
(AIP-195) are separate stories and intentionally out of scope here.

---

## Notes

- [P] = different files, no incomplete dependencies.
- Every task names an exact file path; all under `src/frontend/`.
- Verify T002/T003 fail before implementing T004+.
- Commit after each task or logical group.
- Honour the constitution: in-memory only, existing catalogue for product data, no
  infra/CI changes.
