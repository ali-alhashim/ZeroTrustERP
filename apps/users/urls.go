package users

import (
	"net/http"

	"zerotrusterp/apps/users/controllers"
	ccore "zerotrusterp/core"
)

func init() {
	ccore.Register(RegisterRoutes)
}

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /users/list", controllers.ListUsers)
	
}



