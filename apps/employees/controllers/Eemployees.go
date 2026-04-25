package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
)

func ListEmployees(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
    totalRecords := core.GetCountRecords("employees")

	employees:= GetEmployeesFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Employees",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"Employees": employees,
		
	}

	core.RenderPage(w,r, "apps/employees/views/employees-list.html", data)
}



func GetEmployeeById(id string) models.Employee{
	var employee models.Employee
    
	query:=`SELECT e.id, e.badge_id, e.name, e.department_id, e.local_name, e.job_title_id, e.created_at, e.updated_at,
	       d.id, d.name, d.local_name, d.code, d.manager_id, d.created_at, d.updated_at, d.active,
		   j.id, j.name, j.local_name, j.code, j.description, j.crearted_at, j.updated_at  
	       FROM employees e
	       LEFT JOIN departments d ON e.department_id = d.id
		   LEFT JOIN job_titles j  ON e.job_title_id = j.id
		   WHERE e.id = $1
	       `
	rows, err := core.DB.Query(query, id)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return models.Employee{}
    }
    defer rows.Close()


	return employee
}

func GetEmployeesFromDB(search, sort, order, page, pageSize string)[]models.Employee{

	var employees []models.Employee

	query :="select id, badge_id, name, department_id, local_name, job_title_id, grade, created_at, updated_at, birth_date, active, goverment_id FROM employees WHERE 1=1"
    args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (name ILIKE $" + strconv.Itoa(argIndex) +
			     " OR local_name ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":         "id",
		"Name":       "name",
		"LocalName":  "local_name",
		"Active":     "active",
		
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


	for rows.Next() {
        var e models.Employee
        // Temporary variables to hold IDs from the scan
        var deptID, jobID string 

        // Scan into pointers
        err := rows.Scan(
            &e.ID, &e.BadgeID, &e.Name, &deptID, &e.LocalName, 
            &jobID, &e.Grade, &e.CreatedAt, &e.UpdatedAt, 
            &e.BirthDate, &e.Active, &e.GovermentID,
        )
        if err != nil {
            fmt.Printf("Scan Error: %v", err)
            continue
        }

        // Now fetch the related models using the IDs we just scanned
		dept := GetDepartmentById(deptID)
		job  := GetJobTitleById(jobID)
        e.Department = &dept
        e.JobTitle   = &job

        employees = append(employees, e)
    }

	return employees

}


func CreateEmployee(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodGet{
		data := map[string]interface{}{
				"Title": "Employees",
			}

		core.RenderPage(w,r, "apps/employees/views/employees-create.html", data)
	}
	
}