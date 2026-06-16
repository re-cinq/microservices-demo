# Quickstart: Save products and view them later

A validation guide proving the feature works end to end. Implementation details live in
`tasks.md` (after `/speckit-tasks`) and the code itself — this is a run/verify guide.

## Prerequisites

- The Online Boutique frontend running locally (existing dev workflow for
  `src/frontend`), with the product catalogue service reachable as today.
- A browser with cookies enabled (the saved list is keyed on the existing session
  cookie).

## Unit / handler tests

From `src/frontend`:

```sh
go test ./...
```

Expected: the new store tests (`wishlist_store_test.go`) and handler tests pass —
covering add/list/contains, idempotency (no duplicates), per-session isolation, and the
save → redirect and view → render flows.

## Manual end-to-end validation

Maps to the spec's acceptance scenarios:

1. **Save from a product page** (Scenario 1, FR-001/FR-002)
   - Open any product page → click **Save**.
   - Expect: page re-renders with the control showing the product as **saved**.

2. **View saved products and click through** (Scenario 2, FR-003/FR-005)
   - Click the **saved-products link in the header** → the saved view lists the product.
   - Click the product in the saved view → land on its product page.

3. **Persists across the session** (Scenario 3, FR-006)
   - Browse to the home page and other products, then reopen the saved view.
   - Expect: the saved product is still listed, with no sign-in at any point.

4. **Empty state** (Scenario 4, FR-008)
   - In a fresh session (new browser profile / cleared cookies), open the saved view.
   - Expect: a clear "nothing saved yet" message with a prompt to start browsing — not a
     blank or broken page.

5. **No duplicates** (FR-007)
   - Save the same product twice.
   - Expect: it appears once in the saved view.

6. **Unavailable product** (edge case)
   - With a product saved, the saved view still renders the remaining valid products if
     one saved ID is no longer in the catalogue.

## Constitution spot-check

- No new service or datastore was added (Principles I, II): the only new runtime state
  is the in-memory store in the frontend; product data still comes from the catalogue.
- No infra/manifest/CI files changed (Principle IV): `git status` should show changes
  only under `src/frontend/` (plus these spec docs).
