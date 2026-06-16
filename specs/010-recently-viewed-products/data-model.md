# Phase 1 Data Model: Recently Viewed Products

No persistent schema is introduced (constraint C-002). The "data model" is the in-memory structure and the view-model passed to the template.

## In-memory store

### `recentlyViewedStore`
Process-lifetime state held on `frontendServer`.

| Field | Type | Description |
|---|---|---|
| `mu` | `sync.Mutex` | Guards concurrent access (multiple requests per session). |
| `bySession` | `map[string][]string` | Session ID → ordered list of product IDs, **most-recent-first**. |

**Invariants**
- IDs in a session list are ordered most-recent-first (index 0 = most recent).
- The list itself is not de-duplicated in this story (dedup is AIP-201).
- Reads return at most 4 entries **after** excluding the current product ID.
- Empty/unknown session → empty list (never nil-deref; never errors).

### Operations (store interface)

| Operation | Signature (conceptual) | Behaviour |
|---|---|---|
| Record a view | `Record(sessionID, productID string)` | Prepend `productID` to the session's list. O(1) amortized. |
| List for display | `List(sessionID, excludeProductID string, max int) []string` | Return up to `max` IDs from the session list, skipping any equal to `excludeProductID`. |

> Note: `List` returns IDs; the handler hydrates them to products via `getProduct`. Keeping the store ID-only avoids coupling it to the catalogue and keeps it trivially testable.

## View-model (handler → template)

The `productHandler` adds one key to the template data map:

| Template key | Type | Description |
|---|---|---|
| `recently_viewed` | `[]*pb.Product` | Hydrated, ordered (most-recent-first), ≤4 products, current excluded, unresolved IDs skipped. Empty/absent → strip not rendered. |

This mirrors the existing `recommendations` key (`[]*pb.Product`) consumed by `recommendations.html`.

## Entity mapping to spec

| Spec entity | Realization |
|---|---|
| **Recently Viewed List** | `bySession[sessionID]` (ordered IDs, ≤4 shown) |
| **Product** | Existing `pb.Product` from `productCatalogService` (Id, Name, Picture, PriceUsd) |
| **Session** | Existing `sessionID(r)` from the `shop_session-id` cookie |
