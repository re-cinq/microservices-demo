# Feature Specification: Product Search

**Feature Branch**: `003-product-search`
**Created**: 2026-05-11
**Status**: Draft
**Input**: User description: "Add a product search feature to Online Boutique. A search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

## Clarifications

### Session 2026-05-11

- Q: When the visible product count changes as the shopper types, should the new count be announced to screen reader users? → A: Yes — via a polite ARIA live region. Announce "N products match" (or "1 product matches" for N=1) and "No products match your search." for N=0. Empty/whitespace query announces nothing.
- Q: Should the search input include a visible clear (×) affordance, and if so, from where? → A: Yes — via the browser-native `<input type="search">` clear button. No custom DOM × button.

(Both were originally surfaced as Deferred items by `/speckit-clarify`; the decisions were made during `/speckit-plan` and are documented with rationale in `research.md` R5 and R6. They are mirrored here so this spec is self-contained.)

## Problem *(in the user's words)*

> "I'm a shopper on Online Boutique. The product list shows me everything at once and I have to scroll. When I already know what I'm looking for — 'mug', 'shirt', 'glasses' — I want to type it into a box at the top of the list and have the page narrow down to only the matching products. It should feel instant. I don't want a separate search results page; I want the same product grid, just filtered."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter the product list as I type (Priority: P1)

A shopper lands on the product list page, sees a search box above the product grid, types part of a product name, and the grid immediately shrinks to only products whose names match what they typed. No page reload, no spinner, no navigation.

**Why this priority**: This is the entire feature. Without it the search box has no value. It is the smallest slice that delivers the user's goal.

**Independent Test**: Open the product list page, type "mug" into the search box, observe that only products whose names contain "mug" remain visible in the grid. Achievable without any of the other stories.

**Acceptance Scenarios**:

1. **Given** the shopper is on the product list page with all products visible, **When** they type "sunglasses" into the search box, **Then** only products whose name contains "sunglasses" (case-insensitive) remain visible in the grid.
2. **Given** the shopper has typed "mu" into the search box, **When** they extend their input to "mug", **Then** the visible products further narrow to those matching "mug" without any page reload.
3. **Given** the shopper has typed "Sun" into the search box, **When** the grid updates, **Then** the product "Sunglasses" is visible (i.e. matching is case-insensitive).

---

### User Story 2 - Clear the search to see everything again (Priority: P1)

A shopper who has narrowed the list can empty the search box (delete their text) and the full product grid is restored exactly as it was before they started typing.

**Why this priority**: Without this, a user who searches is trapped in a filtered view. P1 because it is the immediate inverse of Story 1 and is needed for the feature to feel safe to use.

**Independent Test**: Type any query, observe filtering, delete all characters from the search box, observe that every product is visible again.

**Acceptance Scenarios**:

1. **Given** the shopper has typed "mug" and the grid is filtered, **When** they delete all characters from the search box, **Then** every product that was originally on the page is visible again in its original order.
2. **Given** the shopper has typed a query that matches nothing, **When** they delete all characters, **Then** every product is visible again.

---

### User Story 3 - See a clear "no results" state (Priority: P2)

When a shopper's query matches zero products, the grid is empty and a short message tells them no products matched, so they don't think the page is broken.

**Why this priority**: P2 because Stories 1 and 2 are functionally usable without it, but an empty grid with no message is confusing.

**Independent Test**: Type a string that does not appear in any product name (e.g. "zzzz"), observe that the grid is empty and a "no results" message is shown.

**Acceptance Scenarios**:

1. **Given** the shopper is on the product list page, **When** they type a string that matches no product name, **Then** the grid is empty and a message indicating "no products match your search" is shown in its place.
2. **Given** the "no results" message is shown, **When** the shopper clears the search box, **Then** the message disappears and the full grid is restored.

---

### Edge Cases

- **Whitespace-only query** ("   "): treated as an empty query — full grid is shown, no "no results" message.
- **Leading/trailing whitespace** ("  mug  "): trimmed; behaves the same as "mug".
- **Special characters** ("mug!@#"): matched literally against product names; if no product name contains that substring, the "no results" message is shown.
- **Very long query** (e.g. 500 characters): does not crash the page; produces the "no results" message.
- **Rapid typing**: the grid update keeps up with the user's typing without visible lag (no debounce-induced "stale" results visible for more than the next keystroke).
- **Navigation away and back**: when the shopper navigates away from the product list page and returns, the search box is empty and the full grid is shown (no persisted filter state).

## Requirements *(mandatory — observable behaviour only)*

### Functional Requirements

- **FR-001**: The product list page MUST display a single search input above the product grid.
- **FR-002**: The search input MUST be reachable and operable by keyboard alone (tab to focus, type to filter).
- **FR-003**: As the shopper types in the search input, the visible set of products in the grid MUST update on every keystroke without a page reload or navigation.
- **FR-004**: A product MUST be visible in the grid if and only if its name contains the shopper's query as a substring, compared case-insensitively, after trimming leading/trailing whitespace from the query.
- **FR-005**: When the query is empty (or whitespace-only), the grid MUST show every product that was originally rendered on the page, in the original order.
- **FR-006**: When the query matches zero products, the grid area MUST show a short, human-readable message indicating no products match. The message MUST disappear as soon as the query matches one or more products or is cleared.
- **FR-006a**: A polite ARIA live region MUST announce the current visible-result count after each filter update: "N products match" (or "1 product matches" when N=1) when N ≥ 1, "No products match your search." when N = 0, and silence (empty announcement) when the query is empty or whitespace-only.
- **FR-006b**: The search input MUST expose a visible clear control via the browser-native `<input type="search">` clear (×) button. No custom × button is added to the DOM.
- **FR-007**: Filtering MUST operate only on the products already loaded into the page; the feature MUST NOT issue any additional network requests to the product catalogue or any other backend as a result of the shopper typing.
- **FR-008**: The order of products that remain visible MUST be the same as their order in the original (unfiltered) grid.
- **FR-009**: The feature MUST NOT alter the product list page's behaviour when the search box is empty (the page MUST look and behave exactly as it did before this feature was added).
- **FR-010**: The search state MUST NOT persist across page navigations: returning to the product list page from elsewhere MUST show an empty search input and the full grid.

### Explicit Constraints *(what this feature MUST NOT do)*

- **C-001**: MUST NOT introduce any new service in this repository.
- **C-002**: MUST NOT introduce Elasticsearch, Solr, a vector database, or any other new datastore.
- **C-003**: MUST NOT change the source of product data: the existing in-memory product catalogue loaded from `productcatalogservice/products.json` remains the only source.
- **C-004**: MUST NOT add or modify Helm charts, Kubernetes manifests, Kustomize overlays, Skaffold configuration, Terraform, or any other infrastructure-as-code.
- **C-005**: MUST NOT add or rename environment variables consumed by any service.
- **C-006**: MUST NOT modify the build pipeline, Cloud Build configuration, or GitHub Actions workflows.
- **C-007**: MUST NOT add a new gRPC method, protobuf message, or backend endpoint. (No backend change of any kind.)
- **C-008**: MUST NOT introduce a new programming language to any edited service. Frontend changes stay in the existing frontend stack; if `productcatalogservice` were ever extended, it would stay in Go and reuse existing protobuf/gRPC patterns — but in this feature it is not edited at all.
- **C-009**: MUST NOT search any field other than the product's `name`. Description, category, and ID are out of scope for matching.
- **C-010**: MUST NOT issue any new network request triggered by typing in the search input.

### Key Entities

- **Product**: An item already present on the product list page, with a human-readable `name` used for matching. No new fields are introduced.
- **Query**: The current text content of the search input, trimmed of leading/trailing whitespace and compared case-insensitively against product names.

## Success Criteria *(mandatory — measurable, technology-agnostic)*

### Measurable Outcomes

- **SC-001**: A shopper who knows the product name can reach a filtered grid containing the desired product in **3 keystrokes or fewer** for any product whose name has a unique 3-letter prefix among the catalogue.
- **SC-002**: The grid updates to reflect the current query within **100 ms of each keystroke**, as perceived by the shopper (no visible lag between keystroke and grid update).
- **SC-003**: Typing in the search input generates **zero additional network requests** to any backend service (measurable in the browser's network panel).
- **SC-004**: For any query that is a case-insensitive substring of one or more product names, the filtered grid shows **100 % of matching products and 0 % of non-matching products** from the original page.
- **SC-005**: When the search input is empty, the product list page is **byte-for-byte equivalent in observable behaviour** to the page before this feature was introduced (same product set, same order, same network calls on load).
- **SC-006**: A shopper who types a non-matching query sees the "no results" message within the same **100 ms** budget; clearing the input restores the full grid within the same budget.

## Assumptions

- The product list page already renders all products in a single response on initial page load (no pagination, no lazy loading). The current Online Boutique frontend does this; filtering "the products already loaded" is therefore the entire catalogue.
- Product names are short, human-readable strings (e.g. "Sunglasses", "Mug"); substring matching on `name` is sufficient for the shopper's goal and no stemming, fuzzy matching, or synonym expansion is required.
- Case-insensitive substring matching is the expected behaviour ("sun" matches "Sunglasses"). No locale-specific collation is required for v1.
- The feature is for the desktop and mobile web product list page only; no native app, no API consumer.
- Accessibility expectations are the same as the rest of the Online Boutique UI: keyboard reachable, screen-reader-labelled input. No new accessibility regression is acceptable, but no new AA-conformance work beyond what already applies to the page is in scope.
- Analytics/telemetry on search usage is out of scope for v1.
