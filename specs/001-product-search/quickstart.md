# Quickstart: Implementing Product Search

**Phase 1 output for**: `specs/001-product-search/plan.md`  
**Date**: 2026-05-11

A concise guide for any developer picking up this task cold.

---

## What to change and where

### 1. `src/productcatalogservice/product_catalog.go` — extend the filter

Find `SearchProducts` (~line 60). Change the single-field filter to an OR:

```go
// Before
if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) {

// After
query := strings.ToLower(req.Query)
if strings.Contains(strings.ToLower(product.Name), query) ||
    strings.Contains(strings.ToLower(product.Description), query) {
```

---

### 2. `src/frontend/rpc.go` — add gRPC wrapper

Add alongside `getProducts` / `getProduct`:

```go
func (fe *frontendServer) searchProducts(ctx context.Context, query string) ([]*pb.Product, error) {
    resp, err := pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).
        SearchProducts(ctx, &pb.SearchProductsRequest{Query: query})
    return resp.GetResults(), err
}
```

---

### 3. `src/frontend/handlers.go` — extend homeHandler + add template function

**a) Add `highlightTerm` to the template FuncMap** (near line 49):

```go
import "html"   // add to imports

// In the Funcs(...) map:
"highlightTerm": highlightTerm,
```

**b) Add the `highlightTerm` function** (new top-level function in the file):

```go
func highlightTerm(term, text string) template.HTML {
    escaped := html.EscapeString(text)
    if term == "" {
        return template.HTML(escaped)
    }
    lowerText := strings.ToLower(escaped)
    lowerTerm := strings.ToLower(html.EscapeString(term))
    var out strings.Builder
    for i := 0; i < len(escaped); {
        idx := strings.Index(lowerText[i:], lowerTerm)
        if idx == -1 {
            out.WriteString(escaped[i:])
            break
        }
        out.WriteString(escaped[i : i+idx])
        out.WriteString("<mark>")
        out.WriteString(escaped[i+idx : i+idx+len(lowerTerm)])
        out.WriteString("</mark>")
        i += idx + len(lowerTerm)
    }
    return template.HTML(out.String())
}
```

**c) Extend `homeHandler`** to read `?q=` and branch:

At the top of `homeHandler`, after the existing `currencies` and `cart` calls,
replace the `fe.getProducts(...)` block with:

```go
q := r.URL.Query().Get("q")
trimmed := strings.TrimSpace(q)

var (
    products     []*pb.Product
    searchError  string
    isSearch     bool
)

switch {
case q == "":
    // No search — show full catalogue
    products, err = fe.getProducts(r.Context())
case len([]rune(trimmed)) < 2:
    // Too short — block, show error, show full catalogue
    searchError = "Please enter at least 2 characters to search."
    products, err = fe.getProducts(r.Context())
default:
    // Valid search
    isSearch = true
    products, err = fe.searchProducts(r.Context(), trimmed)
}
if err != nil {
    renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve products"), http.StatusInternalServerError)
    return
}
```

Then pass the new keys to `ExecuteTemplate`:

```go
"search_query": q,
"search_error": searchError,
"is_search":    isSearch,
```

---

### 4. `src/frontend/templates/home.html` — add search form + results UI

**a) Search form** — add above the `hot-products-row` div:

```html
<div class="row px-xl-6 mb-3">
  <div class="col-12">
    <form method="GET" action="{{ $.baseUrl }}/">
      <div class="input-group">
        <input type="search" name="q" class="form-control"
               placeholder="Search products…"
               value="{{ $.search_query }}"
               minlength="2"
               aria-label="Search products">
        <div class="input-group-append">
          <button class="btn btn-primary" type="submit">Search</button>
          {{ if $.search_query }}
          <a href="{{ $.baseUrl }}/" class="btn btn-outline-secondary">Clear</a>
          {{ end }}
        </div>
      </div>
      {{ if $.search_error }}
      <div class="text-danger mt-1 small">{{ $.search_error }}</div>
      {{ end }}
    </form>
  </div>
</div>
```

**b) Section heading** — make the heading context-aware:

```html
{{ if $.is_search }}
<h3>Results for &ldquo;{{ $.search_query }}&rdquo;</h3>
{{ else }}
<h3>Hot Products</h3>
{{ end }}
```

**c) Empty state** — add inside the products row, after the `{{ range }}` block:

```html
{{ if and $.is_search (not $.products) }}
<div class="col-12 text-center py-5">
  <p class="text-muted">No products found for &ldquo;{{ $.search_query }}&rdquo;.</p>
  <a href="{{ $.baseUrl }}/">View all products</a>
</div>
{{ end }}
```

Note: Go templates don't have a built-in `not` — register it in the FuncMap,
or use the equivalent: `{{ if and $.is_search (eq (len $.products) 0) }}`.

**d) Highlighting** — in the product card name, swap plain text for highlighted output:

```html
{{/* was: {{ .Item.Name }} */}}
<div class="hot-product-card-name">
  {{ if $.search_query }}{{ highlightTerm $.search_query .Item.Name }}{{ else }}{{ .Item.Name }}{{ end }}
</div>
```

For description highlighting (shown on the card if desired):

```html
{{ if and $.is_search $.search_query }}
<div class="hot-product-card-desc small text-muted">
  {{ highlightTerm $.search_query .Item.Description }}
</div>
{{ end }}
```

---

## Testing

### Unit test — description filter fix

In `src/productcatalogservice/product_catalog_test.go`, add a test case that
passes a query matching only a product's `Description` (not its `Name`) and
asserts the product is returned.

### Manual verification checklist

| Scenario | Expected |
|---|---|
| `/?q=sunglasses` | Only "Sunglasses" card shown; "sunglasses" highlighted in name/description |
| `/?q=SUNGLASSES` | Same result (case-insensitive) |
| `/?q=a` | Full catalogue shown; validation message visible |
| `/?q=  ` | Full catalogue shown; validation message visible |
| `/?q=xyzzy123` | Empty state with "xyzzy123" in message |
| `/` (no q) | Full catalogue, no search UI error |
| Clear button | Returns to `/`, full catalogue |
| Term appears in description only | Product appears in results |
| Term contains `<script>alert(1)</script>` | Rendered as escaped literal text, no JS executes |
