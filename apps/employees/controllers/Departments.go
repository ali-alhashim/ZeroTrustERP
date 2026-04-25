package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"strconv"
)

func ListDepartments(w http.ResponseWriter, r *http.Request) {


	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
    totalRecords := core.GetCountRecords("departments")

	departments:= GetDepartmentsFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Departments",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"Departments": departments,
		
	}

	core.RenderPage(w,r, "apps/employees/views/departments-list.html", data)
}


func CreateDepartment(w http.ResponseWriter, r *http.Request){

	data := map[string]interface{}{
		"Title": "Departments",
		
	}

	if r.Method == http.MethodGet {
		core.RenderPage(w,r, "apps/employees/views/departments-create.html", data)
	}

	if r.Method == http.MethodPost {

		code:=r.FormValue("code")
		name:=r.FormValue("name")
		nameAR:=r.FormValue("nameAR")
		
		active:=r.FormValue("active") == "on"


		var managerID interface{} 

		val := r.FormValue("manager")

		if val == "" || val == "0" {
			managerID = nil // This is allowed for interfaces
		} else {
			managerID = val
		}


		
		fmt.Printf("Create Department code:%s with name:%s -Ar Name:%s and the managerId:%s -status:%t", code, name, nameAR, managerID, active)

		query:= "insert into departments (code, name, local_name, active, manager_id) values ($1, $2, $3, $4, $5)"
        
		 _,err := core.DB.Exec(query, code,name,nameAR,active,managerID)

		 if err != nil {
			fmt.Println("Error inserting department:", err)
			http.Error(w, "Error creating department", http.StatusInternalServerError)
			return
		}

		CurrentUser := core.GetCurrentUser(r)


		core.InsertLog(CurrentUser, "Departments", fmt.Sprintf("Created Department code %s with name: %s",code, name))

		http.Redirect(w, r, "/employees/departments", http.StatusSeeOther)

	}
}



func GetDepartmentsFromDB(search, sort, order, page, pageSize string) []models.Department {

	query := "SELECT id, name, local_name, code,  manager_id, active FROM departments WHERE 1=1"
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
		"Code":     "code",
		"Name":    "name",
		"LocalName":  "local_name",
		"Manager":  "manager_id",
		"Active":    "active",
		
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


	



	var departments []models.Department


	for rows.Next() {
		var l models.Department
		err := rows.Scan(&l.ID, &l.Name, &l.LocalName, &l.Code, &l.Manager, &l.Active)
		if err != nil {
			panic(err)
		}
		departments = append(departments, l)
	}

	return departments
}



func DepartmentsDetails(w http.ResponseWriter, r *http.Request){

	departmentID := r.PathValue("id")

	fmt.Print(" \n Get Department Details ID = "+ departmentID +"\n")
   
	
	
  

	data := map[string]interface{}{
		"Title": "Departments",
		"Department":GetDepartmentById(departmentID),
		
	}

	core.RenderPage(w,r, "apps/employees/views/departments-details.html", data)
}



func GetDepartmentById(id string) models.Department {
    var dept models.Department
    
    // Ensure the SELECT statement is formatted correctly
    query := `
        SELECT 
            d.id, d.code, d.name, d.local_name, d.manager_id, d.active,
            e.id, e.badge_id, e.name, e.department_id, e.local_name, e.job_title_id
        FROM departments d 
        LEFT JOIN employees e ON d.id = e.department_id
        WHERE d.id = $1`

    rows, err := core.DB.Query(query, id)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return models.Department{}
    }
    defer rows.Close()

    for rows.Next() {
        var emp models.Employee
        // Using pointers for employee fields to handle LEFT JOIN NULLs
        var empId *int
        var empBadge, empName, empLocal, empJob *string
        var empDeptId *int

        err := rows.Scan(
            &dept.ID, &dept.Code, &dept.Name, &dept.LocalName, &dept.Manager, &dept.Active,
            &empId, &empBadge, &empName, &empDeptId, &empLocal, &empJob,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            return models.Department{}
        }

        // If empId is not null, an employee exists for this row
        if empId != nil {
            emp.ID = *empId
            if empName != nil { emp.Name = *empName }
            if empBadge != nil { emp.BadgeID = *empBadge }
            // ... map other fields ...
            
            dept.Employees = append(dept.Employees, emp)
        }
    }

    return dept
}