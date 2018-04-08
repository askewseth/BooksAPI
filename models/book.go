package model

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	// ErrInvalidStatus is returned whenever someone tried to create or modify a book to have an
	// invalid status
	ErrInvalidStatus = errors.New("The status must either be CheckedIn or CheckedOut")

	// ErrInvalidRating is returned whenever someone tried to create or modify a book to have an
	// invalid rating
	ErrInvalidRating = errors.New("The rating must be 1-3")
)

// Status is an enum that will cover the two different status for books
type Status uint8

// this const block holds the Status enum values
const (
	CheckedIn Status = iota
	CheckedOut
)

// NullUInt8 is the null value that will be used for uint8 fields
// since uint8 doesn't support -1, the null value is 255
const NullUInt8 = 255

// Book is the struct that holds all of the attributes for a book
type Book struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title,omitempty"`
	Author      string     `json:"author,omitempty"`
	Publisher   string     `json:"publisher,omitempty"`
	PublishDate *time.Time `json:"publish_date,omitempty"`
	Rating      uint8      `json:"rating,omitempty"`
	Status      Status     `json:"status,omitempty"`
}

// NewBook returns an initalized Book struct
// with a uuid and a CheckedIn status
func NewBook() Book {
	id, _ := uuid.NewV4()
	return Book{ID: id, Status: CheckedIn}
}

// NewDefaultBook returns a book with all of the
// fields with a negative value so that manager.ModifyBook
// can tell whether or not a field was given by a user
func NewDefaultBook() Book {
	return Book{
		Title:       "-1",
		Author:      "-1",
		Publisher:   "-1",
		PublishDate: nil,
		Rating:      NullUInt8,
		Status:      Status(NullUInt8),
	}
}

// Validate will return an error if any of the feilds of the book are outside of what they should be
func (b Book) Validate() error {
	// check if the rating is between 1 and 3
	if b.Rating < 1 || b.Rating > 3 {
		return ErrInvalidRating
	}

	// check if the status is one of the two valid statuses
	status := int(b.Status)
	if status < 1 || status > 2 {
		return ErrInvalidStatus
	}

	return nil
}
