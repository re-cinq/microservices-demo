# Data Model: Frontend Product Search Filter

**Feature**: 002-product-search-filter
**Date**: 2026-05-11

---

## Overview

This feature operates entirely on the rendered DOM. There is no new persistent data, no new server-side model, and no schema migration. The "data model" described here is the browser-side representation — the DOM structure and attributes that the filter reads and mutates.

---

## Entities

### ProductCard (DOM element)

Represents one product tile on the home page grid.

| Attribute / Property | Type | Source | Role |
|---|---|---|---|
| `data-product-name` | `string` | Go template (`.Item.Name`) at render time | Filter target; compared against the search query |
| `style.display` | `''` \| `'none'` | Mutated by filter JS | Controls visibility; `''` = visible (default), `'none'` = hidden |
| CSS class `.hot-product-card` | static | Template | Selector used by JS to find all cards |

**Validation rules**:
- `data-product-name` must be present on every `.hot-product-card` element; absence causes the card to be treated as non-matching for any non-empty query (safe fallback).
- `style.display` is the only property mutated by the filter; all other attributes and classes are read-only from the filter's perspective.

**State transitions**:

```
[visible: style.display = '']
        |
        | query entered, name does not match
        v
[hidden: style.display = 'none']
        |
        | query cleared, or name starts to match
        v
[visible: style.display = '']
```

---

### SearchInput (DOM element)

The text field the shopper types into.

| Attribute / Property | Type | Value | Role |
|---|---|---|---|
| `id` | `string` | `product-search` | Referenced by JS event listener |
| `type` | `string` | `search` | Enables native browser clear button |
| `aria-label` | `string` | `"Search products"` | Accessibility label |
| `.value` | `string` | User input | Current query; read on every `input` event |

**State**: Stateless between page loads. Value resets to `''` on navigation.

---

### NoResultsMessage (DOM element)

A feedback element displayed only when the filtered set is empty.

| Attribute / Property | Type | Value | Role |
|---|---|---|---|
| `id` | `string` | `product-no-results` | Referenced by JS to toggle visibility |
| `style.display` | `''` \| `'none'` | `'none'` (default) | Controls visibility |

**State transitions**:

```
[hidden: style.display = 'none']   ← default
        |
        | visibleCount === 0 after filter pass
        v
[visible: style.display = '']
        |
        | visibleCount > 0, or query cleared
        v
[hidden: style.display = 'none']
```

---

## Filter Algorithm

```
on input event:
  query ← searchInput.value.trim().toLowerCase()
  visibleCount ← 0

  for each card in ProductCard[]:
    name ← card.getAttribute('data-product-name').toLowerCase()
    if query === '' or name.indexOf(query) !== -1:
      card.style.display = ''
      visibleCount++
    else:
      card.style.display = 'none'

  noResultsMessage.style.display = (visibleCount === 0) ? '' : 'none'
```

**Complexity**: O(n) where n = number of product cards (~10). No optimisation needed.

---

## Contracts

No external interface contract is exposed. The filter is a purely client-side DOM mutation triggered by user interaction. It does not add any route, endpoint, event, or API. The `/contracts/` directory is omitted for this feature.
