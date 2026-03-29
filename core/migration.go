package core

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
)

// toSnakeCase converts CamelCase to snake_case.
// e.g. JobTitle -> job_title, UserRole -> user_role
func toSnakeCase(s string) string {
    var b strings.Builder
    runes := []rune(s)
    for i, r := range runes {
        if i > 0 && unicode.IsUpper(r) {
            // Don't insert underscore if previous char was also uppercase
            // AND next char is uppercase or end of string (e.g. "ID", "URL")
            prevUpper := unicode.IsUpper(runes[i-1])
            nextUpper := i+1 >= len(runes) || unicode.IsUpper(runes[i+1])
            if !(prevUpper && nextUpper) {
                b.WriteByte('_')
            }
        }
        b.WriteRune(unicode.ToLower(r))
    }
    return b.String()
}

// pluralize converts a CamelCase struct name to a snake_case plural table name.
// e.g. JobTitle -> job_titles, Employee -> employees, Status -> statuses
func pluralize(name string) string {
	snake := toSnakeCase(name)
	if strings.HasSuffix(snake, "s") {
		return snake + "es"
	}
	return snake + "s"
}

// derefType dereferences pointer and slice types to get the underlying struct type.
// Handles: *[]Role -> Role, []Role -> Role, *Role -> Role
func derefType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// GenerateSQLFromStruct generates:
//   - tableSQL:    CREATE TABLE statements (no inline FK constraints)
//   - joinSQL:     CREATE TABLE statements for join tables (no inline FK constraints)
//   - deferredSQL: ALTER TABLE ADD CONSTRAINT FK statements (run after all tables exist)
func GenerateSQLFromStruct(model interface{}) (tableSQL []string, joinSQL []string, deferredSQL []string) {
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

		// Handle many2many / one2many relationships (join tables)
		if strings.HasPrefix(tag, "many2many:") || strings.HasPrefix(tag, "one2many:") {
			relTable := strings.Split(tag, ":")[1]
			leftTable := tableName

			elemType := derefType(field.Type)
			if elemType.Name() == "" {
				log.Printf("⚠️  Could not resolve element type for field '%s' in struct '%s' — skipping join table", field.Name, t.Name())
				continue
			}

			rightTable := pluralize(elemType.Name())
			leftName := toSnakeCase(t.Name())
			rightName := toSnakeCase(elemType.Name())

			joinTableName := relTable
			if joinTableName == "" {
				joinTableName = fmt.Sprintf("%s_%s", leftTable, rightTable)
			}

			// Join table with plain INT columns — NO inline REFERENCES
			join := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n    %s_id INT,\n    %s_id INT,\n    PRIMARY KEY (%s_id, %s_id)\n);",
				joinTableName, leftName, rightName, leftName, rightName)
			joinSQL = append(joinSQL, join)

			// Deferred FKs for both sides
			deferredSQL = append(deferredSQL, fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s_id FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE;",
				joinTableName, joinTableName, leftName, leftName, leftTable,
			))
			deferredSQL = append(deferredSQL, fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s_id FOREIGN KEY (%s_id) REFERENCES %s(id) ON DELETE CASCADE;",
				joinTableName, joinTableName, rightName, rightName, rightTable,
			))
			continue
		}

		// Handle many2one relationship
		if strings.HasPrefix(tag, "many2one:") {
			refTable := strings.Split(tag, ":")[1]
			if refTable == "" {
				refTable = pluralize(derefType(field.Type).Name())
			}
			colName := toSnakeCase(field.Name) + "_id"
			cols = append(cols, fmt.Sprintf("%s INT", colName))
			deferredSQL = append(deferredSQL, fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(id) ON DELETE CASCADE;",
				tableName, tableName, colName, colName, refTable,
			))
			continue
		}

		// Handle one2one relationship
		if strings.HasPrefix(tag, "one2one:") {
			refTable := strings.Split(tag, ":")[1]
			if refTable == "" {
				refTable = pluralize(derefType(field.Type).Name())
			}
			colName := toSnakeCase(field.Name) + "_id"
			cols = append(cols, fmt.Sprintf("%s INT UNIQUE", colName))
			deferredSQL = append(deferredSQL, fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(id) ON DELETE CASCADE;",
				tableName, tableName, colName, colName, refTable,
			))
			continue
		}

		// Normal column
		colName := toSnakeCase(field.Name)
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

	mainSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n);", tableName, strings.Join(cols, ",\n\t"))
	tableSQL = append(tableSQL, mainSQL)

	return tableSQL, joinSQL, deferredSQL
}

// RunMigrations migrates all registered models, or only a specific app.
// Opens a single DB connection for all statements.
func RunMigrations(appName string) {
	log.Println("🔹 Starting migrations...")

	if len(RegisteredModels) == 0 {
		log.Println("⚠️ No models registered!")
		return
	}

	// Collect all models for this app
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

	// Sort models by dependencies
	appModels = sortModelsByDependencies(appModels)

	// Collect all SQL in three phases
	var normalSQL, joinSQL, deferredSQL []string
	for _, model := range appModels {
		tables, joins, deferred := GenerateSQLFromStruct(model)
		normalSQL = append(normalSQL, tables...)
		joinSQL = append(joinSQL, joins...)
		deferredSQL = append(deferredSQL, deferred...)
	}

	// Execution order:
	//   1. Normal tables   (plain columns, no FK constraints)
	//   2. Join tables     (plain columns, no FK constraints)
	//   3. Deferred FKs    (ALTER TABLE, after every table exists)
	allSQL := append(append(normalSQL, joinSQL...), deferredSQL...)

	// Write migration file
	writeMigrationFile(allSQL, appName)

	// Open ONE connection for all statements
	db, err := InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	for _, sql := range allSQL {
		log.Printf("⚙️  Executing: %s", firstLine(sql))
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatalf("❌ Failed to execute SQL: %v\nSQL: %s", err, sql)
		}
	}

	log.Println("✅ All migrations complete!")
}

// firstLine returns the first line of a string for concise logging
func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx != -1 {
		return s[:idx]
	}
	return s
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

// getAppNameFromModel extracts app name from struct package path
func getAppNameFromModel(m interface{}) string {
	t := reflect.TypeOf(m)
	pkgPath := t.PkgPath()

	parts := strings.Split(pkgPath, "/")
	for i, p := range parts {
		if p == "apps" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "core"
}

// sortModelsByDependencies sorts models so referenced tables are created first.
// Since ALL FKs are now deferred via ALTER TABLE, table creation order does not
// actually matter for correctness — but we still sort for cleaner output.
// Cycles are broken by simply not following a back-edge; the node is always
// guaranteed to appear in the output via the safety-net pass at the end.
func sortModelsByDependencies(models []interface{}) []interface{} {
	// 1. Build lookup: TableName -> Model
	modelMap := make(map[string]interface{})
	var modelNames []string

	for _, m := range models {
		t := reflect.TypeOf(m)
		tableName := pluralize(t.Name())
		modelMap[tableName] = m
		modelNames = append(modelNames, tableName)
	}

	// 2. Build dependency graph
	dependencies := make(map[string][]string)
	for _, m := range models {
		t := reflect.TypeOf(m)
		tableName := pluralize(t.Name())
		var deps []string

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("f")

			if strings.HasPrefix(tag, "one2one:") || strings.HasPrefix(tag, "many2one:") {
				parts := strings.Split(tag, ":")
				refTable := parts[1]
				if refTable == "" {
					refTable = pluralize(derefType(field.Type).Name())
				}
				if _, exists := modelMap[refTable]; exists {
					deps = append(deps, refTable)
				}
			}
		}
		dependencies[tableName] = deps
	}

	// 3. Topological sort (DFS) with cycle breaking.
	//    When a back-edge is found (temp[node] == true) we skip that edge
	//    but do NOT return — the outer loop ensures every node is visited.
	sorted := []string{}
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(string)
	visit = func(node string) {
		if visited[node] {
			return
		}
		if temp[node] {
			// Back-edge: skip to break the cycle
			log.Printf("⚠️  Circular dependency on table '%s' — FK will be applied via ALTER TABLE.", node)
			return
		}

		temp[node] = true
		for _, dep := range dependencies[node] {
			visit(dep)
		}
		temp[node] = false
		visited[node] = true
		sorted = append(sorted, node)
	}

	for _, name := range modelNames {
		visit(name)
	}

	// 4. Build sorted model slice
	var sortedModels []interface{}
	inSorted := make(map[string]bool)
	for _, name := range sorted {
		inSorted[name] = true
		if m, ok := modelMap[name]; ok {
			sortedModels = append(sortedModels, m)
		}
	}

	// 5. Safety net: any model skipped due to a cycle gets appended at the end
	for _, name := range modelNames {
		if !inSorted[name] {
			log.Printf("⚠️  Table '%s' not reached during sort — appending at end.", name)
			if m, ok := modelMap[name]; ok {
				sortedModels = append(sortedModels, m)
			}
		}
	}

	return sortedModels
}