# Feature Specification: Recently Viewed Strip on the Home Page

**Feature Branch**: `007-recently-viewed-home`  
**Created**: 2026-05-20  
**Status**: Draft  
**Input**: User description: "AIP-157 — A shopper who returns to the Online Boutique home page during a browsing session can see a strip of products they have already viewed, so they can quickly get back to any of them without searching or navigating back through the catalogue."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Recently Viewed Strip on Home Page (Priority: P1)

A shopper who has viewed one or more products during their current browsing session returns to the home page. A "Recently Viewed" strip appears below the main product grid, showing the products they browsed. Clicking any item takes them directly to that product's detail page.

**Why this priority**: This is the entire feature — there is only one user journey. It closes the loop between product browsing and the home page, letting shoppers resume where they left off without retracing navigation steps.

**Independent Test**: View one or more product detail pages, then navigate to the home page. Verify the recently viewed strip appears with the correct products and each links to its product detail page.

**Acceptance Scenarios**:

1. **Given** a shopper has viewed at least one product, **When** they navigate to the home page, **Then** a recently viewed strip is visible showing the products they viewed.
2. **Given** a shopper clicks a product in the recently viewed strip, **When** the click is registered, **Then** they are taken to that product's detail page.
3. **Given** a shopper has not viewed any products during their session, **When** they visit the home page, **Then** no recently viewed strip is shown.

---

### Edge Cases

- What happens when the session cookie holding recently viewed IDs is empty or missing? → No strip is shown; page renders normally.
- What happens when a recently viewed product ID no longer exists in the catalog? → That item is silently omitted from the strip (the catalog lookup returns no result for it).
- What happens when the shopper has viewed many products? → The strip shows as many products as the shared component supports (consistent with the product detail page behavior established in AIP-156).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The home page MUST display a "Recently Viewed" strip when the shopper's session contains at least one viewed product ID.
- **FR-002**: The recently viewed strip MUST show product image, name, and price for each viewed product.
- **FR-003**: Each item in the recently viewed strip MUST link to that product's detail page.
- **FR-004**: The home page MUST NOT display the recently viewed strip when no products have been viewed in the current session.
- **FR-005**: The recently viewed strip on the home page MUST use the same visual component used on the product detail page (established in AIP-156).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper who has viewed at least one product sees the recently viewed strip on the home page on 100% of visits within the same session.
- **SC-002**: A shopper with no viewed products sees no recently viewed strip on 100% of home page visits.
- **SC-003**: Every item in the strip links correctly to its product detail page.
- **SC-004**: The home page renders without error for all session states (empty, one item, multiple items).

## Assumptions

- The session tracking and recently viewed strip UI component are fully established by AIP-156 (Done); this story adds the strip to the home page only.
- Session state is stored in a browser cookie; no server-side session storage is required.
- The strip component renders identically on the home page as it does on the product detail page — no design changes are needed.
- Mobile layout is handled by the existing responsive styles inherited from the shared component.
