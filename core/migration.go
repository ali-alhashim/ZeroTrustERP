package core

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// pluralize converts struct names to table names
func pluralize(name string) string {
	name = strings.ToLower(name)
	if strings.HasSuffix(name, "s") {
		return name + "es"
	}
	return name + "s"
}

// GenerateSQLFromStruct generates SQL for normal fields and join tables
func GenerateSQLFromStruct(model interface{}) (tableSQL []string, joinSQL []string) {
	t := reflect.TypeOf(model)
	tableName := pluralize(t.Name())

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

		// Handle relationships
		if strings.HasPrefix(tag, "many2many:") || strings.HasPrefix(tag, "one2many:") {
			relTable := strings.Split(tag, ":")[1]
			leftTable := tableName
			rightTable := pluralize(field.Type.Elem().Name())

			joinTableName := relTable
			if joinTableName == "" {
				joinTableName = fmt.Sprintf("%s_%s", leftTable, rightTable)
			}

			join := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    %s_id INT REFERENCES %s(id) ON DELETE CASCADE,
    %s_id INT REFERENCES %s(id) ON DELETE CASCADE,
    PRIMARY KEY (%s_id, %s_id)
);`, joinTableName,
				strings.ToLower(t.Name()), leftTable,
				strings.ToLower(field.Type.Elem().Name()), rightTable,
				strings.ToLower(t.Name()), strings.ToLower(field.Type.Elem().Name()))
			joinSQL = append(joinSQL, join)
			continue
		}

		// Handle one2one relationship
		if strings.HasPrefix(tag, "one2one:") {
			refTable := strings.Split(tag, ":")[1]
			if refTable == "" {
				refTable = pluralize(field.Type.Name())
			}

			colName := strings.ToLower(field.Name) + "_id"

			col := fmt.Sprintf("%s INT UNIQUE REFERENCES %s(id) ON DELETE CASCADE",
				colName, refTable)

			cols = append(cols, col)
			continue
		}


		// Normal column
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
			case "default:true":
				colConstraints = append(colConstraints, "DEFAULT TRUE")
			case "default:false":
				colConstraints = append(colConstraints, "DEFAULT FALSE")
			case "default:current_timestamp":
				colConstraints = append(colConstraints, "DEFAULT CURRENT_TIMESTAMP")
			}
		}

		if colType == "" {
			colType = "TEXT"
		}

		cols = append(cols, fmt.Sprintf("%s %s %s", colName, colType, strings.Join(colConstraints, " ")))
	}

	// Build main table SQL
	mainSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", tableName, strings.Join(cols, ",\n\t"))
	tableSQL = append(tableSQL, mainSQL)

	return tableSQL, joinSQL
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

	// Separate normal tables and join tables
	var normalSQL, joinSQL []string
	for _, model := range appModels {
		tables, joins := GenerateSQLFromStruct(model)
		normalSQL = append(normalSQL, tables...)
		joinSQL = append(joinSQL, joins...)
	}

	// Combine SQL: normal tables first, then join tables
	allSQL := append(normalSQL, joinSQL...)

	// Write migration file
	writeMigrationFile(allSQL, appName)

	// Execute SQL
	for _, sql := range allSQL {
		executeSQL(sql)
	}

	log.Println("✅ All migrations complete!")
}

// writeMigrationFile writes SQL to file in app migrations folder
func writeMigrationFile(allSQL []string, app string) {
	dir := fmt.Sprintf("apps/%s/migrations", app)
	os.MkdirAll(dir, os.ModePerm)

	filename := fmt.Sprintf("%s/0001_init.sql", dir)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("❌ Failed to create migration file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(strings.Join(allSQL, "\n\n"))
	if err != nil {
		log.Fatalf("❌ Failed to write migration file: %v", err)
	}
}

// executeSQL runs SQL against PostgreSQL
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
