package matchmaking

import (
	"context"
	"time"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/config"
	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/storage/redis"
)

type Server struct {
	cfg         config.Matchmaking
	redisClient redis.Client
	done        chan struct{}
}

func NewServer(cfg config.Matchmaking, redisClient redis.Client) Server {
	return Server{
		cfg:         cfg,
		redisClient: redisClient,
		done:        make(chan struct{}),
	}
}

func (s Server) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.Interval)

	for {
		select {
		case <-s.done:
			ticker.Stop()
			return
		case <-ticker.C:
			s.execWithTimeout(ctx)
		}
	}
}

func (s Server) Stop() {
	s.done <- struct{}{}
	close(s.done)
}
