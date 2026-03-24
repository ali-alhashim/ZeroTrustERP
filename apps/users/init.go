package users

import (
    "zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)
	

func init() {

    // Register routes
    core.Register(UserListRoutes)

    // Register models for migrations
    core.RegisterModel(models.Users{})
    core.RegisterModel(models.Roles{})
    core.RegisterModel(models.Permissions{})
}
