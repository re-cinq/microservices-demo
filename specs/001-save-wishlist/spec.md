# Feature Specification: Save for Later / Wishlist

**Feature Branch**: `005-save-wishlist`

**Created**: 2026-05-20

**Status**: Draft

**Input**: User description: "AIP-159: Save a product for later and view the saved list."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Save a product from the product detail page (Priority: P1)

A shopper browsing a product they like but aren't ready to buy can mark it with a single click. They receive immediate confirmation that the product has been saved, so they can continue browsing with confidence that the item will be waiting for them.

**Why this priority**: This is the entry point for the entire feature. Without the ability to save, there is nothing to view. Delivering this alone gives shoppers a way to bookmark interest for the first time.

**Independent Test**: Can be fully tested by visiting any product detail page, clicking "Save for later", and verifying the inline confirmation appears and the product is retrievable during the same session.

**Acceptance Scenarios**:

1. **Given** a shopper is on a product detail page, **When** they click "Save for later", **Then** the product is added to their wishlist and an inline confirmation message is displayed on the page.
2. **Given** a shopper clicks "Save for later" on a product already in their wishlist, **When** the action is processed, **Then** the product is not duplicated in the wishlist and the confirmation is still shown.

---

### User Story 2 - View saved products in the wishlist (Priority: P2)

A shopper who has saved one or more products can navigate to their wishlist at any point during the same visit to review everything they have set aside. Each saved product is shown with enough detail (name, image, price) to remind them what they saved and help them decide what to do next.

**Why this priority**: Viewing saved items is the payoff for the save action. Without it, saving has no visible benefit. Together with Story 1, this completes the core loop.

**Independent Test**: Can be fully tested by saving at least one product, navigating to the wishlist view, and verifying the correct product details are displayed.

**Acceptance Scenarios**:

1. **Given** a shopper has one or more saved products, **When** they navigate to their wishlist, **Then** they see each saved product displayed with its name, image, and price.
2. **Given** a shopper has no saved products, **When** they navigate to their wishlist, **Then** they see a clear empty-state message rather than a blank page.

---

### User Story 3 - Wishlist is cleared when the session ends (Priority: P3)

A shopper's saved list exists only for the duration of their current visit. When they close the tab or return in a new session, the wishlist starts empty. No personal data is retained between visits.

**Why this priority**: This is a privacy and scope boundary — it is explicitly required and must be verifiable, but it does not add new user-visible capability beyond Stories 1 and 2.

**Independent Test**: Can be fully tested by saving products in one session, ending the session, starting a new one, and confirming the wishlist is empty.

**Acceptance Scenarios**:

1. **Given** a shopper has saved products during a session, **When** they end the session and start a new one, **Then** their wishlist is empty.

---

### Edge Cases

- What happens when a shopper saves a product that is later out of stock? (Price and availability shown on the wishlist reflect current state at time of viewing, not time of saving.)
- What happens if a shopper navigates to the wishlist without having saved anything? (Empty-state message is shown — see Story 2, Scenario 2.)
- What happens if a shopper saves the same product more than once? (No duplicate is created — see Story 1, Scenario 2.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A shopper MUST be able to save any product to their wishlist from that product's detail page.
- **FR-002**: The system MUST display an inline confirmation to the shopper immediately after a product is saved.
- **FR-003**: A shopper MUST be able to view all their saved products in a single wishlist view during the same session.
- **FR-004**: The wishlist view MUST display each saved product's name, image, and price.
- **FR-005**: The wishlist MUST be automatically cleared when the shopper's session ends.
- **FR-006**: The feature MUST be available to all shoppers without requiring account creation or login.
- **FR-007**: Saving a product that is already in the wishlist MUST NOT create a duplicate entry.
- **FR-008**: When the wishlist is empty, the wishlist view MUST display a user-friendly empty-state message.

### Key Entities

- **Wishlist**: A session-scoped, ordered collection of products a shopper has saved. Exists only for the duration of the current session. Attributes: session identifier, ordered list of saved product references.
- **Saved Product**: A reference to a product the shopper has bookmarked. Attributes: product identifier, product name, image, price (resolved at display time from the current product catalogue).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can save a product and see the confirmation without any perceptible delay (under 1 second from click to confirmation).
- **SC-002**: A shopper can view all products they saved in a single visit on one page without needing to paginate, for a typical session (up to 20 saved items).
- **SC-003**: 100% of wishlist data is discarded at session end — no saved products carry over to a new session.
- **SC-004**: The wishlist feature is accessible to 100% of shoppers without registration or login.

## Assumptions

- The wishlist entry point (the link or navigation element that takes a shopper to their wishlist view) is a dedicated page accessible via the store's main navigation. The exact placement within the navigation is a design decision outside this specification.
- A shopper's wishlist is independent per browser session. Two concurrent sessions (e.g. two browser tabs opened independently) maintain separate wishlists.
- Product price and availability shown in the wishlist reflect the current catalogue state at the time the wishlist page is loaded, not at the time the product was saved.
- Multi-tab behaviour within a single session (e.g. two tabs open in the same browser session) is out of scope for this version.
- The save control appears only on the product detail page. Saving from other pages (e.g. search results, homepage) is out of scope.
