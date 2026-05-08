# Data Model: Browser-Side Product Search

**Phase**: 1 — Design  
**Feature**: Browser-Side Product Search  
**Date**: 2026-05-08

---

## Overview

This feature introduces no new data model, no new entities, and no new storage. It operates exclusively on data that is already present in the rendered HTML page.

---

## Existing Entities (read-only, no change)

### ProductView (Go struct, frontend handler)

Defined in `src/frontend/handlers.go` as an anonymous struct:

```
ProductView
├── Item  *pb.Product   (id, name, description, picture, price_usd, categories)
└── Price *pb.Money     (currency_code, units, nanos)
```

This struct is passed to the `home` template as `$.products`. It is already fully populated on page load. The search feature reads `.Item.Name` (via the rendered DOM) and performs no mutations.

---

## Client-Side Filter State

The only new "state" introduced is the live value of the search input field. It exists transiently in the browser and is never persisted.

| State | Type | Lifecycle |
|-------|------|-----------|
| `searchQuery` | string | In-memory, browser only. Exists for the duration of the page session. Cleared on navigation or page reload. |

---

## DOM Contract (runtime)

The JavaScript filter relies on the following DOM structure being present. This is a read contract on the existing template output — not a change to it.

| Selector | Role |
|----------|------|
| `#product-search-input` | The new search `<input>` element (added by this feature) |
| `.hot-product-card` | Each product card wrapper; toggled visible/hidden by the filter |
| `.hot-product-card-name` | The element containing the product name text within each card |
| `#no-results-message` | A new element (added by this feature) shown when zero cards match |

No other DOM elements are read or mutated.
