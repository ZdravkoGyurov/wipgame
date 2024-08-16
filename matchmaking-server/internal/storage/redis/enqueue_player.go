package redis

import (
	"context"
	"fmt"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"

	"github.com/redis/go-redis/v9"
)

func (c Client) EnqueuePlayer(ctx context.Context, player types.Player) error {
	return c.execInTransaction(ctx, func(pipe redis.Pipeliner) error {
		key := fmt.Sprintf("%s:%s", c.cfg.HashSetName, player.ID)
		hmsetCmd := pipe.HMSet(ctx, key, player)
		if hmsetCmd.Err() != nil {
			return fmt.Errorf("failed to insert player with id '%s' in the hash set: %w", player.ID, hmsetCmd.Err())
		}

		zaddCmd := pipe.ZAdd(ctx, c.cfg.SortedSetName, redis.Z{
			Score:  float64(player.Rating),
			Member: player.ID,
		})
		if zaddCmd.Err() != nil {
			return fmt.Errorf("failed to insert player with id '%s' in the ordered set: %w", player.ID, zaddCmd.Err())
		}

		return nil
	})
}
