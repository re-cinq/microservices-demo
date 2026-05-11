---
description: "Task list for feature: Browser-Side Product Search"
---

# Tasks: Browser-Side Product Search

**Input**: Design documents from `/specs/001-browser-product-search/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/ui-contract.md](./contracts/ui-contract.md), [quickstart.md](./quickstart.md)

**Tests**: TDD was **not** requested for this feature, and the spec's Out-of-Scope section forbids touching the build pipeline (i.e. no new JS test runner in CI). Verification is therefore manual via the quickstart walk-through. Dedicated test tasks are omitted by design — they would not produce a runnable test suite under the existing constraints.

**Organization**: Tasks are grouped by user story (US1, US2, US3) so each story can be delivered as an independent increment. US1 is the MVP.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to (US1, US2, US3); omitted on Setup/Foundational/Polish
- Every task includes an exact file path

## Path Conventions

This is a web app mapped onto the repo's existing `src/<service>/` layout. The only service touched is `src/frontend/`. No new top-level directories. See [plan.md → Source Code](./plan.md) for the full file list.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the one new directory this feature needs. No language init, no toolchain — those are forbidden by the spec's Out-of-Scope section.

- [X] T001 Create the new static-assets JS directory at `src/frontend/static/js/` (e.g. `mkdir -p src/frontend/static/js/`).

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: DOM contract, styling, and JS scaffolding that every user story depends on. Until this phase completes, none of US1/US2/US3 can be implemented or tested.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T002 Modify `src/frontend/templates/home.html` per [contracts/ui-contract.md §1](./contracts/ui-contract.md): (a) add `<input type="search" id="product-search" class="product-search" aria-label="Search products" aria-controls="hot-products-grid" autocomplete="off">` immediately above the `.hot-products-row` element; (b) add `id="hot-products-grid"` to that grid container (existing classes preserved); (c) inside the `{{ range $.products }}` loop, add `data-name="{{ .Item.Name | lower }}"` (or the template-correct lowercased name expression for this codebase) to each product card and ensure it carries a stable class (e.g. `hot-product-card`); (d) immediately after the grid container, add `<p id="product-search-no-results" class="no-results" role="status" aria-live="polite" hidden>No products match your search.</p>`; (e) before `{{ template "footer" . }}` add `<script defer src="/static/js/product-search.js"></script>`.

- [X] T003 [P] Add minimal styles for the new elements to `src/frontend/static/styles/styles.css`: a `.product-search` rule (sized to sit naturally above the grid, inheriting the page's existing visual style — colour, font, spacing) and a `.no-results` rule (centred, muted text in place of the empty grid). No layout breakage of the existing `.hot-products-row` grid.

- [X] T004 [P] Create `src/frontend/static/js/product-search.js` with: (a) an IIFE wrapper so nothing leaks to global scope; (b) a `DOMContentLoaded` guard that early-returns if `#product-search` is absent (script must be safe on pages that don't have the box); (c) the pure-function export per [contracts/ui-contract.md §2.3](./contracts/ui-contract.md): `window.__productSearch = { matches(dataName, query) { const n = query.trim().toLowerCase(); return n === "" ? true : dataName.includes(n); } };`; (d) element selectors for `#product-search`, `#hot-products-grid`, all `.hot-product-card` elements, and `#product-search-no-results` (cached on init). **Do NOT attach the input listener yet — that lands in T005.**

**Checkpoint**: Foundation ready. DOM contract is in place, styles loaded, JS scaffold present with the pure function callable from the console. All three user stories can now proceed.

---

## Phase 3: User Story 1 — Filter the visible product list by typing a name (Priority: P1) 🎯 MVP

**Goal**: A shopper types into the search box and the product grid narrows to only products whose `name` matches (case-insensitive substring). When no product matches, a no-results message replaces the grid.

**Independent Test**: Load the home page, type `sun` → only "Sunglasses" remains. Type `xyzzy` → grid is empty and "No products match your search." is shown. Delete characters → list returns. (Maps to [quickstart.md §2.1–§2.5](./quickstart.md).)

### Implementation for User Story 1

- [X] T005 [US1] In `src/frontend/static/js/product-search.js`, attach an `input` event listener to `#product-search` that, on every event, computes `normalized = value.trim().toLowerCase()` and iterates over the cached `.hot-product-card` list, removing the `hidden` attribute when `card.dataset.name.includes(normalized)` (or when `normalized === ""`) and setting `hidden` otherwise. Implementation must reuse the pure function created in T004 (`window.__productSearch.matches`). Satisfies FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-012.

- [X] T006 [US1] In `src/frontend/static/js/product-search.js`, extend the listener body added in T005 to track `visibleCount` during the iteration and toggle the `hidden` attribute on `#product-search-no-results` accordingly (hidden when `visibleCount > 0` OR `normalized === ""`; visible when `visibleCount === 0` AND `normalized !== ""`). The query MUST remain in the input — never clear `inputElement.value` from inside the listener. Satisfies FR-008.

**Checkpoint**: At this point, User Story 1 is fully functional and independently testable. Run quickstart steps §1, §2.1–§2.5, §3, §4 to verify before moving on. This is a valid shippable MVP — the spec's User Story 2 and 3 are usability refinements on top.

---

## Phase 4: User Story 2 — Resetting the filter (Priority: P2)

**Goal**: A shopper can return to the full unfiltered list with a single action — using the native clear "×" inside `<input type="search">` or pressing Escape while the box is focused.

**Independent Test**: With `watch` typed and the list filtered, press Escape (with the box focused) or click the native clear button. Box becomes empty, full list returns, no-results hidden. (Maps to [quickstart.md §2.6](./quickstart.md).)

### Implementation for User Story 2

- [X] T007 [US2] In `src/frontend/static/js/product-search.js`, add an explicit `keydown` handler on `#product-search`: when `event.key === "Escape"`, set `inputElement.value = ""` and dispatch a new `Event("input", { bubbles: true })` so the existing US1 listener runs and restores the full list. (The native `type="search"` clear "×" already produces an `input` event on its own — no extra wiring needed for the button path; this task adds Escape for browsers that don't map it natively.)

**Checkpoint**: Native clear and explicit Escape both reset the filter. User Stories 1 + 2 are independently functional.

---

## Phase 5: User Story 3 — Filter resets on page reload (Priority: P3)

**Goal**: The filter is transient. Reloading the product list page starts fresh with an empty search box and the full list visible. The query is never persisted to storage, cookies, or URL.

**Independent Test**: Type `tank`, reload (F5 / Cmd-R). Search box is empty, all products visible. (Maps to [quickstart.md §2.7](./quickstart.md).)

### Implementation for User Story 3

- [X] T008 [US3] Audit `src/frontend/static/js/product-search.js` to confirm it does **not** read from or write to `localStorage`, `sessionStorage`, `document.cookie`, `window.location`, or `history` — and does not add an `unload`/`beforeunload`/`pagehide` listener. This story is implemented by *not* persisting state; T008 is the guardrail. If any persistence call is present, remove it. (FR-010, SC-008.)

**Checkpoint**: All three user stories independently functional. The whole feature is ready for the final polish gate.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Run the manual verification suite end-to-end and confirm the spec's Out-of-Scope constraints are honoured before opening a PR.

- [X] T009 Run the smoke check (DOM contract verification) per [specs/001-browser-product-search/quickstart.md §1](./quickstart.md): confirm the search input, grid container with `id`, every card carrying `data-name` and a stable class, and the no-results element are all present in the rendered home page.

- [X] T010 [P] Run all acceptance scenarios for US1, US2, US3 per [specs/001-browser-product-search/quickstart.md §2](./quickstart.md). Every scenario must produce its expected outcome; record any deviation as a defect against the corresponding FR.

- [X] T011 [P] Run the pure-function console assertions per [specs/001-browser-product-search/quickstart.md §3](./quickstart.md). Final console message must be `product-search pure function: OK` with zero `Assertion failed` lines.

- [X] T012 [P] Run the "zero additional server requests per keystroke" check per [specs/001-browser-product-search/quickstart.md §4](./quickstart.md). Verifies SC-006.

- [X] T013 [P] Run the accessibility check per [specs/001-browser-product-search/quickstart.md §5](./quickstart.md) (axe DevTools preferred; minimal inspection acceptable). Verifies SC-007 — zero new violations on the product list page.

- [X] T014 Run the constraint sweep per [specs/001-browser-product-search/quickstart.md §6](./quickstart.md) — `git diff --name-only main...HEAD` filtered against `src/productcatalogservice/`, `kubernetes-manifests/`, `helm-chart/`, `kustomize/`, `terraform/`, `cloudbuild.yaml`, `skaffold.yaml`, `.github/`, and `protos/`. Every one of those filters MUST be empty. Every line of the unfiltered diff MUST begin with `src/frontend/` or `specs/001-browser-product-search/` (CLAUDE.md edit is also acceptable). This is the final gate before opening the PR.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies. Trivial — one `mkdir`.
- **Phase 2 (Foundational)**: Depends on Phase 1. **Blocks all user stories.**
- **Phase 3 (US1 — MVP)**: Depends on Phase 2. Can be shipped independently of US2/US3.
- **Phase 4 (US2)**: Depends on Phase 2. Does **not** depend on US1 being merged, but in practice US1 will land first because it is the MVP.
- **Phase 5 (US3)**: Depends on Phase 2. Independent of US1 and US2 — it is implemented by *not* adding persistence.
- **Phase 6 (Polish)**: Depends on all desired user stories being implemented. T014 (constraint sweep) MUST be the last task run before opening a PR.

### Within Each User Story

- US1: T005 before T006 (the no-results toggle needs the loop body that T005 introduces; both touch the same file so they cannot be `[P]`).
- US2: T007 stands alone (touches the JS file added in T004).
- US3: T008 is an audit task — no implementation file changes expected; if it does flag something, fix it in place before continuing.

### Parallel Opportunities

- **Within Phase 2**: T003 (CSS) and T004 (JS scaffold) can run in parallel — they touch different files. T002 (template) must complete before any user-story work but can overlap with T003/T004 because no other task in Phase 2 edits `home.html`.
- **Within Phase 6**: T010, T011, T012, T013 can run in parallel — they exercise different DevTools facets of the running app. T009 is the smoke check (run first); T014 is the final gate (run last).
- **Across stories**: with multiple developers, US2 (T007) and US3 (T008) can land in parallel once US1 (T005, T006) is in. They touch the same file (`product-search.js`) so coordinate via small focused commits.

---

## Parallel Example: Phase 2

```bash
# After T002 (template edit) is in flight or done, start the other two foundational tasks together:
Task: "Add .product-search and .no-results styles to src/frontend/static/styles/styles.css"
Task: "Create scaffold src/frontend/static/js/product-search.js with IIFE, DOMContentLoaded guard, pure-function export, and element selectors"
```

## Parallel Example: Phase 6 (verification fan-out)

```bash
# Once T009 smoke check passes, fan out the verification:
Task: "Run acceptance scenarios for US1/US2/US3 per quickstart §2"
Task: "Run pure-function console assertions per quickstart §3"
Task: "Run zero-extra-requests Network-tab check per quickstart §4 (SC-006)"
Task: "Run accessibility check per quickstart §5 (SC-007)"
# Then T014 last — git-diff constraint sweep.
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. T001 — make the JS directory.
2. T002, T003, T004 — foundational (template, CSS, JS scaffold). T003 and T004 in parallel.
3. T005, T006 — US1 implementation (input listener + no-results toggle).
4. **STOP and VALIDATE**: run T009, T010 (US1 scenarios only), T011, T012, T013.
5. If green, this is a shippable MVP. Open a PR with just the US1 scope if desired.

### Incremental Delivery

1. MVP (above) → demo and ship.
2. Add US2 (T007) → re-run T010 for US2 scenarios → ship.
3. Add US3 audit (T008) → re-run T010 for US3 scenarios → ship.
4. Final gate (T014) → open PR.

### Single-Developer Linear Path (simplest)

T001 → T002 → T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012 → T013 → T014.

Total: 14 tasks. Estimated effort: a focused half-day for a developer familiar with the codebase.

---

## Notes

- `[P]` tasks = different files, no dependencies on incomplete tasks.
- `[Story]` label maps a task to a specific user story for traceability; setup/foundational/polish phases intentionally omit it.
- Each user story should be independently completable and testable per its acceptance scenarios in [spec.md](./spec.md).
- Commit after each task or logical group (e.g. T005+T006 together as "implement US1 filter").
- The constraint sweep (T014) is the final gate. If it fails, reduce the change — do not relax the constraints.
- This task list contains 14 tasks: **1 Setup**, **3 Foundational**, **2 for US1**, **1 for US2**, **1 for US3**, **6 Polish/verification**.
