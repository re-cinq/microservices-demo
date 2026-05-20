# Spec: Apply a valid promo code at checkout

**Source:** Jira AIP-150 (Story) · AIP-96 (Epic)
**Branch:** feature/002-promo-codes-at-checkout

---

## Overview

A shopper who has received a discount code from a marketing campaign can enter it on the checkout page and see the discount reflected in their order total before they pay. Only one code can be active on an order at a time; applying a new code replaces any previously applied one.

---

## Functional Requirements

1. A promo code input field and apply action are visible on the checkout page.
2. When a valid promo code is submitted, the discount is applied to the order total.
3. An inline confirmation is shown when a valid code is accepted.
4. Only one promo code can be active per order; submitting a new valid code replaces the previous one and updates the total.
5. When no promo code is entered, the checkout page looks and behaves identically to today — no new friction is introduced.
6. Valid promo codes are a small hardcoded set held in memory; no external datastore is used.

---

## Technical Constraints (hard — from AIP-96)

- Use only the services that already exist in this repo. Do not add new services.
- Do not introduce any new datastore (no database, no cache, no search engine). Work in memory.
- Match the language and patterns of the service you change. The frontend and the product catalogue service are written in Go.
- Do not change infrastructure, deployment manifests, or CI configuration. The change must ship through the existing pipeline unmodified.

---

## User Scenarios

### Scenario 1 — Happy path: valid code applied
**Given** I am on the checkout page  
**When** I enter a valid promo code and click apply  
**Then** the discount is applied to my order total and an inline confirmation tells me the code has been accepted

### Scenario 2 — Replace an existing code
**Given** a valid promo code is already applied to my order  
**When** I enter a different valid promo code and click apply  
**Then** the previous code is replaced, the new discount is applied, and the order total updates accordingly

### Scenario 3 — No code entered (non-regression)
**Given** I am on the checkout page and have not entered a promo code  
**When** I proceed through checkout  
**Then** the experience is identical to today — no new fields, prompts, or friction are introduced for shoppers who do not use a promo code

---

## Success Criteria

- A shopper with a valid code can complete a discounted checkout end-to-end.
- The displayed order total reflects the discount before the shopper submits payment.
- Replacing one valid code with another updates the total correctly.
- The checkout flow for shoppers without a promo code is unchanged.
- No new service, datastore, or infrastructure component is introduced.
- All changes are in Go and follow the patterns of the modified service.
