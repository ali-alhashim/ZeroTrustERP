package users

import "zerotrusterp/core"

func init() {
	core.RegisterModels("users", []core.Model{
		{
			Name: "users",
			Schema: `
			CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE,
			password TEXT,
			role VARCHAR(50),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		},
	})
}