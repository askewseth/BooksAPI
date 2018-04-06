package model

import "time"

// - Title
// - Author
// - Publisher
// - Publish Date
// - Rating (1-3)
// - Status (CheckedIn, CheckedOut)

type Status uint8

const (
	CheckedIn Status = iota
	CheckedOut
)

type Book struct {
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
	Rating      uint8     `json:"rating"`
	Status      Status    `json:"status"`
}
