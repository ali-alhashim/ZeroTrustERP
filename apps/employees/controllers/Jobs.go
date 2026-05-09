package controllers

import (
	"fmt"
	"net/http"
    "strconv"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"encoding/json"
	"encoding/csv"
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
		query += " AND (code ILIKE $" + strconv.Itoa(argIndex) +
			     " OR name ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":          "id",
		"Code":        "code",
		"Name":        "name",
		"LocalName":   "local_name",
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
	if ps <= 0 || ps > 5000 {
		ps = 10
	}

	offset := (p - 1) * ps

	query += " LIMIT $" + strconv.Itoa(argIndex) +
		" OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, ps, offset)

		// ✅ Execute
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()


	var jobs []models.JobTitle

	for rows.Next() {
		var l models.JobTitle
		err := rows.Scan(&l.ID, &l.Name, &l.LocalName, &l.Code, &l.Description)
		if err != nil {
			fmt.Print(err)
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

func GetJobTitleListApi(w http.ResponseWriter, r *http.Request) {

	 totalRecords := core.GetCountRecords("job_titles")
	 fmt.Printf("\n GET Jobs Total Record for API= %d \n", totalRecords)
	 jobs := GetJobsFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}



func JobsDetails(w http.ResponseWriter, r *http.Request){

	id:=r.PathValue("id")

	

	theJob := GetJobTitleById(id)

	if r.Method == http.MethodGet {

	data := map[string]interface{}{
		"Title": "Jobs",
		"Job":theJob,
	}
		core.RenderPage(w,r,"apps/employees/views/jobs-details.html", data)
	}


}


func UpdateJob(w http.ResponseWriter, r *http.Request){

	id:=r.PathValue("id")

	if r.Method == http.MethodPost{


		//beforeUpdate := GetJobTitleById(id)

		code      := r.PostFormValue("code")
        name      := r.PostFormValue("name")
		local_name:= r.PostFormValue("local_name")
		description:= r.PostFormValue("description")

		query :=`UPDATE job_titles
		SET code = $1,
		name = $2,
        local_name = $3,
		description = $4
		WHERE id = $5 
		`
        _, err := core.DB.Exec(query, code, name, local_name,description, id)

		if err !=nil {
			fmt.Print(err)
		}

		core.InsertLog(core.GetCurrentUser(r), "JobTitle", "update Job Title id: "+id)

		http.Redirect(w, r, "/job/details/"+id, http.StatusSeeOther)

	}

}


///jobs/download/csv

func DownloadJobsCSV(w http.ResponseWriter, r *http.Request){

	totalRecords := core.GetCountRecords("job_titles")
	 
	jobs := GetJobsFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	// 2. Set Headers for CSV Download
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment;filename=Job_Definitions_export.csv")

	// THE QUICK FIX: Write the UTF-8 BOM
	w.Write([]byte{0xEF, 0xBB, 0xBF})

	// 3. Initialize CSV Writer
    writer := csv.NewWriter(w)
    defer writer.Flush()

	header := []string{"id", "name", "local_name", "code"}
	if err := writer.Write(header); err != nil {
        http.Error(w, "Failed to write header", http.StatusInternalServerError)
        return
    }

	

	

	for _, job := range jobs {
        row := []string{
            strconv.Itoa(job.ID),
            getString(job.Name),
            getString(job.LocalName),
            job.Code,
        }
        if err := writer.Write(row); err != nil {
            http.Error(w, "Failed to write row", http.StatusInternalServerError)
            return
        }
    }

}


func getString(s *string) string {
    if s == nil {
        return "" // or "N/A"
    }
    return *s
}



///jobs/upload/csv


func ImportCSVJobs(w http.ResponseWriter, r *http.Request) {
    // 1. Limit the size of the upload (e.g., 5MB) to prevent server abuse
    r.ParseMultipartForm(5 << 20)

    // 2. Retrieve the file from form data
    file, _, err := r.FormFile("csvFile")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 3. Initialize the CSV Reader
    reader := csv.NewReader(file)
    
    // Read the first line (the header) to skip it
    _, err = reader.Read()
    if err != nil {
        http.Error(w, "The CSV file is empty or invalid", http.StatusBadRequest)
        return
    }

    // 4. Iterate through the records
    records, err := reader.ReadAll()
    if err != nil {
        http.Error(w, "Error reading CSV content", http.StatusInternalServerError)
        return
    }

    for _, column := range records {
        // Based on your download structure: 0:id, 1:name, 2:local_name, 3:code
        // Note: If ID is provided, you might want to UPDATE. If empty, INSERT.
        
        id, _ := strconv.Atoi(column[0])
        name := column[1]
        localName := column[2]
        code := column[3]

        // 5. Database Logic
        if id > 0 {
            // Logic to Update existing department
            // core.UpdateDepartment(id, name, localName, code)

		query := "UPDATE job_titles SET code=$1, name=$2, local_name=$3  WHERE id=$4"

		_, err := core.DB.Exec(query, code, name, localName, id)
		if err != nil {
			fmt.Printf("Error updating jobs: %v\n", err)
			http.Error(w, "Error updating jobs", http.StatusInternalServerError)
			return
		}

		CurrentUser := core.GetCurrentUser(r)

		core.InsertLog(CurrentUser, "Jobs", fmt.Sprintf("Updated jobs ID %s with code: %s and name: %s",id, code, name))




        } else {
            // Logic to Insert new department
            // core.InsertDepartment(name, localName, code)
			query:= "insert into job_titles (code, name, local_name) values ($1, $2, $3)"
        
		 _,err := core.DB.Exec(query, code,name,localName)

		 if err != nil {
			fmt.Println("Error inserting jobs:", err)
			//http.Error(w, "Error creating jobs", http.StatusInternalServerError)
			//return
		}

		CurrentUser := core.GetCurrentUser(r)

		core.InsertLog(CurrentUser, "Jobs", fmt.Sprintf("Created Job code %s with name: %s",code, name))
        }
    }

    
    http.Redirect(w, r, "/employees/Jobs", http.StatusSeeOther)
}




