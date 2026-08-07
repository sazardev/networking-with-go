---
description: Deep, multi-dimensional code review. Auto-detects the language from the path (Go / Shell / HTML-CSS-JS) and applies the matching review skill at maximum power. Usage: /review [path] [lang] — e.g. "/review website/tools/buildreader/main.go", "/review website", "/review . github". Leaving scope empty reviews the current selection or the whole repo.
---

You are a strict code reviewer at maximum power. Review the scope below with
the correct skill loaded, run the mechanical checks first, and do not soften
findings.

## Scope

Scope: $ARGUMENTS

Default to the current file selection, or the whole repo if nothing is given.

## Language detection

Detect the language from the scope and load the matching skill:

- `.go` files, `exercises/part2`, `website/tools/buildreader`, any Go package:
  use the **review-go** skill.
- `.sh` files and `.github/workflows/*.yml` steps, or any shell pipeline: use
  the **review-shell** skill.
- `website/*.html`, `*.css`, `*.js`, `sw.js`, `sw-precache.js`, PWA assets:
  use the **review-frontend** skill.
- Mixed or unknown: inventory the scope first, review each language with its
  skill, and merge the reports. If the user explicitly passed a language in
  the arguments, use it.

For Go code under `exercises/part2`, also apply the **go-exercise-reviewer**
skill (module-mode quirk and chapter pairing).

## Protocol

1. Read every file in scope in full before judging.
2. Run the mechanical checks for the language (vet/build/gofmt, shellcheck or
   `bash -n`, `node --check`, etc.) and report raw results.
3. Review against the skill's dimensions and checklists.
4. Run the skill's "what a test would catch" / failure-thinking pass.
5. Report with `file:line` evidence, prioritized findings, and concrete fixes.
6. If the user asks you to fix, apply edits and re-run the checks before
   reporting done.
