# Data Model: Save for Later / Wishlist

**Feature**: 001-save-wishlist
**Date**: 2026-05-20

## Entities

### WishlistStore (server-side, in-memory)

Lives as a field on `frontendServer`. Not persisted. Cleared implicitly when the process restarts or the session ID is no longer used.

| Field | Type | Description |
|---|---|---|
| `wishlists` | `sync.Map` | Concurrent map: session ID → `[]string` (ordered slice of product IDs) |

**Lifecycle**: Entries are created on first save for a given session ID. Entries are never explicitly deleted (session IDs expire via cookie max-age; orphaned server-side entries are harmless at demo scale).

**Invariants**:
- A product ID appears at most once per session's slice (enforced at save time).
- Product IDs are the canonical identifiers from the product catalogue (`productcatalogservice`).

---

### WishlistItem (view model, not persisted)

Constructed at render time by fetching product details from `productcatalogservice`.

| Field | Type | Source | Description |
|---|---|---|---|
| `ProductID` | `string` | wishlist store | Identifier used to fetch full product |
| `Name` | `string` | productcatalogservice | Display name of the product |
| `Picture` | `string` | productcatalogservice | Relative image path |
| `Price` | `pb.Money` | productcatalogservice + currencyservice | Price in the shopper's selected currency |

**Note**: `Price` is resolved at display time (when `GET /wishlist` is rendered), not at save time, so it always reflects the current catalogue state.

---

## State Transitions

```
[Product detail page]
        │
        │ POST /wishlist/save  (product_id=<id>)
        ▼
[saveWishlistHandler]
        │
        ├─ product already in wishlist? → skip append
        │
        └─ append product_id to session slice
        │
        └─ redirect → GET /product/<id>?saved=1
                           │
                           ▼
               [product.html shows confirmation banner]

[Header nav: wishlist icon]
        │
        │ GET /wishlist
        ▼
[viewWishlistHandler]
        │
        ├─ fetch product details for each saved ID
        │
        └─ render wishlist.html
               │
               ├─ items present → product cards (name, image, price)
               └─ no items → empty-state message
```

---

## Session Scoping

The session ID is set by the `ensureSessionID` middleware (in `main.go`) and stored in a browser cookie named `shop_session-id`. When a shopper opens a new browser session (no cookie), a new session ID is generated and the associated wishlist is empty.
