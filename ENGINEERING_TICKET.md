# Engineering ticket — Duplicate items in recommendations

**Service:** `recommendationservice` (Python)
**File:** `src/recommendationservice/recommendation_server.py`
**Surfaces affected:** product page "You may also like" widget, home page "Hot Products" widget (both call `RecommendationService.ListRecommendations`)
**Source:** customer ticket from Priya — "Same product showing twice in recommendations"

## Summary

`ListRecommendations` samples product indices **with replacement**, so the same product can be returned more than once in a single response. The customer-visible result is the same item appearing in two recommendation slots.

## Root cause

`recommendation_server.py:79`:

```python
indices = random.choices(range(num_products), k=num_return)
```

`random.choices` samples **with replacement** — repeated indices are expected behaviour. The deduplication on line 75 (`set(product_ids) - set(request.product_ids)`) only removes products the caller already knows about; it does not prevent the sampler from picking the same index twice.

## Reproduction

1. Open any product page on the site, e.g. `/product/<id>`.
2. Refresh several times.
3. Observe two identical tiles in the "You may also like" strip.

Probability of at least one duplicate when drawing `k=5` from a pool of `n` distinct products is `1 - n!/((n-k)! * n^k)`. With ~9 candidate products (catalogue size minus the current product), that's roughly **70% per request** — matches the customer's "if I refresh a few times I can usually reproduce it".

## Fix

Use sampling **without replacement**:

```python
indices = random.sample(range(num_products), num_return)
```

`random.sample` guarantees distinct indices. `num_return` is already clamped to `min(max_responses, num_products)` on line 77, so the call is safe when the candidate pool is smaller than `max_responses`.

## Test plan

- Unit: call `ListRecommendations` 1000× with a fixed seed, assert `len(set(resp.product_ids)) == len(resp.product_ids)` for every response.
- Manual: hit a product page 20× in a row, confirm no duplicate tiles in the "You may also like" strip.
- Manual: same check on the home page "Hot Products" widget.

## Risk / blast radius

- One-line change in a single Python file.
- Only the `recommendationservice` image needs to be rebuilt and rolled.
- No API/proto changes, no DB migration, no client changes.
- Behaviourally identical to the previous correct implementation (this is a regression — `git log` on `main` shows `random.sample` is the original code).
