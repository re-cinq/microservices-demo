# Plan: Frontend Product Filter

## Technical Context

| Item | Detail |
|------|--------|
| Frontend language | Go (HTML templates via `html/template`) |
| Template file | `src/frontend/templates/home.html` |
| Stylesheet | `src/frontend/static/styles/styles.css` |
| Product data in template | `{{ range $.products }}` — each product has `.Item.Name`, `.Item.Id`, `.Item.Picture`, `.Price` |
| JS approach | Vanilla JS — no framework, no build step, no new dependencies |
| No-backend constraint | Filtering operates entirely on already-rendered DOM nodes |

---

## Design Decisions

- **Vanilla JS only.** The project has no JS build toolchain. A small inline `<script>` block (or a new static `.js` file) is the right fit. Inline script in the template keeps the change self-contained.
- **DOM-based filtering.** Each `.hot-product-card` div already contains a `.hot-product-card-name` element. The filter reads `textContent` from that element and toggles `display` on the parent card. No data attributes needed.
- **CSS for the "no results" message.** A hidden `<div id="no-results-msg">` is added to the template; JS shows/hides it based on match count.
- **No style file changes required** for core functionality — show/hide via `element.style.display`. A small optional CSS block can be added inline in the template for search box styling to match the existing Bootstrap-based UI.

---

## Architecture

```
src/frontend/templates/home.html   ← only file that changes
```

No Go code changes. No new files required (optionally a new `src/frontend/static/styles/product-filter.css` if styles grow large, but inline is fine for this scope).

---

## Implementation Steps

### Step 1 — Add the search input to `home.html`

Insert above the `<div class="row hot-products-row ...">`:

```html
<div class="row mb-3 px-xl-6">
  <div class="col-12">
    <input
      type="text"
      id="product-search"
      class="form-control"
      placeholder="Search products…"
      autocomplete="off"
    />
  </div>
</div>
```

### Step 2 — Add the "no results" message

Inside the `hot-products-row` div, after the `{{ range }}` block:

```html
<div class="col-12" id="no-results-msg" style="display:none;">
  <p class="text-muted">No products match your search.</p>
</div>
```

### Step 3 — Add the filter script

At the bottom of the template (before `{{ end }}`):

```html
<script>
  (function () {
    const input = document.getElementById('product-search');
    const cards = document.querySelectorAll('.hot-product-card');
    const noResults = document.getElementById('no-results-msg');

    input.addEventListener('input', function () {
      const query = input.value.trim().toLowerCase();
      let visible = 0;

      cards.forEach(function (card) {
        const name = card.querySelector('.hot-product-card-name').textContent.toLowerCase();
        const match = query === '' || name.includes(query);
        card.style.display = match ? '' : 'none';
        if (match) visible++;
      });

      noResults.style.display = visible === 0 ? 'block' : 'none';
    });
  })();
</script>
```

---

## File Change Summary

| File | Change |
|------|--------|
| `src/frontend/templates/home.html` | Add search input, no-results div, inline script |

No other files need to change.
