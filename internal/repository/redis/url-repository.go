package redis

import "github.com/redis/go-redis/v9"

type URLRepository struct {
	client *redis.Client
}

func NewURLRepository(client *redis.Client) *URLRepository {
	return &URLRepository{client: client}
}
