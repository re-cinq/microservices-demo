# Tasks: Recently Viewed Strip on the Home Page

**Input**: Design documents from `specs/007-recently-viewed-home/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, quickstart.md ✓

**Organization**: Tasks grouped by user story. Single user story (P1) — no foundational phase needed (all infrastructure exists from AIP-156).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to

---

## Phase 1: Setup

**Purpose**: Verify the working environment and existing patterns before making changes.

- [X] T001 Verify `recentlyViewedFromCookie`, `getRecentlyViewedProducts`, and `recently_viewed` template all exist in `src/frontend/handlers.go` and `src/frontend/templates/recently_viewed.html`

---

## Phase 2: User Story 1 — Recently Viewed Strip on Home Page (Priority: P1) 🎯 MVP

**Goal**: Shoppers who have viewed products during their session see a recently-viewed strip below the hot-products section on the home page.

**Independent Test**: View one or more products, navigate to the home page — the recently-viewed strip appears with correct products and working links (see quickstart.md Scenarios 1–4).

### Implementation for User Story 1

- [X] T002 [P] [US1] Update `homeHandler` in `src/frontend/handlers.go` to fetch recently-viewed products from cookie using `recentlyViewedFromCookie` and `getRecentlyViewedProducts`, convert prices with `convertCurrency`, and pass `"recently_viewed"` to the template data map (non-fatal: log warning and omit strip on error)
- [X] T003 [P] [US1] Add recently-viewed strip to `src/frontend/templates/home.html` — insert `<div>{{ if $.recently_viewed }}{{ template "recently_viewed" $ }}{{ end }}</div>` after the closing `</div>` of the hot-products row and before the desktop footer row

**Checkpoint**: After T002 and T003, User Story 1 is fully functional. Verify using quickstart.md Scenarios 1–5.

---

## Phase 3: Polish & Validation

**Purpose**: Confirm correctness across all session states and clean up.

- [ ] T004 Manually validate all five scenarios in `specs/007-recently-viewed-home/quickstart.md` against a running frontend service [requires local/dev cluster — complete manually]

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **User Story 1 (Phase 2)**: Depends on Phase 1 completion
- **Polish (Phase 3)**: Depends on Phase 2 completion

### User Story 1 Task Dependencies

- T002 and T003 target different files (`handlers.go` vs `home.html`) and can run in **parallel**
- T004 depends on both T002 and T003 being complete

### Parallel Opportunities

```bash
# T002 and T003 can run simultaneously (different files):
Task T002: "Update homeHandler in src/frontend/handlers.go"
Task T003: "Update src/frontend/templates/home.html"
```

---

## Implementation Strategy

### MVP (only one story — complete the whole thing)

1. Complete Phase 1: Setup (T001 — verification only)
2. Complete Phase 2: User Story 1 (T002 + T003 in parallel)
3. **STOP and VALIDATE**: Run quickstart.md scenarios manually
4. Complete Phase 3: Polish (T004)

---

## Notes

- [P] tasks target different files and have no mutual dependencies
- No new files, no new dependencies, no schema changes — existing patterns only
- Error handling in T002 must be non-fatal (warn + skip strip) per research.md decision
- The `recently_viewed` sub-template is used unchanged — do not modify `recently_viewed.html`
