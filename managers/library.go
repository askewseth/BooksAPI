package managers

import (
	"errors"
	"sort"
	"sync"

	"github.com/askewseth/kubernetes/models"
	uuid "github.com/satori/go.uuid"
)

var (
	// The global libary instance
	library Library

	// sync.Once for the singleton
	once = sync.Once{}

	// ErrNoBookWithThatID is the error returned whenever someone tried to
	// GET, PUT, or DELETE a book with an id that isn't found in the manager
	ErrNoBookWithThatID = errors.New("The given book uuid wasn't found")
)

// Library is the struct that holds all of the books
type Library struct {
	sync.Mutex `json:"-"`
	Books      map[uuid.UUID]model.Book `json:"books"`
}

// GetLibrary is a thread safe singleton model which will, on the first
// time being called, initalize a new libary, and on subsequent calls
// return that same instance of the library struct
func GetLibrary() Library {

	// if this is the first time this function has been called
	// then create a new libary
	once.Do(func() {
		library = newLibrary()
	})

	// return the global library instance
	return library
}

// newLibrary will return a newly initalized library struct, it is only
// to be used by GetLibrary
func newLibrary() Library {
	return Library{Books: make(map[uuid.UUID]model.Book)}
}

// sortBooks will just sort a slice of books in place by title
func sortBooks(books []model.Book) {
	sort.Slice(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})
}

// GetBooks returns a sorted slice of all of the books in the
// library
func (l *Library) GetBooks() []model.Book {
	l.Lock()
	defer l.Unlock()

	// index through the map and get a slice of all of the books
	books := make([]model.Book, len(l.Books))
	i := 0
	for _, book := range l.Books {
		books[i] = book
		i++
	}

	// sort the books before returning them
	sortBooks(books)

	return books
}

// AddBook is a thread safe putter for a key in the library's
// Book map
func (l *Library) AddBook(book model.Book) error {
	l.Lock()
	defer l.Unlock()

	l.Books[book.ID] = book

	return nil
}

// GetBookByID is just a thread safe getter for a key in the library's
// Book map
func (l *Library) GetBookByID(id uuid.UUID) (model.Book, error) {
	l.Lock()
	defer l.Unlock()

	book, found := l.Books[id]
	if !found {
		return book, ErrNoBookWithThatID
	}

	return book, nil
}

// ModifyBook will take an a book and update the given book with the same
// uuid with all of the fields populated
func (l *Library) ModifyBook(newBook model.Book) error {
	l.Lock()
	defer l.Unlock()

	// first see if the book is in the map
	book, found := l.Books[newBook.ID]
	if !found {
		return ErrNoBookWithThatID
	}

	// if the book was found then go through each of the
	// book attributes, and if the newBook's field isn't the default,
	// then overwrite the book's attribute

	defaultBook := model.NewDefaultBook()

	if newBook.Title != defaultBook.Title {
		book.Title = newBook.Title
	}

	if newBook.Author != defaultBook.Author {
		book.Author = newBook.Author
	}

	if newBook.Publisher != defaultBook.Publisher {
		book.Publisher = newBook.Publisher
	}

	if newBook.PublishDate != defaultBook.PublishDate {
		book.PublishDate = newBook.PublishDate
	}

	if newBook.Rating != defaultBook.Rating {
		book.Rating = newBook.Rating
	}

	if newBook.Status != defaultBook.Status {
		book.Status = newBook.Status
	}

	// overwrite the book in the map with the modified book
	// to get the new parameters
	l.Books[book.ID] = book

	return nil
}

// DeleteBook will remove a book from the library if it exists
func (l *Library) DeleteBook(id uuid.UUID) error {
	l.Lock()
	defer l.Unlock()

	if _, found := l.Books[id]; !found {
		return ErrNoBookWithThatID
	}

	delete(l.Books, id)
	return nil
}
