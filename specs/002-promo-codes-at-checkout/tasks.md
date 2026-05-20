# Tasks: Apply a valid promo code at checkout

**Spec:** `specs/002-promo-codes-at-checkout/spec.md`
**Plan:** `specs/002-promo-codes-at-checkout/plan.md`

---

## Phase 1 — Proto (foundation for all other work)

- [x] [T1] ~~Add `string promo_code = 7;` to `PlaceOrderRequest` in `protos/demo.proto`~~ **Revised:** `protoc` not available in this environment; promo code passed via gRPC metadata instead (no proto change required — T2/T3 skipped)
- [x] [T2] ~~Regenerate checkout service proto~~ Skipped — no proto change
- [x] [T3] ~~Regenerate frontend proto~~ Skipped — no proto change

---

## Phase 2 — Checkout service (discount applied to actual charge)

- [x] [T4] Add package-level `promoCodes map[string]int32` hardcoded with `SAVE10→10` and `PROMO20→20` (percent off) in `src/checkoutservice/main.go`
- [x] [T5] Add `applyPromoDiscount(m pb.Money, pct int32) pb.Money` helper using nanos arithmetic in `src/checkoutservice/main.go`
- [x] [T6] In `checkoutService.PlaceOrder()`, extract promo code from gRPC incoming metadata and apply discount before `chargeCard()` *(revised from proto field to gRPC metadata)*

---

## Phase 3 — Frontend handler (cookie state + RPC wiring)

- [x] [T7] Add `cookiePromoCode = cookiePrefix + "promo-code"` constant in `src/frontend/main.go`; add matching `promoCodes map[string]int32` as package-level var in `src/frontend/handlers.go`
- [x] [T8] Add `applyPromoHandler` in `src/frontend/handlers.go`; register `POST /cart/promo` route in `src/frontend/main.go`
- [x] [T9] Update `viewCartHandler` in `src/frontend/handlers.go` to read promo cookie and pass `total_cost`, `promo_applied`, `promo_code`, `promo_savings` to template
- [x] [T10] Update `placeOrderHandler` to attach promo code as gRPC outgoing metadata; clear cookie on successful order

---

## Phase 4 — UI (cart template)

- [x] [T11] Add promo code form (input + Apply button) to `src/frontend/templates/cart.html` between shipping row and total row
- [x] [T12] Add conditional discount savings row, updated total label, and inline confirmation message to `src/frontend/templates/cart.html`

---

## Completion checklist

- [x] `go build ./...` passes in both `src/frontend/` and `src/checkoutservice/`
- [x] Entering a valid code on the cart page shows a discounted total and confirmation
- [x] Replacing one valid code with another updates the total correctly
- [x] Visiting the cart without a code looks identical to the current page
