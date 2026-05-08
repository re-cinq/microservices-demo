# Tasks: Browser-Side Product Search

**Input**: Design documents from `specs/001-browser-product-search/`  
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓  
**Tests**: Not requested in spec — no test tasks generated.  
**Branch**: `001-browser-product-search`

## Format: `[ID] [P?] [Story?] Description with file path`

- **[P]**: Can run in parallel (different files, no incomplete-task dependencies)
- **[Story]**: Which user story this task belongs to

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm the working environment before making any changes. No new dependencies or tooling are introduced for this feature.

- [x] T001 Verify branch is `001-browser-product-search` (`git branch --show-current`)
- [x] T002 [P] Confirm `src/frontend/templates/home.html` exists and matches expected structure (`.hot-product-card`, `.hot-product-card-name` classes present)
- [x] T003 [P] Confirm `src/frontend/static/styles/styles.css` exists and is writable

**Checkpoint**: Environment verified — no new packages, services, or config files are needed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add the CSS for the search box to `styles.css`. This is foundational because the HTML added in Phase 3 references the `.product-search-box` class.

- [x] T004 Add `.product-search-box` CSS rule to `src/frontend/static/styles/styles.css`:
  ```css
  .product-search-box {
    width: 100%;
    max-width: 400px;
    padding: 0.5rem 0.75rem;
    border: 1px solid #ced4da;
    border-radius: 4px;
    font-size: 1rem;
  }
  ```

**Checkpoint**: CSS in place — HTML additions in Phase 3 will render correctly.

---

## Phase 3: User Story 1 — Filter Products by Name (Priority: P1) 🎯 MVP

**Goal**: A shopper types into the search box and the product grid filters in real time to show only name-matching cards. All products reappear when the box is cleared.

**Independent Test**: Open the home page, type a partial product name (e.g. "shirt"), confirm only matching cards are visible. Clear the input, confirm all cards return.

### Implementation for User Story 1

- [x] T005 [US1] In `src/frontend/templates/home.html`, inside `.hot-products-row` and before `{{ range $.products }}`, add the search input and no-results message:
  ```html
  <div class="col-12 mb-3">
    <input
      id="product-search-input"
      type="text"
      class="product-search-box"
      placeholder="Search products..."
      aria-label="Search products"
      autocomplete="off"
    >
    <p id="no-results-message" class="text-muted mt-2" style="display:none">
      No products match your search.
    </p>
  </div>
  ```

- [x] T006 [US1] In `src/frontend/templates/home.html`, at the bottom of the `{{ define "home" }}` block (before `{{ end }}`), add the inline filter script:
  ```html
  <script>
    (function () {
      var input = document.getElementById('product-search-input');
      var noResults = document.getElementById('no-results-message');

      input.addEventListener('input', function () {
        var query = input.value.trim().toLowerCase();
        var cards = document.querySelectorAll('.hot-product-card');
        var visible = 0;

        cards.forEach(function (card) {
          var name = card.querySelector('.hot-product-card-name');
          var matches = !query || (name && name.textContent.toLowerCase().indexOf(query) !== -1);
          card.style.display = matches ? '' : 'none';
          if (matches) visible++;
        });

        noResults.style.display = (query && visible === 0) ? '' : 'none';
      });
    })();
  </script>
  ```

**Checkpoint**: US1 complete. Smoke test: type a name fragment → only matching cards shown. Clear → all cards return. Type unmatched string → "No products match your search." appears.

---

## Phase 4: User Story 2 — Case-Insensitive Matching (Priority: P2)

**Goal**: Filtering works regardless of the capitalisation the user types.

**Independent Test**: Type a product name in all-lowercase (e.g. "sunglasses") and confirm the product card with that name (capitalised or not) is still visible.

### Implementation for User Story 2

- [x] T007 [US2] Verify the script added in T006 uses `.toLowerCase()` on both the query and the card name before comparing. No additional file change is required — this is already encoded in the T006 script. Mark complete after confirming the call `name.textContent.toLowerCase().indexOf(query)` is present in `src/frontend/templates/home.html`.

**Checkpoint**: US2 complete. Smoke test: type "sunglasses" (lowercase) → "Sunglasses" card (or equivalent) is visible.

---

## Phase 5: User Story 3 — Keyboard-Accessible Search Box (Priority: P3)

**Goal**: A keyboard-only user can tab to the search input and filter the grid without a mouse.

**Independent Test**: Tab through the home page without a mouse. The search input receives focus before any product card. Type a query — the grid filters.

### Implementation for User Story 3

- [x] T008 [US3] Verify the `<input>` added in T005 in `src/frontend/templates/home.html` has:
  - `aria-label="Search products"` attribute present
  - It is placed in the DOM **before** `{{ range $.products }}` so it appears in tab order ahead of any product card link
  - No `tabindex="-1"` or other attribute that would remove it from the tab sequence
  Mark complete after confirming both conditions in the file.

**Checkpoint**: US3 complete. Smoke test: Tab key reaches the search box before the first product link. Typing filters the grid without needing Enter or a mouse click.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification that the feature integrates cleanly with the existing page.

- [x] T009 [P] Visually verify the search box aligns with the existing "Hot Products" heading on desktop viewport (≥ 992 px wide).
- [x] T010 [P] Visually verify the search box renders correctly on mobile viewport (< 576 px wide) — it should use `width: 100%` from the CSS rule and not overflow.
- [x] T011 Verify no browser console errors appear on the home page with the search feature active (open DevTools → Console, load page, type a query).
- [x] T012 Verify no regression on the product detail page (`/product/{id}`) — navigation from a filtered card still works correctly.

**Checkpoint**: All stories functional, no regressions, no console errors.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1 completion. Blocks Phase 3+ (HTML references the CSS class).
- **Phase 3 (US1)**: Depends on Phase 2. Core deliverable — ship this alone for an MVP.
- **Phase 4 (US2)**: Depends on Phase 3 (verifies code written in T006). Parallel possible with Phase 5.
- **Phase 5 (US3)**: Depends on Phase 3 (verifies markup written in T005). Parallel possible with Phase 4.
- **Phase 6 (Polish)**: Depends on Phases 3–5 complete.

### User Story Dependencies

| Story | Depends on | Can parallel with |
|-------|-----------|-------------------|
| US1 (P1) | Phase 2 | — |
| US2 (P2) | US1 (T006) | US3 |
| US3 (P3) | US1 (T005) | US2 |

### Within Each Phase

- T005 and T006 must be done in order within Phase 3 (T006 script references IDs added in T005).
- T007 and T008 (Phases 4 & 5) are verification tasks — no file edits required if T005/T006 were done correctly.

### Parallel Opportunities

- T002 and T003 (Phase 1) can run in parallel.
- T009 and T010 (Phase 6) can run in parallel.
- US2 and US3 phases can be verified in parallel once US1 is complete.

---

## Parallel Example: Phases 4 & 5 (after Phase 3 done)

```
Task: "Verify case-insensitive indexOf in home.html" (T007 / US2)
Task: "Verify aria-label and DOM position of input in home.html" (T008 / US3)
```

Both read the same file with no writes — safe to run in parallel.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (confirm environment)
2. Complete Phase 2: Add CSS to `styles.css`
3. Complete Phase 3: Add search box + JS to `home.html`
4. **STOP and VALIDATE**: Load home page, test filtering manually
5. Ship — US1 alone is the full end-user value

### Incremental Delivery

1. Phase 1 + 2 + 3 → MVP (filtering works)
2. Phase 4 → confirm case-insensitivity (likely already working from T006)
3. Phase 5 → confirm keyboard accessibility (likely already working from T005)
4. Phase 6 → polish sign-off

### Single Developer

All tasks can be completed sequentially T001 → T012 in order. Estimated time: 30–60 minutes.

---

## Notes

- T007 and T008 are verification-only — if T005 and T006 were implemented exactly per plan.md, they require no additional edits.
- No Go source files change. No proto files change. No infra files change.
- The IIFE wrapper in T006 is important — it avoids global scope pollution without a module bundler.
- `display: ''` (empty string) is used instead of `'block'` so Bootstrap's `col-md-4` responsive classes remain authoritative over layout.
