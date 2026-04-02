package controllers

import (
	"net/http"
	"zerotrusterp/core"
	
)


func DashboardController(w http.ResponseWriter, r *http.Request) {


	data := map[string]interface{}{
		"Title": "Dashboard",
		
	}

	core.RenderPage(w,r, "apps/dashboard/views/dashboard.html", data)

}