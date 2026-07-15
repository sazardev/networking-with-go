// Command buildreader turns docs/**/*.mdx into the web reader under
// website/read/. It reuses go-pretty-pdf's own "serve" subcommand (the same
// tool that composes the PDF) to get HTML with matching fonts and already-
// styled custom components (<Axiom>, <Warning>, <DeepDive>) for free, then
// slices that HTML into one page per chapter (143 pages, grouped under one
// index page per part) so each chapter is its own indexable URL with real
// SEO metadata, fixes a heading-id collision bug (heading ids are only
// unique per-chapter, not per-book), builds a flat search index for the
// reader's command palette, and (re)generates website/sitemap.xml.
//
// Run from the repository root:
//
//	go run website/tools/buildreader/main.go
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	siteBase      = "https://sazardev.github.io/networking-with-go"
	coverImageURL = "https://raw.githubusercontent.com/sazardev/networking-with-go/master/assets/cover.jpg"
)

type part struct {
	Slug  string
	Title string
	Desc  string
	Dir   string
}

var parts = []part{
	{"networking-foundations", "Networking Foundations", "OSI & TCP/IP models, IP addressing and subnetting, ports, sockets, protocols, firewalls, NAT, and VPNs — theory only, no code, told through real history.", "docs/part1"},
	{"go-fundamentals", "Go From Scratch", "A full Go language crash course before the networking code starts, capped by a flagship chapter on goroutines, channels, and concurrency.", "docs/go-fundamentals"},
	{"core-go-networking", "Core Go Networking", "TCP, UDP, HTTP, JSON, WebSockets, DNS, proxies, chat apps, file transfer, NAT traversal, context, and TLS — theory paired with runnable Go every time.", "docs/part2"},
	{"advanced-specialized", "Advanced & Specialized", "gRPC, WebRTC, MQTT/IoT, SDN, NFV, packet inspection, custom protocols, mesh networks, and real-time networking for games and blockchain.", "docs/advanced"},
	{"cybersecurity-ethical-hacking", "Cybersecurity & Ethical Hacking", "Scanning, sniffing, vulnerability assessment, IDS/IPS, SIEM, malware analysis, red team vs blue team, honeypots, and PKI — for authorized, ethical use.", "docs/part3"},
	{"production-apis-architecture", "Production APIs & Architecture", "REST and gRPC APIs with Gin and Fiber, PostgreSQL and SQLite, Docker and Google Cloud, Clean/Hexagonal Architecture, and Domain-Driven Design.", "docs/part-apis"},
}

var (
	styleRe      = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	sectionOpen  = regexp.MustCompile(`<section id="(section-[0-9.]+)">`)
	headingRe    = regexp.MustCompile(`(?s)<h([1-3]) id="([^"]+)">(.*?)</h[1-3]>`)
	tagStripRe   = regexp.MustCompile(`<[^>]+>`)
	slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)
)

// linkCtx carries the relative-path prefixes needed at a given page depth.
// Only two depths exist: read/ itself (depth0) and read/<part>/ (depth1).
type linkCtx struct {
	toRead string // prefix to reach website/read/
	toSite string // prefix to reach website/ (the main landing page)
}

var (
	ctxDepth0 = linkCtx{toRead: "", toSite: "../"}
	ctxDepth1 = linkCtx{toRead: "../", toSite: "../../"}
)

type headingInfo struct {
	ID      string
	Text    string
	Level   int
	Excerpt string
}

type chapterEntry struct {
	SectionID string
	Title     string
	Slug      string
	BodyHTML  string
	Excerpt   string
}

type partResult struct {
	meta     part
	chapters []chapterEntry
}

type flatRef struct {
	PartSlug     string
	ChapterSlug  string
	ChapterTitle string
}

type searchEntry struct {
	ID      string `json:"id"`
	Page    string `json:"page"`
	Part    string `json:"part"`
	Chapter string `json:"chapter"`
	Heading string `json:"heading"`
	Level   int    `json:"level"`
	Excerpt string `json:"excerpt"`
}

func main() {
	log.SetFlags(0)

	outDir := "website/read"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("create %s: %v", outDir, err)
	}

	var results []partResult
	var search []searchEntry
	var theme string

	for _, p := range parts {
		log.Printf("composing %s ...", p.Dir)
		raw, err := composePartHTML(p.Dir)
		if err != nil {
			log.Fatalf("%s: %v", p.Dir, err)
		}

		if theme == "" {
			if m := styleRe.FindStringSubmatch(raw); m != nil {
				theme = m[1]
			}
		}

		contentStart := strings.Index(raw, `<section id="section-`)
		if contentStart == -1 {
			log.Fatalf("%s: no <section> found in composed HTML", p.Dir)
		}
		content := raw[contentStart:]
		lastClose := strings.LastIndex(content, "</section>")
		if lastClose == -1 {
			log.Fatalf("%s: no closing </section> found", p.Dir)
		}
		content = content[:lastClose+len("</section>")]

		opens := sectionOpen.FindAllStringSubmatchIndex(content, -1)
		if len(opens) == 0 {
			log.Fatalf("%s: no chapter sections found", p.Dir)
		}

		var chapters []chapterEntry
		for i, om := range opens {
			start := om[0]
			end := len(content)
			if i+1 < len(opens) {
				end = opens[i+1][0]
			}
			sectionHTML := content[start:end]
			sectionID := content[om[2]:om[3]]

			openTagEnd := strings.Index(sectionHTML, ">") + 1
			innerBody := sectionHTML[openTagEnd:]

			newBody, chapterTitle, headings := rewriteHeadings(innerBody, sectionID)
			chapterSlug := fmt.Sprintf("%02d-%s", i+1, slugify(chapterTitle))
			bodyHTML := fmt.Sprintf(`<section id="%s">`, sectionID) + newBody

			excerpt := ""
			if len(headings) > 0 {
				excerpt = headings[0].Excerpt
			}

			chapters = append(chapters, chapterEntry{
				SectionID: sectionID,
				Title:     chapterTitle,
				Slug:      chapterSlug,
				BodyHTML:  bodyHTML,
				Excerpt:   excerpt,
			})

			for _, h := range headings {
				search = append(search, searchEntry{
					ID:      h.ID,
					Page:    p.Slug + "/" + chapterSlug + ".html",
					Part:    p.Title,
					Chapter: chapterTitle,
					Heading: h.Text,
					Level:   h.Level,
					Excerpt: h.Excerpt,
				})
			}
		}

		results = append(results, partResult{meta: p, chapters: chapters})
		log.Printf("  -> %d chapters", len(chapters))
	}

	sidebarDepth0 := buildSidebar(results, ctxDepth0)
	sidebarDepth1 := buildSidebar(results, ctxDepth1)

	var flat []flatRef
	for _, r := range results {
		for _, ch := range r.chapters {
			flat = append(flat, flatRef{PartSlug: r.meta.Slug, ChapterSlug: ch.Slug, ChapterTitle: ch.Title})
		}
	}

	totalChapters := 0
	globalPos := 0
	for _, r := range results {
		partDir := filepath.Join(outDir, r.meta.Slug)
		if err := os.MkdirAll(partDir, 0o755); err != nil {
			log.Fatalf("create %s: %v", partDir, err)
		}

		writeFile(filepath.Join(partDir, "index.html"), renderPartIndexPage(r, sidebarDepth1, theme))

		for _, ch := range r.chapters {
			globalPos++
			var prev, next *flatRef
			if globalPos-2 >= 0 {
				prev = &flat[globalPos-2]
			}
			if globalPos < len(flat) {
				next = &flat[globalPos]
			}
			page := renderChapterPage(ch, r.meta, prev, next, globalPos, sidebarDepth1, theme)
			writeFile(filepath.Join(partDir, ch.Slug+".html"), page)
		}
		totalChapters += len(r.chapters)
	}

	writeFile(filepath.Join(outDir, "index.html"), renderIndexPage(results, sidebarDepth0, theme))

	data, err := json.MarshalIndent(search, "", "  ")
	if err != nil {
		log.Fatalf("marshal search index: %v", err)
	}
	writeFile(filepath.Join(outDir, "search-index.json"), string(data))

	writeSitemap(results)

	log.Printf("done: %d parts, %d chapters, %d search entries", len(results), totalChapters, len(search))
}

func composePartHTML(sourceDir string) (string, error) {
	port, err := freePort()
	if err != nil {
		return "", fmt.Errorf("find free port: %w", err)
	}

	cmd := exec.Command("pretty-pdf", "serve", "--source", sourceDir, "--port", strconv.Itoa(port), "--quiet")
	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start pretty-pdf serve: %w", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	url := fmt.Sprintf("http://127.0.0.1:%d/", port)
	if err := waitReady(url, 15*time.Second); err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	return string(body), nil
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func waitReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %s", url)
}

// rewriteHeadings fixes the fact that go-pretty-pdf slugs headings from text
// alone (e.g. every chapter's "Frequently Asked Questions" h2 gets the same
// id), making anchors ambiguous document-wide. H2/H3 ids get prefixed with
// their chapter's own unique section id; the H1 (chapter title) loses its id
// entirely since the enclosing <section id="section-X.Y.Z"> is already the
// correct, unique anchor for "jump to this chapter".
func rewriteHeadings(body, sectionID string) (newBody string, chapterTitle string, headings []headingInfo) {
	matches := headingRe.FindAllStringSubmatchIndex(body, -1)
	var out strings.Builder
	last := 0
	for i, m := range matches {
		out.WriteString(body[last:m[0]])

		level, _ := strconv.Atoi(body[m[2]:m[3]])
		origID := body[m[4]:m[5]]
		inner := body[m[6]:m[7]]
		plainTitle := plainText(inner)

		excerptEnd := len(body)
		if i+1 < len(matches) {
			excerptEnd = matches[i+1][0]
		}
		excerpt := plainText(body[m[1]:excerptEnd])
		if len([]rune(excerpt)) > 240 {
			excerpt = string([]rune(excerpt)[:240]) + "…"
		}

		var id string
		if level == 1 {
			id = sectionID
			chapterTitle = plainTitle
			out.WriteString(fmt.Sprintf("<h1>%s</h1>", inner))
		} else {
			id = sectionID + "--" + origID
			out.WriteString(fmt.Sprintf(`<h%d id="%s">%s</h%d>`, level, id, inner, level))
		}

		headings = append(headings, headingInfo{ID: id, Text: plainTitle, Level: level, Excerpt: excerpt})
		last = m[1]
	}
	out.WriteString(body[last:])
	return out.String(), chapterTitle, headings
}

func plainText(s string) string {
	s = tagStripRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugNonAlnum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return strings.TrimSpace(string(r[:n])) + "…"
}

func buildSidebar(results []partResult, ctx linkCtx) string {
	var b strings.Builder
	b.WriteString(`<nav class="reader-sidebar" aria-label="Table of contents">`)
	for _, r := range results {
		fmt.Fprintf(&b, `<div class="toc-part" data-part="%s"><div class="toc-part-title"><a href="%s%s/index.html">%s</a></div><ul>`,
			r.meta.Slug, ctx.toRead, r.meta.Slug, html.EscapeString(r.meta.Title))
		for _, ch := range r.chapters {
			fmt.Fprintf(&b, `<li><a href="%s%s/%s.html">%s</a></li>`, ctx.toRead, r.meta.Slug, ch.Slug, html.EscapeString(ch.Title))
		}
		b.WriteString(`</ul></div>`)
	}
	b.WriteString(`</nav>`)
	return b.String()
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
	log.Printf("wrote %s", path)
}

const favicon = `data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 64 64'%3E%3Crect width='64' height='64' fill='%23000'/%3E%3Ccircle cx='16' cy='46' r='5' fill='none' stroke='%23fff' stroke-width='3'/%3E%3Ccircle cx='48' cy='46' r='5' fill='none' stroke='%23fff' stroke-width='3'/%3E%3Ccircle cx='32' cy='16' r='5' fill='none' stroke='%23fff' stroke-width='3'/%3E%3Cline x1='19' y1='43' x2='29' y2='20' stroke='%23fff' stroke-width='3'/%3E%3Cline x1='45' y1='43' x2='35' y2='20' stroke='%23fff' stroke-width='3'/%3E%3Cline x1='21' y1='46' x2='43' y2='46' stroke='%23fff' stroke-width='3'/%3E%3C/svg%3E`

// ldNode is a loosely-typed JSON-LD node, marshaled via encoding/json so
// string escaping (quotes, unicode) is always correct.
type ldNode map[string]interface{}

func marshalLD(graph []ldNode) string {
	doc := map[string]interface{}{
		"@context": "https://schema.org",
		"@graph":   graph,
	}
	b, err := json.Marshal(doc)
	if err != nil {
		log.Fatalf("marshal json-ld: %v", err)
	}
	safe := strings.ReplaceAll(string(b), "</", "<\\/")
	return "<script type=\"application/ld+json\">\n" + safe + "\n</script>\n"
}

func chapterJSONLD(ch chapterEntry, p part, position int, canonicalURL string) string {
	graph := []ldNode{
		{
			"@type": "Chapter",
			"name":  ch.Title,
			"isPartOf": ldNode{
				"@type": "Book",
				"name":  "Networking with Go, Made Easy",
				"url":   siteBase + "/",
			},
			"position":            position,
			"url":                 canonicalURL,
			"author":              ldNode{"@type": "Person", "name": "Omar Flores Salazar", "url": "https://github.com/sazardev"},
			"license":             "https://creativecommons.org/licenses/by-nc/4.0/",
			"isAccessibleForFree": true,
			"description":         ch.Excerpt,
			"inLanguage":          "en",
		},
		{
			"@type": "BreadcrumbList",
			"itemListElement": []ldNode{
				{"@type": "ListItem", "position": 1, "name": "Read Online", "item": siteBase + "/read/index.html"},
				{"@type": "ListItem", "position": 2, "name": p.Title, "item": siteBase + "/read/" + p.Slug + "/index.html"},
				{"@type": "ListItem", "position": 3, "name": ch.Title, "item": canonicalURL},
			},
		},
	}
	return marshalLD(graph)
}

func partIndexJSONLD(p part, chapters []chapterEntry) string {
	items := make([]ldNode, len(chapters))
	for i, ch := range chapters {
		items[i] = ldNode{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     ch.Title,
			"url":      siteBase + "/read/" + p.Slug + "/" + ch.Slug + ".html",
		}
	}
	graph := []ldNode{
		{
			"@type": "BreadcrumbList",
			"itemListElement": []ldNode{
				{"@type": "ListItem", "position": 1, "name": "Read Online", "item": siteBase + "/read/index.html"},
				{"@type": "ListItem", "position": 2, "name": p.Title, "item": siteBase + "/read/" + p.Slug + "/index.html"},
			},
		},
		{
			"@type":           "ItemList",
			"name":            p.Title + " — Chapters",
			"itemListElement": items,
		},
	}
	return marshalLD(graph)
}

type pageMeta struct {
	Title       string
	Description string
	Canonical   string
	OGType      string
	JSONLD      string
}

func pageHead(m pageMeta, ctx linkCtx, theme string) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString(`<meta charset="utf-8">` + "\n")
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">` + "\n")
	fullTitle := m.Title + " — Networking with Go, Made Easy"
	fmt.Fprintf(&b, "<title>%s</title>\n", html.EscapeString(fullTitle))
	fmt.Fprintf(&b, `<meta name="description" content="%s">`+"\n", html.EscapeString(m.Description))
	b.WriteString(`<meta name="robots" content="index, follow, max-image-preview:large">` + "\n")
	fmt.Fprintf(&b, `<link rel="canonical" href="%s">`+"\n", m.Canonical)
	fmt.Fprintf(&b, `<link rel="icon" type="image/svg+xml" href="%s">`+"\n", favicon)

	fmt.Fprintf(&b, `<meta property="og:type" content="%s">`+"\n", m.OGType)
	b.WriteString(`<meta property="og:site_name" content="Networking with Go, Made Easy">` + "\n")
	fmt.Fprintf(&b, `<meta property="og:title" content="%s">`+"\n", html.EscapeString(m.Title))
	fmt.Fprintf(&b, `<meta property="og:description" content="%s">`+"\n", html.EscapeString(m.Description))
	fmt.Fprintf(&b, `<meta property="og:url" content="%s">`+"\n", m.Canonical)
	fmt.Fprintf(&b, `<meta property="og:image" content="%s">`+"\n", coverImageURL)
	b.WriteString(`<meta name="twitter:card" content="summary_large_image">` + "\n")

	if m.JSONLD != "" {
		b.WriteString(m.JSONLD)
	}

	b.WriteString("<style>\n" + theme + "\n</style>\n")
	fmt.Fprintf(&b, `<link rel="stylesheet" href="%sreader.css">`+"\n", ctx.toRead)
	b.WriteString("</head>\n")
	return b.String()
}

func pageHeader(ctx linkCtx) string {
	return fmt.Sprintf(`<header class="reader-top">
  <a class="brand" href="%sindex.html">net<span>/</span>go<span>.</span>book</a>
  <div class="reader-top-actions">
    <a href="%sindex.html">All Parts</a>
    <button id="search-open" type="button">Search <kbd>Ctrl K</kbd></button>
    <button id="theme-toggle" type="button" aria-pressed="false">Dark mode</button>
  </div>
</header>
`, ctx.toSite, ctx.toRead)
}

func breadcrumbHTML(ctx linkCtx, p part, ch *chapterEntry) string {
	var b strings.Builder
	b.WriteString(`<nav class="breadcrumb" aria-label="Breadcrumb">`)
	fmt.Fprintf(&b, `<a href="%sindex.html">Read Online</a>`, ctx.toRead)
	b.WriteString(` <span class="crumb-sep">&rsaquo;</span> `)
	fmt.Fprintf(&b, `<a href="%s%s/index.html">%s</a>`, ctx.toRead, p.Slug, html.EscapeString(p.Title))
	if ch != nil {
		b.WriteString(` <span class="crumb-sep">&rsaquo;</span> `)
		fmt.Fprintf(&b, `<span class="crumb-current">%s</span>`, html.EscapeString(ch.Title))
	}
	b.WriteString(`</nav>` + "\n")
	return b.String()
}

func chapterNavHTML(ctx linkCtx, prev, next *flatRef) string {
	var b strings.Builder
	b.WriteString(`<nav class="chapter-nav">`)
	if prev != nil {
		fmt.Fprintf(&b, `<a class="chapter-nav-prev" href="%s%s/%s.html"><span class="chapter-nav-label">&larr; Previous</span><span class="chapter-nav-title">%s</span></a>`,
			ctx.toRead, prev.PartSlug, prev.ChapterSlug, html.EscapeString(prev.ChapterTitle))
	} else {
		b.WriteString(`<span></span>`)
	}
	if next != nil {
		fmt.Fprintf(&b, `<a class="chapter-nav-next" href="%s%s/%s.html"><span class="chapter-nav-label">Next &rarr;</span><span class="chapter-nav-title">%s</span></a>`,
			ctx.toRead, next.PartSlug, next.ChapterSlug, html.EscapeString(next.ChapterTitle))
	}
	b.WriteString(`</nav>`)
	return b.String()
}

func paletteMarkup(ctx linkCtx) string {
	return fmt.Sprintf(`<div id="palette" class="palette" hidden>
  <div class="palette-box">
    <input id="palette-input" type="text" placeholder="Search the book..." autocomplete="off">
    <div id="palette-results"></div>
  </div>
</div>
<script>window.READER_BASE = "%s";</script>
<script src="%sreader.js"></script>
`, ctx.toRead, ctx.toRead)
}

func renderChapterPage(ch chapterEntry, p part, prev, next *flatRef, globalPos int, sidebar, theme string) string {
	canonical := fmt.Sprintf("%s/read/%s/%s.html", siteBase, p.Slug, ch.Slug)
	meta := pageMeta{
		Title:       ch.Title,
		Description: truncateRunes(ch.Excerpt, 155),
		Canonical:   canonical,
		OGType:      "article",
		JSONLD:      chapterJSONLD(ch, p, globalPos, canonical),
	}
	var b strings.Builder
	b.WriteString(pageHead(meta, ctxDepth1, theme))
	fmt.Fprintf(&b, "<body data-current-part=\"%s\">\n", p.Slug)
	b.WriteString(pageHeader(ctxDepth1))
	b.WriteString(`<div class="reader-layout">` + "\n")
	b.WriteString(sidebar)
	b.WriteString(`<main class="reader-content">` + "\n")
	b.WriteString(breadcrumbHTML(ctxDepth1, p, &ch))
	b.WriteString(ch.BodyHTML)
	b.WriteString("\n" + chapterNavHTML(ctxDepth1, prev, next) + "\n")
	b.WriteString("</main>\n</div>\n")
	b.WriteString(paletteMarkup(ctxDepth1))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func renderPartIndexPage(r partResult, sidebar, theme string) string {
	canonical := fmt.Sprintf("%s/read/%s/index.html", siteBase, r.meta.Slug)
	meta := pageMeta{
		Title:       r.meta.Title,
		Description: r.meta.Desc,
		Canonical:   canonical,
		OGType:      "website",
		JSONLD:      partIndexJSONLD(r.meta, r.chapters),
	}
	var b strings.Builder
	b.WriteString(pageHead(meta, ctxDepth1, theme))
	fmt.Fprintf(&b, "<body data-current-part=\"%s\">\n", r.meta.Slug)
	b.WriteString(pageHeader(ctxDepth1))
	b.WriteString(`<div class="reader-layout">` + "\n")
	b.WriteString(sidebar)
	b.WriteString(`<main class="reader-content">` + "\n")
	b.WriteString(breadcrumbHTML(ctxDepth1, r.meta, nil))
	fmt.Fprintf(&b, "<h1>%s</h1>\n<p>%s</p>\n", html.EscapeString(r.meta.Title), html.EscapeString(r.meta.Desc))
	b.WriteString(`<ol class="chapter-list">` + "\n")
	for _, ch := range r.chapters {
		fmt.Fprintf(&b, `<li><a href="%s.html"><span class="chapter-list-title">%s</span><span class="chapter-list-excerpt">%s</span></a></li>`+"\n",
			ch.Slug, html.EscapeString(ch.Title), html.EscapeString(truncateRunes(ch.Excerpt, 160)))
	}
	b.WriteString("</ol>\n</main>\n</div>\n")
	b.WriteString(paletteMarkup(ctxDepth1))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func renderIndexPage(results []partResult, sidebar, theme string) string {
	meta := pageMeta{
		Title:       "Read Online",
		Description: "Read all 143 chapters of Networking with Go, Made Easy free in your browser — same content as the PDF, same fonts, with full-text search.",
		Canonical:   siteBase + "/read/index.html",
		OGType:      "website",
	}
	var b strings.Builder
	b.WriteString(pageHead(meta, ctxDepth0, theme))
	b.WriteString("<body>\n")
	b.WriteString(pageHeader(ctxDepth0))
	b.WriteString(`<div class="reader-layout">` + "\n")
	b.WriteString(sidebar)
	b.WriteString(`<main class="reader-content">` + "\n")
	b.WriteString(`<h1>Read Networking with Go, Made Easy — online</h1>` + "\n")
	b.WriteString(`<p>All 143 chapters, six parts, right in the browser — same content as the PDF, same fonts, always in sync with the latest edit.</p>` + "\n")
	b.WriteString(`<div class="part-grid">` + "\n")
	for i, r := range results {
		fmt.Fprintf(&b, `<a class="part-card" href="%s/index.html"><span class="part-idx">%02d</span><h2>%s</h2><p>%s</p><span class="part-count">%d chapters</span></a>`+"\n",
			r.meta.Slug, i+1, html.EscapeString(r.meta.Title), html.EscapeString(r.meta.Desc), len(r.chapters))
	}
	b.WriteString("</div>\n</main>\n</div>\n")
	b.WriteString(paletteMarkup(ctxDepth0))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func writeSitemap(results []partResult) {
	type urlEntry struct{ Loc, ChangeFreq, Priority string }

	var urls []urlEntry
	urls = append(urls, urlEntry{siteBase + "/", "weekly", "1.0"})
	urls = append(urls, urlEntry{siteBase + "/read/index.html", "weekly", "0.9"})
	for _, r := range results {
		urls = append(urls, urlEntry{siteBase + "/read/" + r.meta.Slug + "/index.html", "monthly", "0.8"})
		for _, ch := range r.chapters {
			urls = append(urls, urlEntry{siteBase + "/read/" + r.meta.Slug + "/" + ch.Slug + ".html", "monthly", "0.7"})
		}
	}

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, u := range urls {
		fmt.Fprintf(&b, "  <url>\n    <loc>%s</loc>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n  </url>\n", u.Loc, u.ChangeFreq, u.Priority)
	}
	b.WriteString(`</urlset>` + "\n")

	writeFile("website/sitemap.xml", b.String())
	log.Printf("sitemap: %d URLs", len(urls))
}
