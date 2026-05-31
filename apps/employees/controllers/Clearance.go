package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/employees/models"
	"strconv"
	"encoding/json"
)

func ListClearance(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)

	
    
    totalRecords := core.GetCountRecords("clearance_documents")

	clearanceDocuments:= GetClearanceDocumentsFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Clearance Documents",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"ClearanceDocuments": clearanceDocuments,
		
	}

	core.RenderPage(w,r, "apps/employees/views/Clearance-list.html", data)
}

func GetClearanceDocumentsFromDB(search, sort, order, page, pageSize string) []models.ClearanceDocument {

	query :="select id, employee_id, template_id, status, requested_date,completed_date, created_at, updated_at from clearance_documents where 1=1"

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (status ILIKE $" + strconv.Itoa(argIndex) +
			     " OR employee_id ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":          "id",
		"Status":        "status",
		"Name":        "name",
		"RequestedDate":   "requested_date",
		"CompletedDate":  "completed_date",
		"CreatedAt":  "created_at",
		"UpdatedAt":  "updated_at",
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


	var clearanceDocuments []models.ClearanceDocument

	for rows.Next() {
		var l models.ClearanceDocument
		var employeeID, templateID string
		err := rows.Scan(&l.ID, &employeeID, &templateID, &l.Status, &l.RequestedDate, &l.CompletedDate, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			fmt.Print(err)
		}

		Employee:= GetEmployeeById(employeeID)
		Template:= GetClearanceTemplateById(templateID)

		l.Template = &Template

		l.Employee = &Employee

		

		clearanceDocuments = append(clearanceDocuments, l)
	}

	return clearanceDocuments


}

func GetClearanceTemplateById(id string) models.ClearanceTemplate {
	var t models.ClearanceTemplate
	var tItems []models.ClearanceTemplateItem

	query := `select t.id, t.name, t.description, i.id, i.template_id, i.name, i.department_id
	from clearance_templates t
	left join clearance_template_items i on t.id = i.template_id
	where t.id = $1`

	rows, err := core.DB.Query(query, id)
    if err != nil {
       fmt.Printf("Error in GetClearanceTemplateById: %v", err)
    }
    defer rows.Close()

	
    for rows.Next() {
		var item models.ClearanceTemplateItem
        var departmentId string
		

		err := rows.Scan(&t.ID, &t.Name, &t.Description, &item.ID, &t.ID, &item.Name, &departmentId)
		if err != nil {
			fmt.Print(err)
		}

		Department := GetDepartmentById(departmentId)
		item.Department = &Department
		

		tItems = append(tItems, item)
	}

	t.Items = &tItems
	return t
}




func ListClearanceTemplates(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)

	
    
    totalRecords := core.GetCountRecords("clearance_templates")

	clearanceTemplates:= GetClearanceTemplatesFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Clearance Templates",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"ClearanceTemplates": clearanceTemplates,
		
	}

	core.RenderPage(w,r, "apps/employees/views/Clearance-Templates-list.html", data)
}




func GetClearanceTemplatesFromDB(search, sort, order, page, pageSize string) []models.ClearanceTemplate {

	query :="select id, name, description from clearance_templates where 1=1"

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (name ILIKE $" + strconv.Itoa(argIndex) +
			     " OR description ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":          "id",
		"Name":        "name",
		"Description": "description",
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


	var clearanceTemplates []models.ClearanceTemplate

	for rows.Next() {
		var l models.ClearanceTemplate
		
		err := rows.Scan(&l.ID, &l.Name, &l.Description)
		if err != nil {
			fmt.Print(err)
		}

		Items := GetClearanceTemplateItemsById(strconv.Itoa(l.ID))

		l.Items = &Items

		

		

		

		clearanceTemplates = append(clearanceTemplates, l)
	}

	return clearanceTemplates


}


func CreateClearanceTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := map[string]interface{}{
			"Title": "Create Clearance Template",
		}
	
		core.RenderPage(w,r, "apps/employees/views/Clearance-Template-create.html", data)
	}

	if r.Method == "POST" {

		name:= r.PostFormValue("name")
		description := r.PostFormValue("description")

		item_names := r.PostForm["item_name"]
		department_ids := r.PostForm["item_department"]

		fmt.Printf("Received form data: name=%s, description=%s\n", name, description)
		fmt.Printf("Received form data: item_names=%v, department_ids=%v\n", item_names, department_ids)

		// Insert template into clearance_templates (name, description) returning id
		// loop through item_names and department_ids and insert into clearance_template_items (template_id, name, department_id)
		
		var templateID string
		err := core.DB.QueryRow(`INSERT INTO clearance_templates (name, description) VALUES ($1, $2) RETURNING id`, name, description).Scan(&templateID)
		if err != nil {
			fmt.Print(err)
		}

		for i := range item_names {
			itemName := item_names[i]
			departmentID := department_ids[i]

			_, err := core.DB.Exec(`INSERT INTO clearance_template_items (template_id, name, department_id) VALUES ($1, $2, $3)`, templateID, itemName, departmentID)
			if err != nil {
				fmt.Print(err)
			}
		}

		http.Redirect(w, r, "/employees/Clearance-Templates", http.StatusSeeOther)

	}
}

func GetClearanceTemplateItemsById(templateId string) []models.ClearanceTemplateItem {
	var items []models.ClearanceTemplateItem

	rows, err := core.DB.Query(`SELECT id, name, department_id FROM clearance_template_items WHERE template_id = $1`, templateId)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ClearanceTemplateItem
		var departmentId string
		err := rows.Scan(&item.ID, &item.Name, &departmentId)
		if err != nil {
			fmt.Print(err)
		}
		Department := GetDepartmentById(departmentId)
		item.Department = &Department
		items = append(items, item)
	}

	return items
}


func DetailsClearanceTemplate(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if r.Method == "GET" {

		data := map[string]interface{}{
			"Title": "Clearance Template Details",
			"Template": GetClearanceTemplateById(id),
		}
	
		core.RenderPage(w,r, "apps/employees/views/Clearance-Template-details.html", data)

	}

	if r.Method == "POST" {

		name:= r.PostFormValue("name")
		description := r.PostFormValue("description")

		_, err := core.DB.Exec(`UPDATE clearance_templates SET name = $1, description = $2 WHERE id = $3`, name, description, id)
		if err != nil {
			fmt.Print(err)
		}

		http.Redirect(w, r, "/employees/Clearance-Templates", http.StatusSeeOther)
	}


}


func GetClearanceTemplateListApi(w http.ResponseWriter, r *http.Request) {

	totalRecords := core.GetCountRecords("clearance_templates")
	 
	clearanceTemplates := GetClearanceTemplatesFromDB("", "ID", "asc", "1", strconv.Itoa(totalRecords))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clearanceTemplates)
}

func GetServicePeriodById(id string) models.ServicePeriod {

	var sp models.ServicePeriod
	var employeeID  string

	err := core.DB.QueryRow(`SELECT id, employee_id, hire_date, termination_date, reason, eosb_paid, note FROM service_periods WHERE id = $1`, id).Scan(&sp.ID, &employeeID, &sp.HireDate, &sp.TerminationDate, &sp.Reason, &sp.EOSBPaid, &sp.Note)
	if err != nil {
		fmt.Print(err)
	}

	Employee := GetEmployeeById(employeeID)
	sp.Employee = &Employee

	

	return sp
}