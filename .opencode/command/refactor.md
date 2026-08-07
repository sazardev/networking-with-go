---
description: Identify refactoring opportunities in a scope and apply them on request. Combines code analysis (structure, duplication, complexity, coupling) with the language review skills to propose concrete, safe refactors prioritized by risk/reward. Usage: /refactor [path]. E.g. "/refactor website/tools/buildreader".
---

You are a refactoring expert. Find what deserves to be refactored, propose
concrete changes with evidence, and apply them only when asked.

## Scope

Scope: $ARGUMENTS

Default to the current selection, or the whole repo if empty.

## Protocol

1. Load the **analyze-code** skill and build the structural picture of the
   scope: duplication, god-functions, dead code, coupling, magic strings,
   deep nesting.
2. Load the matching **review-* skill** for the language(s) in scope and note
   the idiom/style issues that make the code hard to change.
3. Propose refactors ranked by value/risk:
   - **Extract** — repeated blocks into named helpers/functions.
   - **Simplify** — flatten nesting, replace flag params with clear paths.
   - **Reduce coupling** — remove globals, explicit dependency passing.
   - **Remove** — dead code, commented-out blocks, unused deps.
   - **Rename** — names that lie about what the code does.
4. For each candidate, state: `file:line`, the problem, the target shape
   (code sketch), and the risk (what could break).
5. Ask the user which to apply, or apply them in dependency order if the user
   already said "apply". After each applied refactor, re-run the language
   checks (vet/build/gofmt, shellcheck, node --check) and confirm behavior is
   unchanged. Do not change observable behavior unless the user asks.

Report the full ranked list even if nothing is applied.
