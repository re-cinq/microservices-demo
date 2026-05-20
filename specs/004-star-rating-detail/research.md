# Research: Star Rating on Product Detail Page

**Feature**: 004-star-rating-detail
**Date**: 2026-05-20

---

## Decision 1: How does rating data flow from `products.json` to the frontend template?

**Decision**: Extend the protobuf `Product` message with a `rating` field (float32, field 7), seed the value in `products.json`, and let the existing `loadCatalogFromLocalFile` / `jsonpb.Unmarshal` path carry it through the existing `GetProduct` gRPC call. The frontend handler already receives `*pb.Product`; the template simply reads the new field.

**Rationale**:
- The catalog is already unmarshalled into `pb.Product` via `jsonpb`. Adding a proto field is the minimal, pattern-consistent change — no second parse pass, no side-channel map.
- The frontend already consumes `product.Item.*` in the template. A new field on the same struct requires no new handler plumbing, only a template addition.
- The existing `packagingInfo` pattern (a separate optional struct) was considered but rejected: it required a second HTTP call to a separate service. Ratings live in `products.json` already, so no second call is needed.
- No new service, datastore, or infrastructure is introduced — the change is purely additive to the existing gRPC contract between `productcatalogservice` and `frontend`.

**Alternatives considered**:
- **Side-channel map in catalog loader**: Parse `products.json` a second time as plain JSON to build a `map[string]float32`. Rejected — two parse passes on the same file, inconsistent with the existing pattern.
- **Separate `GetRating` RPC**: A new RPC method on the existing service. Rejected — more surface area than needed; `GetProduct` already returns everything required.
- **Frontend reads `products.json` directly**: Rejected — violates service boundary; frontend does not own the catalog data.

---

## Decision 2: Where is the `.proto` file and how are generated files updated?

**Decision**: The proto source lives at `src/productcatalogservice/genproto/demo.proto` (generated output `demo.pb.go` in the same package). The existing `genproto.sh` script runs `protoc`. For this training story, the generated `demo.pb.go` is updated in-place by hand-editing to add the `Rating float32` field, matching the minimal pattern used for all other scalar fields in that struct. The `.proto` source file is updated in sync so the two remain consistent.

**Rationale**:
- Running `protoc` in the training environment may not be set up; hand-editing the generated file is acceptable for this scope.
- Keeping `demo.proto` in sync documents the intent and avoids drift.
- The constraint "no CI configuration changes" is satisfied — no new build step is added.

---

## Decision 3: Rating scale and display format

**Decision**: 0.0–5.0 floating-point value stored as `float32` in the proto/JSON. Displayed in the template as a numeric value (e.g., "4.2 / 5") alongside a simple Unicode star character (★). No filled/empty star SVG rendering in this story — that is deferred to avoid frontend asset changes beyond the template.

**Rationale**:
- 5-star scale is the industry standard; no clarification was provided so the default is used.
- Numeric display + ★ requires no new static assets or CSS classes.
- All values in `products.json` will be seeded with a rating; the zero-value case (0.0) is handled by AIP-155, not this story.
