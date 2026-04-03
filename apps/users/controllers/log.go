package controllers

import (
	"net/http"
	"zerotrusterp/apps/users/models"
	"zerotrusterp/core"
	"strconv"

)



func InsertLog(user *models.User, resource string, action string) {

	// TODO: implement log insertion to database, create a new log record with user id, resource, action and timestamp
	query := "INSERT INTO logs (user_id, username, email, resource, action) VALUES ($1, $2, $3, $4, $5)"
	_, err := core.DB.Exec(query, user.ID, user.Username, user.Email, resource, action)
	if err != nil {
		panic(err)
	}
	
}



func ListLogs(w http.ResponseWriter, r *http.Request) {
	 

	query := r.URL.Query()

	search    := query.Get("q")
	sortBy    := query.Get("sort")
	order     := query.Get("order")
	page := query.Get("page")
	pageSize   := query.Get("pageSize")

	logs := GetLogsFromDB(search, sortBy, order, page, pageSize)

	data := map[string]interface{}{
		"Title": "Logs",
		"Logs": logs,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,

	}

	core.RenderPage(w,r, "apps/users/views/logs-list.html", data)
}


func GetLogsFromDB(search, sort, order, page, pageSize string) []models.Log {

	query := "SELECT id, user_id, username, email,  resource, action, timestamp FROM logs WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (username ILIKE $" + strconv.Itoa(argIndex) +
			" OR email ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}


	// 🔒 Safe sorting
	allowedSort := map[string]string{
		"ID":     "id",
		"Email":  "email",
		"UserID": "user_id",
		"Resource": "resource",
		"Action": "action",
		"Timestamp": "timestamp",
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


	



	var logs []models.Log


	for rows.Next() {
		var l models.Log
		err := rows.Scan(&l.ID, &l.UserID, &l.Username, &l.Email, &l.Resource, &l.Action, &l.Timestamp)
		if err != nil {
			panic(err)
		}
		logs = append(logs, l)
	}

	return logs
}