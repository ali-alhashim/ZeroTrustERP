package controllers

import (
	"encoding/json"
	"fmt"
    "database/sql"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
    
	query:=`SELECT e.id, e.badge_id, e.name, e.department_id, e.local_name, e.job_title_id, e.created_at, e.updated_at,e.image,
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

	query :="select id, badge_id, name, department_id, local_name, job_title_id, grade, created_at, updated_at, birth_date, active, goverment_id, image FROM employees WHERE 1=1"
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
		"Image":      "image",
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
    // 1. Use sql.NullString instead of plain string
    var deptID, jobID sql.NullString 

    // 2. Scan into these NullString variables
    err := rows.Scan(
        &e.ID, &e.BadgeID, &e.Name, &deptID, &e.LocalName, 
        &jobID, &e.Grade, &e.CreatedAt, &e.UpdatedAt, 
        &e.BirthDate, &e.Active, &e.GovermentID, &e.Image,
    )
    if err != nil {
        fmt.Printf("Scan Error: %v\n", err)
        continue
    }

    // 3. Check .Valid to see if the database actually had a value
    if deptID.Valid && deptID.String != "" && deptID.String != "0" {
        dept := GetDepartmentById(deptID.String)
        e.Department = &dept
    } else {
        e.Department = nil // Explicitly null if DB was null
    }

    if jobID.Valid && jobID.String != "" && jobID.String != "0" {
        job := GetJobTitleById(jobID.String)
        e.JobTitle = &job
    } else {
        e.JobTitle = nil // Explicitly null if DB was null
    }

    employees = append(employees, e)
}

	return employees

}


func CreateEmployee(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        data := map[string]interface{}{"Title": "Create Employee"}
        core.RenderPage(w, r, "apps/employees/views/employees-create.html", data)
        return // Good practice to return after rendering
    }

    if r.Method == http.MethodPost {
        // 1. Only call ParseMultipartForm for image uploads
        const maxMemory = 10 << 20 
        if err := r.ParseMultipartForm(maxMemory); err != nil {
            http.Error(w, "File too large or bad request", http.StatusBadRequest)
            return
        }

        // 2. Extract values
        name         := r.PostFormValue("name")
        localName    := r.PostFormValue("local_name")
        departmentId := r.PostFormValue("departmentId")
        jobTitleId   := r.PostFormValue("jobTitleId")
        birthDatestr := r.PostFormValue("birthDate")
        govermentId  := r.PostFormValue("govermentId")
        badgeId      := r.PostFormValue("badgeId")
        grade        := r.PostFormValue("grade")

        // 3. Handle Image
        var img image.Image
        file, _, err := r.FormFile("image")
        if err == nil {
            defer file.Close() // Close only if file exists
            img, _, err = image.Decode(file)
            if err != nil {
                http.Error(w, "Invalid image format", http.StatusBadRequest)
                return
            }
        } else if err != http.ErrMissingFile {
            // Some other error happened during upload
            http.Error(w, "Error uploading file", http.StatusInternalServerError)
            return
        }

        // 4. Pointer Logic for Optional Relations
        var theDepartment *models.Department
        var theJobTitle   *models.JobTitle

        if departmentId != "0" && departmentId != "" {
            dept := GetDepartmentById(departmentId)
            if dept.ID != 0 {
                theDepartment = &dept
            }
        }

        if jobTitleId != "0" && jobTitleId != "" {   
            job := GetJobTitleById(jobTitleId)
            if job.ID != 0 {
                theJobTitle = &job
            }
        }

        // 5. Date Parsing
        birthDate, err := time.Parse("2006-01-02", birthDatestr)
        if err != nil {
            http.Error(w, "Invalid Date (Required: YYYY-MM-DD)", http.StatusBadRequest)
            return
        }

        // 6. Build Model
        employee := models.Employee{
            Name:        name,
            LocalName:   localName,
            Department:  theDepartment,
            JobTitle:    theJobTitle,
            BirthDate:   birthDate,
            GovermentID: govermentId,
            BadgeID:     badgeId,
            Grade:       grade,
            Active:      true, // Usually true for new employees
        }

        // 7. Insert
        if err := InsertEmployeeToDB(employee, img); err != nil {
            fmt.Printf("Insert error: %v\n", err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/employees/list", http.StatusSeeOther)
    }
}

func GenerateBadgeIdApi(w http.ResponseWriter, r *http.Request) {
    // Get next_value from sequence (prefix_sequences)
    var badgeId string
    query := "SELECT next_value FROM prefix_sequences WHERE name = 'badge_id'"
    
    err := core.DB.QueryRow(query).Scan(&badgeId)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Generated Badge ID: %s\n", badgeId)

    // Set header to JSON
    w.Header().Set("Content-Type", "application/json")
    
    // Create the response map
    response := map[string]string{"badge_id": badgeId}

    // Encode and return the JSON response
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        fmt.Printf("Encoding error: %v\n", err)
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
        return
    }
}


func InsertEmployeeToDB(employee models.Employee, img image.Image) error {
    var savedImagePath string

    // 1. Image Handling
    if img != nil {
        dirPath := "./media/employees/images/"
        fileName := fmt.Sprintf("%s.jpg", employee.BadgeID)
        savedImagePath = filepath.Join(dirPath, fileName)

        if err := os.MkdirAll(dirPath, 0755); err != nil {
            return fmt.Errorf("failed to create directory: %v", err)
        }

        out, err := os.Create(savedImagePath)
        if err != nil {
            return fmt.Errorf("failed to create image file: %v", err)
        }
        defer out.Close()

        if err := jpeg.Encode(out, img, nil); err != nil {
            return fmt.Errorf("failed to encode image: %v", err)
        }
        employee.Image = savedImagePath
    }

    // 2. Safe NULL Handling (Prevents Panic)
    var deptID, jobID interface{}
    
    if employee.Department != nil && employee.Department.ID != 0 {
        deptID = employee.Department.ID
    } else {
        deptID = nil // Database will see this as NULL
    }

    if employee.JobTitle != nil && employee.JobTitle.ID != 0 {
        jobID = employee.JobTitle.ID
    } else {
        jobID = nil // Database will see this as NULL
    }

    // 3. Database Insert (Using $ placeholders for Postgres)
    query := `
        INSERT INTO employees (
            badge_id, name, department_id, local_name, 
            job_title_id, grade, birth_date, active, 
            goverment_id, image
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

    _, err := core.DB.Exec(query, 
        employee.BadgeID, 
        employee.Name, 
        deptID,             // Safe interface (could be nil)
        employee.LocalName,
        jobID,              // Safe interface (could be nil)
        employee.Grade, 
        employee.BirthDate, 
        employee.Active,
        employee.GovermentID, 
        employee.Image,
    )

    if err != nil {
        fmt.Printf("Database insert error: %v\n", err)
        return err
    }

    // 4. Success logic
    fmt.Printf("Employee %s inserted successfully\n", employee.Name)
    UpdateSequenceNextValue(employee.BadgeID, "badge_id")

    return nil
}


func UpdateSequenceNextValue(currentValue string, sequenceName string) {
    // 1. Start the transaction
    tx, err := core.DB.Begin()
    if err != nil {
        fmt.Printf("Error starting transaction: %v\n", err)
        return
    }

    // Defer a rollback. If the function returns before tx.Commit(), 
    // the changes are cancelled (safe mode).
    defer tx.Rollback()

    var prefix string
    var step int
    var digits int

    // 2. Select with FOR UPDATE to lock the row
    query := "SELECT prefix, step, digits FROM prefix_sequences WHERE name = $1 FOR UPDATE"
    err = tx.QueryRow(query, sequenceName).Scan(&prefix, &step, &digits)
    if err != nil {
        fmt.Printf("Error fetching sequence (locked): %v\n", err)
        return
    }

    // 3. Logic to calculate the next value (same as before)
    if len(currentValue) <= len(prefix) {
        fmt.Println("Error: currentValue format invalid")
        return
    }
    
    numberPart := currentValue[len(prefix):]
    currentNum, _ := strconv.Atoi(numberPart)
    nextNum := currentNum + step
    nextValue := fmt.Sprintf("%s%0*d", prefix, digits, nextNum)

    // 4. Update using the transaction handle (tx)
    updateQuery := "UPDATE prefix_sequences SET next_value = $1 WHERE name = $2"
    _, err = tx.Exec(updateQuery, nextValue, sequenceName)
    if err != nil {
        fmt.Printf("Error updating sequence: %v\n", err)
        return
    }

    // 5. Commit the transaction to save changes and release the lock
    if err := tx.Commit(); err != nil {
        fmt.Printf("Error committing transaction: %v\n", err)
        return
    }

    fmt.Printf("Sequence '%s' locked, updated, and released: %s\n", sequenceName, nextValue)
}


func EmployeeImageGET(w http.ResponseWriter, r *http.Request) {


    imageName := r.PathValue("imageName")

   
   fmt.Printf("Requested image: %s\n", imageName)

   safeName := filepath.Base(imageName)

   path := filepath.Join("media", "employees", "images", safeName)
  

    // Serve the image file
   http.ServeFile(w, r, path)
}