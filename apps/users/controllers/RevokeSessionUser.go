package controllers

import(
	"fmt"
	"net/http"
	"zerotrusterp/core"
	
)

func RevokeSession(w http.ResponseWriter, r *http.Request){

	//post request /users/RevokeSession/{{.userID}}
    //get user ID from the link
	//  RevokeSession by update users table set session_token to empty & otp_hash empty & online = false where id = userID
    if r.Method !=http.MethodPost{
       fmt.Print("only http POST Allowed")
	   http.Error(w, "Failed only POST Request Allowed", http.StatusBadRequest)
       return
	}

	userID := r.PathValue("id")
    
	

	fmt.Printf("Processing RevokeSession for User ID: %s\n", userID)

    sql := "update users set session_token='', otp_hash='', online = false where id = $1"
	_, err := core.DB.Exec(sql, userID)

	 if err != nil {
        fmt.Printf("Failed to RevokeSession : %v\n", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }


	w.WriteHeader(http.StatusOK)
    w.Write([]byte("Success"))



	  CurrentUser := core.GetCurrentUser(r)


		InsertLog(CurrentUser, "users", fmt.Sprintf("RevokeSession for User ID : %s ",userID))


}


