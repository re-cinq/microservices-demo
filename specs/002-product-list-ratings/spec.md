# Feature Specification: Star ratings on the product list

**Feature Branch**: `002-product-list-ratings` (spec directory; work continues on git branch `attendee/daniel-tufvander`, no new branch created)
**Created**: 2026-05-20
**Status**: Draft
**Input**: Jira story [AIP-161](https://odevo.atlassian.net/browse/AIP-161) — *Star ratings appear on the product list* — under epic [AIP-95](https://odevo.atlassian.net/browse/AIP-95) *Product ratings*.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Star rating visible next to every product on the list (Priority: P1)

As a shopper browsing Online Boutique, I want to see a star rating next to each product on the product list so I can tell at a glance which products other customers like, without having to click into every product to compare.

**Why this priority**: This is the first independently shippable slice of the parent epic's value. The product list is the default landing surface for shoppers; surfacing social proof here turns "I can't tell what's worth a look" into "I can scan and prioritise". It delivers customer-visible value without depending on any other story in the epic.

**Independent Test**: A shopper opens the Online Boutique product list and, without clicking into any product, can see a star rating for every product card. Verification needs nothing beyond the running app on the cohort's environment.

**Acceptance Scenarios** (taken verbatim from AIP-161):

1. **Given** I open the product list, **When** the page finishes loading, **Then** every product card shows a star rating element alongside its existing details.
2. **Given** a product has a rating value in the product catalogue, **When** its card is rendered on the list, **Then** the rating displayed matches that value.
3. **Given** I refresh the product list, **When** I view the same products again, **Then** the same ratings are shown for the same products (ratings are stable across page loads).

### Edge Cases

- **Very long product lists / paginated lists**: ratings render for every product card regardless of position on the list; performance of the list page is not visibly worse than today.
- **Direct navigation to the product list URL (not arriving via search or category)**: ratings still render — the rating is a property of the product, not of the navigation path.
- **A product whose catalogue entry has no rating value** (see Assumption A4 — we expect all seeded products to have one): no rating element is rendered for that product card, and the rest of the card and the list layout are unaffected. (Out of scope for happy-path testing if A4 holds; included here as a defensive expectation.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Every product card on the product list MUST display a star rating element alongside the product's existing details.
- **FR-002**: The star rating shown on a product card MUST reflect the rating value associated with that product in the product catalogue.
- **FR-003**: The rating shown for a given product MUST be stable across page loads — repeated visits to the product list MUST show the same rating for the same product, unless its catalogue value has been changed.
- **FR-004**: The rating MUST be visible without any interaction (no hover, no click, no expand) — it is a default visual element of the product card.
- **FR-005**: The introduction of ratings MUST NOT alter the existing details shown on a product card; existing fields and layout MUST remain.

### Key Entities

- **Product**: an item in the Online Boutique catalogue. Existing attributes (name, price, image, description, etc.) remain unchanged. Gains one new attribute: a rating value.
- **Rating value**: a numeric value associated with a product, used to render the star display on the product list (and on the product page in a later story). Value range and granularity captured in Assumption A1.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of products visible on the product list page show a rating element on first load, with no shopper interaction required.
- **SC-002**: A shopper can identify the highest-rated and lowest-rated products visible on the list in a single visual scan, without opening any product page.
- **SC-003**: Repeated visits to the product list by the same shopper show identical ratings for the same products across page loads within the same session.
- **SC-004**: The change ships through the existing cohort pipeline with zero new services introduced, zero new datastores introduced, and zero changes to deployment manifests or CI configuration.

## Constraints *(mandatory — inherited from parent epic AIP-95, treated as hard constraints)*

These constraints come from the *Technical constraints* section of epic AIP-95 and are non-negotiable for this spec:

- **C-001**: Only services that already exist in this repository may be used. No new service may be introduced.
- **C-002**: No new datastore of any kind may be introduced — no database, no cache, no search engine. Rating data must live in memory over the existing product catalogue file (`productcatalogservice/products.json`).
- **C-003**: Changes must match the language and patterns of the service being changed. The frontend and the product catalogue service are written in Go; any code added to those services must remain Go and follow the existing conventions of each service.
- **C-004**: Infrastructure, deployment manifests, and CI configuration must not be changed. The feature must ship through the existing pipeline unmodified.

The parent epic also states *out of scope* for the wider Product ratings epic: written review text and moderation. Both remain out of scope for this story.

## Assumptions

The story description (AIP-161) and the epic (AIP-95) do not pin down every detail. The following are explicit assumptions made to fill the gaps; each can be challenged before planning.

- **A1 — Rating scale and granularity**: ratings are values on a 0-to-5 scale with at most half-star granularity (e.g. 0, 0.5, 1.0, … 5.0). This is the conventional star-rating scale and matches "star ratings" in the epic. If a different scale is required, this assumption needs revisiting before planning.
- **A2 — Source of the rating value**: each rating value is stored as a field on the product entry in `productcatalogservice/products.json` and seeded by hand into that file as part of this story. This is the only data location compatible with constraint C-002 (no new datastore).
- **A3 — Read-only display**: shoppers cannot submit, change, or remove ratings as part of this story. The story is display-only. The epic explicitly puts moderation out of scope, which only coheres with read-only ratings; user-submitted ratings would also conflict with constraint C-002 (no new datastore for submissions).
- **A4 — All products seeded**: every product entry in `productcatalogservice/products.json` will carry a rating value as part of this story's data change, so the "product with no rating" edge case is not expected to occur on the happy path.
- **A5 — Display content**: each product card on the list shows the star representation of the rating value. No review count, average count, or text label is shown alongside the stars in this story. (The epic talks about *star ratings* only; review count was an open question in the parent PRD and is intentionally not added here.)
- **A6 — Source-of-truth for both surfaces**: the same rating value will later be reused by the sibling story AIP-162 to render on the product page. The product-page rendering is explicitly out of scope for this story but the data added here must be usable by that story without further data work.

## Dependencies

- **Independently shippable.** AIP-161 does not depend on any other story. It is the first slice of the epic.
- **Downstream**: AIP-162 (*The same star rating appears on the product page*) depends on the catalogue-data and catalogue-service work landed by this story.

## Out of scope for this story

- Rendering ratings on the product page (sibling story AIP-162).
- Allowing shoppers to submit, edit, or remove ratings.
- Written review text and review moderation (out of scope for the whole epic per AIP-95).
- Showing a review count, total reviewer count, or any text caption next to the stars.
- Sorting or filtering the product list by rating.
- Any change to fulfilment, shipping, search, or recommendation logic.
- Any change to infrastructure, deployment manifests, or CI configuration (per constraint C-004).
