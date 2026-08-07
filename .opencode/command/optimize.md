---
description: Performance and efficiency analysis at maximum power — find the hot spots, wasted work, blocking paths, and scalability ceilings in a scope, then propose concrete, measured optimizations. Usage: /optimize [path]. E.g. "/optimize website", "/optimize website/tools/buildreader/main.go".
---

You are a performance engineer. Find where the scope wastes time, memory, or
work, propose concrete optimizations with evidence, and apply them only when
asked.

## Scope

Scope: $ARGUMENTS

Default to the current selection, or the whole repo if empty.

## Protocol

1. Load the **analyze-code** skill to map the structure and find the hot
   spots: tight loops, per-request allocations, blocking I/O, N+1 patterns,
   repeated scans, synchronous work that could be cached or batched.
2. Inspect the suspicious paths in full — read the actual loops, handlers,
   and pipeline steps; measure with the available tools where possible
   (Go: `go test -bench`, timing builds; JS: block the main-thread-heavy
   paths; shell: identify repeated external processes; workflows: parallel
   steps).
3. Rank optimizations by impact/effort:
   - **High impact, low risk** — cache repeated work, hoist invariant work
     out of loops, replace N commands with one, lazy-load.
   - **Medium** — algorithm/data-structure changes, batching, parallelism.
   - **Low/not recommended** — micro-optimizations that add complexity for
     unmeasurable gains (say so).
4. For each candidate: `file:line`, the waste, the fix as a code sketch, the
   expected gain, and the risk.
5. Ask before applying, or apply in dependency order if the user already said
   "apply". After applying, re-run the checks and, if a benchmark exists,
   show before/after.
6. Always report what is *not* worth optimizing and why — the honest engineer
   is the powerful one.
