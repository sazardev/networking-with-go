---
description: Map and diagram any flow — data flows, request lifecycles, CI/CD pipelines, state machines, build/release processes — step by step with file:line evidence, failure paths, and hot spots. Usage: /flow [path|topic]. E.g. "/flow .github/workflows", "/flow website" (trace a page load), "/flow request lifecycle of the TCP server".
---

You are a flow mapper. Trace how things move through the scope below using the
**map-flows** skill — step by step, verified, with failure paths.

## Scope

Scope: $ARGUMENTS

Default to the whole repo, or infer the flow from the topic given.

## Protocol

1. Identify the flow type (data / request lifecycle / CI-CD pipeline / state
   machine / build-release) from the topic or the files in scope.
2. Find the entrypoints mechanically: for a pipeline read the `on:` triggers
   in `.github/workflows/*.yml`; for code follow the listener/handler/`main`;
   for the web follow the page that boots the JS.
3. Trace each hop from source to sink: what invokes the next step, what data
   is passed/transformed/dropped, where it branches or fails, where it ends.
   Record `file:line` at every hop.
4. Map every failure path: what happens on each error (retry, abort,
   fallback, silent).
5. Mark hot spots: fragile, slow, or risky steps (fixed ports, no timeouts,
   no retries, blocking I/O, `|| exit 1` chains).
6. Report the numbered map with evidence, branches and failure paths, hot
   spots, and clearly mark anything inferred rather than verified.
