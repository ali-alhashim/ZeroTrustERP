package sequence

import(
	"net/http"
	"zerotrusterp/core"
	"zerotrusterp/apps/sequence/controllers"
)



func SequenceListRoutes(mux *http.ServeMux) {

	mux.Handle("GET /sequence/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListSequence), "prefix_sequences:R"))
	mux.Handle("GET /sequence/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateSequence), "prefix_sequences:W"))
	mux.Handle("POST /sequence/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateSequence), "prefix_sequences:W"))


}