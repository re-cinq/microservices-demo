# Phase 0 Research: Rate a product and see it on the product page

All unknowns from the Technical Context and the one open spec marker are resolved below. Each decision stays inside the epic's hard constraints (in-memory, existing Go services, no infra changes).

## R1. Where rating aggregates live (HC-002: no datastore)

- **Decision**: Hold a per-product aggregate `{sum int64, count int64}` in a `map[string]*ratingAggregate` on the `productCatalog` struct in `productcatalogservice`, guarded by a `sync.Mutex`. The service already declares a `catalogMutex` and already serves the catalogue from memory, so this matches existing patterns.
- **Rationale**: Satisfies "work in memory over the existing product catalogue" with the least new machinery. Average is computed as `sum/count` on read.
- **Alternatives considered**: Persisting to `products.json` (rejected — it's a read-only seed file and writing it back is effectively a flat-file datastore, against the spirit of HC-002); adding Redis/DB (rejected — explicitly forbidden by HC-002).

## R2. How submission reaches the catalogue service

- **Decision**: Add one RPC, `RateProduct(RateProductRequest) returns (Empty)`, to the existing `ProductCatalogService` in `protos/demo.proto`, with `RateProductRequest { string product_id = 1; int32 stars = 2; }`. The frontend calls it from a new POST handler.
- **Rationale**: Extends the service the epic says to change; reuses the existing gRPC client wiring in `frontend/rpc.go`. Adding fields/RPCs to a proto is backward-compatible.
- **Alternatives considered**: A frontend-only store (rejected — the list-view story AIP-197 needs the same data from the catalogue, so the aggregate must live in the catalogue service); a brand-new microservice (rejected — HC-001).

## R3. Exposing the rating on reads

- **Decision**: Add `float rating = 7;` and `int32 num_ratings = 8;` to the `Product` message. `GetProduct` and `ListProducts` populate them from the in-memory aggregate (`rating = sum/count`, `num_ratings = count`); a product with no submissions returns `num_ratings = 0` and `rating = 0`.
- **Rationale**: Carries the data to both the product page (this story) and the list view (AIP-197) through the responses the frontend already consumes. `num_ratings = 0` is the unambiguous "no ratings" signal that AIP-198 will key off.
- **Alternatives considered**: A separate `GetProductRating` RPC (rejected — extra round-trip; the frontend already fetches the product).

## R4. Validation of submitted stars (FR-003)

- **Decision**: `RateProduct` rejects `stars < 1 || stars > 5` with gRPC `InvalidArgument` and leaves the aggregate untouched. The frontend surfaces only a 1–5 control, so invalid values are a defensive server-side guard.
- **Rationale**: Server-side validation is authoritative and matches the existing use of `status.Errorf(codes.…)` in `product_catalog.go`.

## R5. Concurrency (edge case: concurrent submissions)

- **Decision**: Guard every read and write of the aggregate map with the mutex so two simultaneous submissions both count with no lost update.
- **Rationale**: Standard Go in-memory concurrency; the service is already mutex-aware.

## R6. Fractional average rendering — resolves FR-006 [NEEDS CLARIFICATION]

- **Decision**: Render the average **rounded to the nearest whole star** (filled vs. empty star glyphs), with the **number of ratings in parentheses** — exactly the `★★★★☆ (43)` format in the story's acceptance criteria. No numeric average and no half-star glyph in this story.
- **Rationale**: It is the lowest-complexity option that matches the acceptance-criteria example verbatim and needs no new half-star image assets in `product.html`. The precise numeric value is not required by any acceptance scenario.
- **Alternatives considered**: Half-stars (rejected for this story — needs a half-star asset/treatment and adds rounding edge cases); stars + numeric average e.g. `★★★★☆ 4.3 (43)` (rejected — busier UI, not in the AC).
- **Note**: This resolves the open clarification with an informed default. If the product owner prefers half-stars or a numeric average, it is a small follow-up to the rendering only and does not change the data model or contracts. Flagged for confirmation.

## R7. Proto regeneration (HC-004)

- **Decision**: Regenerate `genproto` for both services from the updated `demo.proto` using the repo's existing `genproto.sh` / generation step. This is a build artifact change, not an infrastructure/deploy/CI change.
- **Rationale**: Keeps generated stubs in sync without touching manifests or pipeline config, satisfying HC-004.
