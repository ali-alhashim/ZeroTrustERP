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
	mux.Handle("POST /api/heartbeat", core.AuthMiddleware(http.HandlerFunc(controllers.SetUserOnline), "users:heartbeat"))
	mux.Handle("POST /api/stopHeartbeat", core.AuthMiddleware(http.HandlerFunc(controllers.StopUserHeartbeat), "users:heartbeat"))

}



