package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"zerotrusterp/apps/users/models"
	"zerotrusterp/core"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search    := query.Get("q")
	sortBy    := query.Get("sort")
	order     := query.Get("order")
	page := query.Get("page")
	pageSize   := query.Get("pageSize")

	users := GetUsersFromDB(search, sortBy, order, page, pageSize)
	totalRecords := core.GetCountRecords("users")


	data := map[string]interface{}{
		"Title": "Users",
		"Users": users,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,

	}

	core.RenderPage(w,r, "apps/users/views/users-list.html", data)
}

func GetUsersFromDB(search, sort, order, page, pageSize string) []models.User {

	query := "SELECT id, username, email, active, online, last_login FROM users WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// 🔍 SAFE search
	if search != "" {
		query += " AND (username ILIKE $" + strconv.Itoa(argIndex) +
			" OR email ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	// 🔒 Safe sorting
	allowedSort := map[string]string{
		"id":     "id",
		"email":  "email",
		"active": "active",
		"online": "online",
		"last_login":"last_login",
	}

	if col, ok := allowedSort[sort]; ok {
		query += " ORDER BY " + col
		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	// 📄 Pagination (page + pageSize)
	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	// defaults
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 10
	}

	offset := (p - 1) * ps

	query += " LIMIT $" + strconv.Itoa(argIndex) +
		" OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, ps, offset)

	// ✅ Execute
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.Online, &u.LastLogin)
		if err != nil {
			panic(err)
		}
		users = append(users, u)
	}

	return users
}



func SetUserActive(w http.ResponseWriter, r *http.Request){

	userID := r.PathValue("id")
	fmt.Print("set user active with ID:", userID)
    
	query :="update users set active=true where id = $1"

	core.DB.Exec(query, userID)


	http.Redirect(w, r, "/users/list", http.StatusSeeOther)
	

}

func SetUserInactive(w http.ResponseWriter, r *http.Request){

	userID := r.PathValue("id")
	fmt.Print("set user inactive ID:", userID)
    
	query :="update users set active=false where id = $1"

	core.DB.Exec(query, userID)

	http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		
}



