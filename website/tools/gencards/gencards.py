#!/usr/bin/env python3
"""
Generates the full set of social (OG) cards for the book's website.

It reads the generated web reader under website/read/ (built by
website/tools/buildreader) and produces one branded 1200x630 card per
chapter plus cards for each part, the landing page, the read index and the
topics page, all featuring the real book cover. It also writes small JPEG
thumbnails, an assets-manifest.json for the marketing page, and a ZIP of
every full-size asset.

Run from the repository root (after `go run website/tools/buildreader`):

    python3 website/tools/gencards/gencards.py

Requires: rsvg-convert, ImageMagick (magick), python3.
"""
import base64
import html
import json
import os
import re
import subprocess
import zipfile

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", ".."))
READ = os.path.join(ROOT, "website", "read")
OUT = os.path.join(ROOT, "website", "assets")
COVER = os.path.join(ROOT, "assets", "cover.jpg")
TMP = "/tmp/opencode/ogcards3"

W, H = 1200, 630
TEXT_X, TEXT_W = 516, 608

BG = "#f7f3ee"
INK = "#2e2620"        # primary text
BODY = "#5c4e3e"       # hook / tagline text
MUT = "#8c7563"        # eyebrow / footer
FAINT = "#a99a87"      # secondary label
ACCENT = "#a88d75"     # accent bar / motif
FRAME = "#c8bba8"      # cover frame border

FONT = "Liberation Sans"
EM = 0.52              # rough average glyph width relative to font size
BODY_BASE = {1: 372, 2: 392, 3: 428}  # first body line baseline per title-line count

H1_RE = re.compile(r"<h1>(.*?)</h1>", re.S)
DESC_RE = re.compile(r'<meta name="description" content="([^"]*)"')


def wrap_words(text, chars_per_line):
    words = text.split()
    lines, cur = [], ""
    for w in words:
        if not cur:
            cur = w
        elif len(cur) + 1 + len(w) <= chars_per_line:
            cur += " " + w
        else:
            lines.append(cur)
            cur = w
    if cur:
        lines.append(cur)
    return lines


def title_block(title):
    """Return (svg, number_of_lines) for the big title, auto-sized and wrapped."""
    n = len(title)
    if n <= 22:
        size = 56
        lines = wrap_words(title, int(TEXT_W / (EM * size)))
        if len(lines) == 1:
            ys = [300]
        else:
            size, ys = 48, [272, 330]
    elif n <= 52:
        size, ys = 50, [272, 336]
        lines = wrap_words(title, int(TEXT_W / (EM * size)))
    else:
        size, ys = 42, [258, 316, 374]
        lines = wrap_words(title, int(TEXT_W / (EM * size)))
    out = []
    for i, ln in enumerate(lines[: len(ys)]):
        esc = html.escape(ln, quote=False)
        out.append(f'<text x="{TEXT_X}" y="{ys[i]}" font-size="{size}" '
                   f'font-family="{FONT}" font-weight="bold" fill="{INK}">{esc}</text>')
    return "".join(out), len(ys)


def body_block(text, size, title_lines):
    text = text.strip()
    if not text:
        return ""
    cpl = int(TEXT_W / (EM * size))
    lines = wrap_words(text, cpl)[:2]
    base = BODY_BASE.get(title_lines, 372)
    gap = int(size * 1.35)
    out = []
    for i, ln in enumerate(lines):
        esc = html.escape(ln, quote=False)
        out.append(f'<text x="{TEXT_X}" y="{base + i * gap}" font-size="{size}" '
                   f'font-family="{FONT}" fill="{BODY}">{esc}</text>')
    return "".join(out)


def banner(*parts, y=96, size=20, ls=6, fill=MUT):
    line = html.escape("  ".join(str(p) for p in parts), quote=False)
    return (f'<text x="{TEXT_X}" y="{y}" font-size="{size}" font-family="{FONT}" '
            f'letter-spacing="{ls}" fill="{fill}">{line}</text>')


def motif():
    return (
        '<g stroke="#a88d75" stroke-width="2" fill="none" opacity="0.22">'
        '<line x1="1010" y1="70" x2="1110" y2="130"/>'
        '<line x1="1110" y1="130" x2="1050" y2="220"/>'
        '<line x1="1050" y1="220" x2="1140" y2="270"/>'
        '<line x1="1010" y1="70" x2="1050" y2="220"/>'
        '<line x1="1060" y1="560" x2="1130" y2="510"/>'
        '<line x1="1130" y1="510" x2="1080" y2="430"/>'
        '<line x1="1060" y1="560" x2="1080" y2="430"/>'
        '</g>'
        '<g fill="#a88d75" opacity="0.28">'
        '<circle cx="1010" cy="70" r="6"/>'
        '<circle cx="1110" cy="130" r="4"/>'
        '<circle cx="1050" cy="220" r="6"/>'
        '<circle cx="1140" cy="270" r="4"/>'
        '<circle cx="1060" cy="560" r="6"/>'
        '<circle cx="1130" cy="510" r="4"/>'
        '<circle cx="1080" cy="430" r="5"/>'
        '</g>'
    )


def cover_block():
    return (
        f'<rect x="76" y="55" width="390" height="520" fill="none" stroke="{FRAME}" stroke-width="2"/>'
        f'<image x="78" y="57" width="386" height="516" href="data:image/jpeg;base64,{COVER_B64}"/>'
    )


def header(part_label=""):
    s = banner("NETWORKING WITH GO")
    s += banner("MADE EASY - FREE OPEN-SOURCE BOOK", y=128, size=14, ls=3, fill=FAINT)
    s += f'<rect x="{TEXT_X}" y="160" width="72" height="7" fill="{ACCENT}"/>'
    if part_label:
        s += banner(part_label, y=216, size=15, ls=2.5, fill=MUT)
    return s


def footer():
    left = banner("FREE - NO PAYWALL - OPEN SOURCE", y=576, size=15, ls=2, fill=MUT)
    right = (f'<text x="1124" y="576" font-size="15" font-family="{FONT}" '
             f'text-anchor="end" fill="{FAINT}">net/go.book</text>')
    return left + right


def svg_for(title, part_label, body_text, tagline=False):
    tb, nl = title_block(title)
    size = 24 if tagline else 21
    body = body_block(body_text, size, nl)
    return (
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{W}" height="{H}" viewBox="0 0 {W} {H}">'
        f'<rect width="{W}" height="{H}" fill="{BG}"/>'
        f'{motif()}'
        f'{cover_block()}'
        f'{header(part_label)}'
        f'{tb}'
        f'{body}'
        f'{footer()}'
        '</svg>'
    )


def render(svg_text, out_name):
    svg_path = os.path.join(TMP, out_name + ".svg")
    with open(svg_path, "w") as f:
        f.write(svg_text)
    png = os.path.join(OUT, out_name + ".png")
    subprocess.run(["rsvg-convert", "-w", str(W), "-h", str(H), svg_path, "-o", png],
                   check=True, capture_output=True)
    thumb = os.path.join(OUT, out_name + ".thumb.jpg")
    subprocess.run(["magick", png, "-resize", "400x", "-quality", "82", thumb],
                   check=True, capture_output=True)
    return png


def chapter_data():
    """Return ordered list of (part_slug, part_title, [(pos, slug, title, excerpt), ...])."""
    parts = []
    for name in sorted(os.listdir(READ)):
        d = os.path.join(READ, name)
        if not os.path.isdir(d) or name == "topics":
            continue
        idx = os.path.join(d, "index.html")
        if not os.path.exists(idx):
            continue
        with open(idx) as f:
            idx_html = f.read()
        m = H1_RE.search(idx_html)
        part_title = m.group(1).strip() if m else name
        chapters = []
        for fn in sorted(os.listdir(d)):
            if not fn.endswith(".html") or fn == "index.html":
                continue
            slug = fn[:-5]
            with open(os.path.join(d, fn)) as f:
                page = f.read()
            m = H1_RE.search(page)
            title = m.group(1).strip() if m else slug
            dm = DESC_RE.search(page)
            excerpt = html.unescape(dm.group(1)) if dm else ""
            pos = 0
            mm = re.match(r"(\d+)", slug)
            if mm:
                pos = int(mm.group(1))
            chapters.append((pos, slug, title, excerpt))
        chapters.sort(key=lambda c: c[0])
        parts.append((name, part_title, chapters))
    return parts


def build():
    os.makedirs(TMP, exist_ok=True)
    parts = chapter_data()
    manifest = []

    def add_asset(file, label, group):
        if file.endswith(".png") and file.startswith("og-"):
            thumb = file[:-4] + ".thumb.jpg"
        else:
            thumb = file
        manifest.append({"file": file, "thumb": thumb,
                         "label": label, "group": group})

    render(svg_for("Networking with Go, Made Easy", "", (
        "From TCP/IP, sockets and concurrency to gRPC, TLS, ethical hacking and "
        "production APIs in Go - free, open source."), tagline=True), "og-card")
    add_asset("og-card.png", "Book landing card", "Featured")

    render(svg_for("Read Online", "",
                   "Every chapter of the book on its own page, with full-text search - always in sync."),
           "og-read")
    add_asset("og-read.png", "Read Online", "Featured")

    render(svg_for("Browse by Topic", "",
                   "Jump straight to the chapters that matter to you: TCP/IP, HTTP, WebSockets, gRPC, security and more."),
           "og-topics")
    add_asset("og-topics.png", "Browse by Topic", "Featured")

    add_asset("cover.jpg", "Book cover", "Featured")
    add_asset("icon-192.png", "App icon 192", "Icons")
    add_asset("icon-512.png", "App icon 512", "Icons")

    for part_idx, (part_slug, part_title, chapters) in enumerate(parts, 1):
        label = f"PART {part_idx:02d} - {part_title.upper()}"
        part_hook = f"{len(chapters)} chapters: theory first, then runnable Go."
        render(svg_for(part_title, label, part_hook), "og-" + part_slug)
        add_asset("og-" + part_slug + ".png", part_title, "Parts")
        for pos, slug, title, excerpt in chapters:
            ch_label = f"CHAPTER {pos:02d} - {part_title.upper()}"
            render(svg_for(title, ch_label, excerpt), f"og-{part_slug}-{slug}")
            add_asset(f"og-{part_slug}-{slug}.png", title, part_title)

    with open(os.path.join(OUT, "assets-manifest.json"), "w") as f:
        json.dump({"count": len(manifest), "assets": manifest}, f, indent=2)

    zip_name = os.path.join(OUT, "networking-with-go-og-assets.zip")
    with zipfile.ZipFile(zip_name, "w", zipfile.ZIP_DEFLATED) as z:
        z.write(COVER, "cover.jpg")
        for entry in manifest:
            path = os.path.join(OUT, entry["file"])
            if os.path.exists(path):
                z.write(path, entry["file"])

    nchapters = sum(len(p[2]) for p in parts)
    print(f"done: {len(parts)} parts, {nchapters} chapters, "
          f"{len(manifest)} assets, zip={os.path.getsize(zip_name)} bytes")


if __name__ == "__main__":
    with open(COVER, "rb") as f:
        COVER_B64 = base64.b64encode(f.read()).decode()
    build()
