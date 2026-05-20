# Tasks: Recently Viewed Products Strip

**Input**: Design documents from `specs/001-recently-viewed-strip/`

**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, quickstart.md ✓

**Tests**: Included in Polish phase (plan.md specifies handlers_test.go).

**Organization**: Tasks grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files or no shared dependencies)
- **[Story]**: Which user story this task belongs to ([US1], [US2])

---

## Phase 1: Setup

**Purpose**: Confirm the baseline is clean before any changes.

- [x] T001 Run `cd src/frontend && go test ./...` and confirm all existing tests pass before making any changes ⚠️ Go not found in PATH on this machine — run manually before merging

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared cookie helpers and product-fetching utility used by both user stories. Must be complete before Phase 3 or Phase 4 can begin.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T002 Add `parseRecentlyViewed(r *http.Request) []string` helper function in `src/frontend/handlers.go` — reads the `shop_recently_viewed` cookie, splits on `|`, returns a `[]string` of product IDs (returns nil if cookie is absent or empty)
- [x] T003 [P] Add `writeRecentlyViewedCookie(w http.ResponseWriter, ids []string)` helper function in `src/frontend/handlers.go` — joins IDs with `|`, sets cookie named `shop_recently_viewed` with `MaxAge: 86400` and `SameSite: http.SameSiteLaxMode`
- [x] T004 Add `(fe *frontendServer) recentlyViewedProducts(ctx context.Context, ids []string) ([]*pb.Product, error)` method in `src/frontend/handlers.go` — iterates over `ids`, calls `fe.getProduct` for each, applies `fe.convertCurrency` using `currentCurrency(r)` from the calling handler, silently skips IDs that return `codes.NotFound`, returns the resolved product slice

**Checkpoint**: Cookie helpers and product resolver are available — user story implementation can now begin.

---

## Phase 3: User Story 1 — Browse and Return to Product (Priority: P1) 🎯 MVP

**Goal**: A shopper who has previously viewed products sees a "Recently viewed" strip on the product detail page and can click any entry to return to that product.

**Independent Test**: Manually set `shop_recently_viewed` cookie to two pipe-separated product IDs from `productcatalogservice/products.json`, load any product detail page (`/product/{id}`), and verify the strip renders both tiles with name, image, and price. Click a tile and confirm navigation to the correct product page.

### Implementation for User Story 1

- [x] T005 [P] [US1] Create `src/frontend/templates/recently_viewed.html` — partial template rendering a horizontal strip titled "Recently viewed"; iterates over `.recently_viewed` (a `[]*pb.Product` slice); each tile shows `product.Picture`, `product.Name`, and `renderMoney product.PriceUsd`; wrapped in a guard `{{if .recently_viewed}}` so it is a no-op when the slice is empty; mirror the structure and CSS classes of `src/frontend/templates/recommendations.html`
- [x] T006 [US1] Update `productHandler` in `src/frontend/handlers.go` to call `fe.recentlyViewedProducts` with the IDs from `parseRecentlyViewed(r)` and inject the result as `"recently_viewed"` into the map passed to `injectCommonTemplateData` (depends on T002, T004)
- [x] T007 [US1] Add `{{template "recently_viewed" .}}` include to `src/frontend/templates/product.html` below the recommendations section (depends on T005, T006)

**Checkpoint**: User Story 1 is fully functional — strip renders from a pre-set cookie. End-to-end viewing requires US2 (Phase 4) to auto-populate the cookie.

---

## Phase 4: User Story 2 — Product View is Recorded (Priority: P2)

**Goal**: When a shopper loads a product detail page, that product is automatically recorded in the recently viewed cookie so it appears in the strip on the next page they visit.

**Independent Test**: Clear the `shop_recently_viewed` cookie, load any product detail page (`/product/{id}`), then open DevTools → Application → Cookies and confirm `shop_recently_viewed` contains the viewed product's ID. View a second product and confirm it is prepended. View a product already in the list and confirm it moves to the front without duplication.

### Implementation for User Story 2

- [x] T008 [US2] Update `productHandler` in `src/frontend/handlers.go` to record the current product view: after the product is fetched, call `parseRecentlyViewed(r)` to get the existing list, filter out the current product ID (dedup), prepend the current product ID, truncate to 5, then call `writeRecentlyViewedCookie(w, ids)` (depends on T002, T003)

**Checkpoint**: Both user stories are complete — viewing a product auto-updates the cookie and the strip renders on the next product page visited.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Tests and validation across both user stories.

- [x] T009 [P] Add table-driven unit tests for `parseRecentlyViewed` in `src/frontend/handlers_test.go` — cover: absent cookie (returns nil), empty cookie value (returns nil/empty), single ID, exactly 5 IDs, malformed value with trailing `|`
- [x] T010 [P] Add table-driven unit tests for the dedup/prepend/truncate logic exercised in `productHandler` in `src/frontend/handlers_test.go` — cover: first view (empty list → single entry), new product prepended, duplicate product moves to front, sixth product evicts oldest, same product viewed twice in a row stays at front with no duplicate
- [ ] T011 Run `cd src/frontend && go test ./...` and confirm all tests pass including T009 and T010 ⚠️ Go not found in PATH — run manually
- [ ] T012 Validate all edge cases listed in `specs/001-recently-viewed-strip/quickstart.md` against the running service via `skaffold dev`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — **blocks both user stories**
- **User Story 1 (Phase 3)**: Depends on Phase 2 (T002, T004); T005 can start in parallel with T002/T003/T004
- **User Story 2 (Phase 4)**: Depends on Phase 2 (T002, T003)
- **Polish (Phase 5)**: Depends on Phases 3 and 4; T009 and T010 can run in parallel

### User Story Dependencies

- **User Story 1 (P1)**: Requires T002 and T004 from Foundational
- **User Story 2 (P2)**: Requires T002 and T003 from Foundational; independent of US1
- US1 and US2 can be developed in parallel once Foundational is complete

### Within Each User Story

- US1: T005 (template) and T006 (handler injection) can start in parallel; T007 depends on both
- US2: Single task (T008); no intra-story parallelism needed

### Parallel Opportunities

- T003 (write cookie helper) and T002 (parse cookie helper) can be written in parallel
- T005 (template) can be started in parallel with T002/T003/T004
- T009 and T010 (tests) can be written in parallel

---

## Parallel Example: Foundational Phase

```
# Run in parallel:
Task T002: "Add parseRecentlyViewed helper in src/frontend/handlers.go"
Task T003: "Add writeRecentlyViewedCookie helper in src/frontend/handlers.go"
# T004 depends on T002 being present first — run after T002
```

## Parallel Example: Polish Phase

```
# Run in parallel:
Task T009: "Table-driven tests for parseRecentlyViewed in src/frontend/handlers_test.go"
Task T010: "Table-driven tests for dedup/prepend/truncate in src/frontend/handlers_test.go"
# T011 depends on both T009 and T010 — run after both complete
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002, T003, T004)
3. Complete Phase 3: User Story 1 (T005, T006, T007)
4. **STOP and VALIDATE**: Set cookie manually, confirm strip renders
5. Continue to Phase 4 to make recording automatic

### Incremental Delivery

1. Phase 1 + Phase 2 → helpers available
2. Phase 3 (US1) → strip renders from manually set cookie → demo-able MVP
3. Phase 4 (US2) → strip auto-populates on product view → full feature
4. Phase 5 → tests and edge case validation → ship-ready

### Single Developer Order

T001 → T002 → T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012

---

## Notes

- All changes confined to `src/frontend/` (constraint C-001)
- No new dependencies or imports beyond Go standard library and existing frontend imports
- Cookie format: pipe-separated IDs (`OLJCESPC7Z|66VCHSJNUP`) — no JSON parsing needed
- `recentlyViewedProducts` (T004) silently skips stale product IDs — do not surface errors to the user
- The strip template must guard against an empty slice to avoid rendering an empty section heading
- Currency conversion must use `currentCurrency(r)` so prices match the rest of the page
