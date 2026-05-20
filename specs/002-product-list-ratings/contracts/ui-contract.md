# UI Contract: Product-list rating widget

**Scope**: the rating widget rendered on every product card in `src/frontend/templates/home.html` (the "Hot Products" grid that is the product list). Consumers: shoppers using the live cohort app. Out of scope: the product page (sibling story AIP-162) and any non-list surface.

## Render rules

A single widget rendered inside the existing `.hot-product-card` block, *after* the product name and price, *before* any future card additions. The card layout above the widget is unchanged (FR-005, spec).

### Position on the card

```html
<div class="col-md-4 hot-product-card">
  <a href="…">…</a>
  <div>
    <div class="hot-product-card-name">{{ .Item.Name }}</div>
    <div class="hot-product-card-price">…</div>

    {{/* NEW: render only if .Rating > 0 */}}
    {{ if gt .Rating 0.0 }}
    <div class="hot-product-card-rating"
         aria-label="Rated {{ printf "%.1f" .Rating }} out of 5 stars">
      {{ renderStars .Rating }}
    </div>
    {{ end }}
  </div>
</div>
```

- `renderStars` is a new Go template function on the frontend that returns a safe HTML string of filled / half / empty star spans according to the value. The function MUST be deterministic — same input, same output — to satisfy SC-003 (stability across page loads).
- `.Rating` is sourced from the new `Rating float32` field on the existing `productView` struct in `handlers.go`.

### Visual rules

| Rating value | Visible glyphs |
|---|---|
| `0.0` | (widget not rendered — see "no-rating fallback" below) |
| `0.5` | half · empty · empty · empty · empty |
| `1.0` | filled · empty · empty · empty · empty |
| `1.5` | filled · half · empty · empty · empty |
| `2.0` | filled · filled · empty · empty · empty |
| `2.5` | filled · filled · half · empty · empty |
| `3.0` | filled · filled · filled · empty · empty |
| `3.5` | filled · filled · filled · half · empty |
| `4.0` | filled · filled · filled · filled · empty |
| `4.5` | filled · filled · filled · filled · half |
| `5.0` | filled · filled · filled · filled · filled |

Glyph technique: Unicode (`★` filled, `☆` empty, half-star via CSS overlay or the `⯨` glyph), styled by a new rule in an existing stylesheet under `src/frontend/static/styles/`. No new font, no new image assets, no new build step (research D8).

### No-rating fallback

If `Rating == 0.0`, the `<div class="hot-product-card-rating">` MUST NOT be emitted at all. The card collapses gracefully — there is no empty placeholder. (research D7)

### Accessibility

- The wrapper element carries `aria-label="Rated <X.X> out of 5 stars"` where `<X.X>` is the rating with one decimal place. This is the canonical signal for screen readers (research D6).
- The star glyphs themselves are presentational; they are *not* announced individually.

### Interaction

None. The widget is read-only and not clickable. It does not respond to hover, focus, or tap (research D3, spec FR-004).

## What MUST NOT change

- The order of existing elements above the widget (name, price).
- The product link `href` on the card image and name.
- The grid columns (`col-md-4`), card width, or card height. The widget adds vertical extent inside the card; the card grid wraps as it does today.
- Any other template, handler, or static asset not listed in `plan.md`'s Source Code section.

## Verification (covered by quickstart.md)

- Open the cohort home page. Every product card shows a star widget (subject to no-rating fallback).
- Reload the page twice; the ratings shown for the same products are identical (SC-003).
- A screen reader reads each widget as "Rated X.X out of 5 stars" (manual check; not automated in this story).
