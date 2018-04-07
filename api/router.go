package api

import (
	"github.com/gorilla/mux"
)

// GetRouter returns all of the routes in a pointer to a mux.Router object which
// can be passed to ListenAndServe.
//
// Here is where any interceptors/decorators would be applied
func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		// append each route to the router
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.Function)

	}
	return router
}
