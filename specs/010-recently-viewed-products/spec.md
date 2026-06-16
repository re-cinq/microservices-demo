# Feature Specification: Recently Viewed Products (Product Page)

**Feature Branch**: `010-recently-viewed-products`

**Created**: 2026-06-16

**Status**: Draft

**Source**: Jira story [AIP-199](https://odevo.atlassian.net/browse/AIP-199) — "See and return to recently viewed products", under epic [AIP-97](https://odevo.atlassian.net/browse/AIP-97) — "Recently viewed".

**Input**: Story description and acceptance criteria from AIP-199, with the technical-constraints section of epic AIP-97 treated as hard constraints.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See and return to recently viewed products (Priority: P1)

As a shopper, while I'm browsing, I see a strip of products I've already looked at during this visit, and I can click any of them to jump straight back to its product page — without searching again. Before I've viewed anything, the page looks exactly as it does today.

**Why this priority**: This is the entire value of the story — letting a shopper return to a product they already saw without re-searching. It is independently shippable and customer-visible on its own.

**Independent Test**: View two or more products, then open a product page and confirm a "Recently viewed" strip appears at the bottom listing the others; click one and confirm it opens that product's page. Open a product page in a fresh session and confirm no strip appears and the page is unchanged.

**Acceptance Scenarios**:

1. **Given** I have viewed one or more products this session, **When** I'm on a product page, **Then** a "Recently viewed" strip lists those products at the bottom of the page, most recently viewed first, and the product I'm currently on is not listed.
2. **Given** I have viewed more than 4 products, **When** the strip renders, **Then** only the 4 most recently viewed are shown.
3. **Given** the strip is shown, **When** I select a product in it, **Then** I land on that product's page.
4. **Given** I have not viewed any product this session, **When** the product page renders, **Then** no strip is shown and the layout is unchanged.
5. **Given** the strip is shown, **When** it renders, **Then** its thumbnails are half (50%) the size of the other product thumbnails on the page.

---

### Edge Cases

- **Viewing the same product more than once**: De-duplication / move-to-front is **out of scope** for this story (handled separately in [AIP-201](https://odevo.atlassian.net/browse/AIP-201)). Within this story, repeated views may appear as repeated entries among the 4 shown.
- **Exactly 4 vs. more than 4 viewed**: With exactly 4 prior views, all 4 show; with more, only the 4 most recent show.
- **A previously viewed product is no longer in the catalogue**: The strip MUST still render the remaining valid products and MUST NOT error or block the page.
- **First product of the session**: After viewing one product and navigating to a second, the first appears in the strip; on the very first product page of a session the strip is absent (nothing yet to show).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST record each product a shopper views during their current session.
- **FR-002**: System MUST display a "Recently viewed" strip at the bottom of a product page whenever the shopper has at least one other recently viewed product to show.
- **FR-003**: The strip MUST list products most-recently-viewed first.
- **FR-004**: The strip MUST exclude the product the shopper is currently viewing.
- **FR-005**: The strip MUST show at most the 4 most recently viewed products.
- **FR-006**: Selecting a product in the strip MUST navigate the shopper to that product's page.
- **FR-007**: When the shopper has viewed no products this session, the product page MUST render with no strip and with its layout unchanged from today.
- **FR-008**: The strip's thumbnails MUST be 50% the size of the other product thumbnails on the page where it appears.
- **FR-009**: The recently viewed list MUST be maintained in memory for the duration of the session only, with no new persistent store introduced.
- **FR-010**: The strip MUST render gracefully (showing remaining valid products, no error) if a recorded product can no longer be retrieved from the catalogue.

### Key Entities *(include if feature involves data)*

- **Recently Viewed List**: An ordered, in-memory collection of product references for a single shopper session, ordered most-recent-first, of which at most 4 are displayed.
- **Product**: An existing catalogue item, referenced by its identifier, with a name, price, and image/thumbnail used to render the strip.
- **Session**: The shopper's current browsing session, which scopes the recently viewed list.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper who has viewed at least one prior product can return to it in a single click from the strip, with zero searches.
- **SC-002**: The strip never displays more than 4 products, regardless of how many were viewed.
- **SC-003**: For a shopper who has viewed at least two distinct products, the strip appears on 100% of subsequent product page loads in the session.
- **SC-004**: When no products have been viewed, the product page is visually identical to today (no added strip, no layout shift).
- **SC-005**: The current product never appears in its own strip.

## Constraints *(hard — inherited from epic AIP-97)*

These are non-negotiable constraints carried from the epic and bound any implementation of this feature:

- **C-001**: Use only services that already exist in this repository. Do not add new services.
- **C-002**: Do not introduce any new datastore (no database, no cache, no search engine). Work in memory over the existing product catalogue in `productcatalogservice/products.json`.
- **C-003**: Match the language and patterns of the service being changed. The frontend and the product catalogue service are written in Go.
- **C-004**: Do not change infrastructure, deployment manifests, or CI configuration. The change must ship through the existing pipeline unmodified.

## Assumptions

- **"Viewed" means loading a product's page.** The epic leaves the trigger open; opening a product page is taken as the concrete definition.
- **Session identity reuses the existing frontend session mechanism**, and the recently viewed list is scoped to that session and held in memory only (it does not survive a service restart). This is consistent with C-002 and with the "persistence/scope" question still open at the PRD level.
- **Ordering is most-recently-viewed first.** The PRD lists ordering as open, but the story's acceptance criteria specify most-recent-first, which is adopted here.
- **Anonymous and signed-in shoppers are treated the same** for this story. (The PRD lists this as an open question; same behaviour is the default until decided.)
- **De-duplication is out of scope** for this story and is covered by AIP-201; this story may show duplicate entries within the 4 displayed.
- **Cross-device sync is out of scope** (per the epic).
