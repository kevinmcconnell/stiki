package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"regexp"

	"github.com/yuin/goldmark"
)

var (
	storagePath     string
	titlePattern    = regexp.MustCompile(`^[A-Za-z]+$`)
	wikiWordPattern = regexp.MustCompile(`[A-Z][a-z]+([A-Z][a-z]+)+`)
)

type Page struct {
	Title string
	Body  string
}

func validTitle(title string) bool {
	return titlePattern.MatchString(title)
}

func loadPage(title string) (*Page, error) {
	body, err := os.ReadFile(filepath.Join(storagePath, title))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}

func (p *Page) save() error {
	dest := filepath.Join(storagePath, p.Title)

	f, err := os.CreateTemp(storagePath, p.Title+".*")
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(p.Body)); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return err
	}

	return os.Rename(f.Name(), dest)
}

func existingPages() map[string]bool {
	pages := map[string]bool{}
	entries, err := os.ReadDir(storagePath)
	if err != nil {
		return pages
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			pages[entry.Name()] = true
		}
	}
	return pages
}

func linkWikiWords(raw string) template.HTML {
	pages := existingPages()

	linked := wikiWordPattern.ReplaceAllStringFunc(raw, func(word string) string {
		if pages[word] {
			return "[" + word + "](/pages/" + word + ")"
		}
		return "[" + word + "?](/pages/" + word + "/edit)"
	})

	var buf bytes.Buffer
	goldmark.Convert([]byte(linked), &buf)
	return template.HTML(buf.String())
}
