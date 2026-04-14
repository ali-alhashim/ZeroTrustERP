package core

import (
	"fmt"
    "zerotrusterp/apps/users/models"
    "strings"
    "log"
    "net"
    "net/http"
)


func GetCountRecords(tableName string) int{

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	

	var count int
    err := DB.QueryRow(query).Scan(&count)
    if err != nil {
        fmt.Printf("Error counting records: %v\n", err)
        return 0
    }

    return count
	 
}



func GetUserByID(id string) models.User {
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

    rows, err := DB.Query(query, id)
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


func GetUserIDByEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
    
	fmt.Printf("Getting user ID for email: %s\n", email)
	var id string

	query := "SELECT id FROM users WHERE email = $1"

	err := DB.QueryRow(query, email).Scan(&id)
	if err != nil {
		log.Println("DB error:", err)
		return "", err
	}

	return id, nil
}


func deleteUserByEmail(email string) {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := DB.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		log.Println("DB error (delete user):", err)
	}
}


func InsertLog(user *models.User, resource string, action string) {

	fmt.Printf("\n Insert Log %d  by %s  for %s with action: %s\n", user.ID, user.Username, resource, action)

	// TODO: implement log insertion to database, create a new log record with user id, resource, action and timestamp
	query := "INSERT INTO logs (user_id, username, email, resource, action) VALUES ($1, $2, $3, $4, $5)"
	_, err := DB.Exec(query, user.ID, user.Username, user.Email, resource, action)
	if err != nil {
		panic(err)
	}
	
}


func GetUserByEmail(email string) *models.User {
    
	fmt.Printf("Getting user by email: %s\n", email)

	email = strings.ToLower(strings.TrimSpace(email))

	

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
        WHERE u.email = $1`

	 rows, err := DB.Query(query, email)
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        return &models.User{}
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
            return &models.User{}
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

    return &user
}


func GetRealIP(r *http.Request) string {
	var ip string

	// 1. Priority: Cloudflare-specific header (Trusted if using CF)
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		ip = cfIP
	}

	// 2. Secondary: Standard Proxy header
	if ip == "" {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			// X-Forwarded-For can be a comma-separated list. 
			// The first one is the original client.
			parts := strings.Split(xForwardedFor, ",")
			ip = strings.TrimSpace(parts[0])
		}
	}

	// 3. Tertiary: Other common proxy headers
	if ip == "" {
		if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
			ip = xRealIP
		}
	}

	// 4. Final Fallback: The direct connection IP
	if ip == "" {
		// net.SplitHostPort correctly handles IPv4 and IPv6 [brackets]
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If it fails (no port), just use the raw string
			ip = r.RemoteAddr
		} else {
			ip = host
		}
	}

	// Clean up any remaining whitespace or weird formatting
	return strings.TrimSpace(ip)
}
