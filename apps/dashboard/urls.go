package dashboard

import (
	"net/http"

	"zerotrusterp/apps/dashboard/controllers"
	
)
func dashboardRoute(mux *http.ServeMux)  {

	mux.HandleFunc("GET /dashboard", controllers.DashboardController)
	
}