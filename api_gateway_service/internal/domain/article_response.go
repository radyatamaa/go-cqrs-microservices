package domain

import "time"

type ArticleResponse struct {
	ID        int `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ArticlePaginationResponse struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
	HasMore    bool `json:"has_more"`
	Articles   []*ArticleResponse `json:"articles"`
}


