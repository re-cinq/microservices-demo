# Specification: Product Search

**Feature directory**: `specs/001-product-search`
**Branch**: `001-product-search` (feature branch off `attendee/jost-werdenhoff`; merges back via PR — see C-7)
**Created**: 2026-06-02 · **Rewritten**: 2026-06-02 (lean SpecKit format)
**Input**: "Add a product search feature to Online Boutique. Frontend-only filter. A search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

## Overview

Shoppers landing on the Online Boutique home page see the full product grid but have no way to narrow it down — to find a specific item they must visually scan every product. This feature adds a search box to the product list (home) page that filters the products **already loaded on the page** by **name**, **in the browser**, with no page reload and no backend change. Matching is case-insensitive and substring-based. Clearing the box restores the full grid.

## Functional Requirements

Every requirement below is observable and testable.

- **FR-001**: The product list (home) page MUST display a search input visible to the shopper, above or alongside the product grid.
- **FR-002**: As the shopper changes the search text, the system MUST update the visible products to show only those whose **name** matches the entered text, without a full page reload.
- **FR-003**: Name matching MUST be case-insensitive.
- **FR-004**: Name matching MUST treat the query as a substring (a product matches if its name contains the query text anywhere).
- **FR-005**: When the search text is empty (or only whitespace), the system MUST show all products originally loaded on the page.
- **FR-006**: When the search text matches no product, the system MUST show no products and MUST display a clear "no results" indication.
- **FR-007**: Clearing the search text MUST restore the full originally loaded product set without a page reload and without re-fetching from the server.
- **FR-008**: The system MUST filter only the products already loaded on the page; it MUST NOT issue additional server requests to perform the search.
- **FR-009**: The feature MUST NOT change which products are loaded, their order when unfiltered, or any other existing home-page behaviour.
- **FR-014**: The query MUST be treated as literal text: special characters (e.g. `& < > " % /`, punctuation, emoji) are matched literally against product names and MUST NOT be interpreted as patterns, markup, or executable content, and MUST NOT cause an error.
- **FR-015**: Matching products MUST remain in the same relative order as in the unfiltered grid; filtering only removes non-matches and MUST NOT re-order results.
- **FR-016**: Arbitrarily long queries MUST be handled gracefully — no error and no noticeable lag — yielding the correct matching set (typically the no-results state).

### Accessibility

- **FR-010**: The search input MUST have an accessible name, such as a visible label or equivalent accessible labelling.
- **FR-011**: The no-results indication MUST be text content visible to the shopper and available to assistive technologies.
- **FR-012**: Filtering MUST NOT trap focus or move focus away from the search input while the shopper is typing.
- **FR-013**: Product cards that are hidden by filtering MUST NOT remain reachable by keyboard navigation.

## User Scenarios

Written as Given / When / Then.

1. **Given** the home page is loaded with all catalogue products visible, **When** the shopper types "watch", **Then** only products whose name contains "watch" (case-insensitive) remain visible and all others are hidden.
2. **Given** the grid is filtered, **When** the shopper clears the search box, **Then** the full set of originally loaded products becomes visible again, in original order, with no page reload.
3. **Given** the home page is loaded, **When** the shopper types text matching no product name (e.g. "zzzzz"), **Then** no products are shown and a clear "no results" indication is displayed.
4. **Given** the home page is loaded, **When** the shopper types "WATCH" in uppercase, **Then** the same products match as for "watch".
5. **Given** the home page is loaded, **When** the shopper types a partial term (e.g. "sun"), **Then** every product whose name contains that substring (e.g. "Sunglasses") remains visible.
6. **Given** a keyboard-only shopper has filtered the grid, **When** they Tab through the page, **Then** focus only reaches visible (matching) product cards and never a hidden one (FR-013), and focus is not pulled out of the search input while typing (FR-012).
7. **Given** the home page is loaded, **When** the shopper types special characters (e.g. `&`, `<`, `"`, `%`, an emoji), **Then** the characters are treated as literal text to match against names — the page does not error, does not execute them, and simply shows the products whose name contains that literal text (usually none).
8. **Given** the home page is loaded, **When** the shopper types a very long query (e.g. 500+ characters), **Then** the page filters without error or noticeable lag and shows the matching set (usually the no-results state).
9. **Given** the grid is filtered to several matches, **When** the shopper compares the order of the visible cards, **Then** they appear in the same relative order as in the unfiltered grid (filtering removes non-matches but never re-orders).

**Edge cases**: whitespace-only query behaves like empty; repeated type/clear cycles work without reload; special characters are matched literally (never interpreted as patterns or markup); very long queries are handled gracefully; if in-browser filtering cannot run, the page still renders the full grid exactly as today (graceful degradation).

## Success Criteria

Measurable and technology-agnostic.

- **SC-001**: A shopper can narrow the grid to a specific product by typing its name and sees the filtered result in under 1 second, with no page reload.
- **SC-002**: For any query, the products shown are exactly the originally loaded products whose name contains the query (case-insensitive) — no false positives, no missed matches, across 100% of catalogue products.
- **SC-003**: Clearing the search restores 100% of originally loaded products in original order, with no page reload.
- **SC-004**: A no-match query yields a clearly understandable "no results" state in 100% of no-match cases (never a blank/broken page).
- **SC-005**: All existing home-page behaviour (loading, ordering, navigation, cart, currency) is unchanged — zero shopper-observable regressions.
- **SC-006**: The feature performs no additional server/network requests to filter.
- **SC-007**: The search input is reachable and operable by keyboard and exposes an accessible name; the no-results state is announced to assistive technology (FR-010, FR-011, FR-012, FR-013 verifiable via screen reader + keyboard-only navigation).

## Constraints (what this feature must NOT do)

Hard boundaries for implementation — do not deviate.

- **C-1**: MUST NOT add any new service; use only services already in the repo.
- **C-2**: MUST NOT introduce Elasticsearch, Solr, a vector database, or any new datastore/search engine.
- **C-3**: MUST use the existing in-memory product catalogue (from `productcatalogservice/products.json`) as the product source; filter over products already loaded.
- **C-4**: MUST match the language of any service edited — frontend is Go, `productcatalogservice` is Go. If `productcatalogservice` is extended, follow existing protobuf/gRPC patterns. (This feature is specified as frontend-only.)
- **C-5**: MUST NOT add new infrastructure config, Helm charts, Kubernetes manifests, or environment variables.
- **C-6**: MUST NOT touch the build pipeline / CI configuration.
- **C-7**: MUST stay within this repo. Work happens on the `001-product-search` feature branch (off `attendee/jost-werdenhoff`) and merges back to `attendee/jost-werdenhoff` via Pull Request; CI deploys `attendee/jost-werdenhoff`. Do not work in other repos or branch from elsewhere.
- **C-8**: Search is over product **name** only — not description, categories, or price — for this version.
- **C-9**: No backend change is required or expected; filtering happens in the browser over the already-rendered product list.
- **C-10**: MUST NOT introduce a separate search page or search results page. Search filters the existing product grid **in place** on the home page.
- **C-11**: MUST NOT restructure or redesign the home page layout. The only additive UI is the search control and the no-results message; the existing "Hot Products" grid, header, footer, cart, and currency selector keep their current layout and styling.

## Assumptions

- The home page already loads the full set of products the shopper expects to search across; "search" means filtering that loaded set, not server-side paging.
- Name-only matching is sufficient for this version (C-8).
- The demo catalogue is small enough that in-browser substring filtering is effectively instant.
- The product name as rendered is the authoritative text to match (case-insensitive, whitespace-trimmed).
- Query state is per-page-session only; no persistence across loads.
- Styling and accessibility follow existing home-page conventions; no new design system.
