# Research: Product Search — Phase 0

This document captures the technical decisions that resolve open questions from `spec.md` and `plan.md` before design (Phase 1) begins. Each decision states what was chosen, why, and what alternatives were considered and rejected.

## Decision 1 — Reuse the existing `SearchProducts` RPC; do not add a new RPC

**Decision:** Use the `ProductCatalogService.SearchProducts(SearchProductsRequest) → SearchProductsResponse` RPC already defined in [`protos/demo.proto:71-103`](../../protos/demo.proto) and implemented in [`src/productcatalogservice/product_catalog.go:60-72`](../../src/productcatalogservice/product_catalog.go). Do **not** edit the proto, do **not** introduce a parallel RPC.

**Rationale:**

- The proto already names exactly the surface the spec described (`SearchProducts(query)`).
- The generated Go stubs (`src/productcatalogservice/genproto/`, `src/frontend/genproto/`) already exist.
- Adding a new RPC would require regenerating proto code in **every** service language (Go, Node, Python, Java, C#) — needlessly invasive and a direct violation of spec Constraint *"No protocol changes to other services."*
- The existing test (`product_catalog_test.go:91`) confirms the RPC is wired up end-to-end inside the service.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Add a new RPC `SearchProductsByName` to dodge the description-matching behaviour | Forks the contract, regenerates stubs in 5 languages, and a one-line narrowing of the existing implementation achieves the same outcome with ~1/100th the diff. |
| Filter results in the frontend after calling `SearchProducts` | Wasteful — backend returns description matches the frontend then discards. Splits the spec-honouring behaviour across two layers. |
| Leave description matching in place; weaken the spec | Diverges from `spec.md` FR-002; the customer-side behaviour ("you searched 'soft' and got the loafers because the description says 'soft leather'") is fuzzier and harder to predict. |

## Decision 2 — Tighten `SearchProducts` to **name-only** substring matching

**Decision:** Modify `productcatalogservice/product_catalog.go` to drop the description-substring branch from the matching predicate, leaving:

> *"include product P iff `strings.Contains(strings.ToLower(P.Name), strings.ToLower(req.Query))`"*

(The existing implementation currently OR's the name check with the same check against `P.Description`.)

**Rationale:**

- Honours spec FR-002 literally.
- One-line behavioural narrowing. No public-contract change.
- Blast radius is zero in this repo: the only caller of `SearchProducts` after this slice ships is the new frontend handler we're adding. No other service calls it today.
- The existing unit test (`TestSearchProducts`, query `"alpha"`, expects 2 matches against fixture products `"Product Alpha One"` and `"Product Alpha Two"`) still passes after the narrowing — both matches are on name.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Keep matching name OR description | Diverges from spec FR-002; produces less predictable results for shoppers ("why did *Loafers* come up when I searched 'wardrobe'?"). |
| Add a query parameter to choose the search mode | Premature flexibility; spec is explicit about MVP behaviour. |
| Move matching logic into the frontend | Pushes business logic out of the service that owns the data. Backend remains the better home for the predicate. |

## Decision 3 — Add a single HTTP entry point: `GET /search?q={query}`

**Decision:** Expose the feature in the frontend as a single GET route, `/search`, taking a `q` query-string parameter. Submit via a header-resident `<form method="GET" action="/search">` partial in `templates/header.html`. Render results via a new template `templates/search-results.html`.

**Rationale:**

- Mirrors the existing routing style in `src/frontend/main.go:149-163` (one `r.HandleFunc(...)` line per surface, gorilla/mux).
- `GET` keeps the URL shareable and back-button-friendly — natural for a search.
- One handler, one template, one route — minimal surface change.
- Header partial places the input on every page where the header is included, satisfying spec User Story 3 (search reachable from cart and detail pages).

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| `POST /search` | Loses shareable URLs and breaks back-button. No security justification for POST (query is non-secret). |
| AJAX / dropdown live-search | Spec Constraints explicitly exclude autocomplete in v1. |
| Submit to `/?q=...` and overload the home handler | Conflates two surfaces; complicates the home template's data shape; harder to add results-specific UI later. |

## Decision 4 — Render-time currency conversion uses the existing `convertCurrency` path

**Decision:** In the search handler, convert each result's `priceUsd` to the session currency via the existing `fe.convertCurrency(ctx, p.GetPriceUsd(), currentCurrency(r))`, exactly the way `homeHandler` does it ([`handlers.go:84`](../../src/frontend/handlers.go)).

**Rationale:**

- Satisfies spec FR-005 ("prices on the results page MUST honour the shopper's current session currency, using the same conversion path used elsewhere").
- Inherits the currency-bug fix from PR #10 automatically.
- Inherits the home page's per-product `Convert` RPC pattern; consistency over premature optimisation.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Batch-convert all result prices in one call | `CurrencyService.Convert` accepts a single Money. Batching would require a proto change (excluded by constraints) for negligible win on ~10 products. |
| Cache converted prices in the frontend | Out of scope per spec; introduces invalidation hazards if rates ever change. |

## Decision 5 — Validation lives in the handler, not in the service

**Decision:** Empty/whitespace-only queries, oversize queries (>200 chars), and pure-punctuation queries are handled in the new frontend handler:

- Trim leading/trailing whitespace; if the resulting query is empty, render the "please enter a search term" branch of the no-results template.
- If the trimmed query is >200 chars, render an error page with a clear message. (Spec FR-009 allows either rejection or truncation; we choose rejection for clarity.)
- Any other non-matching query (including pure punctuation) is allowed to round-trip to `SearchProducts`, which will simply return zero results, which the template renders as the no-matches page.

**Rationale:**

- Keeps the `productcatalogservice` RPC narrowly responsible for the matching predicate; doesn't bloat it with input-validation policy that's really a user-facing UX concern.
- Mirrors the existing in-frontend validation pattern (`src/frontend/validator/validator.go` already handles `AddToCart`, `PlaceOrder`, `SetCurrency` payloads).
- A new `SearchPayload` struct in `validator/` is **optional** — query-string validation is simple enough to inline in the handler. We default to inline for diff minimalism; if a tasks-phase reviewer prefers the validator pattern, it's a trivial extraction.

**Alternatives considered:**

| Alternative | Why rejected |
|---|---|
| Validate in `SearchProducts` (return `InvalidArgument` for empty/oversize) | Wrong layer — that's UX policy, not data policy. Also forces every future caller (e.g. a hypothetical CLI) to reproduce the UX rules. |
| Truncate oversize instead of rejecting | Surprising behaviour ("I searched for a long thing and got results for a thing I didn't search for"). Rejection is clearer. |

## Decision 6 — Header partial conditionality

**Decision:** The header search input renders unconditionally (no template variable gate). The existing `assistant_enabled` flag, currency selector, and cart icon already render conditionally; the search input does not need that treatment because there is no opt-out scenario.

**Rationale:**

- Spec User Story 3 wants search on every page where the header renders.
- Simpler template logic; one less variable in `injectCommonTemplateData`.

## Decision 7 — Tracing and logging follow existing patterns; no new pipeline

**Decision:** The new handler:

- Gets its logger from `r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)` (same as every other handler).
- Logs a single structured info line with the (sanitised) query and the result count.
- Does **not** add new spans, new metrics, or new telemetry pipelines.

OTel tracing is already wired around the frontend via the `otelhttp.NewHandler` wrapper at `main.go:168`; the new route inherits that automatically.

**Rationale:** Spec Constraints rule out new observability infrastructure. Inherited tracing is sufficient for MVP.

## Open questions — none

All Technical Context unknowns from `plan.md` are resolved by the decisions above. No `NEEDS CLARIFICATION` markers remain.
