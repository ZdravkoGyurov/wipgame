package redis

import (
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewClient(cfg config.Redis) Client {
	return Client{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
		}),
	}
}
