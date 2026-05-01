package empapi

import (
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"fmt"
)


func GetEmployeeByID(employeeID int) models.Employee {

	query := "select id, badge_id, name, department_id, local_name, job_title_id, grade, created_at, updated_at, birth_date, active, goverment_id, image FROM employees WHERE id = $1"
	row := core.DB.QueryRow(query, employeeID)

	var emp models.Employee
	err := row.Scan(&emp.ID, &emp.BadgeID, &emp.Name, &emp.Department, &emp.LocalName, &emp.JobTitle, &emp.Grade, &emp.CreatedAt, &emp.UpdatedAt, &emp.BirthDate, &emp.Active, &emp.GovermentID, &emp.Image)
	if err != nil {
		fmt.Print(err)
	}

	return emp
}