package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"zerotrusterp/core"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Create User",
		}
		core.RenderPage(w, "apps/users/views/createUser.html", data)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("Username")
		email := r.FormValue("Email")
		active := r.FormValue("Active") == "on"

		relatedEmployeeStr := r.FormValue("RelatedEmployee")

		if username == "" || email == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// ✅ Handle NULL properly
		var relatedEmployee interface{}

		if relatedEmployeeStr == "" {
			relatedEmployee = nil
		} else {
			id, err := strconv.Atoi(relatedEmployeeStr)
			if err != nil {
				http.Error(w, "Invalid Related Employee ID", http.StatusBadRequest)
				return
			}
			relatedEmployee = id
		}

		_, err := core.DB.Exec(
			"INSERT INTO users (username, email, active, related_employee_id) VALUES ($1, $2, $3, $4)",
			username, email, active, relatedEmployee,
		)

		if err != nil {
			fmt.Println("Error inserting user:", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
