# Data Model: Category Jump Links

**Feature**: AIP-187
**Date**: 2026-06-16

## Existing entities (unchanged)

### Product (from productcatalogservice proto)

| Field        | Type       | Notes                                      |
|--------------|------------|--------------------------------------------|
| Id           | string     | Unique product identifier                  |
| Name         | string     | Display name                               |
| Picture      | string     | Relative URL to product image              |
| Categories   | []string   | One or more category slugs, e.g. `["accessories"]` |
| PriceUsd     | Money      | Price in USD                               |

### productView (frontend struct, unchanged)

| Field | Type        | Notes                           |
|-------|-------------|---------------------------------|
| Item  | *pb.Product | Full product proto               |
| Price | *pb.Money   | Converted price in display currency |

## New view model (frontend only, not persisted)

### CategoryGroup

Constructed in `homeHandler` before template rendering. Lives only in request scope.

| Field    | Type           | Notes                                                       |
|----------|----------------|-------------------------------------------------------------|
| Name     | string         | Display name derived from the category slug (e.g. `"accessories"`) |
| Slug     | string         | URL-safe anchor id (e.g. `"accessories"`)                   |
| Products | []productView  | Ordered list of products in this category                   |

## Ordering rules

- **Category order**: determined by first-appearance of `Categories[0]` when iterating products in their existing JSON/gRPC order.
- **Product order within category**: same relative order as the current flat list.
- **Jump link order**: mirrors category order (same slice).

## Validation rules

- A `CategoryGroup` is only created when at least one product maps to it.
- A product with an empty `Categories` slice is placed into a synthetic group named `"Other"` (slug: `"other"`); no jump link is emitted for that group if it would be the only group.
- Jump links reference only groups that exist in the rendered page (no dangling anchors).
