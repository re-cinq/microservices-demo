# BUG: Currency conversion produces incorrect values for high-rate currencies (JPY, TRY, KRW)

**Reported by:** Customer (Sara) via support ticket  
**Service:** `currencyservice`  
**Priority:** P1 — Revenue-impacting, actively exploitable

---

## Steps to Reproduce

1. Open the Online Boutique storefront with the default currency set to USD.
2. Add any product to your view — e.g. the Watch, listed at $109.99 USD.
3. Open the currency selector dropdown and switch to **JPY**.
4. Observe the displayed price for the Watch.
5. Repeat step 3–4 for **TRY** and **KRW** to confirm the pattern.
6. Proceed to checkout with a JPY-priced item and observe the checkout total.

Reproducible on: product listing page, product detail page, and checkout total. Refreshing the page does not change the result — the value is wrong on every load when a high-rate currency is selected.

---

## Expected vs Actual

| Currency | Expected (approx.) | Actual (observed) |
|----------|--------------------|-------------------|
| JPY      | ¥16,000+           | ¥0.92             |
| TRY      | Proportionally high| Near-zero         |
| CAD      | ~$150 CAD          | Approximately correct |
| EUR      | ~€96               | Approximately correct |

**Expected:** When a user selects a non-USD currency, all displayed prices — including on the product page and at checkout — should reflect an accurate conversion of the USD base price into the selected currency, consistent with the exchange rates stored in `currency_conversion.json`.

**Actual:** For currencies whose exchange rate vs EUR is significantly greater than 1 (i.e. currencies "weaker" than the dollar, such as JPY at ~126, TRY, KRW), the converted price is drastically understated — typically less than ¥1 for items worth tens of dollars. Currencies close to EUR parity (CAD, AUD) appear roughly correct because the arithmetic error is proportionally small.

The incorrect value is also used at checkout, meaning a customer can complete a purchase and be charged the wrong (near-zero) amount.

---

## Suspected Area

**File:** [`src/currencyservice/server.js`](src/currencyservice/server.js), lines 153–159  
**Function:** `convert()`

The service converts prices using a two-step EUR pivot: `from_currency → EUR → to_currency`. The `Money` type splits a value into `units` (whole number) and `nanos` (fractional part, stored in billionths). The conversion multiplies both fields independently by the target exchange rate:

```js
// Convert: EUR --> to_currency
const result = _carry({
  units: euros.units * data[request.to_code],
  nanos: euros.nanos * data[request.to_code]   // nanos can be up to ~1,000,000,000
});

result.units = Math.floor(result.units);  // suspected truncation point
result.nanos = Math.floor(result.nanos);
```

When `euros.nanos` (up to ~1,000,000,000) is multiplied by a large rate like JPY (~126), the result is ~126,000,000,000 — which `_carry()` should normalise back into `units`. However, `Math.floor(result.units)` on line 158 truncates fractional units produced by that carry, discarding the bulk of the converted value. For low-rate currencies (CAD ~1.5), the discarded fraction is negligible. For high-rate currencies (JPY ~126), it represents the majority of the correct price.

This is a suspicion based on static code reading. Confirm by logging `euros` and `result` before and after `_carry()` for a JPY conversion request.

---

## Severity: P1 — Revenue-blocking

**Classification:** P1 (not P0 only because the storefront remains operational and the failure is currency-scoped, not a full outage)

**Revenue impact (direct):** The checkout service uses the same converted value for billing. A customer selecting JPY can complete a purchase and be charged ¥0.92 for a $109.99 item. This is an active revenue loss on every non-USD transaction in affected currencies — not a hypothetical.

**Revenue impact (indirect):** Customers who notice the discrepancy before checkout are likely to abandon. Customers who notice after checkout may dispute the charge or request a refund, adding operational cost.

**Trust impact:** The error message makes the store look obviously broken to any user familiar with JPY or TRY denominations. A screenshot of "¥0.92 Watch" is the kind of thing that circulates on social media. Even one viral post creates brand damage disproportionate to the number of affected transactions.

**Scope:** Affects all currencies with an EUR rate significantly above 1. Based on `currency_conversion.json`, this includes at minimum: JPY, KRW, IDR, HUF, CZK, ISK, and TRY — a non-trivial share of potential international users.

**User impact:** Loyalty members and repeat international customers are most exposed; they are also most likely to notice and least likely to quietly accept it.
