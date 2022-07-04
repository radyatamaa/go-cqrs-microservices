package domain

import (
	"context"
	"time"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/database/paginator"
	"gorm.io/gorm"
)

type Article struct {
	ID        int            `gorm:"column:id;primarykey;autoIncrement:true"`
	Author    string         `gorm:"type:text;column:author"`
	Title     string         `gorm:"type:text;column:title"`
	Body      string         `gorm:"type:text;column:body"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName name of table
func (r *Article) TableName() string {
	return "articles"
}

// ArticleUseCase UseCase Interface
type ArticleUseCase interface {
	CreateArticle(c context.Context, command CreateArticleCommand) error
}

// PgArticleRepository Repository Interface
type PgArticleRepository interface {
	SingleWithFilter(ctx context.Context, fields, associate []string, model interface{}, args ...interface{}) error
	FetchWithFilter(ctx context.Context, limit int, offset int, order string, fields, associate []string, model interface{}, args ...interface{}) error
	Update(ctx context.Context, data Article) error
	UpdateSelectedField(ctx context.Context, field []string, values map[string]interface{}, id int) error
	UpdateSelectedFieldWithTx(ctx context.Context, tx *gorm.DB, field []string, values map[string]interface{}, id int) error
	Store(ctx context.Context, data Article) (Article, error)
	StoreWithTx(ctx context.Context, tx *gorm.DB, data Article) (int, error)
	Delete(ctx context.Context, id int) (int, error)
	SoftDelete(ctx context.Context, id int) (int, error)
	DB() *gorm.DB
	FetchWithFilterAndPagination(ctx context.Context, limit int, offset int, order string, fields, associate []string, model interface{}, args ...interface{}) (*paginator.Paginator, error)
}

// MessagingArticleRepository Repository Interface
type MessagingArticleRepository interface {
	PushMessageInsertArticle(ctx context.Context, msg []byte) error
}
