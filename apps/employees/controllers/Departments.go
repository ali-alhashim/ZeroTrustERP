package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/core"
)

func ListDepartments(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Departments",
		
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


		core.InsertLog(CurrentUser, "Departments", fmt.Sprintf("Created Department coade %s with name: %s",code, name))

		http.Redirect(w, r, "/employees/departments", http.StatusSeeOther)

	}
}