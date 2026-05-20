# Quickstart: Recently Viewed Strip on the Home Page

**Feature**: AIP-157 | **Date**: 2026-05-20

## Prerequisites

- Frontend service running locally or in a dev cluster
- At least one product in the catalog

## Scenarios

### Scenario 1: Strip appears after viewing a product

1. Open the home page — confirm no recently viewed strip is visible.
2. Click any product to open its detail page.
3. Navigate back to the home page (click the logo or use the browser back button).
4. **Expected**: A "Recently Viewed" strip appears below the hot-products section showing the product you just viewed, with its image, name, and price.

### Scenario 2: Multiple products appear in the strip

1. View three different products in sequence.
2. Return to the home page.
3. **Expected**: The recently viewed strip shows all three products in reverse-visit order (most recent first), each linking to its detail page.

### Scenario 3: Strip is absent on a fresh session

1. Clear browser cookies (or open a private/incognito window).
2. Visit the home page without viewing any products.
3. **Expected**: No recently viewed strip is shown; the hot-products section and footer render normally.

### Scenario 4: Strip links navigate correctly

1. View a product, return to the home page.
2. Click the product's image in the recently viewed strip.
3. **Expected**: Browser navigates to that product's detail page (`/product/{id}`).

### Scenario 5: Currency selection is respected

1. View a product, then change currency via the header selector.
2. Return to the home page.
3. **Expected**: Prices in the recently viewed strip reflect the selected currency.
