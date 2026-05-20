# Data Model: Promo Code Invalid-Code Error Feedback (AIP-151)

## Entities

### PromoCode (value object, in-memory)

| Field | Type | Notes |
|-------|------|-------|
| Code  | string | Normalised to uppercase on entry. Valid values: `SAVE10`, `PROMO20`. |
| DiscountPct | int32 | Percentage off total (10 or 20). |

Source: `promoCodes` map in `src/frontend/handlers.go` and `src/checkoutservice/main.go`. No datastore.

### PromoErrorState (ephemeral, cookie-carried)

| Cookie name | Value | MaxAge | Purpose |
|-------------|-------|--------|---------|
| `shop_promo-error` | `"1"` | 60 s | Flash flag: set by `applyPromoHandler` on invalid submission; consumed and expired by `viewCartHandler` in the same redirect cycle. |

### PromoCodeState (cookie-carried, from AIP-150)

| Cookie name | Value | MaxAge |
|-------------|-------|--------|
| `shop_promo-code` | valid code string or empty | 48 h (or -1 to expire) |

## State Transitions

```
User submits promo input
        │
        ▼
  applyPromoHandler (POST /cart/promo)
        │
  ┌─────┴──────────────────────────┐
  │ code in promoCodes (valid)?    │
  └─────┬───────────────────┬──────┘
       YES                  NO
        │                   │
set shop_promo-code     clear shop_promo-code
(valid code)            set shop_promo-error=1
        │                   │
        └──────┬────────────┘
               │
        redirect GET /cart
               │
        viewCartHandler
               │
  ┌────────────┴──────────────────┐
  │ shop_promo-code set & valid?  │
  └────────────┬──────────────────┘
            YES │     NO
               │       │
    inject:    │   ┌───┴──────────────────┐
  promo_applied│   │shop_promo-error set? │
  promo_code   │   └───┬──────────────────┘
  promo_savings│      YES       NO
               │       │         │
               │  inject:     no promo state
               │  promo_error  in template
               │  expire cookie
               │
        render cart.html
```
