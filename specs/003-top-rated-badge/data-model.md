# Data Model: Top rated badge

**Feature**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md) · **Research**: [research.md](research.md)

## Entities

No persisted entity is added. The change is a single render-time flag on an existing view-model struct.

### `productView` (frontend view-model, extended)

Lives in `src/frontend/handlers.go`, local to `homeHandler`. Built once per request from the gRPC `Product` returned by the catalogue service. Gains one new field:

| Field | Type | Origin | Notes |
|---|---|---|---|
| `Item` | `*pb.Product` | gRPC response | existing |
| `Price` | `*pb.Money` | currency conversion | existing |
| `Rating` | `float32` | `p.GetRating()` (added by AIP-161) | existing |
| **`TopRated`** | **`bool`** | **derived: `Rating == max(Rating across slice) && max > 0`** | **new — research D1** |

#### Computation rule

```
maxRating := 0.0
for each v in ps: if v.Rating > maxRating { maxRating = v.Rating }
if maxRating > 0:
    for each v in ps: v.TopRated = (v.Rating == maxRating)
else:
    // all flags remain false (zero value)
```

Two linear passes. Order in `ps` is the order in `products.json` (no sort).

#### Properties

- **Stable**: identical input slice → identical `TopRated` distribution. Satisfies SC-003 and FR-006.
- **Tie-safe**: multiple products at the same `maxRating` all get `TopRated = true`. Satisfies FR-002 and AC-2.
- **Defensive at zero**: if every product has `Rating == 0.0`, no flag is set. Satisfies FR-003 and AC-3.
- **Local**: lives entirely inside `homeHandler`. Other handlers (`productHandler`, `cartHandler`, etc.) are not affected — consistent with A2 (home page only).
