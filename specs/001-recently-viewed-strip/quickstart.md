# Quickstart: Recently Viewed Products Strip

**Feature**: `specs/001-recently-viewed-strip`
**Date**: 2026-05-20

---

## Prerequisites

- Docker Desktop running
- `skaffold` installed (`brew install skaffold` or see [skaffold.dev](https://skaffold.dev))
- `kubectl` configured for a local cluster (minikube / Docker Desktop Kubernetes)

---

## Run Locally

```bash
# From repo root — starts all services with hot-reload for the frontend
skaffold dev
```

Once running, the frontend is available at `http://localhost:80` (or the port shown in skaffold output).

---

## Manual Test Flow

1. Open `http://localhost/` in a browser
2. Click any product to open its detail page — this writes the `shop_recently_viewed` cookie
3. Navigate back to the home page or click another product
4. Open a second product detail page
5. Scroll down — the "Recently viewed" strip should appear showing the first product

**Cookie inspection** (Chrome DevTools → Application → Cookies):
- `shop_recently_viewed` should contain a pipe-separated list of product IDs in most-recent-first order

---

## Run Frontend Unit Tests

```bash
cd src/frontend
go test ./...
```

Expected output includes tests for cookie read/write logic in `handlers_test.go`.

---

## Run productcatalogservice Tests

```bash
cd src/productcatalogservice
go test ./...
```

No changes to this service — tests should pass unchanged.

---

## Verify Edge Cases Manually

| Scenario | How to test |
|----------|-------------|
| Single product viewed | View one product; strip should not appear (nothing to return to) |
| Strip at capacity | View 6+ distinct products; strip should show only the 5 most recent |
| Duplicate prevention | View product A, then B, then A again; A should appear once, at front |
| Currency change | View products, then change currency via the dropdown; strip prices should update |
| Stale cookie ID | Manually edit cookie to include a non-existent product ID; page should load normally with that ID omitted from the strip |

---

## Key Files Changed

| File | Change |
|------|--------|
| `src/frontend/handlers.go` | `productHandler`: write/update cookie; new `recentlyViewedProducts()` helper |
| `src/frontend/templates/recently_viewed.html` | New partial template for the strip |
| `src/frontend/templates/product.html` | Include `recently_viewed` partial |
| `src/frontend/handlers_test.go` | New table-driven tests for cookie logic |
