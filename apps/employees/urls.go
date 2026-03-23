package employees

import (
	"net/http"

	"zerotrusterp/apps/employees/controllers"
	ccore "zerotrusterp/core"
)

func init() {
	ccore.Register(RegisterRoutes)
}

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /employees/list", controllers.ListEmployees)

	mux.HandleFunc("GET /employees/departments", controllers.ListDepartments)
	
}