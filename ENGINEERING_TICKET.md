# Engineering Ticket — Currency Conversion Returns Incorrect Values

**Reported by:** Customer (via support ticket)
**Severity:** High — financial impact on every non-USD checkout
**Service:** `currencyservice`
**File:** `src/currencyservice/server.js`

## Summary

Currency conversion produces incorrect (far too small) amounts for currencies
with exchange rates greater than 1 USD (e.g. JPY, TRY). The checkout total
reflects the same wrong values, meaning customers can complete purchases at a
fraction of the correct price.

## Actual Behaviour

When converting a price from EUR to a target currency, the service **divides**
by the exchange rate instead of multiplying:

```js
units: euros.units / data[request.to_code],
nanos: euros.nanos / data[request.to_code]
```

Observed on the live site:
- **JPY** prices are approximately **30% lower** than the correct value
- **TRY** prices are approximately **10x lower** than the correct value
- Currencies close to EUR parity (e.g. CAD ≈ 1.5) show only a small discrepancy

The magnitude of the error scales with the exchange rate: the higher the rate,
the further the result diverges from the correct value.

## Expected Behaviour

The conversion should **multiply** by the exchange rate. `data[request.to_code]`
holds the number of target-currency units per 1 EUR, so multiplying gives the
correct result:

```js
units: euros.units * data[request.to_code],
nanos: euros.nanos * data[request.to_code]
```

## Root Cause

Line 154–155 of `src/currencyservice/server.js` uses `/` instead of `*` in
the EUR → target-currency conversion step.

## Fix

```diff
- units: euros.units / data[request.to_code],
- nanos: euros.nanos / data[request.to_code]
+ units: euros.units * data[request.to_code],
+ nanos: euros.nanos * data[request.to_code]
```

## Impact

- All non-USD currency conversions return values that are too small by a factor
  of `rate²` (divide instead of multiply introduces an inverse error).
- Checkout totals are affected — customers are being charged the wrong amount.
- No data loss or security impact; purely a calculation error.
