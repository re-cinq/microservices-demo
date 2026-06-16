# Contract: Wishlist HTTP/UI routes (frontend)

The frontend exposes server-rendered pages and form endpoints via `gorilla/mux` in
`src/frontend/main.go`. This story adds the routes below. Paths are prefixed with the
existing `baseUrl`. Identity is the existing session cookie (`sessionID(r)`); no auth.

## POST `{baseUrl}/wishlist` — Save a product

Mirrors the add-to-cart form contract.

- **Request**: `application/x-www-form-urlencoded`
  - `product_id` (string, required) — the product to save.
- **Behaviour**: Adds `product_id` to the current session's saved collection
  (idempotent — saving an already-saved product changes nothing). Empty/missing
  `product_id` is ignored.
- **Response**: `302` redirect back to the referring product page (or product page for
  `product_id`), so the page re-renders with the saved indicator set (FR-001, FR-002).
- **Errors**: Unknown `product_id` (not in catalogue) → save is rejected gracefully;
  shopper is redirected without a server error.

## GET `{baseUrl}/wishlist` — Saved-products view

- **Request**: none (session cookie identifies the shopper).
- **Behaviour**: Lists the session's saved products, resolving each ID via the existing
  catalogue helper; product IDs no longer in the catalogue are skipped (edge case).
- **Response**: `200` HTML rendering `wishlist.html`:
  - With saved products → each saved product with a link to its product page
    (`{baseUrl}/product/{id}`) (FR-003, FR-005).
  - With none saved → empty-state message + prompt to browse (FR-008,
    Acceptance Scenario 4).

## Header link (all pages)

- `header.html` gains a link to `{baseUrl}/wishlist`, present across the store
  (FR-004). **No count** in this story (count is AIP-195).

## Product page indicator

- `product.html` shows whether the viewed product is currently saved, via a boolean
  passed into the template (from `Store.Contains(sessionID, productID)`) (FR-002).
- The "Save" control is a `POST` form to `{baseUrl}/wishlist` with a hidden
  `product_id`, styled like the existing add-to-cart button.

## Out of scope (separate stories)

- `POST {baseUrl}/wishlist/remove` — removal (AIP-194).
- Saved-count badge in the header (AIP-195).
- "Save" control on product-list cards (AIP-193).
