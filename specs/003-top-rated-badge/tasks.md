---
description: "Task list for the Top rated badge (AIP-166)"
---

# Tasks: Top rated badge on the highest-rated product card

**Input**: Design documents from [`specs/003-top-rated-badge/`](.)
**Prerequisites**: [`plan.md`](plan.md), [`spec.md`](spec.md), [`research.md`](research.md), [`data-model.md`](data-model.md), [`contracts/ui-contract.md`](contracts/ui-contract.md), [`quickstart.md`](quickstart.md)

**Tests**: Not included. The behaviour is a 5-line view-model computation in one handler; the existing `TestListProductsReturnsRating` already exercises the rating data path. See research D4.

**Organization**: One user story, US1 = [AIP-166](https://odevo.atlassian.net/browse/AIP-166). No sibling stories in this spec.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Different file, no dependency on earlier incomplete tasks — can run in parallel.
- **[US1]**: Belongs to User Story 1.
- Each task names the exact file or command.

---

## Phase 1: Setup (Shared Preflight)

**Purpose**: Confirm the environment matches what `plan.md` assumed.

- [X] T001 Confirmed clean tree on `attendee/daniel-tufvander`; HEAD is `1f49b627 feat: AIP-161 ...`
- [X] T002 [P] Docker v29.4.1 running
- [X] T003 [P] `Rating float32` confirmed on `productView` at handlers.go:82

**Checkpoint**: Preflight clean — US1 implementation can begin. No Foundational phase: this story has no proto/contract/data prerequisites.

---

## Phase 2: User Story 1 - Top rated badge on the highest-rated card (Priority: P1) 🎯 MVP

**Goal**: Every home-page render shows a "Top rated" badge on the product card(s) tied at the maximum rating, and on no other card.

**Independent Test**: Open the cohort home page (`https://daniel-tufvander.training.gcp.re-cinq.com`) after deploy. With the current seed, the Watch card (rating 5.0) shows the "Top rated" badge; no other card does. Reload three times — same card every time. Maps to acceptance scenarios AS-1, AS-2, AS-3 in [`spec.md`](spec.md).

### Implementation for User Story 1

- [X] T004 [US1] Added `TopRated bool` to `productView` struct (handlers.go:83)
- [X] T005 [US1] Added max-rating pass + TopRated flag loop right after the existing product loop (handlers.go:94)
- [X] T006 [US1] Added `{{ if .TopRated }}<span class="hot-product-card-top-rated">Top rated</span>{{ end }}` as first child of `.hot-product-card` in `home.html` (before the `<a>`)
- [X] T007 [P] [US1] Added `.hot-product-card-top-rated` rule block in `styles.css` right after `.hot-product-card-rating`

### Optional compile-check

- [X] T008 [US1] `go build ./...` in `golang:1.25` Docker — clean

**Checkpoint**: US1 fully implemented. Run the Independent Test above on the cohort URL after the PR auto-merges.

---

## Phase 3: Polish & Cross-Cutting Concerns

- [ ] T009 [P] **DEFERRED to cohort URL**: visually verify AS-1 (only Watch shows the badge) after the PR auto-merges and deploy-attendee runs
- [ ] T010 [P] **DEFERRED to cohort URL**: optional AS-2 spot-check (force a tie at 5.0)
- [ ] T011 **DEFERRED**: final spec re-read and checklist tick once cohort URL verified

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies — start immediately
- **User Story 1 (Phase 2)**: depends on Setup; no Foundational phase needed
- **Polish (Phase 3)**: depends on US1 complete + the PR being merged and deployed

### Task-Level Dependencies inside US1

- T004 (productView extension) is the prerequisite for T005 (computation references `productView.Rating` and `.TopRated`) and for T006 (template branches on `.TopRated`)
- T005 depends on T004
- T006 depends on T004
- T007 (CSS) is parallel to T004/T005/T006 — different file
- T008 (compile-check) depends on T004, T005, T006 (it compiles the modified code)

### Parallel Opportunities

- **Phase 1**: T002 and T003 parallel
- **Phase 2**: T007 (CSS) can start in parallel with T004/T005/T006 — different file
- **Phase 3**: T009 and T010 parallel

---

## Implementation Strategy

### MVP scope

This story *is* the MVP for this cycle. There is only one user story; the parent epic AIP-95 has one more deferred story (AIP-162, product-page rating) that lives in a separate Spec Kit feature.

### Suggested execution order

1. **Phase 1 (Setup)**: ~30 seconds of preflight.
2. **Phase 2 (US1)**:
   - Track A (Go): T004 → T005 (sequential; same file).
   - Track B (template): T006 (depends on T004 for `.TopRated` to exist on the view-model).
   - Track C (CSS): T007 — parallel to both above.
   - Then T008 compile-check.
3. **Phase 3 (Polish)**: after PR merge + deploy.

### Stop-and-validate points

- After **T008**: the slice compiles and you can open a PR confidently.
- After **T009**: the deployed app behaves as expected; safe to mark the story Done in Jira.

---

## Notes

- `[P]` tasks touch different files and have no incomplete-task dependency.
- `[US1]` is the only story label in this feature.
- Do not introduce changes outside the files listed in [`plan.md` → Project Structure](plan.md) — in particular, **do not** touch `protos/**`, `src/productcatalogservice/**`, `kubernetes-manifests/**`, `helm-chart/**`, `skaffold.yaml`, `cloudbuild.yaml`, or `.github/**` (constraints C-001/C-004/C-005/C-006).
