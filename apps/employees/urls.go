package employees

import (
	"net/http"

	"zerotrusterp/apps/employees/controllers"
	
)


func EmployeeListRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /employees/list", controllers.ListEmployees)

	mux.HandleFunc("GET /employees/departments", controllers.ListDepartments)
	
}