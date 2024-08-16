package types

type Player struct {
	ID            string `json:"id" redis:"id"`
	Rating        int    `json:"rating" redis:"rating"`
	RatingRange   int    `json:"ratingRange" redis:"ratingRange"`
	OpponentFound bool   `json:"opponentFound" redis:"opponentFound"`
}
