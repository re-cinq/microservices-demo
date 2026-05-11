# Feature Specification: Product Search

**Feature Branch**: `attendee/kate-payne-fix`  
**Created**: 2026-05-11  
**Status**: Draft  

---

## Overview

Shoppers browsing Online Boutique have no way to find a specific product without
scrolling through every item on the catalogue page. If they arrive with a
particular item in mind — or only a vague idea — there is no shortcut: they
must visually scan every product card. This friction increases with catalogue
size and costs conversions.

This feature adds a search input to the catalogue page so shoppers can filter
products by name or description, see the matching text highlighted, and get a
clear message when nothing matches.

---

## Constraints *(what this feature must NOT do)*

- MUST NOT introduce new backend services beyond those already in the repository.
- MUST NOT use Elasticsearch, Solr, vector databases, or any datastore other
  than the existing in-memory product catalogue.
- MUST NOT add new environment variables, Kubernetes manifests, Helm charts, or
  infrastructure configuration files.
- MUST NOT modify the CI/CD pipeline or build configuration.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Basic Keyword Search (Priority: P1)

A shopper arrives at the catalogue page knowing roughly what they want. They
type a word into the search box, submit it, and instantly see only the products
that mention that word — either in the product name or in the description. The
rest of the catalogue disappears until they clear the search.

**Why this priority**: Filtering to relevant results is the core value of the
feature. Nothing else is useful without it.

**Independent Test**: Can be fully tested by loading the catalogue, entering a
term known to match some products and not others, and verifying the result set
changes accordingly. Delivers immediate, standalone value to shoppers.

**Acceptance Scenarios**:

1. **Given** the catalogue page is loaded and all products are visible,  
   **When** the shopper types `"camera"` into the search field and submits,  
   **Then** only products whose name or description contains `"camera"`
   (case-insensitive) are shown, and all non-matching products are hidden.

2. **Given** the catalogue page is loaded,  
   **When** the shopper types `"SUNGLASSES"` (all uppercase) and submits,  
   **Then** the same products are returned as for `"sunglasses"`, confirming
   that matching is case-insensitive.

3. **Given** the shopper is on any page of the site,  
   **When** they navigate to the catalogue page without a search term,  
   **Then** all products are displayed (no filtering applied).

4. **Given** the catalogue page is loaded,  
   **When** the shopper types a single character (e.g., `"a"`) and attempts to
   submit,  
   **Then** the search does not execute, no results change, and an inline
   message is shown indicating that at least 2 characters are required.

5. **Given** the catalogue page is loaded,  
   **When** the shopper types only whitespace (e.g., two spaces) and attempts
   to submit,  
   **Then** the search does not execute and the same inline minimum-length
   message is shown.

---

### User Story 2 — Empty-State "No Results" Message (Priority: P2)

A shopper searches for a term that matches no product. Instead of a blank,
unexplained page they see a clear message telling them what they searched for
and that nothing was found, so they know the search worked and can try again.

**Why this priority**: Without an empty state, a blank catalogue page looks
like a loading error. This is the minimum necessary to make failure legible.

**Independent Test**: Can be fully tested by submitting a search term that is
guaranteed to match nothing (e.g., a random string) and verifying the empty-
state UI appears with the correct term echoed back.

**Acceptance Scenarios**:

1. **Given** the catalogue page is loaded,  
   **When** the shopper submits the search term `"xyzzy123"` (no matching
   products),  
   **Then** no product cards are displayed, and a message is visible on the
   page that includes the literal text `"xyzzy123"` and communicates that no
   products were found.

2. **Given** the empty-state message is displayed,  
   **When** the shopper clears the search field and submits an empty query,  
   **Then** the full product catalogue is restored and the empty-state message
   is no longer visible.

---

### User Story 3 — Search-Term Highlighting in Results (Priority: P3)

A shopper sees their search term visually marked within the product names and
descriptions in the results list. This confirms which part of the text caused
each product to match and helps them quickly assess relevance.

**Why this priority**: Highlighting aids comprehension but is not required for
the search to function. It is a quality-of-experience improvement delivered on
top of the working filter.

**Independent Test**: Can be fully tested by searching for a term that appears
mid-word or mid-sentence in several product descriptions and checking that
exactly those substrings are visually distinct from the surrounding text.

**Acceptance Scenarios**:

1. **Given** results are displayed for the search term `"lens"`,  
   **When** the shopper looks at the product list,  
   **Then** every occurrence of `"lens"` (case-insensitive) within product
   names and descriptions is visually distinct from the surrounding text (e.g.,
   rendered bold, underlined, or with a background highlight).

2. **Given** results are displayed for `"lens"`,  
   **When** the shopper inspects a product whose name contains `"Lens"` and
   whose description also contains `"lens"`,  
   **Then** both occurrences are highlighted, not just the first.

3. **Given** results are displayed for `"lens"`,  
   **When** the shopper inspects a product whose description contains
   `"lenses"`,  
   **Then** the substring `"lens"` within `"lenses"` is highlighted, confirming
   substring matching.

---

### Edge Cases

- What happens when the search term is fewer than 2 non-whitespace characters?
  The search MUST NOT execute; the submit action must be prevented and the
  shopper shown an inline prompt indicating that at least 2 characters are
  required. The catalogue remains in its current state (full list or prior
  results).
- What happens when the search term contains special regex or HTML characters
  (e.g., `<script>`, `"`)?
  The displayed highlight must not execute code or break page layout — the term
  must be treated as a literal string.
- What happens when every product matches (e.g., the user types a letter
  present in all product names)?
  All products are shown with the matching text highlighted; no empty state
  appears.

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The catalogue page MUST display a search input field that is
  visible without scrolling on a standard desktop viewport.
- **FR-001a**: The search field MUST prevent submission and display an inline
  validation message when the entered term contains fewer than 2
  non-whitespace characters. The catalogue MUST remain in its current state
  (unfiltered or showing prior results) when submission is blocked.
- **FR-002**: When a shopper submits a valid search term (2 or more
  non-whitespace characters), the page MUST display only products whose name
  or description contains that term, using case-insensitive, substring
  matching.
- **FR-003**: When no products match the submitted search term, the page MUST
  display an empty-state message that includes the exact search term the
  shopper entered and communicates that no products were found.
- **FR-004**: In search results, every occurrence of the search term within
  product names and descriptions MUST be rendered visually distinct from
  surrounding text (highlight, bold, or equivalent).
- **FR-005**: A shopper MUST be able to clear the search (via an explicit reset
  control, or by emptying the field) and see the full catalogue restored
  without a full page reload being required. Clearing the field does not
  constitute a search submission and MUST NOT trigger the minimum-length
  validation message.
- **FR-006**: Search MUST check both the product name and the product
  description fields; a product matches if the search term appears in either
  field, or in both.
- **FR-007**: The search input MUST preserve the submitted term as its value
  after results are rendered, so the shopper can see what they searched for and
  edit it.

### Key Entities

- **Product**: An item in the catalogue with a name and a description (and
  other attributes such as price and image). The search operates on the name
  and description fields only.
- **Search Term**: The literal string entered by the shopper. Matching is
  case-insensitive and substring-based. An empty or whitespace-only term is
  treated as "no filter".
- **Search Result**: The subset of products whose name or description contains
  the search term. May be empty.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A shopper can go from landing on the catalogue page to seeing
  filtered results for a known product in under 5 seconds on a standard
  connection, provided the search term is at least 2 non-whitespace characters.
- **SC-001a**: Submitting a term shorter than 2 non-whitespace characters never
  changes the displayed product list — verified by attempting single-character
  and whitespace-only submissions and confirming the catalogue is unchanged.
- **SC-002**: Search correctly returns results from both product name and
  description fields — validated by a test product whose keyword appears only
  in the description.
- **SC-003**: When a search term matches zero products, an empty-state message
  is displayed that echoes the search term — verified for at least three
  distinct non-matching terms.
- **SC-004**: 100% of matching substrings in the result list are visually
  highlighted — verified against at least one product where the term appears
  multiple times and at least one where it appears only in the description.
- **SC-005**: Clearing the search restores the full catalogue in a single user
  action — no navigation away from the page required.

---

## Assumptions

- All products are already loaded into memory by the product catalogue service
  at startup; there is no need to query an external datastore.
- The existing catalogue page is the correct and only surface for this feature;
  a dedicated search results page is out of scope.
- Pagination, if present on the catalogue page, is out of scope for this
  iteration — search results are shown as a flat, unpaginated list.
- Search is client-triggered (the shopper submits the term); real-time
  filtering as the user types is out of scope for v1.
- The feature targets logged-in and anonymous shoppers equally; no
  personalisation or per-user search history is required.
- Accessibility beyond standard semantic HTML (e.g., screen-reader
  announcements of result counts) is desirable but not a blocking requirement
  for v1.
