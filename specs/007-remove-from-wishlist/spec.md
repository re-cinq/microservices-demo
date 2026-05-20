# Feature Specification: Remove a product from the wishlist

**Feature Branch**: `007-remove-from-wishlist`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Jira story AIP-167: "Remove a product from the wishlist"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Remove a saved product from the wishlist (Priority: P1)

A shopper reviewing their saved products can remove individual items they no longer want. Each product card on the wishlist page shows a "Remove" button; clicking it removes only that product and refreshes the page to reflect the new state.

**Why this priority**: This is the sole story for this feature. Without it the wishlist grows indefinitely with no way to manage it, which undermines its usefulness.

**Independent Test**: Can be fully tested by saving at least one product, navigating to `/wishlist`, clicking "Remove", and verifying the product no longer appears in the list.

**Acceptance Scenarios**:

1. **Given** a user has one or more products saved in their wishlist, **When** they click "Remove" next to a product, **Then** that product is immediately removed from the list and the wishlist page re-renders without it.
2. **Given** a user has exactly one product saved in their wishlist, **When** they remove it, **Then** the wishlist page shows the empty-state message rather than a product grid.
3. **Given** a user has removed a product from their wishlist, **When** they navigate to that product's detail page, **Then** the page does not show the "Saved to your wishlist" confirmation banner.

---

### Edge Cases

- What happens when the product_id form field is missing or empty? (Handler returns HTTP 400 Bad Request — mirrors `saveWishlistHandler` validation.)
- What happens when a product_id is submitted that isn't in the wishlist? (No-op: filtering produces the same slice; redirect to `/wishlist` proceeds normally.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A "Remove" button MUST appear on each product card in the wishlist view.
- **FR-002**: Clicking "Remove" MUST issue a `POST /wishlist/remove` request with the product's ID.
- **FR-003**: The handler MUST remove the matching product ID from the session's wishlist slice and store the result back.
- **FR-004**: After removal, the user MUST be redirected to `/wishlist` (HTTP 303).
- **FR-005**: The `product_id` form field MUST be validated; a missing value MUST return HTTP 400.

### Key Entities

- **Wishlist** (`fe.wishlists sync.Map`): keyed by session ID (string) → `[]string` of product IDs. Shared with AIP-159. No structural change required.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A product removed from the wishlist is no longer visible on the wishlist page within the same HTTP round-trip (redirect + re-render).
- **SC-002**: Removing the last item results in the empty-state message rendering correctly.
- **SC-003**: No product data is persisted beyond the in-memory `sync.Map` (privacy constraint from AIP-94).

## Assumptions

- AIP-159 is merged and the `wishlists sync.Map`, `viewWishlistHandler`, and `wishlist.html` template are present in `src/frontend/`.
- The `?saved=1` confirmation banner on the product page is URL-param-driven; no changes to its logic are needed to satisfy AC #3.
- Wishlist state is session-scoped and in-memory; no cross-session persistence is in scope.
