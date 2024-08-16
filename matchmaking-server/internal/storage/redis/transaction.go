package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (c Client) execInTransaction(ctx context.Context, transaction func(p redis.Pipeliner) error) error {
	pipe := c.client.TxPipeline()

	if err := transaction(pipe); err != nil {
		return err
	}

	_, err := pipe.Exec(ctx)
	return err
}
