# CLAUDE.md

AGENTS.md at the repo root is the canonical source for project conventions,
commands, and rules. Refer to it for the full guide.

This file exists for Claude Code compatibility and echoes the most essential
content only. Anything stated here must also be present, and kept in sync, in
AGENTS.md. If they conflict, AGENTS.md wins.

## What this repo is

This is **not a Go application or library** -- it's an educational book
("Networking with Go - The Easy Way Guide") teaching network programming in
Go, plus a matching set of standalone runnable code exercises. There is no
top-level `go.mod`, no test suite, and no CI for exercises. CI exists for
website/PDF builds under `.github/workflows/`.

The chapters are compiled into a single PDF with
[go-pretty-pdf](https://github.com/sazardev/go-pretty-pdf). **Read the
`mdx-pdf-format` skill (`.claude/skills/mdx-pdf-format/SKILL.md`) before
creating or editing any `docs/**/*.mdx` file.**

```sh
pretty-pdf check   # validate all docs/**/*.mdx against the format rules
pretty-pdf build   # render docs/ into the PDF configured in go-pretty-pdf.yml
```

## Repository structure

- `README.md` -- the master table of contents for the whole book. If you add,
  remove, rename, or reorder a chapter file in `docs/`, update the
  corresponding entry/link in `README.md` too.
- `docs/part1/` -- pure networking theory (chapters 3-13 theory-only, no Go
  code; chapters 1-2 have one small snippet each at the end).
- `docs/go-fundamentals/` -- Go language crash course (installation through
  concurrency). This is where hands-on Go code begins for the reader.
- `docs/part2/` -- core Go networking topics; each chapter pairs theory with
  Go code. `exercises/part2/` contains runnable counterparts.
- `docs/advanced/` -- specialized/advanced networking topics (gRPC, WebRTC,
  MQTT, SDN, etc.).
- `docs/part3/` -- cybersecurity/offensive-defensive networking topics.
- `docs/part-apis/` -- building modern APIs/backends in Go.
- `exercises/part2/` -- one directory per exercise, numbered to match its
  `docs/part2/` chapter.

## Running exercises

Most have no `go.mod` and depend only on the standard library:

```sh
GO111MODULE=off go run ./exercises/part2/<exercise-dir>
```

Exercises with `gorilla/websocket` (numbered 12 and 13) have their own
`go.mod`:

```sh
cd exercises/part2/12-websocket-echo-server && go run .
```

When adding a new exercise:
- Only add a `go.mod` if the exercise needs a dependency outside the standard
  library.
- Keep each exercise self-contained as a single `main.go`.
- Split client/server pairs into separate directories.

## Writing/editing documentation chapters

- Every chapter needs YAML frontmatter with `id: "[X.Y.Z]"` and `title: "..."`.
- Zero emoji anywhere in docs or exercises -- hard repo-wide rule.
- Part 2 onward mixes theory prose with embedded Go code; Part 1 chapters 3-13
  are theory-only prose.
