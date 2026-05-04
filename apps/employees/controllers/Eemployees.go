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
    "io"
    "mime/multipart"
    
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



func GetEmployeeById(id string) models.Employee {
    fmt.Printf("............Fetching employee with ID: %s\n", id)

    var employee models.Employee

    
    // Null types for the LEFT JOINed tables
    var dID, dManager sql.NullInt64
    var dName, dLocal, dCode sql.NullString
    var dActive sql.NullBool
    var dCreated, dUpdated sql.NullTime
    
    var jID sql.NullInt64
    var jName, jLocal, jCode, jDesc sql.NullString
    var jCreated, jUpdated sql.NullTime

    query := `
        SELECT 
            e.id, e.badge_id, e.name, e.department_id, e.local_name, e.job_title_id, e.created_at, e.updated_at, e.image,e.education, e.major, e.religion,e.goverment_id, e.email, e.nationality, e.gender, e.marital_status, e.phone_number, e.address,e.birth_date,
            d.id, d.name, d.local_name, d.code, d.manager_id, d.created_at, d.updated_at, d.active,
            j.id, j.name, j.local_name, j.code, j.description, j.created_at, j.updated_at  
        FROM employees e
        LEFT JOIN departments d ON e.department_id = d.id
        LEFT JOIN job_titles j  ON e.job_title_id = j.id
        WHERE e.id = $1`

    err := core.DB.QueryRow(query, id).Scan(
        // IMPORTANT: Scan into the ID fields, NOT the struct fields
        &employee.ID, &employee.BadgeID, &employee.Name, &dID, &employee.LocalName, 
        &jID, &employee.CreatedAt, &employee.UpdatedAt, &employee.Image, &employee.Education, &employee.Major, &employee.Religion,&employee.GovermentID, &employee.Email, &employee.Nationality, &employee.Gender, &employee.MaritalStatus, &employee.PhoneNumber, &employee.Address, &employee.BirthDate,
        // Department scan
        &dID, &dName, &dLocal, &dCode, &dManager, &dCreated, &dUpdated, &dActive,
        // Job Title scan
        &jID, &jName, &jLocal, &jCode, &jDesc, &jCreated, &jUpdated,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Printf("❌ No employee found in database with ID: %s\n", id)
        } else {
            fmt.Printf("❌ Database error during Scan: %v\n", err)
        }
        return models.Employee{}
    }

    // Map Department if JOIN found a record
    if dID.Valid {
        dept := models.Department{
            ID:        int(dID.Int64),
            Name:      dName.String,
            LocalName: dLocal.String,
            Code:      dCode.String,
            Active:    dActive.Bool,
            CreatedAt: dCreated.Time,
            UpdatedAt: dUpdated.Time,
        }
        if dManager.Valid {
            // Check if your Employee.ID is string or int. 
            // If ID is string: strconv.Itoa(int(dManager.Int64))
            // If ID is int: int(dManager.Int64)
            dept.Manager = &models.Employee{
                ID: int(dManager.Int64), 
            }
        }
        employee.Department = &dept
    }

    // Map JobTitle if JOIN found a record
    if jID.Valid {
        employee.JobTitle = &models.JobTitle{
            ID:   int(jID.Int64),
            Name: jName.String,
            // ... add other job fields here
        }
    }

    var familyMembers []models.FamilyMember
    var certifications []models.Certification
    var emergencyContacts []models.EmergencyContact
    var employeeDocuments []models.EmployeeDocument
    
    familyMembers = GetEmployeeFamilyMembers(id)
    certifications = GetEmployeeCertifications(id)
    emergencyContacts = GetEmployeeEmergencyContacts(id)
    employeeDocuments = GetEmployeeDocuments(id)

    employee.FamilyMembers = familyMembers
    employee.Certifications = certifications
    employee.EmergencyContacts = emergencyContacts
    employee.EmployeeDocuments = employeeDocuments
    return employee
}

func GetEmployeeDocuments(employeeId string) []models.EmployeeDocument {
    var employeeDocuments []models.EmployeeDocument
    query := `SELECT id, employee_id, name, type, expiry_date, file_path FROM employee_documents WHERE employee_id = $1`

    rows, err := core.DB.Query(query, employeeId)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return employeeDocuments
    }
    defer rows.Close()

    var EmployeeID int

    for rows.Next() {
        var doc models.EmployeeDocument
        err := rows.Scan(
            &doc.ID,
            &EmployeeID,
            &doc.Name,
            &doc.Type,
            &doc.ExpiryDate,
            &doc.FilePath,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            continue
        }
        employeeDocuments = append(employeeDocuments, doc)
    }

    return employeeDocuments
}

func GetEmployeeEmergencyContacts(employeeId string) []models.EmergencyContact {
    var emergencyContacts []models.EmergencyContact
    query := `SELECT id, employee_id, name, relationship, phone FROM emergency_contacts WHERE employee_id = $1` 

    rows, err := core.DB.Query(query, employeeId)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return emergencyContacts
    }
    defer rows.Close()

    var EmployeeID int

    for rows.Next() {
        var ec models.EmergencyContact
        err := rows.Scan(
            &ec.ID,
            &EmployeeID,
            &ec.Name,
            &ec.Relationship,
            &ec.Phone,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            continue
        }
        emergencyContacts = append(emergencyContacts, ec)
    }

    return emergencyContacts
}


func GetEmployeeCertifications(employeeId string) []models.Certification {
    var certifications []models.Certification
    query := `SELECT id, employee_id, name, issuer, issue_date, expiry_date, file_path FROM certifications WHERE employee_id = $1`

    rows, err := core.DB.Query(query, employeeId)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return certifications
    }
    defer rows.Close()

    var EmployeeID int

    for rows.Next() {
        var cert models.Certification
        err := rows.Scan(
            &cert.ID,
            &EmployeeID,
            &cert.Name,
            &cert.Issuer,
            &cert.IssueDate,
            &cert.ExpiryDate,
            &cert.FilePath,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            continue
        }
        certifications = append(certifications, cert)
    }

    return certifications
}

func GetEmployeeFamilyMembers(employeeId string) []models.FamilyMember {

    var familyMembers []models.FamilyMember
    query := `SELECT id, employee_id, name, relationship, contact_number, government_id, birth_date, file_path FROM family_members WHERE employee_id = $1`

    rows, err := core.DB.Query(query, employeeId)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return familyMembers
    }
    defer rows.Close()

    var EmployeeID int

    for rows.Next() {
        var fm models.FamilyMember
        err := rows.Scan(
            &fm.ID,
            &EmployeeID,
            &fm.Name,
            &fm.Relationship,
            &fm.ContactNumber,
            &fm.GovernmentId,
            &fm.BirthDate,
            &fm.FilePath,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            continue
        }
        familyMembers = append(familyMembers, fm)
    }

    return familyMembers
}







func GetEmployeesFromDB(search, sort, order, page, pageSize string)[]models.Employee{

	var employees []models.Employee

	query :="select id, badge_id, name, department_id, local_name, job_title_id, grade, created_at, updated_at, birth_date, active, goverment_id, image, education, major, religion, email, nationality, gender, marital_status, phone_number, address FROM employees WHERE 1=1"
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
    &e.ID,            // 1
    &e.BadgeID,       // 2
    &e.Name,          // 3
    &deptID,          // 4
    &e.LocalName,     // 5
    &jobID,           // 6
    &e.Grade,         // 7
    &e.CreatedAt,     // 8
    &e.UpdatedAt,     // 9
    &e.BirthDate,     // 10
    &e.Active,        // 11
    &e.GovermentID,   // 12
    &e.Image,         // 13
    &e.Education,     // 14
    &e.Major,         // 15
    &e.Religion,      // 16
    &e.Email,         // 17 (Fixed duplicate scan here)
    &e.Nationality,   // 18
    &e.Gender,        // 19
    &e.MaritalStatus, // 20
    &e.PhoneNumber,   // 21
    &e.Address,       // 22
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
        email        := r.PostFormValue("email")
        nationality  := r.PostFormValue("nationality")
        gender      := r.PostFormValue("gender")
        active      := r.PostFormValue("active") == "on" // Checkbox handling
        maritalStatus := r.PostFormValue("maritalStatus")
        phoneNumber  := r.PostFormValue("phoneNumber")
        address      := r.PostFormValue("address")
        education     := r.PostFormValue("education")
        religion      := r.PostFormValue("religion")
        major         := r.PostFormValue("major")


        var Certifications    []models.Certification
        var FamilyMembers     []models.FamilyMember
        var EmergencyContacts []models.EmergencyContact
        var EmployeeDocuments []models.EmployeeDocument

        var filePath string
        var CertName string
        var issuingOrgan string

        //Certifcations maybe there are multiple certification maybe empty
        certifications       := r.PostForm["certificationName[]"]
        issuingOrganizations := r.PostForm["issuingOrganization[]"]
        issuingDate:= r.PostForm["issueDate[]"]
        expirationDate := r.PostForm["expirationDate[]"]

        certificationAttachments := r.MultipartForm.File["certificationAttachment[]"]  

        if len(certifications) > 0 {
            fmt.Printf("Received certifications: %v\n", certifications)
            //loop through certifications and build the slice of Certification models
            for i := range certifications {

                    CertName = certifications[i]
                    

                issuingDate_parsedDate, err := time.Parse("2006-01-02", issuingDate[i])
                if err != nil {
                     fmt.Printf("invalid date format: %v", err)
                }
                expirationDate_parsedDate, err := time.Parse("2006-01-02", expirationDate[i])
                if err != nil {
                     fmt.Printf("invalid date format: %v", err)
                }

                if i < len(issuingOrganizations) {
                    issuingOrgan = issuingOrganizations[i]
                } else {
                    fmt.Printf("No issuing organization for certification: %s\n", certifications[i])
                    issuingOrgan = ""
                }

                if i < len(certificationAttachments) {
                    fileHeader := certificationAttachments[i]
                    // Process the file as needed (e.g., save to disk, get path, etc.)
                   filePath = UploadEmployeeAttachment(fileHeader, badgeId, "certifications") // Reusing the same function for simplicity
                    fmt.Printf("Received certification attachment: %s\n", fileHeader.Filename)
                    // You would typically save the file and get its path to store in the database
                } else {
                    fmt.Printf("No attachment for certification: %s\n", certifications[i])
                    filePath =""
                }

                fmt.Printf("Certification details - Name: %s, Issuer: %s, Issue Date: %s, Expiry Date: %s, File Path: %s\n", CertName, issuingOrgan, issuingDate_parsedDate.Format("2006-01-02"), expirationDate_parsedDate.Format("2006-01-02"), filePath)
                cert := models.Certification{
                    Name  : CertName,
                    Issuer: issuingOrgan,
                    IssueDate: issuingDate_parsedDate,
                    ExpiryDate: expirationDate_parsedDate,
                    FilePath: filePath,
                    // You would also extract other fields like Issuer, IssueDate, ExpiryDate, and FilePath similarly
                }
                Certifications = append(Certifications, cert)
            } //end of loop for certifications
        } else {
            fmt.Println("No certifications received")
        }


        familyMemberName          := r.PostForm["familyMemberName[]"]
        familyMemberRelationship  := r.PostForm["familyMemberRelationship[]"]
        familyMemberContactNumber := r.PostForm["familyMemberContactNumber[]"]
        familyMemberGender        := r.PostForm["familyMemberGender[]"]
        familyMemberGovernmentId  := r.PostForm["familyMemberGovernmentId[]"]
        familyMemberBirthDate     := r.PostForm["familyMemberBirthDate[]"]
        familyMemberAttachment    := r.MultipartForm.File["familyMemberAttachment[]"]
        
        if len(familyMemberName) > 0 {
            fmt.Printf("Received family members: %v , relationships: %v , contact numbers: %v , genders: %v ID: %v birthDate: %v files: %v\n", 
                      familyMemberName, familyMemberRelationship, familyMemberContactNumber, familyMemberGender, familyMemberGovernmentId, familyMemberBirthDate, familyMemberAttachment)
            // Similar loop for family members
                for i := range familyMemberName {   
                    var familyFilePath string
                    if i < len(familyMemberAttachment) {
                        fileHeader := familyMemberAttachment[i]
                        familyFilePath = UploadEmployeeAttachment(fileHeader, badgeId, "familyMembers") // Reusing the same function for simplicity
                        fmt.Printf("Received family member attachment: %s\n", fileHeader.Filename)
                    } else {
                        fmt.Printf("No attachment for family member: %s\n", familyMemberName[i])
                        familyFilePath = ""
                    }

                    birthDateParsed, err := time.Parse("2006-01-02", familyMemberBirthDate[i])
                    if err != nil {
                        fmt.Printf("Invalid birth date format for family member: %v\n", err)
                        birthDateParsed = time.Time{} // Zero value if parsing fails
                    }

                    familyMember := models.FamilyMember{
                        Name:          familyMemberName[i],
                        Relationship:  familyMemberRelationship[i],
                        ContactNumber: &familyMemberContactNumber[i],
                        GovernmentId:  &familyMemberGovernmentId[i],
                        BirthDate:     birthDateParsed,
                        FilePath:      familyFilePath,
                    }
                    FamilyMembers = append(FamilyMembers, familyMember)
                } //end of loop for family members
        } else {
            fmt.Println("No family members received")
        }

         emergencyContactName := r.PostForm["emergencyContactName[]"]
         emergencyContactRelationship := r.PostForm["emergencyContactRelationship[]"]
         emergencyContactNumber := r.PostForm["emergencyContactNumber[]"]
        if len(emergencyContactName) > 0 {
            fmt.Printf("Received emergency contacts: %v relationships: %v numbers: %v\n", emergencyContactName, emergencyContactRelationship, emergencyContactNumber)
            // Similar loop for emergency contacts
                for i := range emergencyContactName {
                        contact:= models.EmergencyContact{
                        Name:         emergencyContactName[i],
                        Relationship: emergencyContactRelationship[i],
                        Phone:        emergencyContactNumber[i],
                    }
                    EmergencyContacts = append(EmergencyContacts, contact)
                } //end of loop for emergency contacts
        } else {
            fmt.Println("No emergency contacts received")
        }

        documentName        := r.PostForm["documentName[]"]
        documentType        := r.PostForm["documentType[]"]
        documentExpiryDate  := r.PostForm["documentExpiryDate[]"]
        documentAttachments := r.MultipartForm.File["documentAttachment[]"]
        if len(documentName) > 0 {
            fmt.Printf("Received employee documents: %v\n", EmployeeDocuments)
            // Similar loop for employee documents
                for i := range documentName {
                    var documentFilePath string
                    if i < len(documentAttachments) {
                        fileHeader := documentAttachments[i]
                        documentFilePath = UploadEmployeeAttachment(fileHeader, badgeId, "employeeDocuments") // Reusing the same function for simplicity
                        fmt.Printf("Received employee document attachment: %s\n", fileHeader.Filename)
                    } else {
                        fmt.Printf("No attachment for employee document: %s\n", documentName[i])
                        documentFilePath = ""
                    }

                    expiryDateParsed, err := time.Parse("2006-01-02", documentExpiryDate[i])
                    if err != nil {
                        fmt.Printf("Invalid expiry date format for employee document: %v\n", err)
                        expiryDateParsed = time.Time{} // Zero value if parsing fails
                    }

                    doc := models.EmployeeDocument{
                        Name:       documentName[i],
                        Type:       documentType[i],
                        ExpiryDate: expiryDateParsed,
                        FilePath:  documentFilePath,
                    }
                    EmployeeDocuments = append(EmployeeDocuments, doc)
                } //end of loop for employee documents
        } else {
            fmt.Println("No employee documents received")
        }


        

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
            Active:      active,
            Email:       &email,
            Nationality: &nationality,
            PhoneNumber: &phoneNumber,
            Address:       &address,
            MaritalStatus: &maritalStatus,
            Gender:          &gender,
            Certifications: Certifications,
            FamilyMembers: FamilyMembers,
            EmergencyContacts: EmergencyContacts,
            EmployeeDocuments: EmployeeDocuments,
            Education: &education,
            Major: &major,
            Religion: &religion,

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

func UploadEmployeeAttachment(fileHeader *multipart.FileHeader, badgeId string, documentType string) string {
    // 1. Use multipart.FileHeader instead of http.File for the parameter
    file, err := fileHeader.Open()
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return ""
    }
    defer file.Close()

    // 2. Ensure directory exists
    dirPath := fmt.Sprintf("./media/employees/%s/%s", documentType, badgeId)
    if err := os.MkdirAll(dirPath, 0755); err != nil {
        fmt.Printf("Error creating directory: %v\n", err)
        return ""
    }

    // 3. CRITICAL: Prevent filename collisions
    // If two employees upload "cert.pdf", the second will overwrite the first.
    // We'll prefix the filename with a timestamp or UUID.
    uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
    savedFilePath := filepath.Join(dirPath, uniqueName)

    out, err := os.Create(savedFilePath)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return ""
    }
    defer out.Close()

    if _, err := io.Copy(out, file); err != nil {
        fmt.Printf("Error saving file: %v\n", err)
        return ""
    }

    return savedFilePath
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

    // 3. Database Insert (Using $ placeholders for Postgres) we need the ID for the employee after insert to link certificates and documents, so we will use RETURNING id
    query := `
        INSERT INTO employees (
            badge_id, name, department_id, local_name, 
            job_title_id, grade, birth_date, active, 
            goverment_id, image, email, nationality, gender, marital_status, phone_number, address, education, major, religion
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id`

    var newID int
    err := core.DB.QueryRow(query, 
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
        employee.Email,
        employee.Nationality,
        employee.Gender,
        employee.MaritalStatus,
        employee.PhoneNumber,
        employee.Address,
        employee.Education,
        employee.Major,
        employee.Religion,
    ).Scan(&newID)

    if err != nil {
        fmt.Printf("Database insert error: %v\n", err)
        return err
    }

    // 4. Success logic
    fmt.Printf("Employee %s inserted successfully\n", employee.Name)
    UpdateSequenceNextValue(employee.BadgeID, "badge_id")

    // 5. Insert related certifications, documents, family members, and emergency contacts if needed
    // This is where you would loop through any certifications, documents, family members, and emergency contacts associated with the employee and call their respective insert functions, passing the new employee ID.
    
    for _, cert := range employee.Certifications {
        fmt.Print("\n Inserting certification: ", cert.Name, " for employee: ", employee.Name, "\n")
        cert.Employee = &models.Employee{ID: newID} // Set the employee ID for the certification
        if err := InsertEmployeeCertificateToDB(cert); err != nil {
            fmt.Printf("Error inserting certification: %v\n", err)
            // Handle error (e.g., continue, rollback, etc.)
        }
    }

    for _, familyM:= range employee.FamilyMembers {
        fmt.Print("\n Inserting family member: ", familyM.Name, " for employee: ", employee.Name, "\n")
        familyM.Employee = &models.Employee{ID: newID} // Set the employee ID for the family member
        if err := InsertEmployeeFamilyMemberToDB(familyM); err != nil {
            fmt.Printf("Error inserting family member: %v\n", err)
            // Handle error (e.g., continue, rollback, etc.)
        }
    }

    for _, contact := range employee.EmergencyContacts {
        fmt.Print("\n Inserting emergency contact: ", contact.Name, " for employee: ", employee.Name, "\n")
        contact.Employee = &models.Employee{ID: newID} // Set the employee ID for the emergency contact
        if err := InsertEmployeeEmergencyContactToDB(contact); err != nil {
            fmt.Printf("Error inserting emergency contact: %v\n", err)
            // Handle error (e.g., continue, rollback, etc.)
        }
    }

    for _, doc := range employee.EmployeeDocuments {
        fmt.Print("\n Inserting employee document: ", doc.Name, " for employee: ", employee.Name, "\n")
        doc.Employee = &models.Employee{ID: newID} // Set the employee ID for the document
        if err := InsertEmployeeDocumentToDB(doc); err != nil {
            fmt.Printf("Error inserting employee document: %v\n", err)
            // Handle error (e.g., continue, rollback, etc.)
        }
    }

    return nil
}



func InsertEmployeeCertificateToDB(certificate models.Certification) error {
    query := `INSERT INTO certifications (employee_id,name,issuer,issue_date,expiry_date,file_path) VALUES ($1, $2, $3, $4, $5, $6)`
    _, err := core.DB.Exec(query, certificate.Employee.ID, certificate.Name, certificate.Issuer,certificate.IssueDate, certificate.ExpiryDate, certificate.FilePath)
    if err != nil {
        fmt.Printf("Database insert error for certificate: %v\n", err)
        return err
    }
    fmt.Printf("Certificate for Employee ID %d inserted successfully\n", certificate.Employee.ID)
    return nil
}

func InsertEmployeeDocumentToDB(document models.EmployeeDocument) error {
    query := `INSERT INTO employee_documents (employee_id, name, type, file_path, expiry_date) VALUES ($1, $2, $3, $4, $5)`
    _, err := core.DB.Exec(query, document.Employee.ID, document.Name, document.Type, document.FilePath, document.ExpiryDate)
    if err != nil {
        fmt.Printf("Database insert error for employee document: %v\n", err)
        return err
    }
    fmt.Printf("Document for Employee ID %d inserted successfully\n", document.Employee.ID)
    return nil
}

func InsertEmployeeFamilyMemberToDB(familyMember models.FamilyMember) error {
    query := `INSERT INTO family_members (employee_id, name, relationship, birth_date, file_path) VALUES ($1, $2, $3, $4, $5)`
    _, err := core.DB.Exec(query, familyMember.Employee.ID, familyMember.Name, familyMember.Relationship, familyMember.BirthDate, familyMember.FilePath)
    if err != nil {
        fmt.Printf("Database insert error for family member: %v\n", err)
        return err
    }
    fmt.Printf("Family member for Employee ID %d inserted successfully\n", familyMember.Employee.ID)
    return nil
}

func InsertEmployeeEmergencyContactToDB(contact models.EmergencyContact) error {
    query := `INSERT INTO emergency_contacts (employee_id, name, relationship, phone) VALUES ($1, $2, $3, $4)`
    _, err := core.DB.Exec(query, contact.Employee.ID, contact.Name, contact.Relationship, contact.Phone)
    if err != nil {
        fmt.Printf("Database insert error for emergency contact: %v\n", err)
        return err
    }
    fmt.Printf("Emergency contact for Employee ID %d inserted successfully\n", contact.Employee.ID)
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


//TODO we need only ID, Name, LocalName, BadgeID to show in the select in user create form, we can create a new struct with only these fields and use it in the query to optimize the performance instead of scanning all fields of employee struct
func GetEmployeesListApi(w http.ResponseWriter, r *http.Request) {

	 totalRecords := core.GetCountRecords("employees")
	 
	employees := GetEmployeesFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}




func GetEmployeeDetails(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Printf("Fetching details for employee ID: %s\n", id)

    employee := GetEmployeeById(id)
    if employee.ID == 0 {
        http.Error(w, "Employee not found", http.StatusNotFound)
        return
    }

    data := map[string]interface{}{
        "Title": "Employee Details",
        "Employee": employee,
    }

    core.RenderPage(w, r, "apps/employees/views/employees-details.html", data)
}


func UpdateEmployee(w http.ResponseWriter, r *http.Request) {

    fmt.Printf("Received %s request for updating employee\n", r.Method)

    if r.Method == http.MethodPost {
    id := r.PathValue("id")
    fmt.Printf("Updating employee ID: %s\n", id)

    employee := GetEmployeeById(id)

    fmt.Printf("Current employee details: %+v\n", employee)

    
    departmentId := r.PostFormValue("departmentId")
    jobTitleId := r.PostFormValue("jobTitleId")

    fmt.Printf("Received update data - Department ID: %s, Job Title ID: %s\n", departmentId, jobTitleId)

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

    //we need to check if the department or job has been updated not the same as old value if new value then we need to update
    // ExJobTitle & 

    if employee.Department != nil && theDepartment != nil {
        if employee.Department.ID != theDepartment.ID {
            fmt.Printf("Department changed from %s to %s\n", employee.Department.Name, theDepartment.Name)
        } else {
            fmt.Printf("Department remains unchanged: %s\n", employee.Department.Name)
        }
    } else if employee.Department == nil && theDepartment != nil {
        fmt.Printf("Department set to %s\n", theDepartment.Name)
    } else if employee.Department != nil && theDepartment == nil {
        fmt.Printf("Department cleared from %s\n", employee.Department.Name)
    } else {
        fmt.Println("Department remains unchanged: nil")
    }

    if employee.JobTitle != nil && theJobTitle != nil {
        if employee.JobTitle.ID != theJobTitle.ID {
            fmt.Printf("Job Title changed from %s to %s\n", employee.JobTitle.Name, theJobTitle.Name)
        } else {
            fmt.Printf("Job Title remains unchanged: %s\n", employee.JobTitle.Name)
        }
    } else if employee.JobTitle == nil && theJobTitle != nil {
        fmt.Printf("Job Title set to %s\n", theJobTitle.Name)
    } else if employee.JobTitle != nil && theJobTitle == nil {
        fmt.Printf("Job Title cleared from %s\n", employee.JobTitle.Name)
    } else {
        fmt.Println("Job Title remains unchanged: nil")
    }

    employee.Department = theDepartment
    employee.JobTitle = theJobTitle

    fmt.Printf("Updated employee details to be saved: %+v\n", employee)



    SaveEmployeeObj(employee)




    http.Redirect(w, r, "/employees/details/"+id, http.StatusSeeOther)
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
}



func SetExDepartment(employee models.Employee, newDepartment *models.Department, oldDepartment *models.Department) { 



}


func SaveEmployeeObj(employee models.Employee) error {
   
    //here you pass employee struct with all fields to the database and update the employee record in the database with the new values, you can use sql UPDATE statement to update the employee record in the database based on employee.ID and set all fields of employee struct to the database record, you can also use transaction if you want to update multiple tables related to employee like certifications, documents, family members, emergency contacts in one transaction to ensure data integrity

    fmt.Printf("Saving employee to database.......: %+v\n", employee)

    query := `UPDATE employees SET 
    badge_id = $1, 
    name = $2, 
    department_id = $3, 
    local_name = $4, 
    job_title_id = $5, 
    grade = $6, 
    birth_date = $7, 
    active = $8, 
    goverment_id = $9, 
    image = $10, 
    email = $11, 
    nationality = $12
    WHERE id = $13`

    _, err := core.DB.Exec(query, employee.BadgeID, employee.Name, employee.Department.ID, employee.LocalName, employee.JobTitle.ID, employee.Grade, employee.BirthDate, employee.Active, employee.GovermentID, employee.Image, employee.Email, employee.Nationality, employee.ID)
    if err != nil {
        return err
    }


    return nil
}