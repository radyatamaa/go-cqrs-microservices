package domain

type CreateArticleRequest struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

func (r CreateArticleRequest)ToCreateArticleCommand() CreateArticleCommand {
	return CreateArticleCommand{
		ID:     0,
		Author: r.Author,
		Title:  r.Title,
		Body:   r.Body,
	}
}
