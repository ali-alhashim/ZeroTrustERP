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
		employeeId:=r.FormValue("manager")
		active:=r.FormValue("active") == "on"

		fmt.Printf("Create Department code:%s with name:%s -Ar Name:%s and the managerId:%s -status:%s", code, name, nameAR, employeeId, active)
	}
}