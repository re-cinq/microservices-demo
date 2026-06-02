# Research: Browser-Side Product Search

**Feature**: `specs/001-product-search`
**Date**: 2026-06-02

---

## Decision 1: Where does filtering happen?

**Decision**: In the browser, against DOM elements already rendered in the page.

**Rationale**: The spec (C-001, C-004) prohibits backend changes and additional network requests. The product catalogue is already fully rendered into the HTML by the Go template at page load time. No server round-trip is needed.

**Alternatives considered**:
- *Server-side via existing `SearchProducts` RPC*: Already implemented in `productcatalogservice` (case-insensitive substring match on name and description). Rejected because it requires a page reload or AJAX call, violating C-004, and would require backend route changes.
- *Fetching products as JSON and re-rendering*: More complexity, same constraint violation. Rejected.

---

## Decision 2: JavaScript approach — vanilla JS vs. framework

**Decision**: Vanilla JavaScript, no external library.

**Rationale**: The frontend has no existing JS files and no build step. The filter logic is approximately 30–40 lines. Adding a library (e.g., Fuse.js, Alpine.js) would introduce a new external dependency and potentially a new CDN tag or a build pipeline, both of which the spec prohibits (C-003). Vanilla `addEventListener` + DOM traversal is sufficient for the performance target (< 300 ms for 500 items is trivially achievable with a simple loop).

**Alternatives considered**:
- *Fuse.js* (fuzzy search library): Richer matching but overkill for name-only substring matching, and adds an external dependency.
- *Alpine.js*: Reactive, declarative, but requires a CDN addition and a mental-model shift for future maintainers unfamiliar with it.

---

## Decision 3: How to expose product names to the JS

**Decision**: Add a `data-name="{{ .Item.Name }}"` attribute to each `.hot-product-card` element in the Go template.

**Rationale**: The product name is already available in the Go template as `.Item.Name`. A `data-*` attribute is the standard, accessible, layout-independent way to pass server-rendered data to client-side JS. It avoids parsing text content from inner elements (fragile) or making an additional data fetch.

**Alternatives considered**:
- *Read from `.hot-product-card-name` text content*: Works but is fragile — any whitespace, formatting, or HTML entity change in the name display breaks the filter. Rejected.
- *Embed a JSON array of products in a `<script>` block*: Works but is more invasive and creates a second representation of the product list in the page. Rejected.

---

## Decision 4: File placement for JavaScript

**Decision**: New file at `src/frontend/static/scripts/product-search.js`.

**Rationale**: The existing static file server in `main.go` already serves everything under `./static/` via `http.FileServer`. No routing change is needed. A `scripts/` subdirectory mirrors the existing `styles/` convention. Keeping JS in a separate file rather than inlining it in `home.html` keeps the template readable and the JS independently editable.

**Alternatives considered**:
- *Inline `<script>` in `home.html`*: Simpler (no new file), but mixes markup and behaviour in the template, and makes the logic harder to find and test.

---

## Decision 5: DOM mutation strategy — hide/show vs. remove/re-insert

**Decision**: Toggle `display: none` / `display: ''` on each `.hot-product-card` element directly.

**Rationale**: The product elements are already in the DOM. Toggling `display` is the lowest-cost operation: no DOM node creation, no innerHTML, no re-layout of unaffected nodes. For up to 500 items this is well within the 300 ms budget. The Bootstrap grid reflows automatically when columns are hidden.

**Alternatives considered**:
- *Clone and re-insert*: Unnecessary overhead.
- *CSS class toggle*: Cleaner but requires adding a CSS rule (`.hidden { display: none }`), adding a file change to the stylesheet. `style.display` is self-contained in the JS.

---

## Resolved unknowns

| Question | Resolution |
|---|---|
| Does the frontend have any existing JS infrastructure? | No — only Bootstrap CDN and one inline `onchange`. Vanilla JS is the right fit. |
| Is Bootstrap's grid compatible with toggling `display` on `.col-md-4`? | Yes — Bootstrap columns are flex items; hiding individual columns does not break the row. |
| Does the static file server need changes? | No — `main.go` already serves `./static/` recursively. |
| Does `SearchProducts` RPC need to be called? | No — client-side only, per spec constraints. |
