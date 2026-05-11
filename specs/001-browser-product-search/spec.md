# Feature Specification: Browser-Side Product Search

**Feature Branch**: `attendee/joel-sanmoogan`
**Created**: 2026-05-11
**Status**: Draft

## Problem

Shoppers on the Online Boutique home page must scroll through every product to find what they want. There is no way to narrow the list down. A search box on the product list page that instantly filters by product name — without reloading the page or calling the server — would let shoppers find products immediately.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Filter products by name (Priority: P1)

A shopper types into a search box on the product list page and sees only the products whose names contain their search text. Products that do not match disappear. Clearing the box restores all products.

**Why this priority**: This is the entire feature. All other stories depend on it.

**Independent Test**: Open the product list page, type a partial product name into the search box, and verify only matching products remain visible.

**Acceptance Scenarios**:

1. **Given** the product list page is loaded with all products visible, **When** the shopper types a word that matches some product names, **Then** only products whose names contain that word (case-insensitive) remain visible; all others are hidden.
2. **Given** the shopper has typed a search term and the list is filtered, **When** the shopper clears the search box, **Then** all products are shown again.
3. **Given** the product list is filtered, **When** the shopper types additional characters that narrow the match to zero products, **Then** no product cards are shown and an empty-state message is displayed.

---

### User Story 2 — Case-insensitive matching (Priority: P2)

A shopper who types in any mix of upper and lower case still finds the products they are looking for.

**Why this priority**: Without case-insensitive matching, the filter is frustrating to use because capitalisation varies across product names.

**Independent Test**: Type the same search term in all-caps and all-lowercase; both must produce identical result sets.

**Acceptance Scenarios**:

1. **Given** a product named "Sunglasses" exists, **When** the shopper types "sunglasses", **Then** the product is visible.
2. **Given** a product named "Sunglasses" exists, **When** the shopper types "SUNGLASSES", **Then** the product is visible.

---

### User Story 3 — Instant feedback while typing (Priority: P2)

The product list updates on every keystroke without the shopper pressing Enter or a button.

**Why this priority**: Live filtering makes the feature feel responsive and is the expected behaviour for in-page search.

**Independent Test**: Type one character at a time and observe that the list updates after each character, with no page reload.

**Acceptance Scenarios**:

1. **Given** the search box is focused, **When** the shopper presses any alphanumeric key, **Then** the product list updates immediately to reflect the new filter text without a page reload.

---

### User Story 4 — Accessibility: keyboard-navigable search (Priority: P3)

A keyboard-only shopper can reach the search box using the Tab key and type without using a mouse.

**Why this priority**: Basic accessibility; the search box must not be a dead end for keyboard users.

**Independent Test**: Tab to the search box, type a term, and confirm the list filters without touching the mouse.

**Acceptance Scenarios**:

1. **Given** the shopper uses only the keyboard, **When** they Tab to the search box and type a term, **Then** the product list filters correctly.

---

### Edge Cases

- What happens when the search term contains only spaces? → The full product list is displayed (treated as an empty query).
- What happens when no products match the search term? → An empty-state message such as "No products found" is displayed instead of a blank page.
- What happens when the product catalogue has not yet finished loading? → The search box is present but disabled (or ignored) until products are available.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product list page MUST display a text input field labelled for product search, visible without scrolling.
- **FR-002**: The product list MUST update to show only products whose names contain the current search text, evaluated on every change to the input value.
- **FR-003**: Name matching MUST be case-insensitive.
- **FR-004**: When the search input is empty or contains only whitespace, ALL products MUST be shown.
- **FR-005**: When no products match the search text, the product list area MUST display a visible empty-state message instead of a blank region.
- **FR-006**: Filtering MUST operate against the set of products already rendered in the page; no additional network requests are made when the user types.
- **FR-007**: The search input MUST be reachable and operable via keyboard alone.

### Constraints (what NOT to do)

- **C-001**: MUST NOT add, modify, or remove any backend service, gRPC endpoint, or protobuf definition.
- **C-002**: MUST NOT introduce any new datastore, search engine, or external dependency (Elasticsearch, Solr, vector databases, Redis, etc.).
- **C-003**: MUST NOT add new Kubernetes manifests, Helm chart entries, kustomize components, or environment variables.
- **C-004**: MUST NOT modify the CI/CD pipeline or build configuration.
- **C-005**: MUST NOT add new services to the repository.
- **C-006**: Changes MUST be confined to the `attendee/joel-sanmoogan` branch of this repository.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can locate any product by name in 3 seconds or fewer from landing on the product list page, by typing into the search box.
- **SC-002**: The product list visibly updates within 100 ms of each keystroke on a standard laptop browser (no network round-trip required).
- **SC-003**: Searching with an all-lowercase term and the same term in all-uppercase produces identical result sets 100% of the time.
- **SC-004**: When no products match, shoppers see an explicit message rather than a blank list, confirmed by manual test on every supported browser.
- **SC-005**: The search box is reachable and fully operable using the keyboard alone (Tab + typing), verified by keyboard-only navigation test.
- **SC-006**: The existing suite of frontend and productcatalogservice tests continues to pass without modification.

---

## Assumptions

- The full product catalogue is loaded and rendered in the page's HTML before the shopper interacts with the search box; no lazy-loading or pagination is in use on the product list page.
- The number of products in `products.json` is small enough (currently 9) that client-side filtering introduces no perceptible performance cost.
- The feature targets desktop browsers; responsive/mobile layout adjustments for the search box are out of scope for this iteration.
- No authentication or personalisation logic is involved; the same products are shown to all users before filtering.
- Matching is limited to product names; searching by description, category, or price is out of scope.
