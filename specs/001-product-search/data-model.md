# Phase 1 Data Model: Product Search

**Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md)

## Server-side entities

**None added.** The catalogue and its `Product` proto are unchanged (C-014). The feature reads only `Product.name` — a field that already exists on every product — and reads it through the existing `homeHandler` → `getProducts` → `productView{Item, Price}` flow (`src/frontend/handlers.go:79-91`). The template already iterates `$.products`; no new fields, no new types.

## The only new contract: DOM ↔ JS

The JavaScript module in `static/js/search.js` depends on three small contracts that the template (`templates/home.html`) is responsible for establishing on every render.

### Contract A — `data-product-name` on each card

Each rendered product card MUST carry a `data-product-name` attribute whose value is the product's `name`:

```html
<div class="col-md-4 hot-product-card" data-product-name="Watch">…</div>
<div class="col-md-4 hot-product-card" data-product-name="Salt &amp; Pepper Shakers">…</div>
```

- **Source**: Go template field `.Item.Name`, auto-escaped for HTML attribute context by `html/template`.
- **Consumer**: `search.js` reads `card.dataset.productName.toLowerCase()`.
- **Invariant**: every rendered product card has a non-empty `data-product-name`. Products with empty / missing `name` in `products.json` are matched against by the search (they'd never match a non-empty query, per Edge Cases — "products with missing names").
- **Selector used by JS**: `document.querySelectorAll('.hot-product-card')` and then per-card `card.dataset.productName`.

### Contract B — three named elements on the page

| Element id            | Element        | Initial state              | Purpose                                              |
|-----------------------|----------------|----------------------------|------------------------------------------------------|
| `search-input`        | `<input type="text">` | rendered, value empty, no focus | The search text field. Has visible label / aria-label. |
| `search-clear`        | `<button type="button">` | rendered, `hidden`            | Click to clear. JS removes `hidden` when input is non-empty; sets it again when empty. |
| `search-empty-state`  | `<div>`        | rendered with `hidden`, `role="status"`, `aria-live="polite"`, text `No products match your search.` | Shown when no cards match and query is non-empty. ARIA live region announces appearance to screen readers. |

- **Source**: hard-coded in `home.html`, inserted directly above the `<div class="row hot-products-row …">` so the existing grid is not displaced (C-010 / FR-001).
- **Consumer**: `search.js`.
- **Invariant**: all three elements exist on every render of the home page. The empty-state element's **text content is never modified by JS** (closes the XSS surface — see Decision 4 in `research.md`); only the `hidden` attribute is toggled.

### Contract C — page-scoped script tag

```html
<script src="{{ $.baseUrl }}/static/js/search.js" defer></script>
```

- Placed at the end of `home.html`'s body — page-scoped, not in `header.html` (Decision 2 in `research.md`).
- `defer` so it runs after DOM parsing.
- If JS fails to load, the search box still renders but doesn't filter — per the spec's "no JavaScript / client filtering disabled" edge case, the page is not broken.

## Client-side state

`search.js` holds **no persistent state** (C-015). Everything is derived from the DOM at the moment of each keystroke:

- **Query**: read fresh from `document.getElementById('search-input').value` each time the `input` event fires; trimmed; lowercased.
- **Match count**: counted by iterating cards and toggling `style.display` as we go.

There is no in-memory cache of products in JS, no debounce timer, no AbortController, no Promise queue. Each keystroke is a single synchronous pass: trim → lowercase → iterate cards → set `display` → if match count is 0 and query is non-empty, show empty-state, else hide it → if input has any text, show clear button, else hide it.

The `input` event from the search box always reflects the **current** value of the input; rapid typing therefore can't produce stale results (FR-002, US1 #7) — there is no race because there is no asynchrony.

## State machine (client-side, derived not stored)

```
        empty input ──── type 1+ chars ────► non-empty input
            ▲                                       │
            │                                       │
            │ clear button                          │ matches change as query changes
            │ Esc                                   │
            │ delete all chars                      ▼
            │                            ┌──────────────────────┐
            └────────────────────────────┤  matches > 0         │
                                         │  →  show subset      │
                                         │  →  hide empty-state │
                                         ├──────────────────────┤
                                         │  matches = 0         │
                                         │  →  hide all cards   │
                                         │  →  show empty-state │
                                         └──────────────────────┘
```

| Logical state         | `.hot-product-card`         | `#search-empty-state` | `#search-clear` |
|-----------------------|-----------------------------|-----------------------|-----------------|
| empty input           | all `display` default       | `hidden`              | `hidden`        |
| non-empty, ≥1 match   | matching subset visible     | `hidden`              | visible         |
| non-empty, 0 matches  | all `display: none`         | visible (announced)   | visible         |

This table is the source of truth for both the Go markup test (initial states) and the manual browser checks in `quickstart.md` (transitions).
