# Tasks: Browser Product Search

**Input**: Design documents from `/specs/001-browser-product-search/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)

---

## Phase 1: Foundation (one task — no new infra needed)

**Purpose**: Add the RPC wrapper that all stories depend on.

- [ ] T001 [US1] Add `searchProducts()` function to `src/frontend/rpc.go`
  - Follow the exact pattern of `getProducts()` (line 45–49)
  - Call `pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).SearchProducts(ctx, &pb.SearchProductsRequest{Query: query})`
  - Return `resp.GetProducts(), err`
  - ~5 lines of code

**Checkpoint**: `searchProducts()` compiles. No visible change yet — frontend doesn't call it.

---

## Phase 2: User Story 1 — Search for a product by name (P1) MVP

**Goal**: Customer types a query, sees matching products inline on the homepage.

**Independent Test**: Go to `/?q=shirt` — only matching products appear in the same card layout.

### T002–T003: Handler logic

- [ ] T002 [US1] Modify `homeHandler()` in `src/frontend/handlers.go` to read query param
  - Read `q := strings.TrimSpace(r.URL.Query().Get("q"))`
  - If `q != ""` → call `fe.searchProducts(r.Context(), q)` instead of `fe.getProducts(r.Context())`
  - If `q == ""` → call `fe.getProducts(r.Context())` (existing behaviour, unchanged)
  - Existing currency conversion loop stays untouched — it runs on whichever product list is returned
  - **Acceptance tests covered**: US1-1 (search returns matches), US1-3 (under 500ms), US1-4 (currency preserved)

- [ ] T003 [US1] Pass search query to template data in `src/frontend/handlers.go`
  - Add `"search_query": q` to the `map[string]interface{}` passed to `templates.ExecuteTemplate` (~line 110)
  - This enables the template to pre-fill the search box and show the "no results" message

### T004: Search form in template

- [ ] T004 [US1] Add search form to `src/frontend/templates/home.html`
  - Add an HTML `<form method="GET" action="{{ $.baseUrl }}/">` above the `<div class="row hot-products-row">` (~line 40)
  - Include `<input type="text" name="q" value="{{ $.search_query }}" placeholder="Search products...">` 
  - Form submits via Enter key (native HTML behaviour — FR-001)
  - Uses GET method with `?q=` parameter (FR-002 — bookmarkable, back-button friendly)
  - **Acceptance test covered**: US1-1, US1-2 (same card layout — no template changes to product cards)

**Checkpoint**: Search works end to end. Type `shirt` → see matching products. Empty query → all products. Currency selector still works with search. This is the MVP.

**Verify**:
- [ ] Visit `/` — all products shown, search box visible above grid (SC-003)
- [ ] Visit `/?q=shirt` — only matching products shown
- [ ] Visit `/?q=shirt` with EUR selected — prices in EUR (US1-4)
- [ ] Visit `/?q=xyznonexistent` — empty product grid (no results handling comes in Phase 3)

---

## Phase 3: User Story 2 — No results found (P1)

**Goal**: When search returns zero products, show a clear "no results" message instead of an empty grid.

**Independent Test**: Search for "xyznonexistent" — see "No products found" message, no empty grid.

- [ ] T005 [US2] Add empty-state conditional to `src/frontend/templates/home.html`
  - Inside the `hot-products-row` div, before the `{{ range $.products }}` loop:
    - `{{ if and (eq (len $.products) 0) (ne $.search_query "") }}`
    - Show: `<div class="col-12"><p>No products found for '{{ $.search_query }}'</p></div>`
    - `{{ else }}` → existing product card loop
    - `{{ end }}`
  - **Acceptance tests covered**: US2-1 (message shown), US2-2 (no broken layout)

**Checkpoint**: Search with no matches shows a clear message. Search with matches still shows product cards.

**Verify**:
- [ ] Visit `/?q=xyznonexistent` — "No products found for 'xyznonexistent'" message visible
- [ ] Visit `/?q=shirt` — product cards shown (no regression)
- [ ] Visit `/` — all products shown (no regression)

---

## Phase 4: User Story 3 — Clear search and return to full catalogue (P2)

**Goal**: Customer can clear search and see all products again.

**Independent Test**: Search for "shirt", then clear the box and submit — all products reappear.

- [ ] T006 [US3] Ensure empty-query submission returns all products
  - Already handled by T002 (empty `q` falls back to `ListProducts`)
  - The search form's `value="{{ $.search_query }}"` already clears when `q` is empty
  - **This task is verification only** — no code changes needed
  - **Acceptance test covered**: US3-1

**Checkpoint**: Full cycle works — search, see results, clear, see all products.

**Verify**:
- [ ] Search "shirt" → see filtered results → clear input → submit → all products return

---

## Phase 5: Edge Cases & Polish

**Purpose**: Verify edge cases from the spec. No code changes expected — these are validation tasks.

- [ ] T007 [P] Verify empty query: submit empty search box → all products shown
- [ ] T008 [P] Verify whitespace query: submit "   " → all products shown (TrimSpace in T002)
- [ ] T009 [P] Verify long query (500+ chars): paste long string → no error, no results or graceful response
- [ ] T010 [P] Verify special characters: search `<script>alert(1)</script>` → no XSS, safe rendering (Go html/template auto-escapes)
- [ ] T011 [P] Verify case insensitivity: "SHIRT", "Shirt", "shirt" all return same results
- [ ] T012 [P] Verify single-character query: search "s" → returns matching products

**Checkpoint**: All edge cases pass. Feature is complete.

---

## Dependencies & Execution Order

```
T001 (rpc wrapper)
  ↓
T002 → T003 (handler logic — sequential, same file)
  ↓
T004 (template search form) — can start after T003
  ↓  
T005 (empty state) — depends on T004
  ↓
T006 (clear/reset verification) — depends on T002
  ↓
T007–T012 (edge case verification — all parallel)
```

## Acceptance Test → Task Mapping

| Acceptance Test | Task(s) |
|---|---|
| US1-1: Search returns matching products | T001, T002, T004 |
| US1-2: Same card layout for results | T004 (no card template changes) |
| US1-3: Results in under 500ms | T001, T002 (existing RPC, in-memory search) |
| US1-4: Currency preserved in search | T002 (existing conversion loop untouched) |
| US2-1: "No results" message shown | T005 |
| US2-2: No broken layout on empty results | T005 |
| US3-1: Clear search shows all products | T002, T006 |
| Edge cases | T007–T012 |
