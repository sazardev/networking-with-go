---
description: Complete structural and architectural analysis of files, directories, or the whole codebase — what exists, how it is organized, dependencies, complexity, duplication, dead code, and technical debt, all grounded in verified file:line evidence. Usage: /analyze [path]. E.g. "/analyze website/tools/buildreader", "/analyze exercises/part2".
---

You are a codebase analyst at maximum power. Build a complete, verified
picture of the scope using the **analyze-code** skill.

## Scope

Scope: $ARGUMENTS

Default to the whole repo if nothing is given.

## Protocol

1. Inventory mechanically first: list every file in scope with `rg --files`,
   group by type and directory, and size the big files with `wc -l`.
2. Read in dependency order: entrypoints first, then what they call. Read the
   files that matter fully; sample the rest.
3. Map the architecture: layers/components, responsibilities, dependency
   graph, data flow, shared state, cross-cutting concerns.
4. Measure quality signals mechanically where possible: duplication, dead
   code, god-functions, coupling, magic strings, TODOs/FIXMEs, stale artifacts.
5. Trace one real end-to-end path through the code to verify the map.
6. Report in the analyze-code output format: overview, architecture map,
   verified facts with `file:line`, quality findings with severity, direct
   answers to the original question, commands run, and open questions.

Rules: no hypotheticals — every structural claim cites `file:line`; mark
anything you could not verify as unverified rather than guessing.
