# Data Model: Save products and view them later

All structures are in-memory in the frontend process (Principle II). No schema, no
migration, no persistence.

## Entities

### SavedCollection (per session)

The set of products one shopper has saved during their current session.

| Field | Type | Notes |
|-------|------|-------|
| sessionID | string | Key; from `sessionID(r)` (existing session cookie). Not stored as a field so much as the map key. |
| productIDs | ordered set of string | Product IDs the shopper has saved, in save order (most recent first for display). Set semantics → no duplicates (FR-007). |

**Behaviour**
- `Add(sessionID, productID)` — adds the ID if absent; no-op if already present (idempotent).
- `List(sessionID)` — returns the saved product IDs in display order; empty slice if none.
- `Contains(sessionID, productID)` — whether a product is saved (drives the product-page saved/unsaved indicator, FR-002).

**Validation rules**
- `productID` must be non-empty; ignore empty form values.
- No fixed maximum is enforced (Assumptions: no fixed limit in this scope).

### Store (container)

The process-wide holder of all sessions' collections.

| Field | Type | Notes |
|-------|------|-------|
| collections | map[sessionID]orderedSet | In-memory only; process-local. |
| mutex | lock | Guards concurrent request access. |

**Lifecycle / state transitions**
- A collection is created lazily on first save for a session.
- Entries live for the life of the session/process; cleared when the process restarts
  or the session ends (session-scoped per spec; precise window left to implementation).
- No removal operation in this story (removal is AIP-194).

## Relationships

- `Store` 1 ── many `SavedCollection` (one per active session ID).
- A `SavedCollection` references products **by ID only**; full product details are
  resolved at render time via the existing product catalogue (`fe.getProduct`), so the
  catalogue (`products.json`) stays the single source of product data.

## Derived values

- **Saved count** (used by AIP-195, not this story) would be `len(List(sessionID))`.
  This story renders only the saved view and the header link, not the count.
