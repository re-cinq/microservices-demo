# Tasks: Promo Code Invalid-Code Error Feedback (AIP-151)

**Input**: Design documents from `specs/006-promo-error-feedback/`

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no shared dependencies)
- **[Story]**: Which user story (US1, US2, US3) — absent in Setup/Foundational/Polish phases

---

## Phase 1: Setup (Port AIP-150 Foundation)

**Purpose**: Cherry-pick the promo-code infrastructure from kate-payne's AIP-150 implementation into this branch. AIP-151 cannot be implemented without this foundation.

- [x] T001 Cherry-pick source-code files from commit `f14efd7a` (origin/attendee/kate-payne): `src/frontend/main.go`, `src/frontend/handlers.go`, `src/frontend/templates/cart.html`, `src/checkoutservice/main.go` — skip `.specify/feature.json`, `CLAUDE.md`, and `specs/002-*` artefacts

**Checkpoint**: `go build ./...` compiles cleanly; promo input appears on the cart page; SAVE10/PROMO20 apply a discount.

---

## Phase 2: Foundational (Blocking Prerequisite)

**Purpose**: Add the `cookiePromoError` constant before any story implementation can reference it.

**⚠️ CRITICAL**: Phase 3 and 4 handler changes cannot compile without this constant.

- [x] T002 Add `cookiePromoError = cookiePrefix + "promo-error"` constant in `src/frontend/main.go` (alongside `cookiePromoCode`)

**Checkpoint**: Project compiles with the new constant.

---

## Phase 3: User Story 1 — Unrecognised Code Shows Error (Priority: P1) 🎯 MVP

**Goal**: When a shopper submits a promo code that is not SAVE10 or PROMO20, they see a clear inline error message and the order total is unchanged.

**Independent Test**: Enter "BADCODE" on the cart page → press Apply → inline error "Promo code not recognised" is visible below the input; total is unchanged.

### Implementation for User Story 1

- [x] T003 [US1] Update `applyPromoHandler` in `src/frontend/handlers.go`: normalise input with `strings.ToUpper(strings.TrimSpace(...))` before map lookup; on valid code set `cookiePromoCode` and expire `cookiePromoError`; on invalid code clear `cookiePromoCode` and set flash cookie `cookiePromoError=1` with `MaxAge: 60`
- [x] T004 [US1] Update `viewCartHandler` in `src/frontend/handlers.go`: after building `templateData`, read `cookiePromoError`; if present and `Value == "1"`, set `templateData["promo_error"] = true` and immediately expire the cookie with `MaxAge: -1`
- [x] T005 [US1] Add inline error display block in `src/frontend/templates/cart.html` immediately below the closing `</form>` tag of the promo input section: `{{ if $.promo_error }}<small class="text-danger mt-1 d-block">Promo code not recognised — please check the code and try again.</small>{{ end }}`

**Checkpoint**: US1 fully functional — BADCODE shows error; SAVE10/PROMO20 show success; total unchanged on error.

---

## Phase 4: User Story 2 — Recovery from Error (Priority: P2)

**Goal**: After an invalid-code error, entering a valid code clears the error and applies the discount.

**Independent Test**: Trigger error with "BADCODE" → type "SAVE10" → press Apply → error gone, 10% discount shown.

### Implementation for User Story 2

- [x] T006 [US2] Verify T003 already clears `cookiePromoError` on the valid-code path (the `MaxAge: -1` set in the valid branch of `applyPromoHandler`); if not yet handled, add the expiry cookie to `src/frontend/handlers.go` valid path
- [x] T007 [P] [US2] Verify case-insensitive recovery: confirm `strings.ToUpper` in T003 means "save10" and "SAVE10" both succeed and clear any prior error — no code change needed if T003 is correct, but verify by inspection

**Checkpoint**: Entering a valid code after an error shows the success confirmation with discount; error block absent.

---

## Phase 5: User Story 3 — Proceed Without Code After Error (Priority: P3)

**Goal**: After dismissing an error (clearing the field and placing the order), no error persists to the confirmation page, and the full original price is charged.

**Independent Test**: Trigger error → do NOT correct the code → click Place Order → confirmation page shows no error and the full original price.

### Implementation for User Story 3

- [x] T008 [US3] Add expiry of `cookiePromoError` in `placeOrderHandler` in `src/frontend/handlers.go` (alongside the existing `cookiePromoCode` expiry) to ensure the flash cookie cannot survive past order placement
- [x] T009 [US3] Verify `src/frontend/templates/cart.html` — the error block is conditioned on `$.promo_error` only (set from the flash cookie, one-shot) and is absent from the order confirmation template — confirm by reading `src/frontend/templates/order.html` (no change needed if confirmation template is separate)

**Checkpoint**: Full checkout completes at original price with no error visible; confirmation page clean.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [x] T010 [P] Run `go build ./...` in `src/frontend/` and `src/checkoutservice/` to confirm both services compile with all changes
- [ ] T011 Run end-to-end walkthrough from `specs/006-promo-error-feedback/quickstart.md` against `https://nejal-patel.training.gcp.re-cinq.com`
- [x] T012 [P] Verify no-promo regression: add item to cart, proceed through checkout without entering any code — confirm experience is identical to pre-feature behaviour

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start here
- **Phase 2 (Foundational)**: Depends on Phase 1 — BLOCKS all user stories
- **Phase 3 (US1)**: Depends on Phase 2
- **Phase 4 (US2)**: Depends on Phase 3 (error path must exist to test recovery)
- **Phase 5 (US3)**: Depends on Phase 3 (error path must exist to test persistence)
- **Phase 6 (Polish)**: Depends on all story phases

### User Story Dependencies

- **US1 (P1)**: Foundational complete → can start immediately
- **US2 (P2)**: US1 complete (error path must exist for recovery test)
- **US3 (P3)**: US1 complete (error path must exist for persistence test); US2 independent

### Parallel Opportunities

- T006 and T007 (US2) can run in parallel once US1 is complete
- T010 and T012 (Polish) can run in parallel

---

## Parallel Example: User Story 1

```bash
# These three tasks can proceed as a group once Phase 2 is complete:
Task T003: Update applyPromoHandler (handlers.go)
Task T004: Update viewCartHandler (handlers.go) — depends on T003 for the constant
Task T005: Add error block to cart.html — independent of T003/T004, different file
```

---

## Implementation Strategy

### MVP (User Story 1 Only)

1. Complete Phase 1 (cherry-pick AIP-150)
2. Complete Phase 2 (add constant)
3. Complete Phase 3 (T003, T004, T005)
4. **STOP and validate**: BADCODE shows error, SAVE10 applies discount
5. Deploy to `https://nejal-patel.training.gcp.re-cinq.com`

### Incremental Delivery

1. Phase 1 + 2 → foundation ready
2. Phase 3 → error feedback live (MVP — AIP-151 core delivered)
3. Phase 4 → recovery path verified
4. Phase 5 → non-blocking checkout verified
5. Phase 6 → build + e2e walkthrough

---

## Notes

- T001 cherry-pick must skip spec/config artefacts to avoid overwriting `CLAUDE.md` and `.specify/feature.json`
- Flash cookie MaxAge=60s is intentionally short; the one-shot expiry in `viewCartHandler` makes the cookie self-cleaning
- No new services, routes (beyond `/cart/promo` from AIP-150), or infrastructure changes
