# Research: Top rated badge

**Feature**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md)

Same Decision / Rationale / Alternatives shape as the AIP-161 research. Smaller surface, so fewer entries.

## D1. Where to compute "top rated"

- **Decision**: in `homeHandler`, after the existing `ps := make([]productView, len(products))` loop builds the slice, walk `ps` once to find `maxRating := max over ps.Rating`, then walk it a second time to set `ps[i].TopRated = (ps[i].Rating == maxRating && maxRating > 0)`.
- **Rationale**: keeps the logic in one place, one function, no template-side computation. Template stays declarative (`{{ if .TopRated }}`). Two linear passes over a 9-element slice — trivially fast.
- **Alternatives considered**:
  - Compute in the template via a `{{ if eq .Rating $maxRating }}` pattern: would require a "max" template func; uglier and reverses the usual separation (Go computes, template renders).
  - Compute in the catalogue service: violates C-006 (no backend change) and means every consumer of `ListProducts` carries a derived field they didn't ask for.
  - Compute via JavaScript on the client: introduces a new technology, fights C-003, and re-shuffles after first paint.

## D2. Tie semantics

- **Decision**: every product tied at `maxRating` gets the badge. Spec FR-002 and AC-2 both make this explicit.
- **Rationale**: "tied at the top" is honest. Arbitrary tie-breaking (alphabetical, by id, by JSON order) is a hidden judgement that a reviewer would rightly question.
- **Alternatives considered**:
  - Pick first product in JSON order: hidden ordering signal, surprises future editors.
  - Pick randomly: violates SC-003 (stability across reloads).

## D3. Zero-rating fallback

- **Decision**: if `maxRating == 0`, no card carries the badge (i.e. the comparison `maxRating > 0` gates the flag).
- **Rationale**: matches FR-003 and AC-3. Avoids the degenerate case where a 0.0 product would be "the top" of a catalogue of 0.0s.
- **Alternatives considered**:
  - Always badge `max(ratings)` even when 0: misleads the shopper into thinking an unrated product is recommended.

## D4. Test strategy

- **Decision**: no new test. The "compute max + flag tied" logic is 5 lines inline in `homeHandler`. Extracting it into a pure helper just to test it would add more code than it would catch.
- **Rationale**: the existing pattern in this codebase is to not test handler-internal view-model wiring; AIP-161 added one focused catalogue test because the data path crossed a service boundary. Here the change is local to one function. Manual verification on the cohort URL (per quickstart) covers AS-1/2/3.
- **Alternatives considered**:
  - Extract `flagTopRated(views []productView)` helper + unit test: defensible, but ~15 lines of test for ~5 lines of behaviour. Not worth it for cohort scope.
  - Add an htmltest snapshot of the rendered home page: would require adding an htmltest dependency — fights C-003 and adds a much heavier test surface for a trivial change.

## D5. Badge markup choice

- **Decision**: a single `<span class="hot-product-card-top-rated">Top rated</span>` rendered inside `.hot-product-card`, after the price and rating elements but before any future card additions. Block-positioned via CSS; absolutely positioned to sit at the top-right corner of the card so it does not affect existing element flow.
- **Rationale**: cheapest possible markup. No new template partial. Absolute positioning means the badge cannot push the price or stars around — satisfies FR-005 (no existing element altered/hidden).
- **Alternatives considered**:
  - Inline before the name: pushes name down, changes layout, fails FR-005.
  - Overlay on the product image: harder to read, requires image-tint logic for contrast.
  - Standalone partial template (`{{ template "top-rated-badge" . }}`): premature abstraction for a one-line element used in one place.

## D6. Accessibility

- **Decision**: badge is plain text inside a span. No `aria-*` needed — screen readers will read "Top rated" naturally inline. No `role="badge"` or `<mark>` element because the visual style does not need to be announced semantically.
- **Rationale**: simplest accessible thing that works. The star widget already gives the rating value via `aria-label`; the badge adds a textual cue that any screen reader reads as part of the card.
- **Alternatives considered**:
  - Hide visually with `aria-hidden` and rely solely on the star `aria-label`: defeats the purpose of the badge for sighted users; worse for screen readers.
  - Use `<mark>` for semantic emphasis: changes default browser styling unpredictably across UAs.

## D7. CSS technique for the badge

- **Decision**: position the badge at the top-right of `.hot-product-card` with `position: absolute; top: 8px; right: 8px;`. The parent `.hot-product-card` already has positioning context (because its child overlay uses `position: absolute`). Style with the existing accent colour `#f5a623` (already used for filled stars) as the background; white text; small rounded pill shape.
- **Rationale**: re-uses the rating colour to visually anchor the badge to the rating system. The parent already has `position` context, so no parent change needed.
- **Alternatives considered**:
  - Inline-block at the bottom of the card: pushes content, fails FR-005.
  - Use the cymbal brand colour instead of accent: would split the visual language (rating colour ≠ badge colour) and lose the implicit "this is part of the rating system" cue.

---

## Resolved spec assumptions

| Spec Assumption | Resolution decision |
|---|---|
| A1 — strict numeric max | **D1, D3** — `max over ps.Rating`, gated by `> 0` |
| A2 — home page only | confirmed; no other surface touched |
| A3 — copy "Top rated" | confirmed; literal English string, no i18n |
| A4 — visual only, not a link | confirmed; the badge is a `<span>`, not an `<a>` |
| A5 — placement inside `.hot-product-card`, no layout impact | **D5, D7** — absolutely positioned top-right |

No `[NEEDS CLARIFICATION]` markers remain.
