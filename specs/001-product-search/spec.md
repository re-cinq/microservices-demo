# Feature Specification: Product Search


## Problem

Online Boutique shoppers landing on the product list page today have no way to narrow the catalogue to items they're looking for. The full list of products is already rendered on the page, but a shopper who knows what they want (e.g. "watch", "candle") has to scan the entire grid to find it.

In the user's own words:

> "Add a product search feature to Online Boutique. A search box on the product list page that filters the products already loaded, by name, in the browser. No backend change."

The job to be done: let a shopper narrow the visible product grid to items whose name matches what they type, instantly, using only the products already on the page.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Filter visible products by name as I type (Priority: P1)

A shopper opens the product list page and sees a single search input directly above the existing product grid. The grid is unchanged: same products, same order, same cards. As soon as the shopper types one or more characters into the input, the grid narrows to only those products whose name contains what they typed. Editing or clearing the input updates the grid in lock-step. The full grid returns when the input is empty.

**Why this priority**: this is the feature. Without it there is nothing to ship.

**Independent Test**: load the product list page (catalogue contains 9 products: Sunglasses, Tank Top, Watch, Loafers, Hairdryer, Candle Holder, Salt & Pepper Shakers, Bamboo Glass Jar, Mug). Type `watch` → exactly one product card visible (Watch). Type `WATCH` → same card. Clear the input → all 9 cards visible again.

**Acceptance Scenarios**:

1. **Given** the product list page is loaded with all 9 products visible, **When** the shopper types `watch` into the search input, **Then** exactly one product card is visible (Watch), and the other 8 are hidden.
2. **Given** the search input contains `watch` and only the Watch card is visible, **When** the shopper deletes all characters so the input is empty, **Then** all 9 product cards are visible again, in their original order.
3. **Given** the search input is empty and all 9 products are visible, **When** the shopper types `WATCH` (uppercase), **Then** exactly the Watch card is visible — matching is case-insensitive.
4. **Given** the search input is empty, **When** the shopper types `ass`, **Then** exactly two product cards are visible (Sunglasses, Bamboo Glass Jar) — substring match, anywhere in the name.
5. **Given** the search input is empty, **When** the shopper types `&`, **Then** exactly one product card is visible (Salt & Pepper Shakers) — special characters in the query are treated as literal characters.
6. **Given** the shopper is filtering by name, **When** they type, edit, or delete any character, **Then** no new network request is made by the page whose purpose is to fetch product catalogue data (verifiable in the browser network panel: zero new requests matching the product-catalogue endpoint pattern between page load and any sequence of typing).
7. **Given** the shopper types characters very quickly in succession, **When** the final keystroke is processed, **Then** the visible grid matches the final query exactly — no card from a stale intermediate query remains visible.

---

### User Story 2 — Clear feedback when nothing matches (Priority: P2)

When the shopper's query matches no product names, the grid area is replaced with a clear empty-state message that (a) explains there are no matches, (b) is announced to assistive technologies, and (c) gives the shopper an obvious way to recover. The shopper must never see a silently empty grid that looks like a broken page.

**Why this priority**: a silently empty grid looks like a broken page. P2 because P1 still delivers value without this, but ship-quality demands this feedback.

**Independent Test**: with the product list page open, type `zzzzz` and confirm an explicit empty-state message appears in place of the grid. Backspace one character at a time and confirm the message disappears as soon as the query matches at least one product, with the matching cards rendered in place.

**Acceptance Scenarios**:

1. **Given** the product list page is loaded, **When** the shopper types `zzzzz` (a query that matches no product name), **Then** the grid area contains an empty-state element with visible text `No products match your search.` and **does not** contain any product cards.
2. **Given** the empty-state message is shown, **When** the shopper deletes characters so the query is `z` (still no match) then `` (empty), **Then** the empty-state message disappears as soon as the query becomes empty, and all 9 product cards are visible again.
3. **Given** the empty-state message is shown, **When** the shopper edits the query so it now matches `Watch`, **Then** the empty-state message disappears and the Watch card is rendered in the grid area.
4. **Given** the empty-state element exists, **When** a shopper using a screen reader reaches this state, **Then** the message is announced by the screen reader without interrupting whatever the shopper was doing.
5. **Given** the shopper types the literal string `<script>alert(1)</script>` into the search input, **When** the empty-state message is rendered, **Then** no script executes, no `alert` dialog appears, and any echoed query text is HTML-escaped (verifiable: the rendered DOM contains the escaped text, not a live `<script>` element).

---

### User Story 3 — Quick clear of the active filter (Priority: P3)

A visible clear control inside the search input becomes active when the input is non-empty. Activating it once empties the input and returns the full grid. The Esc key, while the input has focus, does the same thing.

**Why this priority**: convenience. P1 already lets the shopper backspace; this is the low-effort polish that makes the feature feel finished.

**Independent Test**: with `watch` typed and the grid filtered to one card, click the clear control inside the input — input becomes empty, grid returns to all 9 cards. Repeat the same test, but instead of clicking, press Esc with the input focused — same result.

**Acceptance Scenarios**:

1. **Given** the search input is empty, **When** the page is loaded, **Then** the clear control is either not rendered or not visible / not focusable.
2. **Given** the shopper has typed `watch` and the grid is filtered, **When** the shopper activates the clear control once (mouse click or keyboard activation), **Then** the search input becomes empty and all 9 product cards are visible again.
3. **Given** the search input contains `watch` and has keyboard focus, **When** the shopper presses Esc, **Then** the search input becomes empty and all 9 product cards are visible again.

---

### Edge Cases

- **Empty catalogue at page load**: if `productcatalogservice` returns zero products, the page renders today's zero-products view; the search input is still present but typing matches nothing — the empty-state message is shown for any non-empty query.
- **Special characters in the query** (`.`, `*`, `\`, `?`, `[`, `]`, `&`, `(`, `)`, `'`, `"`): treated as literal characters. `&` matches `Salt & Pepper Shakers`. `.*` matches only product names containing the literal substring `.*`.
- **HTML / script content in the query** (e.g. `<script>alert(1)</script>`, `<img src=x onerror=...>`): never executed; if echoed in the UI, must be HTML-escaped (see User Story 2, Scenario 5).
- **Leading or trailing whitespace** in the query: trimmed before matching. `"  watch  "` behaves the same as `"watch"`. Whitespace inside the query is preserved as a literal part of the substring (so `red mug` matches only names containing the literal substring `red mug`).
- **Very long pasted query** (several hundred characters): no crash, no UI hang; simply matches no products and shows the empty-state message.
- **Rapid typing**: the grid always reflects the final query at rest; a stale intermediate query must not leave cards visible (User Story 1, Scenario 7).
- **Product with missing / empty name in `products.json`**: treated as a non-match for any non-empty query; ignored when computing whether the full grid is shown for an empty query (i.e. it's part of the "all products" set rendered server-side, but it cannot be matched by any search).
- **Non-ASCII characters**: matched using simple Unicode case folding (lowercase comparison). Accent folding is **not** included in v1 (typing `cafe` does **not** match `Café`).
- **Page reload or navigation away and back**: the search input always starts empty; no query is persisted (no URL parameter, no cookie, no localStorage).
- **No JavaScript available**: the page still renders the full product list as it does today (server-rendered grid). The search input may be present but non-functional. The page is not broken; the shopper sees the unfiltered catalogue.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The product list page MUST display a single text input control directly above the existing product grid on initial render. No existing product card, recommendation block, header, footer, breadcrumb, currency selector, or "Hot Products" section may be removed, reordered, resized, or hidden by the addition of the search input.
- **FR-002**: While the trimmed search query is non-empty, the only product cards visible in the grid MUST be those whose product name contains the query as a substring.
- **FR-003**: When the trimmed search query is empty, every product card normally rendered on the page today MUST be visible, in its original order.
- **FR-004**: Matching MUST be against the product **name** only — not description, id, category, or price.
- **FR-005**: Matching MUST ignore case: searching for `Candle`, `candle`, and `CANDLE` MUST return the same products. Matching MUST NOT ignore accents: searching for `cafe` MUST NOT match a product named `Café`.
- **FR-006**: Special characters in the query (including `.`, `*`, `\`, `?`, `[`, `]`, `&`, `(`, `)`, `'`, `"`, `<`, `>`) MUST be treated as literal characters of the search term, not as wildcards or pattern syntax. Searching for `&` MUST match `Salt & Pepper Shakers`.
- **FR-007**: When the trimmed query is non-empty and matches zero products, the grid area MUST display a visible message reading `No products match your search.` (or a clearly-equivalent localised string), MUST NOT show any product cards, and MUST announce the message to screen-reader users.
- **FR-008**: A visible clear control MUST appear inside or directly adjacent to the search input whenever the input contains at least one character. Activating the clear control once (by mouse click or keyboard) MUST empty the input and restore the full grid. Pressing Esc while the input has focus MUST do the same thing.
- **FR-009**: Typing any text into the search input — including text that looks like code or markup, such as `<script>alert(1)</script>` or `<img src=x onerror=...>` — MUST produce only normal search behaviour. No popup or dialog appears, no part of the page is broken or rearranged, and nothing typed in the search input is interpreted as code by the page.
- **FR-010**: The feature MUST render results inline within the existing product list page. Typing in the search input MUST NOT navigate the browser to a different URL, add an entry to browser history (so the Back button is unaffected by typing), open a modal, or render a separate "search results" page.
- **FR-011**: The feature MUST be available to every visitor of the product list page, including visitors who are not signed in. No sign-in, sign-up, session, or other authentication step MUST be required to use the search input.
- **FR-012**: The search input MUST start empty every time the product list page is loaded — including after a full page reload, and after the shopper navigates away from the page and comes back to it.

### Key Entities

- **Product**: a catalogue item already modelled by the existing system. The only attribute this feature reads is `name`. No new fields are introduced.
- **Search Query**: the shopper's transient input. Lives only in the page session; not persisted, not sent to any backend, not shared between shoppers.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A search input is present on the product list page on initial render, positioned directly above the product grid — verifiable by a single page load with no interactions.
- **SC-002**: For every query in the test set `{"", "watch", "WATCH", "Watch", "ca", "ass", "&", "salt", "zzzzz"}`, the set of visible product cards matches the set computed by case-insensitive substring of the trimmed query against the `name` field of each product in `productcatalogservice/products.json` — for every test query, 100% of expected products are visible and 0% of non-expected products are visible.
- **SC-003**: For any single keystroke in the search input, the time from the keystroke event firing to the DOM reflecting the updated set of visible cards is under 100 ms — verifiable by browser performance instrumentation against the current 9-product catalogue.
- **SC-004**: Zero network requests whose response body contains product-catalogue data are made by the browser between the initial page load completing and the page being unloaded, across any sequence of typing, editing, deleting, or clearing the query — verifiable in the browser's network panel.
- **SC-005**: For a query of `zzzzz`, the rendered grid area contains an element with the visible text `No products match your search.` (or its localised equivalent) and contains zero product cards — verifiable by DOM inspection.
- **SC-006**: For a query of `<script>alert(1)</script>` typed into the search input, no `alert` dialog appears, no script executes, and the DOM contains the literal escaped characters `&lt;script&gt;` if the query is echoed anywhere — verifiable by an automated test that fails if the page emits an `alert` or contains an unescaped `<script>` element in the search/empty-state region.
- **SC-007**: With the search input empty, the rendered product list page MUST contain the same set of product cards, in the same order, with the same per-card content (image, name, price), as the page rendered without the feature deployed — verifiable by a side-by-side DOM diff scoped to the product grid region.
- **SC-008**: Clicking the clear control once, or pressing Esc with the input focused, returns the grid to the full 9-card state from any non-empty query state in a single user action.
- **SC-009**: An anonymous (logged-out, no session cookie) visitor can load the product list page and use the search input with the same observable behaviour as a logged-in visitor — verifiable by running the SC-002 test suite in a fresh incognito session with no prior cookies.
- **SC-010**: Typing in the search input never changes `window.location.href`, never adds an entry to browser history, and never causes any element with role `dialog`, `alertdialog`, or similar to appear — verifiable by observing `history.length` and the DOM before and after a sequence of keystrokes.

## Constraints *(what we MUST NOT do)*

These are immutable boundaries inherited from the project brief plus the senior-engineer-review pass. Spec, plan, and implementation must all stay inside them.

- **C-001**: No new services. Only services already present in this repo (`frontend`, `productcatalogservice`, and the rest of the existing microservice set) may be touched.
- **C-002**: No Elasticsearch, Solr, vector database, or any other new datastore — not even a lightweight in-process index introduced as a new component.
- **C-003**: The catalogue is the existing in-memory product list loaded by `productcatalogservice` from `productcatalogservice/products.json`. Any in-memory filter MUST read from that existing data, not a new store.
- **C-004**: Language must match the service being edited. The `frontend` service is Go; `productcatalogservice` is Go. Any change to either service is in Go.
- **C-005**: If `productcatalogservice` is extended, the extension MUST follow the existing protobuf / gRPC patterns already used by that service. No new RPC frameworks, no new REST endpoints.
- **C-006**: No new infrastructure config: no new Helm charts, no new Kubernetes manifests, no new environment variables, no new ConfigMaps, no new Secrets.
- **C-007**: Development for this feature happens on the `001-product-search` branch of this repo. The final feature MUST land on `attendee/jenny-svennefelt` (the branch CI deploys from) via a normal merge / fast-forward. No fork, no separate repo, no submodule. No other branches are created.
- **C-008**: The build pipeline (CI, image build, deploy steps) MUST NOT be modified.
- **C-009**: No backend change is required by the feature description. Any backend-side work must be justified as making the observable behaviour above achievable within the other constraints; if a pure-frontend implementation satisfies the spec, no `productcatalogservice` change should be made.
- **C-010**: The existing product list page layout MUST NOT change when the search input is empty: same product cards, same order, same per-card content, same surrounding elements (header, breadcrumb, currency selector, recommendations / "Hot Products", footer). The search input is an **addition** above the grid, not a replacement or rearrangement of anything.
- **C-011**: No authentication, no login, no session check, no permission gate may be added by this feature. The search input is available to every visitor of the product list page on the same terms as the existing product grid.
- **C-012**: No separate search results page. Results render inline in the existing product list view. The URL MUST NOT change when the shopper types, no new route is added, no modal / lightbox / popover renders results, no new browser-history entries are pushed.
- **C-013**: While the shopper types in the search input, the page MUST NOT issue any additional request to any service to fetch product catalogue data. Filtering uses only the product list already loaded into the page at first render. (The customer-visible result — instant feedback — is captured by SC-003; this constraint pins the engineering means.)
- **C-014**: The product catalogue service contract MUST NOT change. The set of products returned to any caller, by any product-fetching call, MUST be the same after this feature ships as before. This feature only affects which already-loaded product cards are visible on the product list page; it does not alter the catalogue, the service that serves it, or what any other page or feature receives from that service.
- **C-015**: The search query MUST NOT be persisted in any URL parameter, cookie, localStorage, sessionStorage, browser-history entry, or server-side store. (FR-012 captures the customer-visible outcome — input starts empty on reload; this constraint forecloses the engineering shortcuts that would break it.)

## Assumptions

- **A-001**: "Filters … by name" means the product `name` field only. Matching against description, id, category, or price is out of scope unless the spec is updated.
- **A-002**: Filtering is **live** — applied as the shopper types — rather than requiring an explicit submit action (Enter / button click).
- **A-003**: Case-insensitive substring match is the matching rule. Not fuzzy matching, not prefix-only, not regex, not multi-word tokenisation, not accent folding.
- **A-004**: Leading and trailing whitespace in the query are trimmed before matching; whitespace inside the query is preserved as a literal part of the substring.
- **A-005**: The "products already loaded" set is exactly the set the product list page renders today; this feature does not change how many products are loaded onto the page.
- **A-006**: The product list page in this repo's `frontend` service serves the product grid via server-rendered templates plus existing static assets. Filtering is implemented as a client-side enhancement on top of what the server already renders, with no new server endpoint required. (To be confirmed during `/speckit-plan`.)
- **A-007**: The catalogue size is small enough (currently 9 products) that no debouncing, indexing, virtualisation, or worker is required for SC-003 (sub-100 ms keystroke-to-DOM) to hold.
- **A-008**: Accessibility expectations for the new search input: a visible `<label>` or `aria-label`, focus order respects DOM order, the empty-state element is a polite live region (`role="status"` or `aria-live="polite"`). No additional accessibility commitments (e.g. result-count announcements on every keystroke) are made for v1.
- **A-009**: The search input is not auto-focused on page load. The page's existing focus behaviour is preserved.
- **A-010**: No analytics / telemetry events are emitted by typing in the search input for v1. (If analytics are wired up later, that's a separate spec.)
