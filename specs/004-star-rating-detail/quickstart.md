# Quickstart: Star Rating on Product Detail Page

**Feature**: 004-star-rating-detail

---

## What this feature changes

1. `src/productcatalogservice/genproto/demo.proto` — adds `float rating = 7` to `Product`
2. `src/productcatalogservice/genproto/demo.pb.go` — adds `Rating float32` field (updated in sync with proto)
3. `src/productcatalogservice/products.json` — each product entry gains a `"rating"` value (0.0–5.0)
4. `src/frontend/templates/product.html` — renders the rating beneath the product name

No new files, no new services, no infrastructure changes.

---

## Run locally (existing pipeline)

```sh
# From repo root — same command as today, no changes needed
skaffold dev
```

The existing pipeline picks up all changes automatically.

---

## Verify the feature

1. Open the app in a browser.
2. Click any product.
3. Confirm a rating (e.g., "★ 4.3 / 5") is visible beneath the product name and above the description.
4. Confirm the rest of the page (price, image, add-to-cart, recommendations) is unchanged.

---

## Run tests

```sh
# productcatalogservice unit tests
cd src/productcatalogservice && go test ./...

# frontend (if tests exist)
cd src/frontend && go test ./...
```
