package main

import (
	"context"
	"log"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/storage/redis"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(config.Redis{
		Address:  "127.0.0.1:6379",
		Password: "2xT3f4hNnP",
	})

	err := client.EnqueuePlayer(ctx, types.Player{
		ID:            "1",
		Rating:        1700,
		RatingRange:   50,
		OpponentFound: false,
	})
	if err != nil {
		log.Panic(err)
	}

	err = client.EnqueuePlayer(ctx, types.Player{
		ID:            "2",
		Rating:        1600,
		RatingRange:   50,
		OpponentFound: false,
	})
	if err != nil {
		log.Panic(err)
	}

	err = client.EnqueuePlayer(ctx, types.Player{
		ID:            "3",
		Rating:        1500,
		RatingRange:   50,
		OpponentFound: false,
	})
	if err != nil {
		log.Panic(err)
	}

	players, err := client.RetrieveQueue(ctx)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("players: %v\n", players)

	err = client.DequeuePlayer(ctx, "1")
	if err != nil {
		log.Panic(err)
	}

	players, err = client.RetrieveQueue(ctx)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("players: %v\n", players)
}
