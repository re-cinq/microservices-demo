# Phase 0 — Research: Browser-Side Product Search

**Feature**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

This document resolves the open questions that surfaced when filling the plan's Technical Context. Each entry follows: **Decision → Rationale → Alternatives considered**. All NEEDS CLARIFICATION items from the plan template are resolved below; none remain.

---

## R1. Where does the filter live — server, edge, or browser?

**Decision**: Pure browser-side filter.

**Rationale**: The spec is explicit (FR-007, SC-006): no additional server requests per keystroke; the page already renders the whole catalogue. The Out-of-Scope section forbids new services and datastores. A browser-side filter over the already-rendered DOM is the only option that satisfies all three of those constraints simultaneously.

**Alternatives considered**:

- **Server-side filter via a new `frontend` route (e.g. `/?q=…`)** — would require a round-trip per keystroke or debounce, violates "no additional server requests per keystroke" (SC-006), and adds a code path on the server. Rejected.
- **Extend `productcatalogservice.SearchProducts` and call it from the browser via the frontend** — the existing `SearchProducts` RPC already does substring matching across name+description. Using it would mean network calls per keystroke (violates SC-006), would still require a frontend handler change, and would conflate this feature with the separate `search-misses-descriptions` bug. Explicitly out of scope per spec. Rejected.
- **Edge / CDN-level filter** — would require new infra. Forbidden by the Out-of-Scope section. Rejected.

---

## R2. How does the JS find the products to filter?

**Decision**: Augment each product card in `home.html` with a `data-name` attribute (the lowercased product name), and a stable class hook. The JS queries `document.querySelectorAll('.hot-product-card')` (or whatever class the template already uses) and toggles `hidden` based on a substring test against `data-name`.

**Rationale**: The template already loops `{{ range $.products }}` inside `.hot-products-row`. Adding a small data attribute per card is the smallest possible template change and avoids reparsing the visible text of each card (which is fragile against future markup changes). Lowercasing at render time makes the per-keystroke comparison a single `String.prototype.includes` call.

**Alternatives considered**:

- **Read the product name from the visible text node of each card** — works but is fragile (e.g. if a card later wraps the name in extra spans, the search breaks). Rejected.
- **Inject a JSON blob of all products and filter against that, then re-render** — over-engineered for ~10 items, duplicates data already in the DOM, and adds a re-render step that has no upside here. Rejected.
- **Use `:has-text()` / browser-native CSS selectors** — not portable; `:has-text` is not standard. Rejected.

---

## R3. How is "name match" defined precisely?

**Decision**: Trim the query of leading/trailing whitespace, lowercase it, and require `dataName.includes(trimmedLoweredQuery)`. Empty trimmed query → show all (no filter). No regex, no wildcards, no fuzzy match.

**Rationale**: This exactly matches FR-003 ("contains the query as a substring"), FR-004 (case-insensitive), FR-005 (name only), FR-006 (trim whitespace), and the edge cases for special characters and whitespace-only queries. Substring `includes` cannot throw on special characters, which sidesteps the "regex injection" failure mode.

**Alternatives considered**:

- **Prefix match (`startsWith`)** — narrower than what shoppers expect; "tank" would not find "Tank Top" if a shopper typed "top". Rejected.
- **Fuzzy / Levenshtein** — explicitly forbidden ("No autocomplete, fuzzy / typo-tolerant matching" in Out-of-Scope). Rejected.
- **Regex from the user's input** — risk of `SyntaxError` on `(` or `\\`, and risk of catastrophic-backtracking patterns. Rejected.

---

## R4. How is the "no results" state shown?

**Decision**: A single empty-state element (`<p class="no-results" hidden>No products match your search.</p>`) is rendered alongside the product grid and toggled visible when, after a filter pass, every card has been hidden. The grid itself remains in the DOM (just empty); the search box keeps the shopper's query.

**Rationale**: Matches FR-008 verbatim and avoids removing/re-inserting DOM nodes. The empty-state element exists at all times (with `hidden`) so the JS only flips a boolean — no `innerHTML` churn, no XSS surface even if the query contained markup-like characters.

**Alternatives considered**:

- **Render the message dynamically including the user's query** ("No matches for `<query>`") — pulls user input into the DOM and would need escaping. The spec doesn't require echoing the query inside the message; the query is already visible in the search box. Rejected for v1 (lower risk).
- **Hide the entire `.hot-products-row` and show a full-page message** — bigger visual jolt; FR-008 just says "replace the product grid with a short message". Rejected.

---

## R5. What triggers an update — `input`, `keyup`, or `change`?

**Decision**: `input` event on the search box. No debounce.

**Rationale**: `input` fires on every value change (typing, pasting, IME composition end). It is the standard event for live-filter UIs. With ~10 products and a `querySelectorAll + includes` loop, a debounce is unnecessary and would actively hurt the 200 ms SC-002 budget by adding artificial delay. The browser will repaint after the toggle pass; total wall-clock cost is dominated by repaint, which is well under 200 ms for this DOM size.

**Alternatives considered**:

- **Debounce 100–150 ms** — adds latency for no gain at current scale. Rejected.
- **`keyup`** — misses paste and IME compose; worse than `input`. Rejected.
- **`change`** — only fires on blur. Wrong UX (would require the shopper to press Tab). Rejected.

---

## R6. Accessibility — what does the search box need to be a "good citizen"?

**Decision**:

- `<input type="search">` (not `type="text"`) so browsers give a clear-button affordance for free.
- Associated `<label for="product-search">Search products</label>` (visually shown or visually-hidden via a screen-reader-only class — either is acceptable; the template already has the pattern of using Bootstrap utility classes).
- `aria-controls` pointing at the grid container, and the empty-state element marked `role="status" aria-live="polite"` so screen-reader users hear "No products match your search." when applicable. (Live announcement is a nice-to-have per spec's edge-case note; we include it because it costs nothing.)

**Rationale**: Satisfies FR-011 and SC-007 with no new dependencies. `type="search"` gives keyboard users an Escape-to-clear behaviour in most browsers, which doubles as the User Story 2 reset affordance.

**Alternatives considered**:

- **`<input type="text">` with custom clear button** — more code, no real benefit. Rejected.
- **`role="searchbox"` on a `div`** — never as good as a native `<input>`. Rejected.

---

## R7. Where to put the JS file?

**Decision**: New file `src/frontend/static/js/product-search.js`, referenced from `home.html` via a plain `<script defer src="/static/js/product-search.js"></script>` near the bottom of the home template (inside the `{{ define "home" }}` block, before `{{ template "footer" . }}`).

**Rationale**: The frontend already serves `/static/...` via `http.FileServer` (`main.go` registers `/static/`). Adding a file under `static/` makes it served automatically — no Go change required. `defer` ensures the DOM is ready before the script runs, without needing `DOMContentLoaded` boilerplate.

**Alternatives considered**:

- **Inline the script in the template** — works but is harder to test and lint, and mixes concerns. Rejected.
- **Add a module bundler / npm step** — vastly disproportionate, and would touch the build pipeline (forbidden). Rejected.

---

## R8. How is the feature tested without touching CI?

**Decision**: Two layers:

1. **Manual acceptance walk-through** in `quickstart.md`, one numbered step per Given/When/Then scenario in `spec.md` (8 scenarios across 3 stories). The reviewer runs the frontend locally (or on the deployed `attendee/daniel-tufvander` URL), opens the home page, and ticks each step.
2. **Pure-function JS check** — `product-search.js` exports its match function (`export function matches(name, query)` or attaches it to `window.__productSearch` for testability) so a developer can paste the assertions from `quickstart.md` into the browser DevTools console and confirm the function's logic in isolation.

No automated test harness is added; the constraint "do not touch the build pipeline" forecloses adding a JS test runner to CI. The single Go file we'd otherwise add a `go test` for is **not** modified, so no new Go test is needed either.

**Rationale**: Honours the no-pipeline-change constraint while still giving every scenario a deterministic verification step.

**Alternatives considered**:

- **Add Jest / Vitest** — requires a Node toolchain, a new `package.json`, and ideally a CI step. Forbidden. Rejected.
- **Cypress / Playwright e2e** — requires CI integration. Forbidden. Rejected.

---

## Summary

All Technical Context items are resolved with concrete decisions. No `[NEEDS CLARIFICATION]` markers remain. The design is intentionally minimal:

- **1 template edited** (`home.html`)
- **1 new JS file** (`static/js/product-search.js`)
- **1 stylesheet edited** (`styles.css`)
- **0 Go files edited**
- **0 protobuf changes**
- **0 infrastructure changes**
- **0 CI changes**
