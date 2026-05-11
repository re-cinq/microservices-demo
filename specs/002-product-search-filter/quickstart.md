# Quickstart: Frontend Product Search Filter

**Feature**: 002-product-search-filter

---

## What Changed

One file: `src/frontend/templates/home.html`

- A `<input type="search">` element added above the product grid.
- A `data-product-name` attribute added to every `.hot-product-card` div (stamped by Go template).
- A hidden `<div id="product-no-results">` added after the last card.
- A self-contained `<script>` block (IIFE) at the end of the product section wires up the filter.

---

## Run Locally (Docker Compose / Skaffold)

```bash
# From repo root
skaffold dev
```

Or, to build and run just the frontend service:

```bash
cd src/frontend
go build ./...
```

The frontend reads templates from the `templates/` directory at startup. No recompile is needed for template-only changes when using Skaffold's file sync.

---

## Manual Test Checklist

1. Open the home page (`/`).
2. Confirm the search box appears above the product grid with placeholder "Search products…".
3. Type `mug` — confirm only the Mug card is visible.
4. Type `WATCH` (uppercase) — confirm only the Watch card is visible (case-insensitive).
5. Type `zzz` — confirm all cards are hidden and "No products match your search" message appears.
6. Clear the input — confirm all cards reappear and the message disappears.
7. Type a space only — confirm all cards remain visible (whitespace-only = no filter).
8. Tab to the search box and type using keyboard only — confirm full operability.
9. Verify zero network requests are triggered while typing (check browser DevTools → Network tab, filter by XHR/Fetch).

---

## Run Go Tests

```bash
cd src/frontend
go test ./...
```

The filter logic is pure DOM JavaScript and is not covered by Go unit tests. Browser-level acceptance tests (if added) would live in an e2e test suite; none exist in this repo yet.
