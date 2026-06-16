# Tasks: Category Jump Links

**Input**: Design documents from `specs/001-category-jump-links/`

**Prerequisites**: [plan.md](plan.md) · [spec.md](spec.md) · [research.md](research.md) · [data-model.md](data-model.md) · [contracts/ui-contract.md](contracts/ui-contract.md)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm existing code is understood and the working surface is clear before any changes.

- [x] T001 Read `src/frontend/handlers.go` — locate `homeHandler`, understand how `products` slice is built and passed to the template
- [x] T002 [P] Read `src/frontend/templates/home.html` — understand current product loop structure and surrounding markup
- [x] T003 [P] Read `src/productcatalogservice/products.json` — confirm `categories` field shape and all distinct category values

**Checkpoint**: Codebase surface understood — implementation can begin

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add the `CategoryGroup` view struct and grouping helper that both user stories depend on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T004 Add `CategoryGroup` struct (fields: `Name string`, `Slug string`, `Products []productView`) to `src/frontend/handlers.go`
- [x] T005 Implement `groupByCategory(products []productView) []CategoryGroup` helper in `src/frontend/handlers.go` — iterates products in order, uses `Item.Categories[0]` as primary category, falls back to `"other"` for empty categories, preserves first-appearance category order (see research.md D-001–D-003)
- [x] T006 Write unit test for `groupByCategory` in `src/frontend/handlers_test.go` (or new file `src/frontend/category_test.go`) — cover: normal grouping, multi-category product uses first, empty categories fallback, order preservation

**Checkpoint**: `go test ./...` passes in `src/frontend/` with grouping logic covered

---

## Phase 3: User Story 1 — View Category Jump Links on Page Load (Priority: P1) 🎯 MVP

**Goal**: On page load, the shopper sees a jump link bar at the top and products grouped under category headings.

**Independent Test**: Load `http://localhost:8080` — confirm jump link bar and category `<h3>` headings render; confirm products appear under the correct heading.

### Implementation

- [x] T007 [US1] In `homeHandler` (`src/frontend/handlers.go`): call `groupByCategory` on the fetched products slice and add `"categories"` key to the template data map passed to `home.html`
- [x] T008 [US1] Update `src/frontend/templates/home.html` — replace the flat `{{ range $.products }}` loop with a `{{ range $.categories }}` loop that renders: `<h3 id="{{.Slug}}">{{.Name}}</h3>` followed by product cards for `{{.Products}}`
- [x] T009 [US1] Add the jump link `<nav aria-label="Product categories">` block above the category sections in `src/frontend/templates/home.html` — one `<a href="#{{.Slug}}">{{.Name}}</a>` per category, iterating `$.categories` (see contracts/ui-contract.md)

**Checkpoint**: `http://localhost:8080` shows jump link bar + grouped products; `go test ./...` passes

---

## Phase 4: User Story 2 — Jump to Category Section (Priority: P2)

**Goal**: Clicking a jump link scrolls the viewport to the matching category section.

**Independent Test**: Click any jump link — viewport scrolls to matching `<h3 id="...">` heading; URL fragment updates.

### Implementation

- [x] T010 [US2] Verify anchor `id` values on `<h3>` headings in `src/frontend/templates/home.html` exactly match `href` fragment values on jump links (both equal `.Slug`) — fix any mismatch found in T008/T009
- [ ] T011 [US2] Manually validate scroll behaviour per quickstart.md Scenario 2: click each jump link in a running instance and confirm viewport moves to the correct section

**Checkpoint**: All jump links navigate correctly; no JS changes required (native anchor scroll)

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Validation, cleanup, and confirming no out-of-scope files were touched.

- [ ] T012 [P] Run full quickstart.md validation checklist — all 4 scenarios pass
- [ ] T013 [P] Run `go test ./...` in `src/frontend/` — confirm all tests (including T006) pass
- [x] T014 Confirm diff is limited to `src/frontend/` only: `git diff --name-only` must show no files outside `src/frontend/handlers.go`, `src/frontend/handlers_test.go` (or `category_test.go`), and `src/frontend/templates/home.html`
- [x] T015 [P] Check accessibility: `<nav>` has `aria-label="Product categories"`, headings use `<h3>` consistently with rest of page

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately; T001–T003 are parallel
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS both user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — T007 → T008 → T009 (sequential within story)
- **User Story 2 (Phase 4)**: Depends on Phase 3 (anchors must exist before scroll can be verified)
- **Polish (Phase 5)**: Depends on both user stories complete

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational — no dependency on US2
- **US2 (P2)**: Starts after US1 (anchors rendered by US1 are the prerequisite for scroll testing)

### Parallel Opportunities

- T001, T002, T003 — parallel (read-only, different files)
- T012, T013, T015 — parallel (different concerns, no shared state)

---

## Parallel Example: Setup Phase

```
Task: T001 — read handlers.go
Task: T002 — read home.html         ← run simultaneously with T001
Task: T003 — read products.json     ← run simultaneously with T001
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Phase 1: Setup (read the three files)
2. Phase 2: Foundational (add struct + helper + unit test)
3. Phase 3: User Story 1 (wire handler + update template)
4. **STOP and VALIDATE**: load page, confirm jump links and grouped products appear
5. Ship / demo

### Full Delivery

1. MVP above →
2. Phase 4: User Story 2 (verify scroll, fix any anchor mismatches)
3. Phase 5: Polish (full quickstart checklist, diff check, a11y)

---

## Notes

- [P] tasks touch different files or concerns and have no incomplete dependencies
- [Story] label maps each task to its user story for traceability
- No test framework setup needed — `go test ./...` already works in `src/frontend/`
- No infrastructure, manifest, or CI files should appear in `git diff`
- Commit after each checkpoint to keep rollback easy
