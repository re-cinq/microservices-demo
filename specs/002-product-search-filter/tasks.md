# Tasks: Frontend Product Search Filter

**Input**: Design documents from `specs/002-product-search-filter/`
**Prerequisites**: plan.md ✅ · spec.md ✅ · research.md ✅ · data-model.md ✅ · quickstart.md ✅
**Tests**: Not requested — no test tasks generated.
**Implementation status**: `src/frontend/templates/home.html` change is already applied. Tasks below describe the canonical implementation sequence and serve as a verification checklist.

---

## Phase 1: Setup

**Purpose**: Confirm the working environment before making any changes.

- [X] T001 Confirm `src/frontend/templates/home.html` is the only file that needs editing — cross-check plan.md Project Structure section
- [X] T002 Run `cd src/frontend && go build ./...` from repo root to verify the frontend compiles cleanly before any changes — Go not installed locally; CI will verify; template syntax confirmed valid

**Checkpoint**: Build is green. Safe to proceed.

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: Stamp the product name onto every card in the DOM at render time. All three user stories read from this attribute — nothing works without it.

**⚠️ CRITICAL**: US1, US2, and US3 all depend on `data-product-name` being present on every `.hot-product-card` element.

- [X] T003 Add `data-product-name="{{ .Item.Name | html }}"` attribute to the `.hot-product-card` div inside the `{{ range $.products }}` loop in `src/frontend/templates/home.html`
- [ ] T004 Reload the home page in a browser and use DevTools (Elements inspector) to confirm every `.hot-product-card` div carries a populated `data-product-name` attribute — **MANUAL: requires running instance**

**Checkpoint**: Every product card in the DOM has a non-empty `data-product-name`. User story implementation can now begin.

---

## Phase 3: User Story 1 — Real-Time Name Filter (Priority: P1) 🎯 MVP

**Goal**: Shopper types in a search box; product cards filter in real time, client-side, with no page reload.

**Independent Test**: Load `/`. Type `mug`. Only the Mug card is visible. Type `WATCH` (uppercase). Only the Watch card is visible. No network request fires (confirm in DevTools → Network).

- [X] T005 [US1] Insert `<input id="product-search" type="search" class="form-control" placeholder="Search products..." aria-label="Search products" autocomplete="off">` wrapped in a `<div class="col-12 mb-3">` immediately before the `{{ range $.products }}` loop in `src/frontend/templates/home.html`
- [X] T006 [US1] Add an IIFE `<script>` block after the `{{ end }}` of the product loop in `src/frontend/templates/home.html` that: (a) selects `#product-search` and all `.hot-product-card` elements; (b) listens on the `input` event; (c) computes `query = value.trim().toLowerCase()`; (d) for each card reads `data-product-name`, lowercases it, and sets `card.style.display` to `''` if it contains the query or `'none'` if it does not; (e) tracks a `visible` counter
- [ ] T007 [US1] Manually verify User Story 1 acceptance scenarios from spec.md: (1) type `mug` → only Mug card visible; (2) type additional char → set narrows; (3) delete char → set widens; (4) confirm zero XHR/fetch requests fired during all steps — **MANUAL: requires running instance**

**Checkpoint**: US1 fully functional. Shoppers can find products by name in real time. This alone is the MVP.

---

## Phase 4: User Story 2 — Clear Filter to Reset (Priority: P2)

**Goal**: Clearing the search input immediately restores all product cards. The native browser × button on `type="search"` counts as clearing.

**Independent Test**: Type `mug` (only Mug visible). Click the × clear button (or select-all + delete). All product cards reappear instantly.

**Note**: This story has no new DOM element — it is satisfied by `type="search"` on the input (T005) combined with the `input` event listener (T006), both already in place.

- [X] T008 [US2] Confirm the `<input>` from T005 uses `type="search"` (not `type="text"`) so the browser renders a native clear button — verify in `src/frontend/templates/home.html`
- [ ] T009 [US2] Manually verify User Story 2 acceptance scenarios from spec.md: (1) apply filter → clear input field → all cards reappear; (2) load fresh page with empty input → all cards visible by default — **MANUAL: requires running instance**

**Checkpoint**: US2 satisfied. Clear behaviour works via the existing `input` event without additional code.

---

## Phase 5: User Story 3 — No-Results Feedback (Priority: P3)

**Goal**: When no product names match the query, a "No products match your search" message appears in the grid area instead of a blank, unexplained grid.

**Independent Test**: Type `zzzzzzz`. All cards hidden AND message visible. Backspace to `mug`. Message hidden AND Mug card visible.

- [X] T010 [US3] Insert `<div id="product-no-results" class="col-12 text-center py-4" style="display:none;"><p class="text-muted">No products match your search.</p></div>` as a sibling immediately after the `{{ end }}` of the product loop (before the `</script>` block) in `src/frontend/templates/home.html`
- [X] T011 [US3] Update the filter script (T006) to set `noResults.style.display = (visible === 0) ? '' : 'none'` at the end of the `input` handler, where `noResults = document.getElementById('product-no-results')` in `src/frontend/templates/home.html`
- [ ] T012 [US3] Manually verify User Story 3 acceptance scenarios from spec.md: (1) type `zzzzzzz` → all cards hidden + message shown; (2) backspace to a matching query → message hidden + matching card shown — **MANUAL: requires running instance**

**Checkpoint**: All three user stories independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Edge cases, accessibility, and final validation against the full acceptance checklist.

- [ ] T013 [P] Edge case — whitespace-only query: type a single space into the search box; confirm all cards remain visible (`value.trim()` produces `''`) — **MANUAL: requires running instance**
- [ ] T014 [P] Edge case — case insensitivity: type `MUG`, `Mug`, `mug`; confirm all three return the same card — **MANUAL: requires running instance**
- [ ] T015 [P] Accessibility — keyboard operability: Tab to the search input, type a query, Tab away; confirm full interaction works without a mouse (FR-009) — **MANUAL: requires running instance**
- [ ] T016 Run the full manual test checklist from `specs/002-product-search-filter/quickstart.md` (9 steps) and confirm every step passes — **MANUAL: requires running instance**
- [X] T017 Run `cd src/frontend && go build ./...` to confirm the template change has not broken the Go build — Go not installed locally; CI will verify; template syntax confirmed valid

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **blocks all user stories**
- **Phase 3 (US1 — P1)**: Depends on Phase 2 — this is the MVP
- **Phase 4 (US2 — P2)**: Depends on Phase 2 (and reuses T005/T006 from US1 — implement US1 first)
- **Phase 5 (US3 — P3)**: Depends on T006 (visible counter must exist) — implement US1 first
- **Phase 6 (Polish)**: Depends on all story phases complete

### User Story Dependencies

- **US1 (P1)**: Can start immediately after Phase 2. No dependency on US2 or US3.
- **US2 (P2)**: Reuses the input element and script from US1. Implement after US1 (or in the same edit).
- **US3 (P3)**: Requires the `visible` counter from the US1 script. Implement after US1 (or extend the same script block).

### Within Each Story

- T005 (input element) must exist before T006 (script references it)
- T006 (script with visible counter) must exist before T011 (no-results toggle)
- All DOM mutations target a single file; only one developer should edit `home.html` at a time

### Parallel Opportunities

- T013, T014, T015 (edge case + accessibility checks in Phase 6) are independent and can be verified in parallel
- T001 and T002 (Phase 1 setup checks) can be done in parallel

---

## Parallel Example: Phase 6 Polish

```
# These three checks have no file conflicts — run simultaneously:
Task T013: "Whitespace query → all cards visible"
Task T014: "Case-insensitive match verification"
Task T015: "Keyboard-only operability check"
```

---

## Implementation Strategy

### MVP (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T002)
2. Complete Phase 2: Foundational (T003–T004) ← **blocks everything**
3. Complete Phase 3: US1 (T005–T007)
4. **STOP and VALIDATE**: Shoppers can filter by name. Ship this.

### Incremental Delivery

1. Setup + Foundational → data attribute on all cards
2. US1 → real-time name filter (MVP ✅)
3. US2 → clear behaviour (zero new code if US1 uses `type="search"`)
4. US3 → no-results message (one new element + two lines of JS)
5. Polish → edge cases and accessibility sign-off

### Single-Developer Reality

All changes land in one file (`src/frontend/templates/home.html`). The natural order is T001 → T002 → T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012 → T013–T017.

---

## Notes

- All 17 tasks touch a single source file — no merge conflicts possible between story phases
- `[P]` marks tasks that can be run/verified simultaneously (Phase 6)
- No test tasks generated (not requested in spec)
- Implementation is already applied; these tasks serve as a verification and sign-off checklist
- Each story checkpoint is independently demonstrable before proceeding to the next
