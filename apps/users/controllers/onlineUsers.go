package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "zerotrusterp/core"
)

// ✅ GET /api/online-users
// Returns real-time online user IDs from Hub (not DB)

func OnlineUsersAPI(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Received request for online users")

    // ✅ Authentication check (session + email cookies)
    sessionCookie, err := r.Cookie("session")
    if err != nil {
        fmt.Println("No session cookie found:", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    emailCookie, err := r.Cookie("email")
    if err != nil {
        fmt.Println("No email cookie found:", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    fmt.Printf("Session cookie: %s, Email cookie: %s\n",
        sessionCookie.Value, emailCookie.Value)

    // ✅ Validate the session token against your auth system
    if !core.IsValidSessionToken(emailCookie.Value, sessionCookie.Value) {
        fmt.Println("Invalid session for user:", emailCookie.Value)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // ✅ Get real-time online users directly from Hub
	var MainHub = core.GetMainHub() // Access the global Hub instance
    onlineUsers := MainHub.GetOnlineUserIDs()

    fmt.Println("Online users:", onlineUsers)

    // ✅ Respond with JSON list of online user IDs
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(onlineUsers)
}