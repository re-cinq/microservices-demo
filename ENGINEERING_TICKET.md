# Search ignores product descriptions

## Reproduce
1. Open the boutique. Focus the header search.
2. Search `modern` → 0 results.
3. Search `summer` → 0 results.
4. Search `watch` → returns Watch. Search `tank` → returns Tank Top.

## Expected vs actual
- **Expected:** `modern` returns Sunglasses (description: *"Add a modern touch…"*). `summer` returns Loafers (*"summer wardrobe"*). Match should cover **name and description**.
- **Actual:** Only name tokens match. Description tokens always return zero.

## Suspected area
`src/productcatalogservice/product_catalog.go:65`. `SearchProducts` does a single `strings.Contains(product.Name, query)` — `product.Description` is never read. Evidence: name tokens (`watch`, `tank`) hit; description tokens (`modern`, `summer`) miss. Frontend `searchHandler` and `rpc.go` just forward the query, no filtering.

The unit test `TestSearchProducts` (`src/productcatalogservice/product_catalog_test.go:91`) uses fixtures with empty `Description`, so it can't catch this regression.

## Severity — P2
Not P0/P1: checkout, browse, cart all work; users can still find items via the home page (9 products).
Not P3: silently wrong, not cosmetic. Search is the natural path for intent-driven queries ("modern", "summer", "gift") — exactly the descriptors that fail. A 200-OK zero-result is indistinguishable from "we don't sell this," so shoppers leave and we have no telemetry to detect the drop-off. Three reports in one morning suggests it's noticed quickly — trust cost compounds.
