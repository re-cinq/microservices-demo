# Feature Specification: Recently Viewed Products Strip

**Feature Branch**: `006-aip-163`

**Created**: 2026-05-20

**Status**: Draft

**Input**: User description: "AIP-163 — Shopper can return to a recently viewed product via the strip"

**Jira**: [AIP-163](https://odevo.atlassian.net/browse/AIP-163)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse and Return to Product (Priority: P1)

After viewing one or more product pages during a session, the shopper sees a "Recently viewed" strip on the page. They can click any item in the strip to go directly back to that product's detail page without having to search for it again.

**Why this priority**: This is the core value proposition of the feature — enabling frictionless return navigation to products the shopper has already shown interest in.

**Independent Test**: Can be fully tested by visiting two or more product pages in sequence and verifying the strip appears and that clicking an entry navigates back to the correct product.

**Acceptance Scenarios**:

1. **Given** I have viewed at least one product page in this session, **When** the strip is present on the page, **Then** each recently viewed product appears in the strip with enough detail to identify it (e.g., product name, image, and/or price).
2. **Given** the strip is visible, **When** I click an entry in the strip, **Then** I am taken directly to that product's detail page.

---

### User Story 2 - Product View is Recorded (Priority: P2)

When a shopper loads a product detail page, that product is automatically recorded as viewed and subsequently appears in the recently viewed strip.

**Why this priority**: This is the data capture mechanism that makes the strip possible. Without it, the strip has nothing to show. It is independently testable as a background behaviour.

**Independent Test**: Can be tested by visiting a single product page and confirming the product appears in the strip on the next page visited.

**Acceptance Scenarios**:

1. **Given** I load a product detail page, **When** the page completes loading, **Then** that product is recorded as viewed and will subsequently appear in the strip on any page that displays it.
2. **Given** I have viewed a product that already appears in the strip, **When** I view it again, **Then** it is not duplicated in the strip (it moves to the most recent position or stays in place).

---

### Edge Cases

- What happens when the shopper has viewed only one product — does the strip appear with a single entry?
- How does the system handle a product that has since been removed from the catalogue — does it still appear in the strip?
- What happens when the shopper opens the same product in a new tab — is it recorded twice?
- How does the strip behave if session data is lost or the browser is refreshed — is history preserved within the session?
- What is the maximum number of products shown in the strip? The strip shows a maximum of 5 most recently viewed products.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST record each product detail page a shopper visits during a session as "recently viewed".
- **FR-002**: System MUST display a "Recently viewed" strip to shoppers who have viewed at least one product in the current session.
- **FR-003**: The strip MUST show each recently viewed product with its name, image, and current price.
- **FR-004**: Users MUST be able to click any entry in the strip to navigate directly to that product's detail page.
- **FR-005**: System MUST NOT show duplicate entries for the same product in the strip.
- **FR-006**: System MUST limit the strip to the 5 most recently viewed products.

### Key Entities

- **Session**: A shopper's current browsing session; scopes the recently viewed history.
- **Recently Viewed Product**: A product detail page visited by the shopper in the current session; has at minimum a product identifier and timestamp of viewing.
- **Recently Viewed Strip**: The UI component that renders the list of recently viewed products in order of recency.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Shoppers who have viewed at least one product see the recently viewed strip on subsequent pages without any additional action.
- **SC-002**: Clicking a strip entry takes the shopper to the correct product detail page in under 1 second under normal conditions.
- **SC-003**: The strip accurately reflects the shopper's view history for the current session with no missing or duplicated entries.
- **SC-004**: The feature is delivered as a complete minimum-viable version with all three acceptance criteria from AIP-163 verified.

## Constraints *(inherited from Epic [AIP-97](https://odevo.atlassian.net/browse/AIP-97))*

These are hard constraints that must not be violated by any implementation decision:

- **C-001**: The solution MUST use only services that already exist in the repository. No new services may be introduced.
- **C-002**: No new datastore of any kind may be introduced (no database, cache, or search engine). View history must be held in memory, drawing product data from the existing product catalogue.
- **C-003**: Any code changes MUST match the language and patterns of the service being changed. The frontend and product catalogue service are written in Go.
- **C-004**: Infrastructure, deployment manifests, and CI configuration MUST NOT be changed. The feature must ship through the existing pipeline without modification.

## Assumptions

- Recently viewed history is scoped to the current browser session only; it does not persist across sessions or devices (MVP scope as stated in AIP-163).
- The strip is shown on product detail pages at minimum; exact page placement is a layout/design decision.
- A product is "viewed" when its detail page completes loading (not on hover or partial scroll).
- The strip displays a maximum of 5 most recently viewed products, each showing name, image, and current price.
- No authentication is required — the feature works for anonymous shoppers.
- Dependencies: none (as stated in AIP-163).
