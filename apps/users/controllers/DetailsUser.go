package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/apps/users/models"
	"zerotrusterp/core"
)

func UserDetails(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		fmt.Print("only GET Request Allowed")
		http.Error(w, "only GET Request Allowed", http.StatusBadRequest)
	}

	userID := r.PathValue("id")

	data := map[string]interface{}{
		"Title": "User Details",
		"User":  core.GetUserByID(userID),
	}

	core.RenderPage(w, r, "apps/users/views/user-details.html", data)
}



//---------------------------End Get user by ID

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		fmt.Print("only Patch Request Allowed")
		http.Error(w, "only Patch Request Allowed", http.StatusBadRequest)
	}

	fmt.Print("\n...............Update User with Patch Request .....\n")

	userID := r.PathValue("id")

	fmt.Print("\n  Ok we have request to update user ID: ", userID, "\n")

	query := "select id, email, username, active, related_employee_id from users where id = $1"

	var user1 models.User

	err := core.DB.QueryRow(query, userID).Scan(&user1.ID, &user1.Email, &user1.Username, &user1.Active, &user1.RelatedEmployee)
	if err != nil {
		fmt.Print("error in query !", err)
	}

	// ok log old data brfore the update in log and after the update
	Username := r.FormValue("Username")
	Email := r.FormValue("Email")
	Active := r.FormValue("Active") == "on"
	selectedRoles := r.Form["Role"]

	CurrentUser := core.GetCurrentUser(r)

	fmt.Print(CurrentUser.ID, " : ", CurrentUser.Username, "  sent the following ", Username, " : ", Email, " : ", Active, " to update the User ID : ", userID)
    fmt.Print("\n selectedRoles=", selectedRoles, "\n")


	// 4. Update the User basic info
	tx, err := core.DB.Begin()
    updateUserQuery := "UPDATE users SET username=$1, email=$2, active=$3 WHERE id=$4"
    _, err = tx.Exec(updateUserQuery, Username, Email, Active, userID)
    if err != nil {
        fmt.Println("User update error:", err)
        return
    }


	for _, roleID := range selectedRoles {
        if roleID == "" { continue }
        _, err = tx.Exec("INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
        if err != nil {
            fmt.Println("Role insertion error:", err)
            return
        }
    }
    


	// Commit everything
    if err := tx.Commit(); err != nil {
        http.Error(w, "Failed to save changes", http.StatusInternalServerError)
        return
    }
	//active is a checkbox input checked = true not checked false
    // if request has Role [] then assign the role to the user by insert record in users_roles (user_id, role_id)


	//also see the roles and update if the user made any update

	url := fmt.Sprintf("/users/details/%s", userID)
	http.Redirect(w, r, url, http.StatusSeeOther)

}
