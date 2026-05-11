# UI Contract: Product Search

**Feature**: 003-product-search
**Date**: 2026-05-11

This feature exposes no new public API, no new gRPC method, and no new HTTP endpoint. The only external interface it adds is a **UI/DOM contract** between the server-rendered `home.html` and the client-side `product-search.js`. This document is the canonical statement of that contract — what the template must render so that the script can do its job, and vice versa.

## DOM contract

### 1. Search input

```html
<div class="product-search">
  <label for="product-search-input" class="visually-hidden">Search products</label>
  <input
    id="product-search-input"
    class="product-search-input"
    type="search"
    autocomplete="off"
    placeholder="Search products"
    aria-controls="hot-products-row"
  />
  <div
    id="product-search-status"
    class="visually-hidden"
    role="status"
    aria-live="polite"
    aria-atomic="true"
  ></div>
</div>
```

**Required** (must be present and unchanged for the JS to work):

| Selector / attribute | Purpose |
|---|---|
| `#product-search-input` | The element the JS listens to. `input` events drive the filter. |
| `type="search"` | Yields the browser-native clear (×) button; emits `input` events when cleared. |
| `aria-controls="hot-products-row"` | Associates the input with the grid for assistive tech. |
| `#product-search-status` with `role="status" aria-live="polite"` | The live region the JS writes the count into. Must be visually hidden but present in the accessibility tree. |
| `<label for="product-search-input">` | Required for screen-reader labelling. The visible placeholder is **not** sufficient. |

**Forbidden:**

- No `name` attribute on the input (we don't want it submitted in any form action).
- No `autocomplete` value other than `"off"`.

### 2. Product grid

```html
<div class="row hot-products-row px-xl-6" id="hot-products-row">
  …
  {{ range $.products }}
  <div class="col-md-4 hot-product-card" data-name="{{ lower .Item.Name }}">
    …
    <div class="hot-product-card-name">{{ .Item.Name }}</div>
    …
  </div>
  {{ end }}
</div>
```

**Required:**

| Selector / attribute | Purpose |
|---|---|
| `#hot-products-row` (added) | Target of `aria-controls` on the input. |
| `.hot-product-card` (existing) | The JS query-selects all of these to filter. |
| `data-name` on each `.hot-product-card` | Lowercased product name, used by the JS for substring match. Must be set at render time. |

**Forbidden:**

- The JS must **not** read the visible `.hot-product-card-name` text content. It uses `data-name` only. (This keeps the contract decoupled from any future visible-name formatting changes.)

### 3. Empty-state message

```html
<div class="product-search-empty" hidden>No products match your search.</div>
```

- Lives **inside** `#hot-products-row` or immediately after it.
- The JS sets `hidden=false` (i.e. removes the attribute) only when the query is non-empty and zero cards match. It must be `hidden` in all other states (including initial page load).

## Behaviour contract (what the JS does, in order, on every `input` event)

1. Compute `q = input.value.trim().toLowerCase()`.
2. If `q === ""`:
   - For every `.hot-product-card`: remove the `hidden` attribute.
   - Hide `.product-search-empty`.
   - Set `#product-search-status` text to `""`.
   - Return.
3. Otherwise, for every `.hot-product-card`:
   - If `card.dataset.name.includes(q)`: remove `hidden`.
   - Else: set `hidden` to true.
4. Count `n` = number of cards without `hidden`.
5. If `n === 0`:
   - Set `.product-search-empty` to **not** hidden.
   - Set `#product-search-status` text to `"No products match your search."`.
6. Else:
   - Set `.product-search-empty` to hidden.
   - Set `#product-search-status` text to `"<n> products match"` (or `"1 product matches"` when `n === 1`).

## Non-contract (what this feature deliberately does NOT promise)

- **No URL/query-param round-trip.** Searching does not change the URL; reloading the page does not preserve the query (FR-010).
- **No history entry.** No `pushState`/`replaceState`.
- **No fetch.** No `XMLHttpRequest`, no `fetch()`, no WebSocket. (SC-003.)
- **No global state.** The JS attaches one event listener and operates on the current page only.
- **No matching on description, category, or ID.** Only `data-name`. (C-009.)

## Network contract (negative)

This feature MUST cause **zero** new network requests in response to any user input on the search box. Browser DevTools "Network" panel must show no rows when a shopper types or clears the search input. (Verifiable in `quickstart.md`.)
