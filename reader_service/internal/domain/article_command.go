package domain

import "time"

type CreatedArticleCommand struct {
	ID        int `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r CreatedArticleCommand)ToArticle() Article {
	return Article{
		ID:        r.ID,
		Author:    r.Author,
		Title:     r.Title,
		Body:      r.Body,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
