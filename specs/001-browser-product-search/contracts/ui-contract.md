# UI Contract — Browser-Side Product Search

**Feature**: [spec.md](../spec.md) | **Plan**: [plan.md](../plan.md)

This feature has no API surface (no new HTTP route, no new gRPC method, no new database schema). Its "contract" is therefore a **UI contract** — the agreement between the server-rendered template and the client-side script about what elements exist in the DOM, what attributes they carry, and what events flow between them.

Anything outside this document is implementation detail and may be refactored freely; anything inside it is load-bearing and must be preserved by both the template and the script.

---

## 1. Required DOM elements

### 1.1 Search input

```html
<input
  type="search"
  id="product-search"
  class="product-search"
  aria-label="Search products"
  aria-controls="hot-products-grid"
  autocomplete="off"
/>
```

| Property | Required value | Why |
|---|---|---|
| `type` | `search` | Native clear button; Escape-to-clear in most browsers (User Story 2). |
| `id` | `product-search` | Stable JS selector. |
| `class` | includes `product-search` | CSS selector hook (styling). |
| `aria-label` | non-empty (e.g. `"Search products"`) | FR-011, SC-007 — accessible name. |
| `aria-controls` | the grid container's id | Associates input with the region it filters. |
| `autocomplete` | `off` | Browsers should not surface unrelated history items in a filter UI. |

The element MUST live inside the home template and MUST be rendered above (in DOM order, before) the product grid container.

### 1.2 Product grid container

```html
<div id="hot-products-grid" class="row hot-products-row px-xl-6">
  ...product cards...
</div>
```

| Property | Required value | Why |
|---|---|---|
| `id` | `hot-products-grid` | Target of `aria-controls` on the search input. |
| Existing classes (`row hot-products-row px-xl-6`) | preserved | The existing layout depends on them. |

### 1.3 Product card

Each card produced by `{{ range $.products }}` MUST carry:

| Attribute | Value | Why |
|---|---|---|
| `class` includes `hot-product-card` | (existing class hook) | JS selector. If the template doesn't already use exactly this class name, the JS selector in `product-search.js` MUST be updated to match the actual class — the contract is "one stable class shared by all cards", not specifically this string. |
| `data-name` | `{{ .Item.Name | lower }}` (or equivalent) — the lowercased product name | The only field the filter matches against. Pre-lowercased so the JS does not re-lowercase on every keystroke. |

The card's existing markup, image, link, price, etc. are unchanged.

### 1.4 No-results message

```html
<p id="product-search-no-results" class="no-results" role="status" aria-live="polite" hidden>
  No products match your search.
</p>
```

| Property | Required value | Why |
|---|---|---|
| `id` | `product-search-no-results` | Stable JS selector. |
| `role` | `status` | Empty-state announcement. |
| `aria-live` | `polite` | Screen readers hear the announcement when the JS removes `hidden`. |
| `hidden` (initial state) | present | The page starts unfiltered → no empty state shown. |

This element MUST be a sibling of the grid container (immediately after it in DOM order is conventional).

---

## 2. Client-side behaviour contract

The script at `src/frontend/static/js/product-search.js` MUST implement exactly the following behaviour. Implementations are free to choose internal helpers; the **observable** behaviour below is the contract.

### 2.1 Initialization

On `DOMContentLoaded` (or equivalent — `defer` script is sufficient):

1. Locate the search input by `#product-search`. If absent, do nothing and return. (This script must be safe to load on pages that don't have the box.)
2. Locate all product cards by `.hot-product-card` (or whatever class hook the template uses).
3. Locate the no-results element by `#product-search-no-results`.
4. Attach an `input` event listener to the search input. (No `keyup`, no debounce.)

### 2.2 On every `input` event

```text
query     = inputElement.value
normalized = query.trim().toLowerCase()

if normalized === "":
    for each card: remove `hidden` attribute (show all)
    set no-results element `hidden`
    return

visibleCount = 0
for each card:
    if card.dataset.name.includes(normalized):
        remove `hidden`
        visibleCount += 1
    else:
        set `hidden`

if visibleCount === 0:
    remove `hidden` on no-results element
else:
    set `hidden` on no-results element
```

This pseudocode is the single source of truth for the filter algorithm. Any deviation requires updating this contract first.

### 2.3 Pure-function export (testability)

The script MUST expose its match decision as a pure function, attached to `window.__productSearch.matches`:

```js
window.__productSearch.matches = function (dataName, query) {
  const normalized = query.trim().toLowerCase();
  if (normalized === "") return true;
  return dataName.includes(normalized);
};
```

This lets the quickstart verification step paste assertions into the browser console (e.g. `console.assert(window.__productSearch.matches("sunglasses", "SUN"))`).

### 2.4 Events that this feature does NOT listen for

- `keyup`, `keydown`, `change`, `submit`, `focus`, `blur`. The contract is `input`-only.
- No `form` element wraps the search input. There is no submit semantic.

### 2.5 Side effects this feature is NOT allowed to have

- MUST NOT make any `fetch`, `XMLHttpRequest`, or other network call.
- MUST NOT modify `localStorage`, `sessionStorage`, cookies, or the URL.
- MUST NOT change any DOM outside the home page's product grid and the no-results element.
- MUST NOT alter the order of product cards. Hidden cards remain in the DOM at their original position.

---

## 3. Backwards/forwards compatibility

- **JS disabled**: The page MUST still render the full product list. The search input either is not rendered at all, or is rendered as inert (it would simply do nothing on input). Either is acceptable for v1; the simplest implementation — render it always — is preferred.
- **Adding a new product** to `products.json` later: requires no JS change. As soon as the new product card is rendered with a `data-name`, the filter picks it up.
- **Removing the feature** later: deleting the input, the no-results element, the script tag, and the `data-name` attribute restores the original page. No data migration, no rollback complexity.

---

## 4. Non-contract (deliberately unspecified)

The following are implementation choices the implementer may make freely without violating this contract:

- Exact CSS styling (colour, spacing) of `.product-search` and `.no-results` — as long as they are visible and don't break the existing grid.
- Whether the `<label>` is visually shown or screen-reader-only.
- Whether the no-results message is wrapped in any additional container.
- Internal variable names, function names, and module shape inside `product-search.js`.
