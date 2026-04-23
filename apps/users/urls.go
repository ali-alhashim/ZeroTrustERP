package users

import (
	"net/http"

	"zerotrusterp/apps/users/controllers"
	"zerotrusterp/core"
)



func UserListRoutes(mux *http.ServeMux) {
	
               //HTTP method + path, handler + auth middleware with resource name and the Permission to access  for permission checking
	mux.Handle("GET /users/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListUsers), "users:R"))
	mux.Handle("GET /users/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateUser), "users:W"))
	mux.Handle("POST /users/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateUser), "users:W"))
	mux.Handle("GET /api/online-users", core.AuthMiddleware(http.HandlerFunc(controllers.OnlineUsersAPI), "users:R"))
	mux.Handle("GET /users/audit-trail", core.AuthMiddleware(http.HandlerFunc(controllers.ListLogs), "logs:R"))
	mux.Handle("GET /users/roles", core.AuthMiddleware(http.HandlerFunc(controllers.ListRoles), "roles:R"))
	mux.Handle("GET /roles/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateRole), "roles:W"))
	mux.Handle("POST /roles/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateRole), "roles:W"))
	mux.Handle("GET /api/fetch-roles", core.AuthMiddleware(http.HandlerFunc(controllers.FetchRolesAPI), "roles:R"))
	mux.Handle("POST /users/RevokeSession/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.RevokeSession),"users:U"))

	mux.Handle("GET /users/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UserDetails),"users:R"))
	mux.Handle("PATCH /users/update/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateUser), "users:U"))
    mux.Handle("DELETE /api/users/delete-role/{roleID}/{userID}", core.AuthMiddleware(http.HandlerFunc(controllers.DeleteRoleFromUser),"roles:D"))

	mux.Handle("POST /roles/update/{roleID}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateRole), "roles:U"))

	mux.Handle("DELETE /api/roles/delete-permission/{roleID}/{permissionID}", core.AuthMiddleware(http.HandlerFunc(controllers.DeletePermissionFromRole),"roles:D"))

	mux.Handle("GET /users/status/{id}/Active", core.AuthMiddleware(http.HandlerFunc(controllers.SetUserActive), "users:U"))
	mux.Handle("GET /users/status/{id}/Inactive", core.AuthMiddleware(http.HandlerFunc(controllers.SetUserInactive), "users:U"))

	mux.Handle("GET /role/details/{roleID}", core.AuthMiddleware(http.HandlerFunc(controllers.RoleDeatils), "roles:U"))
}



