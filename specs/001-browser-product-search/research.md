# Research: Browser-Side Product Search

**Feature**: `specs/001-browser-product-search`
**Branch**: `attendee/joel-sanmoogan`
**Date**: 2026-05-11

## Decisions

### 1. Filtering mechanism

**Decision**: Pure DOM filtering via a vanilla JavaScript `input` event listener. No framework, no build step.

**Rationale**: The existing frontend has zero JavaScript. Introducing a framework (React, Alpine, HTMX) would add a build dependency and contradict the constraint against new infrastructure. Vanilla JS is the minimal viable approach and consistent with the repo's current posture.

**Alternatives considered**:
- Server-side search endpoint on `frontend` + form submit: rejected — requires a backend change and a page reload, both ruled out by spec constraints.
- HTMX: rejected — new dependency, no existing usage in the repo.
- Alpine.js: rejected — same reason.

---

### 2. Script placement

**Decision**: Inline `<script>` block at the bottom of `home.html`, inside the `{{ define "home" }}` block.

**Rationale**: Keeps the change self-contained to one template file. The script is only needed on the home page; a shared static `.js` file would be wasteful and harder to scope. Go's `html/template` passes the rendered product name text directly into `.hot-product-card-name`; reading `textContent` at runtime avoids any data encoding issues.

**Alternatives considered**:
- New `src/frontend/static/js/search.js` file + `<script src="...">` tag in `header.html`: functional, but adds a file and a global script load on every page. Rejected in favour of localised inline script.

---

### 3. Search box placement

**Decision**: Inject a `<div class="col-12">` search row immediately above the "Hot Products" heading row inside `.hot-products-row`.

**Rationale**: Placing the search box above the heading is the standard e-commerce convention (user expects to search before browsing). It sits within the same Bootstrap grid column so it inherits existing responsive padding.

**Alternatives considered**:
- Inside the header navbar: rejected — the header is shared across all pages; a product-list-only filter does not belong there.
- Below the heading: functional but non-standard; users expect the control above the content it filters.

---

### 4. CSS approach

**Decision**: Add search-box styles to `src/frontend/static/styles/styles.css`.

**Rationale**: All existing custom CSS lives in `styles.css`. Adding a new file would require a `<link>` tag in `header.html` (shared across all pages) or `home.html` (unsupported — `<head>` is rendered by `header.html` before `home.html` content). The simpler path is a small, clearly-labelled section in the existing stylesheet.

---

### 5. Empty-state message

**Decision**: A hidden `<p id="no-results-msg">` paragraph injected once into the product grid; toggled visible by JS when zero cards match.

**Rationale**: Keeps the empty state declaration in HTML (server-rendered, accessible) rather than injecting DOM nodes from JS (fragile). The JS only toggles CSS visibility.

---

### 6. Existing tests

**Decision**: No changes to Go test files. The filter is a pure HTML/CSS/JS addition; it does not alter any Go handler logic or template data pipeline.

**Rationale**: `money_test.go`, `validator_test.go`, `shippingservice_test.go`, `product_catalog_test.go` all test Go logic untouched by this feature. `go test ./...` in any Go service directory must continue to pass.
