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

	fmt.Print("Update User with Patch Request .....")

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

	var CurrentUser *models.User

	if user, ok := r.Context().Value(core.UserKey).(*models.User); ok {
		CurrentUser = user
	} else {
		fmt.Println("No user in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Print(CurrentUser.ID, " : ", CurrentUser.Username, "  sent the following ", Username, " : ", Email, " : ", Active, " to update the User ID : ", userID)

	//also see the roles and update if the user made any update

	url := fmt.Sprintf("/users/details/%s", userID)
	http.Redirect(w, r, url, http.StatusSeeOther)

}
