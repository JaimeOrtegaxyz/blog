# Project notes for Claude

This is a small static site built from a custom Go SSG. It is governed
by an explicit, opinionated architecture spec.

**Before making any changes, read `docs/architecture.md`.** It lists
the non-negotiables (build time, dependency surface, page weight,
what's deliberately absent) and a checklist of drift signals to watch
for. The spec is load-bearing: many feature requests that look small
in isolation are answered there with a default "no" and a reason.

If the user asks you to "audit drift" or "check the project against
the spec," follow the explicit instructions under
*"How to use this document → As an LLM tasked with auditing this
project"* in `docs/architecture.md`. Walk the checklists mechanically;
measure, don't guess; don't soften the verdict.

## Quick map

- `cmd/build/main.go` — the entire SSG, ~300 lines, one direct dep
  (goldmark).
- `templates/` — HTML templates. `home.html` is standalone (the
  desktop landing); every other page extends `base.html`.
- `content/<section>/*.md` — posts. Front matter is `title`, `date`,
  optional `summary`. Filename = slug = URL path component.
- `static/` — `style.css`, `fonts/plex-*.woff2`, optional
  `wallpaper.<ext>` (webp / avif / jpg / jpeg / png).
- `dist/` — build output, gitignored.

## Common operations

- `make build` — generate `dist/`.
- `make serve` — build and serve on port 8000 via `python3 -m
  http.server`.
- New post: drop `content/<section>/<slug>.md` with front matter.
- New section: edit the `sections` slice in `cmd/build/main.go`,
  `mkdir content/<name>/`, add a matching icon `<li>` block in
  `templates/home.html`.

## Voice when responding to feature requests

The site rejects most things by default — see *"What is deliberately
absent"* and *"Common temptations and the answer"* in the spec. If a
request matches one of those, surface the default answer with the
reason rather than starting to implement.
