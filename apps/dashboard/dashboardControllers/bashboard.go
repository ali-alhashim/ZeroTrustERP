package dashboardControllers

import (
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/apps/employees/controllers"
	"fmt"
	"strconv"
)



func DashboardController(w http.ResponseWriter, r *http.Request) {


	User := core.GetCurrentUser(r)
	if User == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var Employee models.Employee

	if User.RelatedEmployee != nil && User.RelatedEmployee.ID != 0 {
		// User has a related employee, you can access employee details here
		fmt.Printf("User %s is related to employee ID: %d\n", User.Username, User.RelatedEmployee.ID)
        employeeIDStr := strconv.Itoa(User.RelatedEmployee.ID)
		Employee = controllers.GetEmployeeById(employeeIDStr)

		fmt.Printf("Retrieved employee details for ID: %d\n", Employee.ID)
	}




	data := map[string]interface{}{
		"Title": "Dashboard",
		"User":  User,
		"Employee": Employee,
	}

	core.RenderPage(w,r, "apps/dashboard/views/dashboard.html", data)

}