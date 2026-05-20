package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// the interface where handlers render templates
type TemplateRenderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
	RenderError(w http.ResponseWriter, message string, code int)
}

type templateRenderer struct {
	cache map[string]*template.Template
	mu    sync.RWMutex
}

// parse every template file paired with layout.html
func NewTemplateRenderer(dir string) (TemplateRenderer, error) {
	pages := []string{
		"index.html",
		"login.html",
		"register.html",
		"admin_register.html",
		"admin.html",
		"patient_slots.html",
		"book.html",
		"error.html",
	}

	layoutPath := filepath.Join(dir, "layout.html")
	cache := make(map[string]*template.Template, len(pages))

	for _, page := range pages {
		pagePath := filepath.Join(dir, page)
		tmpl, err := template.ParseFiles(layoutPath, pagePath)
		if err != nil {
			return nil, fmt.Errorf("parse template %q: %w", page, err)
		}
		cache[page] = tmpl
	}

	return &templateRenderer{cache: cache}, nil
}

// look up the pre-parsed template by name and executes it
func (tr *templateRenderer) Render(w http.ResponseWriter, name string, data interface{}) {
	tr.mu.RLock()
	tmpl, ok := tr.cache[name]
	tr.mu.RUnlock()

	if !ok {
		http.Error(w, fmt.Sprintf("template %q not found", name), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Printf("template execution error (%s): %v\n", name, err)
	}
}

// render the error.html template with a user-friendly message
func (tr *templateRenderer) RenderError(w http.ResponseWriter, message string, code int) {
	tr.mu.RLock()
	tmpl, ok := tr.cache["error.html"]
	tr.mu.RUnlock()

	w.WriteHeader(code)
	if !ok {
		http.Error(w, message, code)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Code":    code,
		"Message": message,
	})
}
