# Tasks: Product Search

**Input**: Design documents from `specs/001-product-search/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [quickstart.md](./quickstart.md)

**Tests**: This task list includes ONE Go test task (T016, in the Polish phase) because [research.md](./research.md) Decision 6 commits to a Go markup-contract test as regression protection. The spec does **not** require TDD; client-side behaviour is verified manually via the 12-row table in [quickstart.md](./quickstart.md). No JS test framework is added (would violate the spirit of C-006 / C-008).

**Organization**: Tasks are grouped by user story from spec.md to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: parallelizable — different file, no dependency on an incomplete task in the same phase
- **[Story]**: maps to a user story in spec.md (`[US1]`, `[US2]`, `[US3]`)
- Every task includes an exact file path

## Path Conventions

All paths are repository-relative from `~/microservices-demo/` (this branch is `001-product-search`). The feature is confined to `src/frontend/` per [plan.md](./plan.md).

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the new files and wire them into the existing templates. After this phase the page loads `search.css` and `search.js` but neither does anything yet (skeleton files only). All five user-visible behaviours are added in later phases.

- [X] T001 Create new directory `src/frontend/static/js/` (first JS directory in the repo; mirrors the existing `src/frontend/static/styles/` per-area separation pattern — see [research.md](./research.md) Decision 2)
- [X] T002 [P] Create new skeleton file `src/frontend/static/styles/search.css` containing only the Apache 2.0 license header that matches the other CSS files in `src/frontend/static/styles/` (no rules yet — rules are added by US1/US2/US3 tasks)
- [X] T003 [P] Create new skeleton file `src/frontend/static/js/search.js` containing an IIFE wrapper `(function(){ 'use strict'; document.addEventListener('DOMContentLoaded', function(){ /* implementation added by US1/US2/US3 */ }); })();` and the Apache 2.0 license header
- [X] T004 [P] In `src/frontend/templates/header.html`, add `<link rel="stylesheet" type="text/css" href="{{ $.baseUrl }}/static/styles/search.css">` immediately after the existing `<link>` for `bot.css` so it sits with the other per-area stylesheet links
- [X] T005 [P] In `src/frontend/templates/home.html`, add `<script src="{{ $.baseUrl }}/static/js/search.js" defer></script>` on its own line immediately before the final `{{ end }}` that closes the `{{ define "home" }}` block — page-scoped, not site-wide (see [research.md](./research.md) Decision 2)

**Checkpoint**: Page loads with no visible change; DevTools Network panel confirms both new asset files are requested with HTTP 200.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add the per-card data attribute that every user story's JS will read. Single edit; no JS or CSS change in this phase.

**⚠️ CRITICAL**: T006 must complete before any [US1], [US2], or [US3] task that touches `search.js`.

- [X] T006 In `src/frontend/templates/home.html`, extend the opening tag of the product card from `<div class="col-md-4 hot-product-card">` to `<div class="col-md-4 hot-product-card" data-product-name="{{ .Item.Name }}">` (single change inside the `{{ range $.products }}` block — applies to every rendered card)

**Checkpoint**: View `/` in DevTools Elements. Every `.hot-product-card` has a `data-product-name` attribute populated with the product's name. Layout is unchanged.

---

## Phase 3: User Story 1 — Filter visible products by name as I type (Priority: P1) 🎯 MVP

**Goal**: A shopper types into a search input above the product grid and only matching products remain visible, instantly, with no extra network traffic.

**Independent Test**: Load `/`. Catalogue should be the 9 products listed in [spec.md](./spec.md) US1. Type `watch` → exactly Watch card visible. Type `WATCH` → same. Type `ass` → Sunglasses + Bamboo Glass Jar only. Type `&` → Salt & Pepper Shakers only. Delete → all 9 cards visible again. DevTools Network panel shows zero new catalogue requests during any typing.

### Implementation for User Story 1

- [X] T007 [P] [US1] In `src/frontend/templates/home.html`, insert above the `<div class="row hot-products-row px-xl-6">` row, inside the same `col-12 col-lg-12 px-10-percent` column, a `<div class="search-wrapper">` containing a single `<input type="text" id="search-input" name="search" autocomplete="off" placeholder="Search products" aria-label="Search products">` element
- [X] T008 [P] [US1] In `src/frontend/static/styles/search.css`, add rules for `.search-wrapper` (block-level container, margin-bottom to separate from the grid below) and `#search-input` (full-width of wrapper, padding for legibility, font size matching site convention, focus outline visible). The wrapper MUST NOT change the position or size of any existing element when the input is empty (verifiable per C-010)
- [X] T009 [US1] In `src/frontend/static/js/search.js`, inside the existing `DOMContentLoaded` handler, implement the filter:
  1. Cache references to `document.getElementById('search-input')` and `Array.from(document.querySelectorAll('.hot-product-card'))`.
  2. Attach an `input` event listener to the search input.
  3. In the handler, compute `const q = searchInput.value.trim().toLowerCase();`.
  4. Iterate cards: `card.style.display = card.dataset.productName.toLowerCase().includes(q) ? '' : 'none';` (empty query matches every name because `''.includes('')` is true and any string `s` satisfies `s.includes('')`).

**Checkpoint**: All 6 acceptance scenarios of [US1] in [spec.md](./spec.md) pass (rows 1–6 of the manual-verification table in [quickstart.md](./quickstart.md)). Feature is MVP-shippable at this point — US2 and US3 are pure polish on top.

---

## Phase 4: User Story 2 — Clear feedback when nothing matches (Priority: P2)

**Goal**: When the shopper's query matches no products, a clear, accessible message replaces the grid instead of a silently empty page.

**Independent Test**: With US1 working, type `zzzzz`. The grid area shows the message `No products match your search.` and zero product cards. Backspace until the query matches at least one product — message disappears, matching cards reappear. Backspace fully — full grid returns. Also: paste `<script>alert(1)</script>` — no `alert` dialog appears; DevTools Elements panel shows no live `<script>` tag in the search/empty-state region.

### Implementation for User Story 2

- [X] T010 [P] [US2] In `src/frontend/templates/home.html`, insert `<div id="search-empty-state" class="search-empty-state" role="status" aria-live="polite" hidden>No products match your search.</div>` as a sibling of `<div class="row hot-products-row px-xl-6">`, positioned **after** the row in the DOM but **before** the desktop footer block — it occupies the same visual area as the grid when shown (the grid's row is hidden by its own cards being hidden, the empty-state element appears in flow)
- [X] T011 [P] [US2] In `src/frontend/static/styles/search.css`, add `.search-empty-state` rules: typography that reads as informational (not error), centered or grid-aligned text, vertical padding so it sits where the grid would have been, visible only when the `hidden` attribute is absent (the attribute already hides it by default — no explicit `display: none` needed; the rules apply only when the attribute is removed by JS)
- [X] T012 [US2] In `src/frontend/static/js/search.js`, extend the `input` handler from T009:
  1. Track match count during the per-card iteration: `let matchCount = 0; … if (matches) { card.style.display = ''; matchCount++; } else { card.style.display = 'none'; }`.
  2. After the loop: `const emptyState = document.getElementById('search-empty-state'); if (q !== '' && matchCount === 0) { emptyState.removeAttribute('hidden'); } else { emptyState.setAttribute('hidden', ''); }`.
  3. The empty-state element's `textContent` MUST NOT be modified — the message is static (closes the XSS surface per [research.md](./research.md) Decision 4).

**Checkpoint**: Acceptance scenarios US2 #1–#5 in [spec.md](./spec.md) pass. Manual-verification rows 7 (`zzzzz` empty state) and 8 (XSS) from [quickstart.md](./quickstart.md) pass.

---

## Phase 5: User Story 3 — Quick clear of the active filter (Priority: P3)

**Goal**: A visible clear control and Esc keybinding both empty the input and restore the full grid in a single action.

**Independent Test**: With US1 + US2 working, type `watch`. A clear button (×) becomes visible inside the search input. Click it once — input empties, all 9 cards return, clear button hides again. Type `watch` again. Press Esc with the input focused — same result. With input empty on initial page load, the clear button is hidden (or not rendered).

### Implementation for User Story 3

- [X] T013 [P] [US3] In `src/frontend/templates/home.html`, inside the `<div class="search-wrapper">` added in T007, append a `<button type="button" id="search-clear" class="search-clear" aria-label="Clear search" hidden>×</button>` as a sibling of `<input id="search-input">`
- [X] T014 [P] [US3] In `src/frontend/static/styles/search.css`, add `.search-clear` rules: position the button on the right edge of the wrapper (relative-positioned wrapper + absolutely-positioned button, or a flex layout); button at least 32×32px hit target; visible only when `hidden` attribute is absent; cursor pointer; no border / minimal styling to match the site
- [X] T015 [US3] In `src/frontend/static/js/search.js`, extend the existing handler from T009/T012:
  1. Cache reference: `const clearBtn = document.getElementById('search-clear');`.
  2. In the `input` handler, after computing the query: `if (searchInput.value.length > 0) { clearBtn.removeAttribute('hidden'); } else { clearBtn.setAttribute('hidden', ''); }`.
  3. Bind a click listener on `clearBtn` that does: `searchInput.value = ''; searchInput.dispatchEvent(new Event('input', { bubbles: true })); searchInput.focus();`.
  4. Bind a `keydown` listener on `searchInput` that on `event.key === 'Escape'` does the same thing as the click handler (factor into a shared `clearSearch()` function so both call sites share one code path — [research.md](./research.md) Decision 5).

**Checkpoint**: All US3 acceptance scenarios in [spec.md](./spec.md) pass. Manual-verification rows 9–11 from [quickstart.md](./quickstart.md) pass. The feature is now complete per the spec.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Lock in regression protection, walk the full verification suite, and confirm constraint compliance before opening the merge to `attendee/jenny-svennefelt`.

- [X] T016 [P] Create new file `src/frontend/home_test.go` in package `main`. Test name: `TestHomeRendersSearchMarkup`. Approach: parse the home + header + footer templates directly via `template.New("").Funcs(funcMap).ParseFiles(...)`, build a minimal data map (one product named `Watch`, empty currencies/cart, `baseUrl: ""`, plus any other template variables required to render — discover by running and fixing missing-field errors), execute the `"home"` template into a `bytes.Buffer`, and assert with `strings.Contains` that the body contains each of: `id="search-input"`, `id="search-clear"`, `id="search-empty-state"`, `role="status"`, `aria-live="polite"`, the exact text `No products match your search.`, and `data-product-name="Watch"`. If header.html requires non-trivial stubbing (e.g. session helpers), scope down to parse only the `home` template's body and substitute simpler harness templates for `header` and `footer` — the test's purpose is the markup contract from [data-model.md](./data-model.md), not full integration.
- [ ] T017 [P] Walk all 12 rows of the manual-verification table in [quickstart.md](./quickstart.md) against a running build (skaffold dev or the frontend-only run path). Each row maps to a specific spec scenario / success criterion; any failing row is a regression and must be fixed before merge.
- [X] T018 [P] Run `git diff --stat $(git merge-base attendee/jenny-svennefelt HEAD)...HEAD`. Confirm the only files in the diff are:
  - `src/frontend/templates/header.html` (modified, ~1 line)
  - `src/frontend/templates/home.html` (modified)
  - `src/frontend/static/styles/search.css` (new)
  - `src/frontend/static/js/search.js` (new)
  - `src/frontend/home_test.go` (new)
  - Optionally: files under `specs/001-product-search/` (spec, plan, etc.)
  
  Any other modified path is a constraint violation and must be reverted (see the constraint-compliance self-check in [quickstart.md](./quickstart.md)).

**Checkpoint**: T016 passes via `go test ./...` from `src/frontend/`. T017 passes all 12 manual rows. T018 confirms the diff is inside the lines.

---

## Dependencies & Execution Order

### Phase dependencies

- **Phase 1 (Setup)**: no dependencies; can start immediately.
- **Phase 2 (Foundational)**: depends on Phase 1 completion (skeleton JS file must exist before T006 has anything to feed at runtime — though T006's HTML edit alone doesn't need it, the runtime contract does).
- **Phase 3 (US1)**: depends on Phase 2 completion. **MVP** — stop here for the smallest shippable increment.
- **Phase 4 (US2)**: depends on Phase 3 (the empty-state logic in T012 extends T009's iteration; it does **not** depend on US3).
- **Phase 5 (US3)**: depends on Phase 3 (the clear-button visibility logic in T015 piggy-backs on the existing `input` handler from T009; the Esc handler clears the input and dispatches `input`, which relies on T009's listener). Does **not** depend on Phase 4.
- **Phase 6 (Polish)**: depends on Phases 3 + (4 if shipping US2) + (5 if shipping US3). T016 should reflect whichever stories are landed.

### User story dependencies

- US1 is the MVP — fully independent once Phase 2 is done.
- US2 logically depends on US1 (you cannot demonstrate the empty-state without filtering working) but US2's code touches additive lines in the same files. If you implement US2 without US1, the filter just doesn't run and the empty-state element is hidden forever — no benefit.
- US3 logically depends on US1 for the same reason. US3 does **not** depend on US2 — they're independent polish layers on top of US1.

### Within each user story

- T007 ↔ T008 ↔ T009 (within US1) modify three different files (`home.html`, `search.css`, `search.js`) → can be done in parallel by three developers (marked [P]) but they're linked by a runtime contract: the JS in T009 binds to elements added in T007 and styled in T008. Logical order: T007 → T009 → T008 (the engineer should at minimum verify by hand after each).
- Same pattern for US2 (T010/T011/T012) and US3 (T013/T014/T015).
- **Cross-story file conflicts**: T009, T012, T015 all modify `src/frontend/static/js/search.js`. They must run **sequentially** in story-priority order. T007, T010, T013 all modify `src/frontend/templates/home.html` — same. T008, T011, T014 all modify `src/frontend/static/styles/search.css` — same. The [P] markers within a phase are accurate; the [P] markers do **not** carry across phases for the same file.

### Parallel opportunities (concrete)

- Phase 1: T002 + T003 + T004 + T005 in parallel after T001 (four different files).
- Phase 3 (US1): T007 + T008 + T009 in parallel (three different files).
- Phase 4 (US2): T010 + T011 + T012 in parallel (three different files).
- Phase 5 (US3): T013 + T014 + T015 in parallel (three different files).
- Phase 6: T016 + T017 + T018 in parallel (different concerns).

---

## Parallel Example: User Story 1

```text
# Once Phase 2 (T006) is complete, three developers can pick up US1 in parallel:
Developer A: T007 — insert search input markup in src/frontend/templates/home.html
Developer B: T008 — write CSS rules in src/frontend/static/styles/search.css
Developer C: T009 — implement filter in src/frontend/static/js/search.js
```

They sync at the checkpoint to walk the US1 acceptance scenarios together in a browser.

---

## Implementation Strategy

### MVP first (User Story 1 only)

1. Phase 1 (Setup): T001–T005.
2. Phase 2 (Foundational): T006.
3. Phase 3 (US1): T007–T009.
4. **STOP** — walk rows 1–6 of [quickstart.md](./quickstart.md)'s manual-verification table. If they all pass, the MVP is shippable: merge `001-product-search` back to `attendee/jenny-svennefelt` and let CI deploy.

### Incremental delivery (add polish in two more shippable increments)

1. Ship MVP per above.
2. Phase 4 (US2): T010–T012. Walk rows 7 (empty state) and 8 (XSS). Merge.
3. Phase 5 (US3): T013–T015. Walk rows 9–11 (clear button, Esc, layout unchanged). Merge.
4. Phase 6 (Polish): T016–T018. Lock in regression test + full 12-row verification + constraint-compliance diff check. Merge.

### Parallel-team strategy

With three developers and clean coordination:

- Day 0: everyone completes Phase 1 + Phase 2 together (T001–T006).
- Day 1: A takes T007 + T010 + T013 (home.html); B takes T008 + T011 + T014 (search.css); C takes T009 + T012 + T015 (search.js). Each developer works through their file from US1 → US2 → US3 in priority order. Re-sync at each US checkpoint to verify in a browser.
- Day 2: anyone picks up T016–T018 in parallel.

(Note: the parallel-team strategy assumes serialised edits per file across stories — not all-at-once. The [P] markers within a phase still hold across developers.)

---

## Notes

- All [P] markers within a phase represent file-level parallelism. They do not imply that all stories can run simultaneously.
- Every task includes an exact file path; every user-story task is labelled `[USn]`; Setup, Foundational, and Polish phases have no story label per the SpecKit format.
- Tests are scoped to the markup contract only (T016). Client-side behaviour is verified manually per the design decision in [research.md](./research.md) Decision 6 — adding a JS test framework would violate the spirit of C-006 / C-008.
- Commit after each completed task or at minimum after each phase's checkpoint. Suggested commit prefixes: `feat(frontend/search): …` for US1/US2/US3 tasks; `test(frontend/search): …` for T016; `chore(frontend/search): …` for Phase 1.
- The course / instructor branch convention is `attendee/jenny-svennefelt` — see C-007 and [quickstart.md](./quickstart.md) for the merge-back instructions.
