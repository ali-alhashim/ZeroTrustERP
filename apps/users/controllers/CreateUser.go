package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Create User",
		}
		core.RenderPage(w,r, "apps/users/views/users-create.html", data)
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
           
		 var userID int64
		
		 err := core.DB.QueryRow(
		"INSERT INTO users (username, email, active, related_employee_id) VALUES ($1, $2, $3, $4) RETURNING id",
		username, 
		email, 
		active, 
		relatedEmployee,
		).Scan(&userID)


		

		
		if err != nil {
			fmt.Println("Error inserting user:", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}


		// insert user roles in users_roles table set user_id & role_id
		// roles ids from Role
		RolesIds := r.Form["Role"]

		fmt.Printf("set the following Roles Ids %s to new user with Id %d", RolesIds, userID)

		for i:=0; i<len(RolesIds); i++ {

			sql := "INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)"
			err := core.DB.QueryRow(sql, userID, RolesIds[i])

			if err !=nil{
				fmt.Printf("Failed %v \n", err)
			}
		}




		// log the user creation action

		var CurrentUser *models.User

		if user, ok := r.Context().Value(core.UserKey).(*models.User); ok {
			CurrentUser = user
		} else {
			fmt.Println("No user in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}


		InsertLog(CurrentUser, "Users", fmt.Sprintf("Created user %s with email: %s",username, email))



		http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
