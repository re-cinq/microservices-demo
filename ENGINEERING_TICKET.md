# Engineering Ticket — Visa Card Payments Rejected at Checkout

**Severity:** P0

## Steps to Reproduce

1. Add item to cart
2. Add delivery address
3. Attempt to pay with a Visa card

## Expected vs Actual

- **Expected:** Payment processes successfully
- **Actual:** Payment fails with incorrect card type error

## Suspected Area

`simple-card-validator` returns the card type as `'visa'` (lowercase), but the check in `src/paymentservice/charge.js` compares against `'Visa'` (capitalised). This case mismatch causes the condition to never match for Visa cards, triggering the `UnacceptedCreditCard` error. The string should be lowercase `'visa'` to match the library's output.
