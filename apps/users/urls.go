package users

import (
	"net/http"

	"zerotrusterp/apps/users/controllers"
	"zerotrusterp/core"
)



func UserListRoutes(mux *http.ServeMux) {
	
               //HTTP method + path, handler + auth middleware with resource name for permission checking
	mux.Handle("GET /users/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListUsers), "users:list"))
	mux.Handle("GET /users/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateUser), "users:create"))
	mux.Handle("POST /users/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateUser), "users:create"))
	mux.Handle("GET /api/online-users", core.AuthMiddleware(http.HandlerFunc(controllers.OnlineUsersAPI), "users:online"))
	mux.Handle("GET /users/audit-trail", core.AuthMiddleware(http.HandlerFunc(controllers.ListLogs), "users:audit"))
	mux.Handle("GET /users/roles", core.AuthMiddleware(http.HandlerFunc(controllers.ListRoles), "roles:list"))
	mux.Handle("GET /roles/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateRole), "roles:create"))
	mux.Handle("POST /roles/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateRole), "roles:create"))
	mux.Handle("GET /api/fetch-roles", core.AuthMiddleware(http.HandlerFunc(controllers.FetchRolesAPI), "roles:list"))
	mux.Handle("POST /users/RevokeSession/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.RevokeSession),"users:session"))

}



