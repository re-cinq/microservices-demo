/*
 * Browser-side product search.
 *
 * Feature: 001-browser-product-search
 * Spec:    /specs/001-browser-product-search/spec.md
 * Plan:    /specs/001-browser-product-search/plan.md
 * Contract: /specs/001-browser-product-search/contracts/ui-contract.md
 *
 * What this script does:
 *   - On every `input` event, narrows the visible product grid to cards
 *     whose data-name contains the trimmed/lowercased query (substring match).
 *   - Shows a "no results" message when nothing matches.
 *   - Pressing Escape (or using the native <input type="search"> clear control)
 *     resets the filter.
 *
 * What this script does NOT do (by design, per spec Out-of-Scope):
 *   - No network calls.
 *   - No localStorage, sessionStorage, cookies, URL writes, or history writes.
 *   - No reordering of products.
 *   - No reading of any product field other than `name`.
 */
(function () {
  "use strict";

  // ---- Pure function export (testable from DevTools console) -------------
  // Contract §2.3 — exposed on window.__productSearch for console assertions
  // and so the per-event handler below uses the exact same logic.
  function matches(dataName, query) {
    var normalized = String(query == null ? "" : query).trim().toLowerCase();
    if (normalized === "") return true;
    // dataName may be undefined if a card is missing its data-name attribute;
    // treat that as "no name to match against".
    return String(dataName == null ? "" : dataName).toLowerCase().indexOf(normalized) !== -1;
  }

  window.__productSearch = { matches: matches };

  // ---- DOM wiring --------------------------------------------------------
  function init() {
    var input = document.getElementById("product-search");
    if (!input) {
      // This page doesn't have the search box (e.g. cart, checkout). Safe no-op.
      return;
    }
    var noResults = document.getElementById("product-search-no-results");
    var cards = Array.prototype.slice.call(
      document.querySelectorAll(".hot-product-card")
    );

    // Cache each card's lowercased name once, on init, so the per-keystroke
    // hot path is a single .indexOf() — no per-keystroke toLowerCase on the
    // card side. The query is still trimmed+lowercased per keystroke inside
    // `matches`, which is fine (it's one short string).
    cards.forEach(function (card) {
      card._searchName = String(card.dataset.name || "").toLowerCase();
    });

    function applyFilter() {
      var query = input.value;
      var normalized = String(query == null ? "" : query).trim().toLowerCase();
      var visibleCount = 0;

      if (normalized === "") {
        // Empty (or whitespace-only) query: show everything.
        for (var i = 0; i < cards.length; i++) {
          cards[i].hidden = false;
        }
        visibleCount = cards.length;
      } else {
        for (var j = 0; j < cards.length; j++) {
          var card = cards[j];
          var hit = card._searchName.indexOf(normalized) !== -1;
          card.hidden = !hit;
          if (hit) visibleCount += 1;
        }
      }

      if (noResults) {
        // Show the "no results" message only when there's a non-empty query
        // AND zero cards survived the filter. The query stays in the box
        // (we never mutate input.value from inside the listener).
        noResults.hidden = !(normalized !== "" && visibleCount === 0);
      }
    }

    // US1 — filter as the shopper types. `input` covers typing, pasting,
    // IME compose end, and the native clear "×" inside <input type="search">.
    input.addEventListener("input", applyFilter);

    // US2 — Escape clears the box. Some browsers (notably Safari) don't bind
    // Escape to clear by default; this makes the behaviour reliable everywhere.
    // We re-fire `input` so the existing US1 listener does the actual reset.
    input.addEventListener("keydown", function (event) {
      if (event.key === "Escape" && input.value !== "") {
        input.value = "";
        input.dispatchEvent(new Event("input", { bubbles: true }));
        event.preventDefault();
      }
    });

    // US3 — filter resets on reload. There's nothing to *do* here: we
    // deliberately do not read query state from localStorage / sessionStorage
    // / cookies / URL on init. The search box starts empty (the template
    // renders no `value`), so the page begins unfiltered. Reloading produces
    // a fresh script execution and a fresh empty input.
    //
    // Audit hook for T008: this file MUST NOT reference any of the following.
    // Grep this file for: localStorage | sessionStorage | document.cookie |
    // window.location = | history.pushState | history.replaceState
    // None of those should appear outside this comment block.
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    // Script was loaded with `defer`, so by the time it runs the DOM is
    // parsed. This branch handles the edge case where the script is loaded
    // later by some other path.
    init();
  }
})();
