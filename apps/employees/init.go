package employees

import (
	"zerotrusterp/core"
	"zerotrusterp/apps/employees/models"
)

func init() {

	// Register routes
	core.Register(EmployeeListRoutes)

	// Register models for migrations
	core.RegisterModel(models.Employee{})
	core.RegisterModel(models.Department{})
}