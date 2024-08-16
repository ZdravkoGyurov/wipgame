package matchmaking

import (
	"fmt"
	"math"
	"time"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/gamemock"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/rating"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
	"golang.org/x/exp/rand"
)

type match struct {
	player1 types.Player
	player2 types.Player
}

var matches = []match{}

const (
	interval             = 2 * time.Second
	defaultRatingRange   = 50
	ratingRangeIncrement = 10
)

func EnqueuePlayer(p types.Player) {
	queue = insertPlayer(queue, p)
}

func Run() {
	for {
		snapshotQueue := make([]types.Player, len(queue))
		copy(snapshotQueue, queue)
		fmt.Printf("queue snapshot: %v\n", snapshotQueue)
		fmt.Printf("ratingRanges: %v\n", ratingRanges)

		for idx := range snapshotQueue {
			player := snapshotQueue[idx]
			if player.OpponentFound {
				continue
			}
			potentialOponent1 := types.Player{Rating: math.MinInt}
			if idx-1 >= 0 && !snapshotQueue[idx-1].OpponentFound {
				potentialOponent1 = snapshotQueue[idx-1]
			}
			potentialOponent2 := types.Player{Rating: math.MaxInt}
			if idx+1 <= len(snapshotQueue)-1 && !snapshotQueue[idx+1].OpponentFound {
				potentialOponent2 = snapshotQueue[idx+1]
			}

			playerRatingF := float64(player.Rating)
			potentialOponent1RatingF := float64(potentialOponent1.Rating)
			potentialOponent2RatingF := float64(potentialOponent2.Rating)
			ratingRange := getPlayerRatingRange(player.ID)
			if math.Abs(playerRatingF-potentialOponent1RatingF) <= float64(ratingRange) && math.Abs(playerRatingF-potentialOponent2RatingF) <= float64(ratingRange) {
				chosenOpponent := rand.Intn(2)
				if chosenOpponent == 0 {
					player.OpponentFound = true
					snapshotQueue[idx] = player
					potentialOponent1.OpponentFound = true
					snapshotQueue[idx-1] = potentialOponent1
					matches = append(matches, match{player1: player, player2: potentialOponent1})

					continue
				}
				if chosenOpponent == 1 {
					player.OpponentFound = true
					snapshotQueue[idx] = player
					potentialOponent2.OpponentFound = true
					snapshotQueue[idx+1] = potentialOponent2
					matches = append(matches, match{player1: player, player2: potentialOponent2})

					continue

				}
			}
			if math.Abs(playerRatingF-potentialOponent1RatingF) <= float64(ratingRange) {
				player.OpponentFound = true
				snapshotQueue[idx] = player
				potentialOponent1.OpponentFound = true
				snapshotQueue[idx-1] = potentialOponent1
				matches = append(matches, match{player1: player, player2: potentialOponent1})

				continue
			}
			if math.Abs(playerRatingF-potentialOponent2RatingF) <= float64(ratingRange) {
				player.OpponentFound = true
				snapshotQueue[idx] = player
				potentialOponent2.OpponentFound = true
				snapshotQueue[idx+1] = potentialOponent2
				matches = append(matches, match{player1: player, player2: potentialOponent2})

				continue
			}

			updatePlayerRatingRange(player.ID)
		}

		for idx := range matches {
			queue = removePlayer(queue, playerIdx(queue, matches[idx].player1))
			queue = removePlayer(queue, playerIdx(queue, matches[idx].player2))
			delete(ratingRanges, matches[idx].player1.ID)
			delete(ratingRanges, matches[idx].player2.ID)
			// start match
			outcome := gamemock.Run(matches[idx].player1, matches[idx].player2)
			newRating1, newRating2 := rating.CalculateNew(matches[idx].player1.Rating, matches[idx].player2.Rating, outcome)
			fmt.Printf("newRating1(%s), newRating2(%s): %v, %v\n", matches[idx].player1.ID, matches[idx].player2.ID, newRating1, newRating2)
		}

		matches = []match{}

		time.Sleep(interval)
	}
}

func getPlayerRatingRange(playerID string) int {
	ratingRange, found := ratingRanges[playerID]
	if !found {
		return defaultRatingRange
	}

	return ratingRange + ratingRangeIncrement
}

func updatePlayerRatingRange(playerID string) {
	ratingRange, found := ratingRanges[playerID]
	if !found {
		ratingRanges[playerID] = defaultRatingRange
	} else {
		ratingRanges[playerID] = ratingRange + ratingRangeIncrement
	}
}
