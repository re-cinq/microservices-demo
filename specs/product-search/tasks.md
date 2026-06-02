---

description: "Task list for Product Search (client-side name filter)"
---

# Tasks: Product Search (client-side name filter)

**Input**: Design documents from `specs/product-search/`

**Prerequisites**: [plan.md](plan.md) (required), [spec.md](spec.md) (required for user stories)

**Tests**: NOT requested in the spec — the frontend service has no JS test harness and the plan
specifies manual acceptance against `products.json`. No automated test tasks are included;
verification is captured as explicit manual tasks in Polish.

**Organization**: Tasks are grouped by user story (US1 = MVP filter, US2 = no-results message)
so each story is independently implementable and testable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to (US1, US2)
- All paths are relative to the repository root.

## Path Conventions

Change is confined to the **frontend service**: `src/frontend/`. No backend, proto, or other
service is touched (per plan Structure Decision and constraints C-001…C-010).

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the static asset location for the new client-side script.

- [x] T001 Create the directory `src/frontend/static/js/` (new — no JS assets exist there yet; it will be auto-served by the existing file server at `/static/js/`, see `src/frontend/main.go:159`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The page scaffolding both user stories build on — the search input, its script
include, and base styling. No filtering behaviour yet.

**⚠️ CRITICAL**: User story work (Phase 3+) cannot begin until this phase is complete.

- [x] T002 Add a search input row **above** the `hot-products-row` grid inside the existing `col-12 col-lg-12` container in `src/frontend/templates/home.html`. Reuse `static/icons/Hipster_SearchIcon.svg` as the icon. Give the input a stable id (e.g. `id="product-search-input"`) and an accessible label. Do NOT modify the existing grid markup or its Bootstrap classes (constraint C-008: no layout change).
- [x] T003 Add a `<script src="{{ $.baseUrl }}/static/js/product-search.js" defer></script>` include for the new file in `src/frontend/templates/home.html` (place near the bottom of the `home` block so the grid exists when it runs).
- [x] T004 [P] Add CSS for the search input (and a hidden style for the no-results message used in US2) to `src/frontend/static/styles/styles.css`. Reuse Bootstrap's `d-none` for hiding cards; only add cosmetic styling here — no changes to existing grid rules (C-008).

**Checkpoint**: Search box renders on the home page (inert); script loads with no errors.

---

## Phase 3: User Story 1 - Filter the product grid by name as I type (Priority: P1) 🎯 MVP

**Goal**: Typing in the search box narrows the visible product cards by name in real time, with
no page reload and no backend request; clearing restores the full grid in original order.

**Independent Test**: Load the home page, type a known product-name fragment, confirm only
matching `.hot-product-card` elements remain visible and all others are hidden — with zero
network requests (browser network panel) and no reload.

### Implementation for User Story 1

- [x] T005 [US1] Create `src/frontend/static/js/product-search.js`. On DOM ready, query all `.hot-product-card` elements and, for each, read its product name from the `.hot-product-card-name` child (`textContent`, trimmed). Cache the card→name pairs once.
- [x] T006 [US1] In `src/frontend/static/js/product-search.js`, bind an `input` listener on `#product-search-input` that, on each keystroke, computes `query = value.trim().toLowerCase()` and shows/hides each card by toggling Bootstrap's `d-none` class based on `name.toLowerCase().includes(query)` (case-insensitive substring; literal match via `includes`). Satisfies FR-002, FR-003, FR-004, FR-005, FR-008. (Depends on T002, T003, T005.)
- [x] T007 [US1] In the same listener, handle the empty/whitespace query: when `query === ""`, remove `d-none` from all cards so every product is shown in its original DOM order (FR-006, FR-009). (Depends on T006.)

**Checkpoint**: US1 fully functional — filtering, case-insensitivity, whitespace trimming, and
clear-to-restore all work with zero network traffic. This is a shippable MVP.

---

## Phase 4: User Story 2 - Tell me clearly when nothing matches (Priority: P2)

**Goal**: When a non-empty query matches no product, show an inline "no products match" message
instead of an empty grid; remove it when the query matches again or is cleared.

**Independent Test**: Type a string matching no product name (e.g. "zzzzz") → inline message
appears and all cards hidden; clear or type a matching query → message disappears.

### Implementation for User Story 2

- [x] T008 [US2] Add a single hidden no-results message element (e.g. `<div id="product-search-no-results" class="d-none">No products match your search.</div>`) inside the existing grid container in `src/frontend/templates/home.html`, so results render inline (FR-007, constraint C-010: no separate results page).
- [x] T009 [US2] In `src/frontend/static/js/product-search.js`, after the show/hide pass count the visible cards: if the query is non-empty AND zero cards are visible, remove `d-none` from `#product-search-no-results`; otherwise add `d-none` back to it (FR-007, SC-006). (Depends on T006, T007, T008.)

**Checkpoint**: US1 and US2 both work independently — filtering plus clear no-match feedback.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Verify the feature against the spec and confirm no constraint is violated.

- [x] T010 [P] Accessibility: `#product-search-input` has an `aria-label` plus a visually-hidden `<label for>`; it is a native text input (keyboard-operable). Verified in `src/frontend/templates/home.html` markup.
- [ ] T011 Run the spec's acceptance scenarios in a browser (US1 #1–5, US2 #1–2): filter, clear-to-restore in original order, uppercase + padded query, no-match message, and re-show on clear. Confirm against the names in `src/productcatalogservice/products.json` (SC-002). **NOT run locally** — requires the full stack (frontend + productcatalogservice) to render products; `go` is not installed and no stack is running here. To be validated against the CI deploy of `attendee/andrea-backstrom`. JS logic verified by review + clean `bun` transpile.
- [ ] T012 [P] Open the browser network panel and confirm typing in the search box generates **zero** network requests (FR-008, SC-004), and that a filtered product's link, price, image, and recommendations still work (FR-010, SC-007). **NOT run locally** (same reason as T011). Design guarantees zero requests: the filter is pure DOM (`classList.toggle`), no `fetch`/`XHR` in `product-search.js`; existing card markup is untouched so links/price/image/recommendations are unchanged.
- [x] T013 [P] Constraint audit: `git status` confirms the only feature changes are `src/frontend/templates/home.html`, `src/frontend/static/styles/styles.css` (modified) and `src/frontend/static/js/product-search.js` (new). No Go code, no proto, no `productcatalogservice`, no Helm/manifest/env/pipeline files (constraints C-001…C-007). ✓

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup. BLOCKS all user stories.
- **User Story 1 (Phase 3)**: Depends on Foundational. The MVP.
- **User Story 2 (Phase 4)**: Depends on Foundational; its JS logic (T009) depends on US1's listener (T006/T007). Independently testable once done.
- **Polish (Phase 5)**: Depends on the user stories you intend to ship.

### Within / Across Stories

- T006 depends on T002, T003, T005. T007 depends on T006.
- T009 depends on T006, T007, T008.
- US1 (T005–T007) is a complete, shippable increment on its own.

### Parallel Opportunities

- T004 ([P]) can proceed alongside T002/T003 (different file: `styles.css` vs `home.html`).
- In Polish, T010 / T012 / T013 ([P]) are independent checks.
- US1 and US2 markup (T002/T008 are in the same file `home.html`) are **not** parallel with each other.

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Phase 1: Setup (T001)
2. Phase 2: Foundational (T002–T004)
3. Phase 3: User Story 1 (T005–T007)
4. **STOP and VALIDATE**: filter + clear-to-restore work with zero network requests.
5. Deploy/demo — this is a usable MVP.

### Incremental Delivery

1. Setup + Foundational → search box renders.
2. Add US1 → test → deploy/demo (MVP: live filtering).
3. Add US2 → test → deploy/demo (no-results feedback).
4. Polish → verify all spec scenarios and constraints.

---

## Notes

- [P] = different files, no dependencies.
- No backend, proto, `productcatalogservice`, infra, or pipeline changes (C-001…C-007).
- The new static JS file is served automatically by the existing file server — no Go change.
- Commit after each task or logical group.
- Next: run **`/speckit.implement`** to execute these tasks.
