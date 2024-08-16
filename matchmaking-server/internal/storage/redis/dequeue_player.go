package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (c Client) DequeuePlayer(ctx context.Context, playerID string) error {
	return c.execInTransaction(ctx, func(pipe redis.Pipeliner) error {
		key := fmt.Sprintf("%s:%s", hashSetName, playerID)
		delCmd := pipe.Del(ctx, key)
		if delCmd.Err() != nil {
			return fmt.Errorf("failed to remove player with id '%s' from the hash set: %w", playerID, delCmd.Err())
		}

		zremCmd := pipe.ZRem(ctx, sortedSetName, playerID)
		if zremCmd.Err() != nil {
			return fmt.Errorf("failed to remove player with id '%s' from the ordered set: %w", playerID, zremCmd.Err())
		}

		return nil
	})
}
