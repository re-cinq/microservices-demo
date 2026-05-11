# Quickstart: Product Search

**Feature**: 003-product-search
**Date**: 2026-05-11

This document describes how to (a) run Online Boutique with the product-search feature applied, and (b) verify each acceptance scenario from `spec.md` by hand in a browser.

## Prerequisites

- The same toolchain Online Boutique already requires (one of):
  - `skaffold dev` against a local Kubernetes cluster, or
  - `docker compose` for local dev (if present), or
  - The training cluster the cohort is deployed to: the `deploy-attendee` workflow rolls out whatever lands on `attendee/anna-simon`. The frontend is reachable at the host noted in your handout.
- A modern browser (Chrome / Edge / Safari / Firefox).
- Browser DevTools available (specifically the **Network** panel).

## Running locally

The feature only touches the `frontend` service. To iterate quickly:

```bash
# From the repo root
cd src/frontend
go run .
```

(or whichever incantation the team already uses for `skaffold dev`/`docker compose`). The new files are served by the existing static handler — no extra build step.

If you are on the cohort cluster, merging this branch back into `attendee/anna-simon` triggers the GitHub Actions `deploy-attendee` workflow, which rebuilds only the `frontend` service image and rolls it out. Then open the host from your handout in a browser.

## Verification checklist (mapping each spec scenario to a concrete action)

For each row below: open the home page (`/`) in a fresh tab and follow the steps. The leftmost column links the row back to the spec.

### Story 1 — Filter the product list as I type (P1)

| # | Spec scenario | Action | Pass criteria |
|---|---|---|---|
| 1.1 | S1#1 | Type `sunglasses` into the search box | Only the "Sunglasses" card is visible; all others are `hidden` |
| 1.2 | S1#2 | Type `mu`, then `g` (extending to `mug`) | Grid narrows on each keystroke; no spinner, no page reload |
| 1.3 | S1#3 (case-insensitive) | Type `Sun` | "Sunglasses" remains visible (uppercase `S` still matches) |

### Story 2 — Clear the search to see everything again (P1)

| # | Spec scenario | Action | Pass criteria |
|---|---|---|---|
| 2.1 | S2#1 | After filtering by `mug`, delete all characters from the search input | Every product is visible again, in the original order |
| 2.2 | S2#2 | Type `zzzz` (no match), then click the native × clear button in the input | Every product is visible again; "no products match" message is gone |

### Story 3 — See a clear "no results" state (P2)

| # | Spec scenario | Action | Pass criteria |
|---|---|---|---|
| 3.1 | S3#1 | Type `zzzz` | Grid is empty; a "No products match your search." message is shown in its place |
| 3.2 | S3#2 | Then clear the search input | The empty-state message disappears; full grid is restored |

### Edge cases

| # | Edge case | Action | Pass criteria |
|---|---|---|---|
| E1 | Whitespace-only query | Type 3 spaces | Full grid is shown; no empty-state message |
| E2 | Leading/trailing whitespace | Type `  mug  ` | Same result as `mug` |
| E3 | Special characters | Type `mug!@#` | Empty-state message is shown |
| E4 | Very long query | Paste 500 characters of `x` | Page does not crash; empty-state message is shown |
| E5 | Rapid typing | Hold a key down for ~1 second | Grid keeps up; no stuck/stale results |
| E6 | Navigate away and back | Type a query → click any product → click browser Back | Search input is empty; full grid is shown (no persistence) |

### Non-functional verification

| # | SC | How to verify | Pass criteria |
|---|---|---|---|
| N1 | SC-002 | Type a query; watch grid update | Grid update perceptible in the same frame as the keystroke (no visible lag) |
| N2 | SC-003 | Open DevTools → Network → filter "All", clear the log, then type a query of any length | **Zero** new requests recorded — the network panel stays empty as you type and clear |
| N3 | SC-004 | Try several queries that are substrings of multiple product names | 100% of matching cards remain visible; 0% of non-matching cards remain visible |
| N4 | SC-005 | Reload the page with an empty input | Page renders byte-equivalently to before this feature — same product set, same order, same `ListProducts` request fired once on load |
| N5 | Accessibility (R5) | Enable VoiceOver / NVDA, type a query | The screen reader announces "N products match" / "No products match your search." after each filter change |

If every row passes, the feature is ready for `/speckit-tasks` → `/speckit-implement` to lock in.
