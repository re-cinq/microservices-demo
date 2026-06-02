# Implementation Plan: Product Search

**Feature directory**: `specs/001-product-search`
**Spec**: [spec.md](spec.md)
**Branch**: `001-product-search` (feature branch off `attendee/jost-werdenhoff`; merges back via PR — C-7)
**Created**: 2026-06-02

## Technical Context

### Tech stack (existing — do not add to)
- **Service touched**: `src/frontend` — a Go HTTP server (`net/http`) rendering server-side HTML via Go `html/template` (`templates/*.html`), styled with Bootstrap 4.1.1 (loaded from CDN in `templates/footer.html`).
- **Product data path**: `homeHandler` (`src/frontend/handlers.go:59`) calls `fe.getProducts()` (gRPC → `productcatalogservice`, whose catalogue is loaded from `productcatalogservice/products.json`). Each product is wrapped in a `productView{Item *pb.Product, Price *pb.Money}` and passed to the `home` template as `products`.
- **Grid markup**: `templates/home.html:46-57` — `{{ range $.products }}` emits one `<div class="col-md-4 hot-product-card">` per product. The product name is rendered at `home.html:53` as `<div class="hot-product-card-name">{{ .Item.Name }}</div>`. The only keyboard-focusable element in a card is the `<a href=...>` wrapping the image (`home.html:48`).
- **Inline scripts are already used** (`templates/assistant.html:50`) and **no Content-Security-Policy** is set (checked `handlers.go`, `middleware.go`, `main.go`) — so a small inline `<script>` on the home page is safe and consistent with existing patterns.

### Constraints honored (from spec)
- Frontend-only, in-browser filter over products **already rendered** — no new gRPC call, no server round-trip (C-9, C-2, C-8, FR-008).
- No new service, datastore, infra config, env var, or pipeline change (C-1, C-5, C-6).
- The only files changed are inside `src/frontend/` (Go service language unchanged — no Go logic change actually required; this is a template + client-side change). No `.proto` change (C-4 N/A since backend is untouched).
- Work on the `001-product-search` feature branch, merged back to `attendee/jost-werdenhoff` via PR (C-7).

## Design Decisions

### D-1: Pure client-side filter, no backend change
The home page already renders the entire catalogue into the DOM. Filtering = showing/hiding existing `.hot-product-card` nodes based on their name. This satisfies "filter the products already loaded, in the browser" with zero backend or gRPC work. **No change to `handlers.go` or any `.proto` is needed.**

### D-2: Match against a `data-product-name` attribute, not scraped text
Add `data-product-name="{{ .Item.Name }}"` to each product card in the template. The script lowercases this attribute and the query and does a substring test. Using a data attribute (rather than reading `.hot-product-card-name` text) keeps matching robust against any future markup/whitespace changes inside the card. (Go's `html/template` auto-escapes the attribute value — safe against injection.)

### D-3: Hide non-matching cards with a CSS class that sets `display:none`
Setting `display:none` on a non-matching card:
- removes it visually (FR-002),
- **removes its `<a>` from the tab order and from the accessibility tree automatically** — satisfying FR-013 (no keyboard reach to hidden cards) with no extra `tabindex` juggling.
Restoring = removing the class (FR-005, FR-007). No nodes are destroyed, so clearing is instant and order is preserved (SC-003).

### D-4: Accessible search input with a visible label
Add a `<label for="product-search">` (visible, e.g. "Search products") bound to the input (FR-010). The input is `type="search"` so browsers provide a native clear affordance (supports User Story 2 / FR-007).

### D-5: Accessible no-results message via `aria-live`
Add a hidden `<p id="search-no-results" role="status" aria-live="polite">No products match your search.</p>`. The script toggles its visibility. `role="status"`/`aria-live="polite"` means assistive tech announces it when it appears (FR-006, FR-011, SC-004, SC-007). It is real text content, visible to sighted users too.

### D-6: Filter on the `input` event; never move focus
The script listens to the search input's `input` event and only toggles classes on cards + the no-results message. It **never calls `.focus()` elsewhere**, so focus stays in the search box while typing (FR-012). Empty/whitespace-only query → trim → show all (FR-005, edge case).

### D-7: Graceful degradation
All filtering lives in the inline `<script>`. If JS is unavailable, the page renders the full grid exactly as today — nothing breaks (spec edge case / graceful degradation).

## File Structure (changes)

All changes confined to `src/frontend/` (templates + optional CSS). **No Go code change required.**

| File | Change |
|------|--------|
| `src/frontend/templates/home.html` | (1) Add a search `<label>` + `<input id="product-search" type="search">` and the `aria-live` no-results `<p>` just above the product grid (around `home.html:40-45`). (2) Add `data-product-name="{{ .Item.Name }}"` to the `hot-product-card` div (`home.html:47`). (3) Add an inline `<script>` near the end of the home content implementing D-2/D-3/D-5/D-6. |
| `src/frontend/static/styles/styles.css` *(or the existing home stylesheet)* | Add a `.hot-product-card--hidden { display: none; }` rule and light styling for the search box / no-results text, matching existing Bootstrap look. *(Optional — could also reuse Bootstrap's `d-none`; decide in tasks.)* |

### Sketch of the home.html additions

```html
<!-- above the grid (near home.html:42) -->
<div class="col-12">
  <label for="product-search">Search products</label>
  <input type="search" id="product-search" autocomplete="off"
         placeholder="Search products by name">
  <p id="search-no-results" role="status" aria-live="polite" hidden>
    No products match your search.
  </p>
</div>

<!-- card gets the data attribute (home.html:47) -->
<div class="col-md-4 hot-product-card" data-product-name="{{ .Item.Name }}">

<!-- inline script near end of home content -->
<script>
  (function () {
    var input = document.getElementById('product-search');
    if (!input) return;                         // graceful degradation
    var cards = Array.prototype.slice.call(
      document.querySelectorAll('.hot-product-card'));
    var noResults = document.getElementById('search-no-results');
    input.addEventListener('input', function () {
      var q = input.value.trim().toLowerCase();
      var visible = 0;
      cards.forEach(function (card) {
        var name = (card.getAttribute('data-product-name') || '').toLowerCase();
        var match = q === '' || name.indexOf(q) !== -1;
        card.classList.toggle('hot-product-card--hidden', !match);  // display:none
        if (match) visible++;
      });
      if (noResults) noResults.hidden = !(q !== '' && visible === 0);
    });
  })();
</script>
```

## Requirements → Design traceability

| Requirement | Covered by |
|-------------|------------|
| FR-001 search input visible | D-4 (label + input above grid) |
| FR-002 live filter, no reload | D-1, D-6 (`input` event toggles classes) |
| FR-003 case-insensitive | D-2 (`toLowerCase`) |
| FR-004 substring | D-2 (`indexOf`) |
| FR-005 empty → show all | D-6 (trim → `q === ''`) |
| FR-006 / FR-011 no-results message | D-5 (`aria-live` `<p>`) |
| FR-007 clear restores all | D-3 (remove class), `type="search"` clear (D-4) |
| FR-008 / C-9 no server request | D-1 (DOM-only) |
| FR-009 no other behaviour changed | additive markup only; no Go change |
| FR-010 accessible name | D-4 (`<label for>`) |
| FR-012 focus not moved/trapped | D-6 (no `.focus()` calls) |
| FR-013 hidden cards unreachable by keyboard | D-3 (`display:none` drops from tab order) |

## Out of scope / explicitly not doing
- No changes to `productcatalogservice`, `.proto` files, or any gRPC call (C-1, C-2, C-3, C-4, C-9).
- No new datastore, search engine, infra config, env var, Helm/manifest, or CI change (C-2, C-5, C-6).
- No search over description, category, or price — name only (C-8).
- Work on the `001-product-search` feature branch only; no branching from elsewhere; merge back to `attendee/jost-werdenhoff` via PR (C-7).

## Verification approach (feeds /speckit.tasks)
- Manual/browser: load home page, type "watch" / "WATCH" / "sun" / "zzzzz" / clear → assert scenarios 1–5.
- Keyboard: Tab through after filtering → focus skips hidden cards (FR-013); typing keeps focus in box (FR-012).
- Screen reader / accessibility tree: input has an accessible name (FR-010); no-results announced (FR-011).
- Network panel: typing triggers zero new requests (FR-008/SC-006).
- Regression: unfiltered grid identical to current (order, cart, currency) (FR-009/SC-005).
