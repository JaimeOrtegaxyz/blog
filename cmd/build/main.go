package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	texttemplate "text/template"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

const (
	contentDir   = "content"
	templatesDir = "templates"
	staticDir    = "static"
	distDir      = "dist"

	siteTitle = "Jaime Ortega"
	siteURL   = "https://jaimeortega.xyz"
	siteEmail = "jaime.ortc@gmail.com"
	siteDesc  = "Musings on past, present, and future."
)

var basePath string

type Section struct {
	Slug, Title, Desc string
}

var sections = []Section{
	{"thoughts", "Thoughts", "Long-form essays."},
	{"notes", "Notes", "Short takes and observations."},
	{"projects", "Projects", "Things I've built."},
	{"reading", "Reading", "Books and articles."},
}

type Post struct {
	Section Section
	Slug    string
	Title   string
	Date    time.Time
	Summary string
	HTML    template.HTML
}

func (p Post) URL() string     { return fmt.Sprintf("%s/%s/%s.html", basePath, p.Section.Slug, p.Slug) }
func (p Post) AbsURL() string  { return siteURL + p.URL() }
func (p Post) RFC822() string  { return p.Date.Format(time.RFC1123Z) }
func (p Post) DateStr() string { return p.Date.Format("2006-01-02") }

type Site struct {
	Title, URL, Email, Desc, BasePath string
	Sections                          []Section
}

func site() Site {
	return Site{Title: siteTitle, URL: siteURL, Email: siteEmail, Desc: siteDesc, BasePath: basePath, Sections: sections}
}

func main() {
	start := time.Now()
	basePath = os.Getenv("BLOG_BASE_PATH")

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	all, bySection := discover(md)

	mustReset(distDir)

	xmlT := texttemplate.Must(
		texttemplate.New("").Funcs(texttemplate.FuncMap{"xml": xmlEscape}).
			ParseGlob(filepath.Join(templatesDir, "*.xml")),
	)

	renderStandalone("home.html", filepath.Join(distDir, "index.html"), map[string]any{
		"Site":      site(),
		"Sections":  sections,
		"Wallpaper": findWallpaper(),
	})

	for _, sec := range sections {
		ps := bySection[sec.Slug]
		secDir := filepath.Join(distDir, sec.Slug)
		mustMkdir(secDir)

		renderHTML("section.html", filepath.Join(secDir, "index.html"), map[string]any{
			"Site":    site(),
			"Section": sec,
			"Posts":   ps,
		})

		renderXML(xmlT, filepath.Join(secDir, "feed.xml"),
			feedData(sec.Title+" — "+siteTitle, sec.Desc, basePath+"/"+sec.Slug+"/feed.xml", ps))

		for _, p := range ps {
			renderHTML("post.html", filepath.Join(secDir, p.Slug+".html"), map[string]any{
				"Site": site(),
				"Post": p,
			})
		}
	}

	renderXML(xmlT, filepath.Join(distDir, "feed.xml"),
		feedData(siteTitle, siteDesc, basePath+"/feed.xml", all))

	renderHTML("about.html", filepath.Join(distDir, "about.html"), map[string]any{
		"Site": site(),
	})

	if err := copyDir(staticDir, filepath.Join(distDir, "static")); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("built %d posts, %d sections in %s\n",
		len(all), len(sections), time.Since(start).Round(time.Microsecond))
}

func discover(md goldmark.Markdown) ([]Post, map[string][]Post) {
	var all []Post
	by := map[string][]Post{}
	for _, sec := range sections {
		dir := filepath.Join(contentDir, sec.Slug)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			log.Fatalf("read %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			path := filepath.Join(dir, e.Name())
			p, err := parsePost(path, sec, md)
			if err != nil {
				log.Fatalf("%s: %v", path, err)
			}
			all = append(all, p)
			by[sec.Slug] = append(by[sec.Slug], p)
		}
	}
	desc := func(ps []Post) {
		sort.Slice(ps, func(i, j int) bool { return ps[i].Date.After(ps[j].Date) })
	}
	desc(all)
	for k := range by {
		desc(by[k])
	}
	return all, by
}

func parsePost(path string, sec Section, md goldmark.Markdown) (Post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}
	fm, body, err := splitFrontMatter(raw)
	if err != nil {
		return Post{}, err
	}
	if fm["title"] == "" {
		return Post{}, fmt.Errorf("missing title")
	}
	if fm["date"] == "" {
		return Post{}, fmt.Errorf("missing date")
	}
	date, err := parseDate(fm["date"])
	if err != nil {
		return Post{}, fmt.Errorf("bad date %q: %w", fm["date"], err)
	}
	var buf bytes.Buffer
	if err := md.Convert(body, &buf); err != nil {
		return Post{}, err
	}
	return Post{
		Section: sec,
		Slug:    strings.TrimSuffix(filepath.Base(path), ".md"),
		Title:   fm["title"],
		Date:    date,
		Summary: fm["summary"],
		HTML:    template.HTML(buf.String()),
	}, nil
}

func splitFrontMatter(raw []byte) (map[string]string, []byte, error) {
	r := bufio.NewReader(bytes.NewReader(raw))
	first, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(first) != "---" {
		return nil, nil, fmt.Errorf("expected --- on first line")
	}
	fm := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, nil, fmt.Errorf("unexpected EOF in front matter")
		}
		if strings.TrimSpace(line) == "---" {
			break
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		fm[key] = val
	}
	body, _ := io.ReadAll(r)
	return fm, body, nil
}

func parseDate(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02 15:04"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown format")
}

type FeedItem struct {
	Title, URL, Summary, Date string
}

type FeedDoc struct {
	Title, URL, FeedURL, Desc, BuildDate string
	Items                                []FeedItem
}

func feedData(title, desc, feedPath string, ps []Post) FeedDoc {
	doc := FeedDoc{
		Title:     title,
		URL:       siteURL,
		FeedURL:   siteURL + feedPath,
		Desc:      desc,
		BuildDate: time.Now().Format(time.RFC1123Z),
	}
	for _, p := range ps {
		doc.Items = append(doc.Items, FeedItem{
			Title:   p.Title,
			URL:     p.AbsURL(),
			Summary: p.Summary,
			Date:    p.RFC822(),
		})
	}
	return doc
}

func xmlEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func renderHTML(page, out string, data any) {
	t := template.Must(template.ParseFiles(
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, page),
	))
	mustMkdir(filepath.Dir(out))
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := t.ExecuteTemplate(f, "base", data); err != nil {
		log.Fatalf("render %s: %v", out, err)
	}
}

func renderStandalone(page, out string, data any) {
	t := template.Must(template.ParseFiles(filepath.Join(templatesDir, page)))
	mustMkdir(filepath.Dir(out))
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := t.Execute(f, data); err != nil {
		log.Fatalf("render %s: %v", out, err)
	}
}

func findWallpaper() string {
	for _, ext := range []string{".webp", ".avif", ".jpg", ".jpeg", ".png"} {
		path := filepath.Join(staticDir, "wallpaper"+ext)
		if _, err := os.Stat(path); err == nil {
			return basePath + "/" + filepath.ToSlash(path)
		}
	}
	return ""
}

func renderXML(t *texttemplate.Template, out string, data any) {
	mustMkdir(filepath.Dir(out))
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := t.ExecuteTemplate(f, "feed.xml", data); err != nil {
		log.Fatalf("render %s: %v", out, err)
	}
}

func mustReset(d string) {
	if err := os.RemoveAll(d); err != nil {
		log.Fatal(err)
	}
	mustMkdir(d)
}

func mustMkdir(d string) {
	if err := os.MkdirAll(d, 0755); err != nil {
		log.Fatal(err)
	}
}

func copyDir(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		mustMkdir(filepath.Dir(target))
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}
