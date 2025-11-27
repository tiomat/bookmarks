package archiver

import (
	"bookmarks/internal/models"
	"time"

	"github.com/go-shiori/go-readability"
)

type Archiver struct{}

func New() *Archiver {
	return &Archiver{}
}

func (a *Archiver) Archive(url string) (*models.Bookmark, error) {
	var article readability.Article
	var err error

	// Retry logic: 3 attempts
	for i := 0; i < 3; i++ {
		article, err = readability.FromURL(url, 30*time.Second)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return &models.Bookmark{
		URL:     url,
		Title:   article.Title,
		Excerpt: article.Excerpt,
		Content: article.Content,
	}, nil
}
