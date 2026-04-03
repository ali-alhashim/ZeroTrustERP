package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"zerotrusterp/core"
)


func OnlineUsersAPI(w http.ResponseWriter, r *http.Request) {
	

		fmt.Printf("Received request for online users\n")


		//first check if the user is authenticated
		//if not authenticated, return empty list or 401
		// get cookies {session, email} from request , if no cookies, return 401 

	   	 cookies, err := r.Cookie("session")
		 if err != nil {
			fmt.Printf("No session cookie found: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		 }
		 emailCookie, err := r.Cookie("email")
		 if err != nil {
			fmt.Printf("No email cookie found: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		 }	

		 fmt.Printf("Session cookie: %s, Email cookie: %s\n", cookies.Value, emailCookie.Value)

		 // Here you would typically validate the session cookie against your session store
		

		


		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		onlineUsers := getOnlineUsers()
		fmt.Printf("Online users: %v\n", onlineUsers)
		w.Write([]byte(fmt.Sprintf(`["%s"]`, strings.Join(onlineUsers, `","`))))




	
}



func getOnlineUsers() []string {

	//select id from users where online = true
	query := "SELECT id FROM users WHERE online = true"

	rows, err := core.DB.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	var onlineUsers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			panic(err)
		}
		onlineUsers = append(onlineUsers, id)
	}

	return onlineUsers	
}



// this function will be called when we recived /api/heartbeat request, it will update the user's online status in the database
func SetUserOnline(w http.ResponseWriter, r *http.Request) {

	// get cookies {session, email} from request , if no cookies, return 401
	cookies, err := r.Cookie("session")
	if err != nil {
		fmt.Printf("No session cookie found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	emailCookie, err := r.Cookie("email")
	if err != nil {
		fmt.Printf("No email cookie found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("Session cookie: %s, Email cookie: %s\n", cookies.Value, emailCookie.Value)

	// Here you would typically validate the session cookie against your session store
	// If valid, update the user's online status in the database
	// For example: UPDATE users SET online = true WHERE email = emailCookie.Value
	if core.IsValidSessionToken(emailCookie.Value, cookies.Value) {
		fmt.Printf("Valid session for user: %s\n", emailCookie.Value)
		// Update user's online status in the database
		core.SetUserIsOnline(emailCookie.Value, true) 
		
	} else {
		fmt.Printf("Invalid session for user: %s\n", emailCookie.Value)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

}


func StopUserHeartbeat(w http.ResponseWriter, r *http.Request) {

	// get cookies {session, email} from request , if no cookies, return 401
	cookies, err := r.Cookie("session")
	if err != nil {
		fmt.Printf("No session cookie found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	emailCookie, err := r.Cookie("email")
	if err != nil {
		fmt.Printf("No email cookie found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("Session cookie: %s, Email cookie: %s\n", cookies.Value, emailCookie.Value)

	// Here you would typically validate the session cookie against your session store
	// If valid, update the user's online status in the database
	// For example: UPDATE users SET online = false WHERE email = emailCookie.Value
	if core.IsValidSessionToken(emailCookie.Value, cookies.Value) {
		fmt.Printf("Valid session for user: %s\n", emailCookie.Value)
		// Update user's online status in the database
		core.SetUserIsOnline(emailCookie.Value, false)
		
	} else {
		fmt.Printf("Invalid session for user: %s\n", emailCookie.Value)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

}