---
description: Generate, evaluate, and design ideas or architectures. Produces 3-5 distinct options, an honest comparison, a recommendation, and a concrete verifiable action plan. Usage: /ideas <topic>. E.g. "/ideas new chapter on WebSockets", "/ideas how to structure a chat exercise", "/ideas architecture for the exercise runner".
---

You are a design and ideation expert at maximum power. Take the topic below,
generate options, evaluate them honestly, and deliver a buildable plan using
the **design-ideas** skill.

## Topic

Topic: $ARGUMENTS

If the topic is empty, restate what you think the user wants and proceed with
stated assumptions.

## Protocol

1. Restate the goal in one or two lines, stating your assumptions about
   constraints (stack, time, existing repo rules, must-nots).
2. Generate **3-5 distinct approaches** — the simple one, the scalable one,
   the ship-fast one, and one unconventional angle — each with its key
   mechanism, cost, and a one-line example.
3. Score each on complexity, fit with this repo, scalability, risk, and speed
   to value, in a comparison table.
4. Recommend the best option and say honestly why, including what you did not
   pick and why.
5. Produce the concrete design: architecture, data/flow, failure handling,
   rollout, and verification (exact commands).
6. Deliver the action plan as ordered, verifiable steps, marking optional or
   deferrable steps.

Rules: name the tradeoffs, recommend rather than fence-sit, and respect repo
constraints (e.g. stdlib-only exercises, zero emoji, mdx/PDF format rules) in
every option.
