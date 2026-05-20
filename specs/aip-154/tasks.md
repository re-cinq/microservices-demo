# Tasks: Star ratings on the product list (AIP-154)

## Phase 1 — Data contract

- [x] [T1] Add `float rating = 7;` to the `Product` message in `protos/demo.proto`
- [x] [T2] Add `Rating float32` field and `GetRating()` accessor to `Product` struct in `src/frontend/genproto/demo.pb.go`
- [x] [T3] Add `Rating float32` field and `GetRating()` accessor to `Product` struct in `src/productcatalogservice/genproto/demo.pb.go`

## Phase 2 — Data

- [x] [T4] Add `"rating": X.X` seed values (1–5 scale) to every product entry in `src/productcatalogservice/products.json`

## Phase 3 — Catalog service

- [x] [T5] Extend `loadCatalogFromLocalFile` in `src/productcatalogservice/catalog_loader.go` with a secondary `encoding/json` pass to extract and merge `rating` values by product ID into the loaded `pb.Product` entries

## Phase 4 — Frontend rendering

- [x] [T6] Add `renderStars(rating float32) template.HTML` function to `src/frontend/handlers.go` — returns filled/half/empty star HTML; zero rating renders hollow stars + "no reviews" label
- [x] [T7] Register `renderStars` in the `template.FuncMap` in `src/frontend/handlers.go`
- [x] [T8] Add `.star-rating` CSS rules to `src/frontend/static/styles/styles.css`

## Phase 5 — Templates

- [x] [T9] Add `{{ renderStars $.product.Item.Rating }}` to the product detail page in `src/frontend/templates/product.html` (AIP-153 surface — prerequisite for consistency check)
- [x] [T10] Add `{{ renderStars .Item.Rating }}` to each `.hot-product-card` in `src/frontend/templates/home.html` (AIP-154 primary deliverable)
