package employees

import (
	"net/http"

	"zerotrusterp/apps/employees/controllers"
	"zerotrusterp/core"
)


func EmployeeListRoutes(mux *http.ServeMux) {

	mux.Handle("GET /employees/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListEmployees), "employees:R"))

	mux.Handle("GET /employees/departments", core.AuthMiddleware(http.HandlerFunc(controllers.ListDepartments), "departments:R"))

	mux.Handle("GET /employees/Jobs", core.AuthMiddleware(http.HandlerFunc(controllers.ListJobs), "jobs:R"))

	mux.Handle("GET /employees/departments/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateDepartment), "departments:R"))
	
}