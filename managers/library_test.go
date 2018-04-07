package managers

import (
	"testing"

	"github.com/satori/go.uuid"

	model "github.com/askewseth/kubernetes/models"
)

func TestAddBook(t *testing.T) {
	library := newLibrary()

	// verify library is empty
	if len(library.Books) != 0 {
		t.Errorf("Library wasn't empty to begin with")
		t.FailNow()
	}

	book := model.NewBook()
	book.Title = "MyBook"

	library.AddBook(book)

	if library.Books[book.ID].Title != "MyBook" {
		t.Errorf("Added in a book using AddBook but the book wasn't found after the fact")
	}
}

func TestSortBooks(t *testing.T) {
	books := []model.Book{
		model.Book{Title: "B"},
		model.Book{Title: "A"},
	}

	sortBooks(books)

	if books[0].Title != "A" {
		t.Error("")
	}
}

func TestGetBook(t *testing.T) {
	library := newLibrary()

	book := model.NewBook()
	book.Title = "MyBook"

	library.AddBook(book)

	// verify the book is there
	if len(library.Books) != 1 {
		t.Error("Didn't have 1 book in the library after adding 1 book")
		t.FailNow()
	}

	books := library.GetBooks()
	if len(books) != 1 && books[0].Title != book.Title {
		t.Error("GetBooks didn't return the book that was added")
	}
}

func TestModify(t *testing.T) {
	library := newLibrary()

	// create and add a known book
	book := model.NewBook()
	book.Title = "MyBook"
	book.Author = "me"

	library.AddBook(book)

	// verify the book was added
	if len(library.Books) != 1 {
		t.Errorf("Didn't have 1 book in the library after adding 1 book")
		t.FailNow()
	}

	// try to modify 1 of the 2 fields
	modBook := model.NewDefaultBook()
	modBook.ID = book.ID
	modBook.Title = "MyNewBook"
	err := library.ModifyBook(modBook)
	if err != nil {
		t.Errorf("Error modifing book: %v", err)
		t.FailNow()
	}

	newBook, err := library.GetBookByID(book.ID)
	if err != nil {
		t.Errorf("Error getting book by ID: %v", err)
		t.FailNow()
	}

	if newBook.Title != modBook.Title {
		t.Errorf("ModifyBook failed to modify a given attribute")
	}

	if newBook.Author != book.Author {
		t.Errorf("ModifyBook modified a field that wasn't given")
	}
}

func TestDelete(t *testing.T) {
	library := newLibrary()

	// create and add a known book
	book := model.NewBook()
	book.Title = "MyBook"

	library.AddBook(book)

	// verify the book was added
	if len(library.Books) != 1 {
		t.Errorf("Didn't have 1 book in the library after adding 1 book")
		t.FailNow()
	}

	// try to delete the book
	err := library.DeleteBook(book.ID)
	if err != nil {
		t.Errorf("Got an error while trying to delete a book: %v", err)
	}

	if len(library.Books) != 0 {
		t.Errorf("Didn't correctly delete the book from the library, still had 1 book after delete")
	}

	// try to delete a non-existing book and make sure it errors
	bogusID, _ := uuid.NewV4()
	err = library.DeleteBook(bogusID)
	if err != ErrNoBookWithThatID {
		t.Errorf("Expected to get %v error when calling DeleteBook with bogus ID but got %v", ErrNoBookWithThatID, err)
	}
}
