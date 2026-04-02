package controllers

import (
	"fmt"
	"net/http"
	"strings"
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
	
	return []string{"1", "2", "3"}	
}