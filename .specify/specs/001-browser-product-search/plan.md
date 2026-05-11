# Implementation Plan: Browser Product Search

**Branch**: `001-browser-product-search` | **Date**: 2026-05-11 | **Spec**: `specs/001-browser-product-search/spec.md`
**Input**: Feature specification from `/specs/001-browser-product-search/spec.md`

## Summary

Add a search box to the Online Boutique homepage that sends a GET request with a `?q=` parameter. The frontend reads the query parameter, calls the **existing** `SearchProducts` gRPC RPC on `productcatalogservice` (no backend changes), converts prices to the selected currency, and renders matching products using the existing product card layout. Empty/whitespace queries fall back to `ListProducts` (show all).

## Technical Context

**Language/Version**: Go (matching existing frontend and productcatalogservice)
**Primary Dependencies**: gorilla/mux (routing), html/template (rendering), gRPC (service communication) — all already in use
**Storage**: N/A — existing in-memory product catalogue from `products.json`
**Testing**: Manual testing against training URL
**Target Platform**: Linux container (Kubernetes), accessed via browser
**Project Type**: Multi-service web application (microservices)
**Performance Goals**: Results render in under 500ms for catalogue of up to 100 items
**Constraints**: No new services, datastores, infra config, or environment variables

## Constitution Check

No violations. This plan:
- Uses only existing services (frontend, productcatalogservice)
- Uses only existing languages (Go)
- Uses only existing RPC patterns (gRPC + protobuf)
- Adds no new dependencies, infrastructure, or datastores
- Changes exactly 3 files in one service (frontend)

## Changes Required

### Service: `src/frontend/` (Go) — 3 files touched

#### 1. `src/frontend/rpc.go` — Add `searchProducts` helper

Add a new function that wraps the existing `SearchProducts` RPC, following the exact same pattern as `getProducts` and `getProduct`:

```go
func (fe *frontendServer) searchProducts(ctx context.Context, query string) ([]*pb.Product, error) {
    resp, err := pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).
        SearchProducts(ctx, &pb.SearchProductsRequest{Query: query})
    return resp.GetProducts(), err
}
```

**Why**: The proto and generated client already include `SearchProducts` and `SearchProductsRequest`. The frontend just never calls it. This adds the thin wrapper matching the existing pattern.

#### 2. `src/frontend/handlers.go` — Modify `homeHandler`

In the existing `homeHandler` function (~line 59–120), add query parameter reading and conditional search:

- Read `q` from `r.URL.Query().Get("q")`
- Trim whitespace with `strings.TrimSpace()`
- If `q` is non-empty → call `fe.searchProducts(ctx, q)` instead of `fe.getProducts(ctx)`
- If `q` is empty → call `fe.getProducts(ctx)` as before (existing behaviour unchanged)
- Pass `q` to the template data as `"search_query"` so the template can:
  - Pre-fill the search box with the current query
  - Show "No products found for 'query'" when results are empty

**Existing currency conversion is untouched** — the `for i, p := range products` loop that calls `fe.convertCurrency()` runs on whatever product list was returned (all or filtered). This ensures search results respect the selected currency with zero additional code.

#### 3. `src/frontend/templates/home.html` — Add search form + empty state

Above the `<div class="row hot-products-row">` section (~line 40), add:

- An HTML `<form>` with `method="GET"` and `action="{{ $.baseUrl }}/"` containing:
  - A text input `name="q"` with `value="{{ $.search_query }}"` (preserves query on page reload)
  - Submission via Enter key (native HTML form behaviour)
- Inside the product loop area, add a conditional:
  - If `products` is empty AND `search_query` is non-empty → show "No products found for '{{ $.search_query }}'"
  - Otherwise → render product cards as before (existing `{{ range $.products }}` block unchanged)

### Service: `src/productcatalogservice/` — NO changes

The `SearchProducts` RPC already exists in `product_catalog.go` (line ~60). It does case-insensitive substring matching on both product name and description. Use it as-is.

### Protobuf: `protos/demo.proto` — NO changes

`SearchProductsRequest` (with `query` field) and `SearchProductsResponse` (with `results` field) are already defined. The generated Go client code in `src/frontend/genproto/` already includes the client stubs.

## Project Structure

### Files changed

```text
src/frontend/
├── rpc.go           # ADD searchProducts() wrapper (~5 lines)
├── handlers.go      # MODIFY homeHandler() to read ?q= and branch logic (~10 lines changed)
└── templates/
    └── home.html    # ADD search form + empty state conditional (~15 lines added)
```

### Files NOT changed

```text
src/productcatalogservice/    # No changes — SearchProducts RPC already exists
protos/demo.proto             # No changes — messages already defined
src/frontend/genproto/        # No changes — already generated
```

## Risk Assessment

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Search form breaks existing homepage layout | Low | Form is added above product grid, CSS unchanged. Test visually. |
| Currency conversion skipped for search results | Low | Same conversion loop runs regardless of product source. Verify with EUR selection. |
| XSS via search query echoed in template | Medium | Go's `html/template` auto-escapes by default. Verify with `<script>alert(1)</script>` query. |
| Empty string passed to SearchProducts RPC | Low | Handler trims whitespace and falls back to ListProducts. RPC is never called with empty query. |

## Complexity Tracking

No constitution violations. No complexity justification needed.
