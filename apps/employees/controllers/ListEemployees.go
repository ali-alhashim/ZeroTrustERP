package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Employees",
		
	}

	core.RenderPage(w, "apps/employees/views/list.html", data)
}