package controllers

import(
	"fmt"
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/users/models"
)


func UserDetails(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
		fmt.Print("only GET Request Allowed")
		http.Error(w, "only GET Request Allowed", http.StatusBadRequest)
	}

	userID := r.PathValue("id")


	data := map[string]interface{}{
		"Title": "User Details",
		"User":getUserByID(userID),
	}


	core.RenderPage(w,r, "apps/users/views/user-details.html", data)
}


func getUserByID(id string) models.User {

	fmt.Print(" \n Get User By ID", id ,"\n")


    query := "SELECT id, username, email, active, online, last_login, created_at, updated_at FROM users WHERE id=$1"

    // Use QueryRow for a single result
    row := core.DB.QueryRow(query, id)

    var user models.User
    
    // Scan the result into the user struct
    // Ensure LastLogin in models.User is *time.Time to handle NULLs
    err := row.Scan(
        &user.ID, 
        &user.Username, 
        &user.Email, 
        &user.Active, 
        &user.Online, 
        &user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
    )

    if err != nil {
        // If no user is found or database error occurs, log it
        fmt.Printf("Error fetching user %s: %v\n", id, err)
        return models.User{} // Return an empty struct
    }

    return user
}



func UpdateUser(w http.ResponseWriter, r *http.Request){

    if r.Method != http.MethodPatch {
        fmt.Print("only Patch Request Allowed")
		http.Error(w, "only Patch Request Allowed", http.StatusBadRequest)
    }

    fmt.Print("Update User with Patch Request .....")

    userID := r.PathValue("id")

    fmt.Print("\n  Ok we have request to update user ID: ", userID ,"\n")

    query:= "select id, email, username, active, related_employee_id from users where id = $1"

    var user1 models.User

     err := core.DB.QueryRow(query, userID).Scan(&user1.ID, &user1.Email, &user1.Username, &user1.Active, &user1.RelatedEmployee)
    if err !=nil{
        fmt.Print("error in query !", err)
    }
   

    // ok log old data brfore the update in log and after the update
    Username := r.FormValue("Username")
    Email    := r.FormValue("Email")
    Active   := r.FormValue("Active") == "on"



    var CurrentUser *models.User

		if user, ok := r.Context().Value(core.UserKey).(*models.User); ok {
			CurrentUser = user
		} else {
			fmt.Println("No user in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}


    fmt.Print(CurrentUser.ID ," : " , CurrentUser.Username ,"  sent the following ", Username, " : ", Email, " : ", Active, " to update the User ID : ", userID)

     url := fmt.Sprintf("/users/details/%s", userID)
     http.Redirect(w, r, url, http.StatusSeeOther)
     

}