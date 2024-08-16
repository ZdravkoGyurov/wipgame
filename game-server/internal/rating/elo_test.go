package rating_test

import (
	"testing"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/rating"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRating(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rating Suite")
}

var _ = Describe("Rating", func() {
	It("Player 1 wins", func() {
		newRating1, newRating2 := rating.CalculateNew(1200, 1000, rating.OutcomePlayer1Win)
		Expect(newRating1).To(Equal(1208))
		Expect(newRating2).To(Equal(992))
	})
	It("Player 2 wins", func() {
		newRating1, newRating2 := rating.CalculateNew(1500, 1500, rating.OutcomePlayer2Win)
		Expect(newRating1).To(Equal(1484))
		Expect(newRating2).To(Equal(1516))
	})
	It("Player 2 wins vs stronger player", func() {
		newRating1, newRating2 := rating.CalculateNew(2000, 1500, rating.OutcomePlayer2Win)
		Expect(newRating1).To(Equal(1970))
		Expect(newRating2).To(Equal(1530))
	})
	It("Draw with equal ratings", func() {
		newRating1, newRating2 := rating.CalculateNew(1700, 1700, rating.OutcomeDraw)
		Expect(newRating1).To(Equal(1700))
		Expect(newRating2).To(Equal(1700))
	})
})
