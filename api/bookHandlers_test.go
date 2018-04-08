package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/askewseth/kubernetes/managers"
	model "github.com/askewseth/kubernetes/models"
	uuid "github.com/satori/go.uuid"
)

var (
	server *httptest.Server
)

func init() {
	// initilize the server
	server = httptest.NewServer(GetRouter())
}

// sendRequest creates and sends and http request to the httptest server created in
// the init block, returns the response from that request and an error
func sendRequest(url, method, data string) (*http.Response, error) {
	reader := strings.NewReader(data)

	// make the url with the given server.URL
	url = server.URL + url

	// create the request given the verb and url as well as a blank io.Reader
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Error making request: %v", err)
	}
	defer request.Body.Close()

	// send the request
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Error trying to execute the request: %v", err)
	}

	// return the response
	return res, nil
}

func cleanLibrary() {
	library := managers.GetLibrary()
	library.Books = make(map[uuid.UUID]model.Book)
}

func TestGetBooksAPI(t *testing.T) {
	defer cleanLibrary()

	library := managers.GetLibrary()
	library.Books[uuid.UUID{}] = model.Book{Title: "MyBook"}

	res, err := sendRequest("/books", "GET", "")
	if err != nil {
		t.Errorf("Got error when sending request for GET /books: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected status 200 from GET /books, got %v", res.Status)
	}

	var books []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		t.Errorf("Error trying to read the body from GET /devices: %v", err)
	}

	if len(books) != 1 && books[0]["title"].(string) != "MyBook" {
		t.Errorf("Didn't get the one book in the library back on GET /books")
	}
}

func TestGetBookByID(t *testing.T) {
	defer cleanLibrary()
	cleanLibrary()

	library := managers.GetLibrary()
	library.Books[uuid.UUID{}] = model.Book{Title: "MyBook"}

	res, err := sendRequest("/books/"+uuid.UUID{}.String(), "GET", "")
	if err != nil {
		t.Errorf("Got error when sending request for GET /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected status 200 from GET /book/{id}, got %v", res.Status)
	}

	var book map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&book)
	if err != nil {
		t.Errorf("Error trying to read the body from GET /book/{id}")
	}

	if book["title"].(string) != "MyBook" {
		t.Errorf("Didn't get the correct book on GET /books")
	}
}

func TestGetBookByIDBadBook(t *testing.T) {
	defer cleanLibrary()
	cleanLibrary()

	badID, _ := uuid.NewV4()
	res, err := sendRequest("/books/"+badID.String(), "GET", "")
	if err != nil {
		t.Errorf("Got error when sending request for GET /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 404 {
		t.Errorf("Expected status 404 from GET /book/{id} with bad id, got %v", res.Status)
	}
}

func TestGetBookByIDBadUUID(t *testing.T) {
	defer cleanLibrary()

	res, err := sendRequest("/books/4", "GET", "")
	if err != nil {
		t.Errorf("Got error when sending request for GET /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 400 {
		t.Errorf("Expected status 404 from GET /book/{id} with bad id, got %v", res.Status)
	}
}

func TestDeleteBook(t *testing.T) {
	defer cleanLibrary()
	cleanLibrary()

	library := managers.GetLibrary()
	library.Books[uuid.UUID{}] = model.Book{Title: "MyBook"}

	res, err := sendRequest("/books/"+uuid.UUID{}.String(), "DELETE", "")
	if err != nil {
		t.Errorf("Got error when sending request for DELETE /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 202 {
		t.Errorf("Expected status 202 from DELETE /book/{id}, got %v", res.StatusCode)
		fmt.Println(res.Status)
	}

	if len(library.Books) != 0 {
		fmt.Printf("Had books %+v\n", library.Books)
		t.Errorf("DELETE /books/{id} didn't actually delete book from library")
	}
}

func TestDeleteBookBadBook(t *testing.T) {
	defer cleanLibrary()

	res, err := sendRequest("/books/"+uuid.UUID{}.String(), "DELETE", "")
	if err != nil {
		t.Errorf("Got error when sending request for DELETE /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 404 {
		t.Errorf("Expected status 404 from DELETE /book/{id} with bad id, got %v", res.StatusCode)
		fmt.Println(res.Status)
	}
}

func TestDeleteBookBadUUID(t *testing.T) {
	defer cleanLibrary()

	res, err := sendRequest("/books/4", "DELETE", "")
	if err != nil {
		t.Errorf("Got error when sending request for DELETE /book/{id}: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 400 {
		t.Errorf("Expected status 400 from DELETE /book/{id} with bad uuid, got %v", res.StatusCode)
		fmt.Println(res.Status)
	}
}

func TestPostBook(t *testing.T) {
	defer cleanLibrary()

	library := managers.GetLibrary()

	book := map[string]interface{}{
		"title":  "MyPostBook",
		"rating": 1,
	}
	b, err := json.Marshal(book)
	if err != nil {
		t.Errorf("Error marshing book for PUT /books %v", err)
	}

	res, err := sendRequest("/books", "POST", string(b))
	if err != nil {
		t.Errorf("Got error when sending request for GET /books: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 201 {
		t.Errorf("Expected status 201 from POST /books, got %v", res.Status)
	}

	if len(library.Books) != 1 {
		t.Errorf("Didn't have 1 book in the library after calling POST /books")
	}
}

func TestPutBook(t *testing.T) {
	defer cleanLibrary()

	library := managers.GetLibrary()

	id, _ := uuid.NewV4()
	book := model.Book{Title: "MyPutBook", ID: id}
	library.Books[id] = book

	newBook := map[string]interface{}{
		"title":  "MyNewPutBook",
		"rating": 1,
		"status": 0,
	}
	b, err := json.Marshal(newBook)
	if err != nil {
		t.Errorf("Error marshing book for PUT /books %v", err)
	}

	res, err := sendRequest("/books/"+id.String(), "PUT", string(b))
	if err != nil {
		t.Errorf("Got error when sending request for GET /books: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 202 {
		t.Errorf("Didn't get status accepted on good PUT request, got status %v", res.StatusCode)
	}

	if len(library.Books) != 1 && library.Books[id].Title != newBook["title"] {
		t.Errorf("Didn't PUT /book/{id} correctly")
	}
}

func TestPutBadBook(t *testing.T) {
	defer cleanLibrary()

	newBook := map[string]interface{}{
		"title":  "MyNewPutBook",
		"rating": 1,
		"status": 0,
	}
	b, err := json.Marshal(newBook)
	if err != nil {
		t.Errorf("Error marshing book for PUT /books %v", err)
	}

	id, _ := uuid.NewV4()
	res, err := sendRequest("/books/"+id.String(), "PUT", string(b))
	if err != nil {
		t.Errorf("Got error when sending request for GET /books: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 404 {
		t.Errorf("Didn't get status not found on good PUT request")
	}
}

func TestPutBadUUID(t *testing.T) {
	defer cleanLibrary()

	res, err := sendRequest("/books/4", "PUT", "")
	if err != nil {
		t.Errorf("Got error when sending request for GET /books: %v", err)
		t.FailNow()
	}

	if res.StatusCode != 400 {
		t.Errorf("Didn't get status not found on good PUT request")
	}
}
