package controllers

import(
	"net/http"
	"fmt"
	"zerotrusterp/core"
	"zerotrusterp/apps/employees/models"
	"strconv"
)

func ListSalaryComponentTypes(w http.ResponseWriter, r *http.Request) {


	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)

	
    
    totalRecords := core.GetCountRecords("salary_component_types")

	SalaryComp:= GetSalaryComponentTypesFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Salary Component Types",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"SalaryComp": SalaryComp,
		
	}

	core.RenderPage(w,r, "apps/employees/views/componentTypes-list.html", data)
}


func GetSalaryComponentTypesFromDB(search, sort, order, page, pageSize string) []models.SalaryComponentType{

	query :="select id, name,  code from salary_component_types where 1=1"

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


	var components []models.SalaryComponentType

	for rows.Next() {
		var compo models.SalaryComponentType
		err := rows.Scan(&compo.ID, &compo.Name, &compo.Code)
		if err != nil {
			fmt.Print(err)
		}
		components = append(components, compo)
	}

	return components


}


func CreateSalaryComponentTypes(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodGet{
        data := map[string]interface{}{
		"Title": "Salary Component Types",
		}

		core.RenderPage(w,r, "apps/employees/views/componentTypes-create.html", data)
	}

	if r.Method == http.MethodPost {

		code:=r.PostFormValue("code")
		name:=r.PostFormValue("name")

			query := "insert into salary_component_types (name,  code) values ($1, $2)"

		_,err:= core.DB.Exec(query, name,  code)

		if err!=nil{
			fmt.Print(err)
		}

		http.Redirect(w, r, "/salary/componentTypes", http.StatusSeeOther)

	}

}


func DetailsSalaryComponentTypes(w http.ResponseWriter, r *http.Request){

}

