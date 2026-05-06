# TICKET: Currency conversion returns wrong amounts for non-USD currencies

**Severity:** High — financial impact, affects all non-USD checkouts  
**Reported by:** Customer (Sara) via support  
**Service:** `currencyservice` (`src/currencyservice/server.js`)

---

## Summary

All prices displayed and used at checkout in non-USD currencies are incorrect.
Currencies with high exchange rates versus USD (e.g. JPY, TRY) show near-zero
values; customers can complete purchases at a tiny fraction of the correct price.

## Steps to Reproduce

1. Load Online Boutique and add any item to the cart.
2. Switch currency dropdown to **JPY** (or TRY, IDR, or any currency with rate >> 1).
3. Observe: a $109.99 Watch shows as **¥0.92** instead of ~¥16,000.

## Root Cause

`server.js:154-155` — the EUR→target-currency step uses **division** instead of
**multiplication** by the exchange rate:

```js
// BUG (introduced in commit f0682c52)
units: euros.units / data[request.to_code],
nanos: euros.nanos / data[request.to_code]
```

The `data` map stores rates as *units of currency per 1 EUR*. To convert from EUR
to a target currency you must **multiply** by the rate. Dividing inverts the
conversion, producing values that are correct only when the rate ≈ 1 (e.g. CHF).

The first half of the function (from_currency → EUR) correctly uses division
(`/ data[from.currency_code]`), so only the outbound leg is broken.

## Fix

Restore the two `/` operators on lines 154-155 to `*`:

```js
units: euros.units * data[request.to_code],
nanos: euros.nanos * data[request.to_code]
```

## Impact

- All non-USD currency display values are wrong.
- Checkout totals use the same broken numbers — financial loss on every
  non-USD transaction until fixed.
- USD-only purchases are unaffected (no conversion performed).
