package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"

	"github.com/askewseth/kubernetes/managers"
	model "github.com/askewseth/kubernetes/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func writeJsonSuccess(i interface{}, w http.ResponseWriter) error {
	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("Unable to marshal to json: %v", err)
	}

	w.WriteHeader(200)
	w.Write(b)

	return nil
}

func writeJsonFail(w http.ResponseWriter, code int, s string) {
	w.WriteHeader(code)
	w.Write([]byte(s))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func libraryHandler(w http.ResponseWriter, r *http.Request) {
	library := managers.GetLibrary()
	err := writeJsonSuccess(library.GetBooks(), w)
	if err != nil {
		log.Error(err)
	}
}

func postBookHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	book := new(model.Book)
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, book)
	if err != nil {
		log.Errorf("%v", err)
		writeJsonFail(w, 400, "error unmarshaling: "+err.Error())
	}

	library := managers.GetLibrary()
	library.AddBook(*book)

	writeJsonSuccess("added", w)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET")

	router.HandleFunc("/library", libraryHandler).Methods("GET")

	router.HandleFunc("/book", postBookHandler).Methods("POST")

	fmt.Println("Listening on http://localhost:5555/")
	err := http.ListenAndServe(":5555", router)
	if err != nil {
		log.Errorf("Error on ListenAndServe: %v", err)
	}
}
