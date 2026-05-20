---
description: "Task list for AIP-167: Remove a product from the wishlist"
---

# Tasks: Remove a product from the wishlist

**Input**: Design documents from `/specs/007-remove-from-wishlist/`

**Prerequisites**: plan.md ✓, spec.md ✓

**Organization**: Single user story (P1). No setup or foundational phase needed — all dependencies (routing infrastructure, `wishlists sync.Map`, `wishlist.html`) are already present from AIP-159.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to

---

## Phase 1: User Story 1 - Remove a saved product from the wishlist (Priority: P1) 🎯 MVP

**Goal**: A shopper can remove individual products from their wishlist; the page reflects the change immediately, and the empty-state is shown when the last item is removed.

**Independent Test**: Save at least one product, navigate to `/wishlist`, click "Remove", verify the product no longer appears. Remove the last saved product, verify the empty-state message renders.

### Implementation for User Story 1

- [x] T001 [US1] Add `removeWishlistHandler` to `src/frontend/handlers.go` after `saveWishlistHandler`: read `product_id` from form value (return HTTP 400 if empty), load `[]string` from `fe.wishlists` for the session, filter out the matching ID, store the result back, redirect to `/wishlist` with HTTP 303
- [x] T002 [P] [US1] Register route `r.HandleFunc(baseUrl+"/wishlist/remove", svc.removeWishlistHandler).Methods(http.MethodPost)` in `src/frontend/main.go` after the existing `/wishlist/save` route (line 169)
- [x] T003 [P] [US1] Add Remove button form to each product card in `src/frontend/templates/wishlist.html` inside the `{{ range $.items }}` loop after the existing "Add to Cart" form: `<form method="POST" action="{{ $.baseUrl }}/wishlist/remove"><input type="hidden" name="product_id" value="{{ .ProductID }}" /><button type="submit" class="cymbal-button-secondary btn-sm">Remove</button></form>`

**Checkpoint**: All three tasks complete → User Story 1 is fully functional. AC #1 (remove product), AC #2 (empty state via existing `{{ else }}` branch), and AC #3 (no `?saved=1` banner on normal navigation) are all satisfied.

---

## Phase 2: Polish & Cross-Cutting Concerns

**Purpose**: Verification and tidying after core implementation

- [x] T004 Verify `go build ./...` succeeds in `src/frontend/` (no compilation errors)
- [x] T005 [P] Verify existing `wishlist_test.go` tests still pass with `go test ./...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: Can start immediately — all prerequisites from AIP-159 are present
- **Phase 2**: Depends on Phase 1 completion

### Within Phase 1

- T001 (handler) has no file dependency on T002 or T003 — all three can be written in parallel
- T002 and T003 each touch different files (`main.go` and `wishlist.html`) — parallelizable
- T001 must exist before T002 can be tested end-to-end, but both can be written independently

### Parallel Opportunities

```
# All three implementation tasks can be written in parallel:
T001: Add removeWishlistHandler in src/frontend/handlers.go
T002: Register route in src/frontend/main.go
T003: Add Remove button form in src/frontend/templates/wishlist.html
```

---

## Implementation Strategy

### MVP (this story is the entire feature)

1. Complete T001, T002, T003 (can be done in a single pass — all small changes)
2. Build and verify: `go build ./...`
3. Run existing tests: `go test ./...`
4. Manual smoke test: save a product, visit `/wishlist`, click Remove, verify it's gone

---

## Notes

- T001 pattern: mirror `saveWishlistHandler` — same `sync.Map` load/store pattern, same `sessionID(r)` call, same logrus logger retrieval from context if logging is needed
- The empty-state (AC #2) is already handled by the `{{ else }}` branch in `wishlist.html` — no template change needed for that case
- AC #3 requires no code: the `?saved=1` query param is only set by `saveWishlistHandler`'s redirect; normal page navigation never carries it
