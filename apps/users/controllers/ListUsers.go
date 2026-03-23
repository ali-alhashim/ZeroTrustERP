package controllers

import (
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search := query.Get("q")
	sortBy := query.Get("sort")
	order := query.Get("order")

	users := GetUsersFromDB(search, sortBy, order)

	data := map[string]interface{}{
		"Title": "Users",
		"Users": users,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
	}

	core.RenderPage(w, "apps/users/views/list.html", data)
}

func GetUsersFromDB(search, sort, order string) []models.User {
	query := "SELECT id, username, email, role FROM users WHERE 1=1"

	if search != "" {
		query += " AND (username LIKE '%" + search + "%' OR email LIKE '%" + search + "%')"
	}

	if sort != "" {
		query += " ORDER BY " + sort
		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	// Execute query...
	var users []models.User
	// TODO: Execute the query and populate users slice
	return users
}
