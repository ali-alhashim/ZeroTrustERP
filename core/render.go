package core

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Render normal page (with layout)
func RenderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	basePath := "core/templates"

	t, err := template.ParseFiles(
		filepath.Join(basePath, "base.html"),
		filepath.Join(basePath, "sidebar.html"),
		filepath.Join(basePath, "header.html"),
		tmpl, // already full path from app
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Render report (no layout, A4 style)
func RenderReport(w http.ResponseWriter, tmpl string, data interface{}) {
	basePath := "core/templates"

	t, err := template.ParseFiles(
		filepath.Join(basePath, "report.html"),
		tmpl,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
