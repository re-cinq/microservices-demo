# Data Model: Star Rating on Product Detail Page

**Feature**: 004-star-rating-detail
**Date**: 2026-05-20

---

## Entities

### Product (extended)

Existing entity in `productcatalogservice`. This feature adds one field.

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| id | string | products.json / proto field 1 | Unchanged |
| name | string | products.json / proto field 2 | Unchanged |
| description | string | products.json / proto field 3 | Unchanged |
| picture | string | products.json / proto field 4 | Unchanged |
| price_usd | Money | products.json / proto field 5 | Unchanged |
| categories | []string | products.json / proto field 6 | Unchanged |
| **rating** | **float32** | **products.json / proto field 7** | **New — 0.0–5.0** |

**Validation rules**:
- `rating` must be in the range 0.0–5.0 inclusive.
- A value of 0.0 indicates no rating (handled by AIP-155, not this story).
- All products seeded in `products.json` must have an explicit `rating` value before this story is deployed.

---

## State Transitions

No state transitions — `rating` is a static, read-only value for this story. It is loaded once at catalog startup and held in memory.

---

## Proto Change

File: `src/productcatalogservice/genproto/demo.proto`

```proto
message Product {
  string id = 1;
  string name = 2;
  string description = 3;
  string picture = 4;
  Money price_usd = 5;
  repeated string categories = 6;
  float rating = 7;  // NEW: 0.0–5.0; 0.0 = no rating
}
```

Generated file `demo.pb.go` is updated in sync.

---

## products.json Change

Each product entry gains a `"rating"` key:

```json
{
  "id": "OLJCESPC7Z",
  "name": "Sunglasses",
  "rating": 4.3,
  ...
}
```
