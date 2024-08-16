package redis

import (
	"context"
	"fmt"
	"strconv"

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

func (c Client) RetrieveQueue(ctx context.Context) ([]types.Player, error) {
	playersRaw, err := c.client.Eval(ctx, retrieveQueueLuaScript, []string{sortedSetName, hashSetName}).Result()
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

func unmarshalPlayers(playersRawValues []any) ([]types.Player, error) {
	players := []types.Player{}

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
	idFieldIdx            = 1
	ratingFieldIdx        = 3
	ratingRangeFieldIdx   = 5
	opponentFoundFieldIdx = 7
)

func unmarshalPlayer(rawValues []any) (types.Player, error) {
	rating, err := strconv.Atoi(rawValues[ratingFieldIdx].(string))
	if err != nil {
		return types.Player{}, fmt.Errorf("failed to convert rating from string to int: %w", err)
	}

	ratingRange, err := strconv.Atoi(rawValues[ratingRangeFieldIdx].(string))
	if err != nil {
		return types.Player{}, fmt.Errorf("failed to convert ratingRange from string to int: %w", err)
	}

	opponentFound, err := strconv.ParseBool(rawValues[opponentFoundFieldIdx].(string))
	if err != nil {
		return types.Player{}, fmt.Errorf("failed to convert opponentFound from string to bool: %w", err)
	}

	return types.Player{
		ID:            rawValues[idFieldIdx].(string),
		Rating:        rating,
		RatingRange:   ratingRange,
		OpponentFound: opponentFound,
	}, nil
}
