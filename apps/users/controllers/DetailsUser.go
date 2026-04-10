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

	fmt.Print("Get User By ID", id ,"\n")


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