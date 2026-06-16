# Feature Specification: Rate a product and see it on the product page

**Feature Branch**: `010-product-rating-submission`  
**Created**: 2026-06-16  
**Status**: Draft  
**Input**: Jira story AIP-196 (parent epic AIP-95). Story description and acceptance criteria used as the spec; the epic's technical-constraints section is treated as hard constraints.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Submit a star rating from the product page (Priority: P1)

A shopper viewing a product wants to give it a 1–5 star rating and immediately see the product's overall score, so that quality reflects real shopper input rather than guesswork.

**Why this priority**: This is the whole feature for a single product surface, end-to-end: it creates the rating data (there is none today) and makes it visible. Without it there is nothing to display anywhere, so every other rating story depends on it. On its own it is a shippable, customer-visible slice.

**Independent Test**: Open a product page, select a star value, submit, and confirm the displayed average and rating count update to include the submission — no other story required.

**Acceptance Scenarios**:

1. **Given** a product page, **When** I select 1–5 stars and submit, **Then** my rating is recorded and the displayed average and count update to include it.
2. **Given** a product with at least one rating, **When** I view its page, **Then** I see the average shown as stars with the number of ratings in parentheses (e.g. ★★★★☆ (43)).
3. **Given** the product catalogue service restarts, **When** I view a product, **Then** ratings have reset to empty, because ratings are held in memory with no datastore (accepted limitation).

---

### Edge Cases

- **Out-of-range or empty selection**: A submission outside 1–5 stars (or with no star selected) MUST be rejected and MUST NOT change the recorded average or count.
- **First rating on a product**: The first valid submission moves the product from "no ratings" to an average equal to that single rating with a count of 1.
- **Fractional average display**: When the average is not a whole number, how it renders against whole-star icons is unresolved — see [NEEDS CLARIFICATION] in FR-006.
- **Repeat submissions**: The same shopper can submit multiple times and each counts; there is no dedup, per the epic (moderation out of scope).
- **Concurrent submissions**: Two ratings submitted at the same time MUST both be counted, with no lost update to the running total.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: A shopper MUST be able to select a whole-number rating from 1 to 5 stars on the product page and submit it.
- **FR-002**: The system MUST record each valid submission and include it in that product's running average and rating count.
- **FR-003**: The system MUST reject a submission whose value is outside 1–5 (or is absent) without altering the product's average or count.
- **FR-004**: After a successful submission, the product page MUST show the updated average and count without the shopper needing to re-find the product.
- **FR-005**: The product page MUST display the average as stars followed by the number of ratings in parentheses (e.g. ★★★★☆ (43)) for any product that has at least one rating.
- **FR-006**: The system MUST render the average rating to the shopper as stars. [NEEDS CLARIFICATION: how should a fractional average render — round to nearest whole star, show half-stars, or show stars plus a numeric average?]
- **FR-007**: Rating data MUST be derived solely from shopper submissions; the system MUST NOT use seeded or placeholder rating values.

### Hard Constraints *(from epic AIP-95 — non-negotiable)*

- **HC-001**: MUST use only services that already exist in the repository. No new service may be added.
- **HC-002**: MUST NOT introduce any new datastore (no database, no cache, no search engine). Rating data is held in memory over the existing product catalogue.
- **HC-003**: MUST match the language and patterns of the service being changed (the frontend and product catalogue service are written in Go).
- **HC-004**: MUST NOT change infrastructure, deployment manifests, or CI configuration. The change must ship through the existing pipeline unmodified.
- **HC-005** (epic scope): Written review text and moderation are out of scope.

### Key Entities *(include if feature involves data)*

- **Product Rating**: The aggregate rating state for one product — its running total/average of submitted star values and the count of submissions. Held in memory; not persisted. A product with no submissions has no rating.
- **Rating Submission**: A single shopper action contributing one whole-number star value (1–5) to a product's aggregate. Not individually stored or attributed beyond its effect on the aggregate.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can submit a rating and see the product's average and count reflect it within the same product-page visit (no manual refresh or navigation away required).
- **SC-002**: 100% of submissions outside the 1–5 range leave the product's average and count unchanged.
- **SC-003**: Every product that has received at least one rating displays both an average (as stars) and a rating count on its page.
- **SC-004**: No product displays a rating value that did not originate from a shopper submission (zero seeded data).

## Assumptions

- **Display-only beyond submission**: Beyond submitting and viewing, there is no rating history, editing, or deletion in scope.
- **No identity/auth on rating**: Submissions are anonymous; there is no per-user vote limit or dedup, consistent with moderation being out of scope (HC-005).
- **In-memory persistence accepted**: Ratings resetting on service restart, and potentially differing across multiple service replicas, is an accepted limitation for this feature (HC-002). When ratings should move to a persistent datastore is a known open question carried in the PRD, not resolved here.
- **Single surface in this story**: This spec covers only the product page. Showing ratings on the product list view and the "no ratings yet" empty state are separate stories (AIP-197, AIP-198).
- **Star scale is 1–5**: A five-star whole-number scale is assumed from the acceptance criteria.
