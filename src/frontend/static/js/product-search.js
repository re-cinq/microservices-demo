// Product search: client-side, in-browser filter of the product cards already
// rendered on the home page. No backend request, no page reload. Filters the
// `.hot-product-card` elements by their product name (case-insensitive,
// whitespace-trimmed substring match) and shows an inline message when nothing
// matches. See specs/product-search/.
(function () {
  'use strict';

  function init() {
    var input = document.getElementById('product-search-input');
    if (!input) {
      return; // search box not on this page
    }

    var noResults = document.getElementById('product-search-no-results');

    // Cache each product card together with its lower-cased name, read once
    // from the name already in the DOM.
    var cards = Array.prototype.slice
      .call(document.querySelectorAll('.hot-product-card'))
      .map(function (card) {
        var nameEl = card.querySelector('.hot-product-card-name');
        var name = nameEl ? nameEl.textContent.trim().toLowerCase() : '';
        return { el: card, name: name };
      });

    function applyFilter() {
      var query = input.value.trim().toLowerCase();
      var visibleCount = 0;

      cards.forEach(function (card) {
        // Empty query => show everything (in original order, since we only
        // toggle visibility and never reorder/remove nodes).
        var matches = query === '' || card.name.indexOf(query) !== -1;
        card.el.classList.toggle('d-none', !matches);
        if (matches) {
          visibleCount++;
        }
      });

      // Inline "no results" message: only when a non-empty query matches nothing.
      if (noResults) {
        var showMessage = query !== '' && visibleCount === 0;
        noResults.classList.toggle('d-none', !showMessage);
      }
    }

    input.addEventListener('input', applyFilter);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
