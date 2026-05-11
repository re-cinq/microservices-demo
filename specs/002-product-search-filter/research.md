# Research: Frontend Product Search Filter

**Feature**: 002-product-search-filter
**Date**: 2026-05-11
**Status**: Complete — no NEEDS CLARIFICATION items remain

---

## Decision Log

### 1. Filter trigger event

**Decision**: Listen on the `input` event, not `keyup` or `keydown`.

**Rationale**: `input` fires for every content change regardless of how it happens — keyboard, paste, cut, browser auto-fill, and the native `<input type="search">` clear button. `keyup` misses paste and the clear button; `keydown` fires before the value updates.

**Alternatives considered**:
- `keyup`: Missed paste events and the clear (×) button on `<input type="search">`.
- Debounce + `keyup`: Unnecessary for a DOM-only filter with ~10 cards; debounce adds perceived lag with no benefit at this scale.

---

### 2. Show/hide mechanism

**Decision**: Toggle `element.style.display` between `''` (restore to CSS default) and `'none'`.

**Rationale**: FR-005 requires cards to be *hidden*, not removed from the DOM. Toggling `display` satisfies this and is reversible in a single property assignment. It also preserves any existing Bootstrap grid classes and event listeners on child elements.

**Alternatives considered**:
- `visibility: hidden` / `opacity: 0`: Hides visually but the card still occupies space, collapsing the grid awkwardly.
- Removing and re-inserting DOM nodes: Violates FR-005; also more complex (requires holding a reference list).
- CSS class toggle (e.g., `.hidden { display: none }`): Equivalent outcome but requires an extra CSS rule; inline style is self-contained with no new stylesheet dependency.

---

### 3. Match algorithm

**Decision**: Case-insensitive substring match — `name.toLowerCase().indexOf(query.toLowerCase()) !== -1`.

**Rationale**: Matches FR-003 exactly: case-insensitive, substring-based. Simple to reason about, no regex edge cases with special characters, no external library.

**Alternatives considered**:
- `String.prototype.includes()`: Equivalent result; `indexOf !== -1` chosen for ES5 compatibility with the project's existing template JS style.
- Regex match: More powerful but introduces edge-case risk when the query contains regex metacharacters (e.g., `.`, `*`, `(`). Substring match is safer and sufficient.
- Fuzzy / ranked match: Out of scope per spec constraints and no-new-dependency constraint.

---

### 4. Product name data source

**Decision**: Stamp the product name into a `data-product-name` attribute on each card at template render time (server-side Go template).

**Rationale**: The name is already available in the Go template loop variable (`.Item.Name`). Stamping it as a `data-*` attribute keeps the JS entirely attribute-driven — the filter doesn't need to parse inner text or traverse child nodes. It also avoids any mismatch if the display name is later formatted differently (e.g., truncated).

**Alternatives considered**:
- Read `textContent` of `.hot-product-card-name`: Works but is fragile — any inner markup change (e.g., wrapping in a `<span>`) could silently break the filter.
- Fetch product names via a JSON endpoint: Requires an extra network round-trip and a backend endpoint; violates the frontend-only constraint.

---

### 5. No-results feedback element

**Decision**: Pre-render a hidden `<div id="product-no-results">` sibling to the cards; show/hide it via `style.display`.

**Rationale**: Avoids dynamic DOM creation (which requires `document.createElement` boilerplate and risks XSS if any string interpolation is involved). The element is always in the DOM; only its visibility changes.

**Alternatives considered**:
- Dynamically insert a `<p>` node: More JS code, harder to style consistently, unnecessary complexity.
- Show an alert/toast: Inconsistent with the page's existing UI patterns.

---

### 6. Accessibility

**Decision**: Use `aria-label="Search products"` on the input and `type="search"` to expose native browser clear behaviour.

**Rationale**: `type="search"` gives users the native × clear button at no JS cost, satisfying FR-006 (native clear) and FR-009 (keyboard operability). `aria-label` satisfies FR-002 without requiring a visible `<label>` element that would need grid layout adjustments.

**Alternatives considered**:
- Visible `<label>` + `for` attribute: Equally valid but requires layout adjustment; `aria-label` achieves the same accessibility result with less template change.

---

## Summary: No Unknowns Remaining

All six design decisions are resolved. The feature requires a single file change (`src/frontend/templates/home.html`) with no new dependencies, no backend changes, and no new files in the source tree.
