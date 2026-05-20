# Plan: Star ratings on the product list (AIP-154)

## Context

### Current state
AIP-153 is marked Done in Jira (PR #29) but those changes are not present on the `attendee/kate-payne` branch. This implementation therefore delivers the complete vertical slice: proto field, data, catalog loader, template component, and product-list rendering.

### Tech stack
- **Product catalogue service**: Go, reads `products.json` via `jsonpb.Unmarshal` into generated proto structs, serves via gRPC
- **Frontend service**: Go, `html/template`, fetches products from catalogue service via gRPC, renders `home.html` (product list) and `product.html` (detail page)
- **Data contract**: `protos/demo.proto` → compiled into `src/frontend/genproto/demo.pb.go` and `src/productcatalogservice/genproto/demo.pb.go`
- **Styling**: `src/frontend/static/styles/styles.css`

### Key constraint
`jsonpb.Unmarshal` (used by the catalog loader) resolves JSON fields via the proto binary descriptor, not Go struct tags. Manually adding a field to the Go struct is not enough — a secondary `encoding/json` pass is needed to populate `Rating` until the protos are regenerated with `genproto.sh`.

---

## Design decisions

### D1 — Add `rating` as a proto field (field number 7)
Keeps the rating on the `Product` message, flows naturally through the existing gRPC path, and requires no new service or RPC.

### D2 — Two-step JSON loading in the catalog service
`loadCatalogFromLocalFile` will: (1) `jsonpb.Unmarshal` as today for all base fields, then (2) a secondary `encoding/json` pass over `products.json` to extract and merge `rating` values by product ID. This avoids regenerating the proto binary descriptor while keeping the change local and reversible.

### D3 — `renderStars` template function
A new `renderStars(rating float32) template.HTML` function registered on the template engine returns ready-to-render star HTML. This keeps all star-rendering logic in one Go function (not scattered across templates), is easily testable, and is the reusable component AIP-155 will also use.

### D4 — Static seed ratings in `products.json`
Ratings are static values baked into `products.json` (1–5 scale, one decimal place). No vote count. No write path. Fully satisfies the spec without a new datastore.

---

## File changes

| File | Change |
|------|--------|
| `protos/demo.proto` | Add `float rating = 7;` to `Product` message |
| `src/frontend/genproto/demo.pb.go` | Add `Rating` field + `GetRating()` accessor to `Product` struct |
| `src/productcatalogservice/genproto/demo.pb.go` | Same as above |
| `src/productcatalogservice/catalog_loader.go` | Add secondary `encoding/json` pass to populate `Rating` on each product |
| `src/productcatalogservice/products.json` | Add `"rating": X.X` to every product entry |
| `src/frontend/handlers.go` | Register `renderStars` template function; add `"renderStars"` to `template.FuncMap` |
| `src/frontend/handlers.go` | Add `renderStars(rating float32) template.HTML` function |
| `src/frontend/templates/home.html` | Call `renderStars .Item.Rating` inside each `.hot-product-card` div |
| `src/frontend/templates/product.html` | Call `renderStars $.product.Item.Rating` on the product detail page (AIP-153 surface) |
| `src/frontend/static/styles/styles.css` | Add `.star-rating` styles |

---

## Implementation sequence

```
1. proto field       — protos/demo.proto
2. genproto structs  — both demo.pb.go files (frontend + productcatalogservice)
3. products.json     — seed rating values
4. catalog_loader    — two-step JSON parse to populate Rating
5. renderStars func  — handlers.go (function + FuncMap registration)
6. home.html         — star display on product list cards
7. product.html      — star display on product detail page
8. styles.css        — star rating CSS
```

---

## Star rendering spec

```
renderStars(0)   → <span class="star-rating no-rating">☆☆☆☆☆ <small>no reviews</small></span>
renderStars(4.5) → <span class="star-rating" aria-label="4.5 out of 5 stars">★★★★½</span>
```

- Full star: `★` (U+2605) for each whole unit ≤ floor(rating)
- Half star: `½` appended if fractional part ≥ 0.5
- Empty star: `☆` (U+2606) for remaining slots up to 5
- Scale: 1–5; values outside that range are clamped
- Zero / absent rating: hollow stars + "no reviews" label
