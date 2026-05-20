# Research: Promo Code Invalid-Code Error Feedback (AIP-151)

## Prerequisite gap: AIP-150 not yet on attendee/nejal-patel

**Decision**: Cherry-pick the source-code portion of AIP-150 (`f14efd7a` from `origin/attendee/kate-payne`) before implementing AIP-151.

**Rationale**: `attendee/nejal-patel` does not contain the promo-code infrastructure (route, handler, cookie constant, validation map, cart template) that AIP-151 builds on. The most minimal path is a targeted cherry-pick of just the `src/` files; we skip kate-payne's spec/config artefacts.

**Alternatives considered**: Implement the full AIP-150 + AIP-151 scope from scratch. Rejected because it duplicates verified work and diverges the diff from the story boundary.

---

## Flash-message mechanism for error state

**Decision**: Use a short-lived cookie (`shop_promo-error`) as a one-shot flash signal from `applyPromoHandler` (POST) to `viewCartHandler` (GET after redirect).

**Rationale**: The app is server-side-rendered Go with a POST → redirect → GET checkout flow (mux.Router, `net/http`). There is no session store, no in-process state between requests. The existing `shop_promo-code` cookie establishes the pattern. A flash cookie (MaxAge=60s) carries the error flag across the redirect; `viewCartHandler` reads it, injects `promo_error=true` into template data, then immediately expires the cookie in the same response so it behaves as one-shot.

**Alternatives considered**:
- Query parameter (`/cart?promo_error=1`): avoids a cookie but leaks intent into bookmarkable URLs; inconsistent with the existing cookie-only pattern.
- In-memory state (sync.Map keyed on session): adds concurrency complexity and cross-pod inconsistency in a Kubernetes environment.

---

## Case-insensitive matching

**Decision**: Normalise the submitted code to `strings.ToUpper()` before map lookup.

**Rationale**: FR-009 requires case-insensitive matching. `strings.ToUpper` is the simplest single-line approach and is already used elsewhere in the Go codebase. All entries in `promoCodes` are uppercase keys, so `ToUpper` is sufficient without a full case-folding comparison.

**Alternatives considered**: `strings.EqualFold` per key in a loop — more explicit but O(n) and less idiomatic for map lookups.

---

## Error message copy

**Decision**: Inline message: `"Promo code not recognised — please check the code and try again."`

**Rationale**: Matches the acceptance criterion wording ("clear error message"), is self-describing, and avoids jargon. Placed immediately below the promo input form using the existing `text-danger` Bootstrap utility class already present in the cart template.

---

## Files to change (net of cherry-pick)

From cherry-pick (AIP-150 foundation):
- `src/frontend/main.go` — `cookiePromoCode` constant + route registration
- `src/frontend/handlers.go` — `promoCodes` map, `applyPromoHandler`, `viewCartHandler` discount logic, `placeOrderHandler` gRPC metadata
- `src/frontend/templates/cart.html` — promo input form, success confirmation, discount row
- `src/checkoutservice/main.go` — `promoCodes` map, `applyPromoDiscount`, gRPC metadata read

New for AIP-151:
- `src/frontend/main.go` — `cookiePromoError` constant
- `src/frontend/handlers.go` — set flash cookie on invalid code; read and expire it in `viewCartHandler`
- `src/frontend/templates/cart.html` — inline error display block
