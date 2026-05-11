# Implementation Plan: Product Search

**Branch**: `001-product-search` | **Date**: 2026-05-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-product-search/spec.md`

## Summary

Add a customer-facing product search to Online Boutique: a header-resident search input on every page, a `GET /search?q=…` results page in the Go frontend, backed by the **already-existing** `ProductCatalogService.SearchProducts` gRPC RPC. The bulk of the work is frontend glue (one handler, one rpc-client wrapper, one template, one route, one header partial). One server-side change is required: tighten `productcatalogservice/product_catalog.go` so `SearchProducts` matches **product name only** (case-insensitive substring), to honour spec FR-002. No proto changes. No new services, datastores, dependencies, env vars, manifests, Helm charts, or CI config.

## Technical Context

**Language/Version**: Go (matching the existing `frontend` and `productcatalogservice` modules — versions pinned in `src/frontend/go.mod` and `src/productcatalogservice/go.mod`; the change must compile against those toolchains unchanged).
**Primary Dependencies**: `google.golang.org/grpc`, `github.com/gorilla/mux`, `html/template`, `github.com/sirupsen/logrus`, `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` — all already imported. **No new modules.**
**Storage**: None. Catalog source-of-truth remains `productcatalogservice/products.json`, loaded into memory via the existing `parseCatalog()`.
**Testing**: `go test ./...` from each service directory. The `productcatalogservice` module already has a `_test.go`; this is the natural home for the new name-only assertion.
**Target Platform**: Same as today — the existing Kubernetes deployment under `kubernetes-manifests/`. **No deployment-config changes required.**
**Project Type**: Multi-service monorepo (existing). This plan touches two services: `src/frontend/` and `src/productcatalogservice/`.
**Performance Goals**: Spec SC-003 — median search→render under 500 ms. The current catalog has ~10 products; in-memory substring matching is sub-millisecond. The dominant latency will be the existing frontend↔catalog gRPC round-trip plus currency conversion (one `Convert` RPC per result, matching home-page behaviour).
**Constraints**: Spec Constraints section is binding — no new services, no new datastores, no Elasticsearch/Solr/vector DBs, no infra config edits, no build-pipeline edits. Spec SC-008 codifies this as a `git diff` invariant on the merged PR.
**Scale/Scope**: Demo app, ~10 products today, single-digit dev-cluster traffic. No horizontal-scale concerns introduced.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` is an unfilled template (placeholder principles only). With no project-specific principles to gate against, the check is effectively informational. I'll proxy with the spec's own Constraints section, treated as the de-facto constitution for this slice:

| Constraint (proxy gate) | Status | Notes |
|---|---|---|
| No new services (11→11) | **PASS** | Plan extends `productcatalogservice` and `frontend` only. |
| No new datastore / search engine | **PASS** | Catalog stays as the existing static `products.json`; in-memory filter only. |
| Match service language | **PASS** | Both edited services are Go; existing Go patterns reused (handler signature, rpc-client wrapper, template glob). |
| No new infra / Helm / manifest / env var | **PASS** | No edits planned under `kubernetes-manifests/`, `kustomize/`, `helm-chart/`, `istio-manifests/`, `terraform/`. No new env vars in any service Deployment. |
| Stay in branch/repo, no build-pipeline edits | **PASS** | No edits planned under `.github/`, `cloudbuild.yaml`, `skaffold.yaml`, or any service `Dockerfile`. |
| No proto-contract changes | **PASS** | The `SearchProducts` RPC already exists in `protos/demo.proto`; this plan does not modify the proto, so the generated `genproto/` artefacts also remain untouched. |

**Outcome:** no violations. Complexity Tracking section below is empty.

## Project Structure

### Documentation (this feature)

```text
specs/001-product-search/
├── plan.md              # This file (/speckit-plan command output)
├── spec.md              # Already exists (/speckit-specify output)
├── research.md          # Phase 0 output — decisions + rationale
├── data-model.md        # Phase 1 output — entities (Product, SearchQuery, SearchResult)
├── quickstart.md        # Phase 1 output — manual + automated verification recipe
├── contracts/           # Phase 1 output — RPC + HTTP contract documentation
│   ├── searchproducts.rpc.md
│   └── search-route.http.md
├── checklists/
│   └── requirements.md  # Already exists (/speckit-specify output)
└── tasks.md             # Phase 2 output (/speckit-tasks command — NOT created here)
```

### Source Code (repository root)

This feature touches files in two existing services. No new directories are created at the source-tree level.

```text
src/
├── frontend/
│   ├── main.go                    # +1 route registration (/search)
│   ├── handlers.go                # +1 handler: searchHandler
│   ├── rpc.go                     # +1 client wrapper: searchProducts
│   ├── validator/
│   │   └── validator.go           # +1 payload struct: SearchPayload (optional; query-string validation may live in the handler)
│   └── templates/
│       ├── header.html            # +1 form: search input in the global nav
│       └── search-results.html    # NEW file — mirrors home.html's product-grid block
└── productcatalogservice/
    ├── product_catalog.go         # ~1 line: drop description-substring branch from SearchProducts
    └── product_catalog_test.go    # +1 test asserting description-only matches return zero
```

**Structure Decision:** monorepo, in-place edits. The feature is a thin frontend layer on top of an already-existing backend RPC, plus a one-line semantic narrowing of that RPC. No restructuring required.

## Complexity Tracking

> *No constitution violations — section intentionally empty.*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| *(none)* | — | — |
