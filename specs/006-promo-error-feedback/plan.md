# Implementation Plan: Promo Code Invalid-Code Error Feedback (AIP-151)

**Branch**: `006-promo-error-feedback` | **Date**: 2026-05-20 | **Spec**: [spec.md](spec.md)

## Summary

Add inline error feedback to the checkout cart page when a shopper submits an unrecognised promo code. Builds on the promo-code foundation from AIP-150 (cherry-picked from `origin/attendee/kate-payne`). Uses a flash cookie to carry the error state across the POST→redirect→GET flow, matching the existing cookie-only pattern. No new services, datastores, or infrastructure changes.

## Technical Context

**Language/Version**: Go 1.21 (frontend service, checkoutservice)

**Primary Dependencies**: `net/http`, `gorilla/mux`, `html/template`, Bootstrap 4 (cart template)

**Storage**: In-memory `promoCodes` map (no datastore). Short-lived `shop_promo-error` flash cookie (MaxAge=60s).

**Testing**: Manual walkthrough against `https://nejal-patel.training.gcp.re-cinq.com` (no unit test harness in scope per project constraints).

**Target Platform**: GKE Autopilot, namespace `attendee-nejal-patel`

**Project Type**: Server-rendered Go web service

**Performance Goals**: No additional latency — flash cookie is set/read in-process.

**Constraints**: No new services, datastores, infra changes, or CI modifications. Match Go patterns of the services being changed.

**Scale/Scope**: Single attendee namespace, training only.

## Constitution Check

No active constitution defined for this project. No gates to evaluate.

## Project Structure

### Documentation (this feature)

```text
specs/006-promo-error-feedback/
├── plan.md              # This file
├── research.md          # Phase 0 complete
├── data-model.md        # Phase 1 complete
├── quickstart.md        # Phase 1 complete
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (repository root)

```text
src/frontend/
├── main.go              # +cookiePromoError constant
├── handlers.go          # +promoCodes map, applyPromoHandler (valid+error paths), viewCartHandler (flash cookie read + expire)
└── templates/
    └── cart.html        # +promo input form, +inline error block, +success block, +discount row

src/checkoutservice/
└── main.go              # +promoCodes map, applyPromoDiscount, gRPC metadata read for discount
```

## Implementation Steps

### Step 0 — Port AIP-150 foundation (prerequisite)

Cherry-pick source-code changes from `f14efd7a` (kate-payne's AIP-150 commit) into `006-promo-error-feedback`. Skip spec/config artefacts (`.specify/feature.json`, `CLAUDE.md`, `specs/002-*`).

Files ported:
- `src/frontend/main.go` — `cookiePromoCode` constant + `/cart/promo` route
- `src/frontend/handlers.go` — `promoCodes` map, `applyPromoDiscount`, `promoSavings`, `applyPromoHandler` (valid-code path only), `viewCartHandler` discount logic, `placeOrderHandler` gRPC metadata
- `src/frontend/templates/cart.html` — promo input form, success confirmation, discount row
- `src/checkoutservice/main.go` — `promoCodes` map, `applyPromoDiscount`, gRPC metadata read

### Step 1 — Add `cookiePromoError` constant

In `src/frontend/main.go`, add:
```go
cookiePromoError = cookiePrefix + "promo-error"
```

### Step 2 — Extend `applyPromoHandler` with error path

In `src/frontend/handlers.go`, within `applyPromoHandler`, after normalising input to uppercase:
- If code is **invalid**: set `shop_promo-error=1` cookie (MaxAge=60) in addition to clearing `shop_promo-code`.

```go
code := strings.ToUpper(strings.TrimSpace(r.FormValue("promo_code")))
if _, ok := promoCodes[code]; ok {
    // valid — set promo-code cookie, clear any error
    http.SetCookie(w, &http.Cookie{Name: cookiePromoCode, Value: code, MaxAge: cookieMaxAge})
    http.SetCookie(w, &http.Cookie{Name: cookiePromoError, Value: "", MaxAge: -1})
} else {
    // invalid — clear promo-code cookie, set flash error
    http.SetCookie(w, &http.Cookie{Name: cookiePromoCode, Value: "", MaxAge: -1})
    http.SetCookie(w, &http.Cookie{Name: cookiePromoError, Value: "1", MaxAge: 60})
}
```

**Note**: also update the `promoCodes` lookup in the existing valid-code branch to use the normalised uppercase `code`.

### Step 3 — Read and expire flash cookie in `viewCartHandler`

After building `templateData`, add:
```go
if c, err := r.Cookie(cookiePromoError); err == nil && c.Value == "1" {
    templateData["promo_error"] = true
    http.SetCookie(w, &http.Cookie{Name: cookiePromoError, Value: "", MaxAge: -1})
}
```

### Step 4 — Add inline error to `cart.html`

Immediately below the promo input form (after the closing `</form>`), add:
```html
{{ if $.promo_error }}
<small class="text-danger mt-1 d-block">
    Promo code not recognised — please check the code and try again.
</small>
{{ end }}
```

## Complexity Tracking

No constitution violations. No complexity justification required.
