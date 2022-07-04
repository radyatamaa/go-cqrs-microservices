package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/radyatamaa/go-cqrs-microservices/pkg/zaplogger"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/config"
	"github.com/radyatamaa/go-cqrs-microservices/reader_service/internal/domain"
)

const (
	redisProductPrefixKey = "reader:product"
)

type redisRepository struct {
	log         zaplogger.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisRepository(log zaplogger.Logger, cfg *config.Config, redisClient redis.UniversalClient) domain.RedisArticleRepository {
	return &redisRepository{log: log, cfg: cfg, redisClient: redisClient}
}

func (r *redisRepository) Put(ctx context.Context, key string, article *domain.Article) {
	productBytes, err := json.Marshal(article)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisArticlePrefixKey(), key, productBytes).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisArticlePrefixKey(), key)
}

func (r *redisRepository) Get(ctx context.Context, key string) (*domain.Article, error) {
	articleBytes, err := r.redisClient.HGet(ctx, r.getRedisArticlePrefixKey(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}

	var article domain.Article
	if err := json.Unmarshal(articleBytes, &article); err != nil {
		return nil, err
	}

	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisArticlePrefixKey(), key)
	return &article, nil
}

func (r *redisRepository) Del(ctx context.Context, key string) {
	if err := r.redisClient.HDel(ctx, r.getRedisArticlePrefixKey(), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisArticlePrefixKey(), key)
}

func (r *redisRepository) DelAll(ctx context.Context) {
	if err := r.redisClient.Del(ctx, r.getRedisArticlePrefixKey()).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisArticlePrefixKey())
}

func (r *redisRepository) getRedisArticlePrefixKey() string {
	if r.cfg.ServiceSettings.RedisArticlePrefixKey != "" {
		return r.cfg.ServiceSettings.RedisArticlePrefixKey
	}

	return redisProductPrefixKey
}
