---
name: mdx-pdf-format
description: Use whenever creating, renaming, or editing any file under docs/ in this repo, or when asked to compile/build the book PDF. Enforces the exact .mdx frontmatter, ID scheme, and heading-depth rules required by the go-pretty-pdf tool (https://github.com/sazardev/go-pretty-pdf) so the whole book keeps compiling.
---

# MDX format for go-pretty-pdf

This repo's book (`docs/`) is built into a PDF with `go-pretty-pdf`
(`go install github.com/sazardev/go-pretty-pdf/cmd/pretty-pdf@latest`). It only
picks up `.mdx` files (recursive walk of the `source` directory in
`go-pretty-pdf.yml`, currently `docs`) — `.md` files (like the `README.md`
indexes) are ignored by the build and can stay plain Markdown.

Every chapter file that should appear in the PDF **must** be `.mdx` and follow
the rules below. Breaking any of these makes `pretty-pdf check`/`build` fail.

## Required frontmatter

```yaml
---
id: "[X.Y.Z]"
title: "Chapter Title"
---
```

- `id` is mandatory and must match the regex `^\[\d+\.\d+\.\d+\]$` exactly —
  brackets included, three dot-separated integers, quoted as a YAML string.
  Documents are ordered in the PDF **by this ID**, not by filename.
- `title` is mandatory, quoted, and is what renders in the TOC/header.
- Optional fields (`subtitle`, `tags`, `difficulty`, `status`,
  `completeness`, `depends_on`) may be added but are not required — this repo
  intentionally omits them to keep frontmatter minimal; don't add them unless
  the user asks.
- No two files anywhere under `docs/` may share the same `id`.

### This repo's ID scheme

`X` is fixed per top-level folder, `Y` is the file's existing two-digit
numeric prefix, `Z` is reserved (always `0` today):

| Folder | X |
|---|---|
| `docs/part1/` | 1 |
| `docs/part2/` | 2 |
| `docs/advanced/` | 3 |
| `docs/part3/` | 4 |
| `docs/part-apis/` | 5 |

Example: `docs/part2/05-tcp-in-depth-...mdx` → `id: "[2.5.0]"`.

When adding a **new** chapter, keep this scheme: pick the `X` for its folder,
and a `Y` that doesn't collide with an existing file in that folder (typically
the next free number, matching the file's numeric filename prefix).

## Heading depth

`go-pretty-pdf`'s linter caps heading depth at **h3** (`max_heading_depth: 3`
in `go-pretty-pdf.yml`). Never use `####`/`#####` in chapter content — if a
sub-point needs its own label, use a bold line (`**Request:**`) instead of a
new heading level, or fold it into the surrounding `###` section.

## Components available in MDX body

Three custom tags are transpiled to styled HTML — use them instead of plain
blockquotes when the intent matches:

- `<DeepDive title="...">...</DeepDive>` — blue info panel
- `<Warning title="...">...</Warning>` — orange warning panel
- `<Axiom>...</Axiom>` — green italic pull-quote

`title` is optional on `DeepDive`/`Warning`. Only simple inline Markdown
(code spans, `**bold**`) is rewritten inside these tags — don't nest complex
Markdown (tables, nested lists, headings) inside them.

## Variable substitution

`{{var}}` in any `.mdx` file is replaced from the `vars:` map in
`go-pretty-pdf.yml` before parsing. Only use this for values that are genuinely
config-driven (versions, product name); don't introduce new `{{...}}` vars
without adding them to `go-pretty-pdf.yml`.

## Keeping the index in sync

`README.md` (root) and `docs/part-apis/README.md` link to every chapter.
Whenever a chapter file is added, renamed, or removed, update the matching
link in whichever of those two index files references it — including the
`.mdx` extension in the link path.

## Build/validate commands

```sh
pretty-pdf check        # validate frontmatter, IDs, heading depth
pretty-pdf build         # render docs/ -> the PDF configured in go-pretty-pdf.yml
```

Run `pretty-pdf check` after adding or editing any `.mdx` file to catch
format violations before they reach the build.
