# Quickstart: Recently Viewed Products

How to build, run, and verify the product-page recently-viewed strip.

## Run the frontend locally

The change is confined to `src/frontend`. Build/test it like any Go service in this repo:

```bash
cd src/frontend
go build ./...
go test ./...        # includes recentlyviewed_test.go
```

To exercise the UI, run the full demo (e.g. the existing `docker compose` / skaffold / kubernetes flow used in this repo — unchanged by this feature) and open the storefront. No new env vars, services, or manifests are required.

## Manual verification (acceptance criteria)

1. **Fresh session, first product (AC4)**
   - Open a product page in a new browser session.
   - ✅ Expect: **no** "Recently viewed" strip; page looks exactly as before.

2. **Strip appears, most-recent-first, current excluded (AC1)**
   - View product A, then product B, then product C.
   - ✅ On C's page: strip at the **bottom** lists **B then A** (most-recent-first); **C is not** in the strip.

3. **Cap at 4 (AC2)**
   - View 5+ distinct products, then open another product page.
   - ✅ Strip shows only the **4** most recently viewed (excluding the current one).

4. **Click-through (AC3)**
   - Click any product in the strip.
   - ✅ You land on that product's page.

5. **Half-size thumbnails (AC5)**
   - ✅ Strip thumbnails are visibly **half** the size of the main/recommended product thumbnails on the page.

6. **Missing product (FR-010 edge case)**
   - If a previously viewed product ID is no longer in the catalogue, the strip ✅ renders the remaining products without error.

## Out-of-scope reminders (do NOT expect these here)

- **Home-page strip** → AIP-200.
- **De-duplication / move-to-front** → AIP-201. Revisiting the same product *may* show duplicate entries in this story.
- **Cross-device sync** → out of scope (epic).
