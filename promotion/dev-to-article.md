---
title: "Rendering 25 Mermaid diagrams into a print PDF (and what broke)"
published: false
description: "The PDF tool printed Mermaid blocks as source text, blocked network and file:// fetches, and printed before JavaScript could run. Here is the build step that fixed it."
tags: go, mermaid, pdf, opensource
cover_image: https://sazardev.github.io/networking-with-go/assets/og-card.png
---

I maintain *[Networking with Go, Made Easy](https://sazardev.github.io/networking-with-go/)*, a free book that ships as a web reader, a PDF and an EPUB from the same Markdown sources. The web reader renders Mermaid diagrams in the browser. The PDF did not: it printed 25 blocks of `flowchart TB` and `sequenceDiagram` as literal source code.

This is the story of the build step that fixed it, including the parts that were not obvious.

## Three constraints, one conclusion

The PDF is produced by [go-pretty-pdf](https://github.com/sazardev/go-pretty-pdf), which composes the MDX into HTML and prints it with headless Chrome. Three details of that pipeline ruled out the obvious fixes:

1. **It has no Mermaid support.** A `mermaid` code fence becomes an ordinary code block.
2. **It blocks subresources by default.** Network *and* `file://` fetches are blocked, so an `<img src="diagram.svg">` inside the composed page never loads.
3. **It prints immediately after the load event.** There is no wait for JavaScript, so pulling in `mermaid.js` and letting the browser draw the diagram does not work either.

That leaves one shape that survives: the SVG must already be inline in the HTML before Chrome prints it.

## Pre-render, then materialize

The book keeps Mermaid blocks in the chapters. That is the source of truth, and it is what the web reader still uses.

A small Go tool, `tools/diagrams`, does two things:

1. It finds every `mermaid` block in `docs/`.
2. It writes a copy of the book to `build/docs/` with each block replaced by the matching SVG from `assets/diagrams/`.

The PDF build then runs against `build/docs` instead of `docs`:

```sh
go run tools/diagrams/main.go
pretty-pdf build --source build/docs --out networking-with-go.pdf
```

A `Makefile` wraps both commands as `make pdf` and `make epub`.

The asset lookup is content-addressed. The SVG file name is the SHA-256 of the diagram source, so editing a diagram without re-rendering it is a hard build error, not a stale picture:

```text
diagrams: 1 diagram(s) have no rendered SVG:
  d77887b698c1  docs/part1/03-types-of-networks.mdx
Run: go run tools/diagrams/main.go -render
```

Rendering itself uses the official `@mermaid-js/mermaid-cli`, pinned. It needs Node, so it is an authoring step only: CI never runs it, and the release pipeline stays Go-only.

## What broke

### 1. Monospace text and sequence diagrams

The book is set in a monospace face, so the first render profile forced the same stack inside Mermaid. Flowcharts were fine. Sequence diagrams were not: actor boxes and note rectangles came out narrower than their labels.

Mermaid measures those boxes with its own font metrics. Forcing a different family makes the measurement and the rendered text disagree, and the labels spill outside the shapes. The fix was to stop overriding the diagram font and let Mermaid measure what it renders. The figures ended up proportional while the prose stays monospace, a common distinction in book design.

### 2. HTML labels are not XHTML-friendly

Mermaid uses `<foreignObject>` with HTML inside for labels by default. Chrome prints that fine, but the EPUB is generated as XHTML, and foreign objects are a portability problem in e-readers. Setting `htmlLabels: false` (and the `classic` look) produces plain `<text>` and keeps the EPUB well-formed.

### 3. A page break ate the figure's top margin

The first PDF build had one diagram sitting flush against the top of the page, with the running header printed over its first row. Chrome drops a block's top margin when a page break lands on it. Margins were the wrong tool; padding survives fragmentation, so the figure wrapper uses padding for the top gap.

### 4. One diagram was 1144 px wide

A `flowchart LR` with four nodes laid out in a single row is 1144 px wide. Scaled to a 649 px text column, its 15 px labels print at about 6.5 pt. It is legible, but it is the one diagram I would redraw vertically if I were editing content rather than tooling. The pipeline does not try to fix layout decisions.

## The check that keeps it honest

A GitHub Actions job runs the materializer on every pull request touching `docs/` or `assets/diagrams/`. If someone edits a diagram and forgets to re-render it, the check fails with the hash it is missing. The release workflow runs the same step before building.

The result: zero raw Mermaid blocks in the PDF, 25 diagrams printed as vectors, and the EPUB still valid. The only raster image left in the PDF is the cover.

## What I would do differently

This logic belongs inside the PDF tool, not in a wrapper around it. A `render.mermaid` option that pre-renders with Chrome and inlines the SVG would remove the extra build step and the duplicated source tree. Until then, the wrapper is small, deterministic and easy to explain, which is the next best thing.

The book is free to read at [sazardev.github.io/networking-with-go](https://sazardev.github.io/networking-with-go/), and the pipeline lives in the [repository](https://github.com/sazardev/networking-with-go). If you are generating PDFs from Markdown with diagrams, the content-addressed asset trick is the part I would steal first.
