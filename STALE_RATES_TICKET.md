# Engineering Ticket — Exchange Rates Stale Since ~2019

**Reported by:** Internal QA (pricing discrepancy investigation)
**Severity:** Medium — incorrect prices shown, financial impact on non-USD checkouts
**Service:** `currencyservice`
**File:** `src/currencyservice/data/currency_conversion.json`

## Summary

The static exchange rate data bundled with `currencyservice` has not been
updated since approximately 2019. Several currencies have moved significantly
against the Euro since then, causing the shop to display prices that are
materially wrong compared to real-world values.

## Actual Behaviour

Prices in high-drift currencies are noticeably off against current rates:

| Currency | Rate in JSON | Approx. 2026 rate | Impact on displayed price |
|----------|-------------|-------------------|--------------------------|
| JPY | 126.40 / EUR | ~162 / EUR | ~22% too cheap |
| TRY | 6.12 / EUR | ~37.50 / EUR | ~6x too cheap |
| RUB | 74.42 / EUR | ~105 / EUR | ~30% too cheap |
| KRW | 1275 / EUR | ~1485 / EUR | ~14% too cheap |
| HRK | 7.42 / EUR | N/A — currency no longer exists | Croatia joined the Eurozone Jan 2023 |

## Expected Behaviour

Displayed prices should reflect exchange rates that are reasonably close to
current market values. The service uses a static file, so rates will always
lag real time, but the file should be refreshed periodically.

## Root Cause

`currency_conversion.json` is a static file bundled at build time and has
never been updated. There is no automated refresh mechanism.

## Fix

Update all rates to approximate mid-2025 values and remove HRK (Croatian
Kuna), which ceased to exist when Croatia adopted the Euro on 1 January 2023.

Notable changes:
- JPY: 126.40 → 162.00
- TRY: 6.1247 → 37.50
- CHF: 1.1360 → 0.9380 (CHF has strengthened against EUR)
- KRW: 1275.05 → 1485.00
- RUB: 74.42 → 105.00
- HRK: removed

## Follow-up Recommendation

Consider fetching live rates from a public ECB or Open Exchange Rates API
endpoint at service startup, with the static file as a fallback. This would
prevent the same drift recurring.
