# Tasks: Star Rating on Product Detail Page

**Input**: Design documents from `specs/004-star-rating-detail/`
**Branch**: `004-star-rating-detail`
**Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no shared dependencies)
- **[US1]**: Belongs to User Story 1 (the only story in this spec)
- Exact file paths included in every task

---

## Acceptance Criteria → Task Mapping

| Acceptance Criterion | Covered By |
|---|---|
| AC1: Star rating visible alongside product name and description | T004 |
| AC2: Displayed rating matches stored value in products.json exactly | T001 + T002 + T003 + T004 |
| AC3: Existing page layout and all other content visually unchanged | T004 (additive-only template edit) |

---

## Phase 1: Setup

**Purpose**: No new project structure needed — all changes are within existing files. No setup tasks required.

*Skipped: existing services, build tools, and pipeline are reused without modification.*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Extend the protobuf contract so the `rating` field exists on `Product` before any consumer or template work begins.

**⚠️ CRITICAL**: T001 and T002 must be complete before T003 and T004 can be verified end-to-end.

- [x] T001 Add `float rating = 7;` field to the `Product` message in `src/productcatalogservice/genproto/demo.proto`
- [x] T002 Add `Rating float32` field and standard proto3 getter/accessor to the `Product` struct in `src/productcatalogservice/genproto/demo.pb.go`, matching the pattern of existing scalar fields (e.g., `Id`, `Name`)

**Checkpoint**: `Product.Rating` is now part of the contract and will carry values from JSON through gRPC to the frontend.

---

## Phase 3: User Story 1 — Star Rating on Product Detail Page (P1) 🎯 MVP

**Goal**: A shopper opens any product detail page and sees a star rating alongside the product name and description, sourced from the product catalogue.

**Independent test**: Open the app, navigate to any product detail page, confirm a rating value (e.g., "★ 4.3 / 5") is visible below the product name and above the description. Price, image, add-to-cart, and recommendations are unaffected.

- [x] T003 [P] [US1] Add a `"rating"` field (float, 0.0–5.0) to every product entry in `src/productcatalogservice/products.json` — assign realistic values between 3.5 and 4.9 for all existing products
- [x] T004 [P] [US1] Add a rating display block to `src/frontend/templates/product.html` immediately after the `<h2>` product name element, rendering `$.product.Item.Rating` as `★ X.X / 5` — use an inline `{{ if gt $.product.Item.Rating 0.0 }}` guard so a zero value renders nothing (zero-value case handled by AIP-155)

**Checkpoint**: Both T003 and T004 can be developed in parallel (different files, no shared state). End-to-end verification requires T001 + T002 complete first.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [x] T005 Manually verify all 10 products in `src/productcatalogservice/products.json` have a non-zero `rating` value — confirm no product was missed during T003
- [x] T006 Run `go build ./...` in both `src/productcatalogservice/` and `src/frontend/` to confirm no compilation errors introduced by T001–T004
- [x] T007 Run existing tests with `go test ./...` in `src/productcatalogservice/` and confirm `product_catalog_test.go` still passes with the new `rating` field present in JSON

---

## Dependency Graph

```
T001 ──► T002 ──┐
                ├──► (T003 + T004 can run in parallel here)
T003 ──────────┘         │
T004 ──────────────────┘
                         ▼
                    T005 ──► T006 ──► T007
```

T003 and T004 are file-independent and can be worked simultaneously once T001/T002 establish the contract.

---

## Parallel Execution

**Sprint pairing option**:
- Developer A: T001 → T002 → T005 → T006 → T007
- Developer B: T003 → T004 (start after T001 is merged)

**Solo option**: T001 → T002 → T003 + T004 (switch files) → T005 → T006 → T007

---

## Implementation Strategy

**MVP = this entire story** (AIP-153 is already the smallest independently shippable unit).
Ship T001–T007 as a single PR on branch `004-star-rating-detail`.
