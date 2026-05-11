# Data Model: Product Search

**Feature**: 003-product-search
**Date**: 2026-05-11

This feature does **not** add or modify any persisted data model. There is no schema change, no new protobuf message, and no new gRPC method. The only "data" introduced by this feature is transient browser state.

## Entities

### 1. Product (existing — unchanged)

The product entity is read-only for this feature and continues to come from `productcatalogservice/products.json` via the existing `productcatalogservice.ListProducts` gRPC call. The frontend then renders each product as one `.hot-product-card` element in `home.html`.

The fields the feature **reads** from each rendered card:

| Field | Source | Used for |
|---|---|---|
| `name` (string, e.g. `"Sunglasses"`) | The existing `{{ .Item.Name }}` Go template binding | Substring match against the user's query |

The feature does **not** read or rely on `id`, `description`, `categories`, `priceUsd`, or `picture`.

To make the name reliably available to the client-side filter, we add **one** new DOM attribute on each card:

| New DOM attribute | Value | Why |
|---|---|---|
| `data-name` on `.hot-product-card` | The product's name, lowercased once at render time (`{{ lower .Item.Name }}` or equivalent) | Allows the filter to do `card.dataset.name.includes(query)` with no per-card `toLowerCase()` cost per keystroke, and avoids relying on parsing nested HTML for the name. |

This is a presentational annotation only — it does not change the product domain model.

### 2. Query (new — client-side only)

`Query` is in-memory browser state held by the search input. It is never persisted, never sent over the network, and never observed by any backend service.

| Field | Type | Source | Constraints |
|---|---|---|---|
| `raw` | string | The current `.value` of the `<input type="search">` element | Any string the browser accepts as input |
| `normalized` | string | `raw.trim().toLowerCase()` (computed per keystroke) | Used for the substring match |

#### State transitions

| From | Event | To | Visible effect |
|---|---|---|---|
| `normalized == ""` | User types a non-whitespace character | `normalized == "<chars>"` | Cards whose `data-name` includes `normalized` remain visible; others get `hidden`. Live region announces "N products match" (or "No products match your search." if N == 0). |
| `normalized == "<chars>"` | User extends/shortens input | `normalized == "<new>"` | Re-filter; live region re-announces. |
| `normalized != ""` | User clears input (Backspace, native × button, or selects-all + Delete) | `normalized == ""` | All cards revert to visible (no `hidden`). Live region clears. |
| any | User navigates away and returns | `raw == ""` | Page re-renders fresh; query state does not persist (FR-010). |

### 3. Result count (derived, client-side only)

`resultCount` is a derived integer: the number of `.hot-product-card` elements not currently `hidden`. It is computed only to drive the polite live-region announcement and the visible empty-state message. It is never stored.

## Validation rules

- The query is treated as opaque text; there is no syntactic validation. (No "min length", no banned characters.)
- Whitespace-only queries are equivalent to the empty query (FR-005).
- Very long queries (≥ a few hundred chars) are handled identically; they simply produce an empty result set (R3 in research.md).

## Relationships

`Query` is a one-to-many filter over the rendered `Product` cards. There are no relationships introduced into the persisted domain.
