# Feature Specification: Category Jump Links on Product List Page

**Jira Story**: [AIP-187](https://odevo.atlassian.net/browse/AIP-187) — Category jump links reduce scrolling on the product list page
**Jira Epic**: [AIP-98](https://odevo.atlassian.net/browse/AIP-98) — Category browse and filter

**Feature Branch**: `attendee/pontus-vepsalainen`

**Created**: 2026-06-16

**Status**: Draft

**Input**: Derived from Jira AIP-187 (story) and AIP-98 (epic technical constraints)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Category Jump Links on Page Load (Priority: P1)

A shopper lands on the product list page and immediately sees a set of category anchor links at the top. Products below are grouped under category headings in the same order as today.

**Why this priority**: This is the foundational UI change that all other stories depend on. Without it, jump navigation has nothing to link to.

**Independent Test**: Load the product list page and confirm category headings and jump links appear at the top, and products are visually grouped under their category.

**Acceptance Scenarios**:

1. **Given** the shopper is on the product list page, **When** the page loads, **Then** category headings and jump links are visible at the top of the page.
2. **Given** the page has loaded, **When** the shopper views the product list, **Then** products are grouped under their category heading in the same order as today.

---

### User Story 2 - Jump to Category Section (Priority: P2)

A shopper clicks a category jump link and the page scrolls directly to that category's section, eliminating the need to manually scroll through the full catalogue.

**Why this priority**: This is the core value of the feature — reducing scrolling for the shopper.

**Independent Test**: Click a jump link and verify the viewport scrolls to the matching category section heading.

**Acceptance Scenarios**:

1. **Given** the shopper clicks a jump link, **When** the page responds, **Then** the view scrolls to the corresponding category section.

---

### Edge Cases

- What happens when a category has no products? (Category heading should still render if returned by the catalogue, or be omitted entirely — no orphaned jump links.)
- What happens if product data has no category field? (Products without a category are grouped under a generic fallback section and no jump link for that section is shown, or they appear at the end without a jump link.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product list page MUST display a list of category jump links at the top of the page on load.
- **FR-002**: Each jump link MUST correspond to a category that has at least one product.
- **FR-003**: Clicking a jump link MUST scroll the viewport to the matching category section heading on the same page.
- **FR-004**: Products MUST be visually grouped under their category heading, preserving the existing relative order of products within each category.
- **FR-005**: The change MUST be implemented in the frontend product list template only — no backend service changes, no new endpoints, no new data stores.
- **FR-006**: The feature MUST use only services that already exist in the repository.
- **FR-007**: No new datastore (database, cache, or search engine) MAY be introduced; category grouping MUST be derived in-memory from the existing product catalogue (`productcatalogservice/products.json`).
- **FR-008**: All changes MUST match the language and patterns of the service being modified (Go for both the frontend and the product catalogue service).
- **FR-009**: Infrastructure, deployment manifests, and CI configuration MUST NOT be changed.

### Key Entities

- **Category**: A named grouping for products, derived from existing product data. Has a display name and acts as an anchor target on the page.
- **Product**: An item in the catalogue with an associated category. Rendered under its category heading.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can reach any product category section without scrolling past unrelated products — one click from the top of the page.
- **SC-002**: All acceptance scenarios from AIP-187 pass on the product list page without any backend or infrastructure changes.
- **SC-003**: Product order within each category is identical to the pre-feature product list order.
- **SC-004**: No new services, datastores, or deployment manifests are added as a result of this change.

## Assumptions

- The existing product data in `productcatalogservice/products.json` already contains a category field for each product; no data migration is needed.
- Product ordering within a category is preserved from the current list (i.e., products already sorted or ordered as-is).
- Mobile/responsive layout is not in scope for this story; desktop layout is the primary target.
- No search or multi-filter capability is added — this is a pure navigation aid (jump links), not a filter.
- The frontend template is rendered server-side in Go; JavaScript is used only for the native anchor-scroll behaviour already supported by browsers.
