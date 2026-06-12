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
- `templates/` — HTML templates. Every page extends `base.html`
  (including `home.html`, the landing: intro + recent posts per
  section).
- `content/<section>/*.md` — posts. Front matter is `title`, `date`,
  optional `summary`. Filename = slug = URL path component.
- `static/` — `style.css`, `fonts/plex-*.woff2`, and `art/<slug>/`
  images for art posts (see "Art posts" in the spec).
- `dist/` — build output, gitignored.

## Common operations

- `make build` — generate `dist/`.
- `make serve` — build and serve on port 8000 via `python3 -m
  http.server`.
- New post: drop `content/<section>/<slug>.md` with front matter.
- New section: edit the `sections` slice in `cmd/build/main.go` and
  `mkdir content/<name>/`. The home page and nav pick it up
  automatically. Set the section's `Stream` flag to render its index
  as a continuous stream (like Notes) instead of a card list.

## Voice when responding to feature requests

The site rejects most things by default — see *"What is deliberately
absent"* and *"Common temptations and the answer"* in the spec. If a
request matches one of those, surface the default answer with the
reason rather than starting to implement.
