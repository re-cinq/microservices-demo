# Data Model: Browser-Side Product Search

**Feature**: `specs/001-product-search`
**Date**: 2026-06-02

> This feature introduces no new persistent data. All entities below describe data as it exists in the browser DOM at page load time — derived from the Go template rendering of the gRPC `ListProductsResponse`.

---

## Entity: ProductCard (DOM)

Represents one rendered product on the home page, as seen by the filter script.

| Attribute | Source | Description |
|---|---|---|
| `data-name` | `{{ .Item.Name }}` (Go template) | Product name; the sole attribute used for filtering |
| `display` (CSS) | Controlled by `product-search.js` | Visible (`''`) or hidden (`none`) based on current query |

**CSS selector**: `.hot-product-card[data-name]`

**Validation rules**:
- `data-name` is always present (set by the template for every product rendered)
- `data-name` is never empty (products without a name are not renderable)

---

## Entity: SearchQuery

Represents the current value of the search input field.

| Attribute | Type | Description |
|---|---|---|
| `raw` | string | The string as typed by the shopper |
| `normalised` | string | `raw.trim().toLowerCase()` — used for comparison |

**Validation rules**:
- Empty or whitespace-only `normalised` value → show all products, hide no-results message
- Any non-empty `normalised` value → apply substring filter

---

## Entity: FilterResult

Represents the outcome of applying a SearchQuery to all ProductCards.

| Attribute | Type | Description |
|---|---|---|
| `visibleCount` | integer | Number of ProductCards where `data-name` contains the query |
| `showNoResults` | boolean | `true` when `visibleCount === 0` and query is non-empty |

**State transitions**:

```
Query empty / whitespace
  → all ProductCards visible
  → no-results message hidden

Query non-empty, visibleCount > 0
  → matching ProductCards visible, non-matching hidden
  → no-results message hidden

Query non-empty, visibleCount === 0
  → all ProductCards hidden
  → no-results message visible
```

---

## Entities NOT in scope

The following `pb.Product` fields are rendered in the page but are not used by the search feature:

- `Id` — used in product detail link href
- `Picture` — product image
- `Description` — not searched (name-only per spec)
- `PriceUsd` / converted `Price` — not searched
- `Categories` — not searched
