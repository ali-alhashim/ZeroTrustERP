package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"zerotrusterp/apps/users/models"
	"zerotrusterp/apps/employees/empapi"
	"zerotrusterp/core"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search    := query.Get("q")
	sortBy    := query.Get("sort")
	order     := query.Get("order")
	page := query.Get("page")
	pageSize   := query.Get("pageSize")

	users := GetUsersFromDB(search, sortBy, order, page, pageSize)
	totalRecords := core.GetCountRecords("users")


	data := map[string]interface{}{
		"Title": "Users",
		"Users": users,
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,

	}

	core.RenderPage(w,r, "apps/users/views/users-list.html", data)
}

func GetUsersFromDB(search, sort, order, page, pageSize string) []models.User {

	query := "SELECT id, username, email, active, online, last_login, related_employee_id FROM users WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// 🔍 SAFE search
	if search != "" {
		query += " AND (username ILIKE $" + strconv.Itoa(argIndex) +
			" OR email ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	// 🔒 Safe sorting
	allowedSort := map[string]string{
		"id":     "id",
		"email":  "email",
		"active": "active",
		"online": "online",
		"last_login":"last_login",
	}

	if col, ok := allowedSort[sort]; ok {
		query += " ORDER BY " + col
		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	// 📄 Pagination (page + pageSize)
	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	// defaults
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 10
	}

	offset := (p - 1) * ps

	query += " LIMIT $" + strconv.Itoa(argIndex) +
		" OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, ps, offset)

	// ✅ Execute
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []models.User
	



	for rows.Next() {
		var u models.User
		var relatedEmployeeID *int
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Active, &u.Online, &u.LastLogin, &relatedEmployeeID)
		if err != nil {
			panic(err)
		}
		

		if relatedEmployeeID != nil {
			employee := empapi.GetEmployeeByID(*relatedEmployeeID)
			u.RelatedEmployee = &employee // Use the ampersand to get the pointer
      }


		users = append(users, u)
	}

	return users
}





func SetUserActive(w http.ResponseWriter, r *http.Request){

	userID := r.PathValue("id")
	fmt.Print("set user active with ID:", userID)
    
	query :="update users set active=true where id = $1"

	core.DB.Exec(query, userID)


	http.Redirect(w, r, "/users/list", http.StatusSeeOther)


	 CurrentUser := core.GetCurrentUser(r)


		InsertLog(CurrentUser, "users", fmt.Sprintf("Set User Active for User ID : %s ",userID))
	

}

func SetUserInactive(w http.ResponseWriter, r *http.Request){

	userID := r.PathValue("id")
	fmt.Print("set user inactive ID:", userID)
    
	query :="update users set active=false where id = $1"

	core.DB.Exec(query, userID)


	 CurrentUser := core.GetCurrentUser(r)


		InsertLog(CurrentUser, "users", fmt.Sprintf("Set User Inactive for User ID : %s ",userID))


	http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		
}



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

		CurrentUser := core.GetCurrentUser(r)


		InsertLog(CurrentUser, "Users", fmt.Sprintf("Created user %s with email: %s",username, email))



		http.Redirect(w, r, "/users/list", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}


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
