package matchmaking

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
	"golang.org/x/exp/rand"
)

type match struct {
	player1 *types.Player
	player2 *types.Player
}

func (s *Server) execWithTimeout(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, s.cfg.Timeout)
	defer cancel()
	s.exec(timeoutCtx)
}

func (s *Server) exec(ctx context.Context) {
	var matches = []match{}
	var matchedPlayers = types.Set[string]{}

	playerQueue, err := s.redisClient.RetrieveQueue(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to retrieve player queue: %s", err))

		return
	}

	for idx := range playerQueue {
		player := playerQueue[idx]
		if _, hasOpponent := matchedPlayers[player.ID]; hasOpponent {
			continue
		}

		opponent1, opponent2 := getOpponents(playerQueue, idx, matchedPlayers)
		match, found := s.findValidMatch(player, opponent1, opponent2, matchedPlayers)
		if found {
			matches = append(matches, match)
		}
	}

	for _, match := range matches {
		if err := s.redisClient.MatchPlayers(ctx, match.player1.ID, match.player2.ID); err != nil {
			slog.Error(fmt.Sprintf("failed to match player '%s' with player '%s': %s",
				match.player1.ID, match.player2.ID, err))

			continue
		}
	}
}

func (s *Server) ratingRange(timeInQueue time.Time) float64 {
	secondsInQueue := time.Since(timeInQueue).Seconds()
	if secondsInQueue <= s.cfg.BaseRatingRangeDuration {
		return s.cfg.BaseRatingRange
	}

	extraSecondsInQueue := math.Round(secondsInQueue - s.cfg.BaseRatingRangeDuration)
	ratingRange := extraSecondsInQueue / s.cfg.RatingRangeIncrementInterval * s.cfg.RatingRangeMultiplier

	return ratingRange
}

func (s *Server) findValidMatch(player, opponent1, opponent2 *types.Player, matchedPlayers types.Set[string]) (match, bool) {
	isOpponent1ValidMatch := s.isValidMatch(player, opponent1)
	isOpponent2ValidMatch := s.isValidMatch(player, opponent2)

	if isOpponent1ValidMatch && isOpponent2ValidMatch {
		if rand.Intn(2) == 0 {
			return matchPlayers(player, opponent1, matchedPlayers), true
		} else {
			return matchPlayers(player, opponent2, matchedPlayers), true
		}
	}

	if isOpponent1ValidMatch {
		return matchPlayers(player, opponent1, matchedPlayers), true
	}

	if isOpponent2ValidMatch {
		return matchPlayers(player, opponent2, matchedPlayers), true
	}

	return match{}, false
}

func (s *Server) isValidMatch(player1, player2 *types.Player) bool {
	return math.Abs(float64(player1.Rating)-float64(player2.Rating)) <= s.ratingRange(player1.QueuedAt)
}

func matchPlayers(player1, player2 *types.Player, matchedPlayers types.Set[string]) match {
	matchedPlayers.Add(player1.ID)
	matchedPlayers.Add(player2.ID)

	return match{player1: player1, player2: player2}
}

func getOpponents(queue []*types.Player, idx int, matchedPlayers types.Set[string]) (*types.Player, *types.Player) {
	potentialOponent1 := &types.Player{Rating: math.MinInt}
	if idx-1 >= 0 && !matchedPlayers.Has(queue[idx-1].ID) {
		potentialOponent1 = queue[idx-1]
	}

	potentialOponent2 := &types.Player{Rating: math.MaxInt}
	if idx+1 <= len(queue)-1 && !matchedPlayers.Has(queue[idx+1].ID) {
		potentialOponent2 = queue[idx+1]
	}

	return potentialOponent1, potentialOponent2
}
