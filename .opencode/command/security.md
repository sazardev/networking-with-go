---
description: Full security audit of code, config, workflows, and secrets exposure — attack surface map, mechanical scans, severity-ranked findings with fixes, and an explicit clean bill for what is safe. Usage: /security [path]. E.g. "/security", "/security website".
---

You are a security auditor at maximum power. Audit the scope below with the
**audit-security** skill — attack surface, secrets, injection, TLS, untrusted
input, CI/CD supply chain, and dependencies.

## Scope

Scope: $ARGUMENTS

Default to the whole repo if nothing is given.

## Protocol

1. Map the attack surface: code that accepts input, code that executes
   things, credentials storage, and anything deployed or published.
2. Run the mechanical scans and report raw hits:
   - Secrets: `rg -ni '(password|secret|token|api[_-]?key|BEGIN (RSA|OPENSSH|EC) PRIVATE KEY|authorization|bearer)' <scope>` — triage each hit.
   - Plain HTTP / weak TLS: `rg -n 'http://' <scope>` (excluding localhost/example).
   - Injection points: `rg -n 'exec\.Command|os/exec|eval\b|innerHTML|document\.write|sh -c' <scope>`.
   - Unbounded reads: `rg -n 'io\.ReadAll|ioutil\.ReadAll|make\(\[\]byte' <scope>`.
3. Audit each area with the relevant review-* lens (Go, shell, frontend),
   plus config/credentials, network, CI/CD (`permissions:` scoped, actions
   pinned, secrets minimal), and dependencies (scanner if available, or state
   the repo is stdlib-only).
4. Report severity-ranked findings (Critical/High/Medium/Low) each with
   `file:line`, the concrete path to harm, and the fix as a snippet.
5. Give an explicit clean bill for what you checked and found safe, and list
   the commands run.
6. If asked to fix, apply the edits and re-scan before reporting done.
