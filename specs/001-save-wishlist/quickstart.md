# Quickstart: Save for Later / Wishlist

**Feature**: 001-save-wishlist
**Date**: 2026-05-20

## Running the frontend locally

```bash
cd src/frontend
go run .
```

The frontend requires several environment variables for service addresses. For local development without the full cluster, the simplest approach is to run via Skaffold (see repo root `skaffold.yaml`) or Docker Compose if available.

## Testing the wishlist manually

Once the frontend is running:

1. Open any product page, e.g. `http://localhost:8080/product/<id>`
2. Click **Save for later** — you should see the inline confirmation banner appear on the same page.
3. Navigate to `http://localhost:8080/wishlist` — the saved product should appear with its name, image, and price.
4. Navigate to `http://localhost:8080/wishlist` without having saved anything — you should see the empty-state message.
5. Save a product twice — navigate to `/wishlist` and confirm the product appears only once (no duplicate).

### Session-clearing verification (US3)

6. Save one or more products in your current browser session.
7. Open a new incognito / private window (this starts a fresh session with no `shop_session-id` cookie).
8. Navigate to `http://localhost:8080/wishlist` in the incognito window — the list must be **empty**.
9. Confirm that the original window still shows your saved products (sessions are independent).

## Running unit tests

```bash
cd src/frontend
go test ./...
```

Existing tests cover `money/` and `validator/` packages. New wishlist handler logic should be tested in a new `wishlist_test.go` file covering:

- `saveWishlistHandler` — adds product ID, prevents duplicates, rejects empty product_id
- `viewWishlistHandler` — returns empty list when no items saved, returns correct items when saved

## Key files for this feature

| File | Change |
|---|---|
| `src/frontend/main.go` | Add `wishlists sync.Map` field to `frontendServer`; register 2 new routes |
| `src/frontend/handlers.go` | Add `saveWishlistHandler`, `viewWishlistHandler`; update `productHandler` for `?saved=1` |
| `src/frontend/templates/product.html` | Add "Save for later" form button + conditional confirmation banner |
| `src/frontend/templates/wishlist.html` | New template for wishlist view |
| `src/frontend/templates/header.html` | Add wishlist nav icon link |
