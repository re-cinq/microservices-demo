# Feature Specification: Browser-Side Product Search

**Feature Branch**: `001-browser-product-search`
**Created**: 2026-05-11
**Status**: Draft
**Input**: User description: "Add a product search feature to Online Boutique. Frontend-only filter: a search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

## Problem *(in the user's words)*

> "I'm browsing Online Boutique and the product list keeps growing. I know the name of the thing I want — I just want to type it into a search box and see only the matching products. I don't want to scroll through everything or wait for the page to reload."

Shoppers landing on the product list page today have no way to narrow the list down. The page renders all products at once. As the catalogue grows, finding a specific item by name becomes slow and frustrating. We want a search box on the product list page that filters the products already on screen, by name, instantly, without a page reload or backend call.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Filter the visible product list by typing a name (Priority: P1)

A shopper is on the product list page. They want to find a specific product whose name they (roughly) know. They type into a search box and the list narrows to only products whose name matches what they typed. They can clear the box to see everything again.

**Why this priority**: This is the entire feature. Without it, there is no product search. It is also the only story required to deliver value end-to-end — once it works, shoppers can find products by name. Everything else (edge cases, accessibility refinements) is a hardening pass on this same flow.

**Independent Test**: Open the product list page in a browser. Type a query in the search box. Observe that the visible products update to match. Clear the box. Observe that all products return. No other story needs to be implemented for this to be useful.

**Acceptance Scenarios**:

1. **Given** the product list page is loaded and shows all products,
   **When** the shopper types `sun` into the search box,
   **Then** only products whose name contains `sun` (case-insensitive) remain visible, and the rest are hidden.

2. **Given** the search box contains `sun` and the list is filtered,
   **When** the shopper clears the search box,
   **Then** every product that was originally on the page is visible again.

3. **Given** the product list page is loaded,
   **When** the shopper types `SUN` (uppercase) into the search box,
   **Then** the result is identical to typing `sun` — matching is case-insensitive.

4. **Given** the search box is empty,
   **When** the shopper types a single character,
   **Then** the list updates within the same interaction (no page reload, no spinner that blocks the page).

5. **Given** the shopper has typed `xyzzy` (no product matches),
   **When** they look at the page,
   **Then** the product grid is empty and a short "no results" message is shown in its place; the search box still contains `xyzzy` so they can edit it.

---

### User Story 2 — Resetting the filter (Priority: P2)

A shopper has filtered the list and now wants to go back to the full catalogue without retyping or refreshing.

**Why this priority**: Story 1 already covers clearing via deleting characters. This adds a faster path (e.g. a clear control or pressing Escape). It is a usability improvement, not a blocker.

**Independent Test**: With a query typed, use the reset affordance. Observe the list returns to the full set and the box is empty.

**Acceptance Scenarios**:

1. **Given** the search box contains `watch` and the list is filtered,
   **When** the shopper activates the reset control (e.g. clicks an "x" inside the box or presses Escape with the box focused),
   **Then** the search box becomes empty and the full list is restored.

---

### User Story 3 — Filter survives in-page interactions but resets on reload (Priority: P3)

A shopper's filter should feel like a transient view, not a saved preference. When they navigate away or reload, the page starts fresh.

**Why this priority**: This is a behaviour-shaping rule, not a new capability. It prevents surprise on return visits.

**Independent Test**: Type a query, reload the page. Search box should be empty, all products visible.

**Acceptance Scenarios**:

1. **Given** the shopper has filtered the list with `tank`,
   **When** they reload the page,
   **Then** the search box is empty and all products are visible.

2. **Given** the shopper has filtered the list,
   **When** they navigate to a product detail page and click the browser Back button,
   **Then** they return to the list page in its default (unfiltered) state. *(Acceptable degradation: filter may also be preserved — either is fine — but the page MUST be usable; nothing must be in a stuck filtered state with a hidden query.)*

---

### Edge Cases

- **Empty query** (box cleared or never typed in): the full product list is shown. No "no results" message.
- **Whitespace-only query** (`"   "`): treated as an empty query — full list is shown.
- **Query with leading/trailing whitespace** (`"  sun "`): whitespace is ignored; matching is on the trimmed value.
- **Query with characters that have special meaning in search engines** (`*`, `?`, `"`, `(`): treated as plain characters; no wildcard, regex, or operator semantics. The shopper sees what they typed, matched literally as a substring of the name.
- **No matches**: the product grid is empty and a short, plain-language "no results" message is shown in its place. The search box keeps the current query so the shopper can edit it.
- **Catalogue size**: the page already loads the full set today; filtering happens over whatever is on the page. The feature does not change how many products are loaded.
- **Slow typist / fast typist**: filtering keeps up with typing; the list reflects the latest keystroke without flicker.
- **Accessibility**: the search box has a visible label or accessible name so screen-reader users can identify it; results updates are perceivable (e.g. via the visible list itself — explicit live-region announcement is a nice-to-have, not required for v1).

## Requirements *(mandatory)*

### Functional Requirements *(observable behaviour, not implementation)*

- **FR-001**: The product list page MUST display a search input control above (or visually adjacent to) the product grid.
- **FR-002**: When the search input is empty, the page MUST display the same set of products it would display today without this feature.
- **FR-003**: When the shopper types a non-empty query, the page MUST hide every product whose name does not contain the query as a substring, and show every product whose name does.
- **FR-004**: Matching MUST be case-insensitive (`Sun`, `sun`, and `SUN` produce the same result).
- **FR-005**: Matching MUST be on the product **name** only. Description, category, and other fields MUST NOT influence whether a product is shown or hidden.
- **FR-006**: Leading and trailing whitespace in the query MUST be ignored.
- **FR-007**: The list MUST update in response to the shopper's typing without a full page reload and without an additional round-trip to the server for each keystroke.
- **FR-008**: When the query matches zero products, the page MUST replace the product grid with a short, plain-language "no results" message and MUST keep the query visible in the search box.
- **FR-009**: The shopper MUST be able to return to the full, unfiltered list by clearing the search box (deleting all characters is sufficient; a one-click reset affordance is desirable).
- **FR-010**: Reloading the product list page MUST reset the filter — the search box starts empty and all products are visible.
- **FR-011**: The search input MUST have an accessible name (label, `aria-label`, or equivalent) so assistive-technology users can identify it.
- **FR-012**: The feature MUST NOT alter the behaviour of any page other than the product list page.
- **FR-013**: The feature MUST NOT add new product fields, new product records, or new categories to the catalogue.

### Key Entities

- **Product (existing)**: An item already present in the catalogue. Relevant attribute for this feature: `name` (the human-readable product name shown on the list page). All other attributes are untouched.
- **Search Query (new, transient)**: The text the shopper has currently typed into the search box. It lives only in the current browser tab and is discarded on reload. It is not persisted, not sent to the server, and not associated with the shopper's account or cart.

## Out of Scope / Do NOT *(explicit constraints)*

These constraints are mandatory. Any plan or implementation that violates them is rejected.

- **DO NOT add new services.** Use only the services already in this repo.
- **DO NOT introduce Elasticsearch, Solr, vector databases, or any new datastore.**
- **DO NOT change the source of product data.** The product list continues to come from the existing in-memory catalogue loaded from `productcatalogservice/products.json`. Filtering happens over what is already on the page.
- **DO NOT change the implementation language of the services being edited.** The frontend is Go; productcatalogservice is Go. If `productcatalogservice` is extended (it shouldn't need to be for v1), use the existing protobuf / gRPC patterns.
- **DO NOT add new infrastructure config, Helm charts, Kubernetes manifests, or environment variables.** CI deploys whatever lands on `attendee/<your-name>`; the feature must work with the existing deployment as-is.
- **DO NOT touch the build pipeline.** Stay inside this branch and this repo.
- **DO NOT search by description, category, price, or any field other than name** in v1. (Description-based search is tracked separately as the `search-misses-descriptions` bug; it is NOT part of this feature.)
- **DO NOT add search history, suggestions, autocomplete, fuzzy / typo-tolerant matching, ranking, or analytics** in v1.
- **DO NOT change pagination, sorting, or the number of products loaded** by the product list page.
- **DO NOT call the backend on every keystroke.** Filtering is over already-loaded data.
- **DO NOT persist the query** across page reloads or across sessions.

## Success Criteria *(mandatory, testable, technology-agnostic)*

### Measurable Outcomes

- **SC-001**: A shopper looking for a product whose exact name they know can locate it on the list page within **5 seconds** of starting to type, without scrolling and without a page reload.
- **SC-002**: For any query that matches at least one product name, the visible list updates to reflect that query within **200 milliseconds** of the last keystroke.
- **SC-003**: Clearing the search box restores the full product list in under **200 milliseconds** with no page reload.
- **SC-004**: For a query that matches no product name, the shopper sees a "no results" message in place of the grid within the same 200 ms budget, and the query they typed remains visible and editable.
- **SC-005**: Across **100%** of acceptance scenarios listed above (User Stories 1–3), the observed behaviour matches the expected outcome.
- **SC-006**: The product list page issues **zero additional server requests per keystroke** compared with the current page (network traffic per keystroke does not increase).
- **SC-007**: The search input has an accessible name verifiable by an automated accessibility check (e.g. axe) — **zero new accessibility violations** are introduced on the product list page.
- **SC-008**: Reloading the product list page resets the filter in **100%** of trials (the search box is empty and the full list is visible after every reload).

## Assumptions

- The product list page today renders the full catalogue server-side; filtering operates over what the browser has already received. (If product pagination is introduced later, this feature's scope explicitly does not extend to fetching additional pages — see Out of Scope.)
- "Name" refers to the existing `name` field on each product in `productcatalogservice/products.json` (e.g. `"Sunglasses"`, `"Tank Top"`, `"Watch"`).
- Shoppers are using a modern browser with JavaScript enabled. Behaviour with JavaScript disabled is unchanged from today (i.e. the page still renders the full list; the search box may be absent or inert).
- The product catalogue size is small enough that filtering the already-rendered list is imperceptible to the shopper. (At time of writing, the catalogue is ~10 products. The feature is designed for "small catalogue on one page", not for thousands of items.)
- The frontend service is the only service that needs to change. `productcatalogservice` is not modified.
- Existing styling and layout of the product grid are reused; the search box adopts the page's existing visual style.
