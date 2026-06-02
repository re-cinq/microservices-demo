# Tasks: Frontend Product Search

**Branch**: `feat/001-frontend-product-search`
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md) | **Quickstart**: [quickstart.md](./quickstart.md)

---

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel with other [P] tasks in the same phase
- **[Story]**: Which user story this task belongs to (US1 = filter by name, US2 = discoverable + accessible)

---

## Phase 1: Setup

**Purpose**: Confirm working baseline before making any changes.

- [x] T001 Verify the frontend service builds cleanly: run `go build ./...` in `src/frontend/`
- [x] T002 Verify the home page renders all products: open `http://localhost:8080` (or run `go run .` in `src/frontend/`) and confirm product cards are visible

**Checkpoint**: Baseline confirmed — all products render, no build errors.

---

## Phase 2: Foundational

**Purpose**: No shared infrastructure is needed for this feature. The existing Go template, Bootstrap grid, and static serving are already in place.

*(No tasks — proceed directly to Phase 3.)*

---

## Phase 3: User Story 1 — Filter products by name (Priority: P1) 🎯 MVP

**Goal**: Typing in a search box on the home page instantly hides non-matching product cards and shows an empty-state message when nothing matches.

**Independent Test**: Load the home page, type `coat` → only matching products shown; clear → all products shown; type `zzznomatch` → no cards shown and empty-state message appears. Zero network requests in DevTools during all of the above.

### Implementation

- [x] T003 [US1] Add search input and empty-state paragraph above the product grid in `src/frontend/templates/home.html`

  Insert the following block immediately after `<div class="col-12"><h3>Hot Products</h3></div>` and before the `{{ range $.products }}` loop:

  ```html
  <div class="col-12 mb-3">
    <input type="search"
           id="product-search"
           class="form-control"
           placeholder="Search products"
           aria-label="Search products"
           autocomplete="off">
  </div>
  <p id="no-results-msg"
     class="col-12 text-muted"
     style="display:none;">
    No products match your search.
  </p>
  ```

- [x] T004 [US1] Add inline filter script at the bottom of the home template in `src/frontend/templates/home.html`

  Insert the following block immediately before `{{ end }}` (the closing tag of the `{{ define "home" }}` block):

  ```html
  <script>
    (function () {
      var input    = document.getElementById('product-search');
      var cards    = document.querySelectorAll('.hot-product-card');
      var noMsg    = document.getElementById('no-results-msg');

      input.addEventListener('input', function () {
        var term    = input.value.toLowerCase();
        var visible = 0;

        cards.forEach(function (card) {
          var name = card.querySelector('.hot-product-card-name')
                        .textContent.toLowerCase();
          var show = name.includes(term);
          card.style.display = show ? '' : 'none';
          if (show) visible++;
        });

        noMsg.style.display = (term.length > 0 && visible === 0) ? '' : 'none';
      });
    })();
  </script>
  ```

- [ ] T005 [US1] Manually verify acceptance scenarios from `specs/001-frontend-product-search/quickstart.md`:
  - Test 1: type `coat` → only matching cards shown
  - Test 2: clear input → all cards return
  - Test 3: type `zzznomatch` → empty-state message shown
  - Test 4: type `COAT` → same results as `coat`
  - Test 5: open DevTools Network tab, type several characters → zero new requests

**Checkpoint**: User Story 1 is fully functional. MVP is shippable.

---

## Phase 4: User Story 2 — Discoverable and keyboard-accessible (Priority: P2)

**Goal**: The search box is visually present, labelled, and reachable by Tab key without a mouse.

**Independent Test**: Load the home page, press Tab until the search box is focused (no mouse click), type `shirt` → filtering occurs. The input has a visible placeholder.

*Note: T003 already adds `placeholder="Search products"` and `aria-label="Search products"`, satisfying the label requirement. T004's `input` listener handles keyboard input. This phase validates and optionally refines those attributes.*

### Implementation

- [ ] T006 [US2] Confirm the search input is reachable by keyboard alone: reload the page, press Tab repeatedly and verify focus reaches the `#product-search` input before the product cards
- [ ] T007 [US2] If focus order is wrong (e.g., product card links receive focus before the search box), add `tabindex="0"` explicitly to the search input in `src/frontend/templates/home.html`

**Checkpoint**: Search box is keyboard-accessible. Both user stories are complete.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T008 [P] Optional — add a CSS rule to constrain search box width on large screens in `src/frontend/static/styles/cymbal.css` (or whichever stylesheet is active): e.g. `.product-search-wrap { max-width: 400px; }` and wrap the input in a `<div class="product-search-wrap">` in `home.html`
- [ ] T009 Run `go build ./...` in `src/frontend/` to confirm no Go compilation errors were introduced
- [ ] T010 Run full quickstart.md checklist end-to-end one final time before pushing

---

## Dependencies & Execution Order

- **T001 → T002**: Baseline verification must pass before any changes
- **T003 → T004**: Search input must exist in the DOM before the script can reference it
- **T004 → T005**: Script must be in place before manual verification
- **T003, T004 → T006, T007**: US1 implementation must be complete before US2 accessibility check
- **T005, T007 → T008, T009, T010**: Both user stories complete before polish

---

## Parallel Opportunities

T001 and T002 can run in parallel (build check vs. browser check).
T008 and T009 can run in parallel (CSS tweak vs. build check).

---

## Implementation Strategy

### MVP (User Story 1 only)

1. T001 + T002 — verify baseline
2. T003 — add search input + empty-state paragraph
3. T004 — add inline filter script
4. T005 — verify all US1 acceptance scenarios
5. **Ship** — the feature is complete and testable

### Full delivery (both stories + polish)

Continue with T006 → T007 → T008 → T009 → T010.

---

## Notes

- All changes are in `src/frontend/templates/home.html` (and optionally one CSS file)
- No Go handler, route, protobuf, infra, or pipeline files change
- No new dependencies introduced
- Tests are not automated — manual verification via `quickstart.md`
