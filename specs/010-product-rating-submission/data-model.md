# Phase 1 Data Model: Rate a product and see it on the product page

All state is in-memory in `productcatalogservice` (HC-002). Nothing is persisted.

## Entity: Rating Aggregate (new, in-memory)

One per product that has received at least one submission.

| Field | Type | Rules |
|-------|------|-------|
| `product_id` | string | Map key; must match an existing catalogue product id. |
| `sum` | int64 | Running total of submitted star values. Increments by `stars` on each valid submission. |
| `count` | int64 | Number of submissions. Increments by 1 on each valid submission. |

- **Derived**: `average = sum / count` (float). Defined only when `count > 0`.
- **Storage**: `map[string]*ratingAggregate` on the `productCatalog` struct, guarded by a `sync.Mutex`.
- **Lifecycle**: Created on first valid submission for a product. No update/delete by shoppers. Entire map is lost on service restart (accepted, HC-002).

## Entity: Rating Submission (transient)

A single shopper action; not stored individually, only folded into the aggregate.

| Field | Type | Rules |
|-------|------|-------|
| `product_id` | string | Must identify an existing product. |
| `stars` | int32 | Whole number 1–5. Outside range → rejected, aggregate unchanged (FR-003). |

## Changes to existing entity: Product (proto message)

Two additive fields on `Product` in `protos/demo.proto`:

| Field | Type | Meaning |
|-------|------|---------|
| `rating` | `float` (field 7) | Average stars for the product. `0` when `num_ratings == 0`. |
| `num_ratings` | `int32` (field 8) | Count of submissions. `0` means "no ratings yet". |

- Populated by `GetProduct` and `ListProducts` from the aggregate map on each read.
- Additive field numbers (7, 8) → backward compatible with existing consumers.

## State transitions

```text
No aggregate (num_ratings = 0)
      │  first valid RateProduct(stars)
      ▼
count = 1, sum = stars, average = stars
      │  each subsequent valid RateProduct(stars)
      ▼
count += 1, sum += stars, average = sum/count   (monotonic in count)

Invalid stars (outside 1–5)  ──► no state change (rejected)
Service restart              ──► all aggregates cleared → num_ratings = 0
```

## Validation rules (consolidated)

- **V1** (FR-001/FR-003): `stars` MUST be an integer in [1,5]; otherwise reject without state change.
- **V2** (FR-002): A valid submission MUST update exactly one product's aggregate atomically (mutex-guarded).
- **V3** (FR-007): `rating`/`num_ratings` are derived only from submissions; no seeded values anywhere.
- **V4** (R3): When `count == 0`, reads return `rating = 0` and `num_ratings = 0` (the "no ratings" signal).
