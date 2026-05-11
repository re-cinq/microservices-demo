/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

(function () {
  'use strict';

  document.addEventListener('DOMContentLoaded', function () {
    var searchInput = document.getElementById('search-input');
    if (!searchInput) {
      return;
    }
    var cards = Array.prototype.slice.call(document.querySelectorAll('.hot-product-card'));
    var emptyState = document.getElementById('search-empty-state');
    var clearBtn = document.getElementById('search-clear');

    function applyFilter() {
      var query = searchInput.value.trim().toLowerCase();
      var matchCount = 0;
      for (var i = 0; i < cards.length; i++) {
        var card = cards[i];
        var name = (card.dataset.productName || '').toLowerCase();
        var visible = name.indexOf(query) !== -1;
        card.style.display = visible ? '' : 'none';
        if (visible) {
          matchCount++;
        }
      }
      if (emptyState) {
        if (query !== '' && matchCount === 0) {
          emptyState.removeAttribute('hidden');
        } else {
          emptyState.setAttribute('hidden', '');
        }
      }
      if (clearBtn) {
        if (searchInput.value.length > 0) {
          clearBtn.removeAttribute('hidden');
        } else {
          clearBtn.setAttribute('hidden', '');
        }
      }
    }

    function clearSearch() {
      searchInput.value = '';
      searchInput.dispatchEvent(new Event('input', { bubbles: true }));
      searchInput.focus();
    }

    searchInput.addEventListener('input', applyFilter);

    if (clearBtn) {
      clearBtn.addEventListener('click', clearSearch);
    }

    searchInput.addEventListener('keydown', function (event) {
      if (event.key === 'Escape') {
        event.preventDefault();
        clearSearch();
      }
    });
  });
})();
