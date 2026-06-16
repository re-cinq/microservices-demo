# Feature Specification: Recently Viewed Strip

**Jira Epic**: [AIP-97](https://odevo.atlassian.net/browse/AIP-97) — Recently viewed
**Jira Story**: [AIP-175](https://odevo.atlassian.net/browse/AIP-175) — Display recently viewed strip

**Branch**: `attendee/james-colebeck`

**Created**: 2026-06-16

**Status**: Draft

---

## Problem Statement

Shoppers lose track of products they looked at during their session and have to hunt for them again. We want to surface a "recently viewed" strip so they can get back to a product without searching for it again. The strip is built entirely in memory using the existing product catalogue — no new service or datastore is introduced.

---

## User Story

As a shopper browsing Online Boutique, I want to see a strip of products I have already looked at during my current visit, so I can return to something I was interested in without searching for it again.

---

## Acceptance Criteria

1. **Given** I have viewed at least one product this session, **when** I am on the page where the strip is shown, **then** I see a "recently viewed" strip containing those products.
2. **Given** I have not viewed any products this session, **when** I am on the page where the strip would appear, **then** no strip is rendered at all.
3. **Given** the strip is shown, **when** I look at its content, **then** every item displayed corresponds to a product currently present in the catalogue.
4. **Given** a product I viewed this session is removed from the catalogue, **when** the strip renders, **then** that product does not appear in it.
5. **Given** all products I viewed have been removed from the catalogue, **when** the strip would render, **then** the strip is hidden entirely.

---

## Requirements

### Functional Requirements

- **FR-001**: The page MUST display a "recently viewed" strip when the user has viewed at least one product in the current session.
- **FR-002**: The strip MUST be absent (not rendered) when no products have been viewed yet.
- **FR-003**: Each item in the strip MUST correspond to a product currently present in the catalogue — stale/removed products MUST be filtered out before rendering.
- **FR-004**: If filtering removes all items, the strip MUST be hidden entirely.
- **FR-005**: Session state MUST be stored in-memory (e.g., cookie) with no new datastore introduced.

### Out of Scope

- Cross-device sync
- Persistence beyond the browser session

---

## Hard Constraints (from AIP-97)

- **C-001**: Use only the services that already exist in the repo. Do NOT add new services.
- **C-002**: Do NOT introduce any new datastore (no database, no cache, no search engine). Work in memory over the existing product catalogue in `productcatalogservice/products.json`.
- **C-003**: Match the language and patterns of the service being changed. The frontend and product catalogue service are written in Go.
- **C-004**: Do NOT change infrastructure, deployment manifests, or CI configuration. The change must ship through the existing pipeline unmodified.

---

## Assumptions

- "Session" means the browser session — a cookie scoped to the browser session (no `Max-Age` / `Expires`) is an acceptable mechanism.
- The strip will be shown on the product detail page and/or the home page (to be confirmed in planning).
- The current product is NOT shown in its own recently-viewed strip.
- A reasonable cap on recently-viewed items (e.g., 5) is acceptable.
