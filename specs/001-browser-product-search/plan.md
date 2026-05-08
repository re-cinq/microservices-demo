# Implementation Plan: Browser-Side Product Search

**Branch**: `001-browser-product-search` | **Date**: 2026-05-08 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `specs/001-browser-product-search/spec.md`

---

## Summary

Add a real-time, client-side product name filter to the Online Boutique home page. A text input above the product grid listens to every keystroke and shows or hides `.hot-product-card` elements based on a case-insensitive match against the product name text already rendered in the DOM. No backend changes, no new dependencies, no new services. Two files change: `home.html` (search box + inline JS) and `styles.css` (search box styling).

---

## Technical Context

**Language/Version**: Go 1.25 (frontend service); vanilla JavaScript (browser, no transpilation)  
**Primary Dependencies**: Go `html/template`, Bootstrap 4 (CDN, already present), no new dependencies  
**Storage**: N/A — filtering operates on in-memory DOM state only  
**Testing**: Manual browser smoke test against the running frontend container  
**Target Platform**: Any modern browser (Chrome, Firefox, Safari, Edge); existing responsive layout  
**Project Type**: Web service (frontend is a Go HTTP server serving server-side-rendered HTML)  
**Performance Goals**: Grid updates within 100 ms per keystroke for up to 500 products (actual: 30)  
**Constraints**: No backend change; no new infra; no new services; no build pipeline changes  
**Scale/Scope**: Single page (home), 30 products in catalogue

---

## Constitution Check

The constitution template has not been filled in for this project. No principle violations are detected:

- The change is purely additive (two files, no deletions).
- No new service, dependency, or infrastructure is introduced.
- The implementation is the simplest possible: inline JS + CSS addition in existing files.

*Post-design re-check*: research.md and data-model.md confirm the above. No violations.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-browser-product-search/
├── plan.md              (this file)
├── research.md          (Phase 0 output)
├── data-model.md        (Phase 1 output)
└── tasks.md             (Phase 2 output - created by /speckit-tasks)
```

### Source Code Changes

```text
src/frontend/
├── templates/
│   └── home.html        ADD: search input, no-results message, inline <script>
└── static/
    └── styles/
        └── styles.css   ADD: .product-search-box CSS rule
```

**No other files are modified.**

---

## Implementation Design

### 1. Search box placement — `home.html`

Insert immediately before the `{{ range $.products }}` loop, inside `.hot-products-row`:

```html
<div class="col-12 mb-3">
  <input
    id="product-search-input"
    type="text"
    class="product-search-box"
    placeholder="Search products..."
    aria-label="Search products"
    autocomplete="off"
  >
  <p id="no-results-message" class="text-muted mt-2" style="display:none">
    No products match your search.
  </p>
</div>
```

### 2. Inline filter script — `home.html`

Insert at the bottom of the `{{ define "home" }}` block, before `{{ end }}`:

```html
<script>
  (function () {
    var input = document.getElementById('product-search-input');
    var noResults = document.getElementById('no-results-message');

    input.addEventListener('input', function () {
      var query = input.value.trim().toLowerCase();
      var cards = document.querySelectorAll('.hot-product-card');
      var visible = 0;

      cards.forEach(function (card) {
        var name = card.querySelector('.hot-product-card-name');
        var matches = !query || (name && name.textContent.toLowerCase().indexOf(query) !== -1);
        card.style.display = matches ? '' : 'none';
        if (matches) visible++;
      });

      noResults.style.display = (query && visible === 0) ? '' : 'none';
    });
  })();
</script>
```

Key properties:
- IIFE wrapper avoids polluting global scope.
- Empty/whitespace-only query (after `trim()`) shows all cards and hides the no-results message.
- Uses `indexOf` for broad browser compatibility.
- Resets `display` to `''` (not `'block'`) so Bootstrap's responsive column classes remain in control.

### 3. CSS addition — `styles.css`

Append to `src/frontend/static/styles/styles.css`:

```css
.product-search-box {
  width: 100%;
  max-width: 400px;
  padding: 0.5rem 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
}
```

---

## Acceptance Test Mapping

| Acceptance Scenario | How to verify |
|---------------------|--------------|
| Type partial name → only matching cards visible | Type "Sunglasses" → only sunglasses card shown |
| No match → zero cards + "no results" message | Type "zzzzzz" → grid empty, message appears |
| Clear input → all products restored, no reload | Delete all characters → all cards reappear |
| Case-insensitive: "sunglasses" matches "Sunglasses" | Type lowercase name → card still visible |
| Whitespace-only input → all products shown | Type spaces only → no filtering applied |
| Keyboard accessible | Tab to input, type query → grid filters without mouse |

---

## Out of Scope

- Filtering by description, category, or price
- Persistent search state across navigation
- URL query parameter for shareable searches
- Any backend or gRPC change
- Any new service, Helm chart, manifest, or environment variable
