package domain

import "time"

type CreateArticleCommand struct {
	ID        int `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

type CreatedArticleCommand struct {
	ID        int `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r CreateArticleCommand)ToArticle() Article {
	return Article{
		ID:        0,
		Author:    r.Author,
		Title:     r.Title,
		Body:      r.Body,
	}
}
