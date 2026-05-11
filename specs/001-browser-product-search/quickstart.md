# Quickstart — Verify Browser-Side Product Search

**Feature**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

This is the manual verification walk-through for the feature. Each step maps to a specific acceptance scenario or success criterion in `spec.md`. A reviewer ticks them off in order; if every step passes, the feature is done.

There is no automated test suite for this feature by design — the spec's Out-of-Scope section forbids touching the build pipeline, and the catalogue size (~10 products) makes manual verification fast and reliable.

---

## 0. Prerequisites

You need the `frontend` service running with the product list page reachable. Pick whichever path is set up for you:

- **Local Skaffold / Docker dev loop**: from the repo root, `skaffold dev` (or whatever the team's standard local-run command is). Open `http://localhost:8080/` (the home page).
- **Pushed to `attendee/daniel-tufvander`**: open whatever URL the cohort CI deploys the branch to. The page you want is the home page (product list).

Open the browser DevTools (F12 / Cmd-Option-I) and switch to:

- **Elements** tab — to verify the DOM contract.
- **Console** tab — for the pure-function assertions in §3.
- **Network** tab — for the "no extra requests per keystroke" check in §4.

---

## 1. Smoke check — page loads and the search box is there

1. Load the home page.
2. **Expect**: the search input is visible above the product grid. The product grid shows every product (currently ~10).
3. In DevTools → Elements, confirm:
   - `<input type="search" id="product-search" aria-label="Search products" …>` exists.
   - `<div id="hot-products-grid" …>` wraps the product cards.
   - Each product card has both a stable class (e.g. `hot-product-card`) **and** a `data-name="<lowercased name>"` attribute.
   - `<p id="product-search-no-results" … hidden>No products match your search.</p>` exists with `hidden` set.

If any of the above is missing, see [`contracts/ui-contract.md`](./contracts/ui-contract.md) §1.

---

## 2. Acceptance scenarios from spec.md

Run each scenario in order. After each one, **clear the search box** before moving to the next.

### User Story 1 — filter the visible list

- **§2.1 Scenario 1 (spec FR-003, FR-005)**
  Type `sun` into the search box.
  **Expect**: only products whose name contains `sun` (case-insensitive) remain visible. With the current catalogue, that is **Sunglasses**. Everything else is `hidden`.

- **§2.2 Scenario 2 (FR-009)**
  With `sun` filtered, delete every character.
  **Expect**: all products are visible again. The no-results element is `hidden`.

- **§2.3 Scenario 3 (FR-004)**
  Type `SUN` (uppercase).
  **Expect**: identical result to §2.1 — case is ignored.

- **§2.4 Scenario 4 (FR-007, SC-002)**
  With the box empty, type a single character (e.g. `t`).
  **Expect**: the list narrows immediately (no spinner, no reload, no flicker). Visible result with the current catalogue: products whose name contains `t` (Tank Top, Watch, etc. — verify visually).

- **§2.5 Scenario 5 (FR-008)**
  Type `xyzzy`.
  **Expect**: the product grid is visually empty. The "No products match your search." message is visible. The search box still contains `xyzzy`.

### User Story 2 — reset

- **§2.6 Scenario 1 (User Story 2)**
  Type something. Then either:
  - Press the native clear "x" inside the `<input type="search">`, **or**
  - Focus the box and press Escape.
  **Expect**: the box becomes empty, the no-results element is `hidden`, all products are visible.

### User Story 3 — filter resets on reload

- **§2.7 Scenario 1 (FR-010, SC-008)**
  Type `tank`, then press F5 / Cmd-R to reload.
  **Expect**: after reload, the search box is empty and all products are visible.

- **§2.8 Scenario 2 (User Story 3 Scenario 2)**
  Type `watch`. Click a product to navigate to its detail page. Click the browser Back button.
  **Expect**: the home page renders. The page is usable — no stuck filtered state with a hidden query. (Either the filter is reset, or the filter and query are both preserved — both are acceptable per the spec.)

---

## 3. Pure-function check (DevTools console)

Paste these into the Console:

```js
console.assert(window.__productSearch.matches("sunglasses", "sun")  === true,  "lowercase substring");
console.assert(window.__productSearch.matches("Sunglasses", "SUN")  === true,  "case-insensitive (note: data-name is already lowercased; the function also lowercases the query)");
console.assert(window.__productSearch.matches("watch", "tank")      === false, "non-match");
console.assert(window.__productSearch.matches("anything", "")       === true,  "empty query matches all");
console.assert(window.__productSearch.matches("anything", "   ")    === true,  "whitespace-only matches all");
console.assert(window.__productSearch.matches("tank top", "  TANK ")=== true,  "trimming + lowercasing");
console.assert(window.__productSearch.matches("watch", "(")         === false, "literal special characters do not throw");
console.log("product-search pure function: OK");
```

**Expect**: no `Assertion failed` lines in the console. Final message `product-search pure function: OK` is printed.

---

## 4. "Zero additional server requests per keystroke" (spec SC-006)

1. Open DevTools → Network tab. Click "Clear" (the ⊘ button).
2. Type `sunglasses` into the search box, one character at a time.
3. **Expect**: zero new entries appear in the Network panel during typing. (The static JS file may show up once on initial page load — that's fine. What must NOT happen is one network request per keystroke.)

---

## 5. Accessibility check (spec SC-007)

Either:

- Run **axe DevTools** (browser extension) on the home page, before and after the change. **Expect**: zero new violations.
- Or, minimally, verify in DevTools → Elements:
  - The `<input>` has either an associated `<label for="product-search">` or an `aria-label` attribute.
  - The no-results element has `role="status"` and `aria-live="polite"`.

---

## 6. Constraint sweep (spec Out-of-Scope)

Run from the repo root, on the feature branch, just before opening the PR:

```bash
# No new services
test ! -d src/product-search-service && echo "OK: no new service"

# Productcatalog untouched
git diff --name-only main...HEAD -- src/productcatalogservice/
# Expect: empty output

# No infra/manifest/helm/terraform changes
git diff --name-only main...HEAD -- kubernetes-manifests/ helm-chart/ kustomize/ terraform/ cloudbuild.yaml skaffold.yaml .github/
# Expect: empty output

# No protobuf changes
git diff --name-only main...HEAD -- protos/
# Expect: empty output

# Only frontend touched
git diff --name-only main...HEAD
# Expect: every line begins with src/frontend/  OR  specs/001-browser-product-search/
```

If any of the "Expect: empty" outputs is non-empty, the constraint sweep has failed — review the spec's Out-of-Scope section and reduce the change.

---

## Done

If every box above is ticked, the feature meets its spec and is ready to ship on `attendee/daniel-tufvander`.
