---
description: Full, maximum-power review of the book (or a scope) covering English grammar and style, chapter coherence and cross-references, book-wide flow and sequencing, mdx/PDF format compliance, and Go exercise correctness with a test-thinking pass. Use with an optional scope argument, e.g. "review-book docs/part2/05-tcp-in-depth-protocol-theory-and-go-implementation.mdx" or "review-book part3".
---

You are the ultimate editorial board for *Networking with Go — The Easy Way
Guide*. Run a deep, multi-dimensional review at maximum power. Do not skip
passes, do not soften findings, and do not stop at "it compiles."

## Scope

Scope is: $ARGUMENTS

If no scope is given, default to the **entire book** (all `docs/**/*.mdx`,
`README.md`, and `exercises/part2/`).

When a scope is given, prefix every command below with it where practical
(single file: run inside its folder; whole part: `rg ... docs/part3/`). The
whole-book forms below cover everything when scope is empty.

## Pass 1 — Inventory (run these first, with ripgrep)

Build the map mechanically before reviewing anything:

```sh
rg --files docs -g '*.mdx' | sort                                   # every chapter, grouped by folder
comm -23 <(rg -o 'docs/[a-z0-9-]+/[0-9]+-[^ )]+\.mdx' README.md | sort -u) <(rg --files docs -g '*.mdx' | sort -u)   # TOC links with no target file
comm -13 <(rg -o 'docs/[a-z0-9-]+/[0-9]+-[^ )]+\.mdx' README.md | sort -u) <(rg --files docs -g '*.mdx' | sort -u)   # files missing from README TOC
rg -o '\.\./\.\./exercises/part2/[A-Za-z0-9_-]+/main\.go' docs/part2 | sort -u   # every [Exercise: ...] link target
ls exercises/part2                                                              # actual dirs; cross-check against the links
```

- List every chapter grouped by part, in `id` order (derive order from the
  frontmatter `id`, not the filename).
- Cross-check against `README.md` (root) and `docs/part-apis/README.md`.
- State any gaps up front from the diff/comm output: missing TOC entries,
  exercise links without target directories, orphaned exercise dirs.

## Pass 2 — Format gate (non-negotiable, ripgrep first)

Run every one-liner below against the scope; ANY output is a finding, even if
it compiles. Do not eyeball these — the commands are faster and more complete.

```sh
rg --no-filename '^id:' docs | sort | uniq -d                                              # duplicate IDs (must be empty)
rg -n '^id:' docs | rg -v 'id: "\[\d+\.\d+\.\d+\]"'                                        # malformed id: not ^\[\d+\.\d+\.\d+\]$
rg -n '^#{4,} ' docs                                                                       # h3 cap violation
rg -U --pcre2 -n '<(DeepDive|Warning)\s+title="[^"]*">(?:(?!</\1>).)*$' docs                # DeepDive/Warning body not on ONE physical line
rg -U --pcre2 -n '<Axiom>(?:(?!</Axiom>).)*$' docs                                          # Axiom body not on ONE physical line
rg --files docs -g '*.mdx' | xargs awk '/^```/{f=!f} f && length($0)>85{print FILENAME":"FNR": "length($0)}'   # code-fence line over 85 chars
rg -n --pcre2 '[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{FE0F}\x{2190}-\x{21FF}\x{2B00}-\x{2BFF}\x{1F000}-\x{1F0FF}]' docs exercises README.md   # any emoji
```

If `rg --pcre2` is unavailable, fall back to the tag-balance count and inspect
the difference: `rg -o '<(DeepDive|Warning|Axiom)(\s|>)' docs | wc -l` versus
`rg -o '</(DeepDive|Warning|Axiom)>' docs | wc -l` (an imbalance means an
open tag is not closed on the same line or a prose `<Warning>` reference).

Cross-check each flagged tag against `.claude/skills/mdx-pdf-format/SKILL.md`,
then run `pretty-pdf check` from the repo root and report its output. Frontmatter
`id` scheme (X per top-level folder: 0 front-matter, 1 part1, 2 go-fundamentals,
3 part2, 4 advanced, 5 part3, 6 part-apis, 7 back-matter) and the one-physical-line
rule for `DeepDive`/`Warning`/`Axiom` are both enforced here.

## Pass 3 — English grammar and style

With the english-grammar-reviewer standard, proofread every chapter in scope:
grammar, punctuation, spelling, syntax, consistency, tone. Report
Critical/Grammar/Style/Nit with `file:line` and corrected text.

## Pass 4 — Coherence and flow

Start with the mechanical detectors, then apply the book-flow-analyst standard:

```sh
rg --no-filename -o '[A-Z][^.!?]*\.' docs | sort | uniq -d              # sentences repeated verbatim (heuristic: filter near-miss and FAQ refrains by hand)
rg -n '\]\([^)]+\)' docs | rg -v 'http|mailto|#|\{\{var'                # relative links; resolve each and confirm the target exists
rg -n '\]\(\.\./\.\./exercises/part2/' docs/part2 | rg -v 'main(_test)?\.go'   # exercise links not pointing at a main.go / main_test.go
```

Then, by hand:
- Verify ordering, forward/backward references, and "next chapter" hooks.
- Verify cross-references point at real files; verify concept dependencies
  (nothing used before it is explained).
- Check for duplicated explanations and contradictory facts across chapters
  (ports, API names, code snippets, terminology) — re-run the repeated-sentence
  detector after rewrites to confirm removals.
- Spot-check hardcoded port/address numbers with
  `rg -o ':\d{2,5}\b' docs/part2 docs/part-apis | sort -u` for accidental
  collisions with the standard ports the book teaches.

## Pass 5 — Go exercises and test-thinking

For every exercise in scope (or every exercise linked from in-scope
chapters):
- Verify with `GO111MODULE=off go vet . && GO111MODULE=off go build -o /tmp/out .`
  from inside the dir (or `go build ./...` for the gorilla dirs 12/13), plus
  `gofmt -l .`.
- Review correctness, error handling, resource cleanup, concurrency, idiom,
  security.
- Think like a test: list the 3-6 tests you would write for each exercise and
  what each would catch (edge inputs, races, timeouts, leaks, port conflicts,
  malformed input panics).
- Cross-check that the exercise file matches the code block embedded in its
  chapter (extract with `rg -n -A999 -m1 '^```go' docs/part2/NN-*.mdx` where
  useful) — the exercise should be the chapter's program, not a rewrite.

## Pass 6 — The report

Deliver one consolidated report:

1. **Summary** — an overall verdict with counts per category (list every
   non-empty ripgrep finding from Passes 1-2-4 first, verbatim).
2. **Critical issues** — anything that breaks the build, breaks a link, is
   factually wrong, leaks/races, or contradicts another chapter. These must be
   fixed.
3. **Important issues** — grammar errors, weak explanations, missing exercise
   pairing, stale TOC entries.
4. **Polish** — style nits, pacing, weak hooks.
5. **Test plans** per exercise reviewed.
6. **Verified commands run** and their results.

For each issue give `file:line`, the problem, and the fix. When the user asks
you to apply fixes, edit in place, re-run the affected verifications
(`pretty-pdf check`, vet/build/gofmt for any touched exercise, and the Pass 2
one-liners for any touched `.mdx`), and confirm the book still passes before
finishing.
