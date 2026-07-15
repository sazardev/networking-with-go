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
  var base = window.READER_BASE || '';

  function loadIndex() {
    if (index) return Promise.resolve(index);
    return fetch(base + 'search-index.json')
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
      return '<a class="palette-result' + (i === 0 ? ' active' : '') + '" data-href="' + base + e.page + '#' + e.id + '">' +
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

(function () {
  var root = document.documentElement;
  var mqMobile = window.matchMedia('(max-width: 900px)');
  var sidebarToggle = document.getElementById('sidebar-toggle');
  var scrim = document.querySelector('.sidebar-scrim');

  function applyState() {
    if (mqMobile.matches) {
      root.classList.remove('sidebar-collapsed');
    } else {
      root.classList.remove('sidebar-open');
      root.classList.toggle('sidebar-collapsed', localStorage.getItem('sidebarCollapsed') === '1');
    }
  }

  applyState();
  mqMobile.addEventListener('change', applyState);

  if (sidebarToggle) {
    sidebarToggle.addEventListener('click', function () {
      if (mqMobile.matches) {
        root.classList.toggle('sidebar-open');
      } else {
        var collapsed = !root.classList.contains('sidebar-collapsed');
        root.classList.toggle('sidebar-collapsed', collapsed);
        localStorage.setItem('sidebarCollapsed', collapsed ? '1' : '0');
      }
    });
  }

  if (scrim) {
    scrim.addEventListener('click', function () { root.classList.remove('sidebar-open'); });
  }
})();

(function () {
  var toggle = document.getElementById('focus-toggle');
  if (!toggle) return;
  var root = document.documentElement;

  function render() {
    var on = root.classList.contains('focus-mode');
    toggle.classList.toggle('active', on);
    toggle.textContent = on ? 'Exit Focus' : 'Focus';
  }

  if (localStorage.getItem('focusMode') === '1') root.classList.add('focus-mode');
  render();

  toggle.addEventListener('click', function () {
    var on = !root.classList.contains('focus-mode');
    root.classList.toggle('focus-mode', on);
    localStorage.setItem('focusMode', on ? '1' : '0');
    render();
  });
})();

(function () {
  var dec = document.getElementById('font-dec');
  var inc = document.getElementById('font-inc');
  if (!dec || !inc) return;
  var root = document.documentElement;
  var MIN = 0.8, MAX = 1.6, STEP = 0.1;

  function apply(scale) {
    scale = Math.min(MAX, Math.max(MIN, scale));
    root.style.setProperty('--reader-font-scale', scale.toFixed(2));
    localStorage.setItem('readerFontScale', scale.toFixed(2));
  }

  var saved = parseFloat(localStorage.getItem('readerFontScale'));
  apply(isNaN(saved) ? 1 : saved);

  dec.addEventListener('click', function () {
    apply((parseFloat(localStorage.getItem('readerFontScale')) || 1) - STEP);
  });
  inc.addEventListener('click', function () {
    apply((parseFloat(localStorage.getItem('readerFontScale')) || 1) + STEP);
  });
})();

(function () {
  var STORAGE_KEY = 'readChapters';
  var LAST_READ_KEY = 'lastRead';

  function escapeText(s) {
    var div = document.createElement('div');
    div.textContent = s;
    return div.innerHTML;
  }

  function getReadMap() {
    try {
      return JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
    } catch (e) {
      return {};
    }
  }

  function applyReadMarks() {
    var map = getReadMap();
    document.querySelectorAll('[data-chapter-id]').forEach(function (a) {
      a.classList.toggle('read', !!map[a.getAttribute('data-chapter-id')]);
    });
  }

  function markRead(chapterId) {
    if (!chapterId) return;
    var map = getReadMap();
    if (map[chapterId]) return;
    map[chapterId] = Date.now();
    localStorage.setItem(STORAGE_KEY, JSON.stringify(map));
    applyReadMarks();
  }

  applyReadMarks();

  var pageType = document.body.getAttribute('data-page-type');
  var chapterId = document.body.getAttribute('data-chapter-id');

  if (pageType === 'chapter' && chapterId) {
    localStorage.setItem(LAST_READ_KEY, JSON.stringify({
      chapterId: chapterId,
      href: window.location.href,
      title: document.title.replace(' — Networking with Go, Made Easy', ''),
      ts: Date.now()
    }));

    var sentinel = document.querySelector('.read-sentinel');
    if (sentinel && 'IntersectionObserver' in window) {
      var observer = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            markRead(chapterId);
            observer.disconnect();
          }
        });
      }, { rootMargin: '0px 0px -10% 0px' });
      observer.observe(sentinel);
    }
  }

  if (pageType === 'index') {
    var last = null;
    try { last = JSON.parse(localStorage.getItem(LAST_READ_KEY) || 'null'); } catch (e) {}
    if (last && last.href) {
      var h1 = document.querySelector('.reader-content h1');
      if (h1) {
        var banner = document.createElement('div');
        banner.className = 'continue-reading';
        banner.innerHTML =
          '<div><span class="cr-label">Continue reading</span>' +
          '<span class="cr-title">' + escapeText(last.title || 'your last chapter') + '</span></div>' +
          '<a class="btn-continue" href="' + last.href + '">Resume &rarr;</a>';
        h1.insertAdjacentElement('afterend', banner);
      }
    }
  }
})();
