# Tasks: Save for Later / Wishlist

**Input**: Design documents from `specs/001-save-wishlist/`

**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓

**Tests**: Unit tests included in Polish phase (not TDD — spec does not require test-first).

**Organization**: Tasks grouped by user story (US1 → US2 → US3) to enable independent implementation and testing.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story this task belongs to (US1, US2, US3)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the in-memory wishlist store to the existing server struct — the single prerequisite all story tasks share.

- [ ] T001 Add `wishlists sync.Map` field to `frontendServer` struct in `src/frontend/main.go`

**Checkpoint**: `frontendServer` has a `wishlists` field of type `sync.Map`. Compiles cleanly with `go build ./...`.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define the `WishlistItem` view-model used by both US1 and US2 handlers and templates.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T002 Define `WishlistItem` struct (fields: `ProductID`, `Name`, `Picture`, `Price pb.Money`) in `src/frontend/handlers.go`

**Checkpoint**: `WishlistItem` struct is defined and compiles. No handler logic yet.

---

## Phase 3: User Story 1 — Save a product from the product detail page (Priority: P1) 🎯 MVP

**Goal**: A shopper can click "Save for later" on any product detail page and see an inline confirmation that the product was saved to their session wishlist.

**Independent Test**: Visit `/product/<any-id>`, click "Save for later", confirm the inline confirmation banner appears on the same page. No wishlist view needed.

### Implementation for User Story 1

- [ ] T003 [US1] Implement `saveWishlistHandler` in `src/frontend/handlers.go`: read `product_id` from POST body, load session's product-ID slice from `fe.wishlists`, append if not present, store back, redirect to `/product/{product_id}?saved=1`. Return 400 if `product_id` is empty.
- [ ] T004 [P] [US1] Register `POST /wishlist/save` route pointing to `saveWishlistHandler` in `src/frontend/main.go` (parallel with T005 — different file)
- [ ] T005 [P] [US1] Add "Save for later" `<form method="POST" action="...">` button and conditional `{{ if $.saved }}` inline confirmation banner to `src/frontend/templates/product.html`; update `productHandler` in `src/frontend/handlers.go` to read `r.URL.Query().Get("saved")` and pass `"saved": true/false` into the template data map

**Checkpoint**: `POST /wishlist/save` with a valid `product_id` redirects to the product page showing the confirmation banner. Saving the same product twice does not duplicate it.

---

## Phase 4: User Story 2 — View saved products in the wishlist (Priority: P2)

**Goal**: A shopper can navigate to `/wishlist` and see all products saved in the current session, each showing name, image, and price. An empty state is shown when nothing is saved.

**Independent Test**: Save one or more products (US1 complete), navigate to `/wishlist`, confirm each saved product appears with name, image, and price. Navigate to `/wishlist` with nothing saved, confirm empty-state message appears.

### Implementation for User Story 2

- [ ] T006 [US2] Implement `viewWishlistHandler` in `src/frontend/handlers.go`: load saved product IDs from `fe.wishlists` for the session, fetch each product via `fe.getProduct`, convert price via `fe.convertCurrency`, build `[]WishlistItem`, render `wishlist.html`
- [ ] T007 [P] [US2] Register `GET /wishlist` route pointing to `viewWishlistHandler` in `src/frontend/main.go`
- [ ] T008 [P] [US2] Create `src/frontend/templates/wishlist.html`: extends header/footer, renders product cards (name, image, price) for each item in `$.items`; shows empty-state message when `$.items` is empty
- [ ] T009 [P] [US2] Add wishlist icon nav link to `src/frontend/templates/header.html` next to the existing cart link, following the same `<a href="{{ $.baseUrl }}/wishlist" class="cart-link">` pattern

**Checkpoint**: `/wishlist` renders saved products with name, image, and price. Navigating to `/wishlist` with no saved items shows the empty-state message. Header shows the wishlist icon on every page.

---

## Phase 5: User Story 3 — Wishlist is cleared when the session ends (Priority: P3)

**Goal**: A shopper's wishlist is empty in a new browser session. No code changes required — this is guaranteed by the `sync.Map` keyed on session ID combined with the existing cookie-based session mechanism.

**Independent Test**: Save products in one session, open an incognito window (new session, no cookie), navigate to `/wishlist` — list must be empty.

### Verification for User Story 3

- [ ] T010 [US3] Verify session-scoped clearing manually per `specs/001-save-wishlist/quickstart.md` steps (incognito window test) and confirm the behaviour matches AC in spec.md; update quickstart.md with explicit session-clearing verification step if not already present

**Checkpoint**: New session always starts with an empty wishlist. No data persists across sessions.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Unit tests, compilation verification, code quality.

- [ ] T011 [P] Create `src/frontend/wishlist_test.go` with unit tests for: `saveWishlistHandler` adds product, prevents duplicate, rejects empty `product_id`; `viewWishlistHandler` returns empty slice when no products saved
- [ ] T012 [P] Run `go build ./...` and `go vet ./...` from `src/frontend/` and confirm both pass with no errors or warnings

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — blocks all user story work
- **Phase 3 (US1)**: Depends on Phase 2 — T004 and T005 can run in parallel after T003
- **Phase 4 (US2)**: Depends on Phase 2 — T007, T008, T009 can run in parallel after T006; can also start in parallel with Phase 3 since they touch different functions
- **Phase 5 (US3)**: Depends on Phase 3 completing (needs save to work for the test)
- **Phase 6 (Polish)**: Depends on all story phases completing

### User Story Dependencies

- **US1 (P1)**: Depends on Foundational only — no dependency on US2 or US3
- **US2 (P2)**: Depends on Foundational only — no dependency on US1 (independently testable, but realistically needs US1 to populate the list)
- **US3 (P3)**: No new code — depends on US1 working so the session-clearing test has something to verify

### Within Each User Story

- Same-file changes are sequential (both handlers.go edits in T003 and T005 must not conflict)
- Different-file changes are parallel ([P] marked)
- Templates are always parallel with Go handler changes

### Parallel Opportunities

```bash
# After T003 completes (US1):
T004  # Register route in main.go
T005  # Add button + banner to product.html
# Both touch different files — run together

# After T006 completes (US2):
T007  # Register route in main.go
T008  # Create wishlist.html
T009  # Add nav icon to header.html
# All touch different files — run together

# Polish phase:
T011  # Unit tests
T012  # Build + vet
# Independent — run together
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002)
3. Complete Phase 3: US1 (T003 → T004 + T005 in parallel)
4. **STOP and VALIDATE**: Shoppers can save products and see confirmation
5. Ship US1 independently if needed

### Incremental Delivery

1. Setup + Foundational → T001, T002
2. US1 complete → Save + confirmation working (MVP)
3. US2 complete → Wishlist view working
4. US3 verified → Session scoping confirmed (no code)
5. Polish → Tests + build check

---

## Notes

- `sync.Map` is already in the Go standard library — no new dependencies
- `productHandler` already fetches product details; reuse `fe.getProduct` and `fe.convertCurrency` in `viewWishlistHandler`
- The `sessionID(r)` helper already exists in `handlers.go` — use it directly in both new handlers
- Cookie `shop_session-id` has a 48-hour max-age — orphaned `sync.Map` entries for expired sessions accumulate but are never accessed; acceptable at demo scale
- All template changes follow the existing `{{ define "name" }} ... {{ end }}` pattern
