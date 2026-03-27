package users

import (
	"net/http"

	"zerotrusterp/apps/users/controllers"
	"zerotrusterp/core"
)



func UserListRoutes(mux *http.ServeMux) {
	
               //HTTP method + path, handler + auth middleware with resource name for permission checking
	mux.Handle("GET /users/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListUsers), "users:list"))
	
}



