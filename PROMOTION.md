# Promotion playbook

The book is free, so distribution is the only lever: get it in front of people
who already search for Go networking material, and give them a reason to share
it. This file tracks where it is listed and keeps copy ready to paste. Keep
claims checkable against the repository.

## Where it is listed

| Channel | Status | Link |
| --- | --- | --- |
| Awesome Go Books (`dariubs/GoBooks`) | PR open | https://github.com/dariubs/GoBooks/pull/181 |
| Free Programming Books (`EbookFoundation`) | PR open | https://github.com/EbookFoundation/free-programming-books/pull/13642 |
| Awesome Networking (`facyber`) | PR open | https://github.com/facyber/awesome-networking/pull/75 |
| GitHub topics | Done | `book`, `go`, `golang`, `networking`, `network-programming`, `tcp-ip`, `websockets`, `grpc`, `tls`, `http`, `cybersecurity`, `hacking`, `education`, `free` |
| Marketing kit | Live | https://sazardev.github.io/networking-with-go/marketing.html |
| dev.to article | Live | https://dev.to/sazardev/rendering-25-mermaid-diagrams-into-a-print-pdf-and-what-broke-55a4 |
| daily.dev post | Live | https://daily.dev/posts/free-open-source-book-networking-and-network-programming-in-go-e7phwhzww |
| Homepage download | Live | Resolves the latest release asset through the GitHub API |

## Channels that need the author's account

Automation is not an option here: Hacker News, Reddit, X, LinkedIn and Mastodon
all require a logged-in account, and automated cross-posting reads as spam.
Post from the account, one channel per day.

- **Hacker News** — Show HN, Tuesday to Thursday, 08:00-10:00 PT. Reply to every
  comment for the first four hours.
- **Reddit r/golang** — read the sidebar rules first; a text post that happens
  to link the book does better than a link post. r/networking is a separate,
  later post framed around the networking content, not Go.
- **dev.to** — technical article with `canonical_url` pointing at the site.
- **LinkedIn** — short, first-person, one question at the end.
- **X / Mastodon** — build notes and diagrams, not a sales pitch.

## Channels that do not fit

- **Go wiki Books** (`go.dev/wiki/Books`) — changes go through Gerrit
  (`go-review.googlesource.com/c/wiki`) and require the Go CLA; a GitHub PR is
  not the path.
- **awesome-go** — it has a Free e-books section, but its PR template and CI
  expect a package repository with `go.mod`, pkg.go.dev and a coverage report.
  A book cannot satisfy that checklist.

## Sequence

1. Day 0: Show HN. Answer questions; do not ask for votes.
2. Day 1: r/golang text post.
3. Day 2: dev.to article, then LinkedIn.
4. Day 3: X thread and Mastodon.
5. Day 5: r/networking.
6. Day 7: Gophers Slack `#showandtell`, Go Forum, newsletter blurbs.
7. Day 10: follow up on directory PRs.

## Copy

### Hacker News (Show HN)

Title: `Show HN: A free book on networking and network programming in Go`

First comment:

> I wrote this over the past year. It starts with how networks actually work
> (OSI/TCP-IP, addressing, ports, NAT, VPNs), then a Go crash course, then
> networking code: TCP, UDP, HTTP, WebSockets, DNS, TLS, proxies. Later parts
> cover gRPC, WebRTC, MQTT, QUIC, security tooling, and production APIs with
> Gin, Fiber, PostgreSQL and Docker.
>
> Free to read as web, PDF and EPUB. Text is CC BY-NC 4.0, code is MIT. No
> signup. The print PDF renders the 25 Mermaid diagrams as vector SVG, which
> took a small build step because the PDF tool prints Mermaid blocks as text.
>
> Happy to answer questions about the content or the pipeline.

### Reddit r/golang

Title: `I wrote a free, open-source book on networking and network programming in Go (156 chapters, TCP/IP to production APIs)`

Body:

> It is theory first, then complete runnable programs: a Go crash course, TCP,
> UDP, HTTP, JSON, WebSockets, DNS, proxies, TLS, then gRPC, WebRTC, MQTT,
> QUIC, security tooling, and production APIs.
>
> Free online, PDF and EPUB, no signup. Text CC BY-NC 4.0, code MIT:
> https://sazardev.github.io/networking-with-go/
>
> Feedback on pacing and gaps is welcome; the chapters and exercises are in the
> repository.

### Reddit r/networking

Title: `Free networking book with hands-on labs (the theory parts are language-agnostic)`

Body: lead with the networking scope (OSI/TCP-IP, subnetting, NAT, VPNs, TLS,
DNS, QUIC), then mention that the labs happen to be in Go because that is the
language the book teaches alongside.

### dev.to

Title: `Rendering 25 Mermaid diagrams into a print PDF (and what broke)`

Canonical URL: the article page on the site or the repository, not dev.to.
Outline: why the book is theory-first; the MDX to PDF pipeline; why the PDF tool
cannot render Mermaid at build time; content-addressed SVGs; the sequence
diagram font measurement bug; what the web reader does instead. End with the
free book link.

### LinkedIn

> I published a free, open-source book on networking and network programming in
> Go: 156 chapters, from TCP/IP fundamentals to production APIs and security
> tooling. Theory first, then complete runnable programs. No paywall, no
> signup, text CC BY-NC 4.0 and code MIT.
>
> https://sazardev.github.io/networking-with-go/
>
> If you mentor junior backend engineers: which networking topic do they
> struggle with most?

### X thread

1. Most Go networking bugs are not Go bugs. They are TCP, DNS, TLS or HTTP
   assumptions. I wrote a free book about those assumptions, with runnable code.
2. Theory first: OSI/TCP-IP, addressing, ports, NAT, VPNs, then a Go crash
   course, then TCP, UDP, HTTP, WebSockets, DNS, TLS.
3. Later parts: gRPC, WebRTC, MQTT, QUIC, security tooling, production APIs.
   156 chapters, every example a complete program.
4. Free online, PDF and EPUB. CC BY-NC 4.0 text, MIT code.
   https://sazardev.github.io/networking-with-go/

### Mastodon

> Free, open-source book: networking and network programming in Go, from
> TCP/IP to production APIs. Theory first, runnable examples, no signup.
> https://sazardev.github.io/networking-with-go/

### Newsletter blurb

> This week's resource: Networking with Go, Made Easy, a free book that goes
> from networking fundamentals to production APIs in Go, with a full Go
> introduction and runnable code in every chapter. Text CC BY-NC 4.0, code MIT.

## Measuring

```sh
# Downloads for the current release (PDF and EPUB)
gh api repos/sazardev/networking-with-go/releases/latest \
  --jq '.assets[] | "\(.name) \(.download_count)"'

# Repository traffic, last 14 days
gh api repos/sazardev/networking-with-go/traffic/views
```

Watch downloads and referrers, not stars. A star costs one click; a download
means someone opened the book.
