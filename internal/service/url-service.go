package service

import (
	"design-url-shortner/internal/repository/postgres"
	"design-url-shortner/internal/repository/redis"
)

type URLService struct {
	postgresRepo *postgres.URLRepository
	redisRepo    *redis.URLRepository
	baseURL      string
}

func NewURLService(postgresRepo *postgres.URLRepository, redisRepo *redis.URLRepository, baseURL string) *URLService {
	return &URLService{
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
		baseURL:      baseURL,
	}

}
