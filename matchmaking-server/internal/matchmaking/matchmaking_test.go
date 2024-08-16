package matchmaking_test

import (
	"testing"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/matchmaking"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRating(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matchmaking Suite")
}

var _ = Describe("Matchmaking", func() {
	It("", func() {
		matchmaking.EnqueuePlayer(types.Player{ID: "1", Rating: 1200})
		matchmaking.EnqueuePlayer(types.Player{ID: "1", Rating: 1260})
		matchmaking.Run()
	})
})
