package redis

import "context"

// MatchPlayers
// 0. Lua script
// 1. Check if both players are still in queue
// 2. Remove them from the queue
// 3. Add them in new match queue as a pair
// 4. Game server will create a game for them
func (c Client) MatchPlayers(ctx context.Context, playerID1, playerID2 string) error {
	return nil // TODO
}
