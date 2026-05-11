# Phase 1 — Data Model: Browser-Side Product Search

**Feature**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

This feature touches almost no data. There is one **existing** entity it reads, and one **new transient** entity that lives only in browser memory. No persistence, no schema migration.

---

## Existing entity (read-only for this feature)

### Product

Already defined in `src/productcatalogservice/products.json` and surfaced to the frontend via the existing `ListProducts` gRPC call. **This feature reads it; it does not modify it.**

| Field | Type | Used by this feature? | Notes |
|---|---|---|---|
| `id` | string | No | Untouched. |
| `name` | string | **Yes — the only field this feature matches against.** | FR-005: only `name` is used to decide whether a card is shown or hidden. |
| `description` | string | No | Explicitly excluded from matching (Out-of-Scope: "DO NOT search by description … in v1"). |
| `picture` | string | No | Untouched. |
| `priceUsd` | object | No | Untouched. |
| `categories` | array<string> | No | Explicitly excluded from matching. |

**Validation rules from requirements that apply at read time**:

- The lowercased `name` is exposed to the browser via a `data-name` attribute on each product card (see `contracts/ui-contract.md`). The template lowercases it at render time so the JS does one less operation per keystroke.

**State transitions**: none. Products are static within a page load.

---

## New transient entity

### SearchQuery

The text the shopper has currently typed into the search box. It lives only in the current browser tab.

| Field | Type | Notes |
|---|---|---|
| `raw` | string | Exactly what the shopper has typed, including whitespace. Held in `inputElement.value`. |
| `normalized` | string | `raw.trim().toLowerCase()`. This is what gets compared against each product's `data-name`. Recomputed on every `input` event. |
| `isEmpty` | boolean | `normalized.length === 0`. When true, every product card is shown and the no-results message is hidden. |

**Lifecycle**:

```
[page load]
   ↓
SearchQuery exists with raw = "" (empty)
   ↓
[shopper types]
   ↓
input event fires → raw updated → normalized recomputed → filter pass runs
   ↓
[shopper reloads]
   ↓
SearchQuery is destroyed and recreated empty (FR-010, SC-008)
```

**Persistence**: none. Not written to `localStorage`, `sessionStorage`, a cookie, the URL, or the server. The query is intentionally lost on reload per FR-010.

**Validation rules**:

- No length limit is enforced (the input is a substring matcher; long queries simply match nothing).
- No character is special. Special characters (`*`, `?`, `(`, `"`, `\\`) are treated as literal substring characters per the spec's edge-case list.
- Trimming and lowercasing happen exactly once per keystroke, inside `normalize(raw)`.

---

## Relationships

```
            (existing, server-rendered)
+---------+    1..N    +-----------------------+
| Product | ---------> | Rendered product card |
+---------+            | (.hot-product-card)   |
                       | data-name = lower(P.name)
                       +-----------------------+
                                  ^
                                  | filtered by
                                  |
                       +-----------------------+
                       | SearchQuery (transient,
                       | browser-only)         |
                       +-----------------------+
```

There is no relationship to the server side beyond the one that already exists (`frontend` renders product cards from `productcatalogservice.ListProducts`). The `SearchQuery` entity has no server-side representation by design.
