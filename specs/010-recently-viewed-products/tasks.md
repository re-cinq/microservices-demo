---

description: "Task list for Recently Viewed Products (Product Page)"
---

# Tasks: Recently Viewed Products (Product Page)

**Input**: Design documents from `/specs/010-recently-viewed-products/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/recently-viewed.md, quickstart.md

**Tests**: Store unit tests are included because the plan (research Decision 6) explicitly calls for them. No other test tasks are added.

**Scope**: This file covers Jira [AIP-199](https://odevo.atlassian.net/browse/AIP-199) only — the product-page strip. The home-page strip (AIP-200) and de-duplication (AIP-201) are separate stories and are intentionally **not** included here; the foundational store and template partial are built to be reused by them.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story the task belongs to (US1)
- All paths are relative to the repository root.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish a green baseline before changes.

- [ ] T001 Establish baseline: from `src/frontend`, run `go build ./...` and `go test ./...` and confirm they pass before making changes. ⚠️ NOT RUN — no Go toolchain and no running Docker daemon in this environment.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The in-memory per-session store and its wiring — the shared component every recently-viewed story depends on.

**⚠️ CRITICAL**: User Story 1 cannot be implemented until this phase is complete.

- [X] T002 Create the in-memory store in `src/frontend/recentlyviewed.go`: define `recentlyViewedStore` with a `sync.Mutex` and `bySession map[string][]string`; implement `Record(sessionID, productID string)` (prepend, most-recent-first, **no de-duplication** — that is AIP-201) and `List(sessionID, excludeProductID string, max int) []string` (return up to `max` IDs, skipping every entry equal to `excludeProductID`; unknown session → empty slice). Add a constructor that initializes the map.
- [X] T003 Wire the store onto the server in `src/frontend/main.go`: add a `recentlyViewed *recentlyViewedStore` field to the `frontendServer` struct and initialize it where `svc` is constructed in `main()` (depends on T002).
- [X] T004 [P] Unit tests in `src/frontend/recentlyviewed_test.go` covering: most-recent-first ordering, cap at `max`, exclusion of the current product ID, duplicates preserved (no dedup), and per-session isolation / unknown-session-returns-empty (depends on T002).

**Checkpoint**: Store compiles, is reachable from `frontendServer`, and unit tests pass.

---

## Phase 3: User Story 1 - See and return to recently viewed products (Priority: P1) 🎯 MVP

**Goal**: A shopper sees a "Recently viewed" strip at the bottom of the product page (≤4 products, most-recent-first, current excluded, 50% thumbnails) and can click any product to return to it. When nothing has been viewed, the page is unchanged.

**Independent Test**: View products A→B→C; on C's page the strip shows B then A (not C); clicking B opens B's page. In a fresh session the strip is absent and the layout is unchanged. (Full steps in [quickstart.md](./quickstart.md).)

### Implementation for User Story 1

- [X] T005 [US1] In `src/frontend/handlers.go` `productHandler`: build the strip view-model — call `fe.recentlyViewed.List(sessionID(r), id, 4)`, hydrate each returned ID via `fe.getProduct(ctx, id)` **skipping any that fail to resolve** (log at debug; never error the page — FR-010), then call `fe.recentlyViewed.Record(sessionID(r), id)`; add the resulting `[]*pb.Product` to the template data under the key `recently_viewed` (depends on T002, T003).
- [X] T006 [P] [US1] Create the template partial `src/frontend/templates/recentlyviewed.html` (mirroring `recommendations.html`): define `recently_viewed`, render a heading + a row of products linking to `{{ $.baseUrl }}/product/{{.Id}}` with `{{.Picture}}`/`{{.Name}}`, using a `recently-viewed` container class for styling.
- [X] T007 [US1] In `src/frontend/templates/product.html`, include the partial at the bottom of the page (after the recommendations include at line ~77), wrapped in `{{ if .recently_viewed }}{{ template "recently_viewed" . }}{{ end }}` so an empty list renders nothing and the layout is unchanged (FR-007) (depends on T005, T006).
- [X] T008 [P] [US1] Add a CSS rule in `src/frontend/static/styles/styles.css` scoping the `.recently-viewed` strip thumbnails to 50% of the page's normal product thumbnail size (AC5), using relative sizing so the same partial scales correctly when reused on other pages.
- [ ] T009 [US1] Verify against [quickstart.md](./quickstart.md): run `go build ./...` + `go test ./...` in `src/frontend`, then manually confirm AC1–AC5 and the FR-010 missing-product edge case (depends on T005–T008). ⚠️ NOT RUN — no Go toolchain / Docker daemon available; code reviewed manually instead. Re-run once a toolchain is available.

**Checkpoint**: User Story 1 is fully functional and independently testable — this is the shippable MVP for AIP-199.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Final tidy-up; nothing here should change behaviour.

- [X] T010 [P] Confirm no new dependencies, services, datastores, manifests, or CI files were introduced (constraints C-001–C-004): `git diff --name-only` should touch only files under `src/frontend/`. ✅ Verified — all application changes are 7 files under `src/frontend/`; no go.mod change, no new service, no manifest/CI/datastore.
- [X] T011 [P] Code review pass for style consistency with surrounding frontend Go/templates (C-003): naming, error/log handling, template idioms. ✅ Store/handler mirror existing patterns; template partial and CSS mirror the recommendations strip; license headers preserved.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: After Setup. **Blocks** User Story 1.
- **User Story 1 (Phase 3)**: After Foundational.
- **Polish (Phase 4)**: After User Story 1.

### Task Dependencies

- T002 → T003, T004 (store must exist first)
- T002, T003 → T005 (handler uses the store)
- T005, T006 → T007 (include needs both the data key and the partial)
- T005–T008 → T009 (validation needs everything wired)

### Within User Story 1

- T006 (template) and T008 (CSS) are independent of T005 (handler) and of each other — they can be built in parallel; T007 and T009 gate on the others.

### Parallel Opportunities

- T004 can run in parallel with T003 (different files).
- T005, T006, and T008 can be developed in parallel (handler / template / CSS — three different files).
- T010 and T011 (polish) can run in parallel.

---

## Parallel Example: User Story 1

```bash
# After Foundational (T002–T004) is complete, start these together:
Task: "T005 build recently_viewed view-model in src/frontend/handlers.go"
Task: "T006 create src/frontend/templates/recentlyviewed.html"
Task: "T008 add .recently-viewed 50%-thumbnail rule in src/frontend/static/styles/styles.css"
# Then:
Task: "T007 include the partial in src/frontend/templates/product.html"
Task: "T009 run quickstart verification"
```

---

## Implementation Strategy

### MVP (User Story 1 only)

1. Phase 1: Setup (T001) — green baseline.
2. Phase 2: Foundational (T002–T004) — store + wiring + unit tests.
3. Phase 3: User Story 1 (T005–T009) — handler, partial, include, CSS, verify.
4. **STOP and VALIDATE** via quickstart. This is shippable AIP-199.
5. Phase 4: Polish (T010–T011) — constraint/style confirmation.

### Reuse for later stories (not in this file)

- **AIP-200 (home page)**: home handler builds the same `recently_viewed` view-model and includes the same partial — no store changes.
- **AIP-201 (dedup)**: change `Record` in `src/frontend/recentlyviewed.go` to move-to-front + de-duplicate.

---

## Notes

- [P] = different files, no incomplete dependencies.
- All changes are confined to `src/frontend/` — honouring constraints C-001–C-004 (no new service, datastore, dependency, or infra/CI change).
- Commit after each task or logical group.
