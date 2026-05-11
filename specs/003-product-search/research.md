# Phase 0 Research: Product Search

**Feature**: 003-product-search
**Date**: 2026-05-11

Each entry resolves an item from the plan's Technical Context section or a Deferred item from `/speckit-clarify`.

---

## R1. Where the search input lives, and what page is "the product list page"

- **Decision**: Add the search input on the home page only (`/`), rendered by `src/frontend/templates/home.html`, placed immediately above the existing `.hot-products-row` `<div>`.
- **Rationale**: Online Boutique's home handler (`src/frontend/handlers.go:homeHandler`) is the single place that renders the product grid. Inspecting `home.html` confirms there is exactly one server-rendered grid (under `<div class="row hot-products-row …">`) and that handler always populates `$.products` with all products fetched from the catalogue service. There is no second product list page in the frontend. "The product list page" in the spec maps unambiguously to `home.html`.
- **Alternatives considered**:
  - *A separate `/products` route*: rejected — adds a new handler/route, conflicts with FR-009 ("MUST NOT alter the product list page's behaviour when the search box is empty"), and there is no existing such page.
  - *Recommendation pane (`recommendations.html`)*: rejected — it shows a curated subset, not the catalogue.

## R2. How to filter without a network request

- **Decision**: Run the filter entirely in the browser. On `input` events, iterate over the rendered `.hot-product-card` elements, read each card's `data-name` attribute (added to the template), and toggle the element's `hidden` attribute based on a case-insensitive substring match against the trimmed query.
- **Rationale**: Satisfies FR-007 ("MUST NOT issue any additional network requests… as a result of the shopper typing") and SC-003 ("zero additional network requests") deterministically. With ~10 cards (current catalogue) and even at 50–100 cards, `O(n)` toggle is sub-millisecond, comfortably inside the 100 ms perceived-latency budget (SC-002). Uses no framework — vanilla JS, ES2017 features (`String.prototype.includes`, `Array.from`, arrow functions) all supported by the evergreen browsers Online Boutique targets.
- **Alternatives considered**:
  - *Server-side filter with `?q=…` query param*: rejected — issues a new HTTP request per keystroke, violating FR-007/SC-003.
  - *XHR/fetch to a new endpoint*: rejected — introduces a new endpoint (violates C-007) and a network round-trip (violates FR-007).
  - *`display: none` instead of `hidden`*: equivalent visually; `hidden` is preferred because it semantically removes the card from the accessibility tree, helping screen readers reflect the filtered state without an explicit `aria-hidden` toggle.

## R3. Matching algorithm

- **Decision**: `card.dataset.name.toLowerCase().includes(query.trim().toLowerCase())`. Empty/whitespace-only query → all cards visible. Compare on the **trimmed** query; do not modify card data.
- **Rationale**: Matches FR-004 (case-insensitive substring of name, trim leading/trailing whitespace) and FR-005 (empty query → full grid). Substring `includes` matches multi-word product names correctly ("yoga mat" matches "Yoga Mat"). No locale-specific collation needed for v1 (per Assumptions).
- **Alternatives considered**:
  - *Prefix match*: rejected — narrower than the spec; "lasses" would not surface "Sunglasses".
  - *Tokenized match* (split on space, all tokens must appear): rejected — over-engineered for the spec, no benefit for a ~10-product catalogue.
  - *Fuzzy match (e.g. Levenshtein)*: rejected — out of scope; spec explicitly says "substring".

## R4. Update timing — debounce or per-keystroke?

- **Decision**: Filter synchronously on every `input` event. **No debounce.**
- **Rationale**: FR-003 says "MUST update on every keystroke", and the DOM-toggle cost over the current catalogue is far below the 100 ms budget (SC-002). Debouncing would add lag for no benefit. If the catalogue ever grew large enough to be noticeable, a `requestAnimationFrame`-based throttle can be added later without changing acceptance behaviour.
- **Alternatives considered**:
  - *150–250 ms debounce*: rejected — visible "stale results" period; would fail Story 1 scenario #2 ("`mu` → `mug`" without lag).
  - *Filtering only on Enter*: rejected — fails FR-003.

## R5. Accessibility — live announcement of the filtered count

- **Decision**: Add a single visually-hidden `<div role="status" aria-live="polite" aria-atomic="true">` element next to the input. After each filter update, set its text content to one of:
  - "" (empty) when the query is empty or whitespace.
  - "N products match" when `N >= 1`.
  - "No products match your search." when `N == 0` (this is also the visible empty-state message).
- **Rationale**: This is the **Deferred** UX item from `/speckit-clarify`. The recommended approach (a polite live region) gives screen-reader users the same outcome-feedback sighted users get, with negligible implementation cost. Matches FR-006's intent and meets the spec's "no new accessibility regression" assumption.
- **Alternatives considered**:
  - *No announcement*: rejected — fails to inform assistive-tech users that the visible set changed, which is the whole point of the feature.
  - *`aria-live="assertive"`*: rejected — interrupts the user's current screen-reader output on every keystroke; "polite" coalesces nicely.

## R6. Visible "clear" affordance

- **Decision**: Use the browser-native `type="search"` clear button. Do **not** add a custom × button to the DOM.
- **Rationale**: `<input type="search">` renders a built-in clear control in every evergreen browser (Chrome, Edge, Safari, Firefox), is keyboard-accessible, and emits an `input` event when used — which already triggers our filter update via the same code path. Zero added DOM, zero added CSS-quirks across browsers, zero new acceptance cases. (User Story 2 — clear the box — is satisfied either way.)
- **Alternatives considered**:
  - *Custom × button overlaid on the input*: rejected — extra DOM, focus management, mobile-touch-target tuning; no user benefit beyond the native control.
  - *No clear affordance*: rejected — relies on the shopper holding backspace; the native search-input clear button is free.

## R7. Where to attach the JavaScript

- **Decision**: New file `src/frontend/static/js/product-search.js`, referenced via a `<script defer>` tag added to the bottom of `home.html`. The frontend's existing `http.FileServer` at `/static/` (see `src/frontend/main.go`) already serves the new path with no Go code change.
- **Rationale**: Keeps the change local to the template + a self-contained static asset. No new routes, no Go code, no Docker layer changes. Matches "no backend change" (FR-007) and avoids touching `handlers.go`/`main.go`.
- **Alternatives considered**:
  - *Inline `<script>` in `home.html`*: rejected — harder to read, no caching, slightly larger HTML response per request.
  - *Add the JS via a new template partial*: rejected — over-engineered for one file.

## R8. CSS scope

- **Decision**: Append a small block (~10–20 lines) to the existing `src/frontend/static/styles/styles.css`: layout for `.product-search` (input width / spacing above grid), and the `.product-search-empty` visible empty-state message. Use Bootstrap utility classes where the existing template already uses them; do not introduce a new CSS framework.
- **Rationale**: Reuses the existing stylesheet, no new HTTP request, no risk of style leakage. Bootstrap is already loaded by the page header (visible in `home.html` via classes like `col-md-4`, `d-lg-none`).
- **Alternatives considered**:
  - *Separate CSS file*: rejected — extra HTTP request and asset to maintain for ~15 lines.
  - *Inline `<style>` in the template*: rejected — same readability/cache concerns as inline JS.

## R9. Test strategy

- **Decision**: Verify the feature via the `quickstart.md` browser checklist (each Given/When/Then in the spec mapped to a concrete browser action). Do **not** add a new test harness to the frontend service.
- **Rationale**: The frontend service today contains no JS test infrastructure and no UI-level test framework (e.g. Cypress, Playwright). Introducing one would expand the PR well beyond the spec, contradicting the "stay inside this branch, no pipeline changes" constraint. Server-side Go tests offer nothing to assert about a client-side DOM filter. If automated regression coverage is wanted later, a small Playwright check in a separate PR is the natural follow-up.
- **Alternatives considered**:
  - *Add a Playwright/Cypress E2E suite*: rejected — adds new dev dependency, new CI step, conflicts with C-006.
  - *Pure-JS unit test for the substring/trim helper*: defensible (one tiny `node`-runnable file) but adds a JS toolchain to a Go-only frontend; deferred to a follow-up.

## Open questions

None. All Phase 0 unknowns and the Deferred items from `/speckit-clarify` have been resolved above.
