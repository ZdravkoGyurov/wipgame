package rating

import (
	"math"
)

type Outcome float64

const (
	OutcomePlayer1Win = 1
	OutcomePlayer2Win = 0
	OutcomeDraw       = 0.5
)

func CalculateNew(rating1, rating2 int, outcome Outcome) (int, int) {
	rating1F := float64(rating1)
	rating2F := float64(rating2)

	winProbability1 := probability(rating2F, rating1F)
	winProbability2 := 1.0 - winProbability1

	kFactor1, kFactor2 := kFactors(rating1F, rating2F, outcome)
	newRating1 := rating1F + kFactor1*(float64(outcome)-winProbability1)
	newRating2 := rating2F + kFactor2*(1-float64(outcome)-winProbability2)

	return int(math.Round(newRating1)), int(math.Round(newRating2))
}

const (
	scaleFactor  float64 = 400.0
	exponentBase float64 = 10.0
)

func probability(rating1, rating2 float64) float64 {
	power := (rating1 - rating2) / scaleFactor
	return 1.0 / (1 + math.Pow(exponentBase, power))
}

// kFactors
// win - own rating / loss - enemy rating.
func kFactors(rating1, rating2 float64, outcome Outcome) (float64, float64) {
	switch outcome {
	case OutcomePlayer1Win:
		return kFactor(rating1), kFactor(rating1)
	case OutcomePlayer2Win:
		return kFactor(rating2), kFactor(rating2)
	case OutcomeDraw:
		return kFactor(rating2), kFactor(rating1)
	}

	panic("invalid outcome provided")
}

func kFactor(rating float64) float64 {
	if rating >= 2400 {
		return 16.0
	}

	if rating >= 2100 {
		return 24.0
	}

	return 32.0
}
