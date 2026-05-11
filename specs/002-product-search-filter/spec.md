# Feature Specification: Frontend Product Search Filter

**Feature Branch**: `002-product-search-filter`
**Created**: 2026-05-11
**Status**: Draft
**Input**: "Add a product search feature to Online Boutique. Frontend-only filter. A search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

---

## Problem

Shoppers browsing the Online Boutique home page must scroll through every product to find what they are looking for. There is no way to narrow the visible list without navigating away from the page. This slows discovery and creates friction — especially as the catalogue grows.

We want a search box on the product list page that lets shoppers filter the visible products by name as they type, entirely in the browser. The products are already loaded when the page renders, so no additional data fetching is needed and no backend service needs to change.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Filter by partial name while typing (Priority: P1)

A shopper sees all products on the home page. They start typing a word they remember from a product name. With each keystroke the product grid narrows to show only cards whose names contain the typed text, and non-matching cards disappear. The shopper reaches their target product without scrolling.

**Why this priority**: This is the core behaviour. Every other story depends on it working first.

**Independent Test**: Load the home page. Type a substring of one product name (e.g. "mug"). Confirm only matching cards are visible. The page must not reload.

**Acceptance Scenarios**:

1. **Given** the home page is loaded with all product cards visible and the search input is empty, **When** the shopper types "mug", **Then** only product cards whose name contains "mug" (case-insensitive) are visible; all others are hidden.
2. **Given** a query of "mug" is active, **When** the shopper types an additional character making "mugs", **Then** the visible set narrows further in real time.
3. **Given** a query is active, **When** the shopper deletes one character, **Then** the visible set immediately widens to include any cards now matching the shorter query.
4. **Given** any query is active, **Then** no network request is made — the page remains on the same URL with no reload.

---

### User Story 2 — Clear the filter to see all products (Priority: P2)

A shopper has applied a filter and now wants to browse everything again. They click an explicit × clear button rendered next to the search box, and immediately all product cards reappear and the button disappears.

**Why this priority**: Without a visible, consistent reset control the filter becomes a trap. The native browser clear button is browser-specific and not always visible; an explicit button gives all users the same affordance.

**Independent Test**: Apply a filter that hides at least one product. Confirm the × button is visible. Click it. Confirm every product card is visible again and the × button is gone.

**Acceptance Scenarios**:

1. **Given** a non-empty query is hiding some product cards, **When** the shopper clicks the × clear button, **Then** the search input is cleared, all product cards become visible, the no-results message (if shown) disappears, and the × button disappears.
2. **Given** the search input is empty, **Then** the × clear button is not visible.
3. **Given** the shopper starts typing a query, **When** the first character is entered, **Then** the × clear button appears.
4. **Given** the page has just loaded, **When** the shopper has not typed anything, **Then** all product cards are visible, the search box is empty, and the × button is not shown.

---

### User Story 3 — Understand when nothing matches (Priority: P3)

A shopper types a query that matches no product name. Instead of an empty grid with no explanation, they see a clear "No products match your search" message so they know the filter is working and their query simply returned nothing.

**Why this priority**: Prevents confusion when the grid goes blank; low effort, high trust impact.

**Independent Test**: Type a string that cannot match any product (e.g. "zzzzzzz"). Confirm the grid is empty and the no-results message is displayed. Backspace to a matching query; confirm the message disappears and matching cards reappear.

**Acceptance Scenarios**:

1. **Given** a query that matches no product names, **When** the shopper finishes typing, **Then** all product cards are hidden and a "No products match your search" message is displayed in the grid area.
2. **Given** a no-results message is showing, **When** the shopper edits the query so at least one product matches, **Then** the message disappears and the matching card becomes visible.

---

### Edge Cases

- **Whitespace-only query**: Treated as an empty query — all products are visible.
- **Mixed-case query**: "MUG", "mug", and "Mug" all produce identical results.
- **Special characters in product names** (e.g. `&`, `<`, `>`): The filter compares displayed text; special characters must not break matching or cause visual errors.
- **Single product in catalogue**: Normal filter behaviour applies — the one card shows or hides based on the query.
- **Navigate away and return**: The search box resets to empty on page load; no filter state persists across navigations.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product list page MUST display a text search input above the product grid at all times.
- **FR-002**: The search input MUST have an accessible label (visible or screen-reader-accessible) that describes its purpose.
- **FR-003**: As the shopper types, the visible set of product cards MUST update after each keystroke with no perceptible delay.
- **FR-004**: Matching MUST be case-insensitive substring matching against the product's display name.
- **FR-005**: Product cards that do not match the current query MUST be hidden from view; they MUST NOT be removed from the page.
- **FR-006**: When the active query matches zero products, a "No products match your search" message MUST appear in the grid area.
- **FR-007**: An explicit × clear button MUST be rendered next to the search input whenever the input contains text, and MUST be hidden when the input is empty.
- **FR-007b**: Clicking the × button MUST clear the search input, restore all product cards to visible, and hide the no-results message.
- **FR-008**: The feature MUST NOT trigger any network request, page reload, or URL change as a result of typing in the search box.
- **FR-009**: The search input MUST be operable by keyboard alone (focus, type, clear).

### Key Entities

- **Product Card**: A tile on the home page representing one product; carries a display name that is the filter target.
- **Search Input**: The text field the shopper types into; drives the visible subset of product cards.
- **No-Results Message**: A UI element visible only when the filtered set is empty; hidden otherwise.

---

## Constraints *(what NOT to do)*

The following are hard boundaries. Crossing any of them is out of scope regardless of any other consideration.

- **No new services**: Use only the services already present in this repository. Do not add, stub, or reference any new microservice.
- **No new datastores**: Do not introduce Elasticsearch, Solr, a vector database, or any other search or storage technology. The product list is already in memory when the page loads — use it as-is.
- **No backend changes**: The productcatalogservice, any gRPC proto, and all server-side Go code for filtering are off-limits. The filter lives entirely in the browser.
- **No infrastructure changes**: Do not add or modify Helm charts, Kubernetes manifests, Terraform files, environment variables, or CI/CD pipeline configuration.
- **No new build dependencies**: Do not introduce npm packages, Go modules, or any external library for the filter logic.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can locate any specific product by typing a partial name; the matching card is visible before they finish typing the full name.
- **SC-002**: The product grid updates with no perceptible delay after each keystroke — the shopper never waits for a result.
- **SC-003**: Clearing the search box restores all product cards immediately; no action other than clearing the input is required.
- **SC-004**: When no products match, a clear explanatory message is displayed 100% of the time — the grid is never left blank without explanation.
- **SC-005**: Shoppers who do not interact with the search box see no change to the existing page experience; the default state (search empty, all products visible) is indistinguishable from the pre-feature page.
- **SC-006**: The feature works correctly when accessed by keyboard only (no mouse required).

---

## Assumptions

- The full product catalogue is rendered into the page on initial load. There is no pagination, infinite scroll, or lazy loading to account for.
- Filtering is by product **name** only. Filtering by category, description, or price is out of scope.
- No filter state is persisted across page reloads or navigations.
- The feature applies to the home page product grid only. A dedicated search route or search on other pages is out of scope.
- All shoppers see the same product list; there are no per-user catalogue variations.
- The existing Bootstrap utility classes already loaded by the page are sufficient for styling the search input without adding new CSS files.
