# Data Model: Browser-Side Product Search

**Feature**: `specs/001-browser-product-search`
**Date**: 2026-05-11

## Overview

This feature introduces no new server-side data model. All filtering operates on the DOM produced by the existing Go template pipeline. The relevant data entities are the rendered HTML structures already in the page.

---

## DOM Entities

### ProductCard

The unit of filtering. One card exists per product rendered by `{{ range $.products }}` in `home.html`.

| Attribute | HTML | Source |
|---|---|---|
| Container | `div.hot-product-card` | `home.html` loop |
| Name text | `div.hot-product-card-name` | `{{ .Item.Name }}` (server-rendered) |
| Link | `a[href]` wrapping the image | `{{ $.baseUrl }}/product/{{.Item.Id}}` |
| Image | `img[src]` | `{{ .Item.Picture }}` |
| Price | `div.hot-product-card-price` | `{{ renderMoney .Price }}` |

**Filtering key**: `div.hot-product-card-name` → `element.textContent.trim()` (lowercase for comparison)

**Visibility state**: toggled via `element.style.display` (`""` = visible, `"none"` = hidden). Bootstrap grid class `col-md-4` remains on the container; only `display` is overridden inline so the grid reflows correctly.

---

### SearchInput

The new UI element introduced by this feature.

| Attribute | Value |
|---|---|
| Element | `<input type="search">` |
| `id` | `product-search-input` |
| `aria-label` | `"Search products"` |
| `placeholder` | `"Search products…"` |
| Event | `input` (fires on every keystroke) |
| Position in DOM | Immediately before the `div.col-12` containing `<h3>Hot Products</h3>` |

---

### EmptyStateMessage

A single paragraph injected once into the product grid.

| Attribute | Value |
|---|---|
| Element | `<p>` |
| `id` | `no-results-msg` |
| Default state | Hidden (`display: none` via CSS class) |
| Shown when | Zero `.hot-product-card` elements are visible after filtering |
| Text | `"No products match your search."` |

---

## State Transitions

```
[All cards visible, search input empty]
        │
        │ user types character(s)
        ▼
[Cards filtered: matching visible, non-matching hidden]
        │                    │
        │ zero matches        │ ≥1 match
        ▼                    ▼
[empty-state msg shown]   [empty-state msg hidden]
        │
        │ user clears input (value == "" or whitespace only)
        ▼
[All cards visible, empty-state msg hidden]
```

---

## Source Layout

Only two existing files are modified. No new files are created.

```text
src/frontend/
├── templates/
│   └── home.html              ← add search input markup + inline <script>
└── static/
    └── styles/
        └── styles.css         ← add .product-search-box styles + #no-results-msg rule
```
