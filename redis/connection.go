package redis

import "github.com/go-redis/redis"

func NewConnection(host string, password string, db int) *redis.Client  {
	return redis.NewClient(&redis.Options{
			Addr: host,
			Password: password,
			DB: db,
	})
}
