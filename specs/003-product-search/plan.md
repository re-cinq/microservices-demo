# Implementation Plan: Product Search

**Branch**: `003-product-search` | **Date**: 2026-05-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-product-search/spec.md`

## Summary

Add a single search input above the existing product grid on the Online Boutique home page (`/`). The grid's product cards are already rendered server-side by the Go frontend service from the in-memory product catalogue. A small vanilla-JS module attached to that page filters the **already-rendered DOM cards** by case-insensitive substring match on the card's product name, updating visibility on every keystroke with no network call. A polite ARIA live region announces the current visible-count or the "no products match" empty state.

No backend service changes. No new services, datastores, infra, env vars, or pipeline edits. Frontend changes are limited to: one template edit (`home.html`), one new static JS file, and small additions to the existing stylesheet for the input and empty state.

## Technical Context

**Language/Version**: Go 1.25 (frontend service, only template edits; no Go code changes required). Vanilla JavaScript (ES2017+, no framework) for the client-side filter.
**Primary Dependencies**: Existing only — Go `html/template`, Bootstrap classes already in `home.html`. **No new dependencies**.
**Storage**: N/A. Filtering operates on the already-rendered DOM. Product data continues to come from `productcatalogservice/products.json` via the existing gRPC `ListProducts` call (unchanged).
**Testing**: Go `testing` (no new Go code, so no new Go tests). Manual / browser-driven verification per `quickstart.md`. Optionally a tiny JS unit test could be added for the substring/trim logic, but no JS test harness exists in the frontend today — out of scope unless requested.
**Target Platform**: Browser (desktop + mobile web). Same evergreen-browser support as Online Boutique today.
**Project Type**: Web application (Go microservices); this feature only touches the `frontend` service's templates and static assets.
**Performance Goals**: Grid update perceived within **100 ms of each keystroke** (SC-002). DOM toggling for ~10–50 cards is sub-millisecond; the budget is comfortably met by direct `element.hidden = …` toggling.
**Constraints**: 0 additional network requests per keystroke (SC-003); substring match on `name` only (FR-004, C-009); no search state persistence across navigations (FR-010); no new env vars, services, datastores, manifests, charts, or pipeline changes (C-001…C-006, C-008).
**Scale/Scope**: Current catalogue is ~10 products. Implementation is O(n) per keystroke over the in-DOM card list and would remain comfortably under the 100 ms budget at 10× the current catalogue size.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` is the default Spec Kit template with unfilled placeholders — no project-specific principles are ratified for this repo. There are therefore no enforceable constitutional gates to check.

**Gate status**: PASS (vacuous — no principles ratified).

**Self-imposed gates** (derived from the spec's explicit constraints, treated as hard rules for this plan):

- G1 — No new services: PASS (frontend edit only).
- G2 — No new datastores or external indexes: PASS (DOM-level filter).
- G3 — No new env vars, manifests, Helm/Kustomize/Terraform/Skaffold edits, no CI changes: PASS (none planned).
- G4 — No new gRPC method, protobuf, or backend endpoint: PASS (no backend change).
- G5 — Match only product `name`, not description/category/ID: PASS (filter logic operates on the rendered name).
- G6 — Zero additional network requests on keystroke: PASS (purely DOM operations).

## Project Structure

### Documentation (this feature)

```text
specs/003-product-search/
├── plan.md              # This file (/speckit-plan output)
├── spec.md              # /speckit-specify output
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── ui-contract.md   # UI/DOM contract (the only external interface this feature exposes)
└── checklists/
    └── requirements.md  # Spec quality checklist
```

### Source Code (repository root)

Only the `frontend` service is touched. Concrete files:

```text
src/frontend/
├── templates/
│   └── home.html              # EDIT: add search input + ARIA live region above .hot-products-row;
│                              #       add data-name attribute to each .hot-product-card for filtering;
│                              #       add <script src=".../static/js/product-search.js"> at end of page.
├── static/
│   ├── js/                    # NEW DIRECTORY
│   │   └── product-search.js  # NEW: client-side filter + live-region announcer
│   └── styles/
│       └── styles.css         # EDIT (small): styles for .product-search input + .product-search-empty state
└── main.go                    # NO CHANGE (existing FileServer at /static already serves new JS file)
```

Files **not** touched:

- `src/productcatalogservice/*` — no backend change (FR-007, C-007).
- `src/frontend/handlers.go`, `rpc.go`, `main.go` — no Go code changes.
- `helm-chart/`, `kubernetes-manifests/`, `kustomize/`, `terraform/`, `skaffold.yaml`, `cloudbuild.yaml`, `.github/workflows/` — no infra or pipeline changes (G3).

**Structure Decision**: Single-service edit confined to `src/frontend/templates/home.html`, a new `src/frontend/static/js/product-search.js`, and a small CSS appendix in `src/frontend/static/styles/styles.css`. This matches the "no new services" constraint and aligns with how Online Boutique already serves static assets (the frontend's `main.go` mounts `/static` from `./static`).

## Phase 0 Output

See [research.md](./research.md). All Technical Context items resolved with no remaining `NEEDS CLARIFICATION`.

## Phase 1 Outputs

- [data-model.md](./data-model.md) — describes the client-side `Query` state and the DOM contract that backs the filter.
- [contracts/ui-contract.md](./contracts/ui-contract.md) — the user-facing contract (selectors, attributes, ARIA roles) that the JS module relies on. This is the analogue of an API contract for a UI-only feature.
- [quickstart.md](./quickstart.md) — how to run the frontend locally with the change and verify each acceptance scenario.

### Post-design Constitution re-check

No constitutional gates exist; the self-imposed gates (G1–G6) remain satisfied after design.

## Complexity Tracking

No violations to track. The implementation strictly reuses existing capability:

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| _(none)_  | _(n/a)_    | _(n/a)_                              |
