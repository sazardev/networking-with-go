# AGENTS.md

## Repo identity

This is an educational book (`.mdx` chapters under `docs/`) plus standalone Go exercises (`exercises/part2/`). There is no top-level `go.mod`, no test suite, and no CI. Do not treat it as a Go application.

## Commands

```sh
pretty-pdf check          # validate all docs/**/*.mdx format (run after any chapter edit)
pretty-pdf build          # render docs/ into networking-with-go.pdf
go vet ./exercises/part2/<dir>   # sanity-check an exercise (no CI enforces this)
```

Install the PDF tool once: `go install github.com/sazardev/go-pretty-pdf/cmd/pretty-pdf@latest`

## Running exercises

- **Stdlib-only exercises** (most): `go run ./exercises/part2/<exercise-dir>`
- **Third-party deps** (exercises 12 & 13, gorilla/websocket): `cd exercises/part2/<dir> && go run .` — these have their own `go.mod`

## Critical rules

- **Zero emoji** anywhere in docs or exercises — hard repo-wide rule.
- **Before editing any `docs/**/*.mdx`**: read `.claude/skills/mdx-pdf-format/SKILL.md` for required frontmatter (`id`/`title`), ID scheme, h3 depth cap, and allowed components.
- **After adding/removing/renaming/reordering a chapter in `docs/`**: update the matching link in `README.md` — they must stay in sync.
- **Part 1 chapters 3-13 are theory-only** — no Go code. Chapters 1-2 are the sole exception (one small snippet each at the end).
- **New exercises**: single `main.go` per directory, no `go.mod` unless a third-party dependency is needed, client/server pairs in separate directories.
