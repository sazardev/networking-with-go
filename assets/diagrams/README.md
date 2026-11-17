# Rendered diagrams

The book keeps its diagrams as <code>```mermaid</code> blocks inside
`docs/**/*.mdx`, which is what the web reader renders and what contributors
edit. The print pipeline cannot use those blocks directly: `go-pretty-pdf`
v0.8.0 prints a `mermaid` fence as literal text, blocks the network and
`file://` subresource loads, and prints immediately after page load, so neither
a CDN script nor a local `<img src="...">` can supply the drawing in time.

This directory holds the pre-rendered SVG for each diagram, named after the
SHA-256 of its source block. `go run tools/diagrams/main.go` inlines them into
`build/docs/` and refuses to build if a chapter's Mermaid source no longer
matches a checked-in file, so a stale image can never reach a release.

## Regenerating

After editing any diagram:

```sh
MERMAID_PUPPETEER_CONFIG=/path/to/puppeteer.json \
  go run tools/diagrams/main.go -render
```

`-render` needs Node.js and downloads the pinned
`@mermaid-js/mermaid-cli@12.2.0` through `npx`. `MERMAID_PUPPETEER_CONFIG` is
optional: point it at a JSON file with
`{"executablePath": "/usr/bin/google-chrome-stable", "args": ["--no-sandbox"]}`
when Mermaid cannot find a browser by itself (the same reason `pretty-pdf` has
`--chrome-path`). CI never runs this step; it only materializes the
already-committed SVGs, so releases stay Go-only and reproducible.

`mermaid.json` is the render profile: grayscale theme, classic look, and
text-based labels (`htmlLabels: false`, needed for XHTML-valid EPUB output).
It deliberately leaves the font family at Mermaid's default: sequence-diagram
actor and note boxes are measured with that default, so forcing the book's
monospace stack there makes labels wider than their boxes and clips them. The
figures therefore use a proportional face while the prose stays monospace.
