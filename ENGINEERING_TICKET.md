# Duplicate products in "You may also like" recommendations

**Source ticket:** `BUG_REPORT.md` (customer Priya, loyalty member)
**Severity:** P3 — low (cosmetic, no data loss, no checkout impact, but customer-visible and brand-eroding; visible on ~60% of product page loads)

## Summary

The recommendation service samples product indices **with replacement** (`random.choices`) instead of without replacement (`random.sample`). When the chosen sample has fewer than 5 unique items, the same product id appears in multiple slots of the "You may also like" section.

## Steps to reproduce

1. Open the storefront (e.g. `https://<your-name>.training.gcp.re-cinq.com`).
2. Click any product to open its product page.
3. Scroll to **"You may also like"** — the frontend caps the section at 4 slots ([`src/frontend/rpc.go:113-115`](src/frontend/rpc.go)).
4. Hard-refresh the page (Ctrl+F5) repeatedly.
5. Within ~5 refreshes you will see the same product image/name appear in two or more slots.

The customer's screenshot shows the **Mug** appearing in slot 2 and slot 4 simultaneously — a textbook instance of this.

## Expected vs actual

- **Expected:** Each of the 4 visible recommendation slots shows a distinct product, none of which is the product currently being viewed.
- **Actual:** Slots can repeat the same product. The current product is correctly excluded; only the de-duplication between slots is broken.

## Suspected service + file

- **Service:** `recommendationservice` (Python, gRPC).
- **File:** [`src/recommendationservice/recommendation_server.py`](src/recommendationservice/recommendation_server.py) — method `RecommendationService.ListRecommendations`, **line 79**:

  ```python
  indices = random.choices(range(num_products), k=num_return)
  ```

  `random.choices` samples *with* replacement, so the same index can be drawn multiple times. The catalog has 9 products; the current product is excluded, leaving 8 to draw 5 from (the frontend then truncates to the first 4 — `src/frontend/rpc.go:113-115`). Probability of at least one duplicate among the 4 visible slots is ≈ 1 − (8·7·6·5)/8⁴ ≈ **59 %** per page load, which matches the customer's "if I refresh a few times I can usually reproduce it" (≈ 93 % chance of seeing it within 3 refreshes).

## Suggested fix

Replace `random.choices` with `random.sample`, which draws *without* replacement:

```python
indices = random.sample(range(num_products), k=num_return)
```

This is safe because `num_return = min(max_responses, num_products)` is already clamped to the population size, so `random.sample` cannot raise `ValueError`.

## Severity rationale

P3 (low):

- **Customer impact:** cosmetic. No incorrect prices, no broken checkout, no data loss. Brand-quality issue only.
- **Reach:** every visitor to any product page (≈ all sessions), with a duplicate visible on roughly 3 out of 5 page loads.
- **Workaround:** none needed; the page is still functional.
- **Effort to fix:** one-character change, one service redeploy.

Not P2 because there is no functional or revenue impact. Not P4 because it is on the most-trafficked page type (product detail) and visible to every customer.

## Notes on the second symptom ("Hot Products" widget)

The customer mentions, second-hand from a colleague, that **"Hot Products"** on the home page also shows duplicates. Investigation:

- [`src/frontend/handlers.go`](src/frontend/handlers.go) — `homeHandler` passes the full product catalog to the template without sampling.
- [`src/frontend/templates/home.html`](src/frontend/templates/home.html) — iterates `$.products` once, no random selection.
- [`src/productcatalogservice/products.json`](src/productcatalogservice/products.json) — 9 unique product ids, no duplicate entries.

There is no code path on the home page that can produce duplicates. Most likely the colleague conflated the product-page recommendations with the home-page "Hot Products" grid (both show 3-across product cards in the same visual style). **No separate fix required**, but worth a follow-up confirmation with the customer once the recommendation fix is live.
