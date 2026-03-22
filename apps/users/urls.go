package users

import (
	"net/http"
	"zerotrusterp/apps/users/controllers"
)

func RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /users/list", controllers.ListUsers)
}



