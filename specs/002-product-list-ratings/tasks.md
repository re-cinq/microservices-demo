---
description: "Task list for star ratings on the product list (AIP-161)"
---

# Tasks: Star ratings on the product list

**Input**: Design documents from [`specs/002-product-list-ratings/`](.)
**Prerequisites**: [`plan.md`](plan.md), [`spec.md`](spec.md), [`research.md`](research.md), [`data-model.md`](data-model.md), [`contracts/service-contract.md`](contracts/service-contract.md), [`contracts/ui-contract.md`](contracts/ui-contract.md), [`quickstart.md`](quickstart.md)

**Tests**: Included. The plan calls for one focused unit test in `productcatalogservice` (research D11); no other tests are part of this story.

**Organization**: Tasks are grouped by user story. This feature has a single user story (US1 → [AIP-161](https://odevo.atlassian.net/browse/AIP-161)). Sibling story [AIP-162](https://odevo.atlassian.net/browse/AIP-162) is a separate Spec Kit feature and is not generated here.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Different file, no dependency on earlier incomplete tasks — can run in parallel.
- **[US1]**: Belongs to User Story 1. Setup, Foundational, and Polish phases carry no story label.
- Each task names the exact file or command.

## Path Conventions

This feature follows the existing Online Boutique web-service layout (`src/<service>/…`, `protos/demo.proto`). See [plan.md → Project Structure](plan.md) for the full list of files this story touches.

---

## Phase 1: Setup (Shared Preflight)

**Purpose**: Confirm the working environment matches what `plan.md` and `research.md` assumed.

- [X] T001 Confirm clean working tree on branch `attendee/daniel-tufvander` (run `git status` from repo root `/Users/danieltufvander/Lovable app Wheel/microservices-demo`; expect no uncommitted changes outside the `specs/002-product-list-ratings/` scaffold)
- [X] T002 [P] Verify the proto toolchain is available — **resolved via Docker**: local install of `protoc`/`go` is not required. T005, T006, and T009 will run inside a `golang:1.22` container with `protoc-gen-go v1.34.2` and `protoc-gen-go-grpc v1.5.1` (matching the versions used to generate the existing stubs). Docker v29.4.1 is present.
- [X] T003 [P] Verify field number `7` is still unused on `message Product` (run `grep -n "= 7;" protos/demo.proto`; the only matches should be in unrelated messages — `Product` must currently use only field numbers 1–6)

**Checkpoint**: Preflight clean — Foundational phase can begin.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Land the gRPC contract change. Nothing in US1 can be implemented before the regenerated Go stubs expose the `Rating` accessor.

**⚠️ CRITICAL**: T004 must complete before T005 and T006. T005 and T006 can run in parallel with each other but both depend on T004.

- [X] T004 Add `float rating = 7;` to `message Product` in `protos/demo.proto` per [`contracts/service-contract.md`](contracts/service-contract.md) (keep the inline comment from that file: `// 0.0 means "no rating"; otherwise 0.5..5.0 in 0.5 steps.`); do not touch any other message
- [X] T005 [P] Regenerate gRPC stubs for the catalog service — done inside `golang:1.25` container with `protoc-gen-go v1.34.2` / `protoc-gen-go-grpc v1.5.1`; `src/productcatalogservice/genproto/demo.pb.go` now exposes `Product.Rating float32` and `GetRating()` (line 450 / 527)
- [X] T006 [P] Regenerate gRPC stubs for the frontend — same container, same plugin versions; `src/frontend/genproto/demo.pb.go` mirrors the catalog stubs

**Checkpoint**: Proto contract landed and regenerated. US1 implementation can begin.

---

## Phase 3: User Story 1 - Star ratings appear on the product list (Priority: P1) 🎯 MVP

**Goal**: Every product card on the Online Boutique home page (the product list) shows a star rating sourced from `products.json`, matches what was seeded, and stays stable across page reloads.

**Independent Test**: Open the cohort home page (`https://daniel-tufvander.training.gcp.re-cinq.com` or `skaffold dev --port-forward` locally). Without clicking into any product, confirm every card shows a star rating; pick one product and verify the rating matches the value you seeded in `products.json`; reload the page three times and confirm the same products show the same ratings each time. Maps to acceptance scenarios AS-1, AS-2, AS-3 in [`spec.md`](spec.md).

### Tests for User Story 1

> Write the catalog test alongside the data seed so it can fail before the proto change is wired through, then pass.

- [X] T007 [P] [US1] Seeded all 9 products in `src/productcatalogservice/products.json`: Watch=5.0, Sunglasses/Candle Holder=4.5, Tank Top/Hairdryer/Mug=4.0, Loafers/Bamboo Jar=3.5, Salt & Pepper=2.5
- [X] T008 [US1] Added `TestListProductsReturnsRating` to `product_catalog_test.go` — loads the real `products.json` via the production loader, asserts every rating is `[0.0, 5.0]` on a `0.5` increment, and asserts at least one product has `5.0`
- [X] T009 [US1] Ran `go test` inside `golang:1.25` container — all 5 tests pass including the new `TestListProductsReturnsRating`

### Implementation for User Story 1

- [X] T010 [US1] Extended `productView` in `homeHandler` with `Rating float32`; populated from `p.GetRating()` in the existing product loop
- [X] T011 [US1] Added `renderStars(rating float32) template.HTML` next to `renderMoney`/`renderCurrencyLogo` in `handlers.go`; registered as `"renderStars": renderStars` in the existing template FuncMap. Output: spans with classes `star-full`, `star-half` (with `star-half-fill` + `star-half-empty` overlay), `star-empty`
- [X] T012 [US1] Added `{{ if gt .Rating 0.0 }} … {{ end }}` block in `templates/home.html` after `hot-product-card-price`, with `aria-label="Rated <X.X> out of 5 stars"` and `{{ renderStars .Rating }}`
- [X] T013 [P] [US1] Added `.hot-product-card-rating` rule set in `static/styles/styles.css` next to existing `.hot-product-card-*` rules (filled stars `#f5a623`, empty `#d0d0d0`, half-star via `position: absolute` 50% overlay)

**Checkpoint**: US1 fully functional. Run the Independent Test above before moving to Polish.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Manual verification of behaviour not covered by the unit test, plus a final consistency sweep.

- [ ] T014 **DEFERRED to cohort URL**: `skaffold` is not installed locally. Verification of AS-1/AS-2/AS-3 happens on `https://daniel-tufvander.training.gcp.re-cinq.com` after the PR auto-merges and Helm-upgrades the namespace (per `training/README.md`)
- [ ] T015 [P] **DEFERRED to cohort URL**: needs the deployed app to flip a product to `0.0` and re-check; do this after the first deploy lands on the cohort URL
- [ ] T016 [P] **DEFERRED to cohort URL**: needs a screen reader on the deployed page
- [ ] T017 Re-read [`spec.md`](spec.md) and confirm every Functional Requirement (FR-001 → FR-005) and Success Criterion (SC-001 → SC-004) is satisfied by the running app; tick the corresponding rows in [`checklists/requirements.md`](checklists/requirements.md) once the cohort URL has been verified

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies — start immediately
- **Foundational (Phase 2)**: T004 → {T005, T006}; T005 and T006 are parallel after T004
- **User Story 1 (Phase 3)**: depends on Foundational complete (regenerated stubs available)
- **Polish (Phase 4)**: depends on US1 complete

### Task-Level Dependencies inside US1

- T007 (seed JSON) is independent of T010/T011/T012/T013 but feeds T008
- T008 depends on T007 (test asserts seeded values)
- T009 depends on T008 (runs the new test)
- T010 (productView extension) is the prerequisite for T012 (template uses `.Rating`)
- T011 (renderStars func) is the prerequisite for T012 (template calls `renderStars`)
- T012 depends on both T010 and T011
- T013 (CSS) is parallel to T010/T011/T012 — different file

### Parallel Opportunities

- **Phase 1**: T002 and T003 parallel
- **Phase 2**: T005 and T006 parallel after T004
- **Phase 3**: T007 and T013 can start in parallel as soon as Foundational completes. T010/T011 cannot be parallel with each other (same file: `handlers.go`).
- **Phase 4**: T015 and T016 parallel

---

## Parallel Example: User Story 1

```bash
# As soon as Foundational completes, seed and CSS can move in parallel:
Task T007: edit src/productcatalogservice/products.json — add rating field to every product
Task T013: edit src/frontend/static/styles/styles.css — add .hot-product-card-rating rule

# Once T007 lands, the catalog test can be written:
Task T008: extend src/productcatalogservice/product_catalog_test.go with TestListProductsReturnsRating
```

T010, T011, T012 all touch `src/frontend/{handlers.go,templates/home.html}` and must be sequential.

---

## Implementation Strategy

### MVP scope

This feature *is* the MVP for the parent epic AIP-95. There is only one user story; there is no smaller slice to deliver first. Sibling story AIP-162 is the second slice but lives in a separate Spec Kit feature.

### Suggested execution order

1. **Phase 1 (Setup)**: ~5 minutes of preflight.
2. **Phase 2 (Foundational)**: T004 → T005 + T006. Stop here and commit the proto + regen as one coherent change.
3. **Phase 3 (US1)**:
   - Track A (data + test): T007 → T008 → T009
   - Track B (frontend): T010 → T011 → T012 (sequential; same file/template chain), in parallel with T013 (CSS)
   - Both tracks unlocked by Phase 2.
4. **Phase 4 (Polish)**: T014 → T015 + T016 → T017.

### Stop-and-validate points

- After **T009**: backend half is provably correct.
- After **T012**: full vertical slice in place locally; can be demoed.
- After **T017**: ready for PR review against `attendee/daniel-tufvander`.

---

## Notes

- `[P]` tasks touch different files and have no incomplete-task dependency.
- `[US1]` is the only story label in this feature.
- Every task names a concrete file path or shell command — none should require additional context to act on.
- Do not introduce changes outside the files listed in [`plan.md` → Project Structure](plan.md). In particular, **do not** touch `kubernetes-manifests/`, `helm-chart/`, `skaffold.yaml`, `cloudbuild.yaml`, or `.github/` (constraint C-004 in the spec).
- Commit cadence is your choice; a sensible split is one commit after T006 (contract landed), one after T009 (backend done), one after T013 (frontend done), one after T017 (verified). The implement skill may bundle these differently.
