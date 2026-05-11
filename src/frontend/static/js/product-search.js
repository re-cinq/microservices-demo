// Copyright 2026 Online Boutique cohort
//
// Client-side product-name search for the home page.
// Filters the already-rendered .hot-product-card elements in the DOM by
// case-insensitive substring match against the card's data-name attribute.
//
// Implements (see specs/003-product-search/):
//   - FR-001 .. FR-010, FR-006a, FR-006b   (see spec.md)
//   - The DOM/JS behaviour contract        (see contracts/ui-contract.md)
//
// Hard rules:
//   - No fetch / XHR / WebSocket of any kind (SC-003, FR-007).
//   - No reads beyond the card's `data-name` attribute (C-009).
//   - No URL / history / storage writes; no persistence across page navigations
//     (FR-010 — the input is re-rendered empty on each page load).
//
// T010 (US2) verification note:
//   The browser-native <input type="search"> × clear button fires the same
//   `input` event we listen for on Chrome, Edge, Safari, and Firefox as of
//   the targeted evergreen set. No `search` event listener fallback is needed.
//   Re-verify if the browser-support matrix changes.

(function () {
  'use strict';

  document.addEventListener('DOMContentLoaded', function () {
    var input = document.getElementById('product-search-input');
    if (!input) {
      // Not on the home page — this script is a no-op everywhere else.
      return;
    }

    var status = document.getElementById('product-search-status');
    var emptyEl = document.querySelector('.product-search-empty');
    var cards = Array.prototype.slice.call(
      document.querySelectorAll('.hot-product-card')
    );

    if (cards.length === 0 || !status || !emptyEl) {
      // Defensive: if the DOM contract isn't met, do nothing.
      return;
    }

    function setStatus(text) {
      if (status.textContent !== text) {
        status.textContent = text;
      }
    }

    function applyFilter() {
      var q = input.value.trim().toLowerCase();

      if (q === '') {
        // FR-005: empty/whitespace query → all cards visible, no announcement,
        // no empty-state message.
        for (var i = 0; i < cards.length; i++) {
          cards[i].hidden = false;
        }
        emptyEl.hidden = true;
        setStatus('');
        return;
      }

      var visibleCount = 0;
      for (var j = 0; j < cards.length; j++) {
        var card = cards[j];
        var name = (card.dataset.name || '').toLowerCase();
        var match = name.indexOf(q) !== -1;
        card.hidden = !match;
        if (match) {
          visibleCount++;
        }
      }

      // FR-006: visible empty-state message only when query is non-empty
      // AND no card matched.
      emptyEl.hidden = visibleCount !== 0;

      // FR-006a: polite live-region announcement.
      if (visibleCount === 0) {
        setStatus('No products match your search.');
      } else if (visibleCount === 1) {
        setStatus('1 product matches');
      } else {
        setStatus(visibleCount + ' products match');
      }
    }

    input.addEventListener('input', applyFilter);
  });
})();
