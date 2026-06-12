---
title: An example note
date: 2026-05-03
summary: Notes are short. Replace this with your own.
---

A note is shorter than a thought. A paragraph. A link. A small
observation.

This one doubles as a demo of art panels: a markdown paragraph that is
exactly one image becomes a panel boundary, and the image pins to the
left of the viewport while the text beside it scrolls. One rule to
keep in mind — an image only stays pinned while its stretch of text is
taller than the viewport. Short stretch, short pin. Give each image
enough words to be worth pinning for.

![Demo panel one](/static/art/example-note/1.jpg)

The first image is on screen from the very top — any text written
before it simply shares its panel. On a desktop viewport it sits
sticky in the left half while this chunk of text scrolls past it on
the right. There is no JavaScript involved — the pinning and the
handoff to the next image are `position: sticky` doing what it was
designed to do.

The panel is a CSS grid with the figure in the left column and the
text in the right. The figure is one viewport tall and sticky at the
top, so once it reaches the top of the screen it stays there while the
rest of its column scrolls underneath the fold.

This is the same trick magazine sites use a scroll library for. The
browser already knows how to do it. The element sticks within its
container; when the container's end arrives, the element is released
and travels off-screen with it. That release is the handoff.

What that means for writing: the image and its text are one unit. The
text is the duration of the image. A long passage gives the reader a
long stretch with the artwork held in place beside it, which is
exactly the feeling this layout is for — the picture keeps you company
while you read.

It also means the layout degrades honestly. If a stretch of text is
short, the image simply scrolls along with it, like an ordinary
illustration. Nothing breaks, nothing jumps, nothing needs measuring
or observing. The behaviour is proportional to the writing.

Keep scrolling. When this stretch runs out, the next image slides up
and takes the left pane over from this one.

Notice that a handoff happens at a paragraph boundary you chose
when writing — wherever the image line sits in the markdown is where
the previous image's territory ends. There is no configuration for
this, no front matter, no shortcode. The position of the image in the
text *is* the configuration.

Inline images inside a sentence are left alone — only paragraphs that
are exactly one image become panels. So you can still drop a small
diagram mid-paragraph in an ordinary post and nothing about this
layout will trigger.

A second thing worth noticing: the text column keeps the same measure
as every other page on the site. The split layout widens the page, not
the prose. Reading width is a constant here, not a casualty of the
design.

And because each panel is at least one viewport tall, even a shorter
stretch like the next one gets a full screen of presence before the
page moves on.

![Demo panel three](/static/art/example-note/3.jpg)

And a third handoff. On narrow screens all of this collapses to a
single column with the images inline, in order — no sticky, no split,
just an illustrated note read top to bottom.

Drop new files into `content/notes/` to add more. Replace the images
in `static/art/example-note/` with your own — hand-compressed, each
under 200KB.
