package repository

import (
	"context"
	"strings"

	"github.com/radyatamaa/go-cqrs-microservices/write_service/internal/domain"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/database/paginator"

	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"

	"gorm.io/gorm"
)

type pgArticleRepository struct {
	zapLogger zaplogger.Logger
	db        *gorm.DB
}

func NewPgArticleRepository(db *gorm.DB, zapLogger zaplogger.Logger) domain.PgArticleRepository {
	return &pgArticleRepository{
		db:        db,
		zapLogger: zapLogger,
	}
}

func (c pgArticleRepository) DB() *gorm.DB {
	return c.db
}

func (c pgArticleRepository) FetchWithFilter(ctx context.Context, limit int, offset int, order string, fields, associate []string, model interface{}, args ...interface{}) error {
	p := paginator.NewPaginator(c.db, offset, limit, model)

	return p.FindWithFilter(ctx, order, fields, associate, args...).Select(strings.Join(fields, ",")).Error
}

func (c pgArticleRepository) SingleWithFilter(ctx context.Context, fields, associate []string, model interface{}, args ...interface{}) error {

	db := c.db.WithContext(ctx)

	if len(fields) > 0 {
		db = db.Select(strings.Join(fields, ","))
	}
	if len(associate) > 0 {
		for _, v := range associate {
			db.Joins(v)
		}
	}

	if err := db.First(model, args...).Error; err != nil {
		return err
	}

	return nil
}

func (c pgArticleRepository) Update(ctx context.Context, data domain.Article) error {

	err := c.db.WithContext(ctx).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (c pgArticleRepository) UpdateSelectedField(ctx context.Context, field []string, values map[string]interface{}, id int) error {

	return c.db.WithContext(ctx).Table("customer_limit").Select(field).Where("id =?", id).Updates(values).Error
}

func (c pgArticleRepository) Store(ctx context.Context, data domain.Article) (domain.Article, error) {

	err := c.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (c pgArticleRepository) Delete(ctx context.Context, id int) (int, error) {

	err := c.db.WithContext(ctx).Exec("delete from customer_limit where id =?", id).Error
	if err != nil {
		return id, err
	}
	return id, nil
}

func (c pgArticleRepository) SoftDelete(ctx context.Context, id int) (int, error) {
	var data domain.Article

	err := c.db.WithContext(ctx).Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return id, err
	}
	return id, nil
}

func (c pgArticleRepository) UpdateSelectedFieldWithTx(ctx context.Context, tx *gorm.DB, field []string, values map[string]interface{}, id int) error {

	return tx.WithContext(ctx).Table("customer_limit").Select(field).Where("id =?", id).Updates(values).Error
}

func (c pgArticleRepository) StoreWithTx(ctx context.Context, tx *gorm.DB, data domain.Article) (int, error) {

	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data.ID, err
	}
	return data.ID, nil
}

func (c pgArticleRepository) FetchWithFilterAndPagination(ctx context.Context, limit int, offset int, order string, fields, associate []string, model interface{}, args ...interface{}) (*paginator.Paginator, error) {

	p := paginator.NewPaginator(c.db, offset, limit, model)
	if err := p.FindWithFilter(ctx, order, fields, associate, args...).Select(strings.Join(fields, ",")).Error; err != nil {
		return p, err
	}
	return p, nil
}
