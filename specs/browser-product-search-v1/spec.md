# Feature Spec: Browser Product Search v1

## Problem (in the user's words)

Right now, if I want to find a specific product in the Online Boutique, I have to scroll through the entire home page and look manually. There is no way to type a word and get matching products. I want to type something like "sunglasses" or "shirt" into a search box, hit enter, and see only the products that match — with the word I searched for highlighted so I can tell at a glance why each result appeared. If nothing matches, I want a clear message telling me so, not a blank page.

---

## Constraints (what NOT to do)

- Do not add new services. Use only the services already in this repo.
- Do not introduce Elasticsearch, Solr, vector databases, or any external datastore.
- Do not read from anything other than the existing in-memory product catalogue loaded from `productcatalogservice/products.json`.
- Do not add new infrastructure config, Helm charts, Kubernetes manifests, or environment variables.
- Do not touch the build pipeline or CI configuration.
- Do not change any service other than `frontend` and (if needed) `productcatalogservice`. Match the existing language (Go) and existing protobuf/gRPC patterns.

---

## Requirements (observable behaviour)

### R-01 — Search entry point
A search input is visible in the navigation area on every page of the site. The user can type a query and submit it from any page.

### R-02 — Keyword matching
Submitting a query returns every product where the exact sequence of characters the user entered appears as a contiguous substring within the product name or within any word in the product description, matched case-insensitively. For example, searching "sun" matches "Sunglasses" (start of name) and any description containing a word like "sunshine" or "sunscreen" (start, middle, or end of a word). Only exact contiguous character sequences count — no fuzzy matching, no stemming.

### R-03 — Highlighted matches
On the results page, every occurrence of the matched character sequence within a product's name and description is visually distinguished (e.g. bold or highlighted) so the user can immediately see why that product appeared.

### R-04 — Empty state
When no products match the query, the page shows a clear "no results" message that includes the original query string, rather than a blank or broken page.

### R-05 — Short query validation
If the user submits a search with fewer than 2 characters, no backend call is made. The search page (or home page) displays an inline message — for example, "Please enter at least 2 characters to search" — directly beneath the search bar, without navigating away. The search input retains whatever the user typed.

### R-06 — Currency consistency
Prices shown on the search results page reflect whichever currency is currently selected. If the user changes currency while on the results page, the prices update to the newly selected currency without losing the search results or query.

### R-07 — Navigation to product
Each search result card allows the user to open the full product detail page either by clicking the card or by using the browser's "open link in new tab" context menu. Both actions must work correctly.

### R-08 — Result ordering
Results are ranked so that the most specific matches appear first, in the following priority order:

1. Products where the query matches at the **start** of the product **name** (or the start of any word in the product name).
2. Products where the query matches in the **middle** of the product **name** (not at the start of any word).
3. Products where the query matches at the **end** of the product **name** (i.e. the matched sequence ends at the end of a word).
4. Products where the query matches at the **start** of any word in the product **description**.
5. Products where the query matches in the **middle** of a word in the product **description**.
6. Products where the query matches at the **end** of a word in the product **description**.

Within each priority tier, products retain their original catalogue order. Products matched by both name and description appear only once, ranked by their name-match position.

---

## Success Criteria (testable)

| # | Criterion | How to verify |
|---|-----------|---------------|
| SC-01 | Search input is present in the HTML `<nav>` on the home, product, and cart pages | Inspect page source / automated HTML assertion |
| SC-02 | `GET /search?q=sun` returns HTTP 200 and includes a product whose name starts with "Sun" as the first result | HTTP request + first result card assertion |
| SC-03 | The matched character sequence is wrapped in a distinguishing HTML element (e.g. `<mark>`) within the rendered product name and description | HTML assertion on search results page |
| SC-04 | `GET /search?q=zzznomatch` returns HTTP 200 and the response body contains "no results" (case-insensitive) and the string "zzznomatch" | HTTP request + response body assertion |
| SC-05 | Submitting a 1-character query does not navigate to `/search`; the current page displays an inline validation message containing "at least 2" | Browser behaviour assertion / HTML assertion |
| SC-06 | Each product card on the results page contains an `<a href="/product/{id}">` link that can be opened in a new tab | HTML assertion; right-click → open in new tab works |
| SC-07 | With currency set to EUR, prices on `/search?q=shirt` are displayed with the EUR symbol; switching to USD updates the prices without losing results | Visual inspection / currency toggle test |
| SC-08 | Name-match results appear before description-only-match results in the returned HTML | HTML order assertion: first card must be a name match when both name and description matches exist |
| SC-09 | No files outside `src/frontend/` and `src/productcatalogservice/` are modified | `git diff --name-only` |

---

## Acceptance Tests (Given / When / Then)

### AT-01 — Search with results and highlighting
**Given** I am on the home page  
**When** I type "sunglasses" into the search bar and submit  
**Then** the browser navigates to `/search?q=sunglasses`  
**And** the page shows one or more product cards  
**And** the character sequence "sunglasses" (case-insensitive) is visually highlighted within the name or description of each result card  

### AT-02 — No results empty state
**Given** I am on any page with the search bar visible  
**When** I type "zzznomatch" into the search bar and submit  
**Then** the browser navigates to `/search?q=zzznomatch`  
**And** the page shows a message containing "no results" (or equivalent) and the text "zzznomatch"  
**And** no product cards are rendered  

### AT-03 — Short query validation (1 character)
**Given** I am on any page with the search bar visible  
**When** I type "a" into the search bar and submit  
**Then** the browser does not navigate away  
**And** an inline message appears near the search bar containing "at least 2"  
**And** the search input still shows "a"  
**And** no request is made to `productcatalogservice`  

### AT-04 — Result card opens product detail in same tab
**Given** I search for "shirt" and results are shown  
**When** I click a product card  
**Then** I am taken to the product detail page for that specific product  

### AT-05 — Result card opens product detail in new tab
**Given** I search for "shirt" and results are shown  
**When** I open a product card link in a new tab  
**Then** the product detail page for that specific product opens in the new tab  

### AT-06 — Currency reflected in results
**Given** I have selected EUR as my currency  
**When** I search for "shirt"  
**Then** prices on the results page are shown in EUR  

### AT-07 — Currency change on results page
**Given** I am viewing search results for "shirt" with USD selected  
**When** I change the currency to EUR  
**Then** the prices on the results page update to EUR  
**And** the search results and query string remain unchanged  

### AT-08 — Search bar available on product page
**Given** I am viewing a product detail page  
**When** I look at the navigation area  
**Then** the search input is present and submitting a query navigates to the results page  

### AT-09 — Partial contiguous match
**Given** I search for "sun"  
**Then** products whose name or description contains the contiguous sequence "sun" (case-insensitive) appear in the results  
**And** products with no such sequence do not appear  

### AT-10 — Result ordering: name matches before description matches
**Given** the catalogue contains Product A (name contains the query) and Product B (description only contains the query)  
**When** I search for a term matching both  
**Then** Product A appears before Product B in the results list  

### AT-11 — Result ordering: word-start matches before mid-word matches
**Given** the catalogue contains Product A (query matches the start of a word in the name) and Product B (query matches the middle of a word in the name)  
**When** I search for that term  
**Then** Product A appears before Product B in the results list  
