# Research: Browser-Side Product Search

**Phase**: 0 — Research  
**Feature**: Browser-Side Product Search  
**Date**: 2026-05-08

---

## Decision 1: Where does filtering run?

**Decision**: Entirely in the browser, using the product name text already rendered in the DOM.

**Rationale**: The spec hard-constrains "no backend change" and "filter in memory". The Go frontend template already renders every product card (`.hot-product-card`) with its name inside `.hot-product-card-name`. A JavaScript `input` event listener can read those text nodes and toggle card visibility without any network round-trip.

**Alternatives considered**:
- Server-side query parameter filter (`?q=sunglasses`) — rejected: requires backend change and page reload.
- Fetch a `/search?q=` endpoint — rejected: requires a new backend route, violates constraints.

---

## Decision 2: How is JavaScript delivered?

**Decision**: Inline `<script>` block at the bottom of `home.html`, inside the `{{ define "home" }}` block.

**Rationale**: The frontend has no existing JS bundler or build step. There is no `package.json`, no Webpack/Vite config, no `node_modules`. Adding a separate `.js` file under `static/` is technically possible but requires a `<script src>` tag anyway — there is no advantage over inline script for a feature this small (< 30 lines). Keeping it inline in the template keeps the change to a single file.

**Alternatives considered**:
- New `static/scripts/search.js` file — rejected: same outcome, more files, no build-pipeline benefit.
- Alpine.js / htmx — rejected: introduces an external dependency, violates the "no new infrastructure config" spirit and adds unnecessary complexity.

---

## Decision 3: Which DOM attributes drive filtering?

**Decision**: The text content of `.hot-product-card-name` inside each `.hot-product-card` card element.

**Rationale**: The template renders product names in `<div class="hot-product-card-name">{{ .Item.Name }}</div>`. This is the canonical, human-readable name. Reading `textContent` (or `innerText`) of that element and comparing case-insensitively against the query is reliable and requires zero data-attribute additions.

**Alternatives considered**:
- `data-product-name` attributes on each card — rejected: requires a template change AND a JS change; no benefit over direct text comparison.
- Hidden JSON blob embedded in the page — rejected: over-engineered; DOM traversal is sufficient.

---

## Decision 4: What triggers the filter?

**Decision**: The `input` event on the search `<input>` element (fires on every keystroke, paste, and clear).

**Rationale**: `input` is the correct event for real-time filtering; it fires immediately on every character change including cut/paste/autocomplete. `keyup` misses IME composition and paste events. `change` only fires on blur.

**Alternatives considered**:
- `keyup` — rejected: misses paste and clear-via-button events.
- Debounced `input` — rejected: catalogue is small (< 30 products in products.json); debounce adds complexity for no perceptible benefit.

---

## Decision 5: CSS styling approach

**Decision**: Add a small CSS rule block to the existing `src/frontend/static/styles/styles.css` file for the search box.

**Rationale**: The search input needs minimal styling to fit the existing Bootstrap-based layout. A few lines in the existing stylesheet (width, margin, border-radius) is cleaner than inline `style=""` attributes.

**Alternatives considered**:
- Inline styles on the `<input>` — rejected: harder to maintain, not consistent with how existing elements are styled.
- New `search.css` file — rejected: unnecessary for 5–10 CSS lines; adds a `<link>` tag overhead.

---

## Findings: Existing code patterns

| Item | Finding |
|------|---------|
| Template engine | Go `html/template`, all templates in `src/frontend/templates/*.html` |
| Product card selector | `.hot-product-card` (Bootstrap col div) |
| Name element selector | `.hot-product-card-name` (div inside each card) |
| Existing JS | None in any template; no JS bundler |
| Existing CSS | Bootstrap (CDN, via header.html) + `static/styles/styles.css` |
| Product count | ≤ 30 items in `productcatalogservice/products.json` |
| Files to change | `src/frontend/templates/home.html`, `src/frontend/static/styles/styles.css` |
| Files NOT to change | Any Go source, any proto, any infra config |
