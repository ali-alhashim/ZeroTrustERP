package users

import (
	"net/http"

	"zerotrusterp/apps/users/controllers"
	
)



func UserListRoutes(mux *http.ServeMux) {
	
                   //HTTP method + path, handler function
	mux.HandleFunc("GET /users/list", controllers.ListUsers)
	
}



