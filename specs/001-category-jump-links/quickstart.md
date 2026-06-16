# Quickstart Validation Guide: Category Jump Links

**Feature**: AIP-187
**Date**: 2026-06-16

## Prerequisites

- Docker and `skaffold` installed, or the app running locally via `skaffold dev` / `docker-compose`.
- The app is accessible at `http://localhost:8080` (default local port).

## Setup

No database migrations or environment changes required. The feature touches only:
- `src/frontend/handlers.go`
- `src/frontend/templates/home.html`

Run the frontend service:

```bash
# From repo root — hot-reload with skaffold
skaffold dev --port-forward

# Or build and run the frontend only (if running services individually)
cd src/frontend && go build ./... && go test ./...
```

## Validation scenarios

### Scenario 1 — Jump links appear on page load

1. Open `http://localhost:8080` in a browser.
2. **Expected**: A navigation bar with category names (e.g. "Accessories", "Clothing", "Kitchen") is visible near the top of the product section — before any product cards.
3. **Expected**: Product cards below are grouped under matching `<h3>` headings.

### Scenario 2 — Clicking a jump link scrolls to the correct section

1. Open `http://localhost:8080`.
2. Click the "Kitchen" jump link (or any visible category).
3. **Expected**: The viewport scrolls to the `<h3>Kitchen</h3>` heading without a full page reload.
4. Verify by checking the URL fragment — it should update to `#kitchen`.

### Scenario 3 — Product order is preserved within categories

1. Open `http://localhost:8080`.
2. Compare the order of products under each category heading against the original product list (reference: `productcatalogservice/products.json` — products appear in JSON order).
3. **Expected**: Within each category, products appear in the same relative order as they appear in `products.json`.

### Scenario 4 — No backend or infrastructure changes

```bash
# Confirm no new services were added
git diff --name-only HEAD | grep -v 'src/frontend/'
```

**Expected**: No files outside `src/frontend/` appear in the diff.

## Go unit tests

```bash
cd src/frontend && go test ./...
```

**Expected**: All existing tests pass. If new grouping logic is added as a helper function, it should have its own unit tests that also pass.

## Acceptance checklist

- [ ] Category jump links visible on page load (Scenario 1)
- [ ] Clicking a jump link scrolls to the matching section (Scenario 2)
- [ ] Product order within categories matches `products.json` order (Scenario 3)
- [ ] No files outside `src/frontend/` changed (Scenario 4)
- [ ] `go test ./...` passes in `src/frontend/`
