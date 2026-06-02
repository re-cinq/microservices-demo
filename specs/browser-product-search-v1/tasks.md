# Tasks: Browser Product Search v1

Feature directory: `specs/browser-product-search-v1`  
Spec: `spec.md` · Plan: `plan.md`

---

## Phase 1 — Foundational: search logic (no UI, no route)

- [x] [T01] Add `searchProducts` gRPC wrapper to `rpc.go` — `src/frontend/rpc.go`
- [x] [T02] Create `search.go` with `wordMatchTier` helper that returns a rank tier (1–6) for a given text and query string — `src/frontend/search.go`
- [x] [T03] Add `rankSearchResults` to `search.go` that stable-sorts `[]*pb.Product` by tier using `wordMatchTier`, name before description — `src/frontend/search.go`
- [x] [T04] Add `highlightQuery` to `search.go` that escapes input text, wraps each case-insensitive match in `<mark>…</mark>`, and returns `template.HTML` — `src/frontend/search.go`
- [x] [T05] Register `highlightQuery` in the template `FuncMap` alongside `renderMoney` and `renderCurrencyLogo` — `src/frontend/handlers.go`

## Phase 2 — Foundational: route and handler

- [x] [T06] Add `searchHandler` to `handlers.go`: read `q`, validate length < 2 (render validation page, no gRPC call), call `searchProducts`, rank results, convert currency per product, build `productView` slice with `HighlightedName` and `HighlightedDescription` fields, render `search` template — `src/frontend/handlers.go`
- [x] [T07] Register `GET /search` route in `main.go` before the `/static/` path prefix handler — `src/frontend/main.go`

## Phase 3 — User stories (priority order)

- [x] [T08] Create `templates/search.html`: `{{ define "search" }}` block with header, in-page search bar (query pre-filled), result count line, product card grid (image + highlighted name + highlighted description + price, each wrapped in `<a href="/product/{id}">`), no-results empty state, validation error message area, footer — `src/frontend/templates/search.html`
- [x] [T09] Add search bar `<form method="GET" action="…/search">` to `header.html` inside `.controls` div before the cart icon; include `onsubmit` JS guard that blocks submission if `q.length < 2` — `src/frontend/templates/header.html`
- [x] [T10] Populate `search_query` in `injectCommonTemplateData` from the request URL `q` parameter so the header search bar stays filled on the results page and is empty on all other pages — `src/frontend/handlers.go`

## Phase 4 — Polish and validation

- [x] [T11] Add minimal CSS for search bar, result cards, `<mark>` highlight, validation error, and no-results message — consistent with existing Bootstrap 4 + Cymbal styles — `src/frontend/static/styles/styles.css`
- [x] [T12] Verify SC-01: search input (`name="q"`) confirmed present in `header.html`, which is included on every page via `{{ template "header" . }}`
- [x] [T13] Verify SC-02 + SC-08: `rankSearchResults` puts tier-1 (name-word-start) results first; `wordMatchTier("Sunglasses","sun",0)` → tier 1; confirmed by code inspection
- [x] [T14] Verify SC-03: `highlightQuery` wraps matches in `<mark>…</mark>`; `search.html` renders `.HighlightedName` and `.HighlightedDescription` as `template.HTML`
- [x] [T15] Verify SC-04: `search.html` has `search-no-results` block rendering `{{ $.query }}` when `$.results` is empty
- [x] [T16] Verify SC-05: `searchHandler` returns 200 + validation message when `len([]rune(q)) < 2`; JS guard in `header.html` prevents navigation client-side
- [x] [T17] Verify SC-06: `search.html` wraps each card in `<a href="{{ $.baseUrl }}/product/{{ .Item.Id }}">` — standard anchor, browser "open in new tab" works natively
- [x] [T18] Verify SC-07: `searchHandler` calls `convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))`; `setCurrencyHandler` redirects to `Referer` so currency change on `/search?q=shirt` reloads the same URL with new cookie
- [x] [T19] Build check: Go not installed in this environment — all symbol references verified by grep against `genproto/demo.pb.go` and `demo_grpc.pb.go`; all imports confirmed present; no new dependencies introduced
