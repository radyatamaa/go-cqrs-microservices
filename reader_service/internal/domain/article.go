package domain

import (
	"context"
	"time"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/utils"
	readerService "github.com/radyatamaa/go-cqrs-microservices/reader_service/proto/article_reader"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Article struct {
	ID        int       `json:"id" bson:"_id,omitempty"`
	Author    string    `json:"author,omitempty" bson:"author,omitempty" validate:"required,min=3,max=250"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty" validate:"required,min=3,max=250"`
	Body      string    `json:"body,omitempty" bson:"body,omitempty" validate:"required,min=3,max=250"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// ArticlesList articles list response with pagination
type ArticlesList struct {
	TotalCount int64      `json:"totalCount" bson:"totalCount"`
	TotalPages int64      `json:"totalPages" bson:"totalPages"`
	Page       int64      `json:"page" bson:"page"`
	Size       int64      `json:"size" bson:"size"`
	HasMore    bool       `json:"hasMore" bson:"hasMore"`
	Articles   []*Article `json:"articles" bson:"articles"`
}

// TableName name of table
func (r *Article) TableName() string {
	return "articles"
}

// ArticleUseCase UseCase Interface
type ArticleUseCase interface {
	CreateArticle(c context.Context, command CreatedArticleCommand) error
	SearchArticle(c context.Context, query SearchArticleQuery) (*ArticlesList, error)
}

// MongoArticleRepository Repository Interface
type MongoArticleRepository interface {
	Create(ctx context.Context, article Article) (*Article, error)
	Update(ctx context.Context, article Article) (*Article, error)
	Delete(ctx context.Context, id int) error

	GetById(ctx context.Context, id int) (*Article, error)
	Search(ctx context.Context, search string, author string, pagination *utils.Pagination) (*ArticlesList, error)
}

// RedisArticleRepository Repository Interface
type RedisArticleRepository interface {
	Put(ctx context.Context, key string, article *Article)
	Get(ctx context.Context, key string) (*Article, error)
	Del(ctx context.Context, key string)
	DelAll(ctx context.Context)
}

// Mapper
func NewArticleListWithPagination(articles []*Article, count int64, pagination *utils.Pagination) *ArticlesList {
	return &ArticlesList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Articles:   articles,
	}
}

func ArticleToGrpcMessage(article *Article) *readerService.Article {
	return &readerService.Article{
		ID:        int32(article.ID),
		Author:    article.Author,
		Title:     article.Title,
		Body:      article.Body,
		CreatedAt: timestamppb.New(article.CreatedAt),
		UpdatedAt: timestamppb.New(article.UpdatedAt),
	}
}

func ArticleListToGrpc(articles *ArticlesList) *readerService.SearchRes {
	list := make([]*readerService.Article, 0, len(articles.Articles))
	for _, product := range articles.Articles {
		list = append(list, ArticleToGrpcMessage(product))
	}

	return &readerService.SearchRes{
		TotalCount: articles.TotalCount,
		TotalPages: articles.TotalPages,
		Page:       articles.Page,
		Size:       articles.Size,
		HasMore:    articles.HasMore,
		Articles:   list,
	}
}
