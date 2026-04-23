package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListJobs(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Title": "Jobs",
		
	}

	core.RenderPage(w,r, "apps/employees/views/jobs-list.html", data)
}