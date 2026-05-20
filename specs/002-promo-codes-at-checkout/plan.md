# Plan: Apply a valid promo code at checkout

**Spec:** `specs/002-promo-codes-at-checkout/spec.md`
**Jira:** AIP-150 (Story) · AIP-96 (Epic)

---

## Tech stack

- **Frontend:** Go (`src/frontend/`), gorilla/mux, Go HTML templates, HTTP cookies
- **Checkout service:** Go (`src/checkoutservice/`), gRPC server
- **Transport:** gRPC with protobuf (`protos/demo.proto`, generated into both services' `genproto/` dirs)
- **Money:** custom `Money` type (units + nanos) with helpers in `src/*/money/money.go`
- **No new services, datastores, or infra changes** (hard constraint from AIP-96)

---

## Architecture

### The core problem

`placeOrderHandler` (frontend) does **not** pass an order total to checkoutservice — it passes user/address/card data and `checkoutservice.PlaceOrder()` independently fetches the cart and computes the total before charging. To apply a real discount to the payment, the promo code must travel through the gRPC call.

### Decision: proto change + cookie-based cart-page preview

| Concern | Approach |
|---|---|
| Actual discount on charge | Add `promo_code` to `PlaceOrderRequest` proto; checkoutservice applies discount before `chargeCard()` |
| Discounted total on cart page (before "Place Order") | Separate `POST /cart/promo` handler validates code and sets `shop_promo-code` cookie; `viewCartHandler` reads cookie and applies discount to displayed total |
| Cookie naming | Follow existing `cookiePrefix = "shop_"` pattern (same as `shop_currency`) |
| Code storage | Hardcoded `map[string]int32` (code → percentage off) in both frontend and checkoutservice — no new file, no new datastore |

### Request flow

```
GET /cart
  └─ viewCartHandler
       reads shop_promo-code cookie
       if valid: applies discount to total_cost, sets promo_applied=true in template data

POST /cart/promo  (new)
  └─ applyPromoHandler
       validates code against frontend promoCodes map
       sets shop_promo-code cookie (max-age = cookieMaxAge)
       redirects to GET /cart

POST /cart/checkout
  └─ placeOrderHandler
       reads shop_promo-code cookie
       passes code in PlaceOrderRequest.PromoCode
       clears promo cookie after successful order

gRPC: CheckoutService.PlaceOrder(PromoCode="SAVE10")
  └─ PlaceOrder in checkoutservice
       validates code against checkoutservice promoCodes map
       computes total (existing logic)
       if valid code: applies percentage discount using nanos arithmetic
       charges discounted total
```

---

## Design decisions

**1. Promo codes hardcoded in both services**
The frontend needs codes to validate for display before the user places the order. The checkoutservice needs them to apply the actual discount. Keeping identical maps in both services (3–5 codes) is the simplest approach with no new infrastructure. A mismatch between the two maps would be caught in testing.

Codes (illustrative — can be changed before ship):

| Code | Discount |
|---|---|
| `SAVE10` | 10% off |
| `PROMO20` | 20% off |

**2. Percentage discount applied using nanos arithmetic**
`pb.Money` stores `units` (int64) and `nanos` (int32). To apply a percentage without floating point:
```go
totalNanos := m.Units*1_000_000_000 + int64(m.Nanos)
discountedNanos := totalNanos * int64(100-pct) / 100
result.Units = discountedNanos / 1_000_000_000
result.Nanos = int32(discountedNanos % 1_000_000_000)
```
This avoids adding a new helper to the money package.

**3. Cookie cleared after order placement**
`placeOrderHandler` clears the `shop_promo-code` cookie on success so the discount does not silently carry over to the shopper's next visit.

**4. Proto field placed at tag 7**
`PlaceOrderRequest` currently uses tags 1–3 and 5–6 (tag 4 is absent). New field `promo_code` takes tag 7 to avoid any field-number collision.

---

## File change map

| File | Change |
|---|---|
| `protos/demo.proto` | Add `string promo_code = 7` to `PlaceOrderRequest` |
| `src/checkoutservice/genproto/demo.pb.go` | Regenerate via `src/checkoutservice/genproto.sh` |
| `src/checkoutservice/genproto/demo_grpc.pb.go` | Regenerate (same script) |
| `src/frontend/genproto/demo.pb.go` | Regenerate (frontend has its own genproto copy) |
| `src/frontend/genproto/demo_grpc.pb.go` | Regenerate |
| `src/checkoutservice/main.go` | Add `promoCodes` map; add `applyPromoDiscount()` helper; call it in `PlaceOrder()` after total is computed |
| `src/frontend/main.go` | Register `POST /cart/promo` → `applyPromoHandler`; add `cookiePromoCode` constant |
| `src/frontend/handlers.go` | Add `applyPromoHandler`; update `viewCartHandler` (read cookie, pass discount to template); update `placeOrderHandler` (read cookie, pass to RPC, clear cookie on success) |
| `src/frontend/templates/cart.html` | Add promo code input + Apply button form; conditional discount row in cart summary; conditional inline confirmation message |

---

## Out of scope for this story

- Error feedback for unrecognised codes (AIP-151)
- Discount on order confirmation page (AIP-152)
- Proto regeneration scripts beyond running the existing `genproto.sh`
- Any change to Kubernetes manifests, Dockerfiles, or CI
