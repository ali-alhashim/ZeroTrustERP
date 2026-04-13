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
		"User":  getUserByID(userID),
	}

	core.RenderPage(w, r, "apps/users/views/user-details.html", data)
}

func getUserByID(id string) models.User {
    // 1. Join all 5 tables to get the full hierarchy
    query := `
        SELECT 
            u.id, u.username, u.email, u.active, u.online, u.last_login, u.created_at, u.updated_at,
            r.id, r.name, r.description,
            p.id, p.resource, p.action
        FROM users u
        LEFT JOIN users_roles ur ON u.id = ur.user_id
        LEFT JOIN roles r ON ur.role_id = r.id
        LEFT JOIN roles_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = $1`

    rows, err := core.DB.Query(query, id)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return models.User{}
    }
    defer rows.Close()

    var user models.User
    // Use maps to track unique Roles and Permissions as we iterate rows
    roleMap := make(map[int]*models.Role)
    // Map to keep track of which permissions belong to which role
    permMap := make(map[int]map[int]bool) 

    for rows.Next() {
        var (
            // Using pointers for nullable JOIN columns
            rID *int
            rName *string
            rDesc *string
            pID *int
            pResource *string
            pAction *string
        )

        err := rows.Scan(
            &user.ID, &user.Username, &user.Email, &user.Active, &user.Online, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
            &rID, &rName, &rDesc,
            &pID, &pResource, &pAction,
        )
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            return models.User{}
        }

        // Handle Roles
        if rID != nil {
            role, exists := roleMap[*rID]
            if !exists {
                role = &models.Role{
                    ID:          *rID,
                    Name:        *rName,
                    Description: *rDesc,
                    Permissions: &[]models.Permission{}, // Initialize the pointer to a slice
                }
                roleMap[*rID] = role
                permMap[*rID] = make(map[int]bool)
            }

            // Handle Permissions inside the Role
            if pID != nil && !permMap[*rID][*pID] {
                *role.Permissions = append(*role.Permissions, models.Permission{
                    ID:       *pID,
                    Resource: *pResource,
                    Action:   *pAction,
                })
                permMap[*rID][*pID] = true
            }
        }
    }

    // Convert the map of roles into the user's *[]Role slice
    rolesSlice := []models.Role{}
    for _, role := range roleMap {
        rolesSlice = append(rolesSlice, *role)
    }
    user.Roles = &rolesSlice

    return user
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
