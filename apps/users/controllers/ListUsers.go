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
	query := "SELECT id, username,  email, active, online FROM users WHERE 1=1"

	// 🔍 Search (Postgres uses ILIKE for case-insensitive)
	if search != "" {
		query += " AND (username ILIKE '%" + search + "%' OR email ILIKE '%" + search + "%')"
	}

	// 🔒 Safe sorting (IMPORTANT)
	allowedSort := map[string]bool{
		"ID":       true,
		"Email":    true,
		"Active":   true,
		"Online":   true,
	}

	if allowedSort[sort] {
		query += " ORDER BY " + sort
		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	// ✅ Execute query
	rows, err := core.DB.Query(query)
	if err != nil {
		panic(err) // later handle better
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID,  &u.Username, &u.Email, &u.Active,  &u.Online)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	return users
}

