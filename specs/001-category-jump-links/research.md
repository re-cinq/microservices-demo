# Research: Category Jump Links

**Feature**: AIP-187 — Category jump links on product list page
**Date**: 2026-06-16

## Decision Log

### D-001: Where category data comes from

**Decision**: Use `categories` field already present on each product in `productcatalogservice/products.json`. No new data source needed.

**Rationale**: Every product proto already exposes `Categories []string`. The frontend receives these via the existing `ListProducts` gRPC call. Grouping happens in the `homeHandler` before template rendering — pure in-memory, zero new infrastructure.

**Alternatives considered**: Adding a dedicated `/categories` endpoint — rejected (violates epic constraint: no new services or datastores).

---

### D-002: Which category to use when a product has multiple

**Decision**: Use the first element of `Categories` as the primary display category. Products are assigned to exactly one rendered group.

**Rationale**: All products in `products.json` have 1–2 categories; none span more than two. Using the first keeps grouping deterministic and avoids a product appearing in multiple sections. Example: Tank Top has `["clothing", "tops"]` → grouped under "clothing".

**Alternatives considered**: Showing a product in every category it belongs to — rejected (creates duplicate rows, confusing UX, out of scope for this story).

---

### D-003: Category display order

**Decision**: Preserve the order categories first appear when iterating `products.json` top-to-bottom.

**Rationale**: Spec requires "same order as today." Since products are currently rendered in JSON order, the first occurrence of a category name in that order defines the category's position in the grouped list.

**Alternatives considered**: Alphabetical sort — rejected (changes existing order, violates spec).

---

### D-004: Scroll mechanism

**Decision**: Use native HTML anchor links (`<a href="#category-slug">`) and `id` attributes on section headings. No JavaScript required beyond what browsers provide natively.

**Rationale**: The epic constraint prohibits infrastructure changes; adding a JS scroll library would require a build step or CDN dependency. Native anchor scroll works in all modern browsers with zero dependencies.

**Alternatives considered**: `window.scrollIntoView()` via JS — unnecessary; native anchors satisfy the acceptance criterion.

---

### D-005: Files changed

**Decision**: Two files change:
1. `src/frontend/handlers.go` — add grouping logic in `homeHandler` to build ordered category-product map before template data injection.
2. `src/frontend/templates/home.html` — replace the flat product loop with a category-section loop that emits headings, anchor `id`s, and jump links.

**Rationale**: Spec and epic both state "frontend product list template only." The handler-side grouping is the minimum server-side work needed to pass structured data to the template without logic in the template itself (Go templates are logic-light by convention).

**Alternatives considered**: Doing grouping inside the Go template — rejected (Go templates lack map-insertion; impossible without a custom template function, which would touch more code).
