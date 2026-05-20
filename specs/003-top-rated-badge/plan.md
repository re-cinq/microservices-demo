# Implementation Plan: Top rated badge on the highest-rated product card

**Branch**: `attendee/daniel-tufvander` (work continues on this branch; spec dir `003-top-rated-badge` is independent of branch name) | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/003-top-rated-badge/spec.md`
**Source story**: [AIP-166](https://odevo.atlassian.net/browse/AIP-166) under epic [AIP-95](https://odevo.atlassian.net/browse/AIP-95)
**Builds on**: [AIP-161](https://odevo.atlassian.net/browse/AIP-161) (shipped) — `Rating float32` is already on `productView` in `handlers.go` and the star widget is already on every home-page card.

## Summary

Mark the home-page product card(s) tied at the maximum rating with a "Top rated" badge. Pure frontend change: compute the max rating in `homeHandler` after the existing product loop, set a new `TopRated bool` flag on each `productView`, render a small badge in `templates/home.html` when the flag is set, style it via a new CSS rule in `static/styles/styles.css`. Zero proto change, zero backend change, zero data change.

## Technical Context

**Language/Version**: Go (existing version pinned by `src/frontend/go.mod`; no change).
**Primary Dependencies**: existing only — `html/template`, `strings`, the in-repo `genproto.Product` type from AIP-161. No new module.
**Storage**: none. The badge is derived at render time from in-memory data.
**Testing**: optional. The "find max + flag" logic is small enough to skip; if extracted into a pure helper it can be unit-tested in a single table-driven test. See research D4.
**Target Platform**: existing cohort deployment — `https://daniel-tufvander.training.gcp.re-cinq.com`. No platform change.
**Project Type**: web-service (frontend Go service). Existing layout under `src/frontend/`.
**Performance Goals**: no measurable regression. The added work per home-page render is one linear pass over a 9-element slice — negligible.
**Constraints**: hard constraints C-001…C-006 from the spec — no new services, no new datastore, Go only, no infra/CI, **no proto change, no backend change**. This is strictly a frontend story.
**Scale/Scope**: same catalogue size (~10 products); no fan-out concerns.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

`.specify/memory/constitution.md` is still the unfilled Spec Kit template — no ratified principles, so no gates to enforce. Initial gate passes vacuously. Post-design gate (below) does the same.

## Project Structure

### Documentation (this feature)

```text
specs/003-top-rated-badge/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0 — decisions & alternatives
├── data-model.md        # Phase 1 — productView extension
├── quickstart.md        # Phase 1 — local verification / implementation loop
├── contracts/
│   └── ui-contract.md   # Phase 1 — badge widget contract
├── checklists/
│   └── requirements.md  # Already created by /speckit.specify
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

No `service-contract.md` this time — no gRPC contract change.

### Source Code (repository root)

This story touches three files in one service. No new directories.

```text
src/frontend/
├── handlers.go               # add `TopRated bool` to productView; compute max rating; flag tied products
├── templates/home.html       # render `<span class="hot-product-card-top-rated">Top rated</span>` when .TopRated
└── static/styles/styles.css  # add `.hot-product-card-top-rated` rule
```

**Structure Decision**: extend the frontend service in place. Three files in one service. Smallest possible vertical slice that satisfies the spec.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations.

## Phase 0 — Outline & Research

See [research.md](research.md). The five spec-level assumptions (A1–A5) are settled there alongside the implementation-level decisions (where to compute max, badge markup choice, accessibility, tie semantics, test strategy).

## Phase 1 — Design & Contracts

Outputs:

- [data-model.md](data-model.md) — `productView` gains a `TopRated bool`. No persistence, no entity change.
- [contracts/ui-contract.md](contracts/ui-contract.md) — badge widget contract: placement inside `.hot-product-card`, copy ("Top rated"), CSS class, ARIA treatment, what does *not* change.
- [quickstart.md](quickstart.md) — implementation steps + local verification path (no skaffold; defer visual check to cohort URL).

## Phase 1 — Re-evaluate Constitution Check (post-design)

Constitution still unfilled; gates still vacuous. No change.
