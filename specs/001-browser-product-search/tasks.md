# Tasks: Browser-Side Product Search

**Input**: Design documents from `specs/001-browser-product-search/`
**Branch**: `attendee/joel-sanmoogan`
**Date**: 2026-05-11

**Prerequisites used**: plan.md ✅ | spec.md ✅ | research.md ✅ | data-model.md ✅ | quickstart.md ✅

**Tests**: Not requested in spec. No test tasks generated. Existing Go test suites must remain green (verified in Polish phase).

**Organization**: Tasks grouped by user story (US1–US4) for independent implementation and verification.

---

## Phase 1: Setup

**Purpose**: Confirm working baseline before any changes.

- [x] T001 Confirm branch is `attendee/joel-sanmoogan` (`git branch --show-current`)
- [x] T002 [P] Run existing Go tests to establish baseline: `cd src/frontend && go test ./...` and `cd src/productcatalogservice && go test ./...`

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: Add shared CSS that all user story phases depend on for visual correctness. Must be complete before any HTML markup is added.

**⚠️ CRITICAL**: Complete before Phase 3.

- [x] T003 Add `.product-search-box` CSS rule to `src/frontend/static/styles/styles.css` — constrains max-width to 400 px, adds border-radius and padding to the search input
- [x] T004 Add `#no-results-msg` CSS rule to `src/frontend/static/styles/styles.css` — centres the empty-state paragraph, sets colour to `#666`

**Checkpoint**: `styles.css` has both new rules; no HTML changes yet.

---

## Phase 3: User Story 1 — Filter Products by Name (Priority: P1) 🎯 MVP

**Goal**: A shopper types into a search box and sees only products whose names contain the typed text; clearing the box restores all products.

**Independent Test**: Load the home page → type `"sun"` → only the Sunglasses card visible → clear input → all 9 cards visible.

- [x] T005 [US1] Add search input `<div class="col-12 product-search-box">` with `<input id="product-search-input" type="search" …>` above the `<div class="col-12"><h3>Hot Products</h3>` row in `src/frontend/templates/home.html`
- [x] T006 [US1] Add empty-state paragraph `<p id="no-results-msg" style="display:none;">No products match your search.</p>` after the `{{ end }}` closing the `{{ range $.products }}` loop in `src/frontend/templates/home.html`
- [x] T007 [US1] Add `filterProducts(query)` vanilla JS function as an inline `<script>` block at the bottom of the `{{ define "home" }}` block in `src/frontend/templates/home.html` — iterates `.hot-product-card` elements, reads `.hot-product-card-name` text content, toggles `display` property, shows `#no-results-msg` when visible count is zero and query is non-empty
- [x] T008 [US1] Wire `oninput="filterProducts(this.value)"` attribute on the `<input>` element added in T005 in `src/frontend/templates/home.html`

**Checkpoint**: US1 fully functional. Type any partial product name → list filters. Clear → all cards return. Zero match → empty-state message appears.

---

## Phase 4: User Story 2 — Case-Insensitive Matching (Priority: P2)

**Goal**: Searching `"sunglasses"` and `"SUNGLASSES"` both return the Sunglasses product.

**Independent Test**: Type `"SUNGLASSES"` → Sunglasses card visible. Type `"sunglasses"` → same result.

- [x] T009 [US2] Verify `filterProducts()` in `src/frontend/templates/home.html` applies `.toLowerCase()` to both `query.trim()` and `name.textContent.trim()` before the `indexOf` comparison — adjust if missing
- [ ] T010 [US2] Manual spot-check: type `"MUG"` (all caps) → Mug card visible; type `"mug"` → same result

**Checkpoint**: US2 satisfied; case normalisation confirmed in code and by manual test.

---

## Phase 5: User Story 3 — Instant Feedback While Typing (Priority: P2)

**Goal**: The product list updates on every single keystroke — no Enter press, no button click required.

**Independent Test**: Type one character at a time in the search box; confirm list updates after each character with no page reload (network tab stays idle).

- [x] T011 [US3] Confirm the event attribute on the `<input>` in `src/frontend/templates/home.html` is `oninput` (not `onchange` or `onkeyup`) — `oninput` fires on every value change including paste and autocomplete; update if needed

**Checkpoint**: US3 satisfied; single-character typing updates the list immediately.

---

## Phase 6: User Story 4 — Keyboard-Navigable Search (Priority: P3)

**Goal**: A keyboard-only user can Tab to the search box and type to filter without using a mouse.

**Independent Test**: Tab to the search box (no mouse), type `"sun"`, confirm Sunglasses card is the only visible result.

- [x] T012 [US4] Confirm `<input>` in `src/frontend/templates/home.html` has `aria-label="Search products"` and no `tabindex` that would disrupt natural tab order — adjust if missing
- [ ] T013 [US4] Manual keyboard-only verification: open home page, Tab until search box is focused (visible focus ring), type `"watch"` → only Watch card visible

**Checkpoint**: US4 satisfied; search box reachable and operable with keyboard alone.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final verification that all stories work together and nothing is broken.

- [ ] T014 [P] Run full quickstart.md manual verification checklist (7 steps) against a locally deployed frontend — covers all acceptance scenarios from spec.md
- [x] T015 [P] Re-run existing Go test suites to confirm no regressions: `cd src/frontend && go test ./...` and `cd src/productcatalogservice && go test ./...` and `cd src/shippingservice && go test ./...`
- [x] T016 Visual review: confirm search box is styled consistently with the rest of the page (correct font, no layout overflow at narrow viewport widths)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — adds CSS before HTML references class names
- **Phase 3 (US1)**: Depends on Phase 2 — HTML markup references CSS classes from T003/T004
- **Phase 4 (US2)**: Depends on Phase 3 — verifies behaviour of the function written in T007
- **Phase 5 (US3)**: Depends on Phase 3 — verifies event binding written in T008
- **Phase 6 (US4)**: Depends on Phase 3 — verifies attributes on element added in T005
- **Phase 7 (Polish)**: Depends on Phases 3–6 complete

### User Story Dependencies

- **US1 (P1)**: Core — all other stories build on or verify its output
- **US2 (P2)**: Independent of US3/US4; depends only on US1 implementation
- **US3 (P2)**: Independent of US2/US4; depends only on US1 implementation
- **US4 (P3)**: Independent of US2/US3; depends only on US1 implementation

> US2, US3, and US4 can be worked in parallel once Phase 3 is complete.

### Within Phase 3 (US1) — Execution Order

```
T005 (markup) → T006 (empty-state) → T007 (script) → T008 (wire event)
```

T005 and T006 can be written in the same editing pass (both in `home.html`). T007 depends on T005/T006 existing in the DOM. T008 finalises T005.

---

## Parallel Example: US2, US3, US4 (after Phase 3 complete)

```
Phase 3 complete
      │
      ├─── T009 [US2] Verify case normalisation in filterProducts()
      │
      ├─── T011 [US3] Confirm oninput event binding
      │
      └─── T012 [US4] Confirm aria-label and tab order
```

All three can be executed simultaneously — they touch no overlapping code paths.

---

## Implementation Strategy

### MVP (User Story 1 only) — 4 tasks, ~30 lines of code

1. Complete Phase 1: confirm baseline
2. Complete Phase 2: add 8 lines to `styles.css`
3. Complete Phase 3 (T005–T008): add search input, empty-state, filter script, and event wiring to `home.html`
4. **STOP and VALIDATE**: manual test — type, clear, zero-match
5. Ship: this alone satisfies the primary user need

### Full Delivery (all 4 user stories)

After MVP is validated:
- Phase 4 (T009–T010): case-insensitive verification — ~5 minutes
- Phase 5 (T011): live-update event confirmation — ~2 minutes
- Phase 6 (T012–T013): keyboard accessibility check — ~5 minutes
- Phase 7 (T014–T016): final polish and regression check

---

## Notes

- [P] = parallelisable (different files or no shared state)
- [USn] = maps task to user story for traceability
- All changes are in exactly two files: `src/frontend/templates/home.html` and `src/frontend/static/styles/styles.css`
- No new files, no backend changes, no build-pipeline changes
- Commit after Phase 3 checkpoint at minimum; commit after each phase is ideal
