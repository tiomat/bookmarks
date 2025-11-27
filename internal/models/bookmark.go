package models

import "time"

type Bookmark struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"` // HTML content
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
