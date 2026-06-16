# Feature Specification: Product Star Rating — Detail Page

**Jira Epic**: [AIP-95](https://odevo.atlassian.net/browse/AIP-95) — Product ratings
**Jira Stories**: [AIP-181](https://odevo.atlassian.net/browse/AIP-181) (seed ratings) · [AIP-182](https://odevo.atlassian.net/browse/AIP-182) (display on detail page)

**Branch**: `feat/aip-182-product-rating-detail`

**Created**: 2026-06-16

**Status**: Draft

---

## User Scenarios & Testing

### User Story 1 — Seed star ratings into the product catalogue (Priority: P1)

Every product in the catalogue carries a numeric star rating (1–5) so that downstream display stories have real data to work with. Shoppers benefit indirectly: any rating they see reflects a real value rather than a placeholder.

**Why this priority**: Nothing in AIP-182 can be displayed without ratings existing in the data. This is the foundational prerequisite.

**Independent Test**: Fetch every product from the product catalogue service and confirm each one returns a numeric rating between 1 and 5, with no product returning null or missing the field.

**Acceptance Scenarios**:

1. **Given** the product catalogue is loaded, **when** any product is retrieved, **then** it includes a numeric rating between 1 and 5.
2. **Given** the catalogue has multiple products, **when** all products are fetched, **then** every product has a rating value — no product returns null or a missing field.
3. **Given** the service starts up, **when** it reads `products.json`, **then** no new datastore, database, or cache is introduced — ratings live in memory only.

---

### User Story 2 — Display star rating on the product detail page (Priority: P2)

As a shopper viewing a single product, I can see its star rating on the product page so I can judge quality before adding it to my cart.

**Why this priority**: Depends on US1 (ratings must exist before they can be shown). Delivers the visible user value of the feature.

**Independent Test**: Navigate to any product detail page and confirm the star rating is visibly displayed, matches the product's seeded rating value, and no price, review text, or submission UI is present.

**Acceptance Scenarios**:

1. **Given** I navigate to a product detail page, **when** the page loads, **then** the product's star rating is displayed visibly on the page.
2. **Given** a product has a rating of N stars, **when** I view its detail page, **then** the displayed rating matches that value.
3. **Given** the page renders, **when** I inspect it, **then** no price, review text, or user-submission element is present alongside the rating.

---

### Edge Cases

- What if a product has no rating in `products.json`? It must not be possible — US1 requires all products to carry a rating before US2 can display one.
- What if the rating is a non-integer (e.g., 3.7)? The display should handle fractional values gracefully (e.g., show one decimal place or round to nearest half-star).

---

## Requirements

### Functional Requirements

- **FR-001**: Every product in the catalogue MUST carry a numeric star rating between 1 and 5 inclusive.
- **FR-002**: Ratings MUST be stored in `productcatalogservice/products.json` and loaded into memory at service start — no new datastore may be introduced.
- **FR-003**: The product detail page MUST display the star rating visibly when the page loads.
- **FR-004**: The displayed rating MUST match the value stored in the catalogue for that product.
- **FR-005**: The rating display MUST NOT include price information, written review text, or any user-submission element (e.g., a form or input to submit a rating).

### Hard Constraints (from AIP-95)

- **C-001**: Use only the services that already exist in the repo. Do NOT add new services.
- **C-002**: Do NOT introduce any new datastore (no database, no cache, no search engine). Ratings live in memory over `productcatalogservice/products.json`.
- **C-003**: Match the language and patterns of the service being changed. The frontend and product catalogue service are written in Go.
- **C-004**: Do NOT change infrastructure, deployment manifests, or CI configuration. The change must ship through the existing pipeline unmodified.

### Key Entities

- **Product rating**: A numeric value (float, 1–5) associated with each product, stored as a field on the product record in `products.json` and carried through the product catalogue service to the frontend.

---

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every product in the catalogue returns a rating value between 1 and 5 — verifiable by inspecting all products via the catalogue API.
- **SC-002**: A shopper can see the star rating on any product detail page without any additional navigation or interaction.
- **SC-003**: The displayed rating on the product page matches the value in `products.json` — verifiable by cross-checking each product manually.
- **SC-004**: No new service, database, cache, or infrastructure component is introduced — verifiable by diffing deployment manifests and service directories.

---

## Assumptions

- Ratings are static seed values — there is no user-facing submission flow in this story (that is explicitly out of scope per AIP-95).
- The `Product` data structure passed from the product catalogue service to the frontend must be extended to carry the rating field; this requires changes to both the proto definition and generated Go code in both services.
- A reasonable set of seed ratings (covering all products in `products.json`) is sufficient; they do not need to reflect real-world opinions.
- Fractional ratings are permitted (e.g., 4.3); the template should render them to one decimal place.
- The star rating is displayed on the product detail page only — the product list page is out of scope for this story (a separate story may cover that).
