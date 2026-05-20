# Service Contract: `Product.rating`

**Scope**: change to the `Product` message in `protos/demo.proto`. Consumers: `productcatalogservice` (producer), `frontend` (consumer in this story; `checkoutservice`, `recommendationservice`, etc. are unaffected because they ignore the new field).

## Change

Add a single optional field to `message Product`:

```proto
message Product {
    string id = 1;
    string name = 2;
    string description = 3;
    string picture = 4;
    Money price_usd = 5;

    // Categories such as "clothing" or "kitchen" that can be used to look up
    // other related products.
    repeated string categories = 6;

    // 0.0 means "no rating"; otherwise 0.5..5.0 in 0.5 steps.
    float rating = 7;
}
```

### Field semantics

| Aspect | Value |
|---|---|
| Proto number | `7` (next unused) |
| Proto type | `float` |
| Default (proto3) | `0.0` |
| Presence | implicit (no `optional` keyword); `0.0` is the sentinel for "no rating" per data-model.md |
| Legal values at runtime | `0.0`, `0.5`, `1.0`, `1.5`, `2.0`, `2.5`, `3.0`, `3.5`, `4.0`, `4.5`, `5.0` |
| Producer | `productcatalogservice` populates `rating` from `products.json` |
| Consumers | `frontend` reads on the product-list render (this story) and on the product page (next story, AIP-162) |

### RPCs affected

- `ListProducts` — return value now carries `rating` for every product.
- `GetProduct` — return value now carries `rating` for the requested product.
- `SearchProducts` — return values now carry `rating` for every matched product.

No RPC is added, removed, renamed, or has its signature changed. No new RPC for setting or updating ratings (research D3).

## Backward compatibility

- Adding a new field with a fresh number is a non-breaking proto3 change. Services compiled against the old proto continue to deserialise messages from a producer that knows the new field — they simply discard it.
- Producers compiled against the new proto, returning data to a consumer compiled against the old proto: the field is silently ignored by the old consumer.
- Therefore the proto change can be rolled out without coordinating producer/consumer versions. In practice we will rebuild both `productcatalogservice` and `frontend` together because they ship from the same cohort branch.

## Versioning

The repo does not version `demo.proto` formally. No version bump is needed. The Jira story key `AIP-161` is the upstream record of this change.

## Regeneration

The proto change requires regeneration of the per-service `genproto/` stubs:

```bash
# from repo root
( cd src/productcatalogservice && ./genproto.sh )
( cd src/frontend            && ./genproto.sh )
```

Both regenerated `genproto/*.pb.go` files MUST be committed alongside the proto change. CI configuration is **not** touched (research D10, spec constraint C-004).

## Out of scope for this contract

- Mutating ratings via an RPC. There is no `SetRating` / `UpdateRating` RPC.
- Streaming or pub/sub of rating changes.
- A separate `Rating` or `Review` message. Rating is a property of `Product`, not its own resource.
