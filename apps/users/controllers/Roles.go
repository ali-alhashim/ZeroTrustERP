package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)


func ListRoles(w http.ResponseWriter, r *http.Request) {
	fmt.Print("List Roles")



	query := r.URL.Query()

	search    := query.Get("q")
	sortBy    := query.Get("sort")
	order     := query.Get("order")
	page := query.Get("page")
	pageSize   := query.Get("pageSize")

	roles := GetRolesFromDB(search, sortBy, order, page, pageSize)

	data := map[string]interface{}{
		"Title": "Roles",
		"Roles": roles,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,

	}

	core.RenderPage(w,r, "apps/users/views/roles-list.html", data)

}

func GetRolesFromDB(search, sortBy, order, page, pageSize string) []models.Role {


	// Placeholder: Replace with actual DB query logic
	var roles []models.Role



	return roles
}