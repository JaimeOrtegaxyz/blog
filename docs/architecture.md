# Site architecture & principles

This document is the project's north star: how the site is built, what
it must remain, and what it must not become. It exists to do three
things —

1. Onboard future contributors (or future you) without making them
   reconstruct the reasoning from the source.
2. Provide a checklist against which any proposed change can be
   evaluated before it ships.
3. Serve as the input to an LLM tasked with auditing drift —
   *"contrast the current state of this repo against this document and
   report where it has diverged."*

If you are an LLM reading this in a future session, treat the rules
below as constraints, not suggestions. The project has been deliberately
built to refuse churn. Any temptation to "modernise" something here
should be re-read against the relevant principle first.

---

## Core philosophy

The site rests on five claims:

**1. The browser is the runtime.**
Everything else is optional. HTML, CSS, and (rarely) a small bit of
inline JavaScript are the substrate. There is no client-side runtime
framework, no hydration, no virtual DOM, no client-side router. Pages
are real documents, requested by URL and rendered by the browser.

**2. Hand-crafted beats off-the-shelf when the off-the-shelf is bigger
than the problem.**
A static site generator does one job: walk a folder of Markdown,
produce a folder of HTML. The mainstream tools (Hugo, Jekyll, 11ty,
Astro, Next, Gatsby) bundle a job's worth of features and a multiple
of the code. The custom Go SSG in this repo is roughly 300 lines and
solves the problem completely.

**3. Standards over abstractions.**
Native HTML elements over custom components. Native CSS over Tailwind
or PostCSS. Native browser APIs (View Transitions, `prefers-color-
scheme`, `prefers-reduced-motion`, `<link rel="prefetch">`) over
JavaScript polyfills. The web platform has matured; reach for it
before reaching for a library.

**4. Speed is correctness.**
A slow build is a bug. A heavy page is a bug. A page that flashes,
re-flows, or hangs while a font loads is a bug. The metrics matter:
we measure them and we don't let them drift.

**5. Subtraction is a feature.**
Adding feels productive; removing rarely does. But removing is
usually how this site gets better, because the goal is the *minimum
machinery* that delivers the writing. If a feature can be cut without
the writing suffering, it should be considered for cutting.

---

## Non-negotiables

These are the concrete checkpoints. Each is phrased so an LLM can
verify it mechanically.

### Build & dependency surface

1. **Build time under one second.**
   `time make build` should complete the actual SSG run in single-
   digit milliseconds. The wall-clock cost on first run is `go run`
   compilation, which is acceptable. There is no incremental build,
   no cache, no watcher. None is needed.

2. **Single direct Go dependency.**
   `go.mod` lists exactly one direct dependency: `github.com/yuin/
   goldmark` (markdown → HTML). No additional deps are permitted
   without a clear, articulated reason that survives the question
   *"could we just write that ourselves in 30 lines?"*

3. **No `package.json` at the repo root, ever.**
   The site is not a Node project. It does not have a build toolchain
   in JavaScript. There must be no `package.json`, no `node_modules/`,
   no `bun.lockb`, no `pnpm-lock.yaml`, no `yarn.lock` in the project
   root.

4. **No CSS preprocessor.**
   No Tailwind, no PostCSS, no Sass, no Less, no CSS-in-JS. The
   stylesheet at `static/style.css` is hand-written and human-
   readable.

5. **No CMS, no headless CMS, no database.**
   Content is markdown files in `content/<section>/<slug>.md`. The
   filesystem is the CMS. No Contentful, no Sanity, no Notion bridge,
   no Strapi, no Supabase, no Postgres.

### Output budget

6. **Per-page weight target: under ~30KB for every page.**
   Every page (`/`, `/<section>/`, `/<section>/<slug>.html`,
   `/about.html`) must transfer under ~30KB total including HTML,
   CSS, and font. Posts that include art (see "Art posts" below)
   are exempted for their images only: each image is one
   hand-compressed file (webp / avif / jpg) under 200KB, and the
   page's budget excluding images still holds.

7. **Total `dist/` size stays in the low hundreds of KB for a
   small site.**
   With ~20 posts, expect somewhere in the 100–200KB range, plus
   whatever art images a handful of posts carry. If it crosses
   500KB excluding art, something has been added that probably
   shouldn't have.

8. **Zero external network requests.**
   In devtools' network panel, every request must be same-origin.
   No `fonts.googleapis.com`, no `cdn.jsdelivr.net`, no analytics
   beacon, no script CDN. Custom fonts are self-hosted woff2 files
   in `static/fonts/`. Icons are inlined SVG.

9. **At most one CSS request, one font request, zero JS requests
   per page.**
   Every page: 1 HTML + 1 CSS + 1 font + 0 JS + 0 third-party.
   Art posts: same plus N same-origin images, each within the
   200KB budget above.

10. **Zero JavaScript by default.**
    No `<script>` tag should appear in any rendered HTML unless it
    delivers a feature that cannot be done with HTML and CSS alone.
    Browser APIs (View Transitions, anchor navigation, system
    `prefers-color-scheme`, system `prefers-reduced-motion`,
    `<link rel="prefetch">`) replace almost every conventional
    JavaScript use case.

### Markup, layout, and content

11. **Every rendered page must be valid HTML5 and human-readable in
    `view-source`.**
    No minified blobs. No `__NEXT_DATA__` script tag. No data
    smuggled as base64. No hashed asset URLs unless a build tool
    needed them, which it doesn't.

12. **Plain `<a href>` to plain `.html` files.**
    No client-side router. No `Link` component. No
    `data-router-link` shenanigans. The browser handles navigation.

13. **No tags, no categories, no full-text search.**
    Sections are the only taxonomy. They partition content into
    `thoughts/`, `notes/`, `projects/`, `reading/`, plus `about`.
    No nested taxonomies, no tag clouds, no search index, no
    Algolia, no Pagefind. Browser find (Cmd-F) is the search.

14. **No comments, no share buttons, no related-posts widget, no
    "estimated reading time," no newsletter modal, no cookie banner.**
    Each of these is a separate temptation; each should be refused
    by default. Email is the channel for direct contact; RSS is the
    channel for subscription.

15. **No image CDN, no responsive `srcset` framework.**
    Sized `<img>`. Every image is one carefully-chosen,
    appropriately-compressed file.

### Type and visual

16. **One typeface across the whole site: IBM Plex Mono.**
    Self-hosted regular and italic woff2 only. Latin subset only.
    No web font request to a third-party.

17. **All headings render at `font-weight: 400`.**
    No bold variants. The visual hierarchy comes from size,
    spacing, and italic — not weight contrast.

18. **Colour palette derives from a small set of CSS variables.**
    Light and dark are auto-selected via `prefers-color-scheme`.
    There is no theme toggle UI. The user's OS preference is
    authoritative.

19. **Animations are CSS-only.**
    No JS-driven animation libraries (Framer Motion, GSAP,
    Anime.js, Lottie, Three.js). Page transitions use the
    cross-document View Transitions API. Per-element fades use
    `@keyframes`. Scroll-linked motion uses CSS scroll-driven
    animations (`view-timeline` / `animation-timeline`). All
    motion — transitions, fades, view-transition cross-fades,
    scroll-linked animation — is gated on
    `prefers-reduced-motion`.

### Project layout

20. **The SSG is one file: `cmd/build/main.go`.**
    Three loops: discover → render → write. No goroutine pool, no
    incremental cache, no plugin system, no theme system. If the
    file passes ~500 lines, that is a smell — something has grown
    that probably shouldn't have.

21. **Templates in `templates/`, content in `content/`, static
    assets in `static/`. Output goes to `dist/`.**
    These four directories should be all the project has at the
    root, plus `cmd/build/`, `Makefile`, `go.mod`, `go.sum`,
    `.gitignore`, `docs/`, `CLAUDE.md`, and any auxiliary repo-wide
    config. Anything else at root is a smell.

---

## Architecture

### Static site generator

Located at `cmd/build/main.go`. ~400 lines. One direct Go dependency:
`github.com/yuin/goldmark` for markdown → HTML.

The pipeline is unforgivingly simple. The `main()` function:

1. **Discovers** all posts by walking each `content/<section>/`
   directory, parsing the `---`-delimited YAML-ish front matter
   (title, date, summary), and converting the markdown body to HTML.
2. **Renders** the home page, each section's index, RSS feed, and
   posts, the combined feed, and the about page.
3. **Writes** every output to `dist/`.

Every page extends `templates/base.html` via Go template inheritance
(`{{define "content"}}` / `{{template "base"}}`) — including the home
page. There is one shared chrome (brand link, section nav, footer)
and no special-cased documents.

### Sections

A section is a directory under `content/`. The list of sections is
declared once, at the top of `cmd/build/main.go`, in a `[]Section`
slice. Adding a section means: create the directory, add a struct
literal to the slice, drop in markdown files. The home page and the
nav pick it up automatically. There is no plugin to install, no
config file to edit, no theme to extend.

The current sections are `thoughts`, `notes`, `projects`, `reading`,
plus `about` (a single page rather than a section, linked from the
nav). These names are not sacred; they are the current choice.
Renaming requires updating the `Section` slice and renaming the
content directory.

### Content

Markdown files at `content/<section>/<slug>.md`. The slug is the
filename. The URL is `/<section>/<slug>.html`. The filesystem layout
*is* the URL layout. There are no rewrites, no `_index.md`, no
`_drafts/` directory.

Front matter is between `---` lines at the top of the file:

```
---
title: Hello, world
date: 2026-05-04
summary: Optional one-liner shown as the card excerpt and post subtitle.
---
```

Three fields: `title`, `date`, `summary`. No more. There is no
`tags`, `category`, `author`, `cover`, `layout`, `permalink`,
`canonical`, `redirect_from`, `weight`, `draft`, or any of the other
fields a fuller SSG would invite. Every additional field is a
permanent contract on every future post.

### Output

A built page is composed of:

- One HTML file (1–5KB).
- One shared CSS file at `/static/style.css` (~11KB).
- One font file (`plex-regular.woff2`, ~15KB), preloaded via
  `<link rel="preload">`.
- (Art posts only.) N self-hosted images from `static/art/<slug>/`,
  each one hand-compressed file under 200KB.

That is the whole network footprint. Nothing else loads.

### Routing

Plain `<a href>` to plain `.html` files. No client-side router.

- `/` — the landing: a short intro plus the most recent posts from
  each section.
- `/<section>/` — chronological index for the section.
- `/<section>/<slug>.html` — one file per post.
- `/<section>/feed.xml` — RSS feed for the section.
- `/feed.xml` — combined RSS feed across all sections.
- `/about.html` — about page.

Every page uses the shared `site-header` with the brand link and
section nav. The brand link points to `/`.

### Hosting

The output of `make build` is a folder of static files. It can be
served by nginx, Caddy, GitHub Pages, Cloudflare Pages, or any
ordinary static host. There is no serverless function, no edge
runtime, no ISR/SSR/RSC distinction, no build-time webhook, no
deploy preview daemon. Push the folder; that is the deploy.

For hosts that serve the site from a sub-path (GitHub Pages project
sites), set `BLOG_BASE_PATH` at build time and every root-absolute
URL is prefixed with it. The deploy workflow in
`.github/workflows/` does this. It is a build-time string prefix,
not runtime configuration.

---

## Art posts

A post may carry artwork beside its text, in the manner of an
illustrated magazine: the image pins to the left half of the viewport
while its stretch of text scrolls on the right, and the next image
takes over when its stretch arrives.

The mechanic is CSS only, and must stay that way. At build time the
SSG (`panelize()` in `cmd/build/main.go`) wraps each image paragraph
plus the text that follows it into a `<section class="panel">`;
`position: sticky` does the pinning and the handoff. No
IntersectionObserver, no scroll listener, no scrollytelling library.
If a future change to this feature requires a `<script>` tag, the
change is wrong.

A dotted track runs down the image/text seam, with a dot that rides
the article's own scroll timeline as you read — CSS scroll-driven
animation, no JavaScript. Like all motion on the site it is gated on
`prefers-reduced-motion`, and browsers without `animation-timeline`
never show it at all.

Authoring rules:

- An image is a panel boundary only when the markdown paragraph is
  *exactly one image* (`![alt](/static/art/<slug>/file.jpg)` on its
  own line). Images inline within a sentence are left untouched.
- Text before the first image joins the first image's panel, so the
  artwork is on screen from the top of the page — never an empty
  rail beside the opening paragraphs.
- Image files live in `static/art/<slug>/`. Self-hosted, like
  everything else. No image CDN, no srcset pipeline, no on-demand
  resizing.
- Each image is one hand-compressed file (webp / avif / jpg) under
  200KB. Compress before committing; the build will not do it for
  you, deliberately.
- No new front matter. Whether a post has art is inferred from its
  body, not declared.
- Every art panel is at least 1.5 viewports tall (`min-height` in
  CSS), so each image pins for at least half a viewport before
  handing off. Beyond that floor, sticky travel is the text's
  length, not a timer: longer stretches pin longer. Give each image
  enough words to be worth pinning for.

Posts without image paragraphs render exactly as before — centered
single column, byte-identical output. On narrow screens art posts
collapse to a single column with images inline, in order.

---

## What is deliberately absent

Things that look obvious or even necessary in a typical 2026 site,
and that this project rejects. The list is long because the
temptations are many:

- **Frameworks**: React, Vue, Svelte, Solid, Lit, Stencil, Qwik.
- **Static-site frameworks**: Next.js, Remix, Gatsby, Nuxt, Astro,
  Hugo, Jekyll, Eleventy.
- **Styling tools**: Tailwind, PostCSS, Sass, Less, Stylus,
  CSS-in-JS, styled-components, Emotion.
- **Bundlers and build tools**: webpack, Rollup, Parcel, Vite,
  esbuild, Turbopack, Bun-as-a-build-tool.
- **CMSs and headless CMSs**: Contentful, Sanity, Strapi, Storyblok,
  Notion-as-CMS, WordPress.
- **Hosting platforms with magic**: Vercel-specific features (ISR,
  edge functions, image optimisation), Netlify Functions, Cloudflare
  Workers — the site runs on plain static hosting.
- **Data layers**: GraphQL, tRPC, Supabase, Firebase, Prisma,
  Drizzle.
- **Animation libraries**: Framer Motion, GSAP, Anime.js, Lottie,
  Three.js, Motion One.
- **Component libraries**: shadcn/ui, Radix, Headless UI, Material
  UI, Chakra, Mantine.
- **Icon libraries**: Lucide, Heroicons, Font Awesome,
  react-icons. Icons are inlined SVG, hand-drawn or hand-picked.
- **Tracking and analytics**: Google Analytics, Plausible, Umami,
  Mixpanel, PostHog. Server logs aggregate to roughly the same
  metrics without sending visitor data to a third party.
- **Comments**: Disqus, Giscus, Utterances, custom comment systems.
  Email is the channel.
- **Newsletter and signup**: ConvertKit, Substack, Buttondown,
  Mailchimp, modal popups asking to subscribe. RSS is the
  subscription channel.
- **Search**: Algolia, Pagefind, lunr, fuse.js, in-browser search
  index. Browser find (Cmd-F) is the search.
- **Content extras**: tags, categories, author profiles, related
  posts, estimated reading time, share buttons, OG-image generators,
  comment counts, view counters, "edit on GitHub" links.
- **Accessibility add-ons**: skip-to-content links and proper semantic
  HTML are fine; "accessibility *libraries*" (toolbar widgets that
  promise WCAG compliance) are not.

A live "site upgrade" idea will usually fall into one of these
buckets. The default answer is "we considered that and chose not to."

---

## Drift signals — checklist for audits

Concrete patterns that, if they appear, indicate the project is
moving away from its principles. Use this section as a mechanical
checklist when auditing.

### File system

- [ ] A `package.json` has appeared at the repo root.
- [ ] A `node_modules/` directory exists anywhere.
- [ ] A `.env`, `.env.local`, or `next.config.*` exists at root.
- [ ] More than ~6 directories at root (excluding `.git`,
  `.claude/`, hidden config).
- [ ] A `vercel.json`, `netlify.toml`, or other platform-specific
  config has appeared.
- [ ] A `vite.config.*`, `webpack.config.*`, `tailwind.config.*`,
  `postcss.config.*`, `tsconfig.*` exists.
- [ ] A `bun.lockb`, `pnpm-lock.yaml`, `yarn.lock`, or
  `package-lock.json` exists.

### Dependencies

- [ ] `go.mod` has more than one direct dependency.
- [ ] An external font URL appears anywhere in CSS or HTML
  (`fonts.googleapis.com`, `cdn.jsdelivr.net/...`, etc.).
- [ ] A CDN URL for any asset other than a developer-explicit
  exception (none currently exist).

### Source code

- [ ] `cmd/build/main.go` has crossed ~500 lines.
- [ ] The SSG has gained a "plugin" or "theme" abstraction.
- [ ] A second Go program exists in `cmd/`.
- [ ] A `lib/` or `pkg/` directory has appeared with non-trivial
  abstractions.
- [ ] CSS has gained a preprocessor of any kind.
- [ ] Templates have multiplied beyond the small fixed set in
  `templates/`.

### Markup

- [ ] Any rendered HTML page contains a `<script>` tag that is not
  strictly necessary.
- [ ] Any rendered HTML page imports an external script.
- [ ] A page references a hashed asset URL like
  `/assets/style.a1b2c3d4.css` (suggests a bundler).
- [ ] `view-source` on any page shows minified output.
- [ ] A `__NEXT_DATA__`, `__NUXT__`, or similar script tag exists.
- [ ] A `<link rel="manifest">` for a PWA appears (probably not
  needed; if it is, must be justified).

### Performance budgets

- [ ] `time make build` (after compilation) takes more than 100ms
  for the SSG run itself.
- [ ] Any page's transfer size exceeds 30KB (excluding art-post
  images that individually respect the 200KB budget).
- [ ] An image under `static/art/` exceeds 200KB.
- [ ] The CSS file exceeds ~12KB.
- [ ] More than one font file is downloaded per page.
- [ ] Any third-party network request appears in devtools.

### Content surface

- [ ] A new front-matter field has appeared on posts.
- [ ] A `tags:` or `categories:` field has appeared.
- [ ] A `_drafts/` or `archive/` directory has appeared.
- [ ] A "related posts" component has been added.
- [ ] An "estimated reading time" calculation has been added.
- [ ] A comments section, share buttons, or newsletter signup has
  been added.
- [ ] A search input or search index has been added.

### Visual

- [ ] A second typeface has been added.
- [ ] A bold (`font-weight: 700`) variant has been pulled in.
- [ ] A theme toggle has appeared (it is unnecessary; OS
  `prefers-color-scheme` already works).
- [ ] An icon library or icon font (Font Awesome, Lucide, etc.)
  has been imported.
- [ ] An animation library has been imported.

If three or more items are checked, the project is meaningfully
off-spec. Even one or two should prompt a review.

---

## Common temptations and the answer

A short FAQ for the requests this document expects to see most often:

- **"Can we add a search box?"** — No. Browser find is the search. If
  the corpus genuinely outgrows that, the corpus is too big and
  should be split.

- **"Can we add tags?"** — No. Sections are the only taxonomy. If a
  post seems to span sections, pick the better fit.

- **"Can we add comments?"** — No. Email is the channel.

- **"Can we add a theme toggle?"** — No. The OS already has one, and
  the site already respects it via `prefers-color-scheme`.

- **"Can we add Tailwind to make styling faster?"** — No. The
  stylesheet is small enough that hand-writing it is faster than
  installing and configuring anything.

- **"Can we move to Astro/Next/Hugo?"** — No. The custom SSG is
  smaller than any of those, builds faster, and has zero churn risk.

- **"Can we add an animation library for nicer transitions?"** —
  No. View Transitions and `@keyframes` are sufficient.

- **"Can we add Plausible/Umami so we can see traffic?"** — No.
  Server logs aggregate to roughly the same metric without sending
  visitor data to a third party.

- **"Can we add a related-posts widget?"** — No. The section index
  is the related-posts widget.

- **"Can we add an OG image generator?"** — Possibly, only if it can
  be done at build time without adding a runtime dependency, and the
  output is a static `.png` written to `dist/` like any other asset.
  No serverless, no on-demand generation.

- **"Can we add code syntax highlighting?"** — Possibly. Goldmark's
  built-in `extension/highlighting` (Chroma) renders highlighted
  HTML at build time without any client-side JS. This is the only
  addition on this list that genuinely improves the writing surface,
  so it is worth considering. If added, it must remain server-
  rendered HTML — never a client-side highlighter like Prism or
  highlight.js. Adding it costs one transitive Go dependency
  (Chroma); the cost should be weighed against the benefit.

---

## How to use this document

### As a human contributor

Before adding a feature, ask:

1. Does this cross any of the non-negotiables?
2. Does this introduce a dependency, a new file at root, or a
   recurring runtime cost?
3. Could this be done with native HTML, CSS, and the existing Go
   pipeline?
4. If removed in a year, would the writing be worse?

If the answers indicate friction, the change probably should not
ship. If you ship it anyway, document the exception in this file.

### As an LLM tasked with auditing this project

Read this entire file before reading the source. Then, when asked to
"audit drift" or "check whether the project still matches the spec":

1. Walk every "Non-negotiables" rule and verify it against the
   current state of the repo. Use `grep`, `wc`, `du`, and `find` to
   get concrete numbers. Don't guess; measure.
2. Walk the "Drift signals" checklist mechanically. Report each
   item as ✅ matched or ❌ violated, with the file path and line
   that triggered the violation.
3. For each violation, propose the smallest change that restores
   compliance. Do not propose larger refactors unless the violation
   indicates a structural problem.
4. Report findings as a short, direct list. Do not soften the
   verdict. The project's whole stance is that drift is bad and
   that incremental drift compounds.

If the project owner asks "how can we add X," and X is on the
"deliberately absent" list or in the "Common temptations" section,
the default answer is "we don't, and here is why." Only override
that default with a clear justification that the owner accepts.

---

## On voice

This document has a deliberate tone: direct, specific, opinionated.
That tone is part of the spec. A "we might consider..." or "it could
be worth exploring..." version of this document would invite drift.
A flat "no, here is why" version refuses it.

Keep it that way.
