package controllers

import (
	"fmt"
	"net/http"
    "strconv"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
)

func ListJobs(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
    totalRecords := core.GetCountRecords("job_titles")

	jobs:= GetJobsFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Job Titles",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"Jobs": jobs,
		
	}

	core.RenderPage(w,r, "apps/employees/views/jobs-list.html", data)
}

func CreateJob(w http.ResponseWriter, r *http.Request){


	

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
		"Title": "Jobs",
		
	}
		core.RenderPage(w,r, "apps/employees/views/jobs-create.html", data)
	}

	if r.Method == http.MethodPost{
		//create job
		code       := r.FormValue("code")
		name       := r.FormValue("name")
		local_name := r.FormValue("local_name")
		description:=r.FormValue("description")

		query := "insert into job_titles (name, local_name, code, description) values ($1, $2, $3, $4)"

		_,err:= core.DB.Exec(query, name, local_name, code, description)

		if err!=nil{
			fmt.Print(err)
		}

		http.Redirect(w, r, "/employees/Jobs", http.StatusSeeOther)


	}



}


func GetJobsFromDB(search, sort, order, page, pageSize string) []models.JobTitle{

	query :="select id, name, local_name, code, description from job_titles where 1=1"

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (username ILIKE $" + strconv.Itoa(argIndex) +
			     " OR email ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":        "id",
		"Code":     "code",
		"Name":    "name",
		"LocalName":  "local_name",
		"Description":  "description",
		
		
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


	var jobs []models.JobTitle

	for rows.Next() {
		var l models.JobTitle
		err := rows.Scan(&l.ID, &l.Name, &l.LocalName, &l.Code, &l.Description)
		if err != nil {
			panic(err)
		}
		jobs = append(jobs, l)
	}

	return jobs


}


func GetJobTitleById(id string) models.JobTitle {

	var jobTitle models.JobTitle

	query :="select id, name, local_name, code, description from job_titles where id = $1"
    
	err := core.DB.QueryRow(query, id).Scan(
        &jobTitle.ID, 
        &jobTitle.Name, 
        &jobTitle.LocalName, 
        &jobTitle.Code, 
        &jobTitle.Description,
    )

	if err !=nil{
		fmt.Print(err)
	}



	return jobTitle
}



