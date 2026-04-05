package controllers

import (
	
	"net/http"
	
	"zerotrusterp/core"
	
)

func CreateRole(w http.ResponseWriter, r *http.Request) {


	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Create Role",
		}
		core.RenderPage(w,r, "apps/users/views/roles-create.html", data)
		return
	}

}