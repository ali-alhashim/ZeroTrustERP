package dashboard

import (
	"net/http"

	"zerotrusterp/apps/dashboard/controllers"
	"zerotrusterp/core"
)
func dashboardRoute(mux *http.ServeMux)  {

	mux.Handle("GET /dashboard", core.AuthMiddleware(http.HandlerFunc(controllers.DashboardController)))
	
}