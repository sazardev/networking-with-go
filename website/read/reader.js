(function () {
  var root = document.documentElement;
  var toggle = document.getElementById('theme-toggle');

  function effectiveTheme() {
    var attr = root.getAttribute('data-theme');
    if (attr) return attr;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function render() {
    var eff = effectiveTheme();
    toggle.textContent = eff === 'dark' ? 'Light mode' : 'Dark mode';
    toggle.setAttribute('aria-pressed', String(eff === 'dark'));
  }

  var saved = localStorage.getItem('theme');
  if (saved === 'light' || saved === 'dark') {
    root.setAttribute('data-theme', saved);
  }
  render();

  toggle.addEventListener('click', function () {
    var next = effectiveTheme() === 'dark' ? 'light' : 'dark';
    root.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    render();
  });
})();

(function () {
  var currentPart = document.body.getAttribute('data-current-part');
  if (currentPart) {
    var el = document.querySelector('.toc-part[data-part="' + currentPart + '"]');
    if (el) el.classList.add('current-part');
  }
})();

(function () {
  var palette = document.getElementById('palette');
  var input = document.getElementById('palette-input');
  var results = document.getElementById('palette-results');
  var openBtn = document.getElementById('search-open');
  var index = null;
  var activeIndex = -1;

  function loadIndex() {
    if (index) return Promise.resolve(index);
    return fetch('search-index.json')
      .then(function (r) { return r.json(); })
      .then(function (data) { index = data; return index; })
      .catch(function (err) {
        console.warn('Could not load search index', err);
        index = [];
        return index;
      });
  }

  function open() {
    palette.hidden = false;
    input.value = '';
    results.innerHTML = '';
    activeIndex = -1;
    loadIndex().then(function () { input.focus(); });
  }

  function close() {
    palette.hidden = true;
  }

  function score(entry, terms) {
    var haystack = (entry.heading + ' ' + entry.chapter + ' ' + entry.excerpt).toLowerCase();
    var total = 0;
    for (var i = 0; i < terms.length; i++) {
      var t = terms[i];
      if (!t) continue;
      if (haystack.indexOf(t) === -1) return -1;
      total += entry.heading.toLowerCase().indexOf(t) !== -1 ? 3 : 1;
      if (entry.level === 1) total += 1;
    }
    return total;
  }

  function renderResults(query) {
    var terms = query.toLowerCase().split(/\s+/).filter(Boolean);
    if (terms.length === 0) {
      results.innerHTML = '<div class="palette-empty">Type to search chapters and headings...</div>';
      return;
    }
    var scored = [];
    for (var i = 0; i < index.length; i++) {
      var s = score(index[i], terms);
      if (s >= 0) scored.push({ entry: index[i], s: s });
    }
    scored.sort(function (a, b) { return b.s - a.s; });
    scored = scored.slice(0, 30);

    if (scored.length === 0) {
      results.innerHTML = '<div class="palette-empty">No matches.</div>';
      return;
    }

    activeIndex = 0;
    results.innerHTML = scored.map(function (item, i) {
      var e = item.entry;
      var label = e.level === 1 ? e.chapter : e.chapter + ' › ' + e.heading;
      return '<a class="palette-result' + (i === 0 ? ' active' : '') + '" data-href="' + e.page + '#' + e.id + '">' +
        '<div class="pr-title">' + escapeHTML(label) + '</div>' +
        '<div class="pr-meta">' + escapeHTML(e.part) + '</div>' +
        '<div class="pr-excerpt">' + escapeHTML(e.excerpt) + '</div>' +
        '</a>';
    }).join('');
  }

  function escapeHTML(s) {
    var div = document.createElement('div');
    div.textContent = s;
    return div.innerHTML;
  }

  function moveActive(delta) {
    var items = results.querySelectorAll('.palette-result');
    if (items.length === 0) return;
    items[activeIndex] && items[activeIndex].classList.remove('active');
    activeIndex = (activeIndex + delta + items.length) % items.length;
    items[activeIndex].classList.add('active');
    items[activeIndex].scrollIntoView({ block: 'nearest' });
  }

  function goToActive() {
    var items = results.querySelectorAll('.palette-result');
    if (items.length === 0 || activeIndex < 0) return;
    window.location.href = items[activeIndex].getAttribute('data-href');
  }

  input.addEventListener('input', function () { renderResults(input.value); });

  results.addEventListener('click', function (e) {
    var a = e.target.closest('.palette-result');
    if (a) window.location.href = a.getAttribute('data-href');
  });

  document.addEventListener('keydown', function (e) {
    var isMod = e.metaKey || e.ctrlKey;
    if (isMod && e.key.toLowerCase() === 'k') {
      e.preventDefault();
      palette.hidden ? open() : close();
      return;
    }
    if (palette.hidden) return;
    if (e.key === 'Escape') { close(); return; }
    if (e.key === 'ArrowDown') { e.preventDefault(); moveActive(1); return; }
    if (e.key === 'ArrowUp') { e.preventDefault(); moveActive(-1); return; }
    if (e.key === 'Enter') { e.preventDefault(); goToActive(); return; }
  });

  openBtn.addEventListener('click', open);
  palette.addEventListener('click', function (e) {
    if (e.target === palette) close();
  });

  renderResults('');
})();
