package core

import (
	"net/http"
	"reflect"
)


type AppRoute func(*http.ServeMux)


// Store all models
var RegisteredModels []interface{}

// RegisterModel adds a model to migration registry
func RegisterModel(m interface{}) {
    RegisteredModels = append(RegisteredModels, m)
}

// Helper: get table name from struct type
func TableName(m interface{}) string {
    t := reflect.TypeOf(m)
    return t.Name() // default table name = struct name
}




var registeredApps []AppRoute

func Register(app AppRoute) {
	registeredApps = append(registeredApps, app)
}



