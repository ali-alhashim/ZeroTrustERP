package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/users/models"
	"strconv"
	"encoding/json"
)


func ListRoles(w http.ResponseWriter, r *http.Request) {
	fmt.Print("List Roles")



	query := r.URL.Query()

	search    := query.Get("q")
	sortBy    := query.Get("sort")
	order     := query.Get("order")
	page := query.Get("page")
	pageSize   := query.Get("pageSize")

	roles := GetRolesFromDB(search, sortBy, order, page, pageSize)

	data := map[string]interface{}{
		"Title": "Roles",
		"Roles": roles,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,

	}

	core.RenderPage(w,r, "apps/users/views/roles-list.html", data)

}

func GetRolesFromDB(search, sort, order, page, pageSize string) []models.Role {
    // 1. Build the base query parts
    selectBase := `
        SELECT 
            r.id, r.name, r.description,r.created_at, r.updated_at,
            COALESCE(json_agg(
                json_build_object(
                    'id', p.id,
                    'resource', p.resource,
                    'action', p.action
                )
            ) FILTER (WHERE p.id IS NOT NULL), '[]') AS permissions
        FROM roles r
        LEFT JOIN roles_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id `

    whereClause := ""
    args := []interface{}{}
    argIndex := 1

    // 2. Add Search (Must use WHERE before GROUP BY)
    if search != "" {
        whereClause = fmt.Sprintf(" WHERE (r.name ILIKE $%d OR r.description ILIKE $%d) ", argIndex, argIndex+1)
        args = append(args, "%"+search+"%", "%"+search+"%")
        argIndex += 2
    }

    // 3. Assemble the query with GROUP BY
    fullQuery := selectBase + whereClause + " GROUP BY r.id "

    // 4. Add Sorting
    allowedSort := map[string]string{
        "id":          "r.id",
        "name":        "r.name",
        "description": "r.description",
    }
    if col, ok := allowedSort[sort]; ok {
        fullQuery += " ORDER BY " + col
        if order == "desc" {
            fullQuery += " DESC "
        } else {
            fullQuery += " ASC "
        }
    }

    // 5. Pagination
    p, _ := strconv.Atoi(page)
    ps, _ := strconv.Atoi(pageSize)
    if p <= 0 { p = 1 }
    if ps <= 0 || ps > 100 { ps = 10 }
    offset := (p - 1) * ps

    fullQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, ps, offset)

    // 6. Execute and Scan
    rows, err := core.DB.Query(fullQuery, args...)
    if err != nil {
        fmt.Println("Query Error:", err)
        return nil
    }
    defer rows.Close()

    var roles []models.Role
    for rows.Next() {
        var r models.Role
        var permissionsJSON []byte // Scan JSON into bytes first

        // Ensure models.Role.Permissions is a slice of Permission structs
        err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt, &permissionsJSON)
        if err != nil {
            fmt.Println("Scan Error:", err)
            continue
        }

        // Unmarshal the JSON bytes into the struct slice
        if err := json.Unmarshal(permissionsJSON, &r.Permissions); err != nil {
            fmt.Println("Unmarshal Error:", err)
        }
        
        roles = append(roles, r)
    }

    return roles
}

// Get Roles ID & Names as JSON for dropdowns and APIs
// GetRolesAsJson returns the roles as a JSON byte slice
func GetRolesAsJson() ([]byte, error) {
    rows, err := core.DB.Query("SELECT id, name FROM roles ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Define a local struct for clean JSON mapping
    type role struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }

    var roles []role
    for rows.Next() {
        var r role
        if err := rows.Scan(&r.ID, &r.Name); err != nil {
            return nil, err
        }
        roles = append(roles, r)
    }

    // Check if the loop finished correctly or hit a connection error
    if err = rows.Err(); err != nil {
        return nil, err
    }

    // Convert the slice to JSON
    return json.Marshal(roles)
}

func FetchRolesAPI(w http.ResponseWriter, r *http.Request) {
    fmt.Print("Fetch Roles API")


    if r.Method != http.MethodGet {
        fmt.Println("Invalid method for FetchRolesAPI:", r.Method)
         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    fmt.Print("FetchRolesAPI GET Request")
    jsonData, err := GetRolesAsJson()
    if err != nil {
        http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)

}


func DeleteRoleFromUser(w http.ResponseWriter, r *http.Request){

    userID := r.PathValue("userID")
    roleID := r.PathValue("roleID")


    fmt.Print("delete Role API userID=", userID, " roleID= ", roleID)


    // DELETE FROM users_roles 
    // WHERE user_id = $1 AND role_id = $2;

    if userID == "" || roleID == "" {
        http.Error(w, "Missing userID or roleID", http.StatusBadRequest)
        return
    }

    sqlStatement := `DELETE FROM users_roles WHERE user_id = $1 AND role_id = $2`
    result, err := core.DB.Exec(sqlStatement, userID, roleID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        fmt.Println("No record found to delete.")
    }

     CurrentUser := core.GetCurrentUser(r)

    InsertLog(CurrentUser, "Delete Role From User", fmt.Sprintf("User ID : %s , Role ID: ",roleID))

    // 5. Respond to frontend
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "success", "message": "Role deleted"}`)

}




func RoleDeatils(w http.ResponseWriter, r *http.Request){
    roleID := r.PathValue("roleID")

    fmt.Print("...Role Details ID:", roleID)

    Role := GetRoleByID(roleID)

    data := map[string]interface{}{
		"Title": "Users",
		"Role": Role,
        "Resources": core.GetAllResources(),
	}


    core.RenderPage(w,r, "apps/users/views/roles-details.html", data)
}


func GetRoleByID(roleID string) models.Role{
    
    var role models.Role
    var permissions []models.Permission
     
    fmt.Print("\n Get Role By ID ", roleID)

    // select role with related permissions 
    // we have roles table with id, name, description
    // and we have permissions table with id, resource, action ,description
    // roles_permissions table with role_id, permission_id
    // Role struct has Permissions *[]Permission

    query := `
        SELECT 
            r.id, r.name, r.description,
            p.id, p.resource, p.action, p.description
        FROM roles r
        LEFT JOIN roles_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
        WHERE r.id = $1`

    rows, err := core.DB.Query(query, roleID)
    if err != nil {
        fmt.Print(err)
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Permission
        // We use pointers/null types if permissions might be empty (LEFT JOIN)
        var pID *int 
        var pResource, pAction, pDesc *string

        err := rows.Scan(
            &role.ID, &role.Name, &role.Description,
            &pID, &pResource, &pAction, &pDesc,
        )
        if err != nil {
            fmt.Print(err)
        }
        
        if pID != nil {
            p.ID = *pID
            p.Resource = *pResource
            p.Action = *pAction
            p.Description = *pDesc
            permissions = append(permissions, p)
        }
    }

    role.Permissions = &permissions

    fmt.Print(" \n selected role name is: ", role.Name , "\n")

    return role

}



func DeletePermissionFromRole(w http.ResponseWriter, r *http.Request){

     roleID       := r.PathValue("roleID")
     permissionID := r.PathValue("permissionID")

     fmt.Print("Delete Permission " +permissionID +" From Role ID: ", roleID)

      if permissionID == "" || roleID == "" {
        http.Error(w, "Missing permissionID or roleID", http.StatusBadRequest)
        return
    }


     sqlStatement := `DELETE FROM roles_permissions WHERE role_id = $1 AND permission_id = $2`
    result, err := core.DB.Exec(sqlStatement, roleID, permissionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        fmt.Println("No record found to delete.")
    }


    CurrentUser := core.GetCurrentUser(r)

    InsertLog(CurrentUser, "Delete Permission From Role", fmt.Sprintf("Role ID : %s ,Permission ID %s: ",roleID, permissionID))

    // 5. Respond to frontend
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "success", "message": "Role deleted"}`)


}


func UpdateRole(w http.ResponseWriter, r *http.Request){
    roleID:= r.PathValue("roleID")

    fmt.Print("...........update Role ID:", roleID)
}