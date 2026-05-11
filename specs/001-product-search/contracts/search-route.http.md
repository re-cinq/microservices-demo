# Contract: `GET /search` (new HTTP route)

This contract documents the new HTTP entry point exposed by the `frontend` service.

## Endpoint

| Property | Value |
|---|---|
| Method | `GET` |
| Path | `{baseUrl}/search` |
| Registered in | `src/frontend/main.go` (mux router, alongside `/`, `/product/{id}`, `/cart`, etc.) |
| Handler | `frontendServer.searchHandler` in `src/frontend/handlers.go` |

The path includes the same `baseUrl` prefix every other frontend route does, picked up from `BASE_URL` (already in the existing env contract — no new env var added).

## Request

### Query parameters

| Name | Type | Required | Notes |
|---|---|---|---|
| `q` | string | yes | The shopper's search term. URL-decoded by `net/http`. The handler trims surrounding whitespace before processing. |

### Headers

No new request-header requirements. The existing session cookies (`shop_session-id`, `shop_currency`) are read as for every other page:

- `shop_session-id` — used for logging (session correlation); not used in the search itself.
- `shop_currency` — drives the currency the result-card prices are converted into, via the existing `convertCurrency` path.

### Validation (handler-side)

| Rule | Outcome on failure |
|---|---|
| `q` missing | Render the no-query branch (or redirect to `/`); never the full catalog. |
| `q` is empty or whitespace-only after trim | Render the no-query branch with a "please enter a search term" message. HTTP 200. |
| `q` length > 200 (after trim) | Render the generic error page with a "query too long" message. HTTP 400. |

(Validation rules taken from `spec.md` FR-007 and FR-009; rationale in `research.md` Decision 5.)

## Response

### Status codes

| Status | When |
|---|---|
| `200 OK` | Successful render — including the no-matches case and the empty-query case. |
| `302 Found` | Reserved; not used by this handler in MVP. |
| `400 Bad Request` | Query too long. |
| `500 Internal Server Error` | Upstream `SearchProducts` or `convertCurrency` returned an error. Rendered via the existing `renderHTTPError` path so the layout matches other 5xx pages. |

### Body

Content-Type: `text/html; charset=utf-8`, rendered from a new template `search-results.html`. The template receives the standard common data injected by `injectCommonTemplateData` plus the following feature-specific fields:

| Key | Type | Meaning |
|---|---|---|
| `query` | `string` | The (normalised) query the shopper submitted. Echoed in the page heading and in the no-matches message. |
| `items` | `[]searchResultView` | One entry per matched product. Each entry is `{ Item: *pb.Product, Price: *pb.Money }` — mirroring the `productView` struct used by `homeHandler`. |
| `no_matches` | `bool` | `len(items) == 0`. The template branches on this to render either the result grid or the no-matches message. |
| `show_currency` | `bool` | `true`. Allows the header currency selector to render on the results page exactly as it does on the home page. |
| `currencies` | `[]string` | Result of `fe.getCurrencies(ctx)`, same as `homeHandler` passes today. |
| `cart_size` | `int` | Result of `cartSize(cart)`, same as other pages — keeps the header cart indicator accurate. |

The handler does **not** invoke the ads service or the recommendations service for this page. Spec Constraints rule out personalised ranking; calling `chooseAd` is fine in principle (ads are public, non-personalised) but unnecessary for the MVP. Leaving them out keeps the handler minimal.

## Caching

No `Cache-Control` headers added in MVP. Behaviour inherits from the existing default response headers set by `net/http` + the existing middleware stack (`logHandler`, `ensureSessionID`, `otelhttp.NewHandler`).

## Observability

- One `log.Info` line per request: `{event: "search", query: "<normalised>", count: <len(items)>}`. No PII; the query is the shopper's free-text input and is logged for triage of "did anyone search for X?" questions.
- Tracing is inherited from `otelhttp.NewHandler` (`main.go:168`); the gRPC client call to `SearchProducts` is also automatically traced via `otelgrpc`. No new spans added by hand.
- No new metrics. No new dashboards.

## Idempotency and safety

- Safe and idempotent (it's a GET). No side effects.
- No mutation of session state, cart state, or any other persisted store.

## Routing addition

A single line is added to the route table in `main.go` alongside the existing routes:

```go
r.HandleFunc(baseUrl+"/search", svc.searchHandler).Methods(http.MethodGet, http.MethodHead)
```

The `Methods` clause matches the home and product-detail routes (HEAD allowed alongside GET so health-checking infrastructure that issues HEADs behaves consistently).
