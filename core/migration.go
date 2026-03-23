package core

import (
	"database/sql"
	"fmt"
)

type Model struct {
	Name   string
	Schema string
}

var registeredModels = map[string][]Model{}

// Register models per app
func RegisterModels(app string, models []Model) {
	registeredModels[app] = models
}

// Run migrations
func Migrate(db *sql.DB, app string) {
	if app == "all" {
		for appName, models := range registeredModels {
			fmt.Println("Migrating:", appName)
			runModels(db, models)
		}
		return
	}

	models, ok := registeredModels[app]
	if !ok {
		fmt.Println("No models found for app:", app)
		return
	}

	runModels(db, models)
}

func runModels(db *sql.DB, models []Model) {
	for _, m := range models {
		fmt.Println("Creating table:", m.Name)

		_, err := db.Exec(m.Schema)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
