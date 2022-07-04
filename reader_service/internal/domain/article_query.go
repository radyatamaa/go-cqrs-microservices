package domain

import "github.com/radyatamaa/go-cqrs-microservices/pkg/utils"

type SearchArticleQuery struct {
	Author     string            `json:"author"`
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

func NewSearchArticleQuery(text string, author string, pagination *utils.Pagination) SearchArticleQuery {
	return SearchArticleQuery{Text: text, Author: author, Pagination: pagination}
}
