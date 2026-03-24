package core

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

func GenerateSQLFromStruct(model interface{}) string {
	t := reflect.TypeOf(model)
	tableName := strings.ToLower(t.Name()) // use struct name as table name

	var cols []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip virtual fields
		if strings.Contains(field.Tag.Get("v"), "true") {
			continue
		}

		tag := field.Tag.Get("f")
		if tag == "" {
			continue
		}

		colName := strings.ToLower(field.Name)
		colType := ""
		colConstraints := []string{}

		parts := strings.Split(tag, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)

			switch p {
			case "number":
				colType = "INTEGER"
			case "text":
				colType = "TEXT"
			case "bool":
				colType = "BOOLEAN"
			case "timestamp":
				colType = "TIMESTAMP"
			case "primary":
				colConstraints = append(colConstraints, "PRIMARY KEY")
			case "auto":
				if colType == "INTEGER" {
					colType = "SERIAL"
				}
			case "unique":
				colConstraints = append(colConstraints, "UNIQUE")
			case "notnull":
				colConstraints = append(colConstraints, "NOT NULL")
			default:
				// ignore unknown for now
			}
		}

		if colType == "" {
			colType = "TEXT"
		}

		cols = append(cols, fmt.Sprintf("%s %s %s", colName, colType, strings.Join(colConstraints, " ")))
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", tableName, strings.Join(cols, ",\n\t"))
	return sql
}

// RunMigrations migrates all registered models, or only a specific app
func RunMigrations(appName string) {
	log.Println("🔹 Starting migrations...")

	if len(RegisteredModels) == 0 {
		log.Println("⚠️ No models registered!")
		return
	}

	// collect all models for this app
	var appModels []interface{}
	for _, model := range RegisteredModels {
		modelApp := getAppNameFromModel(model)
		if appName != "" && modelApp != appName {
			continue
		}
		appModels = append(appModels, model)
	}

	if len(appModels) == 0 {
		log.Println("⚠️ No models found for app:", appName)
		return
	}

	// write all tables in one file
	writeMigrationFile(appModels, appName)

	// execute SQL for all models
	for _, model := range appModels {
		sql := GenerateSQLFromStruct(model)
		executeSQL(sql)
		log.Printf("✅ Migrated %s", reflect.TypeOf(model).Name())
	}

	log.Println("✅ All migrations complete!")
}


// getAppNameFromModel extracts app name from struct package path
func getAppNameFromModel(m interface{}) string {
	t := reflect.TypeOf(m)
	pkgPath := t.PkgPath() // e.g., "zerotrusterp/apps/users/models"

	parts := strings.Split(pkgPath, "/")
	for i, p := range parts {
		if p == "apps" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "core"
}

// writeMigrationFile writes SQL to file in app migrations folder
func writeMigrationFile(models []interface{}, app string) {
	dir := fmt.Sprintf("apps/%s/migrations", app)
	os.MkdirAll(dir, os.ModePerm)

	filename := fmt.Sprintf("%s/0001_init.sql", dir)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("❌ Failed to create migration file: %v", err)
	}
	defer f.Close()

	var allSQL []string
	for _, model := range models {
		sql := GenerateSQLFromStruct(model)
		allSQL = append(allSQL, sql)
	}

	_, err = f.WriteString(strings.Join(allSQL, "\n\n"))
	if err != nil {
		log.Fatalf("❌ Failed to write migration file: %v", err)
	}
}



func executeSQL(sql string) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("❌ Failed to execute SQL: %v\nSQL: %s", err, sql)
	}
}
