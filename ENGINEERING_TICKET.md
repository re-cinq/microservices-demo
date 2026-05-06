# Duplicate items in recommendation widgets

## Summary

The recommendation widgets occasionally render the same product in more than
one slot of a single response. Reported by a customer on the product page
("You May Also Like") and confirmed by marketing on the home page
("Hot Products"). The cart page widget is also affected — all three surfaces
share one backend RPC.

## Steps to reproduce

### Path A — via the UI (closest to customer report)

1. From the repo root, bring up the stack: `skaffold run` (or follow
   [README.md](README.md) for `docker-compose` or `kubectl` alternatives).
2. Wait until all pods/containers are `Ready`.
3. Open `http://localhost:8080` in a browser.
4. Hard-refresh with `Cmd+Shift+R` (macOS) or `Ctrl+Shift+R` (Linux/Windows)
   to bypass any browser cache.
5. Scroll to the "Hot Products" section (4-tile row near the bottom of the
   home page — the heading is rendered at
   [src/frontend/templates/home.html](src/frontend/templates/home.html), look for `<h3>Hot Products</h3>`).
6. Read the 4 product names. If all 4 are distinct, hard-refresh again.
7. Within ~5 hard-refreshes you will see at least one duplicated product in
   the row.

To reproduce on the product page widget instead:

1. From the home page, click any product (e.g. "Vintage Typewriter").
2. Scroll to "You May Also Like" (4 tiles, rendered by
   [src/frontend/templates/recommendations.html](src/frontend/templates/recommendations.html)).
3. Hard-refresh. A duplicate appears in roughly 4 of every 5 loads.

### Path B — deterministic, no browser

1. Bring up the stack as above.
2. Port-forward the recommendation service to localhost:8080
   (or run it directly with `python recommendation_server.py` after
   `pip install -r requirements.txt`).
3. From `src/recommendationservice/`, run the included client in a loop:
   ```bash
   for i in $(seq 1 10); do python client.py; done
   ```
4. Inspect the logged `product_ids` lists. Multiple runs will contain
   repeated IDs within a single response.

## Expected vs. actual

- **Expected:** each `ListRecommendations` response contains up to 5
  *distinct* product IDs.
- **Actual:** each response contains up to 5 product IDs, with repeats
  possible. The frontend renders the list as-is, so duplicates appear in
  adjacent or non-adjacent slots (the customer screenshot shows the Mug in
  slot 2 and slot 4).

## Suspected area

**Service:** `recommendationservice` (Python, gRPC).
**File / line:** [src/recommendationservice/recommendation_server.py:79](src/recommendationservice/recommendation_server.py).

```python
indices = random.choices(range(num_products), k=num_return)
```

Why this is the suspect:

- `random.choices` samples **with replacement** in Python's stdlib — by
  design, the same index can be drawn more than once in a single call. The
  without-replacement counterpart is `random.sample`. The observed symptom
  (occasional duplicate inside a single response) is exactly what
  `random.choices` produces.
- The frontend funnels all three reported surfaces through this one RPC:
  - product page → [src/frontend/handlers.go:178](src/frontend/handlers.go)
  - cart page → [src/frontend/handlers.go:266](src/frontend/handlers.go)
  - home page → [src/frontend/handlers.go:378](src/frontend/handlers.go)

  All three call `getRecommendations` in
  [src/frontend/rpc.go:99](src/frontend/rpc.go), which invokes
  `ListRecommendations`. A bug in this single line parsimoniously explains
  every reported surface.
- No de-duplication happens downstream — the template at
  [src/frontend/templates/recommendations.html:24](src/frontend/templates/recommendations.html)
  iterates the list as-is.
- Frequency math matches the customer's "if I refresh a few times I can
  usually reproduce it":
  - Catalog has 9 products
    ([src/productcatalogservice/products.json](src/productcatalogservice/products.json)).
  - Home page (no filter): pool 9, picks 5 → P(≥1 duplicate) ≈ 74%.
  - Product page (current product filtered out): pool 8, picks 5 → ≈ 80%.

This is a suspicion, not a verified fix — the engineer picking this up should
confirm by running Path B above and asserting non-uniqueness directly.

## Severity

**Proposed: P2 (high-visibility quality bug, no functional or revenue
impact).**

Rationale:

- **User cost:** appears on the home page and every product page — the two
  highest-traffic surfaces of the storefront. Roughly 3 of every 4 page
  loads will show a duplicate. Users notice; the customer ticket and the
  marketing report are independent confirmations.
- **Trust cost:** the customer used the word "unprofessional." A
  recommendation widget that repeats itself signals low quality and erodes
  confidence in the storefront, even though no transaction is broken.
- **Business cost:** the recommendation widget exists to drive
  cross-sell / additional pageviews. A duplicate slot wastes inventory in
  that widget — the user is shown 4 tiles but sees 3 distinct products.
  Effective surface area of the widget is reduced by ~20% on affected
  loads. We do not have CTR data attached, but the loss is non-zero.
- **Not P0/P1:** no checkout, payment, cart, or auth path is affected. No
  data is lost or corrupted. No security implication.
- **Not P3:** the bug is on the most-trafficked pages, reproducible on most
  loads, and externally reported. It is louder than a typical cosmetic
  ticket.

Fix scope is expected to be a single-line change in one Python file plus a
regression test, so cost-of-fix is very low relative to visibility — argues
for prompt prioritization within P2.

## Notes for the engineer picking this up

- The recommendation service logs the picked IDs at INFO
  ([recommendation_server.py:82](src/recommendationservice/recommendation_server.py))
  but does not flag duplicates. That is why no error logs exist for this
  bug — it has been silently producing duplicates the whole time.
- No frontend caching is in play; every page load triggers a fresh RPC, so
  changes to the recommendation service take effect immediately on the next
  request.
- A regression test should call `ListRecommendations` repeatedly and
  assert `len(set(response.product_ids)) == len(response.product_ids)` for
  each response. Without a loop the test will pass ~25% of the time by
  chance.
- Existing tests live alongside services
  (e.g. [src/productcatalogservice/product_catalog_test.go](src/productcatalogservice/product_catalog_test.go)).
  The recommendation service does not currently have a test file —
  adding one is part of the work.
