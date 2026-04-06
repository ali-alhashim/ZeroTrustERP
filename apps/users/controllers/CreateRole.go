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



	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		roleName := r.FormValue("roleName")
		description := r.FormValue("description")
		permissionsCount := r.FormValue("permissionsCount")

		fmt.Printf("Received form data: roleName=%s, description=%s, permissionsCount=%s\n", roleName, description, permissionsCount)

		//ok we received permissionsCount as string, we need to convert it to int
		// also we need to loop through the permissions and create a slice of permissions
		// Resource[], Action[], description[]
		// we will insert the role first and get the role ID, then we will insert the permissions with the role ID



		http.Redirect(w, r, "/users/roles", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

}