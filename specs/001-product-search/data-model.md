# Data Model: Product Search

**Phase 1 output for**: `specs/001-product-search/plan.md`  
**Date**: 2026-05-11

---

## Entities

### Product *(existing, unchanged)*

The unit of data flowing through the search feature. Defined in
`protos/demo.proto` and generated into both services' `genproto/` packages.

| Field | Type | Role in search |
|---|---|---|
| `id` | string | Not searched; used for product page links |
| `name` | string | **Searched** (case-insensitive substring) |
| `description` | string | **Searched** (case-insensitive substring) |
| `picture` | string | Displayed in result card; not searched |
| `price_usd` | Money | Converted to user currency for display; not searched |
| `categories` | []string | Not searched |

No schema changes required. The proto definition is authoritative.

---

### SearchQuery *(ephemeral — not persisted)*

Represents the shopper's intent for a single request cycle.

| Attribute | Source | Validation rule |
|---|---|---|
| `raw` | `r.URL.Query().Get("q")` | May be empty, any string |
| `trimmed` | `strings.TrimSpace(raw)` | Used for length check |
| `normalised` | `strings.ToLower(trimmed)` | Passed to `SearchProducts` RPC |
| `valid` | `len([]rune(trimmed)) >= 2` | If false, search is blocked |

Not stored anywhere. Lives only within a single HTTP request.

---

### SearchResult *(ephemeral — not persisted)*

The filtered subset of products returned by `SearchProducts`. Represented in
the frontend as `[]productView` (same struct used by the full catalogue today),
augmented with the search query for highlighting.

| Attribute | Source |
|---|---|
| products | `SearchProductsResponse.Results` (converted to `[]productView`) |
| query | The validated search term, passed to template for highlighting |
| isEmpty | `len(products) == 0` |

---

## State Transitions

```
  [no ?q param]          [?q absent or cleared]
       │                          ▲
       ▼                          │
  Full catalogue ──────────────────────────
       │                         clear
       │ user types ≥2 chars
       ▼
  Search submitted
       │
       ├── [trimmed len < 2] ──► Validation error
       │                         (full catalogue still shown)
       │
       ├── [len ≥ 2, results found] ──► Result list
       │                               (products highlighted)
       │
       └── [len ≥ 2, no results] ──► Empty state
                                      (message includes query)
```

---

## Validation Rules

| Rule | Enforcement point | Behaviour on failure |
|---|---|---|
| Query must be ≥ 2 non-whitespace runes | `homeHandler` (server-side) + HTML `minlength="2"` (client hint) | Render home template with `searchError` set; product list unchanged |
| Query is treated as a literal string (no regex) | `highlightTerm` function uses `strings.Index`, not `regexp` | N/A — by design |
| HTML special characters in query are escaped | `html.EscapeString` inside `highlightTerm` | Prevents XSS; escaped output is still highlighted correctly |
