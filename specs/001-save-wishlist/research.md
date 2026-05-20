# Research: Save for Later / Wishlist

**Feature**: 001-save-wishlist
**Date**: 2026-05-20

## Decision 1: In-memory storage mechanism

**Decision**: `sync.Map` on the `frontendServer` struct, keyed by session ID (`string`) → `[]string` of product IDs.

**Rationale**: `sync.Map` is in the Go standard library, requires no new dependencies, and is safe for concurrent access without explicit locking. The wishlist workload is read-heavy (displayed on every wishlist page load) and write-light (saved infrequently), which fits `sync.Map`'s optimisation profile. Adding the field directly to `frontendServer` follows the existing pattern for service state (all gRPC connections are already fields on that struct).

**Alternatives considered**:
- `map[string][]string` with `sync.RWMutex`: functionally equivalent but more boilerplate and error-prone (lock/unlock discipline).
- External cache (Redis, memcached): excluded by hard constraint — no new datastores.
- Cookie-stored product IDs: would avoid server-side state entirely, but cookie size limits and lack of server validation make it fragile.

**Memory lifecycle**: Session IDs are cookie-scoped with a 48-hour max-age (`cookieMaxAge = 60 * 60 * 48` in `main.go`). Orphaned map entries after session expiry accumulate but are never accessed again. For a demo-scale deployment this is acceptable without a background purge goroutine.

---

## Decision 2: Save confirmation pattern

**Decision**: POST-redirect-GET pattern — `POST /wishlist/save` redirects back to the product page with a `?saved=1` query parameter; the product template renders an inline confirmation banner when the parameter is present.

**Rationale**: Matches the existing add-to-cart pattern exactly (`addToCartHandler` also does a POST + `http.Redirect`). No JavaScript required. The Go `html/template` system already supports conditional rendering via `{{ if $.saved }}`. Avoids double-submit on browser refresh (standard PRG benefit).

**Alternatives considered**:
- Fetch/HTMX-based partial update: better UX but introduces a JS dependency and departs from the project's no-JS-framework approach.
- Render the confirmation in the POST response directly: causes double-submit on refresh.

---

## Decision 3: Wishlist navigation entry point

**Decision**: `GET /wishlist` route with a bookmark icon link added to the header navigation, placed next to the existing cart icon in `header.html`.

**Rationale**: The cart link in `header.html` provides an exact pattern to follow — an `<a>` tag with an icon image, a `cart-link` CSS class, and a `{{ $.baseUrl }}/cart` href. The wishlist link mirrors this with `{{ $.baseUrl }}/wishlist`. A count badge (like the cart's `cart_size` circle) showing the number of saved items follows naturally from the same pattern.

**Alternatives considered**:
- Dedicated nav section or dropdown: over-engineered for a single link.
- Link only from the product page confirmation: not globally accessible — shopper can't return to wishlist without re-saving.

---

## Decision 4: Duplicate prevention

**Decision**: Before appending a product ID to the session's wishlist slice, check whether it is already present. If it is, skip the append and still redirect with `?saved=1` (idempotent save).

**Rationale**: Spec FR-007 requires no duplicates. A linear scan over the slice is sufficient — a typical wishlist will hold fewer than 20 items (SC-002), so O(n) is negligible.

---

## Decision 5: Product detail resolution at display time

**Decision**: The wishlist stores only product IDs. When rendering `GET /wishlist`, the handler fetches each product's details (name, picture, price) from the existing `productCatalogSvcConn` gRPC connection, using the same `getProduct` helper already used by `productHandler`.

**Rationale**: Keeps the in-memory store minimal (just IDs). Ensures price and availability shown are always current (satisfies spec Assumption 3). Reuses existing gRPC infrastructure with no new service calls or dependencies.

**Alternatives considered**:
- Cache full product structs in the wishlist map: stale price risk; higher memory use; unnecessary for demo scale.
