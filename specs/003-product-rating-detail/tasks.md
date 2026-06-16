---
description: "Task list for Product Star Rating — Detail Page (AIP-181 + AIP-182)"
---

# Tasks: Product Star Rating — Detail Page

**Input**: Design documents from `/specs/003-product-rating-detail/`

**Jira**: [AIP-182](https://odevo.atlassian.net/browse/AIP-182) / [AIP-181](https://odevo.atlassian.net/browse/AIP-181) / Epic [AIP-95](https://odevo.atlassian.net/browse/AIP-95)

**Organization**: Two user stories. US1 (seed ratings into catalogue data) must complete before US2 (display on product page) can be verified end-to-end. US2 template work can be written in parallel with US1.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no unmet dependencies)
- **[US1]**: Seed star ratings into the product catalogue (AIP-181)
- **[US2]**: Display star rating on the product detail page (AIP-182)

---

## Phase 1: Setup

**Purpose**: No new dependencies or project structure needed — setup is a single orientation step.

- [x] T001 Confirm all 9 product IDs in `src/productcatalogservice/products.json` match the data model table in `specs/003-product-rating-detail/plan.md`

**Checkpoint**: IDs confirmed. Implementation can begin.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The `productRatings` map in the frontend is the shared foundation that both US1 (data) and US2 (display) depend on.

- [x] T002 Add `var productRatings = map[string]float32{...}` package-level map for all 9 products to `src/frontend/handlers.go` after the existing `var` block, using the values from the data model in `plan.md`

**Checkpoint**: `go build ./...` passes in `src/frontend/`. Map is compiled and accessible to handlers.

---

## Phase 3: User Story 1 — Seed ratings into the product catalogue (Priority: P1)

**Goal**: Every product in `products.json` carries a `"rating"` field so the canonical data record reflects the ratings that are displayed.

**Independent Test**: Open `src/productcatalogservice/products.json` and confirm every product object has a `"rating"` key with a float value between 1 and 5. All 9 products must have it.

### Implementation

- [x] T003 [US1] Add `"rating": 4.6` to product `OLJCESPC7Z` (Sunglasses) in `src/productcatalogservice/products.json`
- [x] T004 [US1] Add `"rating": 3.8` to product `66VCHSJNUP` (Tank Top) in `src/productcatalogservice/products.json`
- [x] T005 [US1] Add `"rating": 4.9` to product `1YMWWN1N4O` (Watch) in `src/productcatalogservice/products.json`
- [x] T006 [US1] Add `"rating": 4.1` to product `L9ECAV7KIM` (Loafers) in `src/productcatalogservice/products.json`
- [x] T007 [US1] Add `"rating": 3.5` to product `2ZYFJ3GM2N` (Hairdryer) in `src/productcatalogservice/products.json`
- [x] T008 [US1] Add `"rating": 4.3` to product `0PUK6V6EV0` (Candle Holder) in `src/productcatalogservice/products.json`
- [x] T009 [US1] Add `"rating": 4.7` to product `LS4PSXUNUM` (Salt & Pepper Shakers) in `src/productcatalogservice/products.json`
- [x] T010 [US1] Add `"rating": 4.2` to product `9SIQT8TOJO` (Bamboo Glass Jar) in `src/productcatalogservice/products.json`
- [x] T011 [US1] Add `"rating": 3.9` to product `6E92ZMYYFZ` (Mug) in `src/productcatalogservice/products.json`

**Checkpoint**: All 9 products in `products.json` have a `"rating"` field. Values match the `productRatings` map added in T002.

---

## Phase 4: User Story 2 — Display star rating on the product detail page (Priority: P2)

**Goal**: The product detail page shows the star rating visibly. It matches the stored value and contains no price, review text, or submission UI.

**Independent Test**: Navigate to any product detail page (e.g., Sunglasses). Confirm "★ 4.6 / 5" appears below the product name/price. No form, no review text, no price duplication alongside the rating. Check a second product to confirm the value differs and matches `products.json`.

### Implementation

- [x] T012 [P] [US2] In `productHandler` in `src/frontend/handlers.go`, add `"rating": productRatings[id]` to the `ExecuteTemplate` map (after the `packagingInfo` entry)
- [x] T013 [P] [US2] In `src/frontend/templates/product.html`, add the rating display block after the `<p class="product-price">` line: `{{ if $.rating }}<p class="product-rating">★ {{ printf "%.1f" $.rating }} / 5</p>{{ end }}`

**Checkpoint**: Full acceptance criteria walkthrough:
- AIP-182 AC-1: Rating visible on product page ✓
- AIP-182 AC-2: Displayed value matches `productRatings` map / `products.json` ✓
- AIP-182 AC-3: No price/review/submission UI present ✓

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T014 [P] Verify `go build ./...` and `go vet ./...` pass in `src/frontend/`
- [ ] T015 [P] Confirm all 9 products show their correct rating by visiting each product detail page
- [ ] T016 Smoke-test existing flows unaffected: home page, add-to-cart, cart page, recommendations, recently-viewed strip

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — map must exist before T012 can reference it
- **Phase 3 (US1)**: Can run in parallel with Phase 4 — `products.json` edits are independent of template/handler changes
- **Phase 4 (US2)**: T012 depends on Phase 2 (map must exist); T013 is independent of everything
- **Phase 5 (Polish)**: Depends on Phases 3 and 4 complete

### Critical path

```
T001 → T002 → T003–T011 (US1 data) + T012+T013 (US2 display, parallel) → T014–T016
```

### Parallel opportunities

- T003–T011 are all edits to the same file (`products.json`) — do them in one pass, not truly parallel, but fast
- T012 (`handlers.go`) and T013 (`product.html`) touch different files — genuinely parallel
- T014 and T015 are independent verification tasks — parallel

---

## Parallel Example: Phase 4

```
# These two tasks touch different files — do together:
T012: Edit src/frontend/handlers.go  (add "rating": productRatings[id])
T013: Edit src/frontend/templates/product.html  (add rating display block)
```

---

## Implementation Strategy

### MVP (all tasks are in scope — feature is small)

1. T001: Confirm product IDs
2. T002: Add `productRatings` map to `handlers.go`
3. T003–T011: Add `"rating"` to all products in `products.json`
4. T012 + T013 (parallel): Wire handler + template
5. T014–T016: Build check + smoke tests

### Incremental checkpoints

| After | Verifiable outcome |
|-------|--------------------|
| T002 | `go build ./...` passes |
| T011 | All 9 products in `products.json` have a rating field |
| T013 | Rating visible on product detail page in browser |
| T016 | All ACs confirmed, no regressions |

---

## Notes

- No tests were requested in the spec — no test tasks generated.
- T003–T011 are individual tasks for traceability but should be done in one edit pass on `products.json`.
- T012 and T013 are the only genuinely parallel implementation tasks.
- The `{{ if $.rating }}` guard in the template handles the zero-value case (float32 zero = 0.0 = falsy in Go templates) — if a product ID is somehow missing from the map, no broken UI is shown.
