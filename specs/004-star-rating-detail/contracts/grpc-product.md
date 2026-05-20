# Contract: ProductCatalogService — Product message

**Service**: `productcatalogservice` (gRPC)
**Proto package**: `hipstershop`
**Affected RPC**: `GetProduct`, `ListProducts`, `SearchProducts` (all return `Product`)

---

## Change

The `Product` message gains one new optional scalar field:

```
float rating = 7;
```

- **Backward compatible**: existing callers that don't read field 7 are unaffected.
- **Wire format**: proto3 default for `float` is `0.0`; existing serialised messages without field 7 deserialise cleanly with `rating = 0.0`.
- **JSON mapping**: `"rating": 4.3` in `products.json` maps to `Product.Rating` via `jsonpb.Unmarshal`.

---

## Consumer impact

| Consumer | Impact |
|----------|--------|
| `frontend` | Reads `product.Item.Rating` in `product.html` template — **updated** |
| `recommendationservice` | Does not read `Product` fields beyond `Id` — **no change** |
| All other services | Do not consume `Product` message fields — **no change** |

---

## No new RPCs

This story does not add any new RPC methods. The existing `GetProduct` RPC is sufficient.
