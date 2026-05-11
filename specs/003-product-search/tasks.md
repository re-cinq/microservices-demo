---

description: "Task list for Product Search feature"
---

# Tasks: Product Search

**Input**: Design documents from `/specs/003-product-search/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/ui-contract.md, quickstart.md

**Tests**: Not requested. This feature is a small, single-file UI behaviour with no backend code; verification is by the browser checklist in `quickstart.md` (the spec did not request automated tests, and adding a JS test harness was rejected in `research.md` R9). No test tasks are emitted.

**Organization**: Tasks are grouped by user story. Setup + Foundational tasks must complete before any user-story phase begins, because all three stories share the same DOM scaffold and JS module.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: `[US1]`, `[US2]`, `[US3]` for user-story-phase tasks; omitted for Setup/Foundational/Polish
- Include exact file paths in descriptions

## Path Conventions

This feature touches only the `frontend` service:

- `src/frontend/templates/home.html`
- `src/frontend/static/js/product-search.js` (new)
- `src/frontend/static/styles/styles.css`

No other paths are edited.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare directories for the new static asset.

- [X] T001 Create directory `src/frontend/static/js/` (the frontend's existing `http.FileServer` at `/static/` will serve files placed here with no Go code change).

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Lay down the DOM contract from `contracts/ui-contract.md` that every user story relies on (search input, live region, grid identifier, card `data-name`, empty-state element, JS scaffold, CSS, script tag).

**⚠️ CRITICAL**: No user story can be implemented until this phase is complete — they all share the same DOM and the same JS module.

- [X] T002 In `src/frontend/templates/home.html`, add `id="hot-products-row"` to the existing `<div class="row hot-products-row px-xl-6">` element (the grid container).
- [X] T003 In `src/frontend/templates/home.html`, inside the `{{ range $.products }}` loop, add `data-name="{{ lower .Item.Name }}"` to each `<div class="col-md-4 hot-product-card">` card.
- [X] T004 In `src/frontend/templates/home.html`, insert immediately **before** `.hot-products-row` the search-input block from `contracts/ui-contract.md` §1: `<div class="product-search">` containing a visually-hidden `<label for="product-search-input">Search products</label>`, an `<input id="product-search-input" type="search" autocomplete="off" placeholder="Search products" aria-controls="hot-products-row">`, and a `<div id="product-search-status" class="visually-hidden" role="status" aria-live="polite" aria-atomic="true"></div>`.
- [X] T005 In `src/frontend/templates/home.html`, insert immediately **after** `.hot-products-row` the empty-state element `<div class="product-search-empty" hidden>No products match your search.</div>`.
- [X] T006 [P] Create `src/frontend/static/js/product-search.js` with the IIFE scaffold: `document.addEventListener('DOMContentLoaded', () => { … })`; inside, query `#product-search-input`, `#product-search-status`, all `.hot-product-card` elements, and `.product-search-empty`; bail out early if the input is not present (so the script is a no-op on other pages).
- [X] T007 In `src/frontend/templates/home.html`, add `<script defer src="{{ $.baseUrl }}/static/js/product-search.js"></script>` immediately before the closing `{{ end }}` of the home template (after the existing footer block).
- [X] T008 [P] Append a CSS block to `src/frontend/static/styles/styles.css` containing rules for: `.product-search` (margin / wrapper layout above the grid), `.product-search-input` (full-width or grid-aligned width, matches existing Bootstrap input look), `.product-search-empty` (visible empty-state message styling, hidden by default via the `[hidden]` attribute), and a generic `.visually-hidden` utility (`position:absolute; width:1px; height:1px; padding:0; margin:-1px; overflow:hidden; clip:rect(0 0 0 0); white-space:nowrap; border:0;`) used by the label and live-region.

**Checkpoint**: DOM contract is rendered on every page load. The `product-search.js` scaffold attaches but does no filtering yet. The page looks like before, plus an empty search input above the grid.

---

## Phase 3: User Story 1 - Filter the product list as I type (Priority: P1) 🎯 MVP

**Goal**: As a shopper types in the search input, the visible product cards narrow to only those whose names contain the (trimmed, case-insensitive) query — no page reload, no network request.

**Independent Test** (from `quickstart.md` 1.1–1.3): Type `sunglasses` → only the Sunglasses card remains visible. Type `mu` then `g` → grid narrows on each keystroke with no reload. Type `Sun` → "Sunglasses" remains visible (case-insensitive match).

### Implementation for User Story 1

- [X] T009 [US1] In `src/frontend/static/js/product-search.js`, attach an `input` event listener to `#product-search-input` that does the following on every event:
  1. Compute `const q = input.value.trim().toLowerCase()`.
  2. If `q === ""`: remove the `hidden` attribute from every `.hot-product-card`; return without touching `.product-search-empty` or `#product-search-status` in this phase (those are handled in US3).
  3. Otherwise: for each `.hot-product-card`, set `card.hidden = !card.dataset.name.includes(q)` (uses the lowercased `data-name` written by T003, so no per-keystroke `toLowerCase()` on each card).

**Checkpoint**: US1 is fully functional. Story 1 acceptance scenarios pass. SC-002 (perceived 100 ms update) and SC-003 (zero new network requests) are verifiable in DevTools.

---

## Phase 4: User Story 2 - Clear the search to see everything again (Priority: P1)

**Goal**: Emptying the search input restores every product card to visible in its original order.

**Independent Test** (from `quickstart.md` 2.1–2.2): After filtering by `mug`, delete all characters → every product visible. After `zzzz` filter, click the input's native × button → every product visible, empty-state gone.

### Implementation for User Story 2

- [X] T010 [US2] Verify (no new code) that the empty-query branch of the listener written in T009 correctly restores all cards when the input is emptied via either: (a) backspace to empty, (b) select-all + delete, or (c) the browser-native `type="search"` × clear button. The native × button is required to fire an `input` event on every evergreen browser (Chrome, Edge, Safari, Firefox); if any target browser is observed not to fire one in the future, add an explicit `addEventListener('search', …)` fallback in `src/frontend/static/js/product-search.js` that re-runs the same listener body. Document the verification result inline as a code comment in `product-search.js`.

**Checkpoint**: US1 + US2 both pass independently. The search box is now safe to use — typing narrows, clearing restores.

---

## Phase 5: User Story 3 - See a clear "no results" state (Priority: P2)

**Goal**: When the query matches zero products, the grid area shows a "No products match your search." message in place of the (empty) grid, and the polite ARIA live region announces the same.

**Independent Test** (from `quickstart.md` 3.1–3.2): Type `zzzz` → grid empty, "No products match your search." visible. Clear → message gone, full grid restored.

### Implementation for User Story 3

- [X] T011 [US3] Extend the listener in `src/frontend/static/js/product-search.js` (T009) to also, after the visibility pass, (a) compute `const visibleCount = document.querySelectorAll('.hot-product-card:not([hidden])').length`, (b) set `emptyEl.hidden = !(q !== "" && visibleCount === 0)` (i.e. show empty-state only when the query is non-empty AND no card matched), (c) write the appropriate string into `#product-search-status`: `""` when `q === ""`, `"No products match your search."` when `visibleCount === 0`, `"1 product matches"` when `visibleCount === 1`, otherwise `"<visibleCount> products match"`. This satisfies both FR-006 (visible empty-state message) and FR-006a (live-region announcement).

**Checkpoint**: All three user stories are independently functional. Spec acceptance scenarios 1.1 – 3.2 pass. Edge cases E1 – E6 from the spec / quickstart pass.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Verify the feature against the spec's measurable success criteria and the quickstart checklist.

- [ ] T012 [P] Run the `quickstart.md` browser verification checklist end-to-end (scenarios 1.1 – 3.2 and edge cases E1 – E6); record results in the relevant rows.
- [ ] T013 [P] Verify SC-003 (zero additional network requests) in DevTools → Network panel: clear log, type queries of varying length, including the native × clear button — the Network panel must record **0** new rows for these actions.
- [ ] T014 [P] Verify FR-006a (a11y live-region announcement) by enabling a screen reader (VoiceOver on macOS or NVDA on Windows) and confirming announcements for: query → 1 match, query → N matches, query → 0 matches, query cleared (no announcement).
- [X] T015 Re-read `src/frontend/templates/home.html` after edits and confirm: (a) the empty-input case renders byte-equivalently to the pre-feature page (SC-005), (b) `ListProducts` is still called exactly once on page load and not on keystrokes, (c) no new env vars, no new routes, no new dependencies have been added.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies. Trivially fast.
- **Foundational (Phase 2)**: Depends on Setup. Blocks every user story.
- **US1 (Phase 3)**: Depends on Foundational. The MVP.
- **US2 (Phase 4)**: Depends on US1 (Phase 3) because T010 only verifies what T009 already implements. If US2 is treated as a strict sibling MVP slice rather than a verification, it depends only on Foundational.
- **US3 (Phase 5)**: Depends on US1 (T011 extends T009's listener). Logically independent of US2.
- **Polish (Phase 6)**: Depends on all user stories the team chooses to ship.

### User Story Dependencies

- **US1**: Independent. Implements the core filter behaviour.
- **US2**: Logically depends on US1 — without filtering, "clearing" has nothing to undo. T010 is verification only.
- **US3**: Depends on US1 (extends the same listener). Independent of US2 in code.

### Within Each Story

- US1 has one task (T009). US2 has one verification task (T010). US3 has one task (T011) that extends US1's file. There is no cross-file dependency within stories.

### Parallel Opportunities

- T006 (create JS file) can run in parallel with T002–T005 (home.html edits) — different files.
- T008 (CSS append) can run in parallel with T002–T007 — different file.
- T012, T013, T014 (verification) can run in parallel with each other.
- The three user-story phases cannot meaningfully be parallelised because T009, T010, T011 all touch the same JS file in tight sequence.

---

## Parallel Example: Phase 2 Foundational

```bash
# These three can run in parallel because they touch three different files:
Task: "T006 Create src/frontend/static/js/product-search.js with DOM-ready scaffold"
Task: "T008 Append CSS rules to src/frontend/static/styles/styles.css"

# These four are sequential — all edit src/frontend/templates/home.html:
Task: "T002 Add id=hot-products-row to grid container in home.html"
Task: "T003 Add data-name attribute to each .hot-product-card in home.html"
Task: "T004 Insert search input + label + live region above the grid in home.html"
Task: "T005 Insert empty-state element after the grid in home.html"
Task: "T007 Add <script defer> tag for product-search.js in home.html"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. T001 (Setup) → fast.
2. T002 – T008 (Foundational) → the DOM contract is in place, page renders identically to before plus an inert search box.
3. T009 (US1) → MVP. The shopper can now filter by typing. **STOP and VALIDATE** in the browser.
4. Deploy/demo if ready.

### Incremental Delivery

1. Setup + Foundational → demo-able page with inert search input.
2. + US1 (T009) → **MVP shippable**: typing filters. (US2 falls out for free.)
3. + US3 (T011) → "no products match" UX and a11y announcement. Now feature-complete to spec.
4. Polish (T012 – T015) → checklist verified.

### What an LLM/single developer executes

This feature is small enough that one developer can pick up T001 through T015 in a single sitting; there is no need for parallel team strategy.

---

## Notes

- `[P]` tasks operate on different files and have no incomplete dependencies.
- `[Story]` label maps each task to a spec user story for traceability.
- No automated tests are emitted (per the "Tests" note at top of this file). The `quickstart.md` checklist is the verification source of truth.
- Commit after each logical group: after Foundational (visible scaffold), after US1 (MVP), after US3 (feature-complete), and after Polish (verified).
- Avoid: modifying any file under `src/productcatalogservice/`, `helm-chart/`, `kubernetes-manifests/`, `kustomize/`, `terraform/`, `skaffold.yaml`, `cloudbuild.yaml`, or `.github/workflows/` — those are explicitly out of scope (C-001 … C-006, C-007).
