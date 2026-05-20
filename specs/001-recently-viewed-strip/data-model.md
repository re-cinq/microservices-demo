# Data Model: Recently Viewed Products Strip

**Feature**: `specs/001-recently-viewed-strip`
**Date**: 2026-05-20

---

## Entities

### RecentlyViewedList (runtime value, not persisted server-side)

Represents the ordered list of product IDs a shopper has viewed in the current session, as carried by the browser cookie.

| Field | Type | Description |
|-------|------|-------------|
| `productIDs` | `[]string` | Ordered slice of product IDs, most-recent first; max 5 entries |

**Invariants**:
- No duplicate IDs (enforced at write time via filter-before-prepend)
- Length ≤ 5 (enforced at write time via truncation)
- Each ID is a non-empty string matching an entry in `productcatalogservice/products.json` at the time of viewing (IDs from removed products are silently dropped at read time)

**Storage**: browser cookie `shop_recently_viewed`, pipe-separated (`|`) string of IDs, e.g.:
```
OLJCESPC7Z|66VCHSJNUP|L9ECAV7KIM
```

---

### RecentlyViewedProduct (template view model)

A resolved, display-ready product for rendering in the strip. Derived at request time by fetching each ID in `RecentlyViewedList` via the existing gRPC `GetProduct` RPC and applying currency conversion.

| Field | Type | Source |
|-------|------|--------|
| `ID` | `string` | `pb.Product.Id` |
| `Name` | `string` | `pb.Product.Name` |
| `Picture` | `string` | `pb.Product.Picture` (static asset path) |
| `Price` | `pb.Money` | `pb.Product.PriceUsd` after `convertCurrency` |

This is not a new Go struct — it is the existing `*pb.Product` type already used throughout the frontend, supplemented by `pb.Money` after currency conversion. The template receives `[]*pb.Product` with prices already converted, matching the pattern used in `homeHandler`.

---

## Cookie Specification

| Attribute | Value |
|-----------|-------|
| Name | `shop_recently_viewed` |
| Value | Pipe-separated product IDs (e.g., `OLJCESPC7Z\|66VCHSJNUP`) |
| MaxAge | 86400 seconds (24 hours) |
| HttpOnly | `false` |
| SameSite | `http.SameSiteLaxMode` |
| Secure | matches existing cookie settings (off in dev, on behind TLS proxy) |

---

## State Transitions

```
[No cookie / empty list]
        │
        │  shopper loads /product/{id}
        ▼
[cookie = "ID_A"]                 (list has 1 entry)
        │
        │  shopper loads /product/{id_b}
        ▼
[cookie = "ID_B|ID_A"]            (list has 2 entries, newest first)
        │
        │  shopper views 3 more distinct products
        ▼
[cookie = "ID_E|ID_D|ID_C|ID_B|ID_A"]   (list at capacity: 5)
        │
        │  shopper loads /product/{id_f}
        ▼
[cookie = "ID_F|ID_E|ID_D|ID_C|ID_B"]   (ID_A dropped, oldest evicted)
        │
        │  shopper views ID_C again
        ▼
[cookie = "ID_C|ID_F|ID_E|ID_D|ID_B"]   (ID_C moved to front, no duplicate)
```

---

## Template Data Shape

The product detail page handler will inject a new key into the existing `injectCommonTemplateData` map:

```go
map[string]interface{}{
    // ... existing keys (session_id, user_currency, etc.) ...
    "recently_viewed": []*pb.Product,  // may be empty slice, never nil
}
```

The `recently_viewed.html` partial iterates over this slice and is a no-op when the slice is empty.
