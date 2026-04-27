package controllers

import (
	    "net/http"
		"zerotrusterp/core"
	    "strconv"
		"zerotrusterp/apps/sequence/models"
)


func ListSequence(w http.ResponseWriter, r *http.Request){

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	 totalRecords := core.GetCountRecords("prefix_sequences")

	 records := GetSequenceFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Departments",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"Records": records,
		
	}

	core.RenderPage(w,r, "apps/sequence/views/sequence-list.html", data)

}


func GetSequenceFromDB(search, sort, order, page, pageSize string) []models.PrefixSequence{


	query := "SELECT id, name, prefix, next_value,  digits, step FROM prefix_sequences WHERE 1=1"
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
		"ID":        "id",
		"Prefix":     "prefix",
		"Name":    "name",
		
		
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


	



	var departments []models.PrefixSequence


	for rows.Next() {
		var l models.PrefixSequence
		err := rows.Scan(&l.ID, &l.Name, &l.Prefix, &l.NextValue, &l.Digits, &l.Step)
		if err != nil {
			panic(err)
		}
		departments = append(departments, l)
	}

	return departments

}


func CreateSequence(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodGet {

		data := map[string]interface{}{
		"Title": "Create Sequence",
		}


		core.RenderPage(w,r, "apps/sequence/views/sequence-create.html", data)

	}

	if r.Method == http.MethodPost {

		name := r.FormValue("name")
		prefix := r.FormValue("prefix")
		nextValue := r.FormValue("next_value")
		digits := r.FormValue("digits")
		step := r.FormValue("step")

		query := "INSERT INTO prefix_sequences (name, prefix, next_value, digits, step) VALUES ($1, $2, $3, $4, $5)"
		_, err := core.DB.Exec(query, name, prefix, nextValue, digits, step)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/sequence/list", http.StatusSeeOther)
	}

}