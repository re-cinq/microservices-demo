# Implementation Plan: Recently Viewed Strip on the Home Page

**Branch**: `007-recently-viewed-home` | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/007-recently-viewed-home/spec.md`

## Summary

Add the existing recently-viewed strip component (established in AIP-156) to the home page. Two file changes: `homeHandler` in `handlers.go` fetches recently-viewed products from cookie and passes them to the template; `home.html` conditionally renders the `recently_viewed` sub-template. No new infrastructure, dependencies, or data model changes required.

## Technical Context

**Language/Version**: Go 1.21  
**Primary Dependencies**: gorilla/mux (routing), html/template (rendering) — no new dependencies  
**Storage**: Browser cookie `shop_recently-viewed` (pipe-delimited product IDs, 48h MaxAge) — no server-side storage  
**Testing**: Manual integration test via running frontend service  
**Target Platform**: Linux container (GKE / local Docker)  
**Project Type**: Web service (frontend microservice in microservices-demo)  
**Performance Goals**: No additional RPC calls introduced beyond existing `getRecentlyViewedProducts` pattern  
**Constraints**: No JavaScript; no new template files; must use existing `recently_viewed` sub-template unchanged  
**Scale/Scope**: Single handler change + single template change

## Constitution Check

Constitution file contains only a blank template — no binding constraints to evaluate. No gates to check.

## Project Structure

### Documentation (this feature)

```text
specs/007-recently-viewed-home/
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal — no unknowns)
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

No `data-model.md` (no new entities) and no `contracts/` (no new endpoints or API surfaces).

### Source Code (changed files only)

```text
src/frontend/
├── handlers.go          # homeHandler: add recently-viewed fetch + template key
└── templates/
    └── home.html        # Add {{ if $.recently_viewed }}{{ template "recently_viewed" $ }}{{ end }}
```

**Structure Decision**: Single project (existing `src/frontend/`). All patterns already established — this is a wiring change only.

## Implementation Approach

### handlers.go — homeHandler

Inside `homeHandler`, after building `ps` (the product list), read the recently-viewed cookie, fetch the products, and pass them to the template:

```go
recentlyViewedIDs := recentlyViewedFromCookie(r)
var recentlyViewed []productView
if len(recentlyViewedIDs) > 0 {
    rvProducts, err := fe.getRecentlyViewedProducts(r.Context(), recentlyViewedIDs)
    if err != nil {
        log.WithError(err).Warn("could not retrieve recently viewed products for home page")
        // non-fatal: continue without the strip
    } else {
        for _, p := range rvProducts {
            price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
            if err == nil {
                recentlyViewed = append(recentlyViewed, productView{p, price})
            }
        }
    }
}
```

Add `"recently_viewed": recentlyViewed` to the `ExecuteTemplate` data map.

### home.html

After the hot-products section (before the footer blocks), add:

```html
<div>
  {{ if $.recently_viewed }}
    {{ template "recently_viewed" $ }}
  {{ end }}
</div>
```

This mirrors the placement pattern in `product.html`.

## Complexity Tracking

No constitution violations. No complexity justification needed.
