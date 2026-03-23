package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Users",
		
	}

	core.RenderPage(w, "apps/users/views/list.html", data)
}