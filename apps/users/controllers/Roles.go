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
        "ID":          "r.id",
        "Name":        "r.name",
        "Description": "r.description",
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