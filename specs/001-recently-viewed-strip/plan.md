# Implementation Plan: Recently Viewed Products Strip

**Branch**: `006-aip-163` | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-recently-viewed-strip/spec.md`

## Summary

When a shopper views a product detail page, the frontend appends the product's ID to a browser cookie (`shop_recently_viewed`). On every subsequent page load, the frontend reads that cookie, fetches full product details for up to 5 IDs via the existing productcatalogservice gRPC client, and renders a "Recently viewed" horizontal strip — each tile showing the product name, image, and price in the shopper's selected currency. All changes are confined to the `src/frontend/` package; no new services, datastores, or infrastructure changes are required.

## Technical Context

**Language/Version**: Go 1.25 (frontend service)

**Primary Dependencies**: `gorilla/mux` (routing), `google/golang.org/grpc` (gRPC client), `html/template` (rendering), `logrus` (logging) — all already present; no new dependencies

**Storage**: Browser cookie `shop_recently_viewed` — pipe-separated list of product IDs (e.g., `OLJCESPC7Z|66VCHSJNUP|L9ECAV7KIM`), max 5 entries, 24-hour expiry; no server-side storage

**Testing**: Go standard `testing` package; table-driven tests matching existing `money_test.go` pattern

**Target Platform**: Linux container (unchanged from existing frontend deployment)

**Project Type**: Web service — Go HTTP frontend rendering server-side HTML via gRPC calls to backend services

**Performance Goals**: Strip renders with negligible additional latency; productcatalogservice holds products in memory so each `GetProduct` RPC completes in <5 ms; 5 sequential calls add <25 ms to page render time

**Constraints**:
- C-001: No new services — all changes in `src/frontend/`
- C-002: No new datastore — view history in `shop_recently_viewed` cookie; product data from existing in-memory productcatalogservice
- C-003: Match Go language and patterns of the frontend service (gorilla/mux, logrus, html/template, table-driven tests)
- C-004: No infra/deployment/CI changes

**Scale/Scope**: Affects all product detail page (`/product/{id}`) requests; cookie payload is negligible (<200 bytes)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The project constitution (`/.specify/memory/constitution.md`) has not been populated — it contains only the blank template. No constitution-level gates apply. Hard constraints are instead inherited from epic AIP-97 and encoded as C-001–C-004 above.

**Post-design re-check (Phase 1)**: All design decisions confirmed compliant:
- No new services introduced ✓
- No new datastore — cookie only ✓
- All code in Go, matching frontend patterns ✓
- No manifest, Helm, Terraform, or CI file changes ✓

## Project Structure

### Documentation (this feature)

```text
specs/001-recently-viewed-strip/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
src/frontend/
├── handlers.go              # productHandler: write cookie; new recentlyViewedProducts() helper
├── templates/
│   ├── recently_viewed.html # New: strip partial template
│   ├── product.html         # Modified: include recently_viewed partial
│   └── home.html            # Modified: include recently_viewed partial
└── handlers_test.go         # New: table-driven tests for cookie logic
```

**Structure Decision**: Single-service change within the existing `src/frontend/` package. No new packages or directories at the repository root. The strip template follows the same partial-include pattern as `recommendations.html`.
