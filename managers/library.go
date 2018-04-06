package managers

import (
	"fmt"
	"sync"

	"github.com/askewseth/kubernetes/models"
)

var (
	myLibrary Library

	once = sync.Once{}
)

type Library struct {
	sync.Mutex `json:"-"`
	Books      map[string]model.Book `json:"books"`
}

func GetLibrary() Library {
	once.Do(func() {
		myLibrary = newLibrary()
	})

	return myLibrary
}

func newLibrary() Library {
	return Library{Books: make(map[string]model.Book)}
}

func (l *Library) GetBooks() []model.Book {
	l.Lock()
	defer l.Unlock()

	books := make([]model.Book, len(l.Books))

	i := 0
	for _, book := range books {
		books[i] = book
		i++
	}

	fmt.Printf("Got books: %+v\n", books)

	return books
}

func (l *Library) AddBook(book model.Book) error {
	l.Lock()
	defer l.Unlock()

	fmt.Printf("Adding book: %+v\n", book)

	l.Books[book.Title] = book

	fmt.Printf("From Map: %+v", l.Books[book.Title])
	return nil
}
