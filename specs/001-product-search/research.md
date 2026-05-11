# Phase 0 Research: Product Search

**Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md)

The spec resolved every open question into a documented Assumption (A-001..A-010) or an explicit Constraint (C-001..C-015). There are **no unresolved NEEDS CLARIFICATION markers**. This file therefore records the implementation-shape decisions the plan rests on, plus the alternatives we considered and rejected.

---

## Decision 1 — Filter on the client, against the already-rendered DOM

**Decision**: implement the filter as page-scoped vanilla JavaScript that listens for `input` and `keydown` events on the search box and toggles each product card's `style.display` based on a case-insensitive substring match against `data-product-name`.

**Rationale**:
- Spec FR-007 (empty-state visible without an extra fetch), C-013 (no per-keystroke catalogue fetches), and the user's brief ("filters the products already loaded, in the browser") together rule out any server-side filtering loop.
- The home page (`src/frontend/handlers.go:59`, template `templates/home.html`) already server-renders the full product list. The data is already in the DOM by the time the user can type — re-fetching is wasteful and forbidden.
- Catalogue size is 9 (per `src/productcatalogservice/products.json` and A-007). A linear `cards.forEach(c => …)` per keystroke is sub-millisecond on any device — well inside SC-003's 100 ms budget. No debouncing, indexing, virtualisation, or Web Worker is needed.

**Alternatives considered**:
- *Server-side filter via a new `SearchProducts` gRPC method on `productcatalogservice`*: would extend the protobuf surface (technically allowed by C-005 as long as we follow existing patterns) but introduces a round-trip per keystroke (violates C-013), changes a service contract (violates C-014 unless the new method is purely additive — but still adds complexity for zero customer-visible win), and adds backend work the brief explicitly says is unnecessary. **Rejected.**
- *Server-side filter on form submit, re-render `/`*: meets the literal "filter products by name" brief but produces a worse UX (no live filtering — violates A-002) and adds a round-trip the spec forbids. **Rejected.**
- *Add a JSON product-data endpoint and re-fetch on page load*: duplicates data already inline in the rendered HTML and adds a fetch for no reason. **Rejected.**

---

## Decision 2 — JavaScript in a separate file under `static/js/`

**Decision**: new file `src/frontend/static/js/search.js`, included via `<script src="{{ $.baseUrl }}/static/js/search.js" defer></script>` inside `home.html` (page-scoped — **not** in `header.html`).

**Rationale**:
- The project currently has zero JavaScript files (`find src/frontend/static -name "*.js"` is empty). This is a meaningful "where does new JS live" decision, not just a file placement.
- CSS already lives in per-area files under `static/styles/`. Mirroring that pattern for JS (`static/js/`) is the lowest-surprise option and keeps lint/test tooling pointable at a known directory if we ever add either.
- Including the `<script>` tag in `home.html` (not `header.html`) means the script is only loaded on the product list page. Every other page (cart, checkout, product detail) is byte-for-byte unchanged in terms of payload.
- `defer` ensures the script runs after the DOM is parsed but doesn't block rendering — search box is visible immediately, becomes interactive a tick later.

**Alternatives considered**:
- *Inline `<script>` block at the end of `home.html`*: smaller diff, no new directory. Rejected: ~50 lines of behaviour wants its own file (testability, future-proofing, single-responsibility); a separate file costs us one `<script>` tag and pays for itself the first time anyone wants to lint or replace it.
- *`<script>` include in `header.html` so it loads everywhere*: violates the "page-scoped" intent and adds bytes to every page on the site. **Rejected.**

---

## Decision 3 — `data-product-name` on each card; JS reads `dataset.productName`

**Decision**: extend each product card in `home.html` from
```html
<div class="col-md-4 hot-product-card">
```
to
```html
<div class="col-md-4 hot-product-card" data-product-name="{{ .Item.Name }}">
```
JS reads `card.dataset.productName.toLowerCase()` and compares to the trimmed, lowercased query.

**Rationale**:
- Keeps the JS DOM-agnostic about how the name is rendered inside the card (today it's `<div class="hot-product-card-name">`, tomorrow possibly different — and we don't want the filter to break).
- Go's `html/template` auto-escapes `{{ .Item.Name }}` in attribute context, so product names containing quotes or angle brackets are handled by the existing renderer — no extra work needed for that side of the XSS surface.

**Alternatives considered**:
- *Read text content from `.hot-product-card-name` div*: brittle to template changes; more JS. **Rejected.**
- *Embed product list as JSON in a `<script type="application/json">` tag and read from JS*: duplicates the data that's already in the DOM. **Rejected.**

---

## Decision 4 — Pre-rendered, hidden empty-state element; static text (does not echo the query)

**Decision**: in `home.html`, render a sibling element next to the grid:
```html
<div id="search-empty-state" class="search-empty-state" role="status" aria-live="polite" hidden>
  No products match your search.
</div>
```
JS removes the `hidden` attribute when the trimmed query is non-empty and zero cards match; restores it otherwise. The text is static and does **not** include the user's query.

**Rationale**:
- `role="status"` + `aria-live="polite"` satisfies FR-007's screen-reader announcement requirement without interrupting the user.
- Pre-rendering and toggling (vs. creating the element in JS) keeps the DOM contract stable, keeps the diff smaller, and means the element is the same whether JS is present or not.
- **Static text — does not echo the query.** This closes the XSS surface (FR-009 / SC-006) cleanly: if there's no echoed user input, there's nothing to escape and nothing that could be executed. A warmer message that echoed the query ("No products match 'zzzzz'") would require `textContent` (not `innerHTML`) on every update; the simpler approach is to not echo at all.

**Alternatives considered**:
- *Echo the query in the message*: warmer UX, but adds an XSS-relevant code path; brief explicitly favoured the simple shape, and FR-007's text is already informative on its own. **Rejected for v1.**
- *Replace the grid's HTML entirely when zero match*: more diff, more state to keep straight. **Rejected** — toggling `hidden` on cards and on one sibling element is enough.

---

## Decision 5 — Visible clear button + Esc keybinding; both empty the input and fire the same code path

**Decision**: wrap the input + button in `<div class="search-wrapper">`. Button is `<button type="button" id="search-clear" class="search-clear" aria-label="Clear search">×</button>` positioned by CSS on the right edge of the input. JS:
- Toggles `search-clear` visibility based on input value length (hidden when empty, visible when ≥ 1 char).
- On click: empties the input, fires a synthetic `input` event so the filter re-runs.
- `keydown` listener on the input: on `Escape`, same as click.

**Rationale**:
- A `<button>` is keyboard-accessible by default — no `tabindex` gymnastics.
- One code path for both interactions (click handler reused from Esc handler) means SC-008 holds for both means of clearing without duplicate logic.
- `aria-label` keeps the `×` glyph announced sensibly to screen readers.

---

## Decision 6 — Go test asserts the markup contract; client-side behaviour is verified manually

**Decision**: add `src/frontend/home_test.go`. It drives `homeHandler` via `httptest` and asserts that the rendered body contains:
- a search input element with `id="search-input"` and a sensible `aria-label`
- a `#search-clear` button
- a `#search-empty-state` element with `role="status"`, `aria-live="polite"`, the `hidden` attribute, and the exact text `No products match your search.`
- a `data-product-name="<expected name>"` attribute on each rendered product card

Client-side behaviour (filter, empty state, clear, Esc, rapid typing, XSS, layout invariance, no extra fetches, no URL change) is verified by walking the spec's acceptance scenarios in a real browser — the `quickstart.md` document spells out the 12 manual checks against the spec's User Stories and Success Criteria.

**Rationale**:
- The frontend has tests only for `validator` and `money` packages (`find src/frontend -name "*_test.go"`); there is no existing Go-side HTML assertion or HTTP-handler test pattern, and there is no JS test framework. Adding Jest / Playwright / Cypress would require new tooling, browser binaries, and CI changes — which the spec rules out (C-006 / C-008).
- The Go test we *can* add is cheap and rides on the existing `go test ./...` invocation. It locks in the markup contract that `search.js` depends on, so a template refactor that accidentally drops `data-product-name` will fail the build, not silently break filtering at runtime.

**Alternatives considered**:
- *Add Playwright/Cypress*: introduces a new test framework, new CI requirements (browser binaries), and a new toolchain. **Rejected** (violates the spirit of "no new infra").
- *Skip the Go test entirely, rely on manual only*: loses regression protection on the markup contract. **Rejected** — the Go test is small (~50 lines) and worth it.

---

## Open questions

None. Every decision above is grounded in the spec, the existing repo structure (`src/frontend/templates/home.html`, `src/frontend/handlers.go:59-110`, `src/productcatalogservice/products.json`), and the spec's constraints.
