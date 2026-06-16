---
description: "Task list for Recently Viewed Strip (AIP-175)"
---

# Tasks: Recently Viewed Strip

**Input**: Design documents from `/specs/002-recently-viewed/`

**Jira**: [AIP-175](https://odevo.atlassian.net/browse/AIP-175) / Epic [AIP-97](https://odevo.atlassian.net/browse/AIP-97)

**Organization**: Single user story (AIP-175). Tasks proceed sequentially — each step is a small, self-contained change on the critical path.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no unmet dependencies)
- **[US1]**: Belongs to the single user story (AIP-175 — Display recently viewed strip)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the constant and helper functions that all subsequent tasks depend on.

- [x] T001 Add `cookieRecentlyViewed = "recently_viewed"` constant to `src/frontend/main.go` alongside `cookieCurrency`
- [x] T002 Add helper `recentlyViewedIDs(r *http.Request) []string` to `src/frontend/handlers.go` — reads cookie, splits on comma, returns empty slice if absent
- [x] T003 Add helper `addRecentlyViewed(existing []string, id string) []string` to `src/frontend/handlers.go` — prepends id, deduplicates, caps at 5
- [x] T004 Add helper `filterByLiveCatalogue(ids []string, products []*pb.Product) []*pb.Product` to `src/frontend/handlers.go` — returns only products whose ID is in ids, preserving order
- [x] T005 Add helper `excludeProduct(products []*pb.Product, id string) []*pb.Product` to `src/frontend/handlers.go` — removes the product with the given id from the slice

**Checkpoint**: All helpers compile. Run `go build ./...` from `src/frontend/` to verify.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Wire the cookie read/write into `productHandler` before the template can use it.

- [x] T006 [US1] In `productHandler` in `src/frontend/handlers.go`, after the existing product fetch, read the existing cookie and write the updated cookie using `recentlyViewedIDs`, `addRecentlyViewed`, and `http.SetCookie`
- [x] T007 [US1] In `productHandler`, call `fe.getProducts` to get the live catalogue, then build `recentProducts` using `filterByLiveCatalogue` and `excludeProduct`
- [x] T008 [US1] Pass `"recently_viewed": recentProducts` into the `ExecuteTemplate` map in `productHandler`

**Checkpoint**: App compiles and runs. Visiting a product page sets the `recently_viewed` cookie (verify in browser DevTools → Application → Cookies). No visible UI change yet.

---

## Phase 3: User Story 1 — Display recently viewed strip (Priority: P1) 🎯 MVP

**Goal**: Render the "Recently Viewed" strip on the product page, conditionally shown only when validated recently-viewed products exist.

**Independent Test**: Visit 3 different product pages in sequence. On the 3rd page, a "Recently Viewed" strip appears below recommendations showing the 2 previously visited products (not the current one). Navigate to a 4th product — the strip shows 3 items. Clear cookies and visit one product — no strip appears.

### Implementation

- [x] T009 [P] [US1] Create `src/frontend/templates/recently_viewed.html` — define the `recently_viewed` template partial following the structure of `recommendations.html`: section wrapper, container, heading "Recently Viewed", row of `col-md-3` cards each linking to `/product/{{ .Id }}` with image and name; wrap entire content in `{{ if .recently_viewed }}` guard
- [x] T010 [P] [US1] In `src/frontend/templates/product.html`, add `{{ template "recently_viewed" $ }}` inside a `<div>` after the existing recommendations block

**Checkpoint**: Full acceptance criteria walkthrough:
- AC-1: Visit 2+ products → strip appears on subsequent product pages ✓
- AC-2: Clear cookies → visit one product → no strip shown ✓
- AC-3: All strip items are real products from catalogue ✓
- AC-4/AC-5: (manual — requires removing a product from `productcatalogservice/products.json`, restarting, verifying it disappears from strip) ✓

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Verify correctness, clean up, confirm no regressions.

- [ ] T011 [P] Verify `go build ./...` and `go vet ./...` pass cleanly in `src/frontend/`
- [ ] T012 [P] Smoke-test existing flows unaffected: home page loads, add-to-cart works, cart page renders, recommendations still appear on product page
- [ ] T013 Verify strip does not appear on first product visited (cookie not yet set at time of render — strip correctly absent)
- [ ] T014 Verify strip cap: visit 7 different products — strip shows at most 5, oldest evicted

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately. T002–T005 can be written in one sitting as a helper block.
- **Phase 2 (Foundational)**: Depends on Phase 1 completion (needs T001–T005).
- **Phase 3 (User Story 1)**: T009 and T010 depend on Phase 2 completion (T008 must pass `recently_viewed` to template). T009 and T010 touch different files — fully parallel.
- **Phase 4 (Polish)**: Depends on Phase 3 completion.

### Within Phase 3

- T009 (`recently_viewed.html`) and T010 (`product.html`) are independent — different files, no conflict. Mark both [P].

### Critical path

```
T001 → T002–T005 → T006 → T007 → T008 → T009+T010 (parallel) → T011–T014
```

---

## Parallel Example: Phase 3

```
# These two tasks touch different files — launch together:
T009: Create src/frontend/templates/recently_viewed.html
T010: Edit src/frontend/templates/product.html
```

---

## Implementation Strategy

### MVP (this feature is a single story — all tasks are MVP)

1. Complete Phase 1: helpers in `handlers.go` + constant in `main.go`
2. Complete Phase 2: wire cookie into `productHandler`
3. Complete Phase 3: template partial + product page wiring
4. **Validate**: walk all 5 ACs manually in browser
5. Polish: Phase 4 smoke tests

### Incremental checkpoints

| After | Verifiable outcome |
|-------|--------------------|
| T005 | `go build ./...` passes |
| T008 | Cookie set on product page visit (DevTools) |
| T010 | Strip visible after ≥2 product views |
| T014 | All ACs confirmed, no regressions |

---

## Notes

- No tests were requested in the spec — no test tasks generated.
- T009 and T010 are the only parallelisable implementation tasks; all others are sequential.
- Cookie is `HttpOnly` — not accessible from JS, consistent with existing cookie usage in the codebase.
- `recentlyViewedIDs` is called twice in `productHandler` (once before write to get existing IDs, once after to get the updated list for template rendering) — this is intentional: the updated cookie is written to the response but not re-read from the request in the same handler invocation.
