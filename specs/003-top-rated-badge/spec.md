# Feature Specification: Top rated badge on the highest-rated product card

**Feature Branch**: `003-top-rated-badge` (spec directory; work continues on git branch `attendee/daniel-tufvander`)
**Created**: 2026-05-20
**Status**: Draft
**Input**: Jira story [AIP-166](https://odevo.atlassian.net/browse/AIP-166) — *Top rated badge on the highest-rated product card* — under epic [AIP-95](https://odevo.atlassian.net/browse/AIP-95) *Product ratings*. Builds on the rating data shipped by [AIP-161](https://odevo.atlassian.net/browse/AIP-161).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - "Top rated" badge on the single highest-rated card (Priority: P1)

As a shopper landing on Online Boutique, I want a clear "Top rated" badge on the highest-rated product card so the single best-loved item stands out from the rest at a glance, beyond the star rating alone.

**Why this priority**: Only story this cycle; sole increment on the product-ratings epic after AIP-161. The full star widget already gives a shopper the rating *value* for each card; this story makes the single best-loved card *jump out* as the recommended starting point of the browse session. Cheap to ship, low risk, demoable in one screenshot.

**Independent Test**: A shopper opens the home page. Exactly one card carries a "Top rated" badge — the card whose star rating equals the maximum rating across all cards on the page. With the current seed (Watch = 5.0, everything else ≤ 4.5), only the Watch card carries the badge.

**Acceptance Scenarios** (taken verbatim from AIP-166):

1. **Given** I open the home page, **When** the cards finish loading, **Then** exactly one card carries a "Top rated" badge: the card whose rating equals the maximum rating across all visible products.
2. **Given** multiple products share the highest rating, **When** the cards render, **Then** all cards tied at that highest rating carry the "Top rated" badge.
3. **Given** no product on the page has a rating above 0.0 (defensive — not expected with the current seed), **When** the cards render, **Then** no card carries the "Top rated" badge.

### Edge Cases

- **Tie at the top**: a future seed could put two products both at 5.0. The badge appears on both. AS-2 covers this.
- **Single rated product, all others unrated (rating = 0.0)**: that one rated product is also the highest and gets the badge. Consistent with AS-1.
- **No rated products at all** (defensive): no badge anywhere. AS-3 covers this.
- **Catalogue size grows past today's 9 products**: the badge logic is "max rating" not "top N," so adding products does not change the behaviour.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The home page MUST visually mark the product card(s) whose rating equals the maximum rating among all rendered products with a "Top rated" badge.
- **FR-002**: When multiple cards tie at the maximum rating, every tied card MUST carry the badge — there is no tie-breaking heuristic that picks one over another.
- **FR-003**: When no rendered product has a rating greater than `0.0`, no card MUST carry the badge.
- **FR-004**: The badge MUST be readable without hovering, clicking, or expanding — it is a default visual element of the card.
- **FR-005**: The presence of the badge MUST NOT alter or hide any other existing element on the card (name, price, star widget).
- **FR-006**: The "Top rated" determination MUST be stable across page loads for the same product set and the same rating values.

### Key Entities

No new entity. The badge is a derived attribute of an existing product: `TopRated = (rating == max(rating across rendered products))`. It is computed at render time and not persisted.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of home-page renders, when at least one product has a rating > 0.0, show the "Top rated" badge on the card(s) tied at the maximum rating, and on no other card.
- **SC-002**: A first-time shopper can point at the highest-rated product on the home page within 2 seconds of the page loading, without comparing star widgets across the grid.
- **SC-003**: Repeated reloads of the home page by the same shopper show the badge on the same card(s) for the same rating values — no flicker or drift.
- **SC-004**: The change ships through the existing cohort pipeline with zero new services, zero new datastores, zero changes to deployment manifests or CI configuration, and zero changes to the gRPC contract (no proto regen needed).

## Constraints *(mandatory — inherited from parent epic AIP-95, treated as hard constraints)*

These constraints come from the *Technical constraints* section of epic AIP-95 and are non-negotiable for this spec:

- **C-001**: Only services that already exist in this repository may be used. No new service may be introduced.
- **C-002**: No new datastore of any kind. The badge is computed from the in-memory rating already loaded by the catalogue service; no persistence is added.
- **C-003**: Changes must match the language and patterns of the service being changed. The frontend is Go (`html/template`); any code added must remain Go and follow existing conventions.
- **C-004**: Infrastructure, deployment manifests, and CI configuration must not be changed.

Story-specific additions:

- **C-005**: No proto change. The rating field added by AIP-161 is sufficient; the badge is purely a frontend-derived value.
- **C-006**: No backend change. `productcatalogservice` is not touched by this story.

## Assumptions

- **A1 — "Maximum" is a strict numeric maximum** of the `Rating float32` value already exposed on each `Product` (added by AIP-161). No fuzzy "near-top" bucketing.
- **A2 — Scope is the home page only**. The product detail page, cart, search results, and any other surface are out of scope for this story (consistent with AIP-95's separation of list vs. page work).
- **A3 — Badge copy is the literal English string "Top rated"** in title case. No i18n / translation in this story; the rest of the Online Boutique frontend is not internationalised either.
- **A4 — Badge is a visual element only**, not a link or a filter affordance. Clicking it does nothing beyond the standard click target of the underlying card image (the existing `<a>`).
- **A5 — Badge visual placement** is within the existing `.hot-product-card` container so it does not affect grid layout. Exact pixel position is the implementer's call; the contract is "visible without interaction, does not overflow the card."

## Dependencies

- **Upstream (satisfied)**: [AIP-161](https://odevo.atlassian.net/browse/AIP-161) — landed the `Rating` field on `productView` in `handlers.go` and the star widget on the home page. This story builds directly on that work.
- **Downstream**: none.

## Out of scope for this story

- Any change to the product detail page, the cart, the search results, or any other surface (sibling story [AIP-162](https://odevo.atlassian.net/browse/AIP-162) — deferred — covers the product detail page).
- Sorting / filtering the product list by rating.
- Showing a numeric value or review count alongside the badge.
- A "Top 3" or "Top N" highlight pattern — this story is strictly "tied at max."
- Internationalisation of the badge copy.
- Any change to `productcatalogservice`, the proto contract, or the seed data.
- Any change to deployment manifests, CI configuration, or infrastructure (per C-004).
