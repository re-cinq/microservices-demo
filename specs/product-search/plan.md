# Implementation Plan: Product Search (client-side name filter)

**Branch**: `attendee/andrea-backstrom` | **Date**: 2026-06-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/product-search/spec.md`

## Summary

Add a search box to the Online Boutique product list page (`home.html`) that filters the
product cards **already rendered on the page** by product name, entirely in the browser.
No backend call, no new service, no datastore, no layout change, no new results page.
Typing narrows the visible `.hot-product-card` set in real time; clearing restores the full
grid; a no-match query shows an inline "no products match" message.

The whole feature is delivered in the **frontend service only** (`src/frontend`), as one
template edit plus one static JavaScript file (and a small CSS addition). Because static
assets are auto-served by the existing file server, **no Go code changes are required**.

## Technical Context

**Language/Version**: Go 1.x (frontend service) for templates/serving; vanilla browser
JavaScript (ES5/ES6, no framework) for the filter logic. No Go logic change needed.

**Primary Dependencies**: Existing only — Go `html/template`, Gorilla Mux router, Bootstrap
4.1.1 (already loaded via CDN in [footer.html](../../src/frontend/templates/footer.html) /
[header.html](../../src/frontend/templates/header.html)). No new dependencies.

**Storage**: N/A. The catalogue is the existing in-memory data already rendered into the page
by the `{{ range $.products }}` loop (sourced upstream from
`productcatalogservice/products.json`). The feature reads the names already present in the DOM.

**Testing**: Manual acceptance via the Given/When/Then scenarios in the spec (browser +
network panel). Optional lightweight DOM-level JS assertions; no existing JS test harness in
this service, so manual verification against `products.json` names is the primary gate.

**Target Platform**: Browser (the rendered Online Boutique home page), served by the Go
frontend container deployed by CI to `attendee/andrea-backstrom`.

**Project Type**: Web application — change confined to the frontend service's presentation
layer.

**Performance Goals**: Filter result visible within one rendered frame of a keystroke
(SC-003); zero network requests on input (SC-004). Catalogue is small (single-digit/low-tens
of products), so a linear DOM scan per keystroke is trivially fast.

**Constraints**: Honor spec constraints C-001…C-010 — frontend-only, in-memory/in-DOM filter,
no new service/datastore, no Helm/manifest/env changes, no build-pipeline changes, no layout
change, no auth, inline results (no separate search page).

**Scale/Scope**: One page (`home.html`), one new static JS file, one CSS block, one template
edit. No proto, no gRPC, no other service.

## Constitution Check

*GATE: Must pass before design. Re-check after design.*

The project constitution at `.specify/memory/constitution.md` is the **unpopulated template**
(placeholder tokens only) — it defines no concrete principles or gates. There are therefore no
constitutional rules to violate. Self-imposed gates derived from the spec's constraints:

| Gate | Status |
|------|--------|
| Frontend-only; no new service (C-001) | PASS — change limited to `src/frontend` |
| No new datastore / search engine (C-002) | PASS — filters the DOM already rendered |
| Uses existing in-memory catalogue (C-003) | PASS — reads names already in the page |
| Matches service language / existing patterns (C-004) | PASS — Go template + static JS, no proto change |
| No infra/Helm/manifest/env changes (C-005) | PASS — only template + static assets |
| Stays in branch/repo; no pipeline change (C-006) | PASS |
| No layout change (C-008) | PASS — search box inserted above grid; grid markup untouched |
| No authentication (C-009) | PASS |
| Inline results, no separate page (C-010) | PASS — filters cards in place |

**Result: PASS.** No complexity to track.

## Project Structure

### Documentation (this feature)

```text
specs/product-search/
├── spec.md              # Feature specification (done)
├── plan.md              # This file (/speckit.plan output)
└── tasks.md             # /speckit.tasks output (next step — NOT created here)
```

### Source Code (repository root)

Change is confined to the frontend service:

```text
src/frontend/
├── main.go                         # UNCHANGED — static file server already serves /static/* (main.go:159)
├── templates/
│   └── home.html                   # EDIT — add search input above the product grid + <script> include
└── static/
    ├── js/
    │   └── product-search.js       # NEW — client-side filter logic (no framework)
    ├── styles/
    │   └── styles.css              # EDIT (small) — styling for search box + no-results message
    └── icons/
        └── Hipster_SearchIcon.svg  # EXISTING — reuse as the search box icon
```

**Structure Decision**: Web application, presentation-layer-only change inside the existing
`src/frontend` service. No backend, proto, or other service is touched. A new static JS file
is automatically served by the existing `http.FileServer(http.Dir("./static/"))` route
([main.go:159](../../src/frontend/main.go#L159)), so the only wiring needed is a `<script>`
tag in the template. `productcatalogservice` is **not** modified.

## Design Decisions

### D1 — Filter the DOM, not the data

The home template already renders every product as a `.hot-product-card` containing a
`.hot-product-card-name` div ([home.html:46-57](../../src/frontend/templates/home.html#L46)).
The JS reads each card's name text and shows/hides the card via a CSS class. This satisfies
"filter the products already loaded, in the browser" (FR-003, FR-008) with no data fetch.

- **Match rule**: `card.name.toLowerCase().includes(query.trim().toLowerCase())` — case-insensitive
  substring on the name only (FR-003, FR-004, FR-005). `trim()` handles surrounding whitespace.
  `includes` treats the query as a literal string, so special characters can't error or act as
  patterns (Edge Cases).
- **Empty query**: `trim() === ""` → show all cards, hide the no-results message (FR-006).

### D2 — Show/hide via a class, never remove nodes

Toggle a `d-none`-style hidden class (Bootstrap's `d-none` is already available) rather than
removing/re-inserting DOM nodes. This guarantees the original order is preserved on restore
(FR-009) and that product links, prices, images, and recommendations are untouched (FR-010,
SC-007). No layout change to the grid markup itself (C-008).

### D3 — Inline "no results" message

A single hidden message element (e.g. `#product-search-no-results`) sits inside the existing
grid container. When the visible-card count hits zero for a non-empty query, the JS unhides it;
otherwise it stays hidden (FR-007, SC-006). Rendered in place — no separate results page (C-010).

### D4 — Real-time on `input` event

Bind one `input` listener on the search field so filtering runs on every keystroke/paste with
no submit button and no page reload (FR-002, SC-001). Synchronous DOM work only → result within
one frame (SC-003), zero network traffic (SC-004).

### D5 — Search box placement (no layout disruption)

Insert the search input as a new row **above** the `hot-products-row` grid, inside the existing
`col-12 col-lg-12` container, reusing `Hipster_SearchIcon.svg`. The existing grid `<div>`s and
their Bootstrap classes are left exactly as they are, satisfying "do not change the existing
product list layout" (C-008).

### D6 — No persistence / no URL change

Query is not stored and does not alter the URL; the box starts empty on every load (Assumptions).
Keeps scope minimal and avoids any server-side or history wiring.

### Why NOT a backend / productcatalogservice change

The spec's primary directive is browser-side filtering with no backend change (FR-008, C-007).
A `productcatalogservice` gRPC `SearchProducts`-style path would add proto changes, a new RPC,
frontend wiring, and a round-trip per keystroke — violating "no backend change" and adding
latency against SC-003/SC-004 for zero benefit at this catalogue size. Rejected. (If the feature
were ever re-scoped to server-side search, C-003/C-004 would bind: Go + existing protobuf/gRPC
patterns, in-memory filter over the loaded `products.json`.)

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Card name text includes extra whitespace/markup | Read `textContent` and `trim()` it before matching |
| Future pagination/lazy-load would hide some products from the DOM | Out of scope today (Assumptions); home renders the full catalogue. Note for re-scope. |
| CSS class collision with Bootstrap | Use a feature-scoped class/id prefix (`product-search-*`) and reuse `d-none` only for hide |
| Static JS not cache-busted | Acceptable for training scope; file is new so no stale-cache concern on first deploy |

## Verification Strategy (maps to spec)

Run the app and walk the spec's acceptance scenarios in the browser:

1. Type a known fragment → only matching cards visible (US1 #1, SC-002).
2. Clear box → full grid restored in original order (US1 #2, SC-005, FR-009).
3. Uppercase / padded query → still matches (US1 #3-4, FR-004/005).
4. Watch the network panel while typing → zero requests (US1 #5, SC-004).
5. Type "zzzzz" → inline no-results message, no cards (US2 #1, SC-006); clear → message gone (US2 #2).
6. Click a filtered product, check price/image/recommendations → unchanged (FR-010, SC-007).

## Next Step

Run **`/speckit.tasks`** to generate the dependency-ordered `tasks.md` from this plan.
