# UI Contract: Top rated badge

**Scope**: a single badge rendered on the product card(s) tied at the maximum rating, inside the existing `.hot-product-card` in `src/frontend/templates/home.html`. Out of scope: any other surface (product page, cart, search results).

## Markup

A single `<span>` rendered as a direct child of `.hot-product-card`, gated on the new `.TopRated` flag from the view-model:

```html
<div class="col-md-4 hot-product-card">
  {{ if .TopRated }}
  <span class="hot-product-card-top-rated">Top rated</span>
  {{ end }}

  <a href="…">…</a>
  <div>
    <div class="hot-product-card-name">{{ .Item.Name }}</div>
    <div class="hot-product-card-price">…</div>
    {{ if gt .Rating 0.0 }}
    <div class="hot-product-card-rating" …>
      {{ renderStars .Rating }}
    </div>
    {{ end }}
  </div>
</div>
```

- The badge is the **first child** of `.hot-product-card` so its absolute positioning has a clean stacking context inherited from the parent.
- No new template function is needed; the badge is a literal `<span>` with a static text node.
- The badge MUST NOT be inside the `<a>` — clicking it does the same thing as clicking anywhere on the card image (i.e. nothing extra), per A4.

## Visual rules

| Aspect | Value |
|---|---|
| Position | top-right of `.hot-product-card`, `position: absolute; top: 8px; right: 8px;` |
| Background colour | `#f5a623` (matches filled-star colour from AIP-161) |
| Text colour | `#ffffff` |
| Text content | literal `Top rated` (English, title case) |
| Font size | `11px` |
| Font weight | `600` |
| Padding | `2px 8px` |
| Border radius | `999px` (pill shape) |
| Letter spacing | `0.5px` (subtle, mirrors the rating widget) |
| Text transform | none — the copy is already `Top rated`, not `top rated` |
| Z-index | sits above the image overlay; explicit `z-index: 2` |

## Render rules

- The badge MUST appear **only** when `.TopRated` is true on the rendered view-model.
- The badge MUST appear on **every** card whose `.TopRated` is true — no client-side dedupe, no "first one wins" logic.
- The badge MUST NOT alter the position, size, or visibility of any existing element on the card (name, price, image, star widget). Satisfies FR-005. Absolute positioning is how we achieve this.

## Accessibility

- The badge is plain text inside a `<span>` — screen readers will read "Top rated" inline as part of the card. No additional `aria-label` is necessary; no `aria-hidden`. See research D6.
- Colour contrast: white on `#f5a623` measures ~3.1:1 — borderline for small text. The badge text is `11px / 600` weight, which is at the boundary of WCAG AA-Large. If the cohort design later requires strict AA-Normal, swap text to `#222` for ~7:1 contrast.

## Interaction

None. The badge does not respond to hover, focus, or tap.

## What MUST NOT change

- The existing star widget on the same card (its position, copy, ARIA label, CSS classes). Satisfies FR-005.
- The card image, name, price, or product link.
- Any other template, handler, or static asset not listed in `plan.md`'s Source Code section.

## Verification (covered by quickstart.md)

- Open the cohort home page. With the current seed, the Watch card (rating 5.0) carries the badge; no other card does.
- Temporarily edit a second product's rating to `5.0` in `products.json`, redeploy, confirm both cards carry the badge (AS-2). Revert.
- Set all `rating` values to `0.0` temporarily, confirm no card carries the badge (AS-3). Revert.
