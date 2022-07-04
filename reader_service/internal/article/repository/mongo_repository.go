package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/utils"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoArticleRepository struct {
	log zaplogger.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewMongoArticleRepository(log zaplogger.Logger, cfg *config.Config, db *mongo.Client) domain.MongoArticleRepository {
	return &mongoArticleRepository{log: log, cfg: cfg, db: db}
}

func (p *mongoArticleRepository) Search(ctx context.Context, search string, author string, pagination *utils.Pagination) (*domain.ArticlesList, error) {
	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Articles)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
			bson.D{{Key: "body", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
		}},
	}

	if author != "" {
		filter = bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "title", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
				bson.D{{Key: "body", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
			}},
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "author", Value: author}},
			}},
		}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &domain.ArticlesList{Articles: make([]*domain.Article, 0)}, nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.D{{"createdAt", -1}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errcheck

	articles := make([]*domain.Article, 0, pagination.GetSize())

	for cursor.Next(ctx) {
		var article domain.Article
		if err := cursor.Decode(&article); err != nil {
			return nil, errors.Wrap(err, "Find")
		}
		articles = append(articles, &article)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return domain.NewArticleListWithPagination(articles, count, pagination), nil
}

func (p *mongoArticleRepository) Create(ctx context.Context, article domain.Article) (*domain.Article, error) {

	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Articles)

	_, err := collection.InsertOne(ctx, article, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	return &article, nil
}

func (p *mongoArticleRepository) Update(ctx context.Context, article domain.Article) (*domain.Article, error) {

	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Articles)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated domain.Article
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": article.ID}, bson.M{"$set": article}, ops).Decode(&updated); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &updated, nil
}

func (p *mongoArticleRepository) Delete(ctx context.Context, id int) error {

	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Articles)

	return collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Err()
}

func (p *mongoArticleRepository) GetById(ctx context.Context, id int) (*domain.Article, error) {

	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Articles)

	var article domain.Article
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &article, nil
}
