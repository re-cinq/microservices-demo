# UI Contract: Product List Page — Category Jump Links

**Feature**: AIP-187
**Date**: 2026-06-16

## Template data contract

The `home.html` template receives the following top-level data shape (additions to existing keys shown with ★):

```
TemplateData {
  ...existing keys unchanged...
  ★ categories []CategoryGroup   // ordered list of category sections
}

CategoryGroup {
  Name     string         // human-readable label, e.g. "Accessories"
  Slug     string         // anchor-safe id, e.g. "accessories"
  Products []productView  // products in this category, in original order
}
```

The existing `products []productView` key **may be removed** from the template data once the template is migrated to `categories`; the handler still fetches the same product list.

## Rendered HTML contract

### Jump link bar

Rendered once at the top of the product section. One `<a>` per category:

```html
<nav aria-label="Product categories">
  <a href="#{{.Slug}}">{{.Name}}</a>
  ...
</nav>
```

- Each `href` value is `#<slug>` where `<slug>` matches exactly the `id` on the corresponding section heading.
- Order matches `categories` slice order.

### Category section heading

One per `CategoryGroup`, rendered before that group's products:

```html
<h3 id="{{.Slug}}">{{.Name}}</h3>
```

- The `id` attribute value equals `.Slug` — must match the jump link `href` fragment exactly.

### Product cards within a section

Identical markup to the current `hot-product-card` cards; only the surrounding loop structure changes.

## Behavioural contract

| Trigger | Expected behaviour |
|---|---|
| Page load | Jump link bar and all category headings are visible without interaction |
| Click jump link for category X | Viewport scrolls to the `<h3 id="X">` heading |
| No JavaScript | Page still renders correctly; smooth scroll may not apply but jump navigation still works |

## Accessibility contract

- The `<nav>` wrapping jump links carries `aria-label="Product categories"`.
- Category headings use `<h3>` (consistent with current "Hot Products" heading level).
