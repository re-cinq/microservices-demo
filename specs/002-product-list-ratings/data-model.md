# Data Model: Star ratings on the product list

**Feature**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md) · **Research**: [research.md](research.md)

## Entities

### Product (extended)

Existing entity in the catalogue. Gains exactly one new field for this feature.

| Field | Type (proto / Go) | Required? | Validation | Notes |
|---|---|---|---|---|
| `id` | `string` / `string` | yes | unchanged | existing |
| `name` | `string` / `string` | yes | unchanged | existing |
| `description` | `string` / `string` | yes | unchanged | existing |
| `picture` | `string` / `string` | yes | unchanged | existing |
| `price_usd` | `Money` / `*pb.Money` | yes | unchanged | existing |
| `categories` | `repeated string` / `[]string` | yes (may be empty) | unchanged | existing |
| **`rating`** | **`float`** / **`float32`** | **no (proto3 default `0.0`)** | **`0.0 ≤ rating ≤ 5.0`, increments of `0.5`** | **new — D1** |

#### Validation rules

- `rating` MUST be `0.0`, `0.5`, `1.0`, `1.5`, `2.0`, `2.5`, `3.0`, `3.5`, `4.0`, `4.5`, or `5.0`.
- A value of `0.0` is reserved as the "no rating" sentinel — the frontend MUST NOT render a star widget for such a product (research D7).
- Values outside `[0.0, 5.0]` or not on a half-star increment are a seeding error and MUST be rejected by the catalogue service test (research D11). They are not expected at runtime because `products.json` is hand-curated.

#### Persistence

- Stored only in `src/productcatalogservice/products.json`, alongside the existing fields, lowercased as `"rating": <float>` (`jsonpb` matches proto field name → JSON key).
- Held in memory on the existing `productCatalog.catalog` field (`pb.ListProductsResponse`) by the existing `loadCatalog` path. No new struct, no new cache, no new datastore.
- The Cymbal Shops AlloyDB-backed catalogue loader (`loadCatalogFromAlloyDB`) is **out of scope** for this story — the cohort deployment uses the JSON path. If that branch is exercised later it must be amended to project `rating` from the DB column or fall back to `0.0`; tracked as a follow-up, not a blocker.

#### Relationships

None new. The rating is a property of a single product. There is no `Reviewer`, `Review`, or `Rating` aggregate entity in this story (research D4).

#### State transitions

None. The field is set once by an editor at seeding time and is read-only thereafter (research D3).

### Rating value (vocabulary, not an entity)

For clarity in code review and downstream stories, the eleven legal values map to the on-card render as:

| Value | Filled stars | Half star? | Empty stars | Renders? |
|---|---|---|---|---|
| `0.0` | 0 | no | 0 | **no** (suppressed per D7) |
| `0.5` | 0 | yes | 4 | yes |
| `1.0` | 1 | no | 4 | yes |
| `1.5` | 1 | yes | 3 | yes |
| `2.0` | 2 | no | 3 | yes |
| `2.5` | 2 | yes | 2 | yes |
| `3.0` | 3 | no | 2 | yes |
| `3.5` | 3 | yes | 1 | yes |
| `4.0` | 4 | no | 1 | yes |
| `4.5` | 4 | yes | 0 | yes |
| `5.0` | 5 | no | 0 | yes |

This table is referenced by the UI contract.

## Seeding policy

- Every product entry in `products.json` MUST be edited as part of this story to include a `rating` value (research A4 in the spec; constraint here is that seeding is the engineer's responsibility, not the user's).
- Seed values are an editorial choice. Suggested starting distribution: at least one product at `5.0`, at least one at `2.0`–`2.5`, the remainder spread between `3.5` and `4.5`. Concrete values are not part of the contract; they are not part of any test assertion beyond "the value loaded equals the value seeded".
