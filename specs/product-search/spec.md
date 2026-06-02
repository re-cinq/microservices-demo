# Feature Specification: Product Search (client-side name filter)

**Feature Branch**: `attendee/andrea-backstrom`

**Created**: 2026-06-02

**Status**: Draft

**Input**: User description: "Add a product search feature to Online Boutique. A search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

---

## Problem (in the user's words)

> "I'm on the Online Boutique home page looking at the list of products. There's no way to quickly find a product by name — I have to scan the whole grid with my eyes. I want a search box on the product list page that filters the products that are already on the page, by name, right in the browser, without a page reload and without any backend change."

The product list page ([`src/frontend/templates/home.html`](../../src/frontend/templates/home.html)) renders every product in the catalogue as a card (`{{ range $.products }}`, each showing `.Item.Name`). As the catalogue grows, finding a specific product by eye is slow. There is currently no on-page way to narrow the grid.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter the product grid by name as I type (Priority: P1)

A shopper on the product list page types text into a search box. The grid immediately narrows to only the products whose name matches what they typed. Clearing the box restores the full grid. This is the whole feature — implementing just this story delivers a usable MVP.

**Why this priority**: This is the core value of the feature. Without it there is nothing; with it the feature is complete and demonstrable.

**Independent Test**: Load the home page, type a known product name fragment into the search box, and confirm only matching cards remain visible and all others are hidden — with no network request and no page reload.

**Acceptance Scenarios**:

1. **Given** the home page is loaded with all products visible, **When** the shopper types "sunglasses" into the search box, **Then** only product cards whose name contains "sunglasses" remain visible and all non-matching cards are hidden.
2. **Given** the shopper has typed a query that narrowed the grid, **When** the shopper clears the search box (empties it), **Then** every product card is shown again, in its original order.
3. **Given** the home page is loaded, **When** the shopper types "SUN" (uppercase), **Then** products named "Sunglasses" still match (matching is case-insensitive).
4. **Given** the home page is loaded, **When** the shopper types " sun " (with surrounding spaces), **Then** matching behaves as if the query were "sun" (leading/trailing whitespace is ignored).
5. **Given** the shopper is typing a query, **When** each character is entered, **Then** the visible set of cards updates without any page reload and without any request to the backend.

---

### User Story 2 - Tell me clearly when nothing matches (Priority: P2)

When a shopper's query matches no products, the page shows a clear "no results" message instead of an empty, confusing grid.

**Why this priority**: Prevents the dead-end experience of an empty page. Valuable, but the feature is still usable without it.

**Independent Test**: Type a string that matches no product name and confirm a visible "no products match" message appears and disappears when the query is cleared or changed to something that matches.

**Acceptance Scenarios**:

1. **Given** the home page is loaded, **When** the shopper types a query that matches no product name (e.g. "zzzzz"), **Then** no product cards are shown and a visible "no products match your search" message is displayed.
2. **Given** the "no products match" message is showing, **When** the shopper clears the search box, **Then** the message disappears and all product cards are shown again.

---

### Edge Cases

- **Empty / whitespace-only query**: treated as "no filter" — all products shown, no "no results" message.
- **Query matches every product** (e.g. a single common letter): all cards remain visible; no message.
- **Special characters in query** (e.g. `.`, `*`, `(`): treated as literal text to match, not as patterns; they never cause an error.
- **Product with a missing/empty name** (should not occur in `products.json`, but defensively): never matches a non-empty query and is hidden while filtering is active.
- **Page reload while a query is typed**: the search box starts empty and the full grid is shown (no persistence — see Assumptions).

---

## Requirements *(mandatory)*

Requirements are stated as **observable behaviour**, not implementation.

### Functional Requirements

- **FR-001**: The product list page MUST display a search input control, visible and usable when the page first loads.
- **FR-002**: As the shopper changes the search input, the system MUST update which product cards are visible without reloading the page.
- **FR-003**: A product card MUST remain visible if, and only if, the product's name contains the query text as a substring; all other product cards MUST be hidden.
- **FR-004**: Matching MUST be case-insensitive ("SUN", "sun", and "Sun" produce identical results).
- **FR-005**: Leading and trailing whitespace in the query MUST be ignored when matching.
- **FR-006**: An empty or whitespace-only query MUST result in all product cards being visible, in their original order.
- **FR-007**: When a non-empty query matches no product, the system MUST show a visible "no products match" message and hide all product cards; the message MUST be removed when the query again matches at least one product or is cleared.
- **FR-008**: Filtering MUST operate only on the products already rendered on the page; the system MUST NOT issue any network request to a backend service to perform the search.
- **FR-009**: The relative display order of the products that remain visible MUST match their original order on the page (filtering hides, it does not reorder).
- **FR-010**: The feature MUST NOT change, remove, or break any existing product list behaviour (links to product detail pages, prices, images, recommendations).

### Key Entities

- **Product (as displayed)**: a catalogue item already rendered as a card on the product list page, sourced from the existing in-memory catalogue (`productcatalogservice/products.json`). Relevant attribute for this feature: **name**. The feature reads the name that is already present on the page; it does not fetch or persist any new data.

---

## Constraints — what this feature MUST NOT do *(mandatory)*

These are explicit guardrails copied from the request. Treat them as hard boundaries.

- **C-001**: MUST use only the services already in this repo. MUST NOT add new services.
- **C-002**: MUST NOT introduce Elasticsearch, Solr, vector databases, or any new datastore.
- **C-003**: MUST use the existing in-memory product catalogue loaded from `productcatalogservice/products.json`. Any filtering MUST be in memory.
- **C-004**: MUST match the language of the service being edited — the frontend is Go and `productcatalogservice` is Go. If `productcatalogservice` is extended at all, it MUST use the existing protobuf/gRPC patterns.
- **C-005**: MUST NOT add new infrastructure config, Helm charts, manifests, or environment variables. (CI deploys whatever lands on `attendee/<your-name>`.)
- **C-006**: MUST stay inside this branch and this repo. MUST NOT touch the build pipeline.
- **C-007**: Per the primary request, the intended design performs the filter **in the browser over products already loaded**, with **no backend change**. A backend path is out of scope unless explicitly re-scoped; if ever taken, C-003 and C-004 still bind.
- **C-008**: MUST NOT change the existing product list layout.
- **C-009**: MUST NOT add authentication.
- **C-010**: MUST NOT add a separate search results page; results MUST render inline on the existing product list page.

---

## Success Criteria *(mandatory)*

Measurable, technology-agnostic, and testable.

- **SC-001**: A shopper can locate a known product by typing part of its name and end up with only matching cards visible, in a single uninterrupted interaction (no page reload, no navigation).
- **SC-002**: For any query, the set of visible cards is exactly the set of products whose (case-insensitive, trimmed) name contains the query — verifiable by comparing against `products.json` names.
- **SC-003**: Filtering produces a visible result within one rendered frame of a keystroke as perceived by the user (no spinner, no blank flash, no waiting on the network).
- **SC-004**: Zero backend/network requests are generated by typing in the search box (observable in the browser network panel).
- **SC-005**: Clearing the search box restores 100% of the original product cards in their original order.
- **SC-006**: A query with no matches always shows the "no products match" message and never an empty, message-less grid.
- **SC-007**: No existing product-list behaviour regresses: product links, prices, images, and recommendations work exactly as before the feature.

---

## Assumptions

Informed defaults chosen where the request was silent:

- **Primary design is client-side.** The headline ("filters the products already loaded, by name, in the browser. No backend change") takes precedence. The constraints block's notes about extending `productcatalogservice` (Go, protobuf/gRPC) are treated as guardrails that bind *only if* a backend path is later chosen — not as a mandate to add one.
- **Match semantics**: case-insensitive **substring** match on the product **name** only (not description, categories, or price). This is the most common shopper expectation and the simplest testable rule.
- **No persistence**: the query is not saved across reloads and does not alter the URL. The box starts empty on each page load.
- **Scope is the product list page** (`home.html`) only. Search is not added to the header on every page, the cart, or the product detail page in this iteration.
- **No paging/lazy-loading assumption**: all products are already present in the rendered page (consistent with how `home.html` renders the full catalogue today), so "filter what's already loaded" covers the whole catalogue.
- **Accessibility default**: the search input has a visible or accessible label so it is usable by keyboard and screen-reader users; this is assumed, not specified in detail.
