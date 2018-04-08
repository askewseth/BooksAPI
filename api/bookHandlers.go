package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/askewseth/kubernetes/managers"
	model "github.com/askewseth/kubernetes/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrInvalidUUID is the error returned whenever a user gives an id on a
	// GET, PUT, or DELTE command that isn't a valid UUID
	ErrInvalidUUID = errors.New("The given id was not a valid UUID")
)

// GetBooks is the handler for the GET /books api call,
// it just returns a list of all of the books in the library
func GetBooks(w http.ResponseWriter, r *http.Request) {
	library := managers.GetLibrary()
	err := writeJSONSuccess(w, library.GetBooks(), http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

// PostBook is the handler for the POST /books api call,
// it will add a new book to the library
func PostBook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	book := model.NewBook()
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Errorf("%v", err)
		writeJSONFail(w, 400, "The Post Body was invalid")
		return
	}
	defer r.Body.Close()

	// validate that the books attributes are in the appropriate bounds
	err = book.Validate()
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, err.Error())
		return
	}

	library := managers.GetLibrary()
	library.AddBook(book)

	writeJSONSuccess(w, "", http.StatusCreated)
}

// PutBook is the handler for the PUT /books/{id} api call,
// it will modify the given fields in the API
func PutBook(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	id, err := uuid.FromString(parameters["id"])
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, ErrInvalidUUID.Error())
		return
	}

	book := model.NewDefaultBook()
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// validate that the books attributes are in the appropriate bounds
	err = book.Validate()
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, err.Error())
		return
	}

	book.ID = id

	library := managers.GetLibrary()
	err = library.ModifyBook(book)
	if err != nil {
		writeJSONFail(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONSuccess(w, "", http.StatusAccepted)
}

// DeleteBook is the handler for the DELETE /books/{id} call
// it will remove a book from the library
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	id, err := uuid.FromString(parameters["id"])
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, ErrInvalidUUID.Error())
		return
	}

	library := managers.GetLibrary()
	err = library.DeleteBook(id)
	if err != nil {
		writeJSONFail(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONSuccess(w, "", http.StatusAccepted)
}

// GetBookByID is the handler for the GET /books/{id} call
// it will return a specific book from the library given it's uuid
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	// parse the uuid
	parameters := mux.Vars(r)
	id, err := uuid.FromString(parameters["id"])
	if err != nil {
		writeJSONFail(w, http.StatusBadRequest, ErrInvalidUUID.Error())
		return
	}

	// try to find the book in the library
	library := managers.GetLibrary()
	book, err := library.GetBookByID(id)
	if err != nil {
		writeJSONFail(w, http.StatusNotFound, err.Error())
		return
	}

	// marshal and return the book
	writeJSONSuccess(w, book, http.StatusOK)
}
