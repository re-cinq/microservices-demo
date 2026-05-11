# Feature Specification: Product Search

**Feature Branch**: `001-product-search`
**Created**: 2026-05-11
**Status**: Draft
**Input**: User description: "Add a product search feature to Online Boutique: A search box that calls a new SearchProducts(query) RPC on productcatalogservice returning matches by product name (case-insensitive substring)."

## Problem *(in the user's words)*

> "I came to Online Boutique knowing I wanted a watch. The only way to find one was to scroll the home page until I saw it. There's no search box anywhere. If I'd been looking for something the home page doesn't feature, I'd have given up."

Shoppers arrive at the storefront with **intent** — a product name they've heard of, an item they remember from a previous visit, a category they want to scan. The site forces all of them into the same browse-and-scroll flow because there is no search affordance. The MVP closes this gap with the smallest credible slice: a single text box, name-only matching, no ranking heuristics, no infrastructure changes.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Find a product by typing part of its name (Priority: P1)

A shopper opens Online Boutique looking for a watch. From the home page they enter `watch` (or `Watch`, or `WATCH`) into a visible search input and submit. They see a results page listing every product whose name contains that substring, with picture, name, and the localised price. They click a result and arrive on the existing product detail page.

**Why this priority**: This is the entire MVP. Without it there is no feature. Every later story is additive.

**Independent Test**: Can be fully tested by entering a known catalog name fragment (`sun`, `watch`, `loaf`, `dryer`) into the search input and confirming the matching product(s) appear with the correct picture, name, and price, and that the result link reaches the product detail page. Delivers value the first time it runs: shoppers can now find a known item without scrolling.

**Acceptance Scenarios**:

1. **Given** a shopper is on the home page, **When** they type `watch` into the search input and submit, **Then** the results page lists the *Watch* product with its picture, name, and localised price.
2. **Given** the catalog name contains the query as a substring at any position (start, middle, end), **When** the shopper submits the query, **Then** that product appears in results.
3. **Given** the shopper submits the query in any letter casing — `Watch`, `watch`, `WATCH`, `wAtCh` — **When** the search runs, **Then** they all return the same set of results.
4. **Given** the shopper has selected JPY as the session currency, **When** they view the results page, **Then** prices are rendered in JPY using the same conversion path used elsewhere in the storefront.
5. **Given** the shopper clicks a result, **When** the navigation completes, **Then** they land on the existing `/product/{id}` detail page for that product, unchanged from today.

---

### User Story 2 — Empty- and no-match queries are handled gracefully (Priority: P2)

A shopper submits a query that doesn't match any product, or submits an empty query, and is shown an explanatory page rather than a broken or blank one.

**Why this priority**: Required for production-credibility but not for proving the feature works at all. Without it, the feature appears half-finished but the happy path is not blocked.

**Independent Test**: Submit `"zzzzz"` and confirm a "no results" page renders; submit an empty query and confirm the user is either prevented from submitting or shown a "please enter a search term" message. No dependency on Story 1's matching specifics — only its plumbing.

**Acceptance Scenarios**:

1. **Given** no product name contains the query as a substring, **When** the shopper submits the search, **Then** the results page renders with a clear "no matches for *{query}*" message and a link back to the home page.
2. **Given** the shopper submits an empty or whitespace-only query, **When** they hit Enter, **Then** the request is rejected client-side or the server returns a friendly "please enter a search term" page — **never** the full catalog and **never** a 5xx.
3. **Given** the shopper submits a query containing only punctuation, **When** they hit Enter, **Then** the response is the same "no results" page used for any other non-matching query — not a 5xx, not the full catalog.

---

### User Story 3 — Search is reachable from every page (Priority: P3)

The search input lives in the global header so a shopper who navigated to a product detail page, the cart, or checkout can search again without going back to the home page.

**Why this priority**: A quality-of-life improvement. Story 1 covers the primary flow; this raises discoverability and re-search convenience.

**Independent Test**: From the cart page, type a term into the header search and submit; confirm the results page renders identically to a home-page-initiated search.

**Acceptance Scenarios**:

1. **Given** the shopper is on the cart page, **When** they use the header search input, **Then** they reach the same results page they would from the home page.
2. **Given** the shopper is on a product detail page, **When** they search, **Then** the current product is not given preferential treatment in results — ordering is the same regardless of the page the search was initiated from.

---

### Edge Cases

- **Case sensitivity**: `Sunglasses`, `sunglasses`, and `SUNGLASSES` must all return the same results.
- **Substring position**: queries match name substrings at any position (`gun` matches `Sunglasses` — yes, because `gun` is in `Sunglasses`).
- **Leading/trailing whitespace**: trimmed before matching; `   watch   ` behaves like `watch`.
- **Punctuation in queries**: tolerated; if the punctuation isn't in any product name, it simply contributes to a no-results outcome (no 5xx).
- **Unicode**: non-ASCII queries do not crash the service — they either match Unicode characters present in catalog names (unlikely for the current English catalog) or return zero results.
- **Long queries**: queries longer than a documented cap (e.g. 200 characters) are rejected client-side or truncated server-side with a clear message; they do not become a DoS vector.
- **Concurrent searches**: search load does not noticeably degrade the performance of the existing catalog browse / cart / checkout flows.
- **Currency switching on a results page**: re-renders the same results with prices in the newly selected currency, identical to every other product-listing surface.

## Constraints *(explicit — what we will NOT do in this slice)*

These are the user-supplied technical constraints, quoted verbatim and treated as binding:

- **Use only the services already in this repo. Do NOT add new services.** The 11-service topology stays at 11. Search is exposed by extending an existing service (the product catalog service) and the existing frontend.
- **Do NOT introduce Elasticsearch, Solr, vector databases, or any new datastore.** No new persistence layer, no new search engine, no new cache.
- **Use the existing in-memory product catalogue loaded from `productcatalogservice/products.json`. Filter in memory.** The static JSON file remains the sole source of truth for the catalog.
- **Match the language of the service you're editing: frontend is Go, productcatalogservice is Go. Use the existing protobuf/gRPC patterns if you extend productcatalogservice.** No new language or transport gets added to either service.
- **The CI deploys whatever lands on `attendee/<your-name>`. Do NOT add new infrastructure config, Helm charts, manifests, or environment variables.** No edits under `kubernetes-manifests/`, `kustomize/`, `helm-chart/`, `istio-manifests/`, or `terraform/` as part of this slice. No new env vars in deployments.
- **Stay inside this branch and this repo. Do NOT touch the build pipeline.** No CI/CD edits, no GitHub Actions edits, no Cloud Build edits, no Dockerfile changes that aren't strictly required by the in-service code change.

Additional scope bounds (to keep the MVP small):

- **Name-only matching.** Description fields, category fields, and SKU IDs are not searched in v1.
- **No autocomplete / type-ahead / suggestions.** The MVP is a submit-to-search input, not a live-search dropdown.
- **No user accounts, saved searches, or search history.** The site is anonymous and stays that way.
- **No personalised or behaviour-based ranking.** Results are ordered deterministically (e.g. the catalog's existing order, filtered) — not via the recommendation service or any per-user signal.

## Requirements *(observable behaviour, not implementation)*

### Functional Requirements

- **FR-001**: The system MUST present a visible search input on the home page that accepts free-text input and submits on Enter or button click.
- **FR-002**: Given a non-empty query, the system MUST return a results page listing every catalog product whose **name** contains the query as a **case-insensitive substring** after trimming surrounding whitespace. Description, category, and ID fields are NOT searched.
- **FR-003**: The results page MUST show, for each matching product, the same visual elements the home-page grid shows: picture, product name, and localised price.
- **FR-004**: Each result MUST link to the existing product detail page for that product, using the existing URL scheme (`/product/{id}`), with no change to the detail route.
- **FR-005**: Prices on the results page MUST honour the shopper's current session currency, using the same conversion path used elsewhere in the storefront.
- **FR-006**: When a query returns zero matches, the system MUST render an explicit "no matches for *{query}*" page including the submitted query in the message and a link back to the home page.
- **FR-007**: The system MUST reject empty or whitespace-only queries — either client-side or with a friendly server-rendered message — and MUST NOT return the full catalog as a fallback.
- **FR-008**: The system MUST handle Unicode and arbitrary punctuation in queries without raising a 5xx error. The acceptable behaviour for non-matching characters is zero results.
- **FR-009**: The system MUST cap query length at a documented upper bound (e.g. 200 characters) and reject or truncate longer queries with a clear message.
- **FR-010**: Results ordering MUST be deterministic for the same catalog and the same query. The catalog's existing order (filtered to matches) is the v1 ordering rule. No per-session, per-user, or randomised ranking.
- **FR-011**: The search input MUST be reachable from the cart page and from product detail pages (User Story 3 scope) — i.e. it lives in the global header, not only on the home page.
- **FR-012**: The product catalog's source-of-truth (the existing `productcatalogservice/products.json`) MUST remain the sole catalog data source. No parallel index, no duplicated catalog file, no new schema.
- **FR-013**: The search backend MUST be exposed by an existing in-repo service. No new service, no new datastore, no new external dependency (see Constraints).

### Key Entities

- **Product** *(existing)* — already defined by the product catalog service. Has `id`, `name`, `description`, `picture`, `priceUsd`, `categories`. **No schema change required.**
- **SearchQuery** *(new, transient)* — a single string submitted by a shopper. Has no identity, is not persisted, exists only for the duration of one request.
- **SearchResult** *(new, transient)* — an ordered list of references to existing `Product` records, in the catalog's existing order, filtered to those whose `name` matches the `SearchQuery` by case-insensitive substring. Not persisted; rendered, then discarded.

## Success Criteria *(mandatory, testable, technology-agnostic)*

### Measurable Outcomes

- **SC-001**: A shopper who knows the name of a product in the catalog can land on its detail page in **3 interactions or fewer** (open site → type → submit → click result = 3) — measurable by walkthrough timing.
- **SC-002**: For every product in the catalog, searching for any single word that appears as a substring of the product's name returns that product — **100% recall on name-substring matching**, measurable by an automated check that iterates the catalog and runs one search per product.
- **SC-003**: Median end-to-end time from search submission to results-page render is **under 500 ms** on the demo cluster under nominal load.
- **SC-004**: The feature introduces **zero new 5xx error classes** in the existing storefront — measurable by comparing storefront error rates for one hour before and after rollout.
- **SC-005**: A zero-match query renders the "no results" page (not a blank page, not a 5xx) in **100% of attempted bad queries** during a documented acceptance pass (a fixed test set of at least 20 strings, including empty, whitespace-only, oversize, Unicode, and pure-punctuation cases).
- **SC-006**: Switching the session currency on a results page changes the displayed prices to the new currency for **every** result, with no stale prices remaining — verifiable by visual diff.
- **SC-007**: Search adds **no measurable regression** (≤ 5% change, p95) to the latency of `/`, `/product/{id}`, `/cart`, and `/cart/checkout` page loads under the same load profile, before vs. after rollout.
- **SC-008**: Zero files are added or modified under `kubernetes-manifests/`, `kustomize/`, `helm-chart/`, `istio-manifests/`, `terraform/`, `.github/workflows/`, or `cloudbuild.yaml` as part of the change — measurable by `git diff` on the merged PR.

## Assumptions

- **Catalog size remains small.** The current catalog is a single static JSON file with on the order of 10 products. In-memory substring matching against it is cheap; SC-003's 500 ms budget assumes this scale.
- **Shoppers search in English.** The catalog ships with English names; non-English queries returning zero results is acceptable v1 behaviour.
- **The session currency mechanism continues to work correctly** — i.e. the currency conversion fix on `attendee/tim-love-fix` (see [PR #10](https://github.com/re-cinq/microservices-demo/pull/10)) is in place before this feature ships.
- **The product detail page works.** The `/product/{id}` route is already correct; this feature does not need to verify or fix it.
- **The frontend rendering pipeline (Go templates) can host a new results template** without architectural changes.
- **The existing protobuf/gRPC build chain in `productcatalogservice` works as today.** Regenerating `genproto` after a `.proto` edit is a normal in-repo workflow, not a build-pipeline change.
- **No SLA or contractual commitment** about search quality (relevance, ranking precision) — this is an MVP for a demo app, not a search product.
- **Observability** is limited to existing log/trace patterns. A structured log line per search is fine; building dashboards or new metrics pipelines is out of scope.
