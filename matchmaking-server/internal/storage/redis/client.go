package redis

import (
	"fmt"
	"log/slog"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	cfg    config.Redis
}

func NewClient(cfg config.Redis) Client {
	return Client{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
		}),
		cfg: cfg,
	}
}

func (c Client) Close() {
	if err := c.client.Close(); err != nil {
		slog.Error(fmt.Sprintf("failed to close redis client: %s", err))
	}
}
