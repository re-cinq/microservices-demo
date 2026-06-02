(function () {
  var input = document.getElementById('product-search');
  var noResults = document.getElementById('no-results-message');

  if (!input || !noResults) return;

  function applyFilter() {
    var query = input.value.trim().toLowerCase();
    var cards = document.querySelectorAll('.hot-product-card[data-name]');
    var visible = 0;

    cards.forEach(function (card) {
      var name = (card.dataset.name || '').toLowerCase();
      var match = query === '' || name.indexOf(query) !== -1;
      card.style.display = match ? '' : 'none';
      if (match) visible++;
    });

    noResults.style.display = (query !== '' && visible === 0) ? 'block' : 'none';
  }

  input.addEventListener('input', applyFilter);

  input.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
      input.value = '';
      applyFilter();
    }
  });
}());
