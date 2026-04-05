package controllers

import (
	"fmt"
	"net/http"

	"zerotrusterp/core"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {

	fmt.Println("CreateRole called with method:", r.Method)

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Create Role",
			"Resources": core.GetAllResources(),
		}
		core.RenderPage(w,r, "apps/users/views/roles-create.html", data)
		return
	}

}