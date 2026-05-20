# HTTP Endpoint Contracts: Save for Later / Wishlist

**Feature**: 001-save-wishlist
**Date**: 2026-05-20

These are the two new routes added to the frontend service. All existing routes are unchanged.

---

## POST /wishlist/save

**Purpose**: Add a product to the current session's wishlist.

**Method**: `POST`
**Path**: `/wishlist/save`
**Handler**: `saveWishlistHandler`

### Request

| Element | Details |
|---|---|
| Content-Type | `application/x-www-form-urlencoded` |
| Body field | `product_id` — the product's catalogue ID (required, non-empty string) |
| Session | Session ID read from `shop_session-id` cookie (set by middleware) |

### Behaviour

1. Read `product_id` from form body. If empty → `400 Bad Request`.
2. Load existing wishlist slice for the current session ID from `wishlists` (`sync.Map`). If no entry exists, initialise an empty slice.
3. If `product_id` is not already in the slice, append it.
4. Store the updated slice back into `wishlists`.
5. Redirect to `GET /product/{product_id}?saved=1`.

### Responses

| Status | Condition |
|---|---|
| `303 See Other` | Success (with or without duplicate — idempotent) |
| `400 Bad Request` | `product_id` missing or empty |

---

## GET /wishlist

**Purpose**: Display all products saved in the current session's wishlist.

**Method**: `GET`
**Path**: `/wishlist`
**Handler**: `viewWishlistHandler`

### Request

| Element | Details |
|---|---|
| Session | Session ID read from `shop_session-id` cookie (set by middleware) |
| Query params | None |

### Behaviour

1. Read saved product ID slice for the current session from `wishlists`.
2. For each product ID, fetch product details from `productcatalogservice` (name, picture, price).
3. Convert prices to the shopper's selected currency via `currencyservice`.
4. Render `wishlist.html` template with the resolved product list.
5. If the slice is empty (or no entry exists), render `wishlist.html` with an empty list — the template shows an empty-state message.

### Template data shape

```go
map[string]interface{}{
    "items":       []WishlistItem,  // resolved product view models (may be empty)
    "user_currency": string,
    // standard fields injected by existing helpers:
    // "baseUrl", "currencies", "cart_size", "is_cymbal_brand", etc.
}
```

### Responses

| Status | Condition |
|---|---|
| `200 OK` | Always (empty list and populated list both return 200) |
| `500 Internal Server Error` | Product catalogue or currency service unreachable |

---

## Modified: GET /product/{id}

**Change**: The existing `productHandler` must read the `saved` query parameter and pass it to the template so `product.html` can render the inline confirmation banner.

| Query param | Value | Effect |
|---|---|---|
| `saved` | `1` | Template renders inline "Saved to your wishlist" confirmation |
| (absent) | — | No change to existing rendering |
