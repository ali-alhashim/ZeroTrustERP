package controllers

import (
	
	"fmt"
	"net/http"

	"zerotrusterp/core"
    "zerotrusterp/apps/users/models"
)

func CreateRole(w http.ResponseWriter, r *http.Request) {

	fmt.Println("CreateRole called with method:", r.Method)

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
			"Title": "Create Role",
			"Resources": core.GetAllResources(),
		}
		core.RenderPage(w,r, "apps/users/views/roles-create.html", data)
		return
	}

	// Handle POST request to create role

	if r.Method == http.MethodPost {

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    roleName := r.FormValue("roleName")
    roleDescription := r.FormValue("roleDescription")

    // 1. Use $1, $2 for PostgreSQL
    // Also, use "RETURNING id" if LastInsertId() is not supported by your driver
    var roleID int64
    sql := "INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id"
    err = core.DB.QueryRow(sql, roleName, roleDescription).Scan(&roleID)
    
    if err != nil {
        fmt.Printf("Failed to create role: %v\n", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

	fmt.Print("Role created with ID: ", roleID ,"\n")

    pDescriptions := r.Form["description"]
    pResources    := r.Form["Resource"]
    pActions      := r.Form["Action"]
    
    fmt.Printf("Received permissions: %v with Actions: %v\n", pResources, pActions)

    fmt.Printf("Count: %s\n", r.FormValue("permissionsCount"))

    for i := 0; i < len(pResources); i++ {
        var permID int64
        // 2. Again, use $1, $2, $3 for Postgres
        sql = "INSERT INTO permissions (resource, action, description) VALUES ($1, $2, $3) RETURNING id"
        // pDescriptions[i] might be empty, so we can handle that case if needed
        fmt.Printf("Inserting permission with Resource: %s, Action: %s, Description: %s\n", pResources[i], pActions[i], pDescriptions[i])
        err = core.DB.QueryRow(sql, pResources[i], pActions[i], pDescriptions[i]).Scan(&permID)
        
        if err != nil {
            fmt.Printf("Failed to create permission: %v\n", err)
            http.Error(w, "Failed to create permission", http.StatusInternalServerError)
            return
        }

		fmt.Printf("Permission created with ID: %d for resource: %s and action: %s\n", permID, pResources[i], pActions[i])

        // 3. Mapping table insert
        sql = "INSERT INTO roles_permissions (role_id, permission_id) VALUES ($1, $2)"
        _, err = core.DB.Exec(sql, roleID, permID)
        if err != nil {
            fmt.Print(err)
            http.Error(w, "Failed to assign permission", http.StatusInternalServerError)
            return
        }

        fmt.Printf("Assigned permission ID %d to role ID %d\n", permID, roleID)
    }



    var CurrentUser *models.User

		if user, ok := r.Context().Value(core.UserKey).(*models.User); ok {
			CurrentUser = user
		} else {
			fmt.Println("No user in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}


		InsertLog(CurrentUser, "Roles", fmt.Sprintf("Created Role Name : %s ",roleName))




    http.Redirect(w, r, "/users/roles", http.StatusSeeOther)
    return
}


}

