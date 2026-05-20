# Quickstart: AIP-151 Promo Code Error Feedback

## Test the feature locally

**Prerequisites**: Go 1.21+, Docker (for full stack), or use the training URL.

**Training URL**: https://nejal-patel.training.gcp.re-cinq.com

## Manual test walkthrough

### Happy path (error → correction → success)

1. Navigate to the store, add any item to cart.
2. On the cart page, enter `BADCODE` in the promo field and click **Apply**.
3. **Expected**: inline error "Promo code not recognised — please check the code and try again." Order total unchanged.
4. Enter `SAVE10` and click **Apply**.
5. **Expected**: error clears; success message "Code SAVE10 applied — [amount] off"; total reduced by 10%.

### Error is non-blocking

1. Trigger the error with `BADCODE`.
2. **Without** clearing the field, click **Place Order** / proceed through checkout.
3. **Expected**: checkout completes at the full original price with no error banner visible on the confirmation page.

### No promo code (regression)

1. Add an item to cart and proceed without entering any code.
2. **Expected**: checkout and confirmation pages look and behave identically to before this feature.

## Build

```bash
cd src/frontend && go build ./...
cd src/checkoutservice && go build ./...
```
