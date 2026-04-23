package main

import (
	"log"
	"os"
	"zerotrusterp/core"

    // Import app packages to register their routes & Models
	_ "zerotrusterp/apps/users"
	_ "zerotrusterp/apps/employees"
	_ "zerotrusterp/apps/dashboard"
	_ "zerotrusterp/apps/system"
)



func main() {
    
	// Load environment variables
	core.LoadEnv(".env")

	db, err := core.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	
    core.DB = db

	if core.DB == nil {
	panic("DB is nil")
     }
	

	// 🔥 CLI migration support
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
        app := ""
        if len(os.Args) > 2 {
            app = os.Args[2]
        }
        core.RunMigrations(app)
        return
    }

	
   
   

	core.StartServer("8000")
}
