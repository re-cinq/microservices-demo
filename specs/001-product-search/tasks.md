# Tasks: Product Search

**Input**: Design documents from `specs/001-product-search/`  
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/search-rpc.md](contracts/search-rpc.md), [quickstart.md](quickstart.md)

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies on an incomplete task)
- **[Story]**: Which user story this task belongs to (US1 / US2 / US3)

---

## Phase 1: Setup

**Purpose**: No new project structure is needed — this feature modifies existing
files only. Confirm working state before making changes.

- [x] T001 Verify branch is `attendee/kate-payne-fix` and the repo builds (`go build ./...` in both `src/frontend` and `src/productcatalogservice`)

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: Fix the `SearchProducts` RPC so it matches against both `Name` and
`Description`. All three user stories depend on this being correct.

**⚠️ CRITICAL**: No user story work can begin until this task is complete.

- [x] T002 Fix `SearchProducts` in `src/productcatalogservice/product_catalog.go` lines 60–71: extract `query := strings.ToLower(req.Query)` and change the `if` condition to `strings.Contains(strings.ToLower(product.Name), query) || strings.Contains(strings.ToLower(product.Description), query)`

**Checkpoint**: `SearchProducts("sunglasses")` now returns a product whose description (not name) contains "sunglasses". The service still compiles.

---

## Phase 3: User Story 1 — Basic Keyword Search (Priority: P1) 🎯 MVP

**Goal**: Shoppers can type a term into a search box on the catalogue page, submit
it, and see only matching products. Queries shorter than 2 non-whitespace
characters are blocked with an inline message; the catalogue is unchanged.

**Independent Test**: Load `/`, enter `"camera"` → only camera products shown.
Enter `"a"` → full catalogue stays, error message visible. Enter `"xyzzy123"` → no
products (empty list, no error yet — empty-state copy comes in US2). Entering
`"CAMERA"` returns same results as `"camera"`.

### Implementation for User Story 1

- [x] T003 [P] [US1] Add `searchProducts(ctx, query)` gRPC wrapper to `src/frontend/rpc.go` calling `pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).SearchProducts(ctx, &pb.SearchProductsRequest{Query: query})` and returning `resp.GetResults(), err`
- [x] T004 [P] [US1] Extend `homeHandler` in `src/frontend/handlers.go` to: read `q := r.URL.Query().Get("q")`; compute `trimmed := strings.TrimSpace(q)`; use a `switch` to call `fe.getProducts` (no q), block with `searchError` (len < 2), or call `fe.searchProducts` (valid); pass `search_query`, `search_error`, and `is_search` keys to `ExecuteTemplate` alongside existing keys (depends on T003)
- [x] T005 [US1] Add search form to `src/frontend/templates/home.html` above the `hot-products-row` div: `<form method="GET" action="{{ $.baseUrl }}/">` containing a text input `name="q"` with `value="{{ $.search_query }}"` and `minlength="2"`, a Submit button, a conditional Clear link (`{{ if $.search_query }}<a href="{{ $.baseUrl }}/">Clear</a>{{ end }}`), and a conditional error `{{ if $.search_error }}<div class="text-danger ...">{{ $.search_error }}</div>{{ end }}`; also replace the static `<h3>Hot Products</h3>` with `{{ if $.is_search }}<h3>Results for &ldquo;{{ $.search_query }}&rdquo;</h3>{{ else }}<h3>Hot Products</h3>{{ end }}` (depends on T004)

**Checkpoint**: US1 fully functional. `/?q=camera` filters results. `/?q=a` shows full catalogue + error. `/?q=CAMERA` returns same results as `/?q=camera`. Clear link returns to `/`. No highlights yet (US3); no empty-state copy yet (US2).

---

## Phase 4: User Story 2 — Empty-State "No Results" Message (Priority: P2)

**Goal**: When a valid search returns zero products, the page shows a clear
message that includes the exact search term, so the shopper knows the search
worked and can try again.

**Independent Test**: Submit `/?q=xyzzy123` → no product cards shown; message visible
containing the literal text `"xyzzy123"`. Clear → full catalogue restored.

### Implementation for User Story 2

- [x] T006 [US2] Add empty-state block to `src/frontend/templates/home.html` inside the product row, after the `{{ range $.products }}...{{ end }}` block: `{{ if and $.is_search (eq (len $.products) 0) }}<div class="col-12 text-center py-5"><p class="text-muted">No products found for &ldquo;{{ $.search_query }}&rdquo;.</p><a href="{{ $.baseUrl }}/">View all products</a></div>{{ end }}` (depends on T005)

**Checkpoint**: US1 + US2 both functional. `/?q=xyzzy123` shows empty-state with the term echoed. `/?q=camera` still filters correctly.

---

## Phase 5: User Story 3 — Search-Term Highlighting (Priority: P3)

**Goal**: Every occurrence of the search term within product names and
descriptions in the results list is visually distinct (wrapped in `<mark>`).
Highlighting is case-insensitive and covers all occurrences, including substrings
(e.g., `"lens"` within `"lenses"`). XSS is impossible even if the term contains
HTML special characters.

**Independent Test**: Search `/?q=lens` → `"lens"` highlighted wherever it appears
in names and descriptions. Search `/?q=<script>` → rendered as escaped literal
text, no JS executes. `"lenses"` shows the `"lens"` substring highlighted.

### Implementation for User Story 3

- [x] T007 [P] [US3] Implement `highlightTerm(term, text string) template.HTML` function in `src/frontend/handlers.go`: add `"html"` to the import block; write the function that calls `html.EscapeString(text)`, walks the escaped string with `strings.Index` on lowercased copies, and wraps each match in `<mark>...</mark>`, returning `template.HTML`; register it in the `template.FuncMap` as `"highlightTerm"` alongside the existing `renderMoney` and `renderCurrencyLogo` entries
- [x] T008 [US3] Update product name rendering in `src/frontend/templates/home.html`: replace `{{ .Item.Name }}` in the `hot-product-card-name` div with `{{ if $.search_query }}{{ highlightTerm $.search_query .Item.Name }}{{ else }}{{ .Item.Name }}{{ end }}` (depends on T007)
- [x] T009 [US3] Add description rendering with highlighting to the product card in `src/frontend/templates/home.html`: below the name div, add `{{ if and $.is_search $.search_query }}<div class="hot-product-card-desc small text-muted">{{ highlightTerm $.search_query .Item.Description }}</div>{{ end }}` (depends on T007, T008)

**Checkpoint**: All three user stories functional. `/?q=lens` highlights `"lens"` in both names and descriptions. `/?q=<script>alert(1)</script>` renders safely.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Test coverage for the service-level fix and final manual verification.

- [x] T010 [P] Add a Go test case to `src/productcatalogservice/product_catalog_test.go` that calls `SearchProducts` with a query string that matches only a product's `Description` field (not its `Name`) and asserts the product is present in `Results`
- [ ] T011 Run the manual verification checklist in [quickstart.md](quickstart.md) end-to-end, covering all scenarios in the table (camera filter, case-insensitive, single-char block, whitespace block, no-results empty state, XSS term, description-only match, clear button)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **blocks all user stories**
- **Phase 3 (US1)**: Depends on Phase 2 — MVP deliverable
- **Phase 4 (US2)**: Depends on Phase 3 (needs `is_search`, `products`, and the search form in the template)
- **Phase 5 (US3)**: Depends on Phase 3 (needs `search_query` in template context)
- **Phase 6 (Polish)**: Depends on Phases 3–5

### Within Phase 3 (US1)

- **T003** (rpc.go) and **T004** (handlers.go) are independent — different files, run in parallel
- **T005** (home.html) depends on T004 (needs the template keys T004 introduces)

### Parallel Opportunities

```
T003 (rpc.go)      ──┐
                      ├── both done → T005 (home.html form)
T004 (handlers.go) ──┘

T007 (highlightTerm fn)  →  T008 (name highlight)  →  T009 (desc highlight)

T010 (unit test)   ──┐
                      ├── independent of each other
T011 (manual QA)   ──┘
```

---

## Implementation Strategy

### MVP (User Story 1 only)

1. Complete T001 (verify build)
2. Complete T002 (fix SearchProducts)
3. Complete T003 + T004 in parallel, then T005
4. **STOP**: validate `/?q=camera` filters, `/?q=a` blocks, `/?q=CAMERA` matches

### Incremental Delivery

1. T001 → T002 → T003+T004 → T005 → **US1 done** ✓
2. T006 → **US2 done** ✓ (empty state on top of US1)
3. T007 → T008 → T009 → **US3 done** ✓ (highlighting on top of US1+US2)
4. T010 + T011 → **Polish done** ✓

Each step is independently deployable and demonstrable.

---

## Notes

- `[P]` tasks = different files, no incomplete-task dependencies — safe to parallelise
- `[Story]` label maps each task to a user story for traceability
- No new files, routes, proto changes, or build-pipeline changes anywhere in this list
- `highlightTerm` must use `html.EscapeString` on input before building output — never pass raw user input as `template.HTML` directly
- `len([]rune(trimmed))` not `len(trimmed)` for correct Unicode character counting in the min-length guard
