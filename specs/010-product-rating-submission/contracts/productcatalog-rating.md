# Contract: ProductCatalogService rating extensions

Additions to `protos/demo.proto`. All changes are additive and backward-compatible.

## Modified message: Product

```protobuf
message Product {
    string id = 1;
    string name = 2;
    string description = 3;
    string picture = 4;
    Money price_usd = 5;
    repeated string categories = 6;

    float rating = 7;        // NEW: average stars; 0 when num_ratings == 0
    int32 num_ratings = 8;   // NEW: number of submissions; 0 == no ratings yet
}
```

## New RPC: RateProduct

```protobuf
service ProductCatalogService {
    rpc ListProducts(Empty) returns (ListProductsResponse) {}
    rpc GetProduct(GetProductRequest) returns (Product) {}
    rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse) {}
    rpc RateProduct(RateProductRequest) returns (Empty) {}   // NEW
}

message RateProductRequest {
    string product_id = 1;
    int32 stars = 2;         // whole number 1–5
}
```

### Behaviour contract

| Condition | Result |
|-----------|--------|
| `stars` in [1,5] and `product_id` exists | Aggregate updated (`sum += stars`, `count += 1`); returns `Empty`. |
| `stars` outside [1,5] | `INVALID_ARGUMENT`; no state change. |
| `product_id` not found | `NOT_FOUND`; no state change. |
| Concurrent valid calls | All applied; no lost updates (mutex-guarded). |

### Read contract (GetProduct / ListProducts)

- Each returned `Product` carries `rating` and `num_ratings` from the in-memory aggregate.
- `num_ratings == 0` ⇒ `rating == 0` ⇒ product has no ratings yet.

## Frontend HTTP contract (new)

| Method | Path | Body | Behaviour |
|--------|------|------|-----------|
| POST | `/product/{id}/rate` | form field `rating` (1–5) | Calls `RateProduct`; on success redirects back to `/product/{id}` so the updated average/count render. Invalid input → product page re-rendered without altering the aggregate. |

Existing `GET /product/{id}` is unchanged except the template now renders `rating` + `num_ratings`.
