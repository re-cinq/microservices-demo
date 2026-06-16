# Contracts: Recently Viewed Products

This feature exposes **no new external HTTP API**. It reuses the existing product-page route. The relevant contracts are internal: the store interface and the template data contract.

## 1. HTTP behaviour (existing route, augmented)

**Route**: `GET /product/{id}` (unchanged — `main.go:151`)

| Precondition | Response addition |
|---|---|
| Session has ≥1 recorded product other than `{id}` | Page includes a "Recently viewed" strip at the bottom: ≤4 products, most-recent-first, `{id}` excluded, thumbnails at 50% size. |
| Session has no other recorded products | No strip; page byte-for-byte equivalent to today's layout. |
| A recorded product ID no longer resolves in the catalogue | That product is silently omitted; the rest of the strip and the page render normally (no error). |

**Side effect**: rendering `GET /product/{id}` records `{id}` as most-recently-viewed for the session.

## 2. Store interface contract (internal, Go)

```text
Record(sessionID, productID string)
  - Prepends productID to the session's most-recent-first list.
  - No de-duplication (AIP-201 will add it).

List(sessionID, excludeProductID string, max int) []string
  - Returns up to `max` product IDs for the session, most-recent-first,
    skipping every entry equal to excludeProductID.
  - Unknown session => empty slice.
  - Never returns more than `max` (caller passes 4).
```

## 3. Template data contract

`product.html` receives, in addition to today's keys:

```text
recently_viewed : []*pb.Product   // ordered, ≤4, current excluded, resolved-only; empty => no strip
```

The `recentlyviewed.html` partial renders nothing when `recently_viewed` is empty/absent.

## 4. Acceptance-criteria → contract trace

| AC (spec) | Covered by |
|---|---|
| AC1 (strip at bottom, most-recent-first, current excluded) | HTTP behaviour row 1 + `List` exclude/order |
| AC2 (max 4) | `List(..., max=4)` |
| AC3 (click → product page) | Strip links to existing `/product/{id}` route |
| AC4 (none viewed → unchanged) | HTTP behaviour row 2 + empty-list template guard |
| AC5 (50% thumbnails) | CSS rule on `.recently-viewed` thumbnails |
| FR-010 (missing product) | HTTP behaviour row 3 + skip-on-resolve-failure |
