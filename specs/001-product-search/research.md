# Research: Product Search

**Phase 0 output for**: `specs/001-product-search/plan.md`  
**Date**: 2026-05-11

---

## Finding 1 — SearchProducts RPC already exists

**Decision**: Re-use the existing `SearchProducts` RPC and its proto messages
(`SearchProductsRequest`, `SearchProductsResponse`) without modification.

**Rationale**: The RPC is already defined in `protos/demo.proto` (lines 71–75),
implemented in `src/productcatalogservice/product_catalog.go` (lines 60–71),
and the generated gRPC client stub is available in
`src/frontend/genproto/demo_grpc.pb.go`. Creating a new RPC would duplicate
the contract and require proto regeneration. Re-using the existing RPC is
zero-cost and fully consistent with the constraint "use existing protobuf/gRPC
patterns".

**Alternatives considered**:
- New `SearchProductsV2` RPC — rejected: unnecessary proto change, no new
  capability required.
- Filter in the frontend by calling `ListProducts` then filtering client-side —
  rejected: bypasses the service boundary, moves business logic into the
  presentation layer, and wastes a full catalogue transfer on every search.

---

## Finding 2 — Description field missing from filter (bug)

**Decision**: Fix `SearchProducts` in `product_catalog.go` to OR-match against
`product.Description` in addition to `product.Name`.

**Rationale**: The current implementation (line 65) only tests `product.Name`.
FR-006 requires the search to check both fields. The fix is one boolean
expression change; no new code paths, no new types.

**Current code (line 65)**:
```go
if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) {
```

**Fixed code**:
```go
query := strings.ToLower(req.Query)
if strings.Contains(strings.ToLower(product.Name), query) ||
    strings.Contains(strings.ToLower(product.Description), query) {
```

**Alternatives considered**:
- Regex matching — rejected: overkill for substring search, adds a new
  dependency risk, no requirement for pattern syntax.
- Searching additional fields (categories, ID) — rejected: out of spec scope.

---

## Finding 3 — Minimum-length guard: server-side in handler

**Decision**: Enforce the 2-character minimum in `homeHandler` in the frontend
(server-side), supplemented by `minlength="2"` and `required` on the HTML
`<input>` (client-side hint).

**Rationale**: The spec (FR-001a) requires the search to be *prevented* and
the catalogue to remain unchanged when the term is too short. Server-side
enforcement is mandatory because client-side HTML validation can be bypassed
(e.g., direct URL entry `/?q=a`). HTML `minlength` is added as a UX shortcut
only.

**Implementation**: In `homeHandler`, after reading `q := r.URL.Query().Get("q")`:

```go
trimmed := strings.TrimSpace(q)
if q != "" && len([]rune(trimmed)) < 2 {
    // render home with searchError: "Please enter at least 2 characters."
    // products = full catalogue (do not call SearchProducts)
}
```

Note: `len([]rune(...))` rather than `len(string)` to correctly count Unicode
code points (multi-byte characters each count as 1).

**Alternatives considered**:
- Client-side JS validation only — rejected: bypassable, server-side check
  required by FR-001a.
- Validate in productcatalogservice — rejected: leaks UI policy into a backend
  service; the service should remain a general-purpose filter.

---

## Finding 4 — Highlighting via template function

**Decision**: Implement highlighting as a Go template function
`highlightTerm(term, text string) template.HTML` registered in the frontend's
`template.FuncMap`.

**Rationale**: `html/template` auto-escapes all string output by default.
To inject `<mark>` tags around matches the function must return
`template.HTML` (a pre-sanitised type). The function must escape the input
text itself (using `html.EscapeString`) before wrapping matches, so XSS is
impossible even if the search term contains `<script>` or similar. This is
the idiomatic Go approach for trusted template output and requires no
additional dependencies.

**Implementation sketch**:
```go
func highlightTerm(term, text string) template.HTML {
    if term == "" {
        return template.HTML(html.EscapeString(text))
    }
    escaped := html.EscapeString(text)
    lowerEscaped := strings.ToLower(escaped)
    lowerTerm := strings.ToLower(html.EscapeString(term))
    var out strings.Builder
    for i := 0; i < len(escaped); {
        idx := strings.Index(lowerEscaped[i:], lowerTerm)
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

**Alternatives considered**:
- JavaScript client-side highlighting — rejected: requires JS execution, is
  harder to test server-side, and creates a flash of un-highlighted content.
- `strings.ReplaceAll` — rejected: does not handle case-insensitive matching
  while preserving original case of the matched text.

---

## Finding 5 — No new HTTP route needed

**Decision**: Extend the existing `homeHandler` (route `GET /`) to consume an
optional `?q=` query parameter rather than registering a new `/search` route.

**Rationale**: The spec says the catalogue page is the only surface for search.
Using `/?q=term` keeps search results on the catalogue page URL, is bookmarkable,
and requires zero routing changes. The handler already has all the wiring
(currency conversion, cart, template rendering) that search results also need.

**Alternatives considered**:
- New `GET /search` route + separate handler — rejected: duplicates the
  homeHandler setup boilerplate, introduces a second page the spec says is
  out of scope, and splits the catalogue view across two URLs.
