# Spec: Star ratings on the product list (AIP-154)

## Overview

Shoppers browsing Online Boutique currently see only a product name and description on the product list page — there is no quality signal to help them compare products without clicking into each one. This story adds a star rating to every product card on the list, reading rating values from the existing product catalogue (`productcatalogservice/products.json`). It builds directly on AIP-153, which landed the rating field in `products.json` and introduced the star display component on the product detail page.

**Jira:** [AIP-154](https://odevo.atlassian.net/browse/AIP-154)
**Epic:** [AIP-95 — Product ratings](https://odevo.atlassian.net/browse/AIP-95)
**PRD:** [PRD: Product ratings (nejal-patel)](https://odevo.atlassian.net/wiki/spaces/APT/pages/2358509581/PRD+Product+ratings+nejal-patel)
**Depends on:** AIP-153 ✓ Done (PR #29)

---

## Functional requirements

### FR-1 — Rating displayed on every product card
Every product card rendered on the product list page must include a star rating element. The rating must be visible without any user interaction (no hover, no expand).

### FR-2 — Rating value matches the product catalogue
The rating value shown on a product card must exactly match the value stored for that product in `productcatalogservice/products.json`. If the same product is viewed on both the list and the detail page, both surfaces must show the same value.

### FR-3 — Reuse the star display component from AIP-153
The star component introduced on the product detail page (AIP-153) must be reused on the product list. No second implementation of star rendering is introduced.

### FR-4 — No-rating placeholder
If a product has no rating value in `products.json` (field absent or zero), the card must show a placeholder star display and the label "no reviews" — consistent with the detail page behaviour defined in AIP-155. It must not show a broken widget, a blank space, or a misleading zero.

### FR-5 — Layout stability
Adding the star rating to a product card must not break or reflow the existing card layout. All other card content (image, name, price) must remain visually unchanged.

### FR-6 — No new services or datastores
Ratings are served entirely from the existing in-memory product catalogue. No new backend endpoint, database, cache, or search engine is introduced. No changes to infrastructure, deployment manifests, or CI configuration.

---

## User scenarios

### Scenario 1 — Shopper sees ratings while browsing
**Given** a shopper opens the product list page,
**when** the page finishes loading,
**then** every product card displays a star rating element alongside its existing details.

### Scenario 2 — Rating value is consistent across surfaces
**Given** a product has a rating value in `products.json`,
**when** the shopper views that product on the list and then clicks through to its detail page,
**then** the star rating shown on the card is identical to the rating shown on the detail page.

### Scenario 3 — Ratings are stable across page loads
**Given** a shopper views the product list,
**when** they refresh the page,
**then** the same star ratings are shown for the same products.

### Scenario 4 — Unrated product shows placeholder
**Given** a product has no rating value in `products.json`,
**when** a shopper views that product's card on the list,
**then** a placeholder star display and the label "no reviews" are shown in place of a live rating.

### Scenario 5 — Card layout is preserved
**Given** a shopper views the product list with ratings displayed,
**then** the existing card layout (image, name, price) is not broken, truncated, or reflowed.

---

## Success criteria

| # | Criterion | Measurable test |
|---|-----------|----------------|
| SC-1 | Every product card on the list includes a star rating element | Page snapshot shows a rating element on each card |
| SC-2 | Rating on card matches rating in `products.json` | Value rendered in HTML equals the value in the JSON for that product ID |
| SC-3 | Rating on card matches rating on detail page | Both surfaces render the same numeric value for the same product |
| SC-4 | Unrated products show placeholder + "no reviews" | Card for a product with no rating field shows placeholder stars and "no reviews" text |
| SC-5 | No layout regression | Screenshot diff of product list shows no reflow of existing card content |

---

## Constraints

- Language and patterns must match the existing frontend (Go templates).
- The star component from AIP-153 must be reused — not duplicated.
- No new services, datastores, or infrastructure changes.
- Must ship through the existing CI pipeline unmodified.

---

## Out of scope

- Written review text and moderation.
- Vote counts alongside the star score.
- Any new backend endpoint or data pipeline.
- Changes to the product detail page (covered by AIP-153, done).
- The no-rating empty state on the detail page (covered by AIP-155).
