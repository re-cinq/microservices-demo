# Research: Star ratings on the product list

**Feature**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md)

This file resolves every spec-level open item (Assumptions A1–A6) and the technical unknowns surfaced when mapping the change against the existing codebase. Each entry uses the *Decision / Rationale / Alternatives considered* shape.

## D1. Rating scale and on-the-wire representation

- **Decision**: rating is a `float` (proto type `float`, Go `float32`) bounded to `[0.0, 5.0]` with half-star granularity in steps of `0.5`.
- **Rationale**: matches the conventional 0–5 star scale named in the parent epic. `float` is the cheapest proto scalar for half-star precision and is already the established pattern for non-currency numerics. Half-star granularity gives 11 distinct values (0, 0.5, … 5.0), enough nuance for social proof without inviting a "5.000 vs 4.999" UX problem.
- **Alternatives considered**:
  - `int32` 0–10 (representing half-stars × 2): saves four bytes per product but forces every caller to know the scaling convention. Rejected — clarity over bytes for a 10-product catalogue.
  - `int32` 0–5 (whole stars only): loses half-star nuance with no upside.
  - `double`: precision overkill; `float` is sufficient and matches the typical Go pattern for ratings.

## D2. Where the rating value lives at rest

- **Decision**: the rating is a new JSON field `rating` (lowercase, matching proto's `json_name` convention) on each entry in `src/productcatalogservice/products.json`. Loaded into memory by the existing `loadCatalog` path through `jsonpb` and held on `productCatalog.catalog`.
- **Rationale**: the only data location compatible with constraint C-002 ("no new datastore… work in memory over the existing product catalogue"). The existing `catalog_loader.go` already deserialises `products.json` into the `pb.ListProductsResponse` struct via `jsonpb`, so adding a proto field automatically wires up the JSON deserialisation without writing custom code.
- **Alternatives considered**:
  - In-process map keyed by product ID: rejected — duplicates the source of truth that already lives in `products.json` and would drift.
  - Computed at startup (random per-process): rejected — violates SC-003 (ratings must be stable across page loads) for any deployment with more than one frontend pod.
  - External JSON file or env var: rejected — moves data outside the file the epic explicitly names as the in-memory source.

## D3. Submission flow (read-only vs writable)

- **Decision**: read-only. No submit, no edit, no delete in this story. Ratings change only when an engineer edits `products.json` and redeploys.
- **Rationale**: the epic puts moderation explicitly out of scope; user-submitted ratings without moderation is not a defensible product. A write path would also require either a new datastore (violates C-002) or mutation of `products.json` at runtime (violates the "in memory over the existing product catalogue" wording, and would not persist across pod restarts anyway).
- **Alternatives considered**:
  - Allow submission with no persistence (ephemeral): rejected — pointless feature, breaks SC-003.
  - Allow submission with persistence: rejected — needs a new datastore.

## D4. Aggregate vs single value

- **Decision**: a single rating value per product, treated semantically as "the rating an editor has set". Not modelled as an aggregate over reviewers.
- **Rationale**: A3 (read-only) and the absence of any reviewer entity in the parent epic. There is no second value (a count, a histogram, a list of reviewers) to aggregate over. Modelling a hidden aggregate would invent state that the data does not carry.

## D5. Display content on the card

- **Decision**: stars only. No review count, no numeric label ("4.3"), no tooltip text. Whole stars filled for each whole point, a half star for `.5` increments, empty stars for the remainder up to 5.
- **Rationale**: the epic talks about "star ratings visible". Review count and numeric label were Open Questions in the parent PRD and were not added there; they should not appear here. Keeping the widget visual-only keeps the card layout stable (FR-005).
- **Alternatives considered**:
  - Stars plus a "(N reviews)" caption: rejected — there is no review-count data in the catalogue.
  - Stars plus a numeric "4.5/5" label: rejected — accessibility argument is real but better solved with an `aria-label` than visible text; visible numeric duplicates the visual representation.

## D6. Accessibility

- **Decision**: each rendered rating element carries an `aria-label` of the form `"Rated 4.5 out of 5 stars"`. The visual stars are decorative; the label is the canonical signal for screen readers.
- **Rationale**: avoids tying screen-reader UX to per-star image semantics and avoids any reliance on emoji or font glyphs. Costs nothing.
- **Alternatives considered**: no label (rejected — fails basic a11y); per-star `aria-label` (rejected — noisy).

## D7. Missing-rating fallback

- **Decision**: any product whose JSON entry has `"rating"` absent or `0.0` renders **no** rating element on its card; the rest of the card layout is unchanged.
- **Rationale**: A4 says all products are seeded with a rating, so this branch is defensive only. Rendering "0 stars" would misrepresent the absence of data as a low rating. Suppressing the widget keeps the card visually identical to today for any un-seeded product.
- **Alternatives considered**:
  - Always render 5 empty stars: rejected — visually claims a 0-rating where there is no data.
  - Render a "No rating yet" caption: rejected — adds new copy with no upstream decision.

## D8. Star rendering technique

- **Decision**: Unicode glyphs (`★` for filled, `☆` for empty, `⯨` or a half-star CSS technique for half). Wrapped in a `<span class="product-card-rating">` so the CSS owns sizing, spacing, and colour.
- **Rationale**: zero new static assets, zero new fonts, no new build step. The frontend already serves CSS from `src/frontend/static/styles/`; adding rules to an existing stylesheet costs nothing and stays within C-004 (no infra change). Unicode glyphs render consistently on all modern browsers shipped to cohort users.
- **Alternatives considered**:
  - SVG icons: rejected — needs new static files plus a sprite or per-icon `<img>`, more surface for a tiny gain.
  - Font Awesome / icon font: rejected — adds a dependency the frontend does not currently carry.

## D9. Proto field number and backward compatibility

- **Decision**: add `float rating = 7;` to `message Product` in `protos/demo.proto`. Field number 7 is the next unused. The field is optional (proto3 default semantics) — older consumers see `0.0` and the frontend treats that as "no rating" per D7.
- **Rationale**: adding a new field with a new number is the textbook backwards-compatible proto change. No existing field is renamed, renumbered, or repurposed. Services not yet rebuilt against the new proto (e.g. checkoutservice) continue to work because the field is unrecognised-and-ignored by older readers.
- **Alternatives considered**:
  - Reusing a deprecated field: there are none.
  - A separate `ProductRating` message keyed by product id: rejected — needs a new RPC, which fights C-001 (no new services) and adds a round-trip on the list page for no benefit.

## D10. Regenerating proto stubs

- **Decision**: run `./genproto.sh` in each Go service that consumes the changed proto (`src/productcatalogservice` and `src/frontend`); commit the regenerated `genproto/*.pb.go` files alongside the proto change.
- **Rationale**: this is the existing workflow in this repo — `genproto.sh` lives in each service and is the documented regeneration entry point. Committing regenerated stubs (rather than regenerating at build time) is consistent with the rest of the repo.
- **Alternatives considered**:
  - Regenerate in CI: would require a CI change, which violates C-004.
  - Skip the frontend's regen: rejected — the frontend uses its own `genproto` package; without regeneration the `Rating` field is not accessible.

## D11. Tests

- **Decision**: extend `src/productcatalogservice/product_catalog_test.go` with a focused test that loads `products.json` and asserts that `ListProducts` returns the expected rating for at least one product. No new test file, no new test framework, no frontend test.
- **Rationale**: the catalogue is the only place where behaviour can be deterministically asserted in code; frontend rendering is covered by the quickstart manual check. The repo's existing pattern uses table-driven tests in the same file, which this extension follows.
- **Alternatives considered**:
  - HTTP / template snapshot tests for the frontend: rejected — would need to add a snapshot or htmltest dependency, fights "match the language and patterns of the service you change" (C-003 in the spec wording).
  - Contract test between frontend and catalog: rejected — orthogonal value for the change at hand.

## D12. Verification surface

- **Decision**: the quickstart verifies on the running app (`skaffold dev` locally or the cohort URL after deploy) and on the new unit test. No browser automation.
- **Rationale**: small surface; manual visual verification is the cheapest way to satisfy AS-1 (every card has a rating) and AS-3 (rating stable across reloads).

---

## Resolved spec assumptions

| Spec Assumption | Resolution decision |
|---|---|
| A1 — rating scale | **D1** — float, 0.0–5.0, 0.5 increments |
| A2 — source of value | **D2** — new `rating` JSON field in `products.json` |
| A3 — read-only | **D3** — read-only, set by editors only |
| A4 — all products seeded | confirmed; **D7** keeps the missing branch defensive |
| A5 — display content | **D5** — stars only, no count or numeric label |
| A6 — both surfaces reuse same value | confirmed; proto field is the canonical source for both stories |

No `[NEEDS CLARIFICATION]` markers remain.
