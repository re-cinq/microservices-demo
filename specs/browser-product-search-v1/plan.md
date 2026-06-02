# Implementation Plan: Browser Product Search v1

## Tech Stack

No new dependencies. Everything below uses what is already in the repo:

| Concern | Technology |
|---|---|
| Backend language | Go (existing frontend service) |
| HTTP routing | `github.com/gorilla/mux` (already used) |
| gRPC client | `productcatalogservice.SearchProducts` — already exists, no proto changes needed |
| Currency conversion | `currencyservice.Convert` — already used in `homeHandler` |
| Templating | `html/template` (already used); safe HTML injection via `template.HTML` for `<mark>` highlighting |
| Validation message | Inline HTML rendered server-side; JS progressive enhancement for immediate client feedback |
| Styling | Existing Bootstrap 4 + `static/styles/styles.css` — no new CSS framework |

---

## Design Decisions

### 1. Search logic stays in the frontend service

`productcatalogservice.SearchProducts` already performs case-insensitive substring matching on name and description via `strings.Contains`. It returns the correct candidate set. Ranking (R-08) and highlighting (R-03) are presentation concerns — they require knowledge of the query string at render time and operate on the already-filtered result list. Putting them in the frontend avoids a proto change and keeps `productcatalogservice` focused on data access.

**Trade-off:** If the catalogue grows very large, ranking in the frontend service means deserialising all matching products over gRPC before sorting. Acceptable for this catalogue size; revisit if catalogue exceeds ~10 k products.

### 2. No proto changes to productcatalogservice

`SearchProductsRequest` already carries a `Query` string; `SearchProductsResponse` already returns `[]*Product`. The existing gRPC call is used as-is via a new `searchProducts` wrapper in `rpc.go`, matching the pattern of `getProducts` / `getProduct`.

### 3. Short query validation: server-side primary, JS enhancement

The spec requires the validation message to appear inline without navigating away (R-05). The server handles `GET /search?q=<1 char>` by returning `200` with the search page pre-populated with the validation message and the typed query — no redirect, no gRPC call. A small `<script>` block in the search template intercepts form submission client-side for immediate feedback, but the server path is the source of truth.

### 4. Highlighting via a custom template function

`html/template` auto-escapes all output. To inject `<mark>` tags safely, a `highlightQuery(text, query string) template.HTML` function is registered in the template `FuncMap`. It escapes `text`, locates all case-insensitive occurrences of `query` (escaped), wraps each in `<mark>`, and returns `template.HTML`. This is the same pattern already used by `renderMoney` and `renderCurrencyLogo` in `handlers.go`.

### 5. Currency change on results page works automatically

`setCurrencyHandler` redirects to the `Referer` header. Since the search results URL is `/search?q=<query>`, changing currency while on the results page redirects back to the same search URL with the new currency cookie set — no additional code needed (R-06).

### 6. Ranking implemented as a pure sort in the frontend

A `rankSearchResults` function in a new `search.go` file takes `[]*pb.Product` and the query string, assigns each product a rank tier (1–6 per R-08), and returns a stable-sorted slice. Keeping this in its own file makes it independently unit-testable.

---

## Architecture: Files to Change

### New files

| File | Purpose |
|---|---|
| `src/frontend/search.go` | `rankSearchResults` (sorting by R-08 tiers) and `highlightQuery` (safe `<mark>` injection) — pure functions |
| `src/frontend/templates/search.html` | Search results page: query echo, result count, product cards with highlighted name/description, empty state, short-query validation message |

### Modified files

| File | Change |
|---|---|
| `src/frontend/rpc.go` | Add `searchProducts(ctx, query) ([]*pb.Product, error)` — one wrapper call, same pattern as `getProducts` |
| `src/frontend/handlers.go` | Add `searchHandler` — reads `q`, validates length, calls `searchProducts`, ranks results, converts currency, renders `search` template |
| `src/frontend/main.go` | Register `GET /search` route before the static file handler |
| `src/frontend/templates/header.html` | Add search `<form method="GET" action="/search">` with text input and submit button inside the `.controls` div |

---

## Detailed Component Design

### `search.go` — ranking and highlighting

```
// Rank tiers (lower = higher priority)
// 1: query matches start of a word in product Name
// 2: query matches middle of a word in product Name (not start)
// 3: query matches end of a word in product Name
// 4: query matches start of a word in product Description
// 5: query matches middle of a word in product Description
// 6: query matches end of a word in product Description

func rankSearchResults(products []*pb.Product, query string) []*pb.Product
func wordMatchTier(text, query string, baseOffset int) int  // baseOffset: 0 for name, 3 for desc
func highlightQuery(text, query string) template.HTML
```

**Ranking algorithm:**
1. Normalise query and each word to lowercase.
2. For each product, find the lowest tier (best match) across all words in name, then description.
3. `sort.SliceStable` by tier ascending; products with the same tier retain catalogue order.
4. A product matched in both name and description is ranked by its name tier only.

**Word boundary definition:** A "word" is any whitespace- or punctuation-delimited token. For determining start/middle/end:
- Start: `strings.HasPrefix(strings.ToLower(word), lowerQuery)`
- End (not start): `strings.HasSuffix(strings.ToLower(word), lowerQuery)` and not start
- Middle: match exists in word but neither start nor end

**Highlighting:**
1. `html.EscapeString(text)` on the full text.
2. Find all case-insensitive occurrences using `strings.Index` in a loop.
3. Insert `<mark>` / `</mark>` around each match.
4. Return as `template.HTML`.

### `searchHandler` in `handlers.go`

```
GET /search?q=<query>

1. Read q = strings.TrimSpace(r.URL.Query().Get("q"))
2. If len([]rune(q)) < 2:
     render "search" template with {query: q, validation_error: "Please enter at least 2 characters to search"}
     return  (no gRPC call)
3. products, err := fe.searchProducts(r.Context(), q)
4. rank products with rankSearchResults(products, q)
5. for each product: convertCurrency → build productView{Item, Price, HighlightedName, HighlightedDescription}
6. render "search" template with {query, products, result_count, currencies, cart_size}
```

### `search.html` template structure

```
{{ define "search" }}
{{ template "header" . }}
<main>
  <div class="container">

    <!-- Search bar (repeat from header for page context) -->
    <form method="GET" action="{{ $.baseUrl }}/search" class="search-form">
      <input type="text" name="q" value="{{ $.query }}" ... />
      <button type="submit">Search</button>
    </form>

    <!-- Validation message (R-05) -->
    {{ if $.validation_error }}
    <p class="search-validation-error">{{ $.validation_error }}</p>
    {{ end }}

    <!-- Results header -->
    {{ if not $.validation_error }}
    {{ if $.products }}
    <p>{{ $.result_count }} result(s) for "{{ $.query }}"</p>
    {{ range $.products }}
    <div class="search-result-card">
      <a href="{{ $.baseUrl }}/product/{{ .Item.Id }}">
        <img src="{{ $.baseUrl }}{{ .Item.Picture }}" />
        <div class="search-result-name">{{ .HighlightedName }}</div>
        <div class="search-result-desc">{{ .HighlightedDescription }}</div>
        <div class="search-result-price">{{ renderMoney .Price }}</div>
      </a>
    </div>
    {{ end }}
    {{ else }}
    <p class="search-no-results">No products found for "{{ $.query }}".</p>
    {{ end }}
    {{ end }}

  </div>
</main>
{{ template "footer" . }}
{{ end }}
```

### Header search bar (addition to `header.html`)

Inserted inside `.controls` div, before the cart icon:

```html
<form method="GET" action="{{ $.baseUrl }}/search" class="search-form"
      onsubmit="return validateSearch(this)">
  <input type="text" name="q" placeholder="Search products…"
         value="{{ $.search_query }}" class="search-input" />
  <button type="submit" class="search-button">Search</button>
</form>
<script>
function validateSearch(form) {
  var q = form.q.value.trim();
  if (q.length < 2) {
    form.q.focus();
    return false;   // prevent navigation; server handles it if JS is off
  }
  return true;
}
</script>
```

`$.search_query` is populated by `injectCommonTemplateData` — on the search results page it mirrors the current query so the bar stays populated; on all other pages it is empty. `injectCommonTemplateData` reads it from the request URL query param.

---

## Data Flow

```
Browser
  │  GET /search?q=shirt
  ▼
frontend: searchHandler
  │  len(q) < 2? → render validation page
  │  searchProducts(ctx, q)
  ▼
productcatalogservice: SearchProducts(Query: "shirt")
  │  in-memory filter: strings.Contains on Name + Description
  ▼
frontend: rankSearchResults(products, "shirt")
  │  sort by tier (name-start → name-mid → name-end → desc-start → desc-mid → desc-end)
  │  convertCurrency for each product
  │  highlightQuery(name, "shirt") → template.HTML with <mark>
  │  highlightQuery(description, "shirt") → template.HTML with <mark>
  ▼
render search.html → HTTP 200
```

---

## What Is Intentionally Not Changed

- `productcatalogservice` source code — no changes.
- Proto / gRPC definitions — no changes.
- Any other service (`cartservice`, `currencyservice`, etc.) — no changes.
- `Dockerfile`, `kubernetes-manifests/`, `helm-chart/`, CI workflows — no changes.
- CSS framework or static asset pipeline — new styles added inline or as minimal additions to `styles.css` only.
