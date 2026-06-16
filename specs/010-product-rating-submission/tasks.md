---
description: "Task list for Rate a product and see it on the product page"
---

# Tasks: Rate a product and see it on the product page

**Input**: Design documents from `/specs/010-product-rating-submission/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/productcatalog-rating.md, quickstart.md

**Tests**: A single unit-test task is included for the rating aggregate/validation logic, as called for in plan.md and quickstart.md. No broader TDD suite was requested.

**Organization**: One user story (US1) in this feature. Showing ratings on the list view (AIP-197) and the "no ratings yet" empty state (AIP-198) are separate features.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[US1]**: Belongs to User Story 1

## Path Conventions

Web multi-service repo. Paths are repo-root relative: `protos/`, `src/productcatalogservice/`, `src/frontend/`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm the existing toolchain needed for the change

- [!] T001 ⛔ BLOCKED Confirm proto generation tooling and Go builds work — `go`, `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` are all absent from this environment. Cannot build/verify here.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared contract changes that every other task depends on

**⚠️ CRITICAL**: No US1 implementation can begin until the proto is updated and regenerated

- [X] T002 Add `float rating = 7;` and `int32 num_ratings = 8;` to the `Product` message, and add `rpc RateProduct(RateProductRequest) returns (Empty)` plus `message RateProductRequest { string product_id = 1; int32 stars = 2; }` to `protos/demo.proto` (per contracts/productcatalog-rating.md)
- [!] T003 [P] ⛔ BLOCKED Regenerate gRPC/protobuf stubs for the catalogue service into `src/productcatalogservice/genproto/` — requires `protoc` + Go plugins (absent). Run `cd src/productcatalogservice && ./genproto.sh` in a tooled environment.
- [!] T004 [P] ⛔ BLOCKED Regenerate gRPC/protobuf stubs for the frontend into `src/frontend/genproto/` — requires `protoc` + Go plugins (absent). Run `cd src/frontend && ./genproto.sh` in a tooled environment.

**Checkpoint**: Proto + generated stubs in sync across both services

---

## Phase 3: User Story 1 - Submit a star rating from the product page (Priority: P1) 🎯 MVP

**Goal**: A shopper selects 1–5 stars on the product page, submits, and immediately sees the product's average (as stars) plus the rating count.

**Independent Test**: Open a product page, submit a star value, confirm the displayed average and count update; restart the catalogue service and confirm ratings reset (in-memory).

### Implementation for User Story 1

- [X] T005 [US1] Add a mutex-guarded in-memory rating store (`map[string]*ratingAggregate` with `{sum, count}`) to the `productCatalog` struct and implement `RateProduct` (validate `stars` ∈ [1,5] → else `InvalidArgument`; unknown product → `NotFound`; otherwise update aggregate) in `src/productcatalogservice/product_catalog.go`
- [X] T006 [US1] Populate `Rating` (= sum/count) and `NumRatings` (= count) on each returned product in `GetProduct` and `ListProducts`, defaulting to 0/0 when no submissions, in `src/productcatalogservice/product_catalog.go` (clones products via proto.Clone to avoid mutating the cached catalog)
- [X] T007 [P] [US1] Add table-driven unit tests for the aggregate + validation (valid 1–5, out-of-range rejected, multiple submissions average correctly, unknown product) in `src/productcatalogservice/product_catalog_test.go` — written; not yet run (no `go`)
- [X] T008 [P] [US1] Add a `rateProduct(ctx, id, stars)` gRPC client helper in `src/frontend/rpc.go`
- [X] T009 [US1] Register `POST /product/{id}/rate` routed to `rateProductHandler` in `src/frontend/main.go`
- [X] T010 [US1] Implement `rateProductHandler` (parse `rating` form field, call `rateProduct`, redirect back to `/product/{id}`; reject invalid input without altering state) in `src/frontend/handlers.go`
- [X] T011 [US1] Render the average as stars followed by the count in parentheses (★★★★☆ (43); whole-star rounding via new `renderStars` helper) and add a 1–5 star submit form posting to `/product/{id}/rate` in `src/frontend/templates/product.html`. Unrated products show "No ratings yet".

**Checkpoint**: US1 fully functional and independently testable — this is the shippable MVP

---

## Phase 4: Polish & Cross-Cutting Concerns

- [!] T012 [P] ⛔ BLOCKED Run `go test ./...` in `src/productcatalogservice` — no `go` toolchain in this environment.
- [!] T013 ⛔ BLOCKED Run the quickstart.md manual verification against the running product page — requires building/running the services (needs Go + the regenerated stubs).

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies.
- **Foundational (Phase 2)**: T002 → then T003 and T004 (parallel). Blocks all of Phase 3.
- **User Story 1 (Phase 3)**: Starts after Phase 2.
- **Polish (Phase 4)**: After Phase 3.

### Within User Story 1

- T005 → T006 (same file, sequential).
- T007 depends on T005/T006 (tests the aggregate behaviour).
- T008 → T009 → T010 (client helper before route/handler).
- T011 depends on T004 (needs regenerated `Rating`/`NumRatings` fields).

### Parallel Opportunities

- T003 and T004 run in parallel (different service dirs).
- T007 (catalogue test) and T008 (frontend client) run in parallel (different files/services).
- T011 (template) can proceed alongside catalogue work once T004 is done.

---

## Parallel Example: User Story 1

```bash
# After Phase 2, these can run together (different files):
Task: "Unit tests for aggregate in src/productcatalogservice/product_catalog_test.go"  # T007
Task: "rateProduct client helper in src/frontend/rpc.go"                                # T008
Task: "Stars + count + submit form in src/frontend/templates/product.html"             # T011
```

---

## Implementation Strategy

### MVP (User Story 1 only — this feature)

1. Phase 1: Setup → confirm toolchain.
2. Phase 2: Foundational → proto + regen (blocks everything).
3. Phase 3: US1 → catalogue aggregate/RPC, frontend submit + render.
4. Phase 4: Validate via go test + quickstart, then ship.

Honour HC-001..HC-005 throughout: existing two Go services only, in-memory (no datastore), no infra/deploy/CI edits, no review text or moderation.

---

## Notes

- [P] = different files, no incomplete dependencies.
- All implementation tasks carry exact file paths and the [US1] label.
- This feature delivers only the product-page surface; list-view and empty-state are tracked as AIP-197 / AIP-198.
