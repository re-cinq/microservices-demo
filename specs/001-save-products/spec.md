# Feature Specification: Save products and view them later

**Feature Branch**: `001-save-products`

**Created**: 2026-06-16

**Status**: Draft

**Input**: User description: "Story 1" — Jira AIP-192 (epic AIP-94, *Wishlist / Save for later*)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Save a product and come back to it (Priority: P1)

A shopper browsing the store finds a product they're interested in but isn't ready to buy. Without signing in, they save it from the product's page, keep browsing, and later open a saved-products view — reachable from a link in the site header — to find everything they've saved and go straight back to any of those products during the same visit. When they haven't saved anything yet, the saved view tells them so and how to start.

**Why this priority**: This is the entire promise of the epic — *save products and come back to them* — delivered end to end. On its own it turns "I lost that product when I left" into "I can set it aside and return to it," which is customer-visible value with nothing else built. Every other story in the epic (saving from the product list, removing, a saved-count badge) only enhances this one.

**Independent Test**: Save a product from its page, navigate elsewhere in the store, open the saved-products view from the header, and confirm the product is listed and links back to its detail page — all without signing in and within a single browsing session.

**Acceptance Scenarios**:

1. **Given** I'm on a product's page and haven't saved it, **When** I choose "Save", **Then** the product is added to my saved products and the control shows it as saved.
2. **Given** I've saved one or more products, **When** I open the saved-products view from the header link, **Then** I see each saved product and can click through to its product page.
3. **Given** I've saved a product, **When** I move around the store and return during the same session, **Then** the product is still saved, with no sign-in required.
4. **Given** I have not saved any products, **When** I open the saved-products view, **Then** I see a message that nothing is saved yet and guidance on how to start saving.

---

### Edge Cases

- **Saving the same product twice**: a product the shopper already saved is not duplicated in the saved view; choosing "Save" again has no additional effect (the control already reflects the saved state).
- **A saved product becomes unavailable**: if a saved product is no longer in the catalogue, the saved view does not break — the shopper sees the remaining valid saved products.
- **Empty saved view**: covered by Acceptance Scenario 4 — a clear, non-broken empty state rather than a blank page.
- **Session ends**: when the shopper's browsing session ends, their saved products are cleared (saved products are session-scoped, not persisted across sessions or devices — see Assumptions).
- **Many products saved**: the saved view remains usable as the number of saved products grows (no fixed limit is imposed in this scope — see Assumptions).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A shopper MUST be able to save a product from its product page without signing in or creating an account.
- **FR-002**: The system MUST show, on the product page, whether the product the shopper is viewing is currently saved.
- **FR-003**: A shopper MUST be able to open a saved-products view that lists every product they have saved during the current session.
- **FR-004**: The saved-products view MUST be reachable from a link in the site header that is present across the store.
- **FR-005**: From the saved-products view, a shopper MUST be able to navigate directly to the full product page of any saved product.
- **FR-006**: Saved products MUST persist for the shopper as they navigate around the store within the same browsing session.
- **FR-007**: Saving a product the shopper has already saved MUST NOT create a duplicate entry in the saved-products view.
- **FR-008**: When the shopper has no saved products, the saved-products view MUST display a message stating nothing is saved yet and how to start saving.

### Key Entities *(include if feature involves data)*

- **Saved product**: a reference, scoped to a single shopper's current session, indicating a product the shopper has chosen to save. Identifies which product it points to; ordering by when it was saved is sufficient for display.
- **Saved collection**: the set of saved products belonging to one shopper's session — the contents shown in the saved-products view and the basis for the saved/unsaved state on a product page.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can save a product from its page in a single action, with zero sign-in or account-creation steps.
- **SC-002**: 100% of products a shopper saves remain retrievable in the saved-products view for the remainder of that browsing session.
- **SC-003**: A shopper can reach their saved products from any page in the store in one click via the header.
- **SC-004**: From the saved-products view, a shopper can open any saved product's full detail page in one click.
- **SC-005**: A first-time shopper who opens an empty saved-products view can identify how to start saving products without external help (≥90% in usability testing).

## Assumptions

- **No accounts**: the store has no sign-in, so saving works anonymously; this is consistent with the epic's out-of-scope list.
- **Session scope**: "the same session" means the shopper's current browsing session using the store's existing session mechanism. Saved products persist across page navigation within that session but are not retained across a new session or on another device (cross-device is explicitly out of scope). The precise session window is left to implementation; the reasonable default is the lifetime of the existing browsing session.
- **No fixed limit**: no maximum number of saved products is imposed in this scope; a sensible bound is left to implementation.
- **Saving surface**: in this scope, saving is available from the **product detail page** only. Saving directly from the product list (AIP-193) is a separate story and out of scope here.
- **Out of scope (separate stories)**: removing a saved product (AIP-194) and showing a saved-product count in the header (AIP-195). This story provides only the header link, not a count.
- **Presentation**: shoppers return to saved products via a **dedicated saved-products view**; this resolves the PRD's open question on presentation for this story.
- **Catalogue dependency**: the saved-products view relies on the existing product catalogue to display saved products' details.
