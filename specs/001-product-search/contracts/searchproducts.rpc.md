# Contract: `ProductCatalogService.SearchProducts` (existing RPC)

This contract documents the gRPC method the frontend will call. The proto definition is **not modified** by this feature; only the matching predicate inside the implementation is narrowed.

## Proto surface *(unchanged)*

Defined in `protos/demo.proto`:

```proto
service ProductCatalogService {
    ...
    rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse) {}
}

message SearchProductsRequest {
    string query = 1;
}

message SearchProductsResponse {
    repeated Product results = 1;
}
```

- **Service:** `hipstershop.ProductCatalogService`
- **Method:** `SearchProducts`
- **Target service binary:** `productcatalogservice`
- **Wire endpoint:** the existing `PRODUCT_CATALOG_SERVICE_ADDR` already configured in the frontend's Deployment env vars. **No new env var.**

## Request

| Field | Type | Required | Constraints (enforced by caller) |
|---|---|---|---|
| `query` | `string` | yes | Caller (the frontend handler) trims whitespace before sending, rejects empty, and rejects strings >200 chars (see plan Decision 5). The service itself does not re-validate length. |

## Response

| Field | Type | Notes |
|---|---|---|
| `results` | `repeated Product` | Zero or more products whose `name` contains the query as a case-insensitive substring. Ordered as in the catalog (the natural order of `products.json`). An empty array is a valid response (no-match case). |

The service does not return a gRPC error for "no matches"; that is signalled by `results` being empty.

## Behaviour change in this feature

**Before this feature:**

```go
if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) ||
    strings.Contains(strings.ToLower(product.Description), strings.ToLower(req.Query)) {
    ps = append(ps, product)
}
```

**After this feature:**

```go
if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) {
    ps = append(ps, product)
}
```

The description-substring branch is removed. This is a **behavioural narrowing**, not a contract change: the proto surface is identical, but the result set for the same input may be smaller than before. No other service in the repo calls `SearchProducts` today (the frontend will be the first), so the blast radius is contained.

## Error semantics

- **gRPC code `OK`** — including for zero-match queries. `results` is the empty slice.
- **gRPC code `Unavailable` / `DeadlineExceeded`** — possible if the service is down or slow. The frontend handler must surface these as a generic 5xx, matching the pattern used by `homeHandler` when `ListProducts` fails (`src/frontend/handlers.go:67-71`).
- The service does **not** return `InvalidArgument` for empty queries; the frontend handler must reject those before reaching here (per plan Decision 5).

## Test obligations introduced by this contract

Codified in `src/productcatalogservice/product_catalog_test.go`:

1. **Existing — must still pass after the narrowing.** `TestSearchProducts`: query `"alpha"` → 2 results from the fixture. (Both fixture products have `"Alpha"` in their *name*, so the test is name-only-compatible already.)
2. **New — must be added.** A test that proves description-only matches are excluded. Suggested shape: extend the fixture with a product whose name is `"Widget"` and description is `"alpha"`; assert that `SearchProducts("alpha")` does **not** include this widget. This will fail against the current implementation (pre-narrowing) and pass after.

The exact test code is a `/speckit-tasks` artefact, not a `/speckit-plan` artefact.
