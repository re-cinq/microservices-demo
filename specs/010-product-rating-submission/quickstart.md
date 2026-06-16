# Quickstart: Rate a product and see it on the product page

## What this feature does

On a product page, a shopper picks 1–5 stars and submits. The product's average (shown as stars) and the number of ratings update to include the submission. Ratings are created entirely from shopper input and held in memory in the product catalogue service.

## Touch points (for the implementer)

1. **`protos/demo.proto`** — add `rating` + `num_ratings` to `Product`; add `RateProduct` RPC and `RateProductRequest`. Regenerate `genproto` for both services.
2. **`src/productcatalogservice/product_catalog.go`** — add the mutex-guarded `map[string]*ratingAggregate`; implement `RateProduct` (validate 1–5, update aggregate); populate `rating`/`num_ratings` in `GetProduct` and `ListProducts`.
3. **`src/frontend/rpc.go`** — add a `rateProduct(ctx, id, stars)` client helper.
4. **`src/frontend/main.go`** — register `POST /product/{id}/rate`.
5. **`src/frontend/handlers.go`** — add `rateProductHandler`; the existing `productHandler` already passes the product, whose new fields the template reads.
6. **`src/frontend/templates/product.html`** — render stars + `(num_ratings)`; add a 1–5 star submit form posting to `/product/{id}/rate`.

## Manual verification (maps to acceptance scenarios)

1. Open a product page → submit 4 stars → page shows `★★★★☆ (1)`. *(AC1, AC2)*
2. Submit again with 2 stars → count becomes 2, average reflects both. *(AC1)*
3. Restart `productcatalogservice` → reopen product → ratings reset to "no ratings". *(AC3)*
4. Attempt a submission outside 1–5 (e.g. via crafted POST) → average/count unchanged. *(FR-003 / SC-002)*

## Test

- Unit-test the aggregate + validation in `src/productcatalogservice/product_catalog_test.go` (table-driven: valid 1–5, out-of-range, multiple submissions averaging, unknown product).
- Run: `cd src/productcatalogservice && go test ./...`

## Constraints to honour

In-memory only (no DB/cache), existing two Go services only, no infra/deploy/CI edits, no review text or moderation. See `plan.md` Constitution Check (HC-001..HC-005).
