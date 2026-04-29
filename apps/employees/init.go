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

	core.RegisterModel(models.JobTitle{})
	core.RegisterModel(models.ExJobTitle{})

	core.RegisterModel(models.ExManagerDepartment{})
	core.RegisterModel(models.OrgUnit{})

	core.RegisterModel(models.Contract{})
	core.RegisterModel(models.ContractSalaryLine{})
	core.RegisterModel(models.SalaryComponentType{})
	core.RegisterModel(models.SalaryComponentValue{})

	core.RegisterModel(models.InsuranceGrade{})
	core.RegisterModel(models.InsurancePolicy{})

	core.RegisterModel(models.ShiftSchedule{})

	core.RegisterModel(models.Certification{})
	core.RegisterModel(models.FamilyMember{})
	core.RegisterModel(models.EmergencyContact{})
	core.RegisterModel(models.EmployeeDocument{})

}