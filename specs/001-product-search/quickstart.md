# Quickstart: Verifying Product Search

Once `/speckit-implement` lands the change on `001-product-search`, use the recipe below to verify the feature delivers what `spec.md` promised and that no constraint has been violated.

## Prerequisites

- A working Kubernetes cluster with Online Boutique deployed (the existing `skaffold run` / `kubectl apply -f release/kubernetes-manifests.yaml` path — no new infra steps).
- `kubectl` configured to talk to that cluster.
- The new `001-product-search` branch built and deployed (via the CI pipeline that auto-deploys `attendee/<your-name>` — same pipeline that picked up `attendee/tim-love-fix`).

## 1. Smoke test — happy path

```text
1. Open the storefront frontend external IP in a browser.
2. Confirm a search input is visible in the header.
3. Type `watch` and submit.
4. Expect: results page lists the Watch product, with picture, name, and price.
5. Click the result.
6. Expect: the existing /product/{id} detail page renders for the Watch.
```

Repeat with `sun`, `loaf`, `dryer` — confirms the substring match works at the start, middle, and end of names respectively.

## 2. Case-insensitivity

```text
Submit each of: Watch, watch, WATCH, wAtCh.
Expect: all four produce the same result set.
```

Maps to spec User Story 1 Acceptance Scenario 3 and edge-case bullet "case sensitivity".

## 3. Currency switching

```text
1. From any page, switch the currency selector to JPY.
2. Submit a search for `watch`.
3. Expect: result-card price is rendered in JPY (¥-prefixed), not USD.
4. Switch the currency to GBP on the results page.
5. Expect: result-card price re-renders in GBP (£-prefixed), no stale ¥ values.
```

Maps to spec User Story 1 AC 4, FR-005, SC-006.

## 4. No-match and bad-input handling

```text
Submit each of:
  zzzzz          → expect: "no matches for zzzzz" page (HTTP 200).
  (empty)        → expect: "please enter a search term" page (HTTP 200).
  (spaces only)  → expect: same as empty.
  !!!!!          → expect: no-matches page (HTTP 200).
  <201-char string> → expect: query-too-long error page (HTTP 400).
  あ              → expect: no-matches page (HTTP 200), no 5xx.
```

Maps to spec User Story 2 ACs, FR-006, FR-007, FR-008, FR-009, SC-005.

## 5. Reachability from non-home pages

```text
1. Navigate to /cart.
2. Use the header search input to submit `watch`.
3. Expect: results page renders identically to a home-page-initiated search.
4. Navigate to a /product/{id} detail page.
5. Use the header search input to submit `dryer`.
6. Expect: results show the Hairdryer; the previously-viewed product gets no preferential ordering.
```

Maps to spec User Story 3 ACs, FR-011.

## 6. Backend unit test

Run from the repo root or from `src/productcatalogservice/`:

```sh
cd src/productcatalogservice
go test ./...
```

Expectations:

- **`TestSearchProducts`** (existing) — still passes after the name-only narrowing. (The fixture's matches are name-based.)
- **New name-only assertion** — added in this feature; must pass. Asserts that a fixture product whose *description* contains the query but whose *name* does not is **excluded** from results.

## 7. Constraint invariants — `git diff` on the merged PR

Spec SC-008 asserts zero infra/CI changes. Verify with:

```sh
git diff origin/attendee/tim-love...HEAD --name-only | \
  grep -E '^(kubernetes-manifests|kustomize|helm-chart|istio-manifests|terraform|\.github/workflows|cloudbuild\.yaml)' \
  && echo "FAIL: infra files modified" || echo "PASS: no infra files modified"
```

The command should print `PASS`. If it prints `FAIL`, the implementation has violated a spec constraint and the change should not be merged.

Sanity-check the env vars too — no new `env:` entries should appear in any service Deployment manifest:

```sh
git diff origin/attendee/tim-love...HEAD -- kubernetes-manifests/ kustomize/ helm-chart/
```

Expected output: empty.

## 8. Performance sanity check (optional, manual)

Spec SC-003 budgets the median search→render at under 500 ms. With a 10-product catalog and an in-memory substring filter, this should be effortless. If you have `curl` and the storefront's public IP:

```sh
time curl -s -o /dev/null "http://${STOREFRONT_IP}/search?q=watch"
```

Run 10 times, take the median. Should be well under 500 ms — typical observed: ~50-150 ms dominated by template render and the per-result `Convert` RPC.

## 9. Storefront error-rate check (optional, ops)

Spec SC-004 asserts zero new 5xx classes introduced. If the cluster has logging:

```sh
# adjust for your logging stack — illustrative only
kubectl logs deployment/frontend --since=1h | jq -c 'select(.severity=="ERROR")' | wc -l
```

Compare the count for the hour before and after rollout. Should be unchanged (modulo unrelated noise).

## Troubleshooting quickies

| Symptom | Likely cause | Where to look |
|---|---|---|
| Results page returns 500 | `SearchProducts` is failing or unreachable | `kubectl logs deployment/productcatalogservice`; check the service is up. |
| All searches return zero | The name-only narrowing dropped too much, or the case-folding got broken | `product_catalog.go` — confirm `strings.ToLower` is on both sides of `Contains`. |
| Prices render in USD regardless of currency | The handler is forgetting to call `convertCurrency` | `searchHandler` in `handlers.go` — diff against `homeHandler`. |
| Header search input missing on `/cart` or `/product/{id}` | The form was added to a wrong template (e.g. `home.html` instead of `header.html`) | `templates/header.html` — confirm the form is inside the `{{ define "header" }}` block. |
