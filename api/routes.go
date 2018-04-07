package api

import "net/http"

type route struct {
	Pattern     string
	Function    http.HandlerFunc
	Method      string
	Description string
}

var routes = []route{
	route{
		Pattern:     "/books",
		Function:    GetBooks,
		Method:      "GET",
		Description: "/books will print out all of the books",
	},

	route{
		Pattern:     "/books",
		Function:    PostBook,
		Method:      "POST",
		Description: "POST /book will create a new book in the library",
	},

	route{
		Pattern:     "/books/{id}",
		Function:    GetBookByID,
		Method:      "GET",
		Description: "/book/{id} will return a specific book by it's id",
	},

	route{
		Pattern:     "/books/{id}",
		Function:    PutBook,
		Method:      "PUT",
		Description: "PUT /book/{id} will modify the given book if it exists",
	},

	route{
		Pattern:     "/books/{id}",
		Function:    DeleteBook,
		Method:      "DELETE",
		Description: "DELETE /book/{id} will remove the given book if it exists",
	},
}
