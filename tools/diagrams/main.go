// Command diagrams turns the ```mermaid fenced blocks in docs/ into inline
// SVG for the print pipeline.
//
// go-pretty-pdf v0.8.0 has no Mermaid support: it printed a ```mermaid block as
// literal text, which is how all 25 diagrams in the book used to reach the PDF
// as raw "flowchart TB" source. It also blocks network and file:// subresource
// fetches, so an <img src="..."> pointing at an SVG cannot work, and its print
// step runs the instant the page loads, so Mermaid cannot be rendered live in
// the browser either. Pre-rendered inline SVG is the only shape that survives
// the pipeline, and it prints as real vectors in both the PDF and the EPUB.
//
// The tool copies docs/ into build/docs/, replacing each ```mermaid fence with
// the matching SVG from assets/diagrams/. SVGs are content-addressed: the file
// name is the SHA-256 of the diagram source, so editing a diagram without
// re-rendering it is a hard error instead of a stale picture.
//
//	go run tools/diagrams/main.go          # materialize build/docs
//	go run tools/diagrams/main.go -render  # render missing diagrams (needs Node)
//
// The -render step shells out to Mermaid's official CLI, pinned below. It is
// an authoring convenience only: CI never runs it, so the release build needs
// nothing but Go and pretty-pdf.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// mermaidCLI is the exact @mermaid-js/mermaid-cli version used to produce the
// checked-in SVGs. Bump it, run -render, and commit the new assets.
const mermaidCLI = "@mermaid-js/mermaid-cli@12.2.0"

var (
	fenceRe = regexp.MustCompile("(?ms)^```mermaid[ \t]*\r?\n(.*?)^```[ \t]*$")
	xmlDecl = regexp.MustCompile(`(?s)^<\?xml.*?\?>\s*`)
)

type diagram struct {
	hash   string
	source string
	file   string
	line   int
}

func main() {
	docsDir := flag.String("docs", "docs", "source MDX directory")
	assetsDir := flag.String("assets", "assets/diagrams", "directory holding the rendered SVGs")
	outDir := flag.String("out", "build/docs", "directory to write the materialized book into")
	render := flag.Bool("render", false, "render missing diagrams with @mermaid-js/mermaid-cli (requires Node)")
	flag.Parse()

	diagrams, err := scan(*docsDir)
	if err != nil {
		fatal(err)
	}
	if len(diagrams) == 0 {
		fatal(fmt.Errorf("no ```mermaid blocks found under %s", *docsDir))
	}

	if *render {
		if err := renderMissing(diagrams, *assetsDir); err != nil {
			fatal(err)
		}
	}

	if err := materialize(*docsDir, *outDir, *assetsDir, diagrams); err != nil {
		fatal(err)
	}

	fmt.Printf("Materialized %d chapters into %s with %d rendered diagram(s).\n",
		countChapters(*docsDir), *outDir, len(diagrams))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "diagrams:", err)
	os.Exit(1)
}

// scan walks every .mdx file under root and collects the Mermaid blocks it
// contains, in file order.
func scan(root string) ([]diagram, error) {
	var diagrams []diagram
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".mdx") {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, loc := range fenceRe.FindAllStringSubmatchIndex(string(body), -1) {
			source := normalize(string(body[loc[2]:loc[3]]))
			diagrams = append(diagrams, diagram{
				hash:   hashOf(source),
				source: source,
				file:   path,
				line:   strings.Count(string(body[:loc[0]]), "\n") + 1,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning %s: %w", root, err)
	}
	return diagrams, nil
}

// normalize makes the hash insensitive to line endings and trailing
// whitespace, so re-saving a chapter does not invalidate every asset.
func normalize(source string) string {
	source = strings.ReplaceAll(source, "\r\n", "\n")
	return strings.TrimRight(source, " \t\n")
}

func hashOf(source string) string {
	sum := sha256.Sum256([]byte(source))
	return hex.EncodeToString(sum[:])[:12]
}

// renderMissing renders every diagram whose asset does not exist yet. Identical
// diagrams share a hash, so each source is rendered exactly once.
func renderMissing(diagrams []diagram, assetsDir string) error {
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		return err
	}

	done := map[string]bool{}
	for _, d := range diagrams {
		if done[d.hash] {
			continue
		}
		done[d.hash] = true
		out := filepath.Join(assetsDir, d.hash+".svg")
		if _, err := os.Stat(out); err == nil {
			continue
		}
		// Only demand Node when there is actually something to render, so a
		// no-op `make diagrams` works on a machine without it.
		if _, err := exec.LookPath("npx"); err != nil {
			return errors.New("-render needs npx (Node.js) on PATH; CI never needs this step")
		}

		tmp, err := os.CreateTemp("", "diagram-*.mmd")
		if err != nil {
			return err
		}
		name := tmp.Name()
		if _, err := tmp.WriteString(d.source + "\n"); err != nil {
			tmp.Close()
			os.Remove(name)
			return err
		}
		tmp.Close()

		fmt.Printf("Rendering %s (%s:%d)\n", d.hash, d.file, d.line)
		err = runMermaid(name, out, assetsDir, d.hash)
		os.Remove(name)
		if err != nil {
			return fmt.Errorf("rendering %s:%d: %w", d.file, d.line, err)
		}
	}
	return nil
}

func runMermaid(input, output, assetsDir, hash string) error {
	args := []string{
		"-y", mermaidCLI,
		"-q",
		"-i", input,
		"-o", output,
		"-c", filepath.Join(assetsDir, "mermaid.json"),
		"-I", "dg-" + hash,
		"--no-font-embed",
	}
	// A local Chrome avoids Puppeteer downloading its own; only needed when the
	// CLI cannot find a browser on its own. Same reason pretty-pdf has
	// --chrome-path.
	if cfg := os.Getenv("MERMAID_PUPPETEER_CONFIG"); cfg != "" {
		args = append(args, "-p", cfg)
	}
	cmd := exec.Command("npx", args...)
	cmd.Stdout = io.Discard
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if _, err := os.Stat(output); err != nil {
		return fmt.Errorf("mermaid-cli produced no output: %w", err)
	}
	return nil
}

// materialize mirrors docs/ into out/, rewriting .mdx files in place of the
// fenced blocks and copying everything else untouched.
func materialize(docsDir, outDir, assetsDir string, diagrams []diagram) error {
	sources, err := filepath.Abs(docsDir)
	if err != nil {
		return err
	}
	target, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}
	// Never let a typo'd -out delete the book itself: refuse the current
	// directory, the docs directory, or any ancestor of either.
	sep := string(os.PathSeparator)
	if cwd, err := os.Getwd(); err == nil {
		if target == cwd || sources == target || strings.HasPrefix(sources+sep, target+sep) {
			return fmt.Errorf("refusing to write into %s", target)
		}
	}
	if err := os.RemoveAll(target); err != nil {
		return err
	}

	var missing []diagram
	err = filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(docsDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(outDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(strings.ToLower(path), ".mdx") {
			return os.WriteFile(target, body, 0o644)
		}

		svgFor := func(source string) string {
			hash := hashOf(normalize(source))
			svg, err := os.ReadFile(filepath.Join(assetsDir, hash+".svg"))
			if err != nil {
				missing = append(missing, diagram{hash: hash, file: path})
				return ""
			}
			return figure(string(svg))
		}

		body = fenceRe.ReplaceAllFunc(body, func(match []byte) []byte {
			loc := fenceRe.FindSubmatchIndex(match)
			return []byte(svgFor(string(match[loc[2]:loc[3]])))
		})
		return os.WriteFile(target, body, 0o644)
	})
	if err != nil {
		return fmt.Errorf("materializing %s: %w", outDir, err)
	}

	if len(missing) > 0 {
		var b strings.Builder
		fmt.Fprintf(&b, "%d diagram(s) have no rendered SVG:\n", len(missing))
		for _, d := range missing {
			fmt.Fprintf(&b, "  %s  %s\n", d.hash, d.file)
		}
		b.WriteString("Run: go run tools/diagrams/main.go -render")
		return errors.New(b.String())
	}
	return nil
}

// figure wraps a rendered SVG in a block that stays on one page and never
// overflows the text column. The style has to live on the element itself:
// pretty-pdf's --css replaces the whole theme stylesheet, and its themes have
// no rule for figures.
//
// The top spacing is padding, not margin, on purpose: Chrome drops a block's
// top margin when a page break lands on it, which leaves the figure flush with
// the page edge and the running header printed over its first row. Padding
// survives fragmentation, so it keeps the same ~19pt clearance that normal
// body text gets from its own margins.
func figure(svg string) string {
	svg = xmlDecl.ReplaceAllString(svg, "")
	svg = strings.TrimSpace(svg)
	svg = responsive(svg)
	return "\n<figure style=\"margin: 0 0 1.3em; padding: 1.8em 0 0; text-align: center; break-inside: avoid;\">\n" +
		svg + "\n</figure>\n"
}

func responsive(svg string) string {
	const style = `style="max-width: 100%; height: auto;"`
	if strings.Contains(svg, `style="background-color: white;"`) {
		return strings.Replace(svg, `style="background-color: white;"`, style, 1)
	}
	if end := strings.Index(svg, ">"); end > 0 && !strings.Contains(svg[:end], "max-width") {
		return svg[:end] + " " + style + svg[end:]
	}
	return svg
}

func countChapters(root string) int {
	n := 0
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".mdx") {
			n++
		}
		return nil
	})
	return n
}
