// Command buildreader turns docs/**/*.mdx into the web reader under
// website/read/. It reuses go-pretty-pdf's own "serve" subcommand (the same
// tool that composes the PDF) to get HTML with matching fonts and already-
// styled custom components (<Axiom>, <Warning>, <DeepDive>) for free, then
// slices that HTML into per-part pages, fixes a heading-id collision bug
// (heading ids are only unique per-chapter, not per-book), and builds a
// flat search index for the reader's command palette.
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
	styleRe     = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	sectionOpen = regexp.MustCompile(`<section id="(section-[0-9.]+)">`)
	headingRe   = regexp.MustCompile(`(?s)<h([1-3]) id="([^"]+)">(.*?)</h[1-3]>`)
	tagStripRe  = regexp.MustCompile(`<[^>]+>`)
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
}

type partResult struct {
	meta     part
	bodyHTML string
	chapters []chapterEntry
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

		var body strings.Builder
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

			body.WriteString(fmt.Sprintf(`<section id="%s">`, sectionID))
			body.WriteString(newBody)

			chapters = append(chapters, chapterEntry{SectionID: sectionID, Title: chapterTitle})

			for _, h := range headings {
				search = append(search, searchEntry{
					ID:      h.ID,
					Page:    p.Slug + ".html",
					Part:    p.Title,
					Chapter: chapterTitle,
					Heading: h.Text,
					Level:   h.Level,
					Excerpt: h.Excerpt,
				})
			}
		}

		results = append(results, partResult{meta: p, bodyHTML: body.String(), chapters: chapters})
		log.Printf("  -> %d chapters", len(chapters))
	}

	sidebar := buildSidebar(results)

	totalChapters := 0
	for _, r := range results {
		totalChapters += len(r.chapters)
		page := renderPartPage(r, sidebar, theme)
		writeFile(filepath.Join(outDir, r.meta.Slug+".html"), page)
	}

	writeFile(filepath.Join(outDir, "index.html"), renderIndexPage(results, sidebar, theme))

	data, err := json.MarshalIndent(search, "", "  ")
	if err != nil {
		log.Fatalf("marshal search index: %v", err)
	}
	writeFile(filepath.Join(outDir, "search-index.json"), string(data))

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
		if len(excerpt) > 240 {
			excerpt = excerpt[:240] + "…"
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

func buildSidebar(results []partResult) string {
	var b strings.Builder
	b.WriteString(`<nav class="reader-sidebar" aria-label="Table of contents">`)
	for _, r := range results {
		fmt.Fprintf(&b, `<div class="toc-part" data-part="%s"><div class="toc-part-title">%s</div><ul>`, r.meta.Slug, html.EscapeString(r.meta.Title))
		for _, ch := range r.chapters {
			fmt.Fprintf(&b, `<li><a href="%s.html#%s">%s</a></li>`, r.meta.Slug, ch.SectionID, html.EscapeString(ch.Title))
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

func pageHead(title, canonicalSlug, theme string) string {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString(`<meta charset="utf-8">` + "\n")
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">` + "\n")
	fmt.Fprintf(&b, "<title>%s — Networking with Go, Made Easy</title>\n", html.EscapeString(title))
	b.WriteString(`<meta name="robots" content="index, follow">` + "\n")
	fmt.Fprintf(&b, `<link rel="canonical" href="https://sazardev.github.io/networking-with-go/read/%s.html">`+"\n", canonicalSlug)
	fmt.Fprintf(&b, `<link rel="icon" type="image/svg+xml" href="%s">`+"\n", favicon)
	b.WriteString("<style>\n" + theme + "\n</style>\n")
	b.WriteString(`<link rel="stylesheet" href="reader.css">` + "\n")
	b.WriteString("</head>\n")
	return b.String()
}

func pageHeader() string {
	return `<header class="reader-top">
  <a class="brand" href="../index.html">net<span>/</span>go<span>.</span>book</a>
  <div class="reader-top-actions">
    <a href="index.html">All Parts</a>
    <button id="search-open" type="button">Search <kbd>Ctrl K</kbd></button>
    <button id="theme-toggle" type="button" aria-pressed="false">Dark mode</button>
  </div>
</header>
`
}

func paletteMarkup() string {
	return `<div id="palette" class="palette" hidden>
  <div class="palette-box">
    <input id="palette-input" type="text" placeholder="Search the book..." autocomplete="off">
    <div id="palette-results"></div>
  </div>
</div>
<script src="reader.js"></script>
`
}

func renderPartPage(r partResult, sidebar, theme string) string {
	var b strings.Builder
	b.WriteString(pageHead(r.meta.Title, r.meta.Slug, theme))
	fmt.Fprintf(&b, "<body data-current-part=\"%s\">\n", r.meta.Slug)
	b.WriteString(pageHeader())
	b.WriteString(`<div class="reader-layout">` + "\n")
	b.WriteString(sidebar)
	b.WriteString(`<main class="reader-content">` + "\n")
	b.WriteString(r.bodyHTML)
	b.WriteString("\n</main>\n</div>\n")
	b.WriteString(paletteMarkup())
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

func renderIndexPage(results []partResult, sidebar, theme string) string {
	var b strings.Builder
	b.WriteString(pageHead("Read Online", "index", theme))
	b.WriteString("<body>\n")
	b.WriteString(pageHeader())
	b.WriteString(`<div class="reader-layout">` + "\n")
	b.WriteString(sidebar)
	b.WriteString(`<main class="reader-content">` + "\n")
	b.WriteString(`<h1>Read Networking with Go, Made Easy — online</h1>` + "\n")
	b.WriteString(`<p>All 143 chapters, six parts, right in the browser — same content as the PDF, same fonts, always in sync with the latest edit.</p>` + "\n")
	b.WriteString(`<div class="part-grid">` + "\n")
	for i, r := range results {
		firstChapter := ""
		if len(r.chapters) > 0 {
			firstChapter = r.chapters[0].SectionID
		}
		fmt.Fprintf(&b, `<a class="part-card" href="%s.html#%s"><span class="part-idx">%02d</span><h2>%s</h2><p>%s</p><span class="part-count">%d chapters</span></a>`+"\n",
			r.meta.Slug, firstChapter, i+1, html.EscapeString(r.meta.Title), html.EscapeString(r.meta.Desc), len(r.chapters))
	}
	b.WriteString("</div>\n</main>\n</div>\n")
	b.WriteString(paletteMarkup())
	b.WriteString("</body>\n</html>\n")
	return b.String()
}
