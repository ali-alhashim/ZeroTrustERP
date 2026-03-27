package employees

import (
	"net/http"

	"zerotrusterp/apps/employees/controllers"
	"zerotrusterp/core"
)


func EmployeeListRoutes(mux *http.ServeMux) {

	mux.Handle("GET /employees/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListEmployees), "employees:list"))

	mux.Handle("GET /employees/departments", core.AuthMiddleware(http.HandlerFunc(controllers.ListDepartments), "departments:list"))
	
}