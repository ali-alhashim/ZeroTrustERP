package controllers

import (
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"fmt"
)

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Employees",
		
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