package system

import(
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/system/controllers"
)

func SystemRoutes(mux *http.ServeMux) {

	mux.Handle("GET /system", core.AuthMiddleware(http.HandlerFunc(controllers.Health), "system:R"))

}