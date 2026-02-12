package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var funcMap = template.FuncMap{
	"assetPath": assetPath,
}

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(
		"templates/layout.html",
		"templates/"+tmpl+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages/FrontPage", http.StatusFound)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	if !validTitle(title) {
		http.NotFound(w, r)
		return
	}

	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/pages/"+title+"/edit", http.StatusFound)
		return
	}

	data := struct {
		Title string
		Body  template.HTML
	}{
		Title: page.Title,
		Body:  linkWikiWords(page.Body),
	}
	renderTemplate(w, "view", data)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	if !validTitle(title) {
		http.NotFound(w, r)
		return
	}

	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	if !validTitle(title) {
		http.NotFound(w, r)
		return
	}

	body := r.FormValue("body")
	page := &Page{Title: title, Body: body}
	if err := page.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pages/"+title, http.StatusFound)
}

func main() {
	storagePath = os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	if err := os.MkdirAll(storagePath, 0o755); err != nil {
		log.Fatal(err)
	}

	if err := buildAssetMap(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", rootHandler)
	mux.HandleFunc("GET /pages/{title}", viewHandler)
	mux.HandleFunc("GET /pages/{title}/edit", editHandler)
	mux.HandleFunc("POST /pages/{title}", saveHandler)
	mux.HandleFunc("GET /public/", assetHandler)
	mux.HandleFunc("GET /up", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	var csrf http.CrossOriginProtection

	fmt.Println("Stiki running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, csrf.Handler(mux)))
}
