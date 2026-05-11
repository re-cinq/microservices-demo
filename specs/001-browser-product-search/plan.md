# Implementation Plan: Browser-Side Product Search

**Branch**: `attendee/daniel-tufvander` (spec dir: `001-browser-product-search`) | **Date**: 2026-05-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-browser-product-search/spec.md`

## Summary

Add a search box above the product grid on the Online Boutique home page (the product list page). When the shopper types, the page hides products whose `name` does not contain the query (case-insensitive substring match, name field only). Filtering runs entirely in the browser over the already-rendered product cards — no extra server requests per keystroke, no new datastore, no protobuf change.

Technical approach: extend the existing Go `frontend` service's `home.html` template to render a `<input type="search">` and add data attributes to each product card carrying the lowercase product name. Ship one new static asset, `src/frontend/static/js/product-search.js`, which attaches an `input` event listener and toggles a `hidden` attribute on cards whose data-name does not contain the trimmed, lowercased query. Add minimal CSS to `styles.css`. No changes to `productcatalogservice`, the gRPC contract, `products.json`, or any deployment manifest.

## Technical Context

**Language/Version**: Go 1.22 (the version `src/frontend/go.mod` resolves to today); browser-side JavaScript ES2017+ (vanilla — no framework). No language change introduced by this feature.
**Primary Dependencies**: Existing `frontend` service uses Go `html/template` for server-side rendering. Existing `static/styles/styles.css` for CSS. No new runtime libraries.
**Storage**: None. Catalogue is the existing in-memory list loaded by `productcatalogservice` from `src/productcatalogservice/products.json`. No new storage, no client-side persistence (filter is transient and resets on reload — FR-010, SC-008).
**Testing**:
- Unit-level JS test for the substring/case-insensitive filter function: a tiny Node-runnable spec colocated with the JS file (no new test runner — runs with `node --test` if needed, or stays as a manual snippet executable in the browser console).
- Manual acceptance via the quickstart (Phase 1 deliverable). Each of the 8 G/W/T scenarios from `spec.md` maps to a step in the quickstart.
- No new CI test infra is introduced (constraint: "DO NOT touch the build pipeline").
**Target Platform**: Browser (modern, JS-enabled). Server-side: existing `frontend` container image, deployed by existing CI on `attendee/<your-name>`.
**Project Type**: Web application — server-rendered HTML with a thin sprinkle of client-side JS. Mapped to the repo's existing layout (`src/frontend/`, no new service).
**Performance Goals**: SC-002 / SC-003 — filter result visible within 200 ms of the last keystroke. With ~10 products in the current catalogue this is trivial; an upper bound of a few thousand DOM nodes still completes well under 200 ms.
**Constraints**:
- Zero new network requests per keystroke (SC-006).
- Zero new accessibility violations on the product list page (SC-007).
- Feature does not alter any page other than the product list (FR-012).
- Filter operates only over what the page already rendered (FR-002, FR-007).
**Scale/Scope**: Current catalogue ~10 products; feature is sized for "small catalogue rendered on one page", not for tens-of-thousands.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The repository's `.specify/memory/constitution.md` is the unfilled SpecKit placeholder — no project-specific principles have been ratified. There is therefore no constitutional gate to check against in the strict sense.

In the absence of a ratified constitution, the **explicit "Out of Scope / Do NOT" section of [spec.md](./spec.md) is treated as the binding gate for this feature.** Each item below is a project-level constraint this plan must respect:

| Spec constraint | This plan satisfies it by |
|---|---|
| No new services | Changes touch only `src/frontend/` (existing service). |
| No Elasticsearch / Solr / vector DB / new datastore | Filtering is pure client-side over already-rendered DOM. |
| Use the existing in-memory catalogue from `products.json` | Catalogue rendering is unchanged; the new code reads `name` from already-rendered cards. |
| Match the service's language; reuse gRPC/protobuf patterns if extending | Frontend stays Go + HTML templates. Productcatalogservice is **not** touched. No protobuf change. |
| No new infrastructure config, Helm charts, manifests, env vars | Plan adds only source files inside `src/frontend/`. No `kubernetes-manifests/`, `helm-chart/`, `kustomize/`, `terraform/`, or `cloudbuild.yaml` change. |
| Stay in this branch and repo; do not touch build pipeline | All work lands on `attendee/daniel-tufvander`. No CI config edits. |

**Gate status (initial)**: PASS. No violations require justification.

**Gate status (post-Phase-1 re-check, see end of this document)**: PASS — design phase did not introduce any deviation from the constraints above.

## Project Structure

### Documentation (this feature)

```text
specs/001-browser-product-search/
├── plan.md              # This file (/speckit-plan output)
├── spec.md              # Feature spec (already created)
├── research.md          # Phase 0 output — decisions and rationale
├── data-model.md        # Phase 1 output — Product (existing) + SearchQuery (transient)
├── quickstart.md        # Phase 1 output — manual verification walk-through
├── contracts/
│   └── ui-contract.md   # Phase 1 output — UI element contract for the search box and product cards
└── checklists/
    └── requirements.md  # Spec-quality checklist (already created)
```

### Source Code (repository root)

```text
src/frontend/                                   # existing Go service — only service touched
├── templates/
│   └── home.html                               # MODIFIED: add <input type="search"> above .hot-products-row;
│                                                #           add data-name="<lowercased name>" to each product card
├── static/
│   ├── styles/
│   │   └── styles.css                          # MODIFIED: small additions for .product-search and .no-results
│   └── js/                                     # NEW DIRECTORY
│       └── product-search.js                   # NEW: vanilla JS — input listener, substring filter, no-results message
└── handlers.go                                 # UNCHANGED (homeHandler already passes products to the template)

src/productcatalogservice/                      # NOT TOUCHED
└── products.json                               # NOT TOUCHED — read-only data source, behaviour unchanged

kubernetes-manifests/, helm-chart/,             # NOT TOUCHED — no deployment changes
kustomize/, terraform/, cloudbuild.yaml         #
```

**Structure Decision**: This is a server-rendered web application; the chosen structure is the repo's existing `src/<service>/` layout. No new service or top-level directory is created. The only new directory is `src/frontend/static/js/` (the frontend has no JS today — the bot widget aside, there are no `*.js` files under `static/`). Confining all changes to one existing service keeps the blast radius minimal and the "no infrastructure change" constraint trivially satisfied.

## Complexity Tracking

No constitution violations. Table omitted.

---

## Post-Design Constitution Re-check

After producing `research.md`, `data-model.md`, `contracts/ui-contract.md`, and `quickstart.md`:

- **No new service introduced.** ✅
- **No new datastore.** ✅ Filter operates on the DOM that `home.html` already produces.
- **No protobuf / gRPC change.** ✅ `productcatalogservice` is untouched.
- **No infra / Helm / manifest / env var change.** ✅ File list in "Source Code" above is exhaustive.
- **Language match preserved.** ✅ Server-side Go remains Go; the only new language surface is vanilla ES2017 JS, which is the standard browser surface for a server-rendered Go web app and is not a new "service language".
- **Build pipeline untouched.** ✅ No edits to `cloudbuild.yaml`, `skaffold.yaml`, or any `.github/workflows/` file.

**Gate status (final)**: PASS. Ready for `/speckit-tasks`.
