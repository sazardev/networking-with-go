# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

This is **not a Go application or library** — it's an educational book/guide ("Networking with Go - The Easy Way Guide") teaching network programming in Go, plus a matching set of standalone runnable code exercises. There is no top-level `go.mod` and no test suite. Work here is almost always either (a) writing/editing `.mdx` chapters in `docs/`, (b) writing/fixing small standalone Go programs in `exercises/`, or (c) building the book PDF.

The chapters are compiled into a single PDF with [go-pretty-pdf](https://github.com/sazardev/go-pretty-pdf) (`go install github.com/sazardev/go-pretty-pdf/cmd/pretty-pdf@latest`), configured via `go-pretty-pdf.yml` at the repo root (`source: docs`). **Read the `mdx-pdf-format` skill (`.claude/skills/mdx-pdf-format/SKILL.md`) before creating or editing any `docs/**/*.mdx` file** — it covers the required frontmatter (`id`/`title`), the per-folder ID scheme, the h3 heading-depth cap, and the `DeepDive`/`Warning`/`Axiom` components the linter and renderer enforce.

```sh
pretty-pdf check   # validate all docs/**/*.mdx against the format rules
pretty-pdf build   # render docs/ into the PDF configured in go-pretty-pdf.yml
```

## Repository structure

- `README.md` — the master table of contents for the whole book, linking every chapter in reading order. **If you add, remove, rename, or reorder a chapter file in `docs/`, update the corresponding entry/link in `README.md` too** — the two are expected to stay in sync.
- `docs/part1/` — pure networking theory (introduction, history, types of networks, topologies, OSI/TCP-IP, IP/subnetting, ports/sockets, TCP vs UDP, protocols, security fundamentals, firewalls/NAT/VPN, troubleshooting, performance), written in a warm, story-driven, passionate voice — real history, analogies, and fun facts, not dry reference prose. Chapters 3-13 are theory-only with no Go code. Chapters 1-2 are the sole exception (one small diagram snippet each at the end).
- `docs/go-fundamentals/` — a dedicated crash course in the Go language itself (installation through testing, plus a flagship goroutines/channels/concurrency chapter), sitting between Part 1 and Part 2. This is where hands-on Go code actually begins for the reader.
- `docs/part2/` — core Go networking topics (TCP/UDP, HTTP, WebSockets, DNS, concurrency, context, security, etc.); each chapter pairs theory with Go code samples inline. `exercises/part2/` contains the runnable counterparts to these chapters.
- `docs/advanced/` — specialized/advanced networking topics (gRPC, WebRTC, MQTT, SDN, NFV, etc.).
- `docs/part3/` — cybersecurity/offensive-defensive networking topics (scanning, sniffing, TLS/PKI, IDS/IPS, red/blue team, honeypots, etc.).
- `docs/part-apis/` — building modern APIs/backends in Go (REST, gRPC, Gin, Fiber, SQL/SQLite, Docker, Clean/Hexagonal Architecture, DDD, Google Cloud, security, etc.). `docs/part-apis/README.md` is a scoped index for this section only (stays plain `.md` — it's not a book chapter, so go-pretty-pdf ignores it).
- `exercises/part2/` — one directory per exercise, numbered to match its `docs/part2/` chapter (e.g. `08-goroutines-basic`, `10-http-server-routing`). Chapters with a client and server both get sibling directories (e.g. `06-udp-client` / `06-udp-server`, `13-chat-client-gorilla` / `13-chat-server-gorilla`).

## Running exercises

Each exercise directory is a single `package main` with a `main.go`. Most have **no `go.mod`** and depend only on the standard library — run them directly with:

```sh
go run ./exercises/part2/<exercise-dir>
```
This works only from inside an exercise directory with `GO111MODULE=off`, since Go detects the parent `.git` directory. The reliable command is:

```sh
GO111MODULE=off go run ./exercises/part2/<exercise-dir>
```

A few exercises pull in third-party packages (currently just `gorilla/websocket`, used by the WebSocket/chat exercises numbered 12 and 13) and therefore have their own `go.mod`/`go.sum` scoped to that single directory. For those, `cd` into the directory first:

```sh
cd exercises/part2/12-websocket-echo-server && go run .
```

Not all chapter 12 exercises use gorilla/websocket — the native WebSocket exercises (`websocket-hello-server`, `websocket-native-client`, `websocket-native-server`) are stdlib-only.

When adding a new exercise:
- Only add a `go.mod` if the exercise needs a dependency outside the standard library; otherwise leave it out, consistent with the rest of `exercises/part2/`.
- Keep each exercise self-contained in its own directory as a single `main.go` (matching existing exercises) rather than adding shared/library packages.
- Client/server pairs are split into separate directories, not combined into one program.

There is no linter or test suite configured for the exercises; `go vet`/`gofmt` are reasonable sanity checks but aren't wired into any CI.

## Writing/editing documentation chapters

- Chapters are numbered `.mdx` files; filenames encode both order and topic (`NN-topic-slug.mdx`). Preserve this numbering scheme when adding chapters, and don't renumber existing files without updating their `README.md` links.
- Every chapter needs YAML frontmatter with `id: "[X.Y.Z]"` and `title: "..."` — see the `mdx-pdf-format` skill for the ID scheme and the h3 heading-depth limit before writing content.
- Zero emoji anywhere in the book, in any chapter — a hard, repo-wide rule. Follow the existing tone instead: warm, story-driven, example-driven, with real history/analogies/fun-facts, never dry reference style.
- Part 2 onward mixes theory prose with embedded Go code blocks and diagrams; Part 1 in its entirety is theory-only prose with no code at all — see the structure notes above before adding code to the wrong section.
