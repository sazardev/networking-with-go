# AGENTS.md

## Repo identity

This is an educational book (`.mdx` chapters under `docs/`) plus standalone Go exercises (`exercises/part2/`). There is no top-level `go.mod` and no unified test suite. Do not treat it as a Go application. CI exists under `.github/workflows/`: it vets all exercises, builds the web reader, and builds/releases the PDF+EPUB.

## Commands

```sh
pretty-pdf check                              # validate all docs/**/*.mdx format (run after any chapter edit)
make pdf                                      # inline diagrams + render docs/ into networking-with-go.pdf
make epub                                     # same content as an EPUB 3 file
make diagrams                                 # re-render changed Mermaid blocks (needs Node.js)
GO111MODULE=off go vet ./exercises/part2/<dir>  # sanity-check one stdlib exercise
```

- `make pdf`/`make epub` run `tools/diagrams/main.go` first, which copies `docs/` to the generated `build/docs/` with each ```` ```mermaid ```` block replaced by the matching SVG in `assets/diagrams/`. Plain `pretty-pdf build` still works but prints the Mermaid source as text -- use the make targets for anything that ships.
- `make diagrams` renders missing/changed diagrams with the pinned `@mermaid-js/mermaid-cli`; CI never runs it, so a diagram edit must commit its new `assets/diagrams/<hash>.svg` or the release build fails on purpose.

Install the PDF tool once: `go install github.com/sazardev/go-pretty-pdf/cmd/pretty-pdf@v0.8.0`

- The `GO111MODULE=off` flag is required because Go detects the parent `.git` directory with no `go.mod`. The plain form fails with "cannot find main module".
- The same vet is enforced by CI on every push/PR touching `exercises/` (`.github/workflows/vet-exercises.yml`); it skips dirs with a `go.mod` (vets those in normal module mode) and dirs with no `.go` files at root.

## Running exercises

- **Stdlib-only exercises** (most): `GO111MODULE=off go run ./exercises/part2/<exercise-dir>`
- **gorilla/websocket exercises** (7 dirs: 12-websocket-client, 12-websocket-echo-server, 12-websocket-hello-server-gorilla, and the four `13-chat-*` gorilla dirs): `cd exercises/part2/<dir> && go run .` — these have their own `go.mod`. Do not add `GO111MODULE=off` to these. The native WebSocket exercises (12-websocket-hello-server, 12-websocket-native-client, 12-websocket-native-server) are stdlib-only.

## Critical rules

- **Zero emoji** anywhere in docs or exercises -- hard repo-wide rule.
- **Before editing any `docs/**/*.mdx`**: read `.claude/skills/mdx-pdf-format/SKILL.md` for required frontmatter (`id`/`title`), ID scheme (X per top-level docs folder), h3 depth cap, the one-physical-line rule for `DeepDive`/`Warning`/`Axiom` tags, and code-fence line-length limit.
- **After adding/removing/renaming/reordering a chapter in `docs/`**: update the matching link in `README.md` (root) and/or `docs/part-apis/README.md` -- the book is ordered by frontmatter `id`, and links carry the `.mdx` extension.
- **Part 1 chapters 1-13 are theory-only** -- no Go code. Chapters 1-2 contain only ASCII diagrams, and chapters 3-13 contain no code fences at all.
- **New exercises**: single `main.go` per directory, no `go.mod` unless a third-party dependency is needed, client/server pairs in separate directories. Verify with `GO111MODULE=off go vet . && GO111MODULE=off go build -o /tmp/out .` from inside the dir, plus `gofmt -l`.

## Website / web reader

The reader under `website/read/` is generated from `docs/**/*.mdx` by `website/tools/buildreader/main.go`; CI rebuilds it on every deploy (`.github/workflows/deploy-pages.yml`) and generated files are gitignored -- never hand-edit them.

## Tool versions

- `pretty-pdf` is pinned to `@v0.8.0` in both `deploy-pages.yml` and `release-pdf.yml`. When upgrading, update both workflow files and this line.
- `@mermaid-js/mermaid-cli` is pinned to `@12.2.0` in `tools/diagrams/main.go` (the `-render` path). Bumping it means re-rendering every diagram with `make diagrams` after deleting `assets/diagrams/*.svg`, then committing the new files.
