package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"strconv"
	"encoding/json"
	 "database/sql"
	 "encoding/csv"
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
		query += " AND (code ILIKE $" + strconv.Itoa(argIndex) +
			     " OR name ILIKE $" + strconv.Itoa(argIndex+1) + ")"

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
	
    var ManagerID *string

	for rows.Next() {
		var l models.Department
		var EmployeesDepartment [] models.Employee
		err := rows.Scan(&l.ID, &l.Name, &l.LocalName, &l.Code, &ManagerID, &l.Active)
		if err != nil {
			panic(err)
		}

		if ManagerID != nil {
			fmt.Print(" \n Get Department Manager ID = "+ *ManagerID +"\n")
			Manager := GetEmployeeById(*ManagerID)
			fmt.Print(" \n Get Department Manager  = "+ Manager.Name +"\n")
			l.Manager = &Manager
		} else {
			l.Manager = nil
		}


		//get all employees members of this depaerment l.ID
		query = "select id, badge_id, name, local_name from employees where department_id = $1"
		empRows, empErr := core.DB.Query(query,l.ID)
		if empErr !=nil {
			fmt.Print(empErr)
		}
		for empRows.Next() {
            
			var emp models.Employee
		    errE := rows.Scan(&emp.ID, &emp.BadgeID, emp.Name, emp.LocalName)

			if errE !=nil{
				fmt.Print(errE)
			}

			EmployeesDepartment = append(EmployeesDepartment, emp)

		}


		
		l.Employees = EmployeesDepartment

        
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
            e.id, e.department_id
        FROM departments d 
        LEFT JOIN employees e ON d.id = e.department_id
        WHERE d.id = $1`

    rows, err := core.DB.Query(query, id)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return models.Department{}
    }
    defer rows.Close()

    // FIX: Use sql.NullInt64 to handle cases where manager_id is NULL in the DB
    var managerID sql.NullInt64

    for rows.Next() {
        var emp models.Employee
        // Using pointers for employee fields to handle LEFT JOIN NULLs
        var empId *int
       
        var empDeptId *int

        // FIX: Scan managerID (variable) instead of &dept.Manager (struct field)
        err := rows.Scan(
            			&dept.ID, &dept.Code, &dept.Name, &dept.LocalName, &managerID, &dept.Active,
           			    &empId, &empDeptId,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            return models.Department{}
        }

        // FIX: Only assign the manager to the struct if it is valid (not NULL)
        if managerID.Valid {
            // Assuming dept.Manager is of type int. 
            // If it is *int, use: val := int(managerID.Int64); dept.Manager = &val
			 managerStr := strconv.FormatInt(managerID.Int64, 10)
             manager := GetEmployeeById(managerStr)
			dept.Manager = &manager
        }

        // If empId is not null, an employee exists for this row
        if empId != nil {
            
			  idString := strconv.Itoa(*empId)
			  emp = GetEmployeeById(idString)
			
			
            
            dept.Employees = append(dept.Employees, emp)
        }
    }

    return dept
}


func isDepartmentManagerChanged(departmentID string, managerID string) bool {
    var currentManagerID string

	fmt.Printf("Check is Department Manager is Updated for %s with Manager ID %s", departmentID, managerID)

    // 1. Fetch the CURRENT manager assigned to this department
    query := "SELECT manager_id FROM departments WHERE id = $1"
    
    // Assuming 'db' is your database connection (sql.DB)
    err := core.DB.QueryRow(query, departmentID).Scan(&currentManagerID)
    if err != nil {
        // If the department doesn't exist, technically it's a "change" or an error
        return true 
    }

    // 2. Compare. If the IDs are different, the manager has changed.
    return currentManagerID != managerID
}


func UpdateDepartment(w http.ResponseWriter, r *http.Request){

	departmentID := r.PathValue("id")

	code := r.FormValue("code")
	name := r.FormValue("name")
	local_name:= r.FormValue("local_name")
	manager:= r.FormValue("manager")
	active := r.FormValue("active") == "on"

	fmt.Printf("Update Department %s Code:%s, Name:%s, LocalName:%s, Manager:%s, active:%t",departmentID, code, name, local_name, manager, active)

	if isDepartmentManagerChanged(departmentID, manager){
	 
		exManagerDepartment(departmentID, manager)

	}

	query := "UPDATE departments SET code=$1, name=$2, local_name=$3, manager_id=$4, active=$5 WHERE id=$6"

	_, err := core.DB.Exec(query, code, name, local_name, manager, active, departmentID)
	if err != nil {
		fmt.Printf("Error updating department: %v\n", err)
		http.Error(w, "Error updating department", http.StatusInternalServerError)
		return
	}

	CurrentUser := core.GetCurrentUser(r)

	core.InsertLog(CurrentUser, "Departments", fmt.Sprintf("Updated Department ID %s with code: %s and name: %s",departmentID, code, name))


	http.Redirect(w, r, "/employees/departments", http.StatusSeeOther)


}


func exManagerDepartment(departmentID string, manager string){

	// get the exist employee and set 

}


func GetDepartmentListApi(w http.ResponseWriter, r *http.Request) {

	totalRecords := core.GetCountRecords("departments")
	 
	departments := GetDepartmentsFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}




///department/download/csv
func DownloadDepartmentsCSV(w http.ResponseWriter, r *http.Request){

	totalRecords := core.GetCountRecords("departments")
	 
	departments := GetDepartmentsFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	// 2. Set Headers for CSV Download
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment;filename=departments_export.csv")

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

	for _, dept := range departments {
        row := []string{
            strconv.Itoa(dept.ID),
            dept.Name,
            dept.LocalName,
            dept.Code,
        }
        if err := writer.Write(row); err != nil {
            http.Error(w, "Failed to write row", http.StatusInternalServerError)
            return
        }
    }

}


///department/upload/csv  enctype="multipart/form-data" type="file" name="csvFile"

func ImportCSVDepartments(w http.ResponseWriter, r *http.Request) {
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

		query := "UPDATE departments SET code=$1, name=$2, local_name=$3  WHERE id=$4"

		_, err := core.DB.Exec(query, code, name, localName, id)
		if err != nil {
			fmt.Printf("Error updating department: %v\n", err)
			http.Error(w, "Error updating department", http.StatusInternalServerError)
			return
		}

		CurrentUser := core.GetCurrentUser(r)

		core.InsertLog(CurrentUser, "Departments", fmt.Sprintf("Updated Department ID %s with code: %s and name: %s",id, code, name))




        } else {
            // Logic to Insert new department
            // core.InsertDepartment(name, localName, code)
			query:= "insert into departments (code, name, local_name) values ($1, $2, $3)"
        
		 _,err := core.DB.Exec(query, code,name,localName)

		 if err != nil {
			fmt.Println("Error inserting department:", err)
			http.Error(w, "Error creating department", http.StatusInternalServerError)
			return
		}

		CurrentUser := core.GetCurrentUser(r)

		core.InsertLog(CurrentUser, "Departments", fmt.Sprintf("Created Department code %s with name: %s",code, name))
        }
    }

    
    http.Redirect(w, r, "/employees/departments", http.StatusSeeOther)
}


