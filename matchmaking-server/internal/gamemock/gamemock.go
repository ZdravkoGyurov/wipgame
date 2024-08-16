package gamemock

import (
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/rating"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
	"golang.org/x/exp/rand"
)

func Run(player1, player2 types.Player) rating.Outcome {
	winner := rand.Intn(3)
	if winner == 0 {
		return rating.OutcomePlayer1Win
	}
	if winner == 1 {
		return rating.OutcomePlayer2Win
	}
	return rating.OutcomeDraw
}
