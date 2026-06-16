# Phase 0 Research: Recently Viewed Products

The spec carried no unresolved `[NEEDS CLARIFICATION]` markers (open questions were closed with documented assumptions). Research therefore confirms the design decisions that satisfy the spec within the epic's hard constraints, grounded in the existing frontend code.

## Decision 1 — Where the recently-viewed list lives

**Decision**: An in-process, mutex-guarded `map[sessionID][]string` (product IDs) held as a new field on the existing `frontendServer` struct (`src/frontend/main.go:63`), initialized in `main()`.

**Rationale**: Constraint C-002 forbids any new datastore (DB/cache/search). An in-process map is not a datastore — it is process memory, the same category as the existing gRPC connection fields. It is the simplest thing that works and matches how the frontend already holds per-process state.

**Alternatives considered**:
- *Cookie / client-side storage*: avoids server memory but pushes state to the browser and complicates rendering; the catalogue lookups still happen server-side. Rejected as more complex with no benefit at demo scale.
- *Redis / external cache*: explicitly forbidden by C-002.

## Decision 2 — Session identity

**Decision**: Key the store on `sessionID(r)`, the existing per-shopper session ID already derived from the `shop_session-id` cookie by the `ensureSessionID` middleware (`src/frontend/middleware.go:85`) and read via the `sessionID` helper (`handlers.go:581`).

**Rationale**: Reuses the established session mechanism (consistent with the spec assumption). No new cookie or identity concept. Note: when `ENABLE_SINGLE_SHARED_SESSION=true`, all shoppers share one session ID — the recently-viewed list is then shared too, which is acceptable demo behaviour and out of scope to change.

**Alternatives considered**: New dedicated cookie — rejected (unnecessary; C-001/simplicity).

## Decision 3 — Recording and ordering

**Decision**: On each product-page render, (a) build the strip from the current stored list for the session, excluding the current product ID, taking the first 4; then (b) record the current product ID at the front of the list. Most-recent-first ordering is achieved by prepending.

**Rationale**: Matches acceptance scenarios 1–2 (most-recent-first, current excluded, max 4). Building before recording keeps the current product out naturally and avoids it ever appearing in its own strip (SC-005).

**De-duplication is intentionally NOT implemented here** — it is AIP-201. Per the spec edge cases, repeated views may produce duplicate entries among the 4 shown. The store list is left unbounded-but-trimmed-on-read; AIP-201 will add move-to-front/dedup in the same file.

## Decision 4 — Hydrating IDs to products

**Decision**: Resolve each stored ID to a full product via the existing `fe.getProduct(ctx, id)` RPC (`src/frontend/rpc.go:51`), exactly as `getRecommendations` does (`rpc.go:99`).

**Rationale**: Reuses the existing catalogue read path (C-001/C-002). No bulk endpoint needed for ≤4 lookups.

**Graceful degradation (FR-010)**: Unlike `getRecommendations` (which errors if any lookup fails), the recently-viewed builder will **skip** an ID that fails to resolve and continue, logging at debug. The strip is non-critical UI and must never break the page if a previously viewed product was removed from the catalogue.

## Decision 5 — Rendering and thumbnail sizing

**Decision**: Add a new template partial `templates/recentlyviewed.html` defined like `recommendations.html`, included at the bottom of `product.html` (after the recommendations include at `product.html:77`), guarded by `{{ if .recently_viewed }}` so an empty list renders nothing (FR-007). Thumbnails get a dedicated CSS class (e.g. `.recently-viewed img`) in `static/styles/styles.css` sized at 50% of the page's normal product thumbnail.

**Rationale**: Reuses the proven strip markup/pattern; the `{{ if }}` guard delivers the "no strip, unchanged layout" requirement for free. CSS-relative sizing (50% of the existing thumbnail rule) keeps the "half size" requirement correct across pages where the base thumbnail differs (relevant for the AIP-200 home-page reuse).

**Alternatives considered**: Fixed pixel size — rejected because "50% of the other thumbnails on the page" is relative and the home page (AIP-200) uses a different base size.

## Decision 6 — Testing approach

**Decision**: Unit-test the store in `recentlyviewed_test.go` (`go test`): most-recent-first ordering, cap at 4, current-product exclusion, behaviour with duplicates (documents the no-dedup contract), and per-session isolation. UI behaviour verified manually per `quickstart.md`.

**Rationale**: The store holds all the testable logic; the constitution mandates no specific test strategy. Pure-function store logic is cheap to cover and protects AIP-200/AIP-201 from regressions.
