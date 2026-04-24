package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListDepartments(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Departments",
		
	}

	core.RenderPage(w,r, "apps/employees/views/departments-list.html", data)
}


func CreateDepartment(w http.ResponseWriter, r *http.Request){

	data := map[string]interface{}{
		"Title": "Departments",
		
	}

	if r.Method == http.MethodGet {
		core.RenderPage(w,r, "apps/employees/views/departments-create.html", data)
	}
}