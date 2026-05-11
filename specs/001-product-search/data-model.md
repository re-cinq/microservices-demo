# Data Model: Product Search

The feature does not introduce any persisted entities. The catalog source-of-truth remains `productcatalogservice/products.json`, loaded into memory by the existing `loadCatalog()` / `parseCatalog()` path. The feature introduces two transient, request-scoped entities (`SearchQuery`, `SearchResult`) which exist only for the duration of a single HTTP request and are never stored.

## Entities

### Product *(existing — no schema change)*

The catalog product, as already defined in `protos/demo.proto:77-87`.

| Field | Type | Source | Used by this feature for |
|---|---|---|---|
| `id` | `string` | products.json | Link href on each result card (`/product/{id}`). |
| `name` | `string` | products.json | **The only field matched against the query.** Also displayed on the result card. |
| `description` | `string` | products.json | Not searched in v1 (see plan Decision 2). Not displayed on the result card (matches home-page card layout). |
| `picture` | `string` | products.json | Displayed on the result card. |
| `price_usd` | `Money` | products.json | Converted to session currency via existing `convertCurrency` path; displayed on the result card. |
| `categories` | `[]string` | products.json | Not used by this feature. |

**Validation rules:** none added. The product schema is owned by `productcatalogservice` and is unchanged.

### SearchQuery *(new, transient)*

A single user-supplied search string, present only for the duration of one HTTP request to `/search`.

| Property | Type | Origin | Notes |
|---|---|---|---|
| raw query string | `string` | HTTP query parameter `q` | URL-decoded by `net/http`. |
| normalised query | `string` | derived | `strings.TrimSpace(raw)` — trimming is the only normalisation applied before reaching `SearchProducts`. |

**Validation rules (enforced in the frontend handler, see plan Decision 5):**

1. After trimming, the normalised query MUST NOT be empty. Empty → render the "please enter a search term" branch.
2. The normalised query MUST be ≤ 200 characters. Longer → render an error page; do not call `SearchProducts`.
3. Any other content (Unicode, punctuation, ASCII) is allowed and passed through to `SearchProducts` unchanged.

**Lifetime:** request-scoped. Never persisted. Logged once (sanitised) at request entry.

### SearchResult *(new, transient)*

A render-time view-model assembled by the search handler before template execution. Conceptually:

```text
SearchResult {
    query        : string         // the normalised SearchQuery, for echoing into the template
    items        : [ResultItem]   // ordered list of matched products, in catalog order
    item_count   : int            // len(items)
    no_matches   : bool           // item_count == 0
}

ResultItem {
    product      : Product        // pointer back to the catalog product (id, name, picture)
    price        : Money          // priceUsd converted to the session currency
}
```

| Property | Type | Origin |
|---|---|---|
| `query` | `string` | The normalised `SearchQuery`. Re-displayed in the page heading and the no-results message. |
| `items[*].product` | `*pb.Product` | A direct reference to the product returned by `SearchProductsResponse.results`, untransformed. |
| `items[*].price` | `*pb.Money` | The product's `price_usd` converted to `currentCurrency(r)` via `fe.convertCurrency`. |

**Ordering:** `items` is in the order returned by `SearchProducts`, which is the catalog's natural order (the order in which `products.json` lists them, filtered to matches). No re-ranking is applied in the frontend. This satisfies spec FR-010.

**Lifetime:** stack-scoped within the handler. The struct is built, passed to `templates.ExecuteTemplate`, and goes out of scope when the handler returns. Never persisted.

## Relationships

```text
SearchQuery (1) ─── normalises → request to SearchProducts(query)
                                       │
                                       ▼
                              SearchProductsResponse
                                  results: []Product
                                       │
                                       ▼
                                 [for each match]
                                  fe.convertCurrency(...)
                                       │
                                       ▼
SearchResult (1) ─── contains → ResultItem (0..N) ─── references → Product (existing)
```

- A `SearchQuery` produces exactly one `SearchResult` per request.
- A `SearchResult` references zero or more `Product` entities, by pointer; the catalog data is not copied or duplicated.
- The 1:N from `SearchResult` to `ResultItem` is bounded by the catalog size (currently ~10).

## State transitions

None — both new entities are immutable and request-scoped. The catalog `Product` has the same state model it has today (no transitions; static JSON).
