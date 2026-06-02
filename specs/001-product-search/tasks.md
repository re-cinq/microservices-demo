# Tasks: Product Search

**Feature directory**: `specs/001-product-search`
**Spec**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md)
**Branch**: `001-product-search` (feature branch off `attendee/jost-werdenhoff`; merges back via PR — C-7)
**Created**: 2026-06-02

All changes are confined to `src/frontend/`. No Go code, `.proto`, gRPC, infra, env var, or CI change (constraints C-1…C-9). `[P]` = can be done in parallel with other `[P]` tasks in the same phase.

## Phase 1: Setup

- [x] T001 Confirm working tree is on the `001-product-search` feature branch (off `attendee/jost-werdenhoff`) and clean enough to start (`git status`). Implementation work stays on this branch and merges back via PR (C-7).
- [x] T002 Locate the home stylesheet that styles `.hot-product-card` → `src/frontend/static/styles/styles.css` (rules at lines ~331–369). Recorded for T007.

## Phase 2: Foundational (blocks all user stories)

- [x] T003 In `src/frontend/templates/home.html`, added `data-product-name="{{ .Item.Name }}"` to the product card div (`<div class="col-md-4 hot-product-card" ...>`). This is the match source for the filter (plan D-2). Same `.Item.Name` field already rendered in the card, so the template parses in the existing range context.

## Phase 3: User Story 1 — Filter products by typing a name (P1)

**Goal**: Typing in the search box narrows the grid by name, in-browser, no reload. Independently testable = MVP.

- [x] T004 In `src/frontend/templates/home.html`, added the search control above the grid (after the "Hot Products" heading): a visible `<label for="product-search">Search products</label>` and `<input type="search" id="product-search" autocomplete="off" placeholder="Search products by name">` (FR-001, FR-010, plan D-4).
- [x] T005 In `src/frontend/templates/home.html`, added the no-results element near the search control: `<p id="search-no-results" role="status" aria-live="polite" hidden>No products match your search.</p>` (FR-006, FR-011, plan D-5).
- [x] T006 In `src/frontend/templates/home.html`, added the inline `<script>` near the end of the home content (plan D-2/D-3/D-6): listens to the input's `input` event; lowercases + trims the query; for each `.hot-product-card`, substring-matches `data-product-name`; toggles the `hot-product-card--hidden` class; tracks visible count; toggles the no-results `<p>`'s `hidden`. Guarded with `if (!input) return;` for graceful degradation. Never calls `.focus()` (FR-002, FR-003, FR-004, FR-005, FR-008, FR-012).
- [x] T007 In `src/frontend/static/styles/styles.css`, added `.hot-product-card--hidden { display: none; }` (plan D-3 — also removes hidden cards from tab order/AT, satisfying FR-013) plus styling for `.product-search`, `.product-search-label`, `.product-search-input`, and `.product-search-no-results`, consistent with the existing look.

> **Verification tasks T008–T016 require the running app.** `go` is not installed locally and the frontend needs the other microservices (gRPC). These are intended to be run against the **deployed training URL** after the branch is pushed and the PR merges to `attendee/jost-werdenhoff` (CI deploy), or against a full local stack if available. They remain unchecked until then.

### US1 Checkpoint — verify scenarios 1–5 (manual, in browser)
- [ ] T008 Verify FR-002/003/004/005: load home page → type `watch` (only Watch visible) → `WATCH` (same result) → `sun` (Sunglasses visible) → clear (all products return, original order) → `zzzzz` (no cards, no-results message shown). Maps to spec scenarios 1–5, SC-001/002/003.
- [ ] T009 Verify FR-008/SC-006: open the browser network panel, type in the search box → confirm **zero** new network requests are issued.
- [ ] T009a Verify FR-014 (special characters): type `& < > " % /` and an emoji → no error, no markup execution, characters matched literally (typically no-results). Maps to scenario 7. *(Go `html/template` auto-escapes `data-product-name`; the script compares strings, so no extra code expected — confirm only.)*
- [ ] T009b Verify FR-015 (ordering): type a term matching multiple products → visible cards keep the same relative order as the unfiltered grid. Maps to scenario 9 / SC-003.
- [ ] T009c Verify FR-016 (long query): paste a 500+ character query → no error, no noticeable lag, shows the correct (usually empty) result. Maps to scenario 8.

## Phase 4: User Story 2 — Reset to the full grid (P2)

**Goal**: Shopper returns to the full grid without reload.

- [ ] T010 Verify FR-007 / User Story 2: after filtering, use the `type="search"` native clear (the "x") and also manually empty the box → in both cases all originally loaded products reappear, in original order, with no page reload (SC-003). *(No new code expected beyond T004/T006; this is a verification task — add a clear affordance only if the native one is insufficient.)*

## Phase 5: Accessibility & Polish

- [ ] T011 Verify FR-010: inspect the search input in the accessibility tree (or screen reader) → it has an accessible name ("Search products").
- [ ] T012 Verify FR-011: trigger a no-match query with a screen reader active → the "No products match your search." message is announced (aria-live polite) and is visible text.
- [ ] T013 Verify FR-012: while typing characters, focus remains in the search input the whole time (never jumps away).
- [ ] T014 Verify FR-013: after filtering to a subset, Tab through the page → focus reaches only visible (matching) product links, never a hidden card's link.
- [ ] T015 Verify FR-009/SC-005 (regression): with an empty query the home page is identical to current behaviour — product set, order, "Hot Products" heading, cart, currency selector, ad, footer all unchanged.
- [ ] T016 Constraint audit before handing off: `git diff --stat` shows changes only under `src/frontend/templates/home.html` and the one stylesheet; no new files, services, `.proto`, manifests, Helm, env vars, or CI files (C-1…C-7); no separate search page and no home-page layout restructure (C-10, C-11); still on the `001-product-search` feature branch, ready to PR back to `attendee/jost-werdenhoff`.

## Dependencies & ordering

- T001, T002 (Setup) → before everything.
- T003 (Foundational `data-product-name`) → blocks T006 (script needs the attribute) and T008.
- US1 (T004 → T005 → T006 → T007) is the MVP and must precede its checkpoint (T008, T009).
  - T004 and T005 are independent edits to `home.html` but touch the same file — do sequentially to avoid conflicts (not marked `[P]`).
  - T006 depends on T003, T004, T005. T007 depends on T002 and pairs with T006 (the class name must match).
- US2 (T010) depends on US1 being complete.
- Phase 5 (T011–T016) depends on US1 (and US2) complete; these are verification/polish tasks and may run in parallel `[P]` once the implementation tasks land.

## Definition of done
All of FR-001…FR-013 verified (T008–T015), all constraints audited (T016), all SC-001…SC-007 demonstrably met. Ready for `/speckit.implement`.
