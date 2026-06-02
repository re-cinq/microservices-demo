# Spec: Frontend Product Filter

## Overview

Add a search box to the product listing page (homepage) that filters the already-loaded product cards in the browser by name. No backend requests are made during filtering — all matching happens client-side against the DOM or in-memory product data.

---

## Functional Requirements

### FR-1 — Search input
- A text input field is rendered at the top of the product grid on the homepage.
- The input has placeholder text: `Search products…`
- The input is visible on page load with no pre-filled value.

### FR-2 — Real-time filtering
- As the user types, product cards are shown or hidden instantly (no submit button, no page reload).
- Filtering triggers on every keystroke (input event).
- Matching is case-insensitive.
- Matching is against the product **name** only.

### FR-3 — Partial match
- Products whose name **contains** the query string are shown.
- Products whose name does not contain the query string are hidden.

### FR-4 — Empty query
- When the search box is empty, all products are shown.

### FR-5 — No results state
- If no products match, a message is displayed in the grid area: `No products match your search.`
- The message disappears as soon as at least one product matches.

### FR-6 — No backend involvement
- No new HTTP requests are made during filtering.
- The feature works entirely with the HTML already rendered on the page.

---

## User Scenarios

**Scenario A – User searches by exact name**
Given the homepage is loaded with all products visible,
when the user types `Sunglasses` into the search box,
then only product cards whose name contains `Sunglasses` remain visible.

**Scenario B – User searches by partial name**
Given the homepage is loaded,
when the user types `can`,
then all products whose name contains `can` (case-insensitive, e.g. "Candle") are shown and others are hidden.

**Scenario C – User clears the search**
Given some products are filtered out,
when the user clears the search box,
then all product cards are shown again.

**Scenario D – No match**
Given the homepage is loaded,
when the user types `xyznotaproduct`,
then all product cards are hidden and the message `No products match your search.` is displayed.

---

## Success Criteria

| # | Criterion | Testable by |
|---|-----------|-------------|
| SC-1 | Search box is present on the homepage | Visual inspection / DOM query |
| SC-2 | Typing a product name shows only matching cards | Manual test with known product name |
| SC-3 | Filtering is case-insensitive | Type uppercase variant of a known name |
| SC-4 | Clearing the input restores all products | Clear field, count visible cards |
| SC-5 | No network requests fire during filtering | Browser DevTools Network tab |
| SC-6 | "No products match" message appears when zero results | Type a string that matches nothing |
| SC-7 | "No products match" message disappears when results return | Append/remove chars to go from 0 to 1+ results |
