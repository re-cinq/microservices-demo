# Implementation Plan: Product Search

**Branch**: `001-product-search` | **Date**: 2026-05-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/001-product-search/spec.md`

## Summary

Add a client-side search input directly above the existing product grid on the Online Boutique product list page (`/`). The page already server-renders the full catalogue (~9 products); a small, page-scoped JavaScript module attaches an `input` event listener to the search box and toggles the `display` of existing product cards based on a case-insensitive substring match against each card's `data-product-name` attribute. No new services, no new datastore, no new RPC, no new infra config. Purely a presentation-layer enhancement inside `src/frontend`.

## Technical Context

**Language/Version**: Go 1.25 (frontend service, per `src/frontend/go.mod`); HTML5 + CSS3 + ES2015+ vanilla JavaScript (browser-side, no framework)
**Primary Dependencies**: existing Go `html/template` standard library (no new Go deps); existing Bootstrap 4.1.1 CSS via CDN (no new CSS deps); no JS dependencies introduced
**Storage**: N/A — feature is read-only against the already-loaded product list; no persistence in any tier (C-015)
**Testing**: Go `httptest` + standard `testing` package for handler/template assertions on the rendered markup contract; manual browser verification (driven by the spec's concrete acceptance scenarios in `quickstart.md`) for client-side behaviour. No new test framework added (would violate spirit of C-006 / C-008).
**Target Platform**: any modern browser that meets Online Boutique's existing baseline (ES2015+, same as the rest of the site)
**Project Type**: web — server-rendered Go service serving HTML + static assets. Changes confined to `src/frontend`.
**Performance Goals**: SC-003 — keystroke → DOM update under 100 ms. Trivially met by direct DOM `display` toggling over 9 cards; no debouncing, no indexing, no Web Worker, no virtualisation required (A-007).
**Constraints**: full Constraints section of the spec (C-001..C-015). Most load-bearing for the plan: C-010 (no layout displacement), C-013 (no per-keystroke catalogue fetches), C-014 (no service-contract change), C-015 (no query persistence).
**Scale/Scope**: single page (`/`), ~9 products today, single small JS file (~50 lines), single small CSS file (~30 lines), one Go test file. Entirely confined to `src/frontend`.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` is the unfilled SpecKit template — no project principles have been ratified for Online Boutique. The constitution gate is therefore vacuous: there are no principles to check this plan against. **No violations identified.**

In place of constitution gates, this plan is checked against the spec's own Constraints section (C-001..C-015). Compliance verified inline in the Project Structure / Phase 0 sections below.

**Re-check after Phase 1**: still vacuous; Phase 1 design (`data-model.md`, `research.md`, `quickstart.md`) introduces no new gates.

## Project Structure

### Documentation (this feature)

```text
specs/001-product-search/
├── spec.md                          # already written by /speckit-specify
├── plan.md                          # this file
├── research.md                      # Phase 0 — decisions + alternatives
├── data-model.md                    # Phase 1 — DOM ↔ JS contract; no new server entities
├── quickstart.md                    # Phase 1 — how to build, run, verify
└── checklists/
    └── requirements.md              # spec-quality checklist, written by /speckit-specify
```

No `contracts/` directory: this feature introduces no new service-to-service contract (C-014). The only new contract is the DOM ↔ JS one, captured in `data-model.md`.

### Source Code (repository root)

```text
src/frontend/
├── handlers.go                      # READ-ONLY — homeHandler at /, already passes full product list to template
├── main.go                          # READ-ONLY
├── go.mod                           # READ-ONLY — no new Go deps
├── templates/
│   ├── header.html                  # MODIFY — one new <link> for search.css
│   └── home.html                    # MODIFY — search input + empty-state above grid; data-product-name on each card; <script src> for search.js at end of body
├── static/
│   ├── styles/
│   │   └── search.css               # NEW (~30 lines) — search wrapper, clear button, empty-state styling
│   └── js/                          # NEW directory (first JS in the repo)
│       └── search.js                # NEW (~50 lines vanilla JS) — filter + empty-state + clear + Esc
└── home_test.go                     # NEW — Go test asserting rendered markup contract
```

Read-only / untouched everywhere else:
- All other services under `src/` (C-001 / C-009).
- `kubernetes-manifests/`, `helm-chart/`, `kustomize/`, `terraform/`, `istio-manifests/`, `release/` (C-006).
- `cloudbuild.yaml`, `skaffold.yaml`, every `Dockerfile` (C-008).
- `protos/` and `genproto.sh` (C-005 / C-014).
- `src/productcatalogservice/products.json` (C-003 — read but not modified).

**Structure Decision**: single-service change confined to `src/frontend`. Static assets follow the existing separation pattern (CSS is already split into per-area files under `static/styles/`). `static/js/` is a brand-new sibling directory because the project currently contains no JavaScript files — a deliberate, minimal-footprint addition that mirrors the CSS pattern rather than inlining a `<script>` block in the template. The header template gains one new `<link>` (CSS, site-wide is fine — tiny + page only renders when search is shown, but loading it everywhere costs us nothing meaningful); the home template gains one new `<script src=… defer>` so the JS only loads on the product list page.

## Complexity Tracking

> Constitution Check passed (vacuous gate). No spec-constraint violations.

Nothing to justify. The plan introduces:
- one new directory (`src/frontend/static/js/`)
- two new files in `src/frontend/static/` (search.css, search.js)
- one new test file (`src/frontend/home_test.go`)
- two modified template files (`templates/header.html`, `templates/home.html`)

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| *(none)*  | —          | —                                    |
