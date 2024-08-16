package types

import "time"

type Player struct {
	ID string `json:"id" redis:"id"`
	// Rating is used for sorting the Players in the queue
	Rating int `json:"rating" redis:"rating"`
	// QueuedAt is used for calculating the rating range during the matchmaking search
	QueuedAt time.Time `json:"queuedAt" redis:"queuedAt"`
}
