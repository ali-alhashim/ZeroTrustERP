package controllers

import (
	"net/http"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
    
	w.Header().Set("Content-Type", "text/html")
	//pass the list of users to the template here
	http.ServeFile(w, r, "./apps/users/views/list.html")


	
	
}