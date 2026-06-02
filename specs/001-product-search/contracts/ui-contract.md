# UI Contract: Browser-Side Product Search

**Feature**: `specs/001-product-search`
**Date**: 2026-06-02

This contract defines the HTML interface between the Go template (`home.html`) and the client-side filter script (`product-search.js`). Both sides must honour this contract; changes to either must be reviewed against it.

---

## Contract: Search Input Element

| Property | Value |
|---|---|
| Element | `<input>` |
| ID | `product-search` |
| Type | `search` |
| Event listened to | `input` |

**Rule**: `product-search.js` attaches its listener to `document.getElementById('product-search')`. Renaming or removing this element breaks the filter.

---

## Contract: Product Card Elements

| Property | Value |
|---|---|
| CSS class | `hot-product-card` |
| Required attribute | `data-name` (string, the product's display name) |
| Element tag | `div` |
| Selector used by JS | `.hot-product-card[data-name]` |

**Rule**: Every product rendered in the home page grid MUST have `class="... hot-product-card ..."` and a non-empty `data-name` attribute. The JS toggles `display` on this element to show/hide each product.

**Rule**: The JS MUST NOT modify any attribute other than `style.display` on `.hot-product-card` elements. All other card content (image, name text, price, link) is outside the JS's responsibility.

---

## Contract: No-Results Message Element

| Property | Value |
|---|---|
| Element | `div` |
| ID | `no-results-message` |
| Default state | `display: none` (set inline in the template) |

**Rule**: `product-search.js` sets `style.display` to `'block'` when zero cards match and the query is non-empty, and back to `'none'` otherwise. Renaming this ID breaks the no-results display.

---

## Contract: Script Load Order

The script `product-search.js` is loaded after Bootstrap JS (which is in the page footer). It requires the DOM to be fully parsed. It MUST be loaded at the end of `<body>` (or with `defer`).

---

## What the JS must NOT do

- Must NOT modify any element's class list, text content, or attributes other than `style.display`
- Must NOT make any network requests (fetch, XMLHttpRequest)
- Must NOT interact with any element outside `.hot-product-card` elements, `#product-search`, and `#no-results-message`
- Must NOT depend on any CSS class name other than `hot-product-card`
