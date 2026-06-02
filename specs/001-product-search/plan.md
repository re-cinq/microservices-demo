# Implementation Plan: Browser-Side Product Search

**Branch**: `attendee/joel-johannesson` | **Date**: 2026-06-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-product-search/spec.md`

---

## Summary

Add a client-side product name search to the Online Boutique home page. A search input is placed above the product grid in `home.html`. A small vanilla-JavaScript file reads the input on every keystroke and shows/hides the existing product cards by matching the query against each card's `data-name` attribute. No backend service is changed, no new RPC calls are made, no new infrastructure is added — the filtering runs entirely in the browser against DOM elements already rendered by the Go template.

---

## Technical Context

**Language/Version**: Go (frontend service) — template and static asset change only; no Go source files are modified

**Primary Dependencies**: Bootstrap 4.1.1 (CDN, already present); vanilla JavaScript (no new external dependencies)

**Storage**: N/A — all filtering operates on DOM elements already present in the page at load time

**Testing**: Manual browser testing against the locally running frontend; existing Go unit tests are unaffected

**Target Platform**: Web browser (served by the Go HTTP frontend service)

**Project Type**: Web application — change is limited to one HTML template and one new static JS file

**Performance Goals**: Filter response < 300 ms per keystroke for catalogues up to 500 products (SC-001)

**Constraints**:
- No changes to any backend service or proto definition
- No new environment variables, Helm charts, Kubernetes manifests, or CI/CD steps
- No new network requests triggered by typing
- Existing product list layout (Bootstrap grid, `.hot-product-card` cards) must not change
- Results render inline — no separate results page

**Scale/Scope**: Up to 500 products; current catalogue has 9 items (`products.json`)

---

## Constitution Check

The project constitution is a default template (no project-specific principles have been ratified). Standard engineering principles apply:

- **Simplicity**: The change is the minimal modification to deliver the feature — one template edit, one new static file. No new abstractions or services.
- **No unnecessary complexity**: Vanilla JS is sufficient. No build pipeline step is needed.
- **Existing backend left intact**: The `SearchProducts` RPC already exists in `productcatalogservice` and performs server-side filtering; this feature deliberately does not use it, per spec constraints C-001 and C-004.

**Gate result**: PASS — no violations.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-product-search/
├── plan.md              ← this file
├── research.md          ← Phase 0 decisions
├── data-model.md        ← Phase 1 entities
├── contracts/
│   └── ui-contract.md  ← Phase 1 DOM contract
└── tasks.md             ← Phase 2 output (/speckit-tasks — not yet created)
```

### Source Code changes

```text
src/frontend/
├── templates/
│   └── home.html                  ← MODIFIED: search input + data-name on cards
└── static/
    └── scripts/
        └── product-search.js      ← NEW: client-side filter logic
```

No other files are touched. `handlers.go`, `rpc.go`, `main.go`, all proto files, and all other services are unchanged.

**Structure Decision**: Single-service template + static asset modification. The Bootstrap grid layout in `home.html` is preserved in full. JavaScript lives in a new `static/scripts/` directory that mirrors the existing `static/styles/` pattern.

---

## Implementation Approach

### Step 1 — Add `data-name` to product cards in `home.html`

The existing product loop renders each card as `<div class="col-md-4 hot-product-card">`. Add a `data-name` attribute so the filter script has a clean, layout-independent handle on each product's name:

```html
{{ range $.products }}
<div class="col-md-4 hot-product-card" data-name="{{ .Item.Name }}">
  ... existing card content unchanged ...
</div>
{{ end }}
```

### Step 2 — Add search input and no-results message to `home.html`

Insert above the existing product grid row. The Bootstrap `form-control` class ensures the input inherits the existing visual style without new CSS:

```html
<div class="row search-bar-row">
  <div class="col-12">
    <input
      id="product-search"
      type="search"
      class="form-control"
      placeholder="Search products…"
      aria-label="Search products"
      autocomplete="off"
    >
  </div>
</div>
<div id="no-results-message" style="display:none;" class="col-12 text-center py-4">
  No products match your search.
</div>
```

### Step 3 — Create `static/scripts/product-search.js`

Responsibilities:
- Listen for `input` events on `#product-search`
- For each `.hot-product-card[data-name]`, compare `data-name.toLowerCase()` against the trimmed, lower-cased query (substring match)
- Toggle `display` on each card element (visible / hidden)
- After each pass, if zero cards are visible and the query is non-empty, show `#no-results-message`; otherwise hide it
- On empty or whitespace-only query, restore all cards and hide the message

### Step 4 — Load the script in `home.html`

Add after the existing Bootstrap JS CDN tags (near the end of the included footer template or bottom of `home.html`):

```html
<script src="{{ $.baseUrl }}/static/scripts/product-search.js"></script>
```

The static file server in `main.go` already serves everything under `static/` — no routing change needed.

---

## Complexity Tracking

No Constitution Check violations — not applicable.
