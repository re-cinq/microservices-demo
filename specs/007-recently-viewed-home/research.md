# Research: Recently Viewed Strip on the Home Page

**Feature**: AIP-157 | **Date**: 2026-05-20

No external research required. All decisions are determined by the existing codebase (AIP-156).

## Decisions

### Recently-viewed data source

**Decision**: Read from `shop_recently-viewed` cookie via existing `recentlyViewedFromCookie(r *http.Request) []string` helper.  
**Rationale**: Established in AIP-156; consistent with `shop_wishlist` cookie pattern from AIP-159.  
**Alternatives considered**: Server-side session store — rejected; not used anywhere in the service.

### Product fetch

**Decision**: Use existing `fe.getRecentlyViewedProducts(ctx, ids)` method.  
**Rationale**: Already called in `productHandler`; fetches from product catalog gRPC service with the same deduplication and ordering semantics.  
**Alternatives considered**: None — function exists specifically for this purpose.

### Currency conversion

**Decision**: Call `fe.convertCurrency` for each recently-viewed product, same as the hot-products loop in `homeHandler`.  
**Rationale**: All product prices on the page must respect the user's selected currency. Pattern established in every handler that displays prices.  
**Alternatives considered**: N/A.

### Error handling

**Decision**: Log a warning and render the page without the strip if `getRecentlyViewedProducts` fails.  
**Rationale**: Matches the non-fatal pattern used for ads and recommendations — a degraded page is better than an error page.  
**Alternatives considered**: Hard error — rejected; recently-viewed is a non-critical enhancement.

### Template rendering

**Decision**: Use existing `recently_viewed` sub-template (`recently_viewed.html`) unchanged.  
**Rationale**: Already renders correctly on the product detail page; no design changes requested.  
**Alternatives considered**: Inline HTML — rejected; DRY principle, style parity guaranteed.
