package core

import (
	"html/template"
	"net/http"
)

// Render normal page (with layout)
func RenderPage(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(
		"core/templates/base.html",
		"core/templates/sidebar.html",
		"core/templates/header.html",
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

// Render report (no layout, A4 style)
func RenderReport(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(
		"core/templates/report.html",
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
