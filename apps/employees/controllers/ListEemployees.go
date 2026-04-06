package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Employees",
		
	}

	core.RenderPage(w,r, "apps/employees/views/employees-list.html", data)
}