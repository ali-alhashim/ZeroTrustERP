package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListDepartments(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Departments",
		
	}

	core.RenderPage(w,r, "apps/employees/views/listDepartments.html", data)
}