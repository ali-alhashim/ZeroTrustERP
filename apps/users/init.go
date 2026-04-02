package users

import (
    "zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)
	

func init() {

    // Register routes
    core.Register(UserListRoutes)

    // Register models for migrations
    core.RegisterModel(models.User{})
    core.RegisterModel(models.Role{})
    core.RegisterModel(models.Permission{})
    core.RegisterModel(models.Log{})
}
