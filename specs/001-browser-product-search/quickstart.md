# Quickstart: Browser-Side Product Search

**Feature**: `specs/001-browser-product-search`
**Date**: 2026-05-11

## What was built

A client-side search box on the Online Boutique home page (`/`) that filters the visible product cards by name on every keystroke. No backend changes. No new services or dependencies.

## Files changed

| File | Change |
|---|---|
| `src/frontend/templates/home.html` | Added search `<input>` markup above the product grid + inline `<script>` for filtering |
| `src/frontend/static/styles/styles.css` | Added `.product-search-box` styles and `#no-results-msg` rule |

## Running locally

### Option A — Local Kubernetes cluster (minikube or kind)

```sh
# 1. Start a local cluster (if not already running)
minikube start --cpus=4 --memory 4096 --disk-size 32g

# 2. Build and deploy all services
skaffold run

# 3. Forward the frontend port
kubectl port-forward deployment/frontend 8080:8080

# 4. Open the home page
open http://localhost:8080
```

### Option B — Build the frontend image only (faster iteration)

```sh
# From repo root — rebuild and redeploy only the frontend
skaffold run --module app -b frontend

# Then port-forward as above
kubectl port-forward deployment/frontend 8080:8080
```

## Manual verification checklist

1. **Basic filter** — Type `"sun"` → only the Sunglasses card remains visible.
2. **Clear restores all** — Delete the search text → all 9 product cards reappear.
3. **Case-insensitive** — Type `"SUNGLASSES"` → Sunglasses card is still visible.
4. **Zero results** — Type `"xyznotaproduct"` → all cards hidden, "No products match your search." message appears.
5. **Whitespace only** — Type `"   "` → all cards remain visible (whitespace treated as empty).
6. **Live update** — Type one character at a time; the list updates on every keystroke with no page reload.
7. **Keyboard only** — Press Tab until the search box is focused; type without using the mouse; confirm filtering works.

## Running existing tests

No Go test files were modified. Verify nothing is broken:

```sh
# frontend
cd src/frontend && go test ./...

# productcatalogservice
cd src/productcatalogservice && go test ./...

# shippingservice
cd src/shippingservice && go test ./...
```

All suites must pass with no changes.
