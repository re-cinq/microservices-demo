# Contract: SearchProducts RPC

**Phase 1 output for**: `specs/001-product-search/plan.md`  
**Date**: 2026-05-11

This contract documents the existing gRPC interface used — and extended — by
the product search feature. The proto definition is the authoritative source;
this document records the agreed behaviour for this feature.

---

## gRPC Contract

**Service**: `ProductCatalogService`  
**Proto file**: `protos/demo.proto`  
**Generated stubs**: `src/productcatalogservice/genproto/`, `src/frontend/genproto/`

### RPC: SearchProducts

```protobuf
rpc SearchProducts(SearchProductsRequest) returns (SearchProductsResponse) {}

message SearchProductsRequest {
  string query = 1;
}

message SearchProductsResponse {
  repeated Product results = 1;
}

message Product {
  string id          = 1;
  string name        = 2;
  string description = 3;
  string picture     = 4;
  Money  price_usd   = 5;
  repeated string categories = 6;
}
```

### Behaviour (post-fix)

| Input | Behaviour |
|---|---|
| `query = ""` | Returns all products (empty string matches every substring) |
| `query = "lens"` | Returns products where `name` OR `description` contains `"lens"` (case-insensitive) |
| `query = "LENS"` | Same result as `"lens"` — matching is case-insensitive |
| `query = "xyzzy123"` | Returns empty `results` list (no error) |
| Any query | Matching is substring, not whole-word; `"lens"` matches `"lenses"` |

**No proto changes required.** The change is purely in the server implementation
(`product_catalog.go`).

---

## HTTP Contract (frontend)

**Route**: `GET /`  
**Query param**: `q` (optional)

| `q` value | Server behaviour |
|---|---|
| absent or `""` | Full catalogue rendered (existing behaviour) |
| 1 non-whitespace char (e.g. `"a"`) | Blocked — render home with `searchError`; full catalogue shown |
| whitespace only (e.g. `"  "`) | Blocked — same as above |
| ≥ 2 non-whitespace chars | Call `SearchProducts`; render results or empty state |

**Response**: Always HTTP 200 with the `home` template. Errors in the upstream
gRPC call fall through to the existing `renderHTTPError` path (HTTP 500).

### Template data contract

The `home` template is extended with the following additional keys:

| Key | Type | Present when |
|---|---|---|
| `search_query` | `string` | `q` param was non-empty |
| `search_error` | `string` | Query was present but too short |
| `is_search` | `bool` | A valid search was executed |

Existing keys (`products`, `currencies`, `cart_size`, etc.) are unchanged.
