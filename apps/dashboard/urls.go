package dashboard

import (
	"net/http"

	"zerotrusterp/apps/dashboard/dashboardControllers"
	"zerotrusterp/core"
)
func dashboardRoute(mux *http.ServeMux)  {

	mux.Handle("GET /dashboard", core.AuthMiddleware(http.HandlerFunc(dashboardControllers.DashboardController)))
	
}