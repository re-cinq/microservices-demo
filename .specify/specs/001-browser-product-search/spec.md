# Feature Specification: Browser Product Search

**Feature Branch**: `001-browser-product-search`  
**Created**: 2026-05-11  
**Status**: Draft  
**Input**: Add a search box to the Online Boutique homepage that calls the existing SearchProducts RPC on productcatalogservice, returning matches by product name (case-insensitive substring), rendered inline on the homepage.

## Problem

A customer browsing Online Boutique sees all products on the homepage but has no way to filter or find a specific item by name. With 10+ products on the page, the only option is to scroll and scan visually. There is no search box, no filter, and no way to narrow the catalogue.

## User Scenarios & Testing

### User Story 1 - Search for a product by name (Priority: P1)

A customer lands on the homepage, types a word into a search box, and sees only the products whose names match their query — rendered in the same product card layout as the full catalogue.

**Why this priority**: This is the core interaction. Without it, there is no feature.

**Independent Test**: Type "shirt" into the search box. Only products with "shirt" in the name or description appear. The result cards look identical to the default product cards.

**Acceptance Scenarios**:

1. **Given** the homepage is loaded with all products visible, **When** the customer types "shirt" into the search box and submits (via Enter key or search button), **Then** only products with "shirt" in the name or description are displayed in the product grid.
2. **Given** the customer has typed a query, **When** the results render, **Then** each result uses the same product card layout (image, name, price) as the default homepage listing.
3. **Given** the customer has typed a query, **When** results render, **Then** results appear within 500ms for the existing catalogue size (under 100 items).
4. **Given** the customer has selected a non-default currency (e.g. EUR), **When** they search for a product, **Then** search results display prices in the selected currency, not USD.

---

### User Story 2 - No results found (Priority: P1)

A customer searches for something that doesn't exist in the catalogue. Instead of a blank page or confusing silence, they see a clear message telling them nothing matched.

**Why this priority**: Empty states are where bad search lives. A missing empty state makes the feature feel broken.

**Independent Test**: Type "xyznonexistent" into the search box. A "no results" message appears. No product cards are shown.

**Acceptance Scenarios**:

1. **Given** the homepage is loaded, **When** the customer types a query that matches no products and submits, **Then** a message reading "No products found for '[query]'" (or equivalent) is displayed instead of product cards.
2. **Given** no results are found, **When** the message is displayed, **Then** no product cards, broken layouts, or empty grids are visible.

---

### User Story 3 - Clear search and return to full catalogue (Priority: P2)

After searching, the customer wants to go back to seeing all products. They clear the search box or remove their query and the full catalogue returns.

**Why this priority**: Without a way back, search is a one-way trip. Important for usability but the feature is functional without it.

**Independent Test**: Search for "shirt", see filtered results, clear the search box and submit (or press a clear button). All products reappear.

**Acceptance Scenarios**:

1. **Given** the customer has an active search with filtered results, **When** they clear the search box and submit an empty query, **Then** all products are displayed again (same as initial homepage load).

---

### Edge Cases

- **Empty query submitted**: Submitting an empty search box returns all products (same as no search).
- **Query with only whitespace**: Treated as empty query; all products returned.
- **Very long query (500+ characters)**: The search handles it gracefully — returns no results, does not error or crash.
- **Special characters in query** (`<script>`, `"`, `'`, `&`): The search handles them safely — no XSS, no server errors. Results simply show no matches if none exist.
- **Case insensitivity**: Searching "SHIRT", "Shirt", and "shirt" all return the same results.
- **Single-character query**: Allowed; returns all matching products. No minimum query length enforced.

## Requirements

### Functional Requirements

- **FR-001**: The homepage MUST display a search input field visible above the product grid. The search MUST support submission via Enter key. A visible search button is optional.
- **FR-002**: When a customer submits a non-empty query, the frontend MUST call the existing `SearchProducts` RPC on `productcatalogservice` with the query string. The search MUST use an HTTP GET with a query parameter (e.g. `?q=shirt`) so results are bookmarkable and back-button friendly.
- **FR-003**: The frontend MUST display only the products returned by `SearchProducts`, using the same product card layout as the default listing.
- **FR-004**: When `SearchProducts` returns zero results, the frontend MUST display a "no products found" message instead of product cards.
- **FR-005**: When a customer submits an empty or whitespace-only query, the frontend MUST display all products (equivalent to `ListProducts`).
- **FR-006**: The search MUST be case-insensitive (handled by the existing `SearchProducts` RPC).
- **FR-007**: The search input MUST sanitise user input to prevent XSS in rendered output.

### Constraints (what NOT to do)

- Do NOT add a separate search results page; results render inline on the homepage, replacing the product grid.
- Do NOT change the existing product card layout or styling.
- Do NOT add new services, databases, Elasticsearch, Solr, or any new datastore.
- Do NOT modify the existing `SearchProducts` RPC implementation — it already does case-insensitive substring matching on name and description. Use it as-is.
- Do NOT add new infrastructure config, Helm charts, manifests, or environment variables.
- Do NOT add authentication or user-specific search history.
- Do NOT change the existing product list layout beyond adding the search input.
- Stay inside Go for the frontend and productcatalogservice. Use the existing protobuf/gRPC patterns.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A customer can find a specific product by typing part of its name, with results appearing in under 500ms for catalogues up to 100 items.
- **SC-002**: 100% of acceptance scenarios pass in manual testing on the training URL.
- **SC-003**: The search box is visible without scrolling on a 1280x720 viewport.
- **SC-004**: No JavaScript errors or server errors are introduced by the search feature.

## Assumptions

- The existing `SearchProducts` RPC on `productcatalogservice` works correctly and matches on both name and description using case-insensitive substring. No backend changes are needed.
- The product catalogue remains small enough (under 100 items) that in-memory search with no caching or pagination is acceptable.
- The frontend is server-rendered Go templates — the search will be a form submission or GET request, not a client-side JavaScript filter.
- The CI deploys whatever lands on `attendee/sermed-ali`. No additional deploy configuration is needed.
