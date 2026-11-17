(function () {
  var root = document.documentElement;
  var group = document.getElementById('prefs-theme');
  var MODES = ['light', 'sepia', 'dark'];

  function effectiveTheme() {
    var attr = root.getAttribute('data-theme');
    if (attr) return attr;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function render() {
    if (!group) return;
    var current = effectiveTheme();
    Array.prototype.forEach.call(group.querySelectorAll('button[data-value]'), function (b) {
      var on = b.getAttribute('data-value') === current;
      b.classList.toggle('active', on);
      b.setAttribute('aria-pressed', String(on));
    });
  }

  var saved = localStorage.getItem('theme');
  if (MODES.indexOf(saved) !== -1) {
    root.setAttribute('data-theme', saved);
  }
  render();

  if (!group) return;
  group.addEventListener('click', function (e) {
    var btn = e.target.closest('button[data-value]');
    if (!btn) return;
    var next = btn.getAttribute('data-value');
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

  /* Sitelinks search box: when Google (or anyone) lands on /read/?q=terms,
     pre-open the palette with that query so the results are immediately
     visible — wired to the WebSite SearchAction JSON-LD. */
  var q = new URLSearchParams(window.location.search).get('q');
  if (q) {
    open();
    loadIndex().then(function () {
      input.value = q;
      renderResults(q);
      activeIndex = 0;
      input.focus();
    });
  }

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

/* Collapsible TOC: every section (part) can be collapsed/expanded from its
   title, and every topic (chapter row) can be hidden with its own toggle.
   Hidden topics are restored with the "+N" chip that appears in the part
   title. Both states persist per part/topic id in localStorage. */
(function () {
  var PARTS_KEY = 'tocPartsCollapsed';
  var CHAPTERS_KEY = 'tocChaptersCollapsed';

  function loadMap(key) {
    try { return JSON.parse(localStorage.getItem(key) || '{}') || {}; }
    catch (e) { return {}; }
  }
  function saveMap(key, map) {
    try { localStorage.setItem(key, JSON.stringify(map)); } catch (e) {}
  }

  var partsMap = loadMap(PARTS_KEY);
  var chaptersMap = loadMap(CHAPTERS_KEY);
  var refreshers = [];

  document.querySelectorAll('.toc-part').forEach(function (part) {
    var title = part.querySelector(':scope > .toc-part-title');
    var list = part.querySelector(':scope > ul');
    if (!title || !list) return;
    var id = part.getAttribute('data-part') || '';

    var btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'toc-part-toggle';
    btn.textContent = '\u25BE';
    btn.setAttribute('aria-label', 'Toggle section');
    title.appendChild(btn);

    var chip = null;

    function refresh() {
      var collapsed = part.classList.contains('collapsed');
      btn.textContent = collapsed ? '\u25B8' : '\u25BE';
      btn.setAttribute('aria-expanded', String(!collapsed));
      var hidden = list.querySelectorAll('li.collapsed').length;
      if (hidden > 0) {
        if (!chip) {
          chip = document.createElement('button');
          chip.type = 'button';
          chip.className = 'toc-hidden-chip';
          chip.title = 'Show hidden topics';
          title.appendChild(chip);
          chip.addEventListener('click', function (e) {
            e.preventDefault();
            e.stopPropagation();
            list.querySelectorAll('li.collapsed').forEach(function (li) {
              var a = li.querySelector('a');
              var cid = a && a.getAttribute('data-chapter-id');
              if (cid) delete chaptersMap[cid];
              li.classList.remove('collapsed');
            });
            saveMap(CHAPTERS_KEY, chaptersMap);
            refresh();
          });
        }
        chip.textContent = '+' + hidden;
      } else if (chip) {
        chip.remove();
        chip = null;
      }
    }

    btn.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      part.classList.toggle('collapsed');
      partsMap[id] = part.classList.contains('collapsed');
      saveMap(PARTS_KEY, partsMap);
      refresh();
    });

    part.classList.toggle('collapsed', !!partsMap[id]);
    part._tocRefresh = refresh;
    refreshers.push(refresh);
  });

  document.querySelectorAll('.toc-part li').forEach(function (li) {
    var a = li.querySelector('a');
    if (!a) return;
    var id = a.getAttribute('data-chapter-id');
    var btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'toc-item-toggle';
    btn.textContent = '\u2212';
    btn.setAttribute('aria-label', 'Toggle topic');
    btn.title = 'Collapse topic';
    li.appendChild(btn);

    if (id && chaptersMap[id]) li.classList.add('collapsed');

    btn.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      var hidden = !li.classList.contains('collapsed');
      li.classList.toggle('collapsed', hidden);
      if (id) {
        if (hidden) chaptersMap[id] = true;
        else delete chaptersMap[id];
        saveMap(CHAPTERS_KEY, chaptersMap);
      }
      var part = li.closest('.toc-part');
      if (part && part._tocRefresh) part._tocRefresh();
    });
  });

  refreshers.forEach(function (fn) { fn(); });
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
  var root = document.documentElement;
  var panel = document.getElementById('prefs-panel');
  var toggle = document.getElementById('prefs-toggle');
  if (!panel || !toggle) return;

  var DEFAULTS = { font: 'sans', size: 16, leading: 1.7, measure: 72, bold: false, justify: false };
  var FONTS = {
    sans: 'var(--sans)',
    serif: '"Iowan Old Style", "Palatino Linotype", Palatino, "Book Antiqua", Georgia, "Liberation Serif", serif',
    mono: 'var(--mono)'
  };

  var els = {
    font: document.getElementById('prefs-font'),
    size: document.getElementById('prefs-size'),
    sizeOut: document.getElementById('prefs-size-out'),
    leading: document.getElementById('prefs-leading'),
    leadingOut: document.getElementById('prefs-leading-out'),
    measure: document.getElementById('prefs-measure'),
    measureOut: document.getElementById('prefs-measure-out'),
    bold: document.getElementById('prefs-bold'),
    justify: document.getElementById('prefs-justify'),
    reset: document.getElementById('prefs-reset')
  };

  function clamp(n, lo, hi) { return Math.min(hi, Math.max(lo, n)); }

  function load() {
    var saved = {};
    try { saved = JSON.parse(localStorage.getItem('readerPrefs') || '{}') || {}; } catch (e) {}
    var p = Object.assign({}, DEFAULTS, saved);
    // One-time migration from the old A-/A+ scale and the Justify button.
    if (saved.size == null) {
      var legacyScale = parseFloat(localStorage.getItem('readerFontScale'));
      if (!isNaN(legacyScale)) p.size = Math.round(14 * legacyScale);
    }
    if (saved.justify == null && localStorage.getItem('justifyText') === '1') p.justify = true;
    p.size = clamp(p.size, 13, 30);
    p.leading = clamp(p.leading, 1.3, 2.2);
    p.measure = clamp(p.measure, 50, 95);
    return p;
  }

  var prefs = load();

  function save() {
    try { localStorage.setItem('readerPrefs', JSON.stringify(prefs)); } catch (e) {}
  }

  function apply() {
    root.style.setProperty('--reader-font', FONTS[prefs.font] || FONTS.sans);
    root.style.setProperty('--reader-font-size', prefs.size + 'px');
    root.style.setProperty('--reader-leading', String(prefs.leading));
    root.style.setProperty('--reader-measure', prefs.measure + 'ch');
    root.style.setProperty('--reader-weight', prefs.bold ? '600' : '400');
    root.classList.toggle('justify-text', !!prefs.justify);
    sync();
  }

  function sync() {
    if (els.size) { els.size.value = prefs.size; els.sizeOut.textContent = prefs.size + ' px'; }
    if (els.leading) { els.leading.value = prefs.leading; els.leadingOut.textContent = Number(prefs.leading).toFixed(2); }
    if (els.measure) { els.measure.value = prefs.measure; els.measureOut.textContent = prefs.measure + ' ch'; }
    if (els.bold) els.bold.checked = !!prefs.bold;
    if (els.justify) els.justify.checked = !!prefs.justify;
    if (els.font) Array.prototype.forEach.call(els.font.querySelectorAll('button[data-value]'), function (b) {
      var on = b.getAttribute('data-value') === prefs.font;
      b.classList.toggle('active', on);
      b.setAttribute('aria-pressed', String(on));
    });
  }

  function changed() { apply(); save(); }
  apply();

  if (els.size) els.size.addEventListener('input', function () { prefs.size = parseInt(this.value, 10) || DEFAULTS.size; changed(); });
  if (els.leading) els.leading.addEventListener('input', function () { prefs.leading = parseFloat(this.value) || DEFAULTS.leading; changed(); });
  if (els.measure) els.measure.addEventListener('input', function () { prefs.measure = parseInt(this.value, 10) || DEFAULTS.measure; changed(); });
  if (els.bold) els.bold.addEventListener('change', function () { prefs.bold = this.checked; changed(); });
  if (els.justify) els.justify.addEventListener('change', function () { prefs.justify = this.checked; changed(); });
  if (els.font) els.font.addEventListener('click', function (e) {
    var btn = e.target.closest('button[data-value]');
    if (!btn) return;
    prefs.font = btn.getAttribute('data-value');
    changed();
  });
  if (els.reset) els.reset.addEventListener('click', function () {
    prefs = Object.assign({}, DEFAULTS);
    changed();
  });

  function setOpen(open) {
    panel.hidden = !open;
    toggle.setAttribute('aria-expanded', String(open));
    if (open) panel.focus({ preventScroll: true });
  }

  toggle.addEventListener('click', function (e) {
    e.stopPropagation();
    setOpen(panel.hidden);
  });
  var close = document.getElementById('prefs-close');
  if (close) close.addEventListener('click', function () { setOpen(false); toggle.focus(); });
  document.addEventListener('click', function (e) {
    if (panel.hidden) return;
    if (panel.contains(e.target) || toggle.contains(e.target)) return;
    setOpen(false);
  });
  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && !panel.hidden) { setOpen(false); toggle.focus(); }
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

(function () {
  var AUTHOR = 'Flores Salazar, O.';
  var YEAR = '2025';
  var BOOK_TITLE = 'Networking with Go, Made Easy';

  function esc(s) {
    var d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
  }

  function nearestHeadingElement() {
    var sel = window.getSelection();
    if (!sel.rangeCount) return null;
    var node = sel.anchorNode;
    var el = node && node.nodeType === 3 ? node.parentElement : node;
    if (!el) return null;
    if (el.matches && el.matches('h2, h3')) return el;
    var parentHeading = el.closest('h2, h3');
    if (parentHeading) return parentHeading;
    var section = el.closest('section[id^="section-"]');
    if (!section) return null;
    var headings = section.querySelectorAll('h2, h3');
    var nearest = null;
    for (var i = 0; i < headings.length; i++) {
      if (headings[i].compareDocumentPosition(el) & Node.DOCUMENT_POSITION_FOLLOWING) {
        nearest = headings[i];
      }
    }
    return nearest;
  }

  document.addEventListener('copy', function (e) {
    var text = window.getSelection().toString().trim();
    if (!text) return;

    var node = window.getSelection().anchorNode;
    var el = node && node.nodeType === 3 ? node.parentElement : node;
    if (!el || !el.closest('.reader-content')) return;

    var h1 = document.querySelector('.reader-content section h1');
    if (!h1) return;
    var chapter = h1.textContent.trim();

    var headingEl = nearestHeadingElement();
    var heading = headingEl ? headingEl.textContent.trim() : '';
    var link = document.querySelector('link[rel="canonical"]');
    var url = link ? link.getAttribute('href') : '';
    if (url && headingEl && headingEl.id) {
      url += '#' + headingEl.id;
    }

    var source = AUTHOR + ' (' + YEAR + '). ' + chapter;
    if (heading) source += ' (' + heading + ')';
    source += '. In ' + BOOK_TITLE + '. ' + url;

    var plain = '\u201c' + text + '\u201d\n\n\u2014 ' + source;

    var eChapter = esc(chapter);
    var eHeading = esc(heading);
    var eURL = esc(url);
    var html = '<blockquote style="margin:0 0 6px 0;padding:8px 16px;border-left:3px solid #999;color:var(--pdf-muted,#666);font-family:system-ui,sans-serif;font-size:10.5pt;line-height:1.6;"><p style="margin:0;">' + esc(text) + '</p></blockquote><p style="margin:0;font-family:system-ui,sans-serif;font-size:10pt;line-height:1.5;color:var(--pdf-muted,#444);">\u2014 ' + eChapter + (eHeading ? ' <span style="color:#888;">(' + eHeading + ')</span>' : '') + '. In <em>' + BOOK_TITLE + '</em>. <a href="' + eURL + '">' + eURL + '</a></p>';

    e.clipboardData.setData('text/plain', plain);
    e.clipboardData.setData('text/html', html);
    e.preventDefault();

    var toast = document.getElementById('copy-toast');
    if (!toast) {
      toast = document.createElement('div');
      toast.id = 'copy-toast';
      toast.className = 'copy-toast';
      toast.textContent = 'Copied with APA citation';
      document.body.appendChild(toast);
    }
    toast.classList.add('visible');
    clearTimeout(toast._hide);
    toast._hide = setTimeout(function () { toast.classList.remove('visible'); }, 2000);
  });
})();

/* Fixed reading-progress bar along the top edge of the page. */
(function () {
  var bar = document.getElementById('reading-progress');
  if (!bar) return;
  function update() {
    var doc = document.documentElement;
    var max = doc.scrollHeight - doc.clientHeight;
    bar.style.width = (max > 0 ? (window.scrollY / max) * 100 : 0) + '%';
  }
  window.addEventListener('scroll', update, { passive: true });
  window.addEventListener('resize', update);
  update();
})();

/* Keyboard navigation: [ previous chapter, ] next chapter. */
(function () {
  function isTyping(e) {
    var t = e.target;
    if (!t) return false;
    return /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName) || t.isContentEditable;
  }
  document.addEventListener('keydown', function (e) {
    if (e.metaKey || e.ctrlKey || e.altKey || isTyping(e)) return;
    var k = e.key;
    if (k === '[') {
      var prev = document.querySelector('.chapter-nav-prev');
      if (prev) { e.preventDefault(); prev.click(); }
    } else if (k === ']') {
      var next = document.querySelector('.chapter-nav-next');
      if (next) { e.preventDefault(); next.click(); }
    }
  });
})();

/* "Copy link to section" button appended to every h2/h3. */
(function () {
  var content = document.querySelector('.reader-content');
  if (!content) return;

  function toast(msg) {
    var t = document.getElementById('copy-toast');
    if (!t) {
      t = document.createElement('div');
      t.id = 'copy-toast';
      t.className = 'copy-toast';
      document.body.appendChild(t);
    }
    t.textContent = msg;
    t.classList.add('visible');
    clearTimeout(t._hide);
    t._hide = setTimeout(function () { t.classList.remove('visible'); }, 1600);
  }

  content.querySelectorAll('h2, h3').forEach(function (h) {
    if (!h.id) return;
    var a = document.createElement('a');
    a.className = 'heading-link';
    a.href = '#' + h.id;
    a.setAttribute('aria-label', 'Link to ' + h.textContent.trim());
    a.title = 'Copy link to section';
    a.textContent = '\u00B6';
    h.appendChild(a);
  });

  content.addEventListener('click', function (e) {
    var a = e.target && e.target.closest ? e.target.closest('.heading-link') : null;
    if (!a) return;
    e.preventDefault();
    var url = window.location.href.split('#')[0] + a.getAttribute('href');
    function done() {
      toast('Section link copied');
      window.location.hash = a.getAttribute('href');
    }
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(url).then(done, function () { done(); });
    } else {
      done();
    }
  });
})();

/* Estimated reading time for chapter pages (~200 wpm). */
(function () {
  var el = document.getElementById('reading-time');
  if (!el) return;
  var content = document.querySelector('.reader-content');
  if (!content) return;
  var clone = content.cloneNode(true);
  clone.querySelectorAll('pre, code, script, style, .breadcrumb, .share-bar, .chapter-nav, .continue-reading, .heading-link').forEach(function (n) {
    if (n.parentNode) n.parentNode.removeChild(n);
  });
  var words = (clone.textContent || '').trim().split(/\s+/).filter(Boolean).length;
  el.textContent = '\u00b7 ' + Math.max(1, Math.round(words / 200)) + ' min read';
})();

/* Per-page notes panel, autosaved to localStorage. */
(function () {
  var panel = document.getElementById('notes-panel');
  var scrim = document.getElementById('notes-scrim');
  var toggle = document.getElementById('notes-toggle');
  var closeBtn = document.getElementById('notes-close');
  var text = document.getElementById('notes-text');
  var status = document.getElementById('notes-status');
  if (!panel || !text) return;

  var pageKey = document.body.getAttribute('data-chapter-id') ||
                document.body.getAttribute('data-page-type') ||
                'page';
  var KEY = 'nwgNotes';

  function loadAll() {
    try { return JSON.parse(localStorage.getItem(KEY) || '{}') || {}; }
    catch (e) { return {}; }
  }
  function save() {
    var all = loadAll();
    all[pageKey] = text.value;
    localStorage.setItem(KEY, JSON.stringify(all));
    if (status) {
      status.textContent = 'Saved';
      clearTimeout(save._t);
      save._t = setTimeout(function () { status.textContent = ''; }, 1500);
    }
  }

  text.value = loadAll()[pageKey] || '';
  text.addEventListener('input', save);

  function open() {
    panel.hidden = false;
    scrim.hidden = false;
    toggle.classList.add('active');
    text.focus();
  }
  function closePanel() {
    panel.hidden = true;
    scrim.hidden = true;
    toggle.classList.remove('active');
  }
  if (toggle) toggle.addEventListener('click', function () { panel.hidden ? open() : closePanel(); });
  if (closeBtn) closeBtn.addEventListener('click', closePanel);
  if (scrim) scrim.addEventListener('click', closePanel);
  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && !panel.hidden) closePanel();
  });
})();

/* PDF button: point at the actual release asset, not the releases page and
   not window.print(). Asset names are versioned (networking-with-go-v1.1.0.pdf),
   so the URL comes from the GitHub API; the static href is the releases page
   and stays usable if the API call fails or is rate-limited. Cached for six
   hours because every reader page asks for it. */
(function () {
  var btn = document.getElementById('print-pdf');
  if (!btn) return;
  var KEY = 'latestPdf';
  var TTL = 6 * 60 * 60 * 1000;

  function use(url, version) {
    if (!url) return;
    btn.href = url;
    if (version) btn.title = 'Download the book PDF (' + version + ')';
  }

  try {
    var cached = JSON.parse(localStorage.getItem(KEY) || 'null');
    if (cached && cached.url && Date.now() - cached.at < TTL) {
      use(cached.url, cached.version);
      return;
    }
  } catch (e) {}

  fetch('https://api.github.com/repos/sazardev/networking-with-go/releases/latest')
    .then(function (res) { return res.ok ? res.json() : null; })
    .then(function (data) {
      if (!data || !data.assets) return;
      var asset = null;
      for (var i = 0; i < data.assets.length; i++) {
        if (/\.pdf$/i.test(data.assets[i].name)) { asset = data.assets[i]; break; }
      }
      if (!asset) return;
      use(asset.browser_download_url, data.tag_name);
      try {
        localStorage.setItem(KEY, JSON.stringify({
          url: asset.browser_download_url, version: data.tag_name, at: Date.now()
        }));
      } catch (e) {}
    })
    .catch(function () {});
})();

/* Copy canonical link from the chapter share bar. */
(function () {
  var btn = document.getElementById('share-copy');
  if (!btn) return;
  var canonical = document.querySelector('link[rel="canonical"]');
  var url = canonical ? canonical.getAttribute('href') : window.location.href;
  btn.addEventListener('click', function () {
    function done() {
      btn.textContent = 'Copied!';
      setTimeout(function () { btn.textContent = 'Copy link'; }, 1500);
    }
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(url).then(done, function () { done(); });
    } else {
      done();
    }
  });
})();

/* Newsletter signup (reader index). Endpoint is configured per-site via
   window.SUBSCRIBE_ENDPOINT; until then the form directs to the RSS feed. */
(function () {
  var form = document.getElementById('subscribe-form');
  var status = document.getElementById('subscribe-status');
  var ENDPOINT = window.SUBSCRIBE_ENDPOINT || '';
  if (!form || !status) return;

  function show(msg, ok) {
    status.textContent = msg;
    status.hidden = false;
    status.style.color = ok ? '' : '#c0392b';
  }

  form.addEventListener('submit', function (e) {
    e.preventDefault();
    var email = document.getElementById('subscribe-email').value.trim();
    if (!email) { show('Please enter your email.', false); return; }
    if (!ENDPOINT) {
      show('Newsletter coming soon — subscribe to the RSS feed for updates.', false);
      return;
    }
    var btn = form.querySelector('button[type="submit"]');
    btn.disabled = true;
    btn.textContent = '...';
    fetch(ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email })
    }).then(function (res) {
      if (!res.ok) throw new Error('bad status');
      show('You are subscribed! Check your inbox to confirm.', true);
      form.reset();
    }).catch(function () {
      show('Could not subscribe right now — please try again or use the RSS feed.', false);
    }).then(function () {
      btn.disabled = false;
      btn.textContent = 'Subscribe';
    });
  });
})();
