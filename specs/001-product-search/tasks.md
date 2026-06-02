# Tasks: Browser-Side Product Search

**Input**: Design documents from `specs/001-product-search/`

**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ui-contract.md ✓

**Tests**: No test tasks generated — spec does not request TDD; verification is manual browser testing per acceptance scenarios.

**Organization**: Tasks are grouped by user story (US1 → US2 → US3) in priority order. Each story phase is independently completable and testable.

---

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no shared dependency on an incomplete task)
- **[Story]**: User story this task belongs to (US1/US2/US3)

---

## Phase 1: Setup

**Purpose**: No new project, framework, or dependency is needed. The existing static file server in `main.go` already serves any file placed under `src/frontend/static/` recursively. No routing or configuration changes are required.

*No tasks — proceed directly to Phase 2.*

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: Expose product names in the DOM so the JavaScript filter can read them without touching inner HTML. This single change unblocks all three user stories.

**⚠️ CRITICAL**: No user-story work can begin until T001 is complete.

- [x] T001 Add `data-name="{{ .Item.Name }}"` attribute to the `.hot-product-card` div in `src/frontend/templates/home.html` (the outer `<div class="col-md-4 hot-product-card">` inside the `{{ range $.products }}` loop)

**Checkpoint**: Render the home page and confirm each product card has a `data-name` attribute with the correct product name visible in the browser DevTools element inspector.

---

## Phase 3: User Story 1 — Filter Products by Name (Priority: P1) 🎯 MVP

**Goal**: A shopper can type in a search box and the product grid immediately narrows to products whose names match. Clearing the box restores all products.

**Independent Test**: Load the home page, type `shirt` — only matching products visible. Clear input — all products reappear. Type `xyzzzz` — zero products visible and no-results message shown.

### Implementation for User Story 1

- [x] T002 [US1] Add `<input id="product-search" type="search" class="form-control" placeholder="Search products…" aria-label="Search products" autocomplete="off">` and `<div id="no-results-message" style="display:none;" class="col-12 text-center py-4">No products match your search.</div>` above the product grid row in `src/frontend/templates/home.html`
- [x] T003 [US1] Create `src/frontend/static/scripts/product-search.js` with: `input` event listener on `#product-search`; for each `.hot-product-card[data-name]`, compare `element.dataset.name.toLowerCase()` against `query.trim().toLowerCase()` and toggle `element.style.display`; after each pass, show `#no-results-message` when `visibleCount === 0` and query is non-empty, hide it otherwise; treat whitespace-only query as empty (restore all cards)
- [x] T004 [US1] Add `<script src="{{ $.baseUrl }}/static/scripts/product-search.js"></script>` at the end of the `<body>` section in `src/frontend/templates/home.html`, after the existing Bootstrap JS CDN tags

**Checkpoint**: US1 fully functional — filtering, no-results message, and clear-to-restore all work. This is the complete MVP.

---

## Phase 4: User Story 2 — Instant, No-Reload Filtering (Priority: P2)

**Goal**: Filtering happens on every keystroke with no button press, no page navigation, and no network request. The DevTools Network tab stays silent while typing.

**Independent Test**: Open DevTools Network tab, type in the search box — no new requests appear. Check the URL bar — it does not change. Confirm products update within 300 ms.

### Implementation for User Story 2

- [x] T005 [US2] Verify the `#product-search` input added in T002 is NOT wrapped in a `<form>` element in `src/frontend/templates/home.html`; if it is inside a form, move it out or add `onsubmit="return false;"` to the form to prevent Enter-key page reload

**Checkpoint**: US2 satisfied — typing produces no navigation and no network requests. US2 requires no additional JS beyond what T003 already implements (`input` event, not `submit`).

---

## Phase 5: User Story 3 — Keyboard-Accessible Search (Priority: P3)

**Goal**: A keyboard-only user can Tab to the search box, type a query, and clear it without a mouse. Escape clears the input and restores all products.

**Independent Test**: Without touching the mouse: Tab to the search box, type a query, confirm filtered results. Press Escape, confirm all products reappear.

### Implementation for User Story 3

- [x] T006 [P] [US3] Confirm `aria-label="Search products"` is present on `#product-search` in `src/frontend/templates/home.html` (should already be there from T002 — verify and add if missing)
- [x] T007 [US3] Add a `keydown` listener in `src/frontend/static/scripts/product-search.js` that, when `event.key === 'Escape'`, sets `#product-search` value to `''` and triggers the filter logic to restore all products

**Checkpoint**: US3 satisfied — Tab focuses the input, typing filters, Escape restores. No mouse needed.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Visual spacing and final end-to-end validation.

- [x] T008 [P] Add `.search-bar-row { margin-bottom: 1.5rem; }` to `src/frontend/static/styles/styles.css` to give the search box breathing room above the product grid
- [ ] T009 Manually verify all 8 acceptance scenarios from `specs/001-product-search/spec.md` against the running frontend: case-insensitive match (US1-S1), clear restores (US1-S2), no-results message (US1-S3), no network request (US2-S1), no URL change (US2-S1), Tab focus works (US3-S1), Escape clears (US3-S2), whitespace-only query restores all products (Edge Case)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No tasks — start immediately
- **Phase 2 (Foundational — T001)**: No dependencies — start immediately
- **Phase 3 (US1 — T002, T003, T004)**: Depends on T001 ⚠️
- **Phase 4 (US2 — T005)**: Depends on T002, T003, T004 (US1 must be complete)
- **Phase 5 (US3 — T006, T007)**: Depends on T002, T003 (US1 core must be complete); T006 [P] can run alongside T007
- **Phase 6 (Polish — T008, T009)**: T008 [P] can run any time after T001; T009 requires all prior tasks complete

### User Story Dependencies

- **US1 (P1)**: Can start after T001 (Foundational) — no dependency on US2 or US3
- **US2 (P2)**: Depends on US1 implementation (T002, T003, T004) — verifies US1's event approach is correct
- **US3 (P3)**: Depends on T002 (input in DOM) and T003 (JS file exists) — adds Escape key on top of US1

### Within Each Phase

- T002 → T003 → T004 (sequential — same file then new file then same file again)
- T005 is a verification/guard; if the input was correctly placed outside a form in T002, T005 is trivial
- T006 [P] and T007 can be done in parallel (different concerns in different files)

---

## Parallel Opportunities

```text
# After T001, these can all start in parallel:
T002  (home.html — HTML structure)
T008  (styles.css — spacing only, no dependency on JS)

# After T002 + T003:
T006  (verify aria-label — read-only check)
T007  (add Escape handler in product-search.js)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001 — add `data-name` to cards
2. Complete T002 — add search input and no-results message to template
3. Complete T003 — create filter JS file
4. Complete T004 — load script in template
5. **STOP and VALIDATE**: Test all three US1 acceptance scenarios
6. Deploy/demo — this is the complete, shippable MVP

### Incremental Delivery

1. T001 → T002 → T003 → T004: Core filter working (US1) ✅ shippable
2. T005: Harden no-reload guarantee (US2) ✅ shippable
3. T006 + T007: Keyboard accessibility (US3) ✅ shippable
4. T008 + T009: Polish + full validation ✅ done

---

## Notes

- All changes live in `src/frontend/` — no other service is touched
- The existing `SearchProducts` gRPC handler in `productcatalogservice` is intentionally left unused
- [P] tasks operate on different files or are read-only checks — safe to parallelise
- Each story phase checkpoint can be tested independently before proceeding
- Total: 9 implementation tasks across 3 user stories + 1 foundational task
