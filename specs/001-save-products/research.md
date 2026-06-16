# Research: Save products and view them later

Phase 0 research for `001-save-products`. The spec carried no `[NEEDS CLARIFICATION]`
markers; the items below resolve the design choices the constitution constrains.

## Decision 1 — Where saved state lives

- **Decision**: A process-local, in-memory store inside the frontend service: a map
  from session ID to an ordered set of product IDs, guarded by a mutex for concurrent
  request safety.
- **Rationale**: Principle II forbids any new datastore (no DB/cache/search). The cart
  uses a separate `cartservice` backed by Redis; reusing or imitating that would mean a
  new service or datastore, which Principles I and II forbid. An in-memory map in the
  existing frontend is the only compliant home for per-session state.
- **Alternatives considered**:
  - *cartservice-style backend* — rejected: new service / new datastore (violates I, II).
  - *Browser cookie/localStorage holding the saved list* — rejected: pushes state and
    logic into the client beyond the server-rendered pattern of this app, and the
    existing pages are rendered server-side from session context; mismatch with
    Principle III. (Session identity already rides a cookie; the data stays server-side.)

## Decision 2 — Known limitation: process-local state

- **Decision**: Accept that saved products are not shared across frontend replicas, and
  record it as a known limitation rather than engineering around it.
- **Rationale**: Principle II permits only in-memory state. With multiple frontend
  replicas behind a load balancer, a shopper's requests could land on different replicas
  and see different saved sets. This is an inherent consequence of the constraint, not a
  fixable defect within it. For the session-scoped, demo-scale target it is acceptable;
  the training environment runs effectively single-replica per attendee.
- **Alternatives considered**: shared store / sticky sessions — both require infra or a
  datastore change (violate II and IV). Surfaced here so it is a conscious trade, per
  the Governance rule to record tension explicitly.

## Decision 3 — Shopper identity

- **Decision**: Key saved state on the existing session ID via `sessionID(r)`.
- **Rationale**: The app has no accounts (spec out-of-scope; Principle-agnostic). The
  frontend already issues a session cookie in `ensureSessionID` middleware and exposes
  `sessionID(r)`; the cart uses the same identity. Reusing it satisfies "no login" and
  matches existing patterns (Principle III).
- **Alternatives considered**: a new wishlist-specific cookie — rejected as redundant
  and divergent from the established session pattern.

## Decision 4 — Save interaction

- **Decision**: A `POST` form on the product page (hidden `product_id`) to a new
  `/wishlist` route, which records the save and redirects back, mirroring the
  add-to-cart flow (`addToCartHandler`: `r.FormValue("product_id")` → mutate session
  state → `http.Redirect`).
- **Rationale**: Directly imitates the closest existing behaviour, keeping the change
  idiomatic (Principle III) and progressive-enhancement friendly (works without JS).
- **Alternatives considered**: a JSON/AJAX endpoint — rejected for this story as it
  diverges from the server-rendered form pattern; not needed to meet the spec.

## Decision 5 — Resolving product details for the saved view

- **Decision**: Store only product IDs; resolve full product details for the saved view
  through the existing catalogue helper (`fe.getProduct`), skipping any ID no longer in
  the catalogue.
- **Rationale**: Keeps `products.json` (via the catalogue service) the single source of
  product data (Principle II), avoids duplicating product data into the store, and
  handles the "saved product became unavailable" edge case naturally.
- **Alternatives considered**: caching product snapshots in the store — rejected: stale
  data and an implicit cache (brushes against II).

## Decision 6 — Idempotency & empty state

- **Decision**: Saving an already-saved product is a no-op (set semantics, no
  duplicates). The saved view renders an explicit empty-state message with a prompt to
  browse when the session has nothing saved.
- **Rationale**: Satisfies FR-007 (no duplicates) and FR-008 / Acceptance Scenario 4
  (empty state) directly.

## Decision 7 — Testing approach

- **Decision**: Unit-test the store (add/list/contains, idempotency, per-session
  isolation) with Go `testing`; handler-test save and view with `net/http/httptest`.
- **Rationale**: Matches the repo's Go test style and needs no new tooling (Principle
  III/IV). No CI changes required (Principle IV).
