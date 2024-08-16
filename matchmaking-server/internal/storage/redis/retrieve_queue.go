package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
)

const retrieveQueueLuaScript = `
	local ids = redis.call('ZRANGE', KEYS[1], 0, -1)
	local players = {}

	for i, id in ipairs(ids) do
		local player = redis.call('HGETALL', KEYS[2] .. ':' .. id)
		table.insert(players, player)
	end

	return players
`

func (c Client) RetrieveQueue(ctx context.Context) ([]*types.Player, error) {
	args := []string{c.cfg.SortedSetName, c.cfg.HashSetName}
	playersRaw, err := c.client.Eval(ctx, retrieveQueueLuaScript, args).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal players raw: %w", err)
	}

	playersRawValues, ok := playersRaw.([]any)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshal players raw values")
	}

	players, err := unmarshalPlayers(playersRawValues)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal players: %w", err)
	}

	return players, nil
}

func unmarshalPlayers(playersRawValues []any) ([]*types.Player, error) {
	players := []*types.Player{}

	for _, playerRaw := range playersRawValues {
		playerRawValues, ok := playerRaw.([]any)
		if !ok {
			return nil, fmt.Errorf("failed to unmarshal player raw values")
		}

		player, err := unmarshalPlayer(playerRawValues)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal player: %w", err)
		}

		players = append(players, player)
	}

	return players, nil
}

const (
	idFieldIdx       = 1
	ratingFieldIdx   = 3
	queuedAtFieldIdx = 5
)

func unmarshalPlayer(rawValues []any) (*types.Player, error) {
	rating, err := strconv.Atoi(rawValues[ratingFieldIdx].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to convert rating from string to int: %w", err)
	}

	queuedAt, err := time.Parse(time.Layout, rawValues[queuedAtFieldIdx].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to convert queuedAt from string to time.Time: %w", err)
	}

	return &types.Player{
		ID:       rawValues[idFieldIdx].(string),
		Rating:   rating,
		QueuedAt: queuedAt,
	}, nil
}
