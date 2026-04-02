package core

import (
	"html/template"
	"net/http"
	"path/filepath"
	"os"
)

// Render normal page (with layout)
func RenderPage(w http.ResponseWriter,r *http.Request, tmpl string, data interface{}) {



	dataMap, ok := data.(map[string]interface{})
		if ok {
			dataMap["Menus"] = Menus
		}

	wd, _ := os.Getwd()


	//Automatically inject UserEmail from the Request Context
	if email, ok := r.Context().Value(UserEmailKey).(string); ok {
        dataMap["UserEmail"] = email
    }


    

	

	t, err := template.ParseFiles(
	filepath.Join(wd, "core/templates/base.html"),
	filepath.Join(wd, "core/templates/sidebar.html"),
	filepath.Join(wd, "core/templates/header.html"),
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
	wd, _ := os.Getwd()

	t, err := template.ParseFiles(
		filepath.Join(wd, "core/templates/report.html"),
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

// Render page without layout (for login, etc.)
func RenderPageNoLayout(w http.ResponseWriter, tmpl string, data interface{}) {
	wd, _ := os.Getwd()

	t, err := template.ParseFiles(
		filepath.Join(wd, tmpl),
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
