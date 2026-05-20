# Research: Recently Viewed Products Strip

**Feature**: `specs/001-recently-viewed-strip`
**Date**: 2026-05-20
**Status**: Complete — all unknowns resolved

---

## 1. View-History Storage

**Decision**: Browser cookie named `shop_recently_viewed`, value is a pipe-separated (`|`) ordered list of product IDs.

**Rationale**:
- Constraint C-002 prohibits any new datastore (database, cache, search engine).
- The frontend is already stateless — session identity flows through the `shop_session-id` cookie and cart state lives in `cartservice`. Storing a short list of IDs in a cookie is consistent with this architecture.
- Pipe `|` is used as the separator (rather than comma or JSON) to avoid URL-encoding overhead and keep cookie parsing trivial with `strings.Split` / `strings.Join`.
- Max 5 IDs keeps the cookie payload under 200 bytes (product IDs are 10-character alphanumeric strings).

**Alternatives considered**:
- *In-process map keyed by session ID*: rejected — the frontend runs as multiple replicas; in-memory state would not survive pod restarts or be consistent across replicas.
- *JSON encoding*: rejected — adds unnecessary overhead for a flat list of strings; Go standard library handles pipe-split without an import.
- *New Redis key*: rejected — violates C-002 (new datastore) and would require a new service dependency.

---

## 2. Cookie Lifecycle

**Decision**: Cookie `shop_recently_viewed` has a 24-hour `MaxAge`, `HttpOnly: false` (read by client-side is not required but not harmful), `SameSite: Lax` (consistent with existing cookies).

**Rationale**:
- 24 hours provides a useful session-spanning window (a shopper may return the next day) without indefinite persistence.
- The cookie is written/updated on every `productHandler` request, so its expiry is naturally refreshed on each product view.
- Spec scopes history to the current session; 24-hour expiry is a pragmatic implementation of that intent.

**Alternatives considered**:
- *Session-only cookie (no MaxAge)*: acceptable, but would clear history on browser close, which is stricter than the spec requires.
- *7-day expiry*: would match recommendations patterns but is longer than "session" implies.

---

## 3. Product Detail Fetching for Strip

**Decision**: Use the existing `fe.getProduct(ctx, id)` helper in `rpc.go` in a sequential loop over the cookie IDs.

**Rationale**:
- Constraint C-001 prohibits new services; C-003 requires matching existing patterns.
- `getProduct` is a thin wrapper over a unary gRPC call to productcatalogservice, which serves from an in-memory map. Round-trip latency is negligible (<5 ms per call on the same cluster).
- Sequential calls for 5 products adds <25 ms — well within acceptable page render budgets.
- A parallel `errgroup` fetch could reduce latency further but adds complexity that is not warranted for 5 calls.

**Alternatives considered**:
- *New `GetProducts` bulk RPC*: rejected — would require changes to productcatalogservice proto and server, violating C-001 (adding new capability to an existing service goes beyond the constraint's spirit) and requiring proto regeneration.
- *Parallel goroutine fetch with `errgroup`*: valid future optimisation; deferred as premature.

---

## 4. Currency Conversion

**Decision**: Apply the same `fe.convertCurrency` call used in `homeHandler` and `productHandler` to each product in the recently viewed list.

**Rationale**:
- The strip displays price; price must respect `currentCurrency(r)` to be consistent with the rest of the page.
- `convertCurrency` is already called per-product in the home page loop — the same pattern is directly reusable.

---

## 5. Where to Render the Strip

**Decision**: Render the strip on the **product detail page** (`/product/{id}`) only for the MVP. It is included as a partial template (`recently_viewed.html`) injected at the bottom of `product.html`, above the recommendations strip if both are present.

**Rationale**:
- AIP-163 does not specify which pages carry the strip; product detail pages are where return-navigation is most valuable (the shopper just viewed a product and may want to compare with one they saw earlier).
- The home page is a natural candidate for a later iteration but adds template changes beyond MVP scope.
- Using a partial keeps the template change isolated and makes it straightforward to add to other pages later.

---

## 6. Deduplication and Ordering

**Decision**: When recording a new product view, remove any existing occurrence of that product ID from the list, prepend the new ID, then truncate to 5.

**Rationale**:
- Prevents duplicates (FR-005).
- Puts the most recently viewed item first, which is the natural reading order for a recency strip.
- Truncation at 5 enforces FR-006 and keeps the cookie small.

**Algorithm** (pseudo-code):
```
ids = split(cookie, "|")
ids = filter(ids, id != currentProductID)   // deduplicate
ids = prepend(ids, currentProductID)         // most recent first
ids = ids[:min(5, len(ids))]                 // cap at 5
setCookie("shop_recently_viewed", join(ids, "|"))
```

---

## 7. Error Handling

**Decision**: If a product ID in the cookie no longer exists in the catalogue, silently skip it; if the cookie is absent or malformed, treat as an empty list.

**Rationale**:
- Graceful degradation: a stale or invalid cookie should not break page rendering.
- `getProduct` returns `codes.NotFound` for missing IDs — the helper should catch this error and continue rather than surfacing a 500.
- Consistent with the existing pattern where missing recommendations are logged and ignored.
