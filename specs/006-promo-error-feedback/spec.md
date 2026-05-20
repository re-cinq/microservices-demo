# Feature Specification: Promo Code Invalid-Code Error Feedback

**Feature Branch**: `006-promo-error-feedback`

**Created**: 2026-05-20

**Status**: Draft

**Source**: Jira AIP-151 (parent epic: AIP-96 — Promo codes at checkout)

## Clarifications

### Session 2026-05-20

- Q: When a shopper has an invalid-code error showing and has not corrected it, can they still submit the checkout form and place their order? → A: Non-blocking — order can be placed with error showing; the invalid code is simply not applied, full price charged.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Unrecognised Code Shows Error (Priority: P1)

A shopper on the checkout page enters a promo code that is not in the system — whether they mistyped it, used an expired code, or simply guessed. They submit the code and immediately see a clear, inline error message telling them the code was not applied. The order total remains unchanged. The shopper knows exactly what happened and is not misled into thinking they received a discount.

**Why this priority**: This is the primary deliverable of AIP-151. Without an explicit error state, shoppers who enter invalid codes silently proceed at full price while believing they have a discount — a direct trust and revenue risk.

**Independent Test**: On the checkout page, enter a code that is not SAVE10 or PROMO20 (e.g., "BADCODE") and press Apply. An error message appears inline near the promo code field. The order total is identical to what it was before the code was entered.

**Acceptance Scenarios**:

1. **Given** I am on the checkout page, **When** I enter an unrecognised code and press Apply, **Then** a clear inline error message is displayed and the order total is unchanged.
2. **Given** I am on the checkout page, **When** I enter a blank string and press Apply, **Then** the promo field behaves gracefully — either showing an appropriate error or ignoring the submission — and the order total is unchanged.
3. **Given** I have received an invalid-code error, **When** I look at the order total, **Then** it matches the pre-code total exactly with no discount applied.

---

### User Story 2 - Recovery from Error by Entering a Valid Code (Priority: P2)

A shopper who received an invalid-code error corrects their mistake, types a valid code (SAVE10 or PROMO20), and re-applies. The error clears, the discount is applied, and the total updates to reflect the saving.

**Why this priority**: Recovery is essential for usability — an error state that traps the shopper with no path forward would increase checkout abandonment. This story validates the two-state handoff between error and success.

**Independent Test**: After triggering an invalid-code error with "BADCODE", type "SAVE10" and press Apply. The error disappears, a success confirmation appears, and the order total reflects a 10% discount.

**Acceptance Scenarios**:

1. **Given** I have received an invalid-code error, **When** I enter a valid code and press Apply, **Then** the error message clears and the discount is applied.
2. **Given** a valid code is now applied after an earlier error, **When** I view the order total, **Then** it shows the correctly discounted amount.

---

### User Story 3 - Proceeding Without a Code After an Error (Priority: P3)

A shopper who received an invalid-code error decides not to use a promo code at all. They clear the field and proceed through checkout. The error message is gone, the order total is the full original price, and checkout completes normally.

**Why this priority**: Ensures the error state does not persist into the rest of the checkout flow, which would confuse the shopper and potentially block the purchase.

**Independent Test**: After triggering an invalid-code error, clear the promo code field and click the checkout-submit button. Checkout proceeds without any error banner present and the confirmation shows the original full price.

**Acceptance Scenarios**:

1. **Given** I have received an invalid-code error, **When** I clear the promo code field and proceed to place the order, **Then** no error message persists on the page.
2. **Given** checkout completes without a promo code, **When** I reach the order confirmation page, **Then** the total shown is the original full price and no promo section appears.

---

### Edge Cases

- What happens when the shopper submits a code that differs from SAVE10/PROMO20 only by case (e.g., "save10")? *(Assumption: codes are matched case-insensitively; documented in Assumptions.)*
- What happens if the shopper submits whitespace-padded input (e.g., " SAVE10 ")? *(Assumption: input is trimmed before validation.)*
- What happens if the promo code field is empty when Apply is pressed? *(An appropriate error or no-op is shown; total is unchanged.)*
- What if the page is refreshed after an error? *(The error state need not persist across page loads; the field resets.)*

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The checkout page MUST display an inline error message when the shopper submits a promo code that is not recognised.
- **FR-002**: The error message MUST be proximate to the promo code input field so the shopper clearly associates it with their code entry.
- **FR-003**: The order total MUST remain unchanged when an invalid code is submitted.
- **FR-004**: When the shopper subsequently enters a valid code and applies it, the error message MUST clear and the discount MUST be applied.
- **FR-005**: When the shopper clears the promo code field and proceeds to checkout without a code, no error message MUST persist.
- **FR-006**: The promo code validation MUST use only the existing hardcoded set (SAVE10 = 10%, PROMO20 = 20%) — no new data service, database, cache, or search engine may be introduced.
- **FR-007**: No new backend service may be created; only existing services may be extended.
- **FR-008**: The checkout page layout and all existing content MUST remain visually unchanged except for the addition of the error feedback.
- **FR-010**: The error state is non-blocking — a shopper MAY still submit the checkout form while an invalid-code error is displayed; the invalid code is ignored and the order is charged at the full original price.
- **FR-009**: Code comparison MUST be case-insensitive and trim surrounding whitespace before validation.

### Key Entities

- **Promo Code**: A short alphanumeric string entered by the shopper. Valid values: SAVE10 (10% off), PROMO20 (20% off). Any other value is invalid.
- **Error State**: The UI condition shown when an invalid code is submitted. Contains a human-readable message. Cleared when a valid code is applied or the field is emptied and checkout proceeds.
- **Order Total**: The monetary amount the shopper will be charged. Must equal the pre-code total when an invalid code is in effect.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of invalid code submissions result in a visible error message proximate to the promo code field, with no change to the order total.
- **SC-002**: A shopper who corrects an invalid code to a valid one sees the error clear and the discounted total in a single Apply action — zero additional steps required.
- **SC-003**: A shopper who clears the promo field after an error and proceeds reaches the order confirmation page with no error element visible — verified by end-to-end walkthrough.
- **SC-004**: No regression to existing checkout behaviour for shoppers who never enter a promo code — verified by comparison with the pre-feature checkout flow.
- **SC-005**: No new service, datastore, or infrastructure component is introduced — verified by diff review.

## Assumptions

- Promo code validation is handled in the existing frontend or checkout service — no new service endpoint is needed.
- The valid code set (SAVE10, PROMO20) is already implemented as part of AIP-150 and available in the codebase on `attendee/nejal-patel`.
- Code matching is case-insensitive and trims whitespace; this is a reasonable usability default and avoids a frustrating UX where "save10" fails silently.
- The error message copy ("Promo code not recognised" or similar) can be determined at implementation time; exact wording is not specified in the story.
- No changes to infrastructure, deployment manifests, or CI configuration are permitted. The change must ship through the existing pipeline unmodified.
- Out of scope: campaign-management UI, new datastores, server-side promo code validation beyond the hardcoded set.
- The live deployment target is `https://nejal-patel.training.gcp.re-cinq.com`.
